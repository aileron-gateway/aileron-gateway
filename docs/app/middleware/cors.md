# CORS Middleware

## Summary

This is the design document of the app/middleware/cors package that provides CORSMiddleware.
CSRFMiddleware provides ability to secure APIs from CSRF vulnerability.

## Motivation

CORS is one of the basic security features required for protecting APIs.

### Goals

- CORSMiddleware protect APIs from CORS vulnerability.

### Non-Goals

## Technical Design

### CORS

[CORS (Cross-Origin Resource Sharing)](https://www.w3.org/Security/wiki/CORS) is a specification of sharing resources between different web 'origin'. The concept of 'origin' is defined in the [RFC 6454 - The Web Origin Concept](https://datatracker.ietf.org/doc/rfc6454/). CORS Specification is found in [The Fetch standard](https://www.w3.org/TR/cors/).

### Origin identification

The origin of the API requests must be identified with the combination of scheme, host, port based on [RFC6454 - 4. Origin of a URI](https://www.rfc-editor.org/rfc/rfc6454.html#section-4).
AILERON Gateway can be configured to apply CORS policy for the combinations of origins and HTTP methods.

- scheme
    - `http` or `https`
    - other schemes are not allowed
- host
    - FQDN
- port
- HTTP methods

### Headers

#### Policy headers

**[Cross-Origin-Embedder-Policy (COEP)](https://docs.w3cub.com/http/headers/cross-origin-embedder-policy)**

COEP header can be configured with the following syntax.

```text
Cross-Origin-Embedder-Policy: unsafe-none | require-corp
```

**[Cross-Origin-Opener-Policy (COOP)](https://docs.w3cub.com/http/headers/cross-origin-opener-policy)**

COOP header can be configured with the following syntax.

```text
Cross-Origin-Opener-Policy: unsafe-none | same-origin-allow-popups | same-origin
```

**[Cross-Origin-Resource-Policy (CORP)](https://docs.w3cub.com/http/headers/cross-origin-resource-policy)**

CORP header can be configured with the following syntax.

```text:Syntax
Cross-Origin-Resource-Policy: same-site | same-origin | cross-origin
```

#### Access Control Request Headers

**Origin**

Origin header is the origin of the request.
AILERON Gateway return ths value as `Access-Control-Allow-Origin` if the value is not "null".

**Access-Control-Request-Method**

AILERON Gateway does not check this header and always return `Access-Control-Allow-Methods`.
This header is set in the `Vary` header.

**Access-Control-Request-Headers**

AILERON Gateway does not check this header and always return `Access-Control-Allow-Headers` if some values are configured.
This header is set in the `Vary` header.

#### Access Control Response Headers

**Access-Control-Allow-Origin**

Access-Control-Allow-Origin is the allowed origin.
AILERON Gateway can set wildcard origin "*" and "\<origin\>" based on the incoming requests.
AILERON Gateway cannot return null origin.

```text:Syntax
Access-Control-Allow-Origin: *
Access-Control-Allow-Origin: <origin>
Access-Control-Allow-Origin: null
```

`Vary: Origin` header is also set when this header is responded except for the when the value is wildcard.

**Access-Control-Allow-Methods**

Access-Control-Allow-Methods is the list of allowed methods.
CORS-safelisted methods, `GET`, `HEAD`, `POST` are set by default.

```text:Syntax
Access-Control-Allow-Methods: <method>, <method>, …
Access-Control-Allow-Methods: *
```

**Access-Control-Allow-Headers**

Access-Control-Allow-Headers is the list of header names which can use in the actual API requests.
[CORS-safelisted request headers](https://developer.mozilla.org/en-US/docs/Glossary/CORS-safelisted_request_header) (`Accept`, `Accept-Language`, `Content-Language`, `Content-Type`) are not need to set in this header.

```text:Syntax
Access-Control-Allow-Headers: <header-name>, <header-name>, ...
Access-Control-Allow-Headers: *
```

**Access-Control-Expose-Headers**

Access-Control-Expose-Headers is the list of header names which are exposed to the frontend script.

```text:Syntax
Access-Control-Expose-Headers: <header-name>, <header-name>, ...
Access-Control-Expose-Headers: *
```

**Access-Control-Allow-Credentials**

Access-Control-Allow-Credentials is the flag if the frontend javascript can send cookie or authorization header.

```text:Syntax
Access-Control-Allow-Credentials: true
```

**Access-Control-Max-Age**

Access-Control-Max-Age is the age to cache the response (allowed methods and allowed headers) of preflight request.

```text:Syntax
Access-Control-Max-Age: <seconds>
```

### Preflight Passthrough

AILERON Gateway can proxy preflight requests to upstream servers.
This is called `Preflight Passthrough`.

## Test Plan

### Unit Tests

Unit tests are implemented and passed.

- All functions and methods are covered.
- Coverage objective 98%.

### Integration Tests

### e2e Tests

### Fuzz Tests

Not planned.

### Benchmark Tests

Not planned.

### Chaos Tests

Not planned.

## Future works

Not planned.

## References

- [RFC 6454 - The Web Origin Concept](https://datatracker.ietf.org/doc/rfc6454/)
- [CORS Wiki - W3C](https://www.w3.org/wiki/CORS)
- [CORS Security Wiki - W3C](https://www.w3.org/Security/wiki/CORS)
- [The Fetch standard](https://www.w3.org/TR/cors/)
- [Cross-Origin Resource Sharing (CORS) - mdn web docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [Testing Cross Origin Resource Sharing - OWASP](https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/11-Client-side_Testing/07-Testing_Cross_Origin_Resource_Sharing)
