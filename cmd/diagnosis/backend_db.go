package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/diagnosis/pkg/model"
	"gorm.io/gorm"
)

type diagnosisRepoInterface interface {
	ListDiagnosisDataSource(context.Context, uint32, string) (*[]model.DiagnosisDataSource, error)
	GetDiagnosisDataSource(context.Context, uint32, uint32) (*model.DiagnosisDataSource, error)
	UpsertDiagnosisDataSource(context.Context, *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error)
	DeleteDiagnosisDataSource(context.Context, uint32, uint32) error
	ListWpscanSetting(context.Context, uint32, uint32) (*[]model.WpscanSetting, error)
	GetWpscanSetting(context.Context, uint32, uint32) (*model.WpscanSetting, error)
	UpsertWpscanSetting(context.Context, *model.WpscanSetting) (*model.WpscanSetting, error)
	DeleteWpscanSetting(context.Context, uint32, uint32) error
	ListPortscanSetting(context.Context, uint32, uint32) (*[]model.PortscanSetting, error)
	GetPortscanSetting(context.Context, uint32, uint32) (*model.PortscanSetting, error)
	UpsertPortscanSetting(context.Context, *model.PortscanSetting) (*model.PortscanSetting, error)
	DeletePortscanSetting(context.Context, uint32, uint32) error
	ListPortscanTarget(context.Context, uint32, uint32) (*[]model.PortscanTarget, error)
	GetPortscanTarget(context.Context, uint32, uint32) (*model.PortscanTarget, error)
	UpsertPortscanTarget(context.Context, *model.PortscanTarget) (*model.PortscanTarget, error)
	DeletePortscanTarget(context.Context, uint32, uint32) error
	DeletePortscanTargetByPortscanSettingID(context.Context, uint32, uint32) error
	ListApplicationScan(context.Context, uint32, uint32) (*[]model.ApplicationScan, error)
	GetApplicationScan(context.Context, uint32, uint32) (*model.ApplicationScan, error)
	UpsertApplicationScan(context.Context, *model.ApplicationScan) (*model.ApplicationScan, error)
	DeleteApplicationScan(context.Context, uint32, uint32) error
	ListApplicationScanBasicSetting(context.Context, uint32, uint32) (*[]model.ApplicationScanBasicSetting, error)
	GetApplicationScanBasicSetting(context.Context, uint32, uint32) (*model.ApplicationScanBasicSetting, error)
	UpsertApplicationScanBasicSetting(context.Context, *model.ApplicationScanBasicSetting) (*model.ApplicationScanBasicSetting, error)
	DeleteApplicationScanBasicSetting(context.Context, uint32, uint32) error

	//for InvokeScan
	ListAllWpscanSetting(context.Context) (*[]model.WpscanSetting, error)
}

type diagnosisRepository struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func newDiagnosisRepository(conf *DBConfig) diagnosisRepoInterface {
	repo := diagnosisRepository{}
	repo.MasterDB = initDB(conf, true)
	repo.SlaveDB = initDB(conf, false)
	return &repo
}

type DBConfig struct {
	MasterHost     string
	MasterUser     string
	MasterPassword string
	SlaveHost      string
	SlaveUser      string
	SlavePassword  string

	Schema        string
	Port          int
	LogMode       bool
	MaxConnection int
}

func initDB(conf *DBConfig, isMaster bool) *gorm.DB {
	var user, pass, host string
	if isMaster {
		user = conf.MasterUser
		pass = conf.MasterPassword
		host = conf.MasterHost
	} else {
		user = conf.SlaveUser
		pass = conf.SlavePassword
		host = conf.SlaveHost
	}

	dsn := fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
		user, pass, host, conf.Port, conf.Schema)
	db, err := mimosasql.Open(dsn, conf.LogMode, conf.MaxConnection)
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	return db
}
