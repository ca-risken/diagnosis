package main

import (
	"github.com/kelseyhightower/envconfig"
)

type diagnosisConfig struct {
	Port     string `default:"19001"`
	EnvName  string `default:"default" split_words:"true"`
	LogLevel string `split_words:"true" default:"debug"`

	DB  diagnosisRepoInterface
	SQS *sqsClient
}

func newDiagnosisConfig() (*diagnosisConfig, error) {
	config := &diagnosisConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	config.DB = newDiagnosisRepository()
	config.SQS = newSQSClient()
	return config, nil
}
