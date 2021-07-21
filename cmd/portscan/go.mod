module github.com/CyberAgent/mimosa-aws/src/portscan

go 1.15

require (
	github.com/CyberAgent/mimosa-common/pkg/logging v0.0.0-20210713075629-8fad8ef7892d
	github.com/CyberAgent/mimosa-common/pkg/portscan v0.0.0-20210713075629-8fad8ef7892d
	github.com/CyberAgent/mimosa-common/pkg/xray v0.0.0-20210721063343-44cefe7f590e
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20210712081026-7152ed72951d
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20210712081026-7152ed72951d
	github.com/CyberAgent/mimosa-diagnosis/pkg/common v0.0.0-20210714042048-32f6643f79bf
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20210714042048-32f6643f79bf
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20210714042048-32f6643f79bf
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.39.6
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210714021259-044028024a4f // indirect
	google.golang.org/grpc v1.39.0
)
