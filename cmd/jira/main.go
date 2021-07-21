package main

import (
	"context"

	mimosaxray "github.com/CyberAgent/mimosa-common/pkg/xray"
	"github.com/aws/aws-xray-sdk-go/xray"
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
		mimosaxray.MessageTracingHandler(conf.EnvName, "aws.cloudsploit", newHandler()))
}
