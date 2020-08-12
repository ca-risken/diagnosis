package main

import (
	"encoding/json"
	"fmt"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type sqsConfig struct {
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://localhost:9324"`

	DiagnosisJiraQueueURL string `split_words:"true" required:"true"`
}

type sqsAPI interface {
	send(msg *message.DiagnosisQueueMessage) (*sqs.SendMessageOutput, error)
}

type sqsClient struct {
	svc         *sqs.SQS
	queueURLMap map[string]string
}

func newSQSClient() *sqsClient {
	var conf sqsConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	session := sqs.New(session.New(), &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.Endpoint,
	})
	fmt.Printf("%v\n", conf.AWSRegion)
	fmt.Printf("%v\n", conf.Endpoint)

	return &sqsClient{
		svc: session,
		queueURLMap: map[string]string{
			// queueURLMap:
			"diagnosis:jira": conf.DiagnosisJiraQueueURL,
		},
	}
}

func (s *sqsClient) send(msg *message.DiagnosisQueueMessage) (*sqs.SendMessageOutput, error) {
	url := s.queueURLMap[msg.DataSource]
	if url == "" {
		return nil, fmt.Errorf("Unknown data_source, value=%s", msg.DataSource)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to parse message", zap.Error(err))
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	logger.Info("Send message", zap.String("MessageBody", string(buf)), zap.String("QueueUrl", url))
	resp, err := s.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		logger.Error("Failed to send message", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
