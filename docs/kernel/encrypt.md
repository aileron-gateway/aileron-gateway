# Package `kernel/encrypt`

## Summary

This is the design document of kernel/encrypt package.
Encrypt and decrypt, password hash utilities are documented.

## Motivation

Encryption and decryption, password hashing are very important and sensitive to security.
They can be used from anywhere in the gateway codebase.
Making those utility functions make them secure and stable.

### Goals

- Provide encrypt/decrypt utilities.
- Provide password hash utilities.

### Non-Goals

- Support various encryption algorithms.
- Support various hashing algorithms.

## Technical Design

### Password hash

Password hash functions are hash functions that are used for passwords hashing.
kernel/encrypt package provides the following hash functions.

- [BCrypt](golang.org/x/crypto/bcrypt)
- [SCrypt](golang.org/x/crypto/scrypt)
- [PBKDF2](golang.org/x/crypto/pbkdf2)
- [Argon2i](golang.org/x/crypto/argon2)
- [Argon2id](golang.org/x/crypto/argon2)

Function signature for password hashing is defined like below.
Hashed data returned by the hash function includes salt as the prefix
which means `hashed = salt + hash(password)`.
The salt and the hashed password should be separated when decrypting.

```go
// Password hashing function.
// Hash value of the given password will be returned.
func(password []byte) (hashed []byte, err error)
```

```go
// Hash validating function.
// Hash value and the password will be compared. 
// Non-nil error will be returned if the hash and the password were not matched.
func(hashed []byte, password []byte) error
```

### Common key encryption

