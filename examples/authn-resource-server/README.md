# OAuth2 ResourceServer

## Overview

This example shows response compression.

```mermaid
block-beta
  columns 7
  Downstream:1
  space:1
  block:aileron:3
    HTTPServer["🟪</br>HTTP</br>Server"]
    CompressionMiddleware["🟩</br>Compression</br>Middleware"]
    ReverseProxyHandler["🟥</br>ReverseProxy</br>Handler"]
  end
  space:1
  Upstream:1

Downstream --> HTTPServer
HTTPServer --> Downstream
Upstream --> ReverseProxyHandler
ReverseProxyHandler --> Upstream

style Downstream stroke:#888
style Upstream stroke:#888
style ReverseProxyHandler stroke:#ff6961,stroke-width:2px
style CompressionMiddleware stroke:#77dd77,stroke-width:2px
```

**Legend**:

- 🟥 `#ff6961` Handler resources.
- 🟩 `#77dd77` Middleware resources (Server-side middleware).
- 🟦 `#89CFF0` Tripperware resources (Client-side middleware).
- 🟪 `#9370DB` Other resources.

In this example, following directory structure and files are supposed.

Example resources are available at [examples/authn-resource-server/]({{% github-url "" %}}).
If you need a pre-built binary, download from [GitHub Releases](https://github.com/aileron-gateway/aileron-gateway/releases).

```txt
authn-resource-server/   ----- Working directory.
├── aileron              ----- AILERON Gateway binary (aileron.exe on windows).
├── config.yaml          ----- AILERON Gateway config file.
└── docker-compose.yaml  ----- Docker compose file to run
```

## Config

Configuration yaml to run a server with access logging becomes as follows.

```yaml
# config.yaml

{{% github-raw "config.yaml" %}}
```

The config tells:

- Start a `HTTPServer` with port 8080.
- A ReverseProxy handler is applied.
- Response compression is applied with CompressionMiddleware.
  - Minimum compression size is 10 bytes.

This graph shows the resource dependencies of the configuration.

```mermaid
graph TD
  Entrypoint["🟪 **Entrypoint**</br>default/default"]
  HTTPServer["🟪 **HTTPServer**</br>default/default"]
  ReverseProxyHandler["🟥 **ReverseProxyHandler**</br>default/default"]
  CompressionMiddleware["🟩 **CompressionMiddleware**</br>default/default"]

Entrypoint --"Runner"--> HTTPServer
HTTPServer --"HTTP Handler"--> ReverseProxyHandler
HTTPServer --"Middleware"--> CompressionMiddleware

style ReverseProxyHandler stroke:#ff6961,stroke-width:2px
style CompressionMiddleware stroke:#77dd77,stroke-width:2px
```

## Run

First, start keycloak for testing.

```bash
docker compose up
```

Run the AILERON Gateway with the command:

```bash
./aileron -f ./config.yaml
```

## Check

First, we try to send a HTTP request without any tokens.
An unauthorized response will be returned for such a request as shown below.

```json
$ curl http://localhost:8080/test
{"status":401,"statusText":"Unauthorized"}
```

Next, we try to send a HTTP request with a valid token.
The token can be obtained from token endpoint of the keycloak.

This is an example to get a token from the keycloak.
ROPC (Resource Owner Password Credentials Flow) is used here.

```bash
curl -H "Authorization: Basic YWlsZXJvbl9yb3BjOmMzbEFwVmt2aExucFowMVhCeHNoY0p6WTlGR3NxMzBx" -d "username=test&password=password&grant_type=password&scope=openid" http://localhost:18080/realms/aileron/protocol/openid-connect/token
```

The keycloak will respond a json like this (output is formatted for readability).
We can find `access_token` in the response that we wanted to get.

```json
{
   "access_token":"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ4RnptNEJZZ2ktQTB1ek85UHBqZ1htU1M5SHJzTFpQaGMwSXJIOGhlZFg0In0.eyJleHAiOjE3MDkyMTc5ODcsImlhdCI6MTcwOTIxNzY4NywianRpIjoiOWQ4NzUwZjItNDA0OC00N2ZmLTk1ZTQtNDIwYmY0MzI3ZmYxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxODA4MC9yZWFsbXMvYWlsZXJvbiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJlZTMwNjlmNy00ZGMzLTRmZmEtYjI5MC02YjZjOGE4NDg1MzciLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhaWxlcm9uX3JvcGMiLCJzZXNzaW9uX3N0YXRlIjoiMzljYzU2MzItNzAzZi00YzA5LTgwZWYtMjFmMDgwN2Y3MTBlIiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1haWxlcm9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjM5Y2M1NjMyLTcwM2YtNGMwOS04MGVmLTIxZjA4MDdmNzEwZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiZm9vIGJhciIsInByZWZlcnJlZF91c2VybmFtZSI6InRlc3QiLCJnaXZlbl9uYW1lIjoiZm9vIiwiZmFtaWx5X25hbWUiOiJiYXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.hiaSdU7Mi4sAWL_cQhVqze0imJK04N8DveNTig9K5pci-2LtZgp9VhgQC4hFoUD_1CzKPWZbbCcDZ4c8Ce9fZtEZOgG1xvXeC9nf62fPLtaf_RWjB4QQKpuNHJ1500cuhg0aJR9iNlvHXplfmGo9GHXZnZD6p6kV2rgMBOzdIveCDgzBELkMhn5-Z3NoFFvExjAEbjsIuFW8GnKzhA5lJ7Izalgcw9ziqP2u3WhmEfgQY9Z_claMUPEnxQMsrOszsrcJFBkmFWZ4tdJBvME_aV-xebQc8XHA6y-k42adSq49vR7_-ok9jmwgw-D_UP5U4CfxnDvfiL3I2goavqw6_w",
   "expires_in":300,
   "refresh_expires_in":1800,
   "refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzNGRkYjg2MS1hYTlkLTQ3MzEtYjU5Ni1kNjExYzNjNDAwMzUifQ.eyJleHAiOjE3MDkyMTk0ODcsImlhdCI6MTcwOTIxNzY4NywianRpIjoiZjhiN2JiMDYtOWVjMC00ZmE0LTgzZjctYzkxNTBjN2Q0NGVjIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxODA4MC9yZWFsbXMvYWlsZXJvbiIsImF1ZCI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTgwODAvcmVhbG1zL2FpbGVyb24iLCJzdWIiOiJlZTMwNjlmNy00ZGMzLTRmZmEtYjI5MC02YjZjOGE4NDg1MzciLCJ0eXAiOiJSZWZyZXNoIiwiYXpwIjoiYWlsZXJvbl9yb3BjIiwic2Vzc2lvbl9zdGF0ZSI6IjM5Y2M1NjMyLTcwM2YtNGMwOS04MGVmLTIxZjA4MDdmNzEwZSIsInNjb3BlIjoib3BlbmlkIGVtYWlsIHByb2ZpbGUiLCJzaWQiOiIzOWNjNTYzMi03MDNmLTRjMDktODBlZi0yMWYwODA3ZjcxMGUifQ.kbSSiY6iZ7RP3enSEAsgXNA07Mjmgsfd04P2GZ8L6H0",
   "token_type":"Bearer",
   "id_token":"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ4RnptNEJZZ2ktQTB1ek85UHBqZ1htU1M5SHJzTFpQaGMwSXJIOGhlZFg0In0.eyJleHAiOjE3MDkyMTc5ODcsImlhdCI6MTcwOTIxNzY4NywiYXV0aF90aW1lIjowLCJqdGkiOiIxY2FjYWFhMC02M2JmLTQxMzEtYjQ1My1hNjY3OWI4ZDIxOTYiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjE4MDgwL3JlYWxtcy9haWxlcm9uIiwiYXVkIjoiYWlsZXJvbl9yb3BjIiwic3ViIjoiZWUzMDY5ZjctNGRjMy00ZmZhLWIyOTAtNmI2YzhhODQ4NTM3IiwidHlwIjoiSUQiLCJhenAiOiJhaWxlcm9uX3JvcGMiLCJzZXNzaW9uX3N0YXRlIjoiMzljYzU2MzItNzAzZi00YzA5LTgwZWYtMjFmMDgwN2Y3MTBlIiwiYXRfaGFzaCI6IklVZlQwc0tjbXEydTdKSW0ycXBIQnciLCJhY3IiOiIxIiwic2lkIjoiMzljYzU2MzItNzAzZi00YzA5LTgwZWYtMjFmMDgwN2Y3MTBlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJmb28gYmFyIiwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdCIsImdpdmVuX25hbWUiOiJmb28iLCJmYW1pbHlfbmFtZSI6ImJhciIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.JfpI-3Vkpb9wisJ1F5ieyA56jmbKRyJv4b-PhQFDSLa7pn6TWqszYE5k-qS9k-RjHLI6eGInBlT9-c87WTJ25F5BRKzf1vWsgO35GT05hFYoqnVMLGA6pUxcS9zrtRsFf24-NUNg2OHYFEdVFrnQIz5hh72T9vamK1pYFa_pVuv8hcTAtN64qlSEuC3iIupH1_U94cg6bLzM2_NkJ9s-sRDRL08a3_sKdE36osUAWvz-ozYBWRHBtXwEUPVx1SxBxwKTdSyulnbMCVcuwLhfpLpN03Fs7LanxzGQcIG9ZDE2Sg4pXdDu_zdIZ5PgSorg4lpAqHI34TYMpyvMWVQOkQ",
   "not-before-policy":1684309880,
   "session_state":"39cc5632-703f-4c09-80ef-21f0807f710e",
   "scope":"openid email profile"
}
```

Now, send a HTTP request with the access token.

```bash
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJid2htUXhRa3d2cnVJaUE1akRkVlYyYlBaWkw1WGd6SUZpR01FQXlmRUtFIn0.eyJleHAiOjE3NDc1NDY0OTksImlhdCI6MTc0NzU0NjE5OSwianRpIjoib25ydHJvOmQ4NjgwNDY5LTczNjQtNDA2Zi1hZjIxLTQ4OTUxZjZkNmQ5OSIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTgwODAvcmVhbG1zL2FpbGVyb24iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiOTg5NzVlNGQtNDc1MS00MzYwLTk3YmQtYzdiZTM4OWVlZmQyIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYWlsZXJvbl9yZXNvdXJjZV9zZXJ2ZXIiLCJzaWQiOiI3YzkzYjVkZi1jYzExLTRkYWItODY1NC1kZmI0MmVjOWZjNzAiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIi8qIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1haWxlcm9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiZm9vMSBiYXIxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdDEiLCJnaXZlbl9uYW1lIjoiZm9vMSIsImZhbWlseV9uYW1lIjoiYmFyMSIsImVtYWlsIjoidGVzdDFAZXhhbXBsZS5jb20ifQ.mnuduiDsbnY3YVg1C3SHYYC0_esddBq_51fjj0iDEtzqTrXLZAMeY_3_xIQjF4S3qkB1gv4c0fHDvShCNCrPstYX7DrWnuAHHf2K7QBftJi9Z11XZRmRaci9eWw4tMdFjCqgsAEMZUWyr6-cwZXnvPeopF_rZbkGQUbQtgPFZWxZLgPgMLyB6ACVQzq4g0y_xsOtsQkWcJqSz82JH49bsaF_b8oslY1W5mPXAtRrBWJ8OPo6Qnh87wDv33JL4KGhxC1ZC7nHc_sQhL5QT8YyIMfRuxXkIygrslkF34bSA-UUW3VNZ4oRuSuTOt9tW0oFSEqxJsVslrei_X-wQTp7gw
" http://localhost:8080/test
```

The request will be authorized and non-error response will be responded.
The following output is an example.

```text
========== Request ==========
Method : GET
URI    : /test
Path   : /test
Query  :
Remote : [::1]:41508
========== Header ==========
{
  "Accept": [
    "*/*"
  ],
  "Authorization": [
    "Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ4RnptNEJZZ2ktQTB1ek85UHBqZ1htU1M5SHJzTFpQaGMwSXJIOGhlZFg0In0.eyJleHAiOjE3MDkyMTc5ODcsImlhdCI6MTcwOTIxNzY4NywianRpIjoiOWQ4NzUwZjItNDA0OC00N2ZmLTk1ZTQtNDIwYmY0MzI3ZmYxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxODA4MC9yZWFsbXMvYWlsZXJvbiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJlZTMwNjlmNy00ZGMzLTRmZmEtYjI5MC02YjZjOGE4NDg1MzciLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhaWxlcm9uX3JvcGMiLCJzZXNzaW9uX3N0YXRlIjoiMzljYzU2MzItNzAzZi00YzA5LTgwZWYtMjFmMDgwN2Y3MTBlIiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1haWxlcm9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjM5Y2M1NjMyLTcwM2YtNGMwOS04MGVmLTIxZjA4MDdmNzEwZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiZm9vIGJhciIsInByZWZlcnJlZF91c2VybmFtZSI6InRlc3QiLCJnaXZlbl9uYW1lIjoiZm9vIiwiZmFtaWx5X25hbWUiOiJiYXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.hiaSdU7Mi4sAWL_cQhVqze0imJK04N8DveNTig9K5pci-2LtZgp9VhgQC4hFoUD_1CzKPWZbbCcDZ4c8Ce9fZtEZOgG1xvXeC9nf62fPLtaf_RWjB4QQKpuNHJ1500cuhg0aJR9iNlvHXplfmGo9GHXZnZD6p6kV2rgMBOzdIveCDgzBELkMhn5-Z3NoFFvExjAEbjsIuFW8GnKzhA5lJ7Izalgcw9ziqP2u3WhmEfgQY9Z_claMUPEnxQMsrOszsrcJFBkmFWZ4tdJBvME_aV-xebQc8XHA6y-k42adSq49vR7_-ok9jmwgw-D_UP5U4CfxnDvfiL3I2goavqw6_w"
  ],
  "User-Agent": [
    "curl/7.68.0"
  ]
}
========== Body ==========

==========================
```

curl -H "Authorization: Basic YWlsZXJvbl9yZXNvdXJjZV9zZXJ2ZXI6S1dZUEJnclRFRUdaTkg2d1pzUDJ6eUsxNExaSEppNzc=" -d "username=test1&password=password1&grant_type=password&scope=openid" http://localhost:18080/realms/aileron/protocol/openid-connect/token
