# Multi Servers Example

## About this example

This example shows how to configure logger servers.
Template handlers are used to return responses in this example.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── logging/
        └── config.yaml
```

## Run

Run the example with this command.
A HTTP server will listen on [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/logging/
```

## Test

Send HTTP requests to each servers like below.
No logs will be output to the standard output.

```bash
$ curl http://localhost:8080/contents

AILERON Gateway
```

Send another request.
Not found error logs will be output to the standard output.

```bash
$ curl http://localhost:8080/not-found

{"status":404,"statusText":"Not Found"}
```

Next, change the config file to output logs to files.
The following lines should be changed.

```yaml
# outputTarget: Stdout
outputTarget: File
```

Then, send requests many times with this script.
Logs will be output to log files and logs will be archived when then exceeded maximum size specified in the config.

```bash
while true; do 
  curl http://localhost:8080/not-found;
done;
```
