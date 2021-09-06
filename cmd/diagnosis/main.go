package main

import (
	"fmt"
	"net"

	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/aws/aws-xray-sdk-go/xray"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := newDiagnosisConfig()
	if err != nil {
		panic(err)
	}
	mimosaxray.InitXRay(xray.Config{})

	if err := initLogger(conf.LogLevel); err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		logger.Error("Failed to Opening Port", zap.Error(err))
	}

	defer syncLogger()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				xray.UnaryServerInterceptor(),
				mimosaxray.AnnotateEnvTracingUnaryServerInterceptor(conf.EnvName))))
	diagnosisServer := newDiagnosisService(conf.DB, conf.SQS)
	diagnosis.RegisterDiagnosisServiceServer(server, diagnosisServer)

	reflection.Register(server) // enable reflection API
	logger.Info("Starting gRPC server", zap.String("port", conf.Port))
	if err := server.Serve(l); err != nil {
		logger.Error("Failed to gRPC serve", zap.Error(err))
	}
}
