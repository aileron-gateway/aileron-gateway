Generate a self-signed certificate.

```
$ openssl req -new -x509 -nodes -days 365 -subj '/CN=server-ca' -keyout server-ca.key -out server-ca.crt
```

Generate server-side private keys and certificates.

```
$ openssl genrsa -out server.key
$ openssl req -new -key server.key -subj '/CN=localhost' -out server.csr
$ openssl x509 -req -in server.csr -CA server-ca.crt -CAkey server-ca.key -CAcreateserial -days 365 -out server.crt -extfile <(printf "subjectAltName=DNS:localhost")
```

Check the server certificate.

```
$ openssl x509 -in server.crt -text -noout
```
