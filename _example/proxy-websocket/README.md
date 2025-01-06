# Proxy WebSocket Example

## About this example

This example shows reverse proxy handler can proxy WebSocket.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
├── server  !!!!! This binary is built from the server.go in the Test section. 
└── _example/
    └── proxy-websocket/
        ├── index.html
        ├── server.go
        └── config.yaml
```

## Run

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/proxy-websocket/
```

## Test

Before testing the WebSocket proxy, we need to build a WebSocket server.

So, run the following command to build a server.
A binary `server` will be created.

```bash
go build _example/proxy-websocket/server.go
```

Then run the WebSocket server.

```bash
$ ./server

2024/08/24 20:46:16 WebSocket server listens at 0.0.0.0:9999
```

The WebSocket server has 2 handlers of

- `/`: Return index.html.
- `/ws`: Do WebSocket communication.

So, access to the `index.html` through  [http://localhost:8080/_example/proxy-websocket/](http://localhost:8080/_example/proxy-websocket/) from your browser.

Once you accessed to the index.html, a WebSocket communication will start.
You will receive the current datetime.
And you can send some messages.

```text
Hello!! This is a WebSocket server!!
It's Sat, 24 Aug 2024 20:51:17 GMT
It's Sat, 24 Aug 2024 20:51:18 GMT
It's Sat, 24 Aug 2024 20:51:19 GMT
It's Sat, 24 Aug 2024 20:51:20 GMT
It's Sat, 24 Aug 2024 20:51:21 GMT
Your message arrived: hello!
It's Sat, 24 Aug 2024 20:51:22 GMT
It's Sat, 24 Aug 2024 20:51:23 GMT
It's Sat, 24 Aug 2024 20:51:24 GMT
It's Sat, 24 Aug 2024 20:51:25 GMT
It's Sat, 24 Aug 2024 20:51:26 GMT
It's Sat, 24 Aug 2024 20:51:27 GMT
It's Sat, 24 Aug 2024 20:51:28 GMT
Your message arrived: hi!
It's Sat, 24 Aug 2024 20:51:29 GMT
It's Sat, 24 Aug 2024 20:51:30 GMT

..... output omitted
```
