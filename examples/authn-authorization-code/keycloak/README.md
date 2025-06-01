# Keycloak Configs

## Client and users

Registered realm:

- Realm: `aileron`
- Client ID: `aileron_authorization_code`
- Client Secret: `KWYPBgrTEEGZNH6wZsP2zyK14LZHJi77`
- Base64(ID:Secret): `YWlsZXJvbl9yZXNvdXJjZV9zZXJ2ZXI6S1dZUEJnclRFRUdaTkg2d1pzUDJ6eUsxNExaSEppNzc=`
- Client Authentication: On
- Root URL: `http://localhost:8080`

Registered users:

| Realm | Username | Password | Email | Email verified | Temporal password | First name | Last name |
| - | - | - | - | - | - | - | - |
| aileron | test1 | password1 | <test1@example.com> | true | false | foo1 | bar1 |
| aileron | test2 | password2 | <test2@example.com> | true | true | foo2 | bar2 |

## JWT Keys

Private and public keys are registered in the aileron realm.

- [./keys/private.pem](./keys/private.pem)
- [./keys/public.pem](./keys/public.pem)

It is registered in the aileron realm as

- Kid: `QJQdUdHaY_OXC8BfO-3tqVV0s64nvrFSffyfqONNeYk`.
- Type: `RSA`
- Algorithm: `RS256`

Keys are generated with following command.

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem
```

## Export realm

```bash
docker exec -it keycloak bash
cd  /opt/keycloak/bin/
./kc.sh export --dir /opt/keycloak/data/import --realm aileron
```

See [Importing and exporting realms](https://www.keycloak.org/server/importExport).
