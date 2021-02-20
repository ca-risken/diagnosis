module github.com/CyberAgent/mimosa-diagnosis/cmd/jira

go 1.16

require (
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20210218135922-8cb0b0d8730e
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20210218135922-8cb0b0d8730e
	github.com/CyberAgent/mimosa-diagnosis/pkg/common v0.0.0-20210220161919-c2a93d71069f
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20210220161919-c2a93d71069f
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20210220161919-c2a93d71069f
	github.com/aws/aws-sdk-go v1.37.15
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.11.0 // indirect
	github.com/sirupsen/logrus v1.8.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210220033124-5f55cee0dc0d // indirect
	golang.org/x/sys v0.0.0-20210220050731-9a76102bfb43 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210219173056-d891e3cb3b5b // indirect
	google.golang.org/grpc v1.35.0
)
