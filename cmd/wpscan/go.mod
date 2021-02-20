module github.com/CyberAgent/mimosa-diagnosis/cmd/wpscan

go 1.16

require (
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20201130105221-b9659eb5f70a
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20201130105221-b9659eb5f70a
	github.com/CyberAgent/mimosa-diagnosis/pkg/common v0.0.0-20201203074646-21680d2e8ea0
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20201203074646-21680d2e8ea0
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20201203074646-21680d2e8ea0
	github.com/aws/aws-sdk-go v1.37.10
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/h2ik/go-sqs-poller/v3 v3.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sys v0.0.0-20201211002650-1f0c578a6b29 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201210142538-e3217bee35cc // indirect
	google.golang.org/grpc v1.34.0
)
