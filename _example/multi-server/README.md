# Multi Servers Example

## About this example

This example shows how to run multiple servers.
Template handlers are used to return responses in this example.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── multi-server/
        └── config.yaml
```

## Run

Run the example with this command.
Two servers will listen on [http://localhost:8081/](http://localhost:8081/) and [http://localhost:8082/](http://localhost:8082/).

```bash
./aileron -f _example/multi-server/
```

## Test

Send HTTP requests to each servers like below.
`Hello!!` and `Goodbye!!` will be returned from each server.

```bash
$ curl http://localhost:8081/

Hello!!
```

```bash
$ curl http://localhost:8082/

Goodbye!!
```
