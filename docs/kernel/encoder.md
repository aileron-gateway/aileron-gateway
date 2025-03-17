# Package `kernel/encoder`

## Summary

This is the design document of `kernel/encoder` package.

Encoding and decoding, marshaling and un-marshaling utilities are documented.

## Motivation

Encoding and decoding, marshaling and un-marshaling are used widely in the entire codebase.
For example, base64 encoding can be used for converting []byte data into string.
Protocol-buffer and json, yaml and other formats are converted to different formats based on the situation.
encode/decode and marshal/unmarshal utils are provided for reusability.

### Goals

- Provide encode/decode utilities.
- Provide marshal/un-marshal utilities.

### Non-Goals

- Support various encode/decode algorithms.
- Support various marshal/un-marshal algorithms.

## Technical Design

### encode/decode

Encoding and decoding are used for data format conversion.
Converting from binary to text is one of the typical use case of encoding.
Following table shows typical encoding algorithms.
See [Base32](https://en.wikipedia.org/wiki/Base32), [Base64](https://en.wikipedia.org/wiki/Base64), [Base85](https://en.wikipedia.org/wiki/Ascii85) for more variations.

|                  | Specs    | CompressionRatio | Use case examples                |
| ---------------- | -------- | ---------------- | -------------------------------- |
| Base16 (Hex)     | RFC 4648 | 8/4 = 2.00       | CSRF tokens, Session ID          |
| Base32           | RFC 4648 | 8/5 = 1.60       | Session ID, Request ID, Trace ID |
| Base32 Hex       | RFC 4648 | 8/5 = 1.60       | Session ID, Request ID, Trace ID |
| Base64           | RFC 4648 | 8/6 = 1.33...    | Session ID, Request ID, Trace ID |
| Base64 URL       | RFC 4648 | 8/6 = 1.33...    | Session ID, Request ID, Trace ID |
| Base85 (Ascii85) | -        | 5/4 = 1.25       | Data Transfer                    |

By considering the use cases and the support status of Go standard package [encoding](https://pkg.go.dev/encoding),
the package provides the following encoding algorithms.

- [Base16 (Hex)](https://pkg.go.dev/encoding/hex)
- [Base32](https://pkg.go.dev/encoding/base32)
    - Standard encode
    - Hex encode
    - Standard Escaped encode (See [Avoid vowels](#avoiding-vowels))
    - Hex Escaped encode (See [Avoid vowels](#avoiding-vowels))
- [Base64](https://pkg.go.dev/encoding/base64)
    - Standard encode
    - URL encode

Encoding functions have the following signature.

```go
// Signature of encode functions that return result in []byte.
func(data []byte) (encoded []byte)

// Signature of encode functions that return result in string.
func(data []byte) (encoded string)
```

Decoding functions have the following signature.

```go
// Signature of decode functions that accept data in []byte.
func(data []byte) (decoded []byte, err error)

// Signature of decode functions that accept data in string.
func(data string) (decoded []byte, err error)
```

#### Avoiding vowels

Sometimes, it is not pleasant that encoded strings contain unintended words like `ERROR`, `FATAL` or `WARN`.
To avoid such cases, a technique to re-define the characters used for encoding can be employed.

This package provides two encoding methods.

```go
// Base32 standard encoding with vowel escaped.
// Original Base32 uses "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567".
base32.NewEncoding("BCDFGHJKLMNPQRSTUVWXYZ0123456789")
```

```go
// Base32 hex encoding with vowel escaped.
// Original Base32Hex uses "0123456789ABCDEFGHIJKLMNOPQRSTUV".
base32.NewEncoding("0123456789BCDFGHJKLMNPQRSTUVWXYZ")
```

### marshal/unmarshal

[Marshaling](https://en.wikipedia.org/wiki/Marshalling_(computer_science)) has almost the same meaning as [Serialization](https://en.wikipedia.org/wiki/Serialization).
Following table shows some commonly used serialization methods.

The first column indicates if this package provides the marshaling and un-marshaling.

| Implements | Method          | Human readable | Schema driven | Notes           |
|----------- | --------------- | -------------- | ------------- | --------------- |
|            | XML             | Yes            | No            |                 |
| Yes        | JSON            | Yes            | No            |                 |
|            | CBOR            | No             | No            |                 |
| Yes        | YAML            | Yes            | No            |                 |
|            | Massage Pack    | No             | No            |                 |
| Yes        | Protocol Buffer | No             | Yes           |                 |
|            | Apache Avro     | No             | Yes           |                 |
|            | Apache Thrift   | No             | Yes           |                 |
|            | Gob             | No             | Yes           | Golang specific |

Marshal and un-marshal functions should have the following signature.

```go
// Marshal functions.
// Given object "in" will be marshaled and returned in []byte.
func(in any) (b []byte, err error)
```

```go
// Un-marshal functions.
// Given []byte data of "in" will be un-marshaled into the given "into" object.
func(in []byte, into any) error
```

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

- encode functions that accept `[]byte`.
- decode functions that accept `string`.
- un-marshal functions that accept `[]byte`.

### Benchmark Tests

Benchmark tests are implemented for

- encode functions.
- decode functions.
- marshal functions.
- un-marshal functions.

### Chaos Tests

Not planned.

## Future works

Not planned.

## References

- [https://pkg.go.dev/encoding](https://pkg.go.dev/encoding)
- [https://en.wikipedia.org/wiki/Serialization](https://en.wikipedia.org/wiki/Serialization)
- [https://en.wikipedia.org/wiki/Marshalling_(computer_science)](https://en.wikipedia.org/wiki/Marshalling_(computer_science))
