package main

import (
	"fmt"
	"net"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/gassara-kys/envconfig"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AppConfig struct {
	// backend
	Port     string `default:"19001"`
	EnvName  string `default:"local" split_words:"true"`
	LogLevel string `default:"debug" split_words:"true"`

	// db
	DBMasterHost     string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBMasterUser     string `split_words:"true" default:"hoge"`
	DBMasterPassword string `split_words:"true" default:"moge"`
	DBSlaveHost      string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBSlaveUser      string `split_words:"true" default:"hoge"`
	DBSlavePassword  string `split_words:"true" default:"moge"`

	DBSchema        string `required:"true"    default:"mimosa"`
	DBPort          int    `required:"true"    default:"3306"`
	DBLogMode       bool   `split_words:"true" default:"false"`
	DBMaxConnection int    `split_words:"true" default:"10"`

	// sqs
	AWSRegion string `envconfig:"aws_region"   default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	DiagnosisJiraQueueURL            string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-jira"`
	DiagnosisWpscanQueueURL          string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-wpscan"`
	DiagnosisPortscanQueueURL        string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-portscan"`
	DiagnosisApplicationScanQueueURL string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-applicationscan"`

	// grpc
	ProjectSvcAddr string `required:"true" split_words:"true" default:"project.core.svc.cluster.local:8003"`
}

func main() {
	var appConf AppConfig
	err := envconfig.Process("", appConf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	err = mimosaxray.InitXRay(xray.Config{})
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	service := &DiagnosisService{}
	dbConf := &DBConfig{
		MasterHost:     appConf.DBMasterHost,
		MasterUser:     appConf.DBMasterUser,
		MasterPassword: appConf.DBMasterPassword,
		SlaveHost:      appConf.DBSlaveHost,
		SlaveUser:      appConf.DBSlaveUser,
		SlavePassword:  appConf.DBSlavePassword,
		Schema:         appConf.DBSchema,
		Port:           appConf.DBPort,
		LogMode:        appConf.DBLogMode,
		MaxConnection:  appConf.DBMaxConnection,
	}
	service.repository = newDiagnosisRepository(dbConf)
	sqsConf := &SQSConfig{
		AWSRegion:                        appConf.AWSRegion,
		Endpoint:                         appConf.Endpoint,
		DiagnosisJiraQueueURL:            appConf.DiagnosisJiraQueueURL,
		DiagnosisWpscanQueueURL:          appConf.DiagnosisWpscanQueueURL,
		DiagnosisPortscanQueueURL:        appConf.DiagnosisPortscanQueueURL,
		DiagnosisApplicationScanQueueURL: appConf.DiagnosisApplicationScanQueueURL,
	}
	service.sqs = newSQSClient(sqsConf)
	service.projectClient = newProjectClient(appConf.ProjectSvcAddr)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", appConf.Port))
	if err != nil {
		appLogger.Errorf("Failed to Opening Port, error: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				xray.UnaryServerInterceptor(),
				mimosaxray.AnnotateEnvTracingUnaryServerInterceptor(appConf.EnvName))))
	diagnosis.RegisterDiagnosisServiceServer(server, service)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server, port: %v", appConf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Errorf("Failed to gRPC serve, error: %v", err)
	}
}
