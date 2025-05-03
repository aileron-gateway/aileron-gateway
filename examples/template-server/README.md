# Template Handler Example

## About this example

This example shows how to configure template handler that returns static contents.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    ├── template/
    │   └── config.yaml
    └── template.html
```

## Run

A HTTP server with a template handler will be run with the command.
The server listens on [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/template/
```

## Test

Send a HTTP request like below.
Specify the content type with `Accept` header.

```bash
$ curl -H "Accept: text/plain" http://localhost:8080/

AILERON Gateway
Hello World!!
```

When a template content was not defined for the content types
provided by the Accept header, the gateway returns an error.

```bash
$ curl -H "Accept: text/css" http://localhost:8080/

{
  "status":406,
  "statusText":"Not Acceptable"
}
```

Template can use request information in it.

```bash
$ curl -H "Accept: text/plain" http://localhost:8080/

AILERON Gateway
Hello World!!
```

```bash
$ curl -H "Accept: application/json" http://localhost:8080/

{
  "app": "AILERON Gateway",
  "hello": "World!"
}
```

```bash
$ curl -H "Accept: text/html" -H "foo: bar" http://localhost:8080/?hello=world

<!doctype html>
<html>

<head>
  <title>AILERON Gateway</title>
</head>

<body>
  <h1>AILERON Gateway</h1>
  <p>
    proto : HTTP/1.1</br>
    host : </br>
    method : GET</br>
    path : /</br>
    remote : 127.0.0.1:40608</br>
    header : map[Accept:[text/html] Foo:[bar] User-Agent:[curl/7.68.0]]</br>
    query : map[hello:[world]]</br>
  </p>
</body>

</html>
```
