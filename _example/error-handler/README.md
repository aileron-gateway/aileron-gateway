# Error Handler Example

## About this example

This example shows how to configure error handler and replace the default one.
When an error handler is not configured, default error handler will be used.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    ├── error-handler/
    │   └── config.yaml
    └── html/
        └── index.html
```

## Run

A static file server with custom error handler will be run with the command.
The server listens on [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/error-handler/
```

## Test

Send a HTTP request like this.
The static file server returns html/index.html.
No access logs and error logs are output because there is no error on serving the content.

```bash
$ curl http://localhost:8080/

<!doctype html>
<html>
  <head>
    <title>AILERON Gateway</title>
  </head>
  <body>
    <h1>AILERON Gateway</h1>
  </body>
</html>
```

Next, send a HTTP request to wrong URL path.
This will return customized error.

```bash
$ curl http://localhost:8080/wrong

{
  "app": "AILERON Gateway",
  "example": "core/v1/ErrorHandler",
  "code": "E2141",
  "kind": "CoreStaticServer"
}
```

Use `Accept` header to specify the content type to be returned.

```bash
$ curl -H "Accept: application/xml" http://localhost:8080/wrong

<?xml version="1.0" encoding="UTF-8" ?>
<error>
    <app>AILERON Gateway</app>
    <example>core/v1/ErrorHandler</example>
    <code>E2141</code>
    <kind>CoreStaticServer</kind>
</error>
```
