package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

type BackendConfig struct {
	//	Port     string `default:"19001"`
	LogLevel string `split_words:"true" default:"debug"`

	//	DB  diagnosisRepoInterface
	//	SQS *sqsClient
}

func newBackendConfig() (*BackendConfig, error) {
	config := &BackendConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	//	config.DB = newDiagnosisRepository()
	//	config.SQS = newSQSClient()
	return config, nil
}
