package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newFindingClient(ctx context.Context, svcAddr string) (finding.FindingServiceClient, error) {
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection. error: %w", err)
	}
	return finding.NewFindingServiceClient(conn), nil
}

func newAlertClient(ctx context.Context, svcAddr string) (alert.AlertServiceClient, error) {
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection. error: %w", err)
	}
	return alert.NewAlertServiceClient(conn), nil
}

func newDiagnosisClient(ctx context.Context, svcAddr string) (diagnosis.DiagnosisServiceClient, error) {
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection. error: %w", err)
	}
	return diagnosis.NewDiagnosisServiceClient(conn), nil
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	// gRPCクライアントの呼び出し回数が非常に多くトレーシング情報の送信がエラーになるため、トレースはしない
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
