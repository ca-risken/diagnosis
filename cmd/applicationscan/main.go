package main

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/kelseyhightower/envconfig"
)

type serviceConfig struct {
	EnvName string `default:"default" split_words:"true"`
}

func main() {
	var conf serviceConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	//	_ = tempHandler()
	ctx := context.Background()

	mimosaxray.InitXRay(xray.Config{})
	consumer := newSQSConsumer()
	appLogger.Info("Start the ApplicationScan SQS consumer server...")
	consumer.Start(ctx,
		mimosaxray.MessageTracingHandler(conf.EnvName, "diagnosis.application_scan", newHandler()))
}