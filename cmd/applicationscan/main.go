package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "diagnosis"
	serviceName = "applicationscan"
	settingURL  = "https://docs.security-hub.jp/diagnosis/applicationscan_datasource/"
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
	CoreAddr         string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DiagnosisSvcAddr string `required:"true" split_words:"true" default:"diagnosis.diagnosis.svc.cluster.local:19001"`

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
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	handler := &sqsHandler{
		zapPort:         conf.ZapPort,
		zapPath:         conf.ZapPath,
		zapApiKeyName:   conf.ZapApiKeyName,
		zapApiKeyHeader: conf.ZapApiKeyHeader,
	}
	handler.findingClient = newFindingClient(conf.CoreAddr)
	appLogger.Info(ctx, "Start Finding Client")
	handler.alertClient = newAlertClient(conf.CoreAddr)
	appLogger.Info(ctx, "Start Alert Client")
	handler.diagnosisClient = newDiagnosisClient(conf.DiagnosisSvcAddr)
	appLogger.Info(ctx, "Start Diagnosis Client")
	f, err := mimosasqs.NewFinalizer(common.DataSourceNameApplicationScan, settingURL, conf.CoreAddr, nil)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &SQSConfig{
		Debug:                             conf.Debug,
		AWSRegion:                         conf.AWSRegion,
		SQSEndpoint:                       conf.SQSEndpoint,
		DiagnosisApplicationScanQueueName: conf.DiagnosisApplicationScanQueueName,
		DiagnosisApplicationScanQueueURL:  conf.DiagnosisApplicationScanQueueURL,
		MaxNumberOfMessage:                conf.MaxNumberOfMessage,
		WaitTimeSecond:                    conf.WaitTimeSecond,
	}
	consumer := newSQSConsumer(ctx, sqsConf)
	appLogger.Info(ctx, "Start the ApplicationScan SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.TracingHandler(getFullServiceName(),
					mimosasqs.StatusLoggingHandler(appLogger,
						f.FinalizeHandler(handler))))))
}
