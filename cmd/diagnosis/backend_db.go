package main

import (
	"context"
	"fmt"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/gorm"
)

type diagnosisRepoInterface interface {
	ListDiagnosisDataSource(context.Context, uint32, string) (*[]model.DiagnosisDataSource, error)
	GetDiagnosisDataSource(context.Context, uint32, uint32) (*model.DiagnosisDataSource, error)
	UpsertDiagnosisDataSource(context.Context, *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error)
	DeleteDiagnosisDataSource(context.Context, uint32, uint32) error
	ListJiraSetting(context.Context, uint32, uint32) (*[]model.JiraSetting, error)
	GetJiraSetting(context.Context, uint32, uint32) (*model.JiraSetting, error)
	UpsertJiraSetting(context.Context, *model.JiraSetting) (*model.JiraSetting, error)
	DeleteJiraSetting(context.Context, uint32, uint32) error
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

	//for InvokeScan
	ListAllJiraSetting(context.Context) (*[]model.JiraSetting, error)
	ListAllWpscanSetting(context.Context) (*[]model.WpscanSetting, error)
}

type diagnosisRepository struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func newDiagnosisRepository() diagnosisRepoInterface {
	repo := diagnosisRepository{}
	repo.MasterDB = initDB(true)
	repo.SlaveDB = initDB(false)
	return &repo
}

type dbConfig struct {
	MasterHost     string `split_words:"true" required:"true"`
	MasterUser     string `split_words:"true" required:"true"`
	MasterPassword string `split_words:"true" required:"true"`
	SlaveHost      string `split_words:"true"`
	SlaveUser      string `split_words:"true"`
	SlavePassword  string `split_words:"true"`

	Schema  string `required:"true"`
	Port    int    `required:"true"`
	LogMode bool   `split_words:"true" default:"false"`
}

func initDB(isMaster bool) *gorm.DB {
	conf := &dbConfig{}
	if err := envconfig.Process("DB", conf); err != nil {
		panic(fmt.Sprintf("Failed to load DB config., err: %v", err.Error()))
	}
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
	db, err := mimosasql.Open(dsn, conf.LogMode)
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	return db
}
