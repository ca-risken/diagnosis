package main

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func initLogger(level string) error {
	logger, _ = zap.NewDevelopment()

	return nil
}

func syncLogger() {
	logger.Sync()
}
