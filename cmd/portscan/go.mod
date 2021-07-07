module github.com/CyberAgent/mimosa-aws/src/portscan

go 1.15

require (
	github.com/CyberAgent/mimosa-common/pkg/portscan v0.0.0-20210603021529-d703af2cd8c8
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20210430051641-778e594322c3
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20210430051641-778e594322c3
	github.com/CyberAgent/mimosa-diagnosis/pkg/common v0.0.0-20210630065525-78f621413430
	github.com/CyberAgent/mimosa-diagnosis/pkg/message v0.0.0-20210630065525-78f621413430
	github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis v0.0.0-20210630065525-78f621413430
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.38.54
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210603125802-9665404d3644 // indirect
	google.golang.org/genproto v0.0.0-20210603172842-58e84a565dcf // indirect
	google.golang.org/grpc v1.38.0
)
