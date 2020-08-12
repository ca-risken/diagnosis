package main

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var logger *zap.Logger
var sqsLogger = initSQSLogger()

func initLogger(level string) error {
	logger, _ = zap.NewDevelopment()

	return nil
}

func syncLogger() {
	logger.Sync()
}

func initSQSLogger() *logrus.Logger {
	sqsLogger := logrus.New()
	sqsLogger.SetFormatter(&logrus.JSONFormatter{})
	// logger.SetLevel(logrus.DebugLevel)
	return sqsLogger
}
