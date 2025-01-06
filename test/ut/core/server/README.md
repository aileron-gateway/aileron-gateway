# Creating testdata

## server.crt and server.key


```
openssl genrsa 2048 > server.key
openssl req -new -key server.key -out server.csr
openssl x509 -req -days 36500 -in server.csr -signkey server.key -out server.crt
```
