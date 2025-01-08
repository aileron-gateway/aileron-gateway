# Generate Private Key and Public Key

RFC 7518 JSON Web Algorithms (JWA)

<https://datatracker.ietf.org/doc/rfc7518/>

```text
+--------------+--------------------------------------------------+--------------------+
| "alg" Param  | Digital Signature or MAC                         | Implementation     |
| Value        | Algorithm                                        | Requirements       |
+--------------+--------------------------------------------------+--------------------+
| HS256        | HMAC using SHA-256                               | Required           |
| HS384        | HMAC using SHA-384                               | Optional           |
| HS512        | HMAC using SHA-512                               | Optional           |
| RS256        | RSASSA-PKCS1-v1_5 using SHA-256                  | Recommended        |
| RS384        | RSASSA-PKCS1-v1_5 using SHA-384                  | Optional           |
| RS512        | RSASSA-PKCS1-v1_5 using SHA-512                  | Optional           |
| ES256        | ECDSA using P-256 and SHA-256                    | Recommended+       |
| ES384        | ECDSA using P-384 and SHA-384                    | Optional           |
| ES512        | ECDSA using P-521 and SHA-512                    | Optional           |
| PS256        | RSASSA-PSS using SHA-256 and MGF1 with SHA-256   | Optional           |
| PS384        | RSASSA-PSS using SHA-384 and MGF1 with SHA-384   | Optional           |
| PS512        | RSASSA-PSS using SHA-512 and MGF1 with SHA-512   | Optional           |
| none         | No digital signature or MAC performed            | Optional           |
+--------------+--------------------------------------------------+--------------------+
```

## RS and PS

```bash
openssl genrsa -out private-rs.key 2048
openssl rsa -in private-rs.key -pubout -out public-rs.key
```

View content.

```bash
openssl rsa -pubin -in private.key -text
```

openssl ecparam -name secp384r1 -genkey -noout -out sec1_ec_p384_private.pem
rm sec1_ec_p384_private.pem
openssl ec -in ec_p384_private.pem -pubout -out ec_p384_public.pem

## ES

```bash
openssl ecparam -list_curves
```

```bash
openssl ecparam -name $CURVE -genkey -noout -out private.key
openssl ec -in private.key -pubout -out public.key
```

When using [JWT.io](https://jwt.io/), use PKCS#8 for EC private keys and SPKI for EC public keys.
See [How to create private key for JWT ES384](https://stackoverflow.com/questions/71856396/how-to-create-private-key-for-jwt-es384).

```bash
openssl ecparam -name $CURVE -genkey -noout -out private-tmp.key
openssl pkcs8 -topk8 -nocrypt -in private-tmp.key -out private.key
openssl ec -in private.key -pubout -out public.key
```

`$CURVE` depends on the algorithms.

- `ES256`: prime256v1
- `ES384`: secp384r1
- `ES512`: secp521r1

## HS

Use arbitrary key string.

openssl ecparam -name $CURVE -genkey -noout -out private.key
openssl ec -in private.key -pubout -out public.key

## Memos

All commands for`ES256`. This can be used at [JWT.io](https://jwt.io/).

```bash
openssl ecparam -name prime256v1 -genkey -noout -out private-tmp.key
openssl pkcs8 -topk8 -nocrypt -in private-tmp.key -out private-es256.key
openssl ec -in private-es256.key -pubout -out public-es256.key
rm private-tmp.key
```

All commands for`ES384`. This can be used at [JWT.io](https://jwt.io/).

```bash
openssl ecparam -name secp384r1 -genkey -noout -out private-tmp.key
openssl pkcs8 -topk8 -nocrypt -in private-tmp.key -out private-es384.key
openssl ec -in private-es384.key -pubout -out public-es384.key
rm private-tmp.key
```

All commands for`ES512`. This can be used at [JWT.io](https://jwt.io/).

```bash
openssl ecparam -name secp521r1 -genkey -noout -out private-tmp.key
openssl pkcs8 -topk8 -nocrypt -in private-tmp.key -out private-es512.key
openssl ec -in private-es512.key -pubout -out public-es512.key
rm private-tmp.key
```
