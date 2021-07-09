package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type diagnosisRepoInterface interface {
	ListDiagnosisDataSource(uint32, string) (*[]model.DiagnosisDataSource, error)
	GetDiagnosisDataSource(uint32, uint32) (*model.DiagnosisDataSource, error)
	UpsertDiagnosisDataSource(*model.DiagnosisDataSource) (*model.DiagnosisDataSource, error)
	DeleteDiagnosisDataSource(uint32, uint32) error
	ListJiraSetting(uint32, uint32) (*[]model.JiraSetting, error)
	GetJiraSetting(uint32, uint32) (*model.JiraSetting, error)
	UpsertJiraSetting(*model.JiraSetting) (*model.JiraSetting, error)
	DeleteJiraSetting(uint32, uint32) error
	ListWpscanSetting(uint32, uint32) (*[]model.WpscanSetting, error)
	GetWpscanSetting(uint32, uint32) (*model.WpscanSetting, error)
	UpsertWpscanSetting(*model.WpscanSetting) (*model.WpscanSetting, error)
	DeleteWpscanSetting(uint32, uint32) error
	ListPortscanSetting(uint32, uint32) (*[]model.PortscanSetting, error)
	GetPortscanSetting(uint32, uint32) (*model.PortscanSetting, error)
	UpsertPortscanSetting(*model.PortscanSetting) (*model.PortscanSetting, error)
	DeletePortscanSetting(uint32, uint32) error
	ListPortscanTarget(uint32, uint32) (*[]model.PortscanTarget, error)
	GetPortscanTarget(uint32, uint32) (*model.PortscanTarget, error)
	UpsertPortscanTarget(*model.PortscanTarget) (*model.PortscanTarget, error)
	DeletePortscanTarget(uint32, uint32) error
	DeletePortscanTargetByPortscanSettingID(uint32, uint32) error

	//for InvokeScan
	ListAllJiraSetting() (*[]model.JiraSetting, error)
	ListAllWpscanSetting() (*[]model.WpscanSetting, error)
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

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
			user, pass, host, conf.Port, conf.Schema))
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	db.LogMode(conf.LogMode)
	db.SingularTable(true) // if set this to true, `User`'s default table name will be `user`
	return db
}
