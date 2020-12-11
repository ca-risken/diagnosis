package main

import (
	"context"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type findingConfig struct {
	FindingSvcAddr string `required:"true" split_words:"true"`
}

func newFindingClient() finding.FindingServiceClient {
	var conf findingConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Errorf("Faild to load finding config. error: %v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.FindingSvcAddr)
	if err != nil {
		appLogger.Errorf("Faild to get GRPC connection. error: %v", err)
	}
	return finding.NewFindingServiceClient(conn)
}

type alertConfig struct {
	AlertSvcAddr string `required:"true" split_words:"true"`
}

func newAlertClient() alert.AlertServiceClient {
	var conf alertConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Errorf("Faild to load alert config. error: %v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.AlertSvcAddr)
	if err != nil {
		appLogger.Errorf("Faild to get GRPC connection. error: %v", err)
	}
	return alert.NewAlertServiceClient(conn)
}

type diagnosisConfig struct {
	DiagnosisSvcAddr string `required:"true" split_words:"true"`
}

func newDiagnosisClient() diagnosis.DiagnosisServiceClient {
	var conf diagnosisConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Errorf("Faild to load diagnosis config. error: %v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.DiagnosisSvcAddr)
	if err != nil {
		appLogger.Errorf("Faild to get GRPC connection. error: %v", err)
	}
	return diagnosis.NewDiagnosisServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
