package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type diagnosisRepoInterface interface {
	ListDiagnosisDataSource(uint32, string) (*[]DiagnosisDataSource, error)
	GetDiagnosisDataSource(uint32, uint32) (*DiagnosisDataSource, error)
	UpsertDiagnosisDataSource(*DiagnosisDataSource) (*DiagnosisDataSource, error)
	DeleteDiagnosisDataSource(uint32, uint32) error
	ListJiraSetting(uint32, uint32) (*[]JiraSetting, error)
	GetJiraSetting(uint32, uint32) (*JiraSetting, error)
	UpsertJiraSetting(*JiraSetting) (*JiraSetting, error)
	DeleteJiraSetting(uint32, uint32) error
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
		//		logger.Error("Failed to load DB config.", zap.Error(err))
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
