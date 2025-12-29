# CHANGELOG v1.1.x <!-- omit in toc -->

**Table of contents**

- [Versions](#versions)
- [v1.1.1](#v111)
  - [Changes since v1.0.0](#changes-since-v100)
    - [New features](#new-features)
    - [Bug fix, Security fix](#bug-fix-security-fix)
    - [Other changes](#other-changes)
  - [Dependencies](#dependencies)
  - [Migration guides](#migration-guides)
- [v1.1.0](#v110)
  - [Changes since v1.0.4](#changes-since-v104)
    - [New features](#new-features-1)
    - [Breaking changes](#breaking-changes)
    - [Bug fix, Security fix](#bug-fix-security-fix-1)
    - [Other changes](#other-changes-1)
  - [Dependencies](#dependencies-1)
    - [Added](#added)
    - [Changed](#changed)
    - [Removed](#removed)
  - [Migration guides](#migration-guides-1)

## Versions

| AILERON Versoin | Go   | protoc | protoc-gen-go |
| --------------- | ---- | ------ | ------------- |
| v1.1.0          | 1.25 | v29.0  | v1.36.4       |

## v1.1.1

### Changes since v1.0.0

#### New features

_Nothing has changed._

#### Bug fix, Security fix

_Nothing has changed._

#### Other changes

- [#105](https://github.com/aileron-gateway/aileron-gateway/pull/105): Move kernel/testutil package to internal/testutil (@k7a-tomohiro)
- [#108](https://github.com/aileron-gateway/aileron-gateway/pull/108): Move kernel/hash, kernel/encrypt packages to internal/hash, internal/encrypt (@k7a-tomohiro)
- [#109](https://github.com/aileron-gateway/aileron-gateway/pull/109): Move kernel/txtutil package to internal/txtutil (@k7a-tomohiro)
- [#110](https://github.com/aileron-gateway/aileron-gateway/pull/110): Move kernel/kvs package to internal/kvs (@k7a-tomohiro)
- [#111](https://github.com/aileron-gateway/aileron-gateway/pull/111): Move kernel/network package to internal/network (@k7a-tomohiro)
- [#112](https://github.com/aileron-gateway/aileron-gateway/pull/112): Move util/security package to internal/security (@k7a-tomohiro)

### Dependencies

TODO: Fill

### Migration guides

_Migration is not required._

## v1.1.0

### Changes since v1.0.4

#### New features

_Nothing has changed._

#### Breaking changes

_Nothing has changed._

#### Bug fix, Security fix

_Nothing has changed._

#### Other changes

_Nothing has changed._

### Dependencies

#### Added

- buf.build/go/hyperpb: v0.1.3
- buf.build/go/protovalidate: v1.1.0
- github.com/brianvoe/gofakeit/v6: v6.28.0
- github.com/bytecodealliance/wasmtime-go/v39: v39.0.1
- github.com/containerd/containerd/v2: v2.2.0
- github.com/containerd/typeurl/v2: v2.2.3
- github.com/fatih/color: v1.15.0
- github.com/go-viper/mapstructure/v2: v2.4.0
- github.com/gogo/protobuf: v1.3.2
- github.com/hashicorp/golang-lru/v2: v2.0.7
- github.com/huandu/go-clone: v1.7.3
- github.com/huandu/go-sqlbuilder: v1.38.1
- github.com/huandu/xstrings: v1.4.0
- github.com/jordanlewis/gcassert: 389ef75
- github.com/lestrrat-go/dsig-secp256k1: v1.0.0
- github.com/lestrrat-go/dsig: v1.0.0
- github.com/lestrrat-go/httprc/v3: v3.0.2
- github.com/lestrrat-go/jwx/v3: v3.0.12
- github.com/lestrrat-go/option/v2: v2.0.0
- github.com/mattn/go-colorable: v0.1.13
- github.com/mattn/go-isatty: v0.0.19
- github.com/olekukonko/errors: v1.1.0
- github.com/olekukonko/ll: v0.0.9
- github.com/rivo/uniseg: v0.4.7
- github.com/rodaine/protogofakeit: v0.1.1
- github.com/timandy/routine: v1.1.6
- github.com/valyala/fastjson: v1.6.7
- github.com/vektah/gqlparser/v2: v2.5.31
- go.yaml.in/yaml/v2: v2.4.3
- go.yaml.in/yaml/v3: v3.0.4
- golang.org/x/tools/go/expect: v0.1.1-deprecated
- golang.org/x/tools/go/packages/packagestest: v0.1.1-deprecated
- gonum.org/v1/gonum: v0.16.0
- sigs.k8s.io/randfill: v1.0.0

#### Changed

- buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go: 8976f5b → 2a1774d
- cel.dev/expr: v0.24.0 → v0.25.1
- cloud.google.com/go/compute/metadata: v0.6.0 → v0.9.0
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: v1.26.0 → v1.30.0
- github.com/aileron-projects/go: v0.0.0-alpha.14 → v0.0.0-alpha.17
- github.com/alecthomas/units: b94a6e3 → 0f3dac3
- github.com/andybalholm/brotli: v1.1.1 → v1.2.0
- github.com/cenkalti/backoff/v5: v5.0.2 → v5.0.3
- github.com/cncf/xds/go: 2f00578 → 0feb691
- github.com/containerd/platforms: v0.2.1 → v1.0.0-rc.2
- github.com/cpuguy83/go-md2man/v2: v2.0.4 → v2.0.7
- github.com/dgraph-io/badger/v4: v4.7.0 → v4.8.0
- github.com/dgraph-io/ristretto/v2: v2.2.0 → v2.3.0
- github.com/envoyproxy/go-control-plane/envoy: v1.32.4 → v1.35.0
- github.com/envoyproxy/go-control-plane: v0.13.4 → 75eaa19
- github.com/fsnotify/fsnotify: v1.8.0 → v1.9.0
- github.com/go-jose/go-jose/v4: v4.0.4 → v4.1.3
- github.com/golang-jwt/jwt/v5: v5.2.2 → v5.3.0
- github.com/golang/glog: v1.2.4 → v1.2.5
- github.com/google/cel-go: v0.25.0 → v0.26.1
- github.com/google/flatbuffers: v25.2.10+incompatible → v25.9.23+incompatible
- github.com/google/pprof: c008609 → 94a9f03
- github.com/grpc-ecosystem/grpc-gateway/v2: v2.26.3 → v2.27.3
- github.com/klauspost/compress: v1.18.0 → v1.18.2
- github.com/lestrrat-go/blackmagic: v1.0.3 → v1.0.4
- github.com/mattn/go-runewidth: v0.0.9 → v0.0.16
- github.com/olekukonko/tablewriter: v0.0.5 → v1.1.0
- github.com/onsi/ginkgo/v2: v2.23.4 → v2.11.0
- github.com/onsi/gomega: v1.36.3 → v1.27.10
- github.com/open-policy-agent/opa: v1.1.0 → v1.11.0
- github.com/opencontainers/image-spec: v1.1.0 → v1.1.1
- github.com/pelletier/go-toml/v2: v2.1.0 → v2.2.4
- github.com/prometheus/client_golang: v1.22.0 → v1.23.2
- github.com/prometheus/common: v0.64.0 → v0.67.4
- github.com/prometheus/procfs: v0.16.1 → v0.19.2
- github.com/quic-go/qpack: v0.5.1 → v0.6.0
- github.com/quic-go/quic-go: v0.52.0 → v0.57.1
- github.com/redis/go-redis/v9: v9.9.0 → v9.17.2
- github.com/sagikazarmark/locafero: v0.4.0 → v0.11.0
- github.com/segmentio/asm: v1.2.0 → v1.2.1
- github.com/sergi/go-diff: v1.3.1 → v1.4.0
- github.com/sirupsen/logrus: v1.9.3 → dd1b4c2
- github.com/sourcegraph/conc: v0.3.0 → 5f936ab
- github.com/spf13/afero: v1.11.0 → v1.15.0
- github.com/spf13/cast: v1.6.0 → v1.10.0
- github.com/spf13/cobra: v1.9.1 → v1.10.1
- github.com/spf13/pflag: v1.0.6 → v1.0.10
- github.com/spf13/viper: v1.18.2 → v1.21.0
- github.com/spiffe/go-spiffe/v2: v2.5.0 → v2.6.0
- github.com/stoewer/go-strcase: v1.3.0 → v1.3.1
- github.com/stretchr/testify: v1.10.0 → v1.11.1
- github.com/tchap/go-patricia/v2: v2.3.2 → v2.3.3
- github.com/tidwall/match: v1.1.1 → v1.2.0
- go.opentelemetry.io/auto/sdk: v1.1.0 → v1.2.1
- go.opentelemetry.io/contrib/detectors/gcp: v1.34.0 → v1.38.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.59.0 → v0.63.0
- go.opentelemetry.io/contrib/instrumentation/runtime: v0.61.0 → v0.64.0
- go.opentelemetry.io/contrib/propagators/autoprop: v0.61.0 → v0.64.0
- go.opentelemetry.io/contrib/propagators/aws: v1.36.0 → v1.39.0
- go.opentelemetry.io/contrib/propagators/b3: v1.36.0 → v1.39.0
- go.opentelemetry.io/contrib/propagators/jaeger: v1.36.0 → v1.39.0
- go.opentelemetry.io/contrib/propagators/opencensus: v0.61.0 → v0.64.0
- go.opentelemetry.io/contrib/propagators/ot: v1.36.0 → v1.39.0
- go.opentelemetry.io/contrib/zpages: v0.60.0 → v0.62.0
- go.opentelemetry.io/otel/bridge/opencensus: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/stdout/stdoutmetric: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/stdout/stdouttrace: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/exporters/zipkin: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/metric: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/sdk/metric: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/sdk: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel/trace: v1.36.0 → v1.39.0
- go.opentelemetry.io/otel: v1.36.0 → v1.39.0
- go.opentelemetry.io/proto/otlp: v1.7.0 → v1.9.0
- go.uber.org/mock: v0.5.2 → v0.6.0
- golang.org/x/crypto: v0.38.0 → v0.46.0
- golang.org/x/exp: b6e5de4 → 8475f28
- golang.org/x/mod: v0.24.0 → v0.31.0
- golang.org/x/net: v0.40.0 → v0.48.0
- golang.org/x/oauth2: v0.30.0 → v0.32.0
- golang.org/x/sync: v0.14.0 → v0.19.0
- golang.org/x/sys: v0.33.0 → v0.39.0
- golang.org/x/term: v0.32.0 → v0.38.0
- golang.org/x/text: v0.25.0 → v0.32.0
- golang.org/x/time: v0.9.0 → v0.14.0
- golang.org/x/tools: v0.33.0 → v0.40.0
- google.golang.org/genproto/googleapis/api: 200df99 → 97cd9d5
- google.golang.org/genproto/googleapis/rpc: 200df99 → 97cd9d5
- google.golang.org/grpc: v1.72.2 → v1.77.0
- google.golang.org/protobuf: v1.36.6 → v1.36.11
- oras.land/oras-go/v2: v2.3.1 → v2.6.0
- sigs.k8s.io/yaml: v1.4.0 → v1.6.0

#### Removed

- github.com/AdaLogics/go-fuzz-headers: ced1acd
- github.com/OneOfOne/xxhash: v1.2.8
- github.com/bufbuild/protovalidate-go: v0.10.1
- github.com/bytecodealliance/wasmtime-go/v3: v3.0.2
- github.com/cenkalti/backoff/v4: v4.3.0
- github.com/chzyer/readline: v1.5.1
- github.com/containerd/containerd: v1.7.25
- github.com/francoispqt/gojay: v1.2.13
- github.com/go-task/slim-sprig/v3: v3.0.0
- github.com/gorilla/mux: v1.8.1
- github.com/hashicorp/hcl: v1.0.0
- github.com/ianlancetaylor/demangle: f615e6b
- github.com/magiconair/properties: v1.8.7
- github.com/mitchellh/mapstructure: v1.5.0
- github.com/pkg/errors: v0.9.1
- github.com/prashantv/gostub: v1.1.0
- github.com/sagikazarmark/slog-shim: v0.1.0
- github.com/zeebo/errs: v1.4.0
- go.uber.org/atomic: v1.9.0
- golang.org/x/telemetry: bda5523
- gopkg.in/ini.v1: v1.67.0
- gopkg.in/yaml.v2: v2.4.0

### Migration guides

No migration required except for updating package versions.

If you noticed something wrong with migration, please [create an issue](https://github.com/aileron-gateway/aileron-gateway/issues/new).
