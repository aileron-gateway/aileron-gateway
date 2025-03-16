<div align="center">

# ![AILERON Gateway](./docs/logo.svg)

**AILERON Gateway is the secure and high-performance general purpose API Gateway for enterprise systems.**
Visit [our website](https://aileron-gateway.github.io/) to get more information.

[![GoDoc](https://godoc.org/github.com/aileron-gateway/aileron-gateway?status.svg)](http://godoc.org/github.com/aileron-gateway/aileron-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/aileron-gateway/aileron-gateway)](https://goreportcard.com/report/github.com/aileron-gateway/aileron-gateway)
[![License](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](./LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Faileron-gateway%2Faileron-gateway.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Faileron-gateway%2Faileron-gateway?ref=badge_shield&issueType=license)

[![Codecov](https://codecov.io/gh/aileron-gateway/aileron-gateway/branch/main/graph/badge.svg?token=L62XLZNFLE)](https://codecov.io/gh/aileron-gateway/aileron-gateway)
[![Test Suite](https://github.com/aileron-gateway/aileron-gateway/actions/workflows/test-suite.yaml/badge.svg?branch=main)](https://github.com/aileron-gateway/aileron-gateway/actions/workflows/test-suite.yaml?query=branch%3Amain)
[![OpenSourceInsight](https://badgen.net/badge/open%2Fsource%2F/insight/cyan)](https://deps.dev/go/github.com%2Faileron-gateway%2Faileron-gateway)
[![OSS Insight](https://badgen.net/badge/OSS/Insight/orange)](https://ossinsight.io/analyze/aileron-gateway/aileron-gateway)

</div>

## Key Features

- Networking
  - Routing: Path based, query value based, header value based, etc.
  - LoadBalancing: RoundRobin, Random, Maglev, RingHash, etc.
- Authentication/Authorization
  - [OAuth2](https://oauth.net/2/) / [OpenID Connect](https://openid.net/developers/how-connect-works/)
  - [Financial-grade API](https://openid.net/wg/fapi/)
  - [OpenPolicyAgent](https://www.openpolicyagent.org/)
  - [Casbin](https://casbin.org/)
- Observability for cloud native environment.
  - Metrics: [Prometheus](https://prometheus.io/)
  - Tracing: [OpenTelemetry](https://opentelemetry.io/)
  - Profiling: [Profiling](https://go.dev/blog/pprof) endpoint.

![features.svg](docs/features.svg)

## Install AILERON Gateway

> [!TIP]
> See the [Installation](https://aileron-gateway.github.io/docs/installation/) for full documentation.

Install options are:

1. **Pre-built binary**: Download from [Releases](https://github.com/aileron-gateway/aileron-gateway/releases).
2. **Linux package**: Download from [Releases](https://github.com/aileron-gateway/aileron-gateway/releases).
   1. RPM packages `*.rpm`
   2. Debian packages `*.deb`
   3. Alpine packages `*.apk`
   4. Arch Linux packages `*.pkg.tar.zst`
3. **Container image**: Pull from [GitHub Container Registry](https://github.com/aileron-gateway/aileron-gateway/pkgs/container/aileron-gateway%2Faileron).
4. **Build from source code**:
   1. Latest release: `go install github.com/aileron-gateway/aileron-gateway/cmd/aileron@latest`
   2. Main branch: `go install github.com/aileron-gateway/aileron-gateway/cmd/aileron@main`

## Getting Started

> [!TIP]
> See the [Getting Started](https://aileron-gateway.github.io/docs/getting-started/) and the [Tutorials](https://aileron-gateway.github.io/docs/tutorials/) for full documentation.

### Command line options

```bash
$ aileron --help

Options :
  -e, --env stringArray    env file path. each line be 'KEY=VALUE'
  -f, --file stringArray   config file or directory path. absolute or relative
  -h, --help               show help message
  -i, --info               show build information
  -o, --out string         template output format. yaml or json (default "yaml")
  -t, --template string    show template config. value format be 'Group/Version/Kind(/Namespace/Name)'
  -v, --version            show version
```

### Run a reverse proxy example

Example configs are available under [./_example/*](./_example/).
Here we start with a simple reverse proxy using [_example/reverse_proxy/](_example/reverse_proxy/).

Just run the AILERON Gateway with the config.
Config paths can be `relative` or `absolute`, `file path` or `directory path`, `*.yaml` or `*.yml` or `*.json`.
A reverse proxy server will listen on the port `:8080`.

```bash
$ aileron -f _example/reverse-proxy/

{"time":"2025-03-15 18:20:12","level":"INFO","msg":"server started. listening on [::]:8080","datetime":{"date":"2025-03-15","time":"18:20:12.556","zone":"Local"},"location":{"file":"httpserver/server.go","func":"httpserver.(*runner).Run","line":56}}
```

In this example, the reverse proxy server is pointing on the [https://httpbin.org/](https://httpbin.org/).
So, we can access to the [https://httpbin.org/](https://httpbin.org/) via [http://localhost:8080/](http://localhost:8080/) like below.

```bash
$ curl http://localhost:8080/get

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.68.0",
    "X-Amzn-Trace-Id": "Root=1-67d5479b-217d6f167e89a49e1c785584",
    "X-Forwarded-Host": "localhost:8080"
  },
  "origin": "127.0.0.1, 106.73.5.65",
  "url": "http://localhost:8080/get"
}
```

## Contributing

> [!IMPORTANT]
> Thank you for being interested in contributing to our community!!
>
> We apologize, but we could not merge any PRs until we make all legal problems clear.
> For example, preparing code of conduct, contributor license agreement and so on.
> We are now actively working on this hard as we can merge PRs soon as possible.
>
> Thank you for your understanding. March 15th, 2025.

See the [Contributing](https://aileron-gateway.github.io/community/contributing/) for full documentation.

## Support

See the [Support](https://aileron-gateway.github.io/community/support/) for full documentation.

Need enterprise support?

Contact us: `aileron-gateway@nri.co.jp`

## Governance

We are now actively working on this topic.

## License

Apache License 2.0

See the [LICENSE](./LICENSE) file.
