package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type diagnosisRepoInterface interface {
	ListDiagnosis(uint32, string) (*[]Diagnosis, error)
	GetDiagnosis(uint32, uint32) (*Diagnosis, error)
	UpsertDiagnosis(*Diagnosis) (*Diagnosis, error)
	DeleteDiagnosis(uint32, uint32) error
	ListDiagnosisDataSource(uint32, string) (*[]DiagnosisDataSource, error)
	GetDiagnosisDataSource(uint32, uint32) (*DiagnosisDataSource, error)
	UpsertDiagnosisDataSource(*DiagnosisDataSource) (*DiagnosisDataSource, error)
	DeleteDiagnosisDataSource(uint32, uint32) error
	ListRelDiagnosisDataSource(uint32, uint32, uint32) (*[]RelDiagnosisDataSource, error)
	GetRelDiagnosisDataSource(uint32, uint32) (*RelDiagnosisDataSource, error)
	UpsertRelDiagnosisDataSource(*RelDiagnosisDataSource) (*RelDiagnosisDataSource, error)
	DeleteRelDiagnosisDataSource(uint32, uint32) error
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
		//		logger.Error("Failed to open DB.", zap.String("isMaster", strconv.FormatBool(isMaster)), zap.Error(err))
		return nil
	}
	db.LogMode(conf.LogMode)
	db.SingularTable(true) // if set this to true, `User`'s default table name will be `user`
	//	logger.Info("Connected to Database.", zap.String("isMaster", strconv.FormatBool(isMaster)))
	return db
}
