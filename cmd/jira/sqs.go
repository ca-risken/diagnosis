package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gassara-kys/go-sqs-poller/worker/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type sqsConfig struct {
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://localhost:9324"`

	DiagnosisJiraQueueName string `split_words:"true" default:"diagnosis-jira"`
	DiagnosisJiraQueueURL  string `split_words:"true" default:"http://localhost:9324/queue/diagnosis-jira"`
	MaxNumberOfMessage     int64  `split_words:"true" default:"10"`
	WaitTimeSecond         int64  `split_words:"true" default:"20"`
}

func newSQSConsumer() *worker.Worker {
	var conf sqsConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	var sqsClient *sqs.SQS
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	if !zero.IsZeroVal(&conf.Endpoint) {
		sqsClient = sqs.New(sess, &aws.Config{
			Region:   &conf.AWSRegion,
			Endpoint: &conf.Endpoint,
		})
	} else {
		sqsClient = sqs.New(sess, &aws.Config{
			Region: &conf.AWSRegion,
		})
	}
	return &worker.Worker{
		Config: &worker.Config{
			QueueName:          conf.DiagnosisJiraQueueName,
			QueueURL:           conf.DiagnosisJiraQueueURL,
			MaxNumberOfMessage: conf.MaxNumberOfMessage,
			WaitTimeSecond:     conf.WaitTimeSecond,
		},
		Log:       appLogger,
		SqsClient: sqsClient,
	}
}
