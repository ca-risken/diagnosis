package main

import (
	"context"
	"time"

	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/gassara-kys/envconfig"
	"google.golang.org/grpc"
)

type findingConfig struct {
	FindingSvcAddr string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
}

func newFindingClient() finding.FindingServiceClient {
	var conf findingConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Faild to load finding config error: err=%+v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.FindingSvcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return finding.NewFindingServiceClient(conn)
}

type alertConfig struct {
	AlertSvcAddr string `required:"true" split_words:"true" default:"alert.core.svc.cluster.local:8004"`
}

func newAlertClient() alert.AlertServiceClient {
	var conf alertConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Faild to load alert config error: err=%+v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.AlertSvcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return alert.NewAlertServiceClient(conn)
}

type diagnosisConfig struct {
	DiagnosisSvcAddr string `required:"true" split_words:"true" default:"diagnosis.diagnosis.svc.cluster.local:19001"`
}

func newDiagnosisClient() diagnosis.DiagnosisServiceClient {
	var conf diagnosisConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Faild to load diagnosis config error: err=%+v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.DiagnosisSvcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return diagnosis.NewDiagnosisServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	// gRPCクライアントの呼び出し回数が非常に多くトレーシング情報の送信がエラーになるため、トレースは無効にしておく
	//conn, err := grpc.DialContext(ctx, addr,
	//	grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
