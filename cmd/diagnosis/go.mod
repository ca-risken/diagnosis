module github.com/CyberAgent/mimosa-diagnosis/cmd/diagnosis

go 1.16

require (
	github.com/CyberAgent/mimosa-common/pkg/database v0.0.0-20210721063343-44cefe7f590e
	github.com/CyberAgent/mimosa-common/pkg/xray v0.0.0-20210803120909-2cc57e3c75d2
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20210714042048-32f6643f79bf
	github.com/CyberAgent/mimosa-diagnosis/pkg/model v0.0.0-20210714042048-32f6643f79bf
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20210714042048-32f6643f79bf
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.39.6
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.18.1
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210714021259-044028024a4f // indirect
	google.golang.org/grpc v1.39.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/gorm v1.21.12
)
