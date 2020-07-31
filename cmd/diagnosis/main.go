package main

import (
	"fmt"
	"net"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type diagnosisConfig struct {
	Port string `default:"19001"`
}

func main() {
	var conf diagnosisConfig
	err := envconfig.Process("Diagnosis", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	if err := initLogger(c.LogLevel); err != nil {

		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Fatal(err)
	}

	server := grpc.NewServer()
	diagnosisServer := newDiagnosisService()
	diagnosis.RegisterDiagnosisServiceServer(server, diagnosisServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}
