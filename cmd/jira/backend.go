package main

import (
	"github.com/gassara-kys/envconfig"
)

type BackendConfig struct {
	LogLevel string `split_words:"true" default:"debug"`
	EnvName  string `split_words:"true" default:"local"`
}

func newBackendConfig() (*BackendConfig, error) {
	config := &BackendConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	return config, nil
}
