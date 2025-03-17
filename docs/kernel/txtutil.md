# Package `kernel/txtutil`

## Summary

This is the design document of `kernel/txtutil` package.

`kernel/txtutil` package provides text processing utilities.

## Motivation

Text processing is one of the fundamental procession in any application including API Gateways.
It can be used for HTTP header and body processing, generating contents and so on.
It will be useful to isolate and share the utility of functions.

Usecase:

- HTTP header value replace
- HTTP body replace
- HTTP response body generation
- Static content generation

### Goals

- Provides generalized text processing functions.
    - Text matching functions.
    - Text replacing functions.
    - Handing template functions

### Non-Goals

- Add functions that would be used for specific purposes.

## Technical Design

### Matchers

Matcher functions provides `string` or `[]byte` matching functions.

As an generalized interface, the matching functions have the following signature.

```go
// MatchFunc[string] for string matching.
// MatchFunc[[]byte] for []byte matching.
type MatchFunc[T any] func(T) bool
```

`Matcher` interface have a method that satisfy the `MatchFunc[T any]`.

```go
// Matcher[string] for string matching.
// Matcher[[]byte] for []byte matching.
type Matcher[T any] interface {
  Match(T) bool
}
```

Supported matching types are listed below.
Note that the match types can be added or removed in future updates.

