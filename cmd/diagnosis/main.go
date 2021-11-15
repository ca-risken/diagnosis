package main

import (
	"fmt"
	"net"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := newDiagnosisConfig()
	if err != nil {
		panic(err)
	}
	mimosaxray.InitXRay(xray.Config{})

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Errorf("Failed to Opening Port, error: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				xray.UnaryServerInterceptor(),
				mimosaxray.AnnotateEnvTracingUnaryServerInterceptor(conf.EnvName))))
	diagnosisServer := newDiagnosisService(conf.DB, conf.SQS)
	diagnosis.RegisterDiagnosisServiceServer(server, diagnosisServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server, port: %v", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Errorf("Failed to gRPC serve, error: %v", err)
	}
}
