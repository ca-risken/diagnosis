module github.com/ca-risken/diagnosis/cmd/wpscan

go 1.17

require (
	github.com/aws/aws-sdk-go v1.43.15
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20220304031727-c94e2c463b27
	github.com/ca-risken/common/pkg/profiler v0.0.0-20220304031727-c94e2c463b27
	github.com/ca-risken/common/pkg/sqs v0.0.0-20220304031727-c94e2c463b27
	github.com/ca-risken/common/pkg/xray v0.0.0-20220304031727-c94e2c463b27
	github.com/ca-risken/core/proto/alert v0.0.0-20220309052852-c058b4e5cb84
	github.com/ca-risken/core/proto/finding v0.0.0-20220406055419-c44f2ad61c00
	github.com/ca-risken/diagnosis/pkg/common v0.0.0-20220309052839-1394e8f2e4a1
	github.com/ca-risken/diagnosis/pkg/message v0.0.0-20220309052839-1394e8f2e4a1
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20220309052839-1394e8f2e4a1
	github.com/gassara-kys/envconfig v1.4.4
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/golang/protobuf v1.5.2
	github.com/stretchr/testify v1.7.1
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	google.golang.org/grpc v1.45.0
)

require (
	github.com/DataDog/datadog-go/v5 v5.1.0 // indirect
	github.com/DataDog/gostackparse v0.5.0 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/pprof v0.0.0-20220218203455-0368bd9e19a7 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.34.0 // indirect
	golang.org/x/net v0.0.0-20220403103023-749bd193bc2b // indirect
	golang.org/x/sys v0.0.0-20220406163625-3f8b81556e12 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220405205423-9d709892a2bf // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.36.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
