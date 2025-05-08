# SOAPREST Middleware

## Summary

This is the technical document of app/middleware/soaprest package that provides SOAPRESTMiddleware.
SOAPRESTMiddleare provides conversion functionality between SOAP/XML and REST/JSON.

## Motivation

Converting between SOAP/XML and REST/JSON is required in scenarios where you wish to continue supporting legacy clients that can only handle SOAP/XML, while utilizing REST/JSON-based applications on the server side.

### Goals

- SOAPRESTMiddleware can convert SOAP/XML requests into REST/JSON format.
- SOAPRESTMiddleware can convert REST/JSON responses back into SOAP/XML format.

### Non-Goals

## Technical Design

### Converting requests/responses

SOAPRESTMiddleware converts requests in SOAP/XML format into REST/JSON format and transforms responses in REST/JSON format back into SOAP/XML format.

![soaprest-middleware.svg](./img/soaprest-middleware.svg)

SOAPRESTMiddleware is positioned to be the first middleware applied, allowing it to work with other middleware designed for JSON.

SOAPRESTMiddleware implements `core.Middleware` interface to work as middleware.

```go
type Middleware interface {
  Middleware(http.Handler) http.Handler
}
```

### Request Validation

1. Middleware validates whether it's a SOAP1.2 or SOAP1.1 request by checking:
   - The Content-Type header is `application/soap+xml`, treat as SOAP1.2 request
   - The Content-Type header is `text/xml` and the SOAPAction header exists, treat as SOAP1.1 request
2. If the request doesn't meet either condition, a VersionMismatch error with 403 Forbidden is returned.

### Converting Algorithms

SOAPRESTMiddleware supports following converting algorithms.

- Simple
- Rayfish
- BadgerFish

### Request Transformation

When a SOAP request is received, the middleware performs the following transformations:

1. Validates the SOAP version based on Content-Type:
   - SOAP 1.2: Content-Type is `application/soap+xml`
   - SOAP 1.1: Content-Type is `text/xml` with a SOAPAction header
   - If neither condition is met, returns a VersionMismatch error with HTTP 403 Forbidden status
2. Extracts any charset parameters from the Content-Type header for later use in the response.
3. Preserves SOAP action information by adding it as an `X-SOAP-Action` header:
   - For SOAP 1.2 requests: Uses the action parameter from the Content-Type header
   - For SOAP 1.1 requests: Uses the value from the SOAPAction header
4. Preserves the original Content-Type by adding it as an `X-Content-Type` header:
   - This allows backend services to identify the SOAP version (1.1 or 1.2)
   - The complete original Content-Type is preserved, including all parameters
5. Converts the SOAP XML request body to JSON format.
6. Creates a new request with:
   - JSON body
   - Content-Type set to `application/json`
   - Accept set to `application/json` to ensure the backend returns JSON response
   - Original headers preserved
   - SOAP action information consistently preserved in the `X-SOAP-Action` header for both SOAP versions
   - SOAP version information preserved in the `X-Content-Type` header

### Response Transformation

When a response is received from the backend service, the middleware transforms it back to the appropriate SOAP format:

1. The upstream handler's response is fully buffered by the middleware.
2. The middleware validates that the response Content-Type is `application/json`. If not, an InvalidContentType error with HTTP 500 status is returned.
3. The JSON response is converted back to SOAP format using the configured algorithm.
4. The Content-Type header is set to match the original request's SOAP format:
   - For SOAP 1.2 requests: Content-Type is set to `application/soap+xml` with any original charset parameters preserved
   - For SOAP 1.1 requests: Content-Type is set to `text/xml` with any original charset parameters preserved
5. The buffered response is sent as a complete SOAP XML message after all transformations are complete.

This middleware does not support streaming responses. It fully buffers upstream responses to ensure proper error handling and complete XML transformation.
Flush operations from upstream handlers are intercepted to prevent premature header commits that would interfere with proper error status codes if conversion issues occur.

## Test Plan

### Unit Tests

Unit tests are implemented and passed.

- Most functions and methods are covered with a coverage objective of 97%.
- Some components in `api.go` cannot have their coverage measured through unit test.
- For these components, we verify correct behavior through integration tests that confirm proper configuration is applied and the middleware functions as expected.

### Integration Tests

Integration tests are implemented with these aspects.

- SOAPRESTMiddleware works as middleware.
- SOAPRESTMiddleware works with input configuration.

### e2e Tests

e2e tests are implemented with these aspects.

- SOAPRESTMiddleware works as middleware.
- SOAPRESTMiddleware works with input configuration.

### Fuzz Tests

Not planned.

### Benchmark Tests

Not planned.

### Chaos Tests

Not planned.

## Future works

- [ ] Documentation regarding the algorithm.
