package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

type DiagnosisConfig struct {
	Port     string `default:"19001"`
	LogLevel string `split_words:"true" default:"debug"`

	DB  diagnosisRepoInterface
	SQS *sqsClient
}

func newDiagnosisConfig() (*DiagnosisConfig, error) {
	config := &DiagnosisConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	config.DB = newDiagnosisRepository()
	config.SQS = newSQSClient()
	return config, nil
}
