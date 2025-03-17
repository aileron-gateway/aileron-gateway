# Package `kernel/hash`

## Summary

This is the design document of `kernel/hash` package.
`kernel/hash` package provides hash utilities.

## Motivation

- Check integrity.
- Convert strings to fixed length one.
- Calculating signature.

So, it is reasonable to isolate this feature and make it re-usable.

### Goals

- Provide hashing utilities.

### Non-Goals

- Support various hash algorithms.
- Remove insecure hash algorithm. Users must use appropriate ones.

## Technical Design

### Interfaces

The signature to calculate a hash of certain data is defined like below.

- Hash functions have the following signature.
- Input data can be nil.
- No panics in the function.
- The length of the returned slice depends on the hashing algorithms.
- Hash functions MUST return the same result for the same input.

```go
func(data []byte) (hash []byte)
```

Note:
[hash/maphash](https://pkg.go.dev/hash/maphash) requires a seed for a hash table instance.
If the seed were not set, it generates different result for the same input.
Hash functions that have this kind of specification cannot be provided in this package.

Available hash algorithms are defined as `Algorithm` type.

```go
type Algorithm int
```

### Hash algorithms

[NIST approved hash functions](https://csrc.nist.gov/projects/hash-functions) specifies the following algorithms.
They have the property of [cryptographic hash function](https://en.wikipedia.org/wiki/Cryptographic_hash_function).
They contain [FIPS 180-4](https://csrc.nist.gov/pubs/fips/180-4/upd1/final).

- SHA-1
- SHA-2 family: SHA-224, SHA-256, SHA-384, SHA-512, SHA-512/224, SHA-512/256
- SHA-3 family: SHA3-224, SHA3-256, SHA3-384, SHA3-512
- SHAKE-128, SHAKE-256

Other commonly used hashes are

- [MD5](https://www.rfc-editor.org/info/rfc1321)
- [FNV](https://datatracker.ietf.org/doc/draft-eastlake-fnv/) family: FNV1, FNV1a
- [BLAKE2](https://datatracker.ietf.org/doc/rfc7693/) family: BLAKE2s, BLAKE2b

See [List of hash functions](https://en.wikipedia.org/wiki/List_of_hash_functions) for more lists.

By considering the use cases and implementation state of go standard packages,
this package implements the following hash functions.

Because this package provides atomic operations to calculate hash values,
it's users' responsibility to choose appropriate hash algorithm based on a situation.

NEVER use non-cryptographic hash algorithms for securing some data.

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

- all hashing functions that accept `[]byte`.

### Benchmark Tests

Benchmark tests are implemented for

- all hash functions.

### Chaos Tests

Not planned.

## Future works

Not planned.

## References

- [https://pkg.go.dev/hash](https://pkg.go.dev/hash)
- [https://pkg.go.dev/crypto](https://pkg.go.dev/crypto)
- [https://pkg.go.dev/golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)
- [Hash Functions - NIST](https://csrc.nist.gov/projects/hash-functions)
