module github.com/CyberAgent/mimosa-diagnosis/cmd/jira

go 1.15

require (
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20201005093216-9802eda8de17
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20200907104350-788fdb7c85b1
	github.com/aws/aws-sdk-go v1.33.20
	github.com/go-sql-driver/mysql v1.5.0
	github.com/h2ik/go-sqs-poller/v3 v3.0.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.6.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/zap v1.15.0
	google.golang.org/grpc v1.31.0
)
