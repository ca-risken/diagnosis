module github.com/CyberAgent/mimosa-diagnosis/cmd/jira

go 1.15

require (
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20201014033227-7bd4a4b0822e
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20201013111037-bf1eb2b9314c
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20201013111037-bf1eb2b9314c
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.35.7
	github.com/go-sql-driver/mysql v1.5.0
	github.com/h2ik/go-sqs-poller/v3 v3.0.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb // indirect
	golang.org/x/sys v0.0.0-20201013132646-2da7054afaeb // indirect
	google.golang.org/genproto v0.0.0-20201013134114-7f9ee70cb474 // indirect
	google.golang.org/grpc v1.33.0
)