kernel/encrypt supports the following common key encryptions.
As shown in the table, both block cipher and stream cipher are supported.
Go standard packages [crypt/aes](https://pkg.go.dev/crypto/aes),
[crypto/des](https://pkg.go.dev/crypto/des), [crypto/rc4](https://pkg.go.dev/crypto/rc4) are used.

Initial vectors are generated with the same length with the block size or key length for RC4.
Initial vector is generated using [crypt/rand](https://pkg.go.dev/crypto/rand) package.

| Type   | Algorithm | Stream Modes  | Block Modes | Key length | Block size | iv size  |
| ------ | --------- | ------------- | ----------- | ---------- | ---------- | -------- |
| Block  | AES       | CFB, CTR, OFB | CBC, GCM    | 16, 24, 32 | 16         | 16       |
| Block  | DES       | CFB, CTR, OFB | CBC         | 8          | 8          | 8        |
| Block  | 3DES      | CFB, CTR, OFB | CBC         | 24         | 8          | 8        |
| Stream | RC4       | -             | -           | 1 - 256    | -          | len(key) |

Encryption and decryption functions have the following signature.
Input key **MUST** have the valid key length because these functions
do not use any hashing technique to adjust the key length.
Returned ciphertext includes initial vector iv as the prefix which means
`ciphertext = iv + encrypt(key, plaintext)`.

```go
// Common key encryption function.
// ciphertext = iv + encrypt(key, plaintext)
func(key []byte, plaintext []byte) (ciphertext []byte, err error)
```

```go
// Common key decryption function.
func(key []byte, ciphertext []byte) (plaintext []byte, err error)
```

When using a block cipher
Padding and un-padding are required for block cipher.
Following three padding algorithms are implemented.

- [PKCS#7 - RFC 5652](https://datatracker.ietf.org/doc/rfc5652/)
- [ISO7816](https://en.wikipedia.org/wiki/Padding_(cryptography)#ISO/IEC_7816-4)
- [ISO10126](https://en.wikipedia.org/wiki/Padding_(cryptography)#ISO_10126)

Padding and un-padding functions are defined to have the following signatures.

```go
// Padding function.
// Append padding to the given data.
func(blockSize int, data []byte) (padded []byte, err error)
```

```go
// Un-adding function.
// Remove padding from the given data.
func(blockSize int, data []byte) (unpadded []byte, err error)
```

### Public key encryption

Currently not implemented.

## Test Plan

### Unit Tests

Unit tests are implemented and passed.

- All functions and methods are covered.
- Coverage objective 98%.

### Integration Tests

Not planned.

### e2e Tests

Not planned.

### Fuzz Tests

Fuzz tests are implemented for

- password hash functions that accept `[]byte`.
- password compare functions that accept `[]byte`.
- encryption functions that accept `[]byte`.
- decryption functions that accept `[]byte`.
- padding functions that accept `[]byte`.
- unpadding functions that accept `[]byte`.

### Benchmark Tests

Benchmark tests are implemented for

- password hash functions.
- password compare functions.
- encryption functions.
- decryption functions.
- padding functions.
- unpadding functions.

### Chaos Tests

Not planned.

## Future works

- [ ] Add public key encryption.

## References

- [FIPS 140-1 Security Requirements for Cryptographic Modules](https://csrc.nist.gov/pubs/fips/140-1/upd1/final)
- [FIPS 140-2 Security Requirements for Cryptographic Modules](https://csrc.nist.gov/pubs/fips/140-2/upd2/final)
- [FIPS 140-3 Security Requirements for Cryptographic Modules](https://csrc.nist.gov/pubs/fips/140-3/final)
- [Docker Container with FIPS-140 Cryptographic Module - KrakenD](https://www.krakend.io/docs/enterprise/security/fips-140/)
- [APISIX: Be FIPS 140-2 Compliant](https://api7.ai/blog/apisix-fips-140-2)
- [About FIPS 140-2 Compliance in Kong Gateway](https://docs.konghq.com/gateway/latest/kong-enterprise/fips-support/)

## Appendix

### Decrypt with openssl

Openssl can be used to decrypt ciphertext generated by the encryption
function this package provides.

Example:

- ciphertext: `66697865642076616c75652077696c6cc9e33d47997b54533831bf3447c9db50762bbf3e064ec9c35c840659f3f93ead`
- algorithm: `AES CBC`
- key (password): `16bytes password` or `313662797465732070617373776f7264` in hex

**Step 1.**:

Separate the ciphertext into iv and data part.
The block size of the AES is 16 bytes, which is 32 bytes in Hex encoded string.

iv, data := ciphertext[:16], ciphertext[17:]

It results in

- iv: `66697865642076616c75652077696c6c`
- data: `c9e33d47997b54533831bf3447c9db50762bbf3e064ec9c35c840659f3f93ead`

**Step 2.**:

Decrypt on the terminal.
Before running openssl command, data have to be converted to binary file.
It is done using `xxd` command.
After that, the ciphertext can be decrypted to plaintext with the key (in Hex) and iv (in Hex) as shown below.
Finally, the original data `plaintext message` is obtained.

```bash
$ echo -n "c9e33d47997b54533831bf3447c9db50762bbf3e064ec9c35c840659f3f93ead" > data.txt
$ xxd -r -p data.txt > data.bin
$ openssl enc -d -aes-128-cbc -nosalt -in data.bin -out decrypted.txt -K "313662797465732070617373776f7264" -iv "66697865642076616c75652077696c6c"

$ cat decrypted.txt
plaintext message
```

### Decrypt with python

Following code shows the decryption using [pyca/cryptography](https://cryptography.io/en/latest/).
See [https://cryptography.io/en/latest/hazmat/primitives/symmetric-encryption/](https://cryptography.io/en/latest/hazmat/primitives/symmetric-encryption/).

Note that the padding bytes `\x0f` are left in the output.

```python
>>> from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
>>> from cryptography.hazmat.backends import default_backend
>>> 
>>> key = bytes.fromhex('313662797465732070617373776f7264')
>>> iv = bytes.fromhex('66697865642076616c75652077696c6c')
>>> data = bytes.fromhex('c9e33d47997b54533831bf3447c9db50762bbf3e064ec9c35c840659f3f93ead')
>>> 
>>> cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
>>> decryptor = cipher.decryptor()
>>> decryptor.update(data) + decryptor.finalize()

b'plaintext message\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f\x0f'
```
