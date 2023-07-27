package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/applicationscan"
	"github.com/ca-risken/diagnosis/pkg/grpc"
	"github.com/ca-risken/diagnosis/pkg/sqs"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "diagnosis"
	serviceName = "applicationscan"
	settingURL  = "https://docs.security-hub.jp/diagnosis/applicationscan_datasource/"
)

var (
	appLogger            = logging.NewLogger()
	samplingRate float64 = 0.3000
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// sqs
	Debug string `default:"false"`

	AWSRegion   string `envconfig:"aws_region"   default:"ap-northeast-1"`
	SQSEndpoint string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	DiagnosisApplicationScanQueueName string `split_words:"true" default:"diagnosis-applicationscan"`
	DiagnosisApplicationScanQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-applicationscan"`
	MaxNumberOfMessage                int32  `split_words:"true" default:"1"`
	WaitTimeSecond                    int32  `split_words:"true" default:"20"`

	// grpc
	CoreAddr             string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DataSourceAPISvcAddr string `required:"true" split_words:"true" default:"datasource-api.datasource.svc.cluster.local:8081"`

	// zap
	ZapPort         string `split_words:"true" default:"8080"`
	ZapPath         string `split_words:"true" default:"/zap/zap.sh"`
	ZapApiKeyName   string `split_words:"true" default:"apikey"`
	ZapApiKeyHeader string `split_words:"true" default:"X-ZAP-API-Key"`
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
		ServiceName:  getFullServiceName(),
		Environment:  conf.EnvName,
		Debug:        conf.TraceDebug,
		SamplingRate: &samplingRate,
	}
	tracer.Start(tc)
	defer tracer.Stop()

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
	appc, err := applicationscan.NewApplicationScanClient(conf.ZapPort, conf.ZapPath, conf.ZapApiKeyName, conf.ZapApiKeyHeader, appLogger)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create diagnosis client, err=%+v", err)
	}
	appLogger.Info(ctx, "Start ApplicationScan Client")
	handler := applicationscan.NewSqsHandler(fc, ac, dc, appc, appLogger)

	f, err := mimosasqs.NewFinalizer(message.DataSourceNameApplicationScan, settingURL, conf.CoreAddr, nil)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &sqs.SQSConfig{
		Debug:              conf.Debug,
		AWSRegion:          conf.AWSRegion,
		SQSEndpoint:        conf.SQSEndpoint,
		QueueName:          conf.DiagnosisApplicationScanQueueName,
		QueueURL:           conf.DiagnosisApplicationScanQueueURL,
		MaxNumberOfMessage: conf.MaxNumberOfMessage,
		WaitTimeSecond:     conf.WaitTimeSecond,
	}
	consumer, err := sqs.NewSQSConsumer(ctx, sqsConf, appLogger)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create SQS consumer, err=%+v", err)
	}
	appLogger.Info(ctx, "Start the ApplicationScan SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.TracingHandler(getFullServiceName(),
					mimosasqs.StatusLoggingHandler(appLogger,
						f.FinalizeHandler(handler))))))
}
