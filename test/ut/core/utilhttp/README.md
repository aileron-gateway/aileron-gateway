# Notes for creating testdata

## test.crt/test.key

```bash
openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out test.crt -keyout test.key
---
Country Name (2 letter code) [AU]:
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:test
Email Address []:
```
