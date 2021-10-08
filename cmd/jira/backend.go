package main

import (
	"github.com/kelseyhightower/envconfig"
)

type BackendConfig struct {
	LogLevel string `split_words:"true" default:"debug"`
	EnvName  string `default:"default" split_words:"true"`
}

func newBackendConfig() (*BackendConfig, error) {
	config := &BackendConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	return config, nil
}
