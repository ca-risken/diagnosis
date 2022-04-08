package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
)

type SQSConfig struct {
	AWSRegion string
	Endpoint  string

	DiagnosisWpscanQueueURL          string
	DiagnosisPortscanQueueURL        string
	DiagnosisApplicationScanQueueURL string
}

type sqsAPI interface {
	sendWpscanMessage(ctx context.Context, msg *message.WpscanQueueMessage) (*sqs.SendMessageOutput, error)
	sendPortscanMessage(ctx context.Context, msg *message.PortscanQueueMessage) (*sqs.SendMessageOutput, error)
	sendApplicationScanMessage(ctx context.Context, msg *message.ApplicationScanQueueMessage) (*sqs.SendMessageOutput, error)
}

type sqsClient struct {
	svc         *sqs.SQS
	queueURLMap map[string]string
}

func newSQSClient(conf *SQSConfig) *sqsClient {
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
			common.DataSourceNameWPScan:          conf.DiagnosisWpscanQueueURL,
			common.DataSourceNamePortScan:        conf.DiagnosisPortscanQueueURL,
			common.DataSourceNameApplicationScan: conf.DiagnosisApplicationScanQueueURL,
		},
	}
}

func (s *sqsClient) sendWpscanMessage(ctx context.Context, msg *message.WpscanQueueMessage) (*sqs.SendMessageOutput, error) {
	url := s.queueURLMap[msg.DataSource]
	if url == "" {
		return nil, fmt.Errorf("Unknown data_source, value=%s", msg.DataSource)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		appLogger.Errorf("Failed to parse message, error: %v", err)
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	appLogger.Infof("Send message, MessageBody: %v, QueueURL: %v", string(buf), url)
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		appLogger.Errorf("Failed to send message, error: %v", err)
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
		appLogger.Errorf("Failed to parse message, error: %v", err)
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	appLogger.Infof("Send message, MessageBody: %v, QueueURL: %v", string(buf), url)
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		appLogger.Errorf("Failed to send message, error: %v", err)
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
		appLogger.Errorf("Failed to parse message, error: %v", err)
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	appLogger.Infof("Send message, MessageBody: %v, QueueURL: %v", string(buf), url)
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		appLogger.Errorf("Failed to send message, error: %v", err)
		return nil, err
	}
	return resp, nil
}
