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
	serviceName = "portscan"
	settingURL  = "https://docs.security-hub.jp/diagnosis/portscan_datasource/"
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

	DiagnosisPortscanQueueName string `split_words:"true" default:"diagnosis-portscan"`
	DiagnosisPortscanQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-portscan"`
	MaxNumberOfMessage         int64  `split_words:"true" default:"5"`
	WaitTimeSecond             int64  `split_words:"true" default:"20"`

	// grpc
	FindingSvcAddr   string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
	AlertSvcAddr     string `required:"true" split_words:"true" default:"alert.core.svc.cluster.local:8004"`
	DiagnosisSvcAddr string `required:"true" split_words:"true" default:"diagnosis.diagnosis.svc.cluster.local:19001"`
}

func main() {
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	ctx := context.Background()

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	handler := &sqsHandler{}
	handler.findingClient = newFindingClient(conf.FindingSvcAddr)
	handler.alertClient = newAlertClient(conf.AlertSvcAddr)
	handler.diagnosisClient = newDiagnosisClient(conf.DiagnosisSvcAddr)
	f, err := mimosasqs.NewFinalizer(common.DataSourceNamePortScan, settingURL, conf.FindingSvcAddr, nil)
	if err != nil {
		appLogger.Fatalf("Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &SQSConfig{
		Debug:                      conf.Debug,
		AWSRegion:                  conf.AWSRegion,
		SQSEndpoint:                conf.SQSEndpoint,
		DiagnosisPortscanQueueName: conf.DiagnosisPortscanQueueName,
		DiagnosisPortscanQueueURL:  conf.DiagnosisPortscanQueueURL,
		MaxNumberOfMessage:         conf.MaxNumberOfMessage,
		WaitTimeSecond:             conf.WaitTimeSecond,
	}
	consumer := newSQSConsumer(sqsConf)
	appLogger.Info("Start the portscan SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.StatusLoggingHandler(appLogger,
					mimosasqs.TracingHandler(getFullServiceName(),
						f.FinalizeHandler(handler))))))
}
