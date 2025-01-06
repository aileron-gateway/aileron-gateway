# HTTP3 Sever Example

This folder is the example of a http3 server.
Browsers do not always use http3 but also use HTTP2.
If you have noticed that the browser you are using does not use http3, try to change your browsers.
We have checked that the Firefox uses http3.

**NOTICE: HTTP3 will be used for untrusted certificates.**
**Self-signed certs won't work with HTTP3
HTTP3 uses Quic protocol at transport layer.
AILERON uses [quic-go](https://github.com/quic-go/quic-go) for communication.

For logging and debugging quic, see the link below.
Environmental variable `QUIC_GO_LOG_LEVEL=debug` shows debug logs about Quic.

[https://github.com/quic-go/quic-go/wiki/Logging](https://github.com/quic-go/quic-go/wiki/Logging)

## Run

Before starting the http3 server, you have to place your SSL certificates.
Certificates have to be trusted by the browser.

Default certs paths are (see config.yaml)

- ./pki/cert.crt
- ./pki/cert.key

[https://github.com/FiloSottile/mkcert](https://github.com/FiloSottile/mkcert) is one of the util tools that can generate SSL certificate.

Run the command to start a http3 server.

```bash
./aileron -f _example/http3-server/
```

## Test

To test this example, access to [https://localhost:8443/](https://localhost:8443/) from tour browser.
Hostname should be the same with the value in SSL certificate.

## Notice

When you get an error like below, visit the url as mentioned and follow the step there.

- [https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes](https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes).

```text
failed to sufficiently increase receive buffer size (was: 208 kiB, wanted: 2048 kiB, got: 416 kiB). See https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes for details.
```
