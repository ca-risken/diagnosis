package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type sqsConfig struct {
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://localhost:9324"`

	DiagnosisJiraQueueURL            string `split_words:"true" required:"true"`
	DiagnosisWpscanQueueURL          string `split_words:"true" required:"true"`
	DiagnosisPortscanQueueURL        string `split_words:"true" required:"true"`
	DiagnosisApplicationScanQueueURL string `split_words:"true" required:"true"`
}

type sqsAPI interface {
	send(ctx context.Context, msg *message.JiraQueueMessage) (*sqs.SendMessageOutput, error)
	sendWpscanMessage(ctx context.Context, msg *message.WpscanQueueMessage) (*sqs.SendMessageOutput, error)
	sendPortscanMessage(ctx context.Context, msg *message.PortscanQueueMessage) (*sqs.SendMessageOutput, error)
	sendApplicationScanMessage(ctx context.Context, msg *message.ApplicationScanQueueMessage) (*sqs.SendMessageOutput, error)
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
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	session := sqs.New(sess, &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.Endpoint,
	})
	xray.AWS(session.Client)

	return &sqsClient{
		svc: session,
		queueURLMap: map[string]string{
			// queueURLMap:
			"diagnosis:jira":             conf.DiagnosisJiraQueueURL,
			"diagnosis:wpscan":           conf.DiagnosisWpscanQueueURL,
			"diagnosis:portscan":         conf.DiagnosisPortscanQueueURL,
			"diagnosis:application-scan": conf.DiagnosisApplicationScanQueueURL,
		},
	}
}

func (s *sqsClient) send(ctx context.Context, msg *message.JiraQueueMessage) (*sqs.SendMessageOutput, error) {
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
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
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

func (s *sqsClient) sendWpscanMessage(ctx context.Context, msg *message.WpscanQueueMessage) (*sqs.SendMessageOutput, error) {
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
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
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

func (s *sqsClient) sendPortscanMessage(ctx context.Context, msg *message.PortscanQueueMessage) (*sqs.SendMessageOutput, error) {
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
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
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

func (s *sqsClient) sendApplicationScanMessage(ctx context.Context, msg *message.ApplicationScanQueueMessage) (*sqs.SendMessageOutput, error) {
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
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
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