| Match type | Used method                                                   | Try on Go Playground                                                   |
| ---------- | ------------------------------------------------------------- | ---------------------------------------------------------------------- |
| Exact      | ==                                                            | [https://go.dev/play/p/tzXNQYFeEbm](https://go.dev/play/p/tzXNQYFeEbm) |
| Prefix     | [strings#HasPrefix](https://pkg.go.dev/strings#HasPrefix)     | [https://go.dev/play/p/f_eU7-K49ZV](https://go.dev/play/p/f_eU7-K49ZV) |
| Suffix     | [strings#HasSuffix](https://pkg.go.dev/strings#HasSuffix)     | [https://go.dev/play/p/dOpLnzu74nv](https://go.dev/play/p/dOpLnzu74nv) |
| Contains   | [strings#Contains](https://pkg.go.dev/strings#Contains)       | [https://go.dev/play/p/tbFRqJTL7vt](https://go.dev/play/p/tbFRqJTL7vt) |
| Path       | [path#Match](https://pkg.go.dev/path#Match)                   | [https://go.dev/play/p/xY56ZBgCGrO](https://go.dev/play/p/xY56ZBgCGrO) |
| FilePath   | [filepath#Match](https://pkg.go.dev/path/filepath#Match)      | [https://go.dev/play/p/dHl5dof11ZF](https://go.dev/play/p/dHl5dof11ZF) |
| Regex      | [regexp#Regexp.Match](https://pkg.go.dev/regexp#Regexp.Match) | [https://go.dev/play/p/AjWEC9C_YIh](https://go.dev/play/p/AjWEC9C_YIh) |
| RegexPOSIX | [regexp#Regexp.Match](https://pkg.go.dev/regexp#Regexp.Match) | [https://go.dev/play/p/yP1LIgg0PAe](https://go.dev/play/p/yP1LIgg0PAe) |

### Replacers

Replacer functions provides `string` or `[]byte` replacing functions.

As an generalized interface, the replace functions have the following signature.

```go
// ReplaceFunc[string] for replacing string values.
// ReplaceFunc[[]byte] for replacing []byte values.
type ReplaceFunc[T any] func(T) T
```

`Replacer` interface have a method that satisfy the `ReplaceFunc[T any]`.

```go
// Replacer[string] for replacing string values.
// Replacer[[]byte] for replacing []byte values.
type Replacer[T any] interface {
  Replace(T) T
}
```

Supported replacing types are listed below.
Note that the replace types can be added or removed in future updates.

| Replace type | Used method                                                             | Try on Go Playground                                                   |
| ------------ | ----------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| Noop         | -                                                                       | [https://go.dev/play/p/pEaudzyUeOX](https://go.dev/play/p/pEaudzyUeOX) |
| Fixed        | -                                                                       | [https://go.dev/play/p/gIBOBblqe4w](https://go.dev/play/p/gIBOBblqe4w) |
| Value        | [strings#ReplaceAll](https://pkg.go.dev/strings#ReplaceAll)             | [https://go.dev/play/p/Qz5wmj2MKca](https://go.dev/play/p/Qz5wmj2MKca) |
| Left         | [strings#Repeat](https://pkg.go.dev/strings#Repeat)                     | [https://go.dev/play/p/_XnzkUO5DNE](https://go.dev/play/p/_XnzkUO5DNE) |
| Right        | [strings#Repeat](https://pkg.go.dev/strings#Repeat)                     | [https://go.dev/play/p/nDmuJ1F-TX2](https://go.dev/play/p/nDmuJ1F-TX2) |
| Trim         | [strings#Trim](https://pkg.go.dev/strings#Trim)                         | [https://go.dev/play/p/p-O6aMuV_R0](https://go.dev/play/p/p-O6aMuV_R0) |
| TrimLeft     | [strings#TrimLeft](https://pkg.go.dev/strings#TrimLeft)                 | [https://go.dev/play/p/8tSm5oex608](https://go.dev/play/p/8tSm5oex608) |
| TrimRight    | [strings#TrimRight](https://pkg.go.dev/strings#TrimRight)               | [https://go.dev/play/p/xuH7DdM0yN_y](https://go.dev/play/p/xuH7DdM0yN_y) |
| TrimPrefix   | [strings#TrimPrefix](https://pkg.go.dev/strings#TrimPrefix)             | [https://go.dev/play/p/jhIFvHB8FoH](https://go.dev/play/p/jhIFvHB8FoH) |
| TrimSuffix   | [strings#TrimSuffix](https://pkg.go.dev/strings#TrimSuffix)             | [https://go.dev/play/p/uY2WNAPlZ9M](https://go.dev/play/p/uY2WNAPlZ9M) |
| Encode       | Regexp + Encode                                                         | [https://go.dev/play/p/yy1QvqHOjue](https://go.dev/play/p/yy1QvqHOjue) |
| Hash         | Regexp + Hash                                                           | [https://go.dev/play/p/Hjc8Z2rR1wz](https://go.dev/play/p/Hjc8Z2rR1wz) |
| Regexp       | [regexp#Regexp.ReplaceAll](https://pkg.go.dev/regexp#Regexp.ReplaceAll) | [https://go.dev/play/p/T6Mu9usjstw](https://go.dev/play/p/T6Mu9usjstw) |
| Expand       | [regexp#Regexp.Expand](https://pkg.go.dev/regexp#Regexp.Expand)         | [https://go.dev/play/p/Pw6tLUVZYcO](https://go.dev/play/p/Pw6tLUVZYcO) |
| Encrypt      | Regexp + Encryption                                                     | [https://go.dev/play/p/mBvv3DGInhg](https://go.dev/play/p/mBvv3DGInhg) |

### Fasttemplate

Fasttemplate is the simple and very fast but limited feature template.
It primarily intended to be used for

- Error message formatting.
- Log message formatting.

This feature is the replacement of [valyala/fasttemplate](https://github.com/valyala/fasttemplate).
[valyala/fasttemplate](https://github.com/valyala/fasttemplate) is a nice package but raises many panics because of the limited type support.

This feature has no generalized signatures or interfaces.
Use the template struct directly.
This is an overview usage of the template.

```go
input := map[string]any{
  "foo": "alice",
  "bar": "bob",
}

// func(tpl string, start string, end string) *txtutil.FastTemplate
tpl := txtutil.NewFastTemplate(`Hello {{foo}} and {{bar}}!`, "{{", "}}")

fmt.Println(string(tpl.Execute(input))) //  Hello alice and bob!
```

### Template

Template functions provides document template features.
It contains a function that generate a content from template using external information.

`Template` is defined a the common interface for template content.
Its Content method returns the content generated from the template.
It accepts external, or embed-able, information.

```go
type Template interface {
  Content(map[string]any) []byte
}
```

Following types of templates are defined.
Note that the template types can be added or removed in future development.

| Template type | Description                             | Embed values  | Used Package                                      |
| ------------- | --------------------------------------- | ------------- | ------------------------------------------------- |
| Text          | Plain text template                     | Not supported | -                                                 |
| GoText        | Text template using go standard package | Supported     | [text/template](https://pkg.go.dev/text/template) |
| GoHTML        | HTML template using go standard package | Supported     | [html/template](https://pkg.go.dev/html/template) |

For example,

template (this format is just an example)

```json
{
  "status": {{ status }},
  "code": "{{ code }}",
  "message": "{{ message }}" 
}
```

and external information

```go
map[string]any{
  "status": 500,
  "code": "E1234",
  "message": "cannot connect to the session storage",
}
```

will produce the following content.

```json
{
  "status": 500,
  "code": "E1234",
  "message": "cannot connect to the session storage" 
}
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

Fuzz tests are supposed to be implemented for

- functions that accepts `[]byte` data.
- functions that accepts `string` data.

### Benchmark Tests

Not planned.

### Chaos Tests

Not planned.

## Future works

Not planned.

## References

- [API Specification - Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/reference/spec/)
- [location - Nginx](https://nginx.org/en/docs/http/ngx_http_core_module.html#location)
- [ApisixRoute](https://apisix.apache.org/docs/ingress-controller/concepts/apisix_route/)
- [Routers - traefik](https://doc.traefik.io/traefik/routing/routers/)
- [About Expressions Language - Kong](https://docs.konghq.com/gateway/latest/reference/expressions-language/)
- [Module ngx_http_rewrite_module](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html)
