module github.com/aileron-gateway/aileron-gateway

go 1.23.0

toolchain go1.23.4

godebug default=go1.23

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.1-20241127180247-a33202765966.1
	github.com/bufbuild/protovalidate-go v0.8.2
	github.com/google/go-cmp v0.6.0
	github.com/pion/dtls/v3 v3.0.4
	github.com/quic-go/quic-go v0.48.2
	github.com/spf13/pflag v1.0.5
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
	golang.org/x/crypto v0.31.0
	golang.org/x/net v0.33.0
	google.golang.org/grpc v1.67.1
	google.golang.org/grpc/examples v0.0.0-20240821223602-0a5b8f7c9b41
	google.golang.org/protobuf v1.36.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cel.dev/expr v0.19.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/cel-go v0.22.1 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/onsi/ginkgo/v2 v2.22.2 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/transport/v3 v3.0.7 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/exp v0.0.0-20250103183323-7d7fa50e5329 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250102185135-69823020774d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250102185135-69823020774d // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
