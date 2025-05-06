```
# 自己署名証明書の生成
$ openssl req -new -x509 -nodes -days 365 -subj '/CN=client-ca' -keyout client-ca.key -out client-ca.crt
 
# クライアント側の秘密鍵、証明書の生成
$ openssl genrsa -out client.key
$ openssl req -new -key client.key -subj '/CN=localhost' -out client.csr
$ openssl x509 -req -in client.csr -CA client-ca.crt -CAkey client-ca.key -CAcreateserial -days 365 -out client.crt -extfile <(printf "subjectAltName=DNS:localhost")

# クライアント証明書の確認
$ openssl x509 -in client.crt -text -noout
```
