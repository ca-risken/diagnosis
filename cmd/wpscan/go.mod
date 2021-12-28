module github.com/ca-risken/diagnosis/cmd/wpscan

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/aws/aws-sdk-go v1.42.22
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/sqs v0.0.0-20211210074045-79fdb4c61950
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20211014091419-ba2a979c7659
	github.com/ca-risken/core/proto/finding v0.0.0-20211227095930-ef25be878432
	github.com/ca-risken/diagnosis/pkg/common v0.0.0-20211014145120-f1682296ef05
	github.com/ca-risken/diagnosis/pkg/message v0.0.0-20211014145120-f1682296ef05
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20211014145120-f1682296ef05
	github.com/gassara-kys/envconfig v1.4.4
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20211014222326-fd004c51d1d6 // indirect
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
	google.golang.org/grpc v1.41.0
)
