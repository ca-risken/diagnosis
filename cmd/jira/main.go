package main

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
)

func main() {
	conf, err := newBackendConfig()
	if err != nil {
		panic(err)
	}
	if err := initLogger(conf.LogLevel); err != nil {
		panic(err)
	}
	ctx := context.Background()
	mimosaxray.InitXRay(xray.Config{})
	consumer := newSQSConsumer()
	logger.Info("Start the jira SQS consumer server...")
	consumer.Start(ctx,
		mimosaxray.MessageTracingHandler(conf.EnvName, "diagnosis.jira", newHandler()))
}
