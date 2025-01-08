
## Bundle signing key

public key and private key for signing bundles are generated with these command.

```bash
openssl genrsa -out private.key 2048
openssl rsa -in private.key -pubout -out public.key
```

## Build bundle

```bash
opa sign --signing-key private.key --bundle bundle/
mv .signatures.json bundle/
```

```bash
opa build --verification-key public.key --signing-key private.key --bundle bundle/
```
