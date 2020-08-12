package main

import (
	"context"
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
	consumer := newSQSConsumer()
	logger.Info("Start the jira SQS consumer server...")
	consumer.Start(ctx, newHandler())
}
