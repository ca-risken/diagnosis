package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/grpc"
	"github.com/ca-risken/diagnosis/pkg/sqs"
	"github.com/ca-risken/diagnosis/pkg/wpscan"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "diagnosis"
	serviceName = "wpscan"
	settingURL  = "https://docs.security-hub.jp/diagnosis/wpscan_datasource/"
)

var appLogger = logging.NewLogger()

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// sqs
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	DiagnosisWpscanQueueName string `split_words:"true" default:"diagnosis-wpscan"`
	DiagnosisWpscanQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-wpscan"`
	MaxNumberOfMessage       int32  `split_words:"true" default:"10"`
	WaitTimeSecond           int32  `split_words:"true" default:"20"`
	// grpc
	CoreAddr             string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DataSourceAPISvcAddr string `required:"true" split_words:"true" default:"datasource-api.core.svc.cluster.local:8081"`
	// wpscan
	ResultPath string `split_words:"true" required:"true" default:"/tmp"`
}

func main() {
	ctx := context.Background()
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	wc := wpscan.NewWpscanConfig(conf.ResultPath, appLogger)
	fc, err := grpc.NewFindingClient(conf.CoreAddr)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create finding client, err=%+v", err)
	}
	appLogger.Info(ctx, "Start Finding Client")
	ac, err := grpc.NewAlertClient(conf.CoreAddr)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create alert client, err=%+v", err)
	}
	appLogger.Info(ctx, "Start Alert Client")
	dc, err := grpc.NewDiagnosisClient(conf.DataSourceAPISvcAddr)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create diagnosis client, err=%+v", err)
	}
	appLogger.Info(ctx, "Start Diagnosis Client")
	handler := wpscan.NewSqsHandler(wc, fc, ac, dc, appLogger)

	f, err := mimosasqs.NewFinalizer(message.DataSourceNameWPScan, settingURL, conf.CoreAddr, &mimosasqs.DataSourceRecommnend{
		ScanFailureRisk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", message.DataSourceNameWPScan),
		ScanFailureRecommendation: fmt.Sprintf(`Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the network is reachable to the target host.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- %s
		- And please also check the FAQ page.
		- https://docs.security-hub.jp/contact/faq/#wpscan
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`, settingURL),
	})
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &sqs.SQSConfig{
		AWSRegion:          conf.AWSRegion,
		SQSEndpoint:        conf.Endpoint,
		QueueName:          conf.DiagnosisWpscanQueueName,
		QueueURL:           conf.DiagnosisWpscanQueueURL,
		MaxNumberOfMessage: conf.MaxNumberOfMessage,
		WaitTimeSecond:     conf.WaitTimeSecond,
	}
	consumer, err := sqs.NewSQSConsumer(ctx, sqsConf, appLogger)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create SQS consumer, err=%+v", err)
	}
	appLogger.Info(ctx, "Start the wpscan SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.TracingHandler(getFullServiceName(),
					mimosasqs.StatusLoggingHandler(appLogger,
						f.FinalizeHandler(handler))))))
}
