module github.com/CyberAgent/mimosa-diagnosis/cmd/jira

go 1.15

require (
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20201028054340-f4dee5a77a75
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20201028054340-f4dee5a77a75
	github.com/CyberAgent/mimosa-diagnosis/pkg/common v0.0.0-20201102104958-713d5ff0f488
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20201102104958-713d5ff0f488
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20201102104958-713d5ff0f488
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.35.19
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/h2ik/go-sqs-poller/v3 v3.0.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102 // indirect
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201030142918-24207fddd1c3 // indirect
	google.golang.org/grpc v1.33.1
)
