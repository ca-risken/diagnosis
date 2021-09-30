module github.com/ca-risken/diagnosis/cmd/wpscan

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.40.52
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20210927112235-42730386bf2a
	github.com/ca-risken/common/pkg/sqs v0.0.0-20210927112235-42730386bf2a
	github.com/ca-risken/common/pkg/xray v0.0.0-20210927112235-42730386bf2a
	github.com/ca-risken/core/proto/alert v0.0.0-20210924100500-e1499111345b
	github.com/ca-risken/core/proto/finding v0.0.0-20210924100500-e1499111345b
	github.com/ca-risken/diagnosis/pkg/common v0.0.0-20210928110806-54123bffb7e9
	github.com/ca-risken/diagnosis/pkg/message v0.0.0-20210928110806-54123bffb7e9
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20210928110806-54123bffb7e9
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/valyala/fasthttp v1.30.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20210929193557-e81a3d93ecf6 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/genproto v0.0.0-20210929214142-896c89f843d2 // indirect
	google.golang.org/grpc v1.41.0
)
