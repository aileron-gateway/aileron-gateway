module github.com/aileron-gateway/aileron-gateway

go 1.23.0

toolchain go1.23.4

godebug default=go1.23

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.4-20241127180247-a33202765966.1
	github.com/andybalholm/brotli v1.1.1
	github.com/bufbuild/protovalidate-go v0.8.2
	github.com/casbin/casbin/v3 v3.0.0-beta.7
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/lestrrat-go/jwx/v2 v2.1.3
	github.com/open-policy-agent/opa v1.0.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pion/dtls/v3 v3.0.4
	github.com/prometheus/client_golang v1.20.5
	github.com/quic-go/quic-go v0.49.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/spf13/pflag v1.0.5
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/vmihailenco/msgpack/v5 v5.4.1
	go.opencensus.io v0.24.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.59.0
	go.opentelemetry.io/contrib/propagators/autoprop v0.59.0
	go.opentelemetry.io/contrib/propagators/aws v1.34.0
	go.opentelemetry.io/contrib/propagators/b3 v1.34.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.34.0
	go.opentelemetry.io/contrib/propagators/opencensus v0.59.0
	go.opentelemetry.io/contrib/propagators/ot v1.34.0
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.34.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.34.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.34.0
	go.opentelemetry.io/otel/exporters/zipkin v1.34.0
	go.opentelemetry.io/otel/metric v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.opentelemetry.io/otel/sdk/metric v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
	go.opentelemetry.io/proto/otlp v1.5.0
	golang.org/x/crypto v0.32.0
	golang.org/x/net v0.34.0
	google.golang.org/grpc v1.69.4
	google.golang.org/grpc/examples v0.0.0-20240821223602-0a5b8f7c9b41
	google.golang.org/protobuf v1.36.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cel.dev/expr v0.19.1 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.2.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.5 // indirect
	github.com/dgraph-io/ristretto v0.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cel-go v0.22.1 // indirect
	github.com/google/flatbuffers v24.12.23+incompatible // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/ginkgo/v2 v2.22.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/transport/v3 v3.0.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.61.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/bridge/opencensus v1.34.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
