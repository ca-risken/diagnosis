package main

import (
	"context"
	"fmt"
	"net"

	"github.com/ca-risken/common/pkg/profiler"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/gassara-kys/envconfig"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	nameSpace   = "diagnosis"
	serviceName = "diagnosis"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	// backend
	Port            string   `default:"19001"`
	EnvName         string   `default:"local" split_words:"true"`
	LogLevel        string   `default:"debug" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

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

	DiagnosisWpscanQueueURL          string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-wpscan"`
	DiagnosisPortscanQueueURL        string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-portscan"`
	DiagnosisApplicationScanQueueURL string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-applicationscan"`

	// grpc
	CoreAddr string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
}

func main() {
	ctx := context.Background()
	var appConf AppConfig
	err := envconfig.Process("", &appConf)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(appConf.ProfileTypes)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(appConf.ProfileExporter)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      appConf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: appConf.EnvName,
		Debug:       appConf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

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
		DiagnosisWpscanQueueURL:          appConf.DiagnosisWpscanQueueURL,
		DiagnosisPortscanQueueURL:        appConf.DiagnosisPortscanQueueURL,
		DiagnosisApplicationScanQueueURL: appConf.DiagnosisApplicationScanQueueURL,
	}
	service.sqs = newSQSClient(sqsConf)
	service.projectClient = newProjectClient(appConf.CoreAddr)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", appConf.Port))
	if err != nil {
		appLogger.Errorf(ctx, "Failed to Opening Port, error: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				grpctrace.UnaryServerInterceptor())))
	diagnosis.RegisterDiagnosisServiceServer(server, service)

	reflection.Register(server) // enable reflection API
	appLogger.Infof(ctx, "Starting gRPC server, port: %v", appConf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Errorf(ctx, "Failed to gRPC serve, error: %v", err)
	}
}
