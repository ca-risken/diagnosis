module github.com/ca-risken/diagnosis/cmd/diagnosis

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/aws/aws-sdk-go v1.41.3
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/database v0.0.0-20211223025030-6bfdc45e906c
	github.com/ca-risken/common/pkg/rpc v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/diagnosis/pkg/message v0.0.0-20211014145120-f1682296ef05
	github.com/ca-risken/diagnosis/pkg/model v0.0.0-20211014145120-f1682296ef05
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20211014145120-f1682296ef05
	github.com/gassara-kys/envconfig v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20211014222326-fd004c51d1d6 // indirect
	google.golang.org/grpc v1.42.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/gorm v1.21.16
)
