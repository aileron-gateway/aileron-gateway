# Header Cert Example

## About this example

This example shows how to configure header cert.

Files that required to run this example is shown below.
Make sure to build or download AILERON Gateway binary, `aileron`.

```txt
./
├── aileron
└── _example/
    └── header-cert/
        ├── pki/
        │   ├── client.crt
        │   ├── client.key
        │   ├── fingerprint.txt (Optional)
        │   └── rootCA.crt
        └── config.yaml
```

## Run

Run the example with this command.
Reverse proxy server will listen on  [http://localhost:8080/](http://localhost:8080/).

```bash
./aileron -f _example/header-cert/
```

## Test

Send a HTTP request like below.
The default upstream server of the reverse proxy is [http://httpbin.org/](http://httpbin.org/)

This request returns 200 OK status.
No error logs in the gateway log.

```bash
$ curl -v --cert pki/client.crt --key pki/client.key \
     -H "X-SSL-Client-Cert: $(base64 -w 0 _example/header-cert/pki/client.crt)" \
     -H "X-SSL-Client-Fingerprint: $(cat _example/header-cert/pki/fingerprint.txt)" \
     http://localhost:8080/

> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.11.1
> Accept: */*
> X-SSL-Client-Cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMrakNDQWVJQ0ZIb2xCZ2Z2VzEydm1vTG9lc3k3T1Z2N3BNMnhNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1ENHgKQ3pBSkJnTlZCQVlUQWtwUU1SUXdFZ1lEVlFRS0RBdGxlR0Z0Y0d4bExtTnZiVEVaTUJjR0ExVUVBd3dRUlZoQgpUVkJNUlNCUVMwa2dVazlQVkRBZUZ3MHlOVEF6TWpZd016VTRNREJhRncwek5UQXpNalF3TXpVNE1EQmFNRFV4CkN6QUpCZ05WQkFZVEFrcFFNUlF3RWdZRFZRUUtEQXRsZUdGdGNHeGxMbU52YlRFUU1BNEdBMVVFQXd3SFEweEoKUlU1VU1UQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUtIeldHeERydU1VNklVWQpyWUhHdHk0SkIrNCtSN0w4ckttSzlYNlF3WkhnT1JSM0ZWU2pBYkZ0Y3J1TmdqNUR0OXFXblg2OWJUd0h5S0orCkVtRW5OeDJYZE1Gd0FvcEUzSmR4OXhTa2ovSnQwZjQvWVAxVXRiZFNhVmJPek5MUUpNOWhDMWYzTE41Yk8yVTcKamJWTEtkYW5sSU9LRlBUcyswbDFIYVdEdC9mUlVTY2JvSE01Z3I4SC9Yc2J5blAxcTFqbTNQaVVHakwyL0xhdgpxUDVaUjJRbm9nSjZxNm4ydC9GLzVCQTRyMWJ3Tm9wT1k3L3ZqR0FNK2pSdTVkWFFDbXZteEVSUkY0NFEzdmxqCktwMHdLSGE2MnF5Nit3SXc1MVVXQU5RbHNWWUg4WHU0M09LWDRXZ0JvY2E2akVSdE9ZMTByNzVCRGlmNXRlMEgKTWM4ci95c0NBd0VBQVRBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQUlYbnpERGNVZmZ5Uk5hMTIxTitaSWFQbwpBZitCUGpObWZoNlJ6KzdDMmNERmovVUNZRmlXWnNDbzhoUDVuUTFybFh5TXkxcitmRjJldlB4d0N0a05jQnNVClNaRVJnZmRPLytGbzhNeFRSZDZtcC9VN1RQcndsZERpekIrcDRRUXZzTEdPaTlkRTFHdWZwM3VncER1WHd0RmgKeWZnMVI5ekgxeUxvOS9MVzNXUzVFVGRZWDgzc0RyOHMyQWRwRmhZbVIxWXV2cnJVZTh5Yk8wSzAvYzJ3bUp3MAp1OGxGMHdGOVpaOGwyRzhGZ0FYMHY2aU1tSTNSMDc1SVdhb0t5b1ZZckxZZDNFeFpLRHZmdTE2ajhwaXJqbURYCncxcHFmNUNwemZ2ZEVoNmViajA0VENLYTh4VHdUczc2YzA2Q3lWRkEvYTRCeitDK2FhUWVvVXRTakNIN0JBPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==
> X-SSL-Client-Fingerprint: f5f170f7773d0620c276f1e74a6d36f9f32c42759a86d45dcbe38c4eb91d20cf

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 9593
< Content-Type: text/html; charset=utf-8
< Date: Mon, 31 Mar 2025 00:52:49 GMT
< Server: gunicorn/19.9.0
```

Setting a fingerprint to the header is optional.

This request also returns 200 OK status.

```bash
$ curl -v --cert pki/client.crt --key pki/client.key \
     -H "X-SSL-Client-Cert: $(base64 -w 0 _example/header-cert/pki/client.crt)" \
     http://localhost:8080/

> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.11.1
> Accept: */*
> X-SSL-Client-Cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMrakNDQWVJQ0ZIb2xCZ2Z2VzEydm1vTG9lc3k3T1Z2N3BNMnhNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1ENHgKQ3pBSkJnTlZCQVlUQWtwUU1SUXdFZ1lEVlFRS0RBdGxlR0Z0Y0d4bExtTnZiVEVaTUJjR0ExVUVBd3dRUlZoQgpUVkJNUlNCUVMwa2dVazlQVkRBZUZ3MHlOVEF6TWpZd016VTRNREJhRncwek5UQXpNalF3TXpVNE1EQmFNRFV4CkN6QUpCZ05WQkFZVEFrcFFNUlF3RWdZRFZRUUtEQXRsZUdGdGNHeGxMbU52YlRFUU1BNEdBMVVFQXd3SFEweEoKUlU1VU1UQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUtIeldHeERydU1VNklVWQpyWUhHdHk0SkIrNCtSN0w4ckttSzlYNlF3WkhnT1JSM0ZWU2pBYkZ0Y3J1TmdqNUR0OXFXblg2OWJUd0h5S0orCkVtRW5OeDJYZE1Gd0FvcEUzSmR4OXhTa2ovSnQwZjQvWVAxVXRiZFNhVmJPek5MUUpNOWhDMWYzTE41Yk8yVTcKamJWTEtkYW5sSU9LRlBUcyswbDFIYVdEdC9mUlVTY2JvSE01Z3I4SC9Yc2J5blAxcTFqbTNQaVVHakwyL0xhdgpxUDVaUjJRbm9nSjZxNm4ydC9GLzVCQTRyMWJ3Tm9wT1k3L3ZqR0FNK2pSdTVkWFFDbXZteEVSUkY0NFEzdmxqCktwMHdLSGE2MnF5Nit3SXc1MVVXQU5RbHNWWUg4WHU0M09LWDRXZ0JvY2E2akVSdE9ZMTByNzVCRGlmNXRlMEgKTWM4ci95c0NBd0VBQVRBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQUlYbnpERGNVZmZ5Uk5hMTIxTitaSWFQbwpBZitCUGpObWZoNlJ6KzdDMmNERmovVUNZRmlXWnNDbzhoUDVuUTFybFh5TXkxcitmRjJldlB4d0N0a05jQnNVClNaRVJnZmRPLytGbzhNeFRSZDZtcC9VN1RQcndsZERpekIrcDRRUXZzTEdPaTlkRTFHdWZwM3VncER1WHd0RmgKeWZnMVI5ekgxeUxvOS9MVzNXUzVFVGRZWDgzc0RyOHMyQWRwRmhZbVIxWXV2cnJVZTh5Yk8wSzAvYzJ3bUp3MAp1OGxGMHdGOVpaOGwyRzhGZ0FYMHY2aU1tSTNSMDc1SVdhb0t5b1ZZckxZZDNFeFpLRHZmdTE2ajhwaXJqbURYCncxcHFmNUNwemZ2ZEVoNmViajA0VENLYTh4VHdUczc2YzA2Q3lWRkEvYTRCeitDK2FhUWVvVXRTakNIN0JBPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 9593
< Content-Type: text/html; charset=utf-8
< Date: Mon, 31 Mar 2025 02:09:25 GMT
< Server: gunicorn/19.9.0
```

A Request without client certificate is configured to be rejected.

```bash
$ curl -v  http://localhost:8080/

{"status":400,"statusText":"Bad Request"}
```

## About Certificates

The default certificates are generated by the following procedure.

We used [OpenSSL](https://www.openssl.org/) to create a range of certificates.

- [./pki/client.crt](./pki/client.crt)
- [./pki/client.key](./pki/client.key)


1. Create private key and certificate for the root CA
    ```bash
    $ openssl req -x509 -new \
                -newkey rsa:2048 -keyout rootCA.key -nodes \
                -sha256 \
                -days 3650 \
                -out rootCA.crt \
                -subj "/C=JP/O=example.com/CN=EXAMPLE PKI ROOT"

    ```
2. Create private key and client signature reqest for client certificate
    ```bash
    openssl req -newkey rsa:2048 -keyout client.key -nodes -out client.csr \
                -subj "/C=JP/O=example.com/CN=CLIENT"
    ```

3. Sign client certificate signature request with private key and certificate for root CA.
   ```bash
    openssl x509 -req \
                -in client.csr \
                -CA rootCA.crt -CAkey rootCA.key -CAcreateserial \
                -sha256 \
                -days 3650 \
                -out client.crt 
   ```

## About Fingerprint

The default fingerprint is generated by the following code.
- [./pki/fingerprint.txt](./pki/fingerprint.txt)

```go
package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	cert, err := os.ReadFile("pki/client.crt") // Set the path for the client.crt
    if err != nil {
        log.Fatalf("Failed to read client certificate: %v", err)
    }

	certBlock, _ := pem.Decode(cert)
	if certBlock == nil {
        log.Fatal("Failed to decode certificate")
	}
	fmt.Println("certBlock:", certBlock)

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
        log.Fatalf("Failed to decode client certificate: %v", err)
	}

	fp := sha256.Sum256(cert.Raw)
	fpHex := hex.EncodeToString(fp[:])
	fmt.Println(fpHex) // Use this fingerprint
}
```



