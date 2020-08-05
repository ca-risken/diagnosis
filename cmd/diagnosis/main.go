package main

import (
	"fmt"
	"net"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := newDiagnosisConfig()
	if err != nil {
		panic(err)
	}

	if err := initLogger(conf.LogLevel); err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		logger.Error("Failed to Opening Port", zap.Error(err))
	}

	defer syncLogger()

	server := grpc.NewServer()
	diagnosisServer := newDiagnosisService(conf.DB)
	diagnosisDataSourceServer := newDiagnosisDataSourceService(conf.DB)
	relDiagnosisDataSourceServer := newRelDiagnosisDataSourceService(conf.DB, conf.SQS)
	diagnosis.RegisterDiagnosisServiceServer(server, diagnosisServer)
	diagnosis.RegisterDiagnosisDataSourceServiceServer(server, diagnosisDataSourceServer)
	diagnosis.RegisterRelDiagnosisDataSourceServiceServer(server, relDiagnosisDataSourceServer)

	reflection.Register(server) // enable reflection API
	logger.Info("Starting gRPC server", zap.String("port", conf.Port))
	if err := server.Serve(l); err != nil {
		logger.Error("Failed to gRPC serve", zap.Error(err))
	}
}
