package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

type SQSConfig struct {
	AWSRegion   string `envconfig:"aws_region" default:"ap-northeast-1"`
	SQSEndpoint string `envconfig:"sqs_endpoint" default:"http://localhost:9324"`

	DiagnosisQueueURL string `split_words:"true" required:"true"`
}

type BackendConfig struct {
	DB  DBConfig  `envconfig:"DB"`
	SQS SQSConfig `envconfig:"SQS"`
}

func newBackendConfig() (*BackendConfig, error) {
	backendConfig := &BackendConfig{}
	if err := envconfig.Process("", &backendConfig); err != nil {
		return nil, err
	}

	return &backend, nil
}
