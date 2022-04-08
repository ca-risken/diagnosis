package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "diagnosis"
	serviceName = "wpscan"
	settingURL  = "https://docs.security-hub.jp/diagnosis/wpscan_datasource/"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`

	// sqs
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	DiagnosisWpscanQueueName string `split_words:"true" default:"diagnosis-wpscan"`
	DiagnosisWpscanQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-wpscan"`
	MaxNumberOfMessage       int64  `split_words:"true" default:"10"`
	WaitTimeSecond           int64  `split_words:"true" default:"20"`
	// grpc
	FindingSvcAddr   string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
	AlertSvcAddr     string `required:"true" split_words:"true" default:"alert.core.svc.cluster.local:8004"`
	DiagnosisSvcAddr string `required:"true" split_words:"true" default:"diagnosis.diagnosis.svc.cluster.local:19001"`
	// wpscan
	ResultPath         string `split_words:"true" required:"true" default:"/tmp"`
	WpscanVulndbApikey string `split_words:"true"`
}

func main() {
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	ctx := context.Background()
	err = mimosaxray.InitXRay(xray.Config{})
	if err != nil {
		appLogger.Fatal(err.Error())
	}

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

	handler := &sqsHandler{}
	handler.wpscanConfig = WpscanConfig{
		ResultPath:         conf.ResultPath,
		WpscanVulndbApikey: conf.WpscanVulndbApikey,
	}
	appLogger.Info("Start Wpscan Client")
	handler.findingClient = newFindingClient(conf.FindingSvcAddr)
	appLogger.Info("Start Finding Client")
	handler.alertClient = newAlertClient(conf.AlertSvcAddr)
	appLogger.Info("Start Alert Client")
	handler.diagnosisClient = newDiagnosisClient(conf.DiagnosisSvcAddr)
	appLogger.Info("Start Diagnosis Client")
	f, err := mimosasqs.NewFinalizer(common.DataSourceNameWPScan, settingURL, conf.FindingSvcAddr, &mimosasqs.DataSourceRecommnend{
		ScanFailureRisk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", common.DataSourceNameWPScan),
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
		appLogger.Fatalf("Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &SQSConfig{
		AWSRegion:                conf.AWSRegion,
		Endpoint:                 conf.Endpoint,
		DiagnosisWpscanQueueName: conf.DiagnosisWpscanQueueName,
		DiagnosisWpscanQueueURL:  conf.DiagnosisWpscanQueueURL,
		MaxNumberOfMessage:       conf.MaxNumberOfMessage,
		WaitTimeSecond:           conf.WaitTimeSecond,
	}
	consumer := newSQSConsumer(sqsConf)
	appLogger.Info("Start the wpscan SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.StatusLoggingHandler(appLogger,
					mimosaxray.MessageTracingHandler(conf.EnvName, getFullServiceName(),
						f.FinalizeHandler(handler))))))
}
