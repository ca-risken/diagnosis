module github.com/ca-risken/diagnosis/cmd/diagnosis

go 1.17

require (
	github.com/aws/aws-sdk-go v1.41.3
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/database v0.0.0-20211223025030-6bfdc45e906c
	github.com/ca-risken/common/pkg/logging v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/rpc v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/diagnosis/pkg/message v0.0.0-20211014145120-f1682296ef05
	github.com/ca-risken/diagnosis/pkg/model v0.0.0-20211014145120-f1682296ef05
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20211224114208-28900968b46b
	github.com/gassara-kys/envconfig v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	google.golang.org/grpc v1.42.0
	gorm.io/gorm v1.21.16
)

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/ca-risken/common/pkg/sqs v0.0.0-20220113015330-0e8462d52b5b // indirect
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.31.0 // indirect
	golang.org/x/net v0.0.0-20211014222326-fd004c51d1d6 // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211013025323-ce878158c4d4 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/driver/mysql v1.1.2 // indirect
)
