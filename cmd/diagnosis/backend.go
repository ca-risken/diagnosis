package main

import (
	"github.com/gassara-kys/envconfig"
)

type diagnosisConfig struct {
	Port     string `default:"19001"`
	EnvName  string `default:"local" split_words:"true"`
	LogLevel string `default:"debug" split_words:"true"`

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
