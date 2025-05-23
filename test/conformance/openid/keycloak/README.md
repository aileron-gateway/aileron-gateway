# Keycloak

## Overview

This folder contains docker-compose for keycloak that are used from examples.

## Run the instance

Instances can be run in two different modes: `mTLS disabled` and `mTLS enabled`.

|  | mTLS disabled | mTLS enabled |
| - | - | - |
| Default port | `18080` | `18443` |
| Default protocol | HTTP | mTLS |
| Run Keycloak | `docker-compose up` | `docker-compose up -f docker-compose-mtls.yml` |
| Stop Keycloak | `docker-compose down` | `docker-compose down` |

## Accessing the Keycloak

If started with `mTLS disabled`, Admin and well-known endpoints are

- `Admin console`: [http://localhost:18080/admin/](http://localhost:18080/admin/)
- `Well-Known URL`: [http://localhost:18080/realms/aileron/.well-known/openid-configuration](http://localhost:18080/realms/aileron/.well-known/openid-configuration)

If started with `mTLS enabled`, Admin and well-known endpoints are

- `Admin console`: [https://localhost:18443/admin/](https://localhost:18443/admin/)
- `Well-Known URL`: [https://localhost:18443/realms/aileron/.well-known/openid-configuration](https://localhost:18443/realms/aileron/.well-known/openid-configuration)

For more details, please refer to the keycloak official documentation such as [getting-started-docker](https://www.keycloak.org/getting-started/getting-started-docker).

Some OAuth clients and uses are also loaded in the instance.

**Clients.**

| Realm | Client ID | Client PW | Base64(ID:PW) | mTLS enabled |
| - | - | - | - | - |
| aileron | aileron_authorization_code | NNJxpasHJgwGmgTYMoVFTKGUWZxp12bl | YWlsZXJvbl9hdXRob3JpemF0aW9uX2NvZGU6Tk5KeHBhc0hKZ3dHbWdUWU1vVkZUS0dVV1p4cDEyYmw= |  |
| aileron | aileron_client_credentials | j9wEw6Zj3dhGIhVzGCD1jpsDLDkMl8wD | YWlsZXJvbl9jbGllbnRfY3JlZGVudGlhbHM6ajl3RXc2WmozZGhHSWhWekdDRDFqcHNETERrTWw4d0Q= |  |
| aileron | aileron_resource_server | 2Ry7fU4a71TEbv3C1vBkMiJvHFUh4jzq | YWlsZXJvbl9yZXNvdXJjZV9zZXJ2ZXI6MlJ5N2ZVNGE3MVRFYnYzQzF2QmtNaUp2SEZVaDRqenE= |  |
| aileron | aileron_ropc | c3lApVkvhLnpZ01XBxshcJzY9FGsq30q | YWlsZXJvbl9yb3BjOmMzbEFwVmt2aExucFowMVhCeHNoY0p6WTlGR3NxMzBx |  |
| aileron | client-secret-jwt | MkHyKDoZRM0r5hUN3me5fVJQWkWrRQc9 | - |  |
| aileron | private-key-jwt | - | - |  |
| aileron | self-signed-tls-client-auth | - | - | ✅ |
| aileron | jarm | 5xgbiDOE9qGvKIJa0XzKhTr1vsom4ZXF | amFybTo1eGdiaURPRTlxR3ZLSUphMFh6S2hUcjF2c29tNFpYRg== |  |
| aileron | request-object | mFofYoEehKmF67qJpSGZmHfZ5ez0BGTr | cmVxdWVzdC1vYmplY3Q6bUZvZllvRWVoS21GNjdxSnBTR1ptSGZaNWV6MEJHVHI= |  |
| aileron | certificate-bound-access-tokens | - | - | ✅ |



**Uses.**

| Realm | Username | Password | email | email verified | First name | Last name |
| - | - | - | - | - | - | - |
| admin | admin | password | <admin@example.com> | true | - | - |
| aileron | test | password | <test@example.com> | true | foo | bar |
| aileron | test1 | password1 | <test1@example.com> | true | foo1 | bar1 |
| aileron | test2 | password2 | <test2@example.com> | true | foo2 | bar2 |
| aileron | test3 | password3 | <test3@example.com> | false | foo3 | bar3 |

**Additional scopes.**

- management
- read
- write

**Signing keys.**

- `RSA Private Key`: [./keys/private.key](./keys/private.key)
- `RSA Public Key`: [./keys/public.key](./keys/public.key)

## Export configuration

To export configuration, run the following commands.

```bash
docker exec -it keycloak bash
cd  /opt/keycloak/bin/
./kc.sh export --dir /opt/keycloak/data/import
```

## Appendix

### Client Credentials with curl

```bash
curl -k --request POST \
--header "Authorization: Basic YWlsZXJvbl9jbGllbnRfY3JlZGVudGlhbHM6ajl3RXc2WmozZGhHSWhWekdDRDFqcHNETERrTWw4d0Q=" \
--data "grant_type=client_credentials&scope=openid" \
"http://localhost:18080/realms/aileron/protocol/openid-connect/token"
```

### ROPC with curl

```bash
curl -k --request POST \
--header "Authorization: Basic YWlsZXJvbl9yb3BjOmMzbEFwVmt2aExucFowMVhCeHNoY0p6WTlGR3NxMzBx" \
--data "username=test1&password=password1&grant_type=password&scope=openid" \
"http://localhost:18080/realms/aileron/protocol/openid-connect/token"
```
