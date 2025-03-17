# Package `kernel/mac`

## Summary

This is the design document of `kernel/mac` package.

`kernel/mac` package provides [MAC (Message authentication code)](https://en.wikipedia.org/wiki/Message_authentication_code) utilities.

## Motivation

MAC is used for API key, securing data and more.
Function for calculating MAC does not depend on any business features.
Therefore, isolating those functions from any other features make them re-usable and stable.

### Goals

- Provide mac utilities.

### Non-Goals

- Support various mac algorithms.

## Technical Design

MAC has some [variations](https://en.wikipedia.org/wiki/Message_authentication_code).
For example, HMAC, KMAC and CMAC are the kind of MAC those are [approved by the NIST](https://csrc.nist.gov/projects/message-authentication-codes).

This package provides HMAC functions by utilizing go standard packages.

### Interface

The signature to calculate a HMAC of certain data is defined like below.

- HMAC functions have the following signature.
- Input message can be nil.
- Allowed key length depends on the underlying hash algorithm.
- No panics in the function.
- The length of the returned slice depends on the underlying hash algorithms.
- HMAC functions MUST return the same result for the same input.

```go
func(message []byte, key []byte) (hash []byte)
```

Available hash algorithms are defined as `Algorithm` type.

```go
type Algorithm int
```

### HMAC algorithms

Because the hash functions have same signature in the AILERON Gateway,
HMAC functions are defined for all of them.

Because this package provides atomic operations to calculate HMAC values,
it's users' responsibility to choose appropriate hash algorithm based on a situation.

Because the HMAC is basically used for security purpose,
non-cryptographic hash algorithms are not used even they can be chosen.

| Algorithm   | Cryptographic | Byte Length | Used package                                      |
| ----------- | ------------- | ----------- | ------------------------------------------------- |
| SHA1        | Yes           | 20          | [crypto/sha1](https://pkg.go.dev/crypto/sha1)     |
| SHA224      | Yes           | 28          | [crypto/sha256](https://pkg.go.dev/crypto/sha256) |
| SHA256      | Yes           | 32          | [crypto/sha256](https://pkg.go.dev/crypto/sha256) |
| SHA384      | Yes           | 48          | [crypto/sha512](https://pkg.go.dev/crypto/sha512) |
| SHA512      | Yes           | 64          | [crypto/sha512](https://pkg.go.dev/crypto/sha512) |
| SHA512_224  | Yes           | 28          | [crypto/sha512](https://pkg.go.dev/crypto/sha512) |
| SHA512_256  | Yes           | 32          | [crypto/sha512](https://pkg.go.dev/crypto/sha512) |
| SHA3_224    | Yes           | 28          | [golang.org/x/crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3) |
| SHA3_256    | Yes           | 32          | [golang.org/x/crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3) |
| SHA3_384    | Yes           | 48          | [golang.org/x/crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3) |
| SHA3_512    | Yes           | 64          | [golang.org/x/crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3) |
| SHAKE128    | Yes           | 32          | [golang.org/x/crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3) |
| SHAKE256    | Yes           | 64          | [golang.org/x/crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3) |
| MD5         | Yes           | 16          | [crypto/md5](https://pkg.go.dev/crypto/md5)     |
| FNV1_32     | No            | 4           | [hash/fnv](https://pkg.go.dev/hash/fnv)         |
| FNV1a_32    | No            | 4           | [hash/fnv](https://pkg.go.dev/hash/fnv)         |
| FNV1_64     | No            | 8           | [hash/fnv](https://pkg.go.dev/hash/fnv)         |
| FNV1a_64    | No            | 8           | [hash/fnv](https://pkg.go.dev/hash/fnv)         |
| FNV1_128    | No            | 16          | [hash/fnv](https://pkg.go.dev/hash/fnv)         |
| FNV1a_128   | No            | 16          | [hash/fnv](https://pkg.go.dev/hash/fnv)         |
| CRC32       | No            | 4           | [hash/crc32](https://pkg.go.dev/hash/crc32)     |
| CRC64ISO    | No            | 8           | [hash/crc64](https://pkg.go.dev/hash/crc64)     |
| CRC64ECMA   | No            | 8           | [hash/crc64](https://pkg.go.dev/hash/crc64)     |
| BLAKE2s_256 | Yes           | 32          | [golang.org/x/crypto/blake2s](https://pkg.go.dev/golang.org/x/crypto/blake2b) |
| BLAKE2b_256 | Yes           | 32          | [golang.org/x/crypto/blake2b](https://pkg.go.dev/golang.org/x/crypto/blake2b) |
| BLAKE2b_384 | Yes           | 48          | [golang.org/x/crypto/blake2b](https://pkg.go.dev/golang.org/x/crypto/blake2b) |
| BLAKE2b_512 | Yes           | 64          | [golang.org/x/crypto/blake2b](https://pkg.go.dev/golang.org/x/crypto/blake2b) |
| XX          | No            | 8           | [github.com/cespare/xxhash](https://pkg.go.dev/github.com/cespare/xxhash/v2)  |

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

- all functions that accept `[]byte`.

### Benchmark Tests

Benchmark tests are implemented for

- all functions.

### Chaos Tests

Not planned.

## Future works

Not planned.

## References

- [https://pkg.go.dev/hash](https://pkg.go.dev/hash)
- [https://pkg.go.dev/crypto](https://pkg.go.dev/crypto)
- [https://pkg.go.dev/golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)
- [Hash Functions - NIST](https://csrc.nist.gov/projects/hash-functions)
- [FIPS 198-1 The Keyed-Hash Message Authentication Code (HMAC)](https://csrc.nist.gov/pubs/fips/198-1/final)
- [NIST SP 800-224 Keyed-Hash Message Authentication Code (HMAC): Specification of HMAC and Recommendations for Message Authentication](https://csrc.nist.gov/pubs/sp/800/224/ipd)
