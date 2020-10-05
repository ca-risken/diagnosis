module github.com/CyberAgent/mimosa-diagnosis/cmd/diagnosis

go 1.15

require (
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20200907104350-788fdb7c85b1
	github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis v0.0.0-20200907104350-788fdb7c85b1
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.33.19
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.15
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.uber.org/zap v1.15.0
	google.golang.org/grpc v1.31.1
)
