package main

import (
	"fmt"

	//	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	//	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	//	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type diagnosisRepoInterface interface {
	ListDiagnosis(uint32, uint32, string) (*[]model.Diagnosis, error)
	GetDiagnosis(uint32, string) (*model.Diagnosis, error)
	UpsertDiagnosis(*model.Diagnosis) (*model.Diagnosis, error)
	DeleteDiagnosis(uint32, uint32) error
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
		appLogger.Fatalf("Failed to load DB config. err: %+v", err)
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
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
		return nil
	}
	db.LogMode(conf.LogMode)
	db.SingularTable(true) // if set this to true, `User`'s default table name will be `user`
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}

func (a *diagnosisRepository) ListDiagnosis(projectID, diagnosisID uint32, diagnosisAccountID string) (*[]model.Diagnosis, error) {
	query := `
select
  *
from
  diagnosis
where
  project_id = ?
`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnosisID) {
		query += " and diagnosis_id = ?"
		params = append(params, diagnosisID)
	}
	if !zero.IsZeroVal(diagnosisAccountID) {
		query += " and diagnosis_account_id = ?"
		params = append(params, diagnosisAccountID)
	}

	data := []model.Diagnosis{}
	if err := a.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetDiagnosisByAccountID = `select * from diagnosis where project_id = ? and diagnosis_account_id = ?`

func (a *diagnosisRepository) GetDiagnosisByAccountID(projectID uint32, diagnosisAccountID string) (*model.Diagnosis, error) {
	data := model.Diagnosis{}
	if err := a.SlaveDB.Raw(selectGetDiagnosisByAccountID, projectID, diagnosisAccountID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertDiagnosis = `
INSERT INTO diagnosis
  (diagnosis_id, name, project_id, diagnosis_account_id)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  project_id=VALUES(project_id)
`

func (a *diagnosisRepository) UpsertDiagnosis(data *model.Diagnosis) (*model.Diagnosis, error) {
	if err := a.MasterDB.Exec(insertUpsertDiagnosis,
		data.DiagnosisID, data.Name, data.ProjectID, data.DiagnosisAccountID).Error; err != nil {
		return nil, err
	}

	updated, err := a.GetDiagnosisByAccountID(data.ProjectID, data.DiagnosisAccountID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

const deleteAws = `delete from diagnosis where project_id = ? and diagnosis_id = ?`

func (a *diagnosisRepository) DeleteDiagnosis(projectID, diagnosisID uint32) error {
	if err := a.MasterDB.Exec(deleteAws, projectID, diagnosisID).Error; err != nil {
		return err
	}
	return nil
}

type dataSource struct {
	DiagnosisDataSourceID uint32
	DataSource            string
	MaxScore              float32
	DiagnosisID           uint32 `gorm:"column:diagnosis_id"`
	ProjectID             uint32
	AssumeRoleArn         string
	ExternalID            string
}

func (a *diagnosisRepository) ListDataSource(projectID, diagnosisID uint32, ds string) (*[]dataSource, error) {
	var params []interface{}
	query := `
select
  ads.diagnosis_data_source_id
  , ads.data_source
  , ads.max_score
  , ards.diagnosis_id
  , ards.project_id
  , ards.assume_role_arn
  , ards.external_id
from
  diagnosis_data_source ads
  left outer join (
    select * from diagnosis_rel_data_source where project_id = ? `
	params = append(params, projectID)

	if !zero.IsZeroVal(diagnosisID) {
		query += " and diagnosis_id = ?"
		params = append(params, diagnosisID)
	}
	if !zero.IsZeroVal(ds) {
		query += " and dataSource = ?"
		params = append(params, ds)
	}
	query += `
  ) ards using(diagnosis_data_source_id)
order by
  ads.diagnosis_data_source_id
`
	data := []dataSource{}
	if err := a.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertDiagnosisRelDataSource = `
INSERT INTO diagnosis_rel_data_source
  (diagnosis_id, diagnosis_data_source_id, project_id, assume_role_arn, external_id)
VALUES
  (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  project_id=VALUES(project_id),
  assume_role_arn=VALUES(assume_role_arn),
  external_id=VALUES(external_id)
`

func (a *diagnosisRepository) UpsertDiagnosisRelDataSource(data *diagnosis.DataSourceForAttach) (*model.DiagnosisRelDataSource, error) {
	if err := a.MasterDB.Exec(insertUpsertDiagnosisRelDataSource, data.AwsId, data.AwsDataSourceId, data.ProjectId, data.AssumeRoleArn, data.ExternalId).Error; err != nil {
		return nil, err
	}
	return a.GetDiagnosisRelDataSourceByID(data.AwsId, data.AwsDataSourceId, data.ProjectId)
}

const selectGetDiagnosisRelDataSourceByID = `select * from diagnosis_rel_data_source where diagnosis_id = ? and diagnosis_data_source_id = ? and project_id = ?`

func (a *diagnosisRepository) GetDiagnosisRelDataSourceByID(diagnosisID, diagnosisDataSourceID, projectID uint32) (*model.DiagnosisRelDataSource, error) {
	data := model.DiagnosisRelDataSource{}
	if err := a.SlaveDB.Raw(selectGetDiagnosisRelDataSourceByID, diagnosisID, diagnosisDataSourceID, projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteDiagnosisRelDataSource = `delete from diagnosis_rel_data_source where project_id = ? and diagnosis_id = ? and diagnosis_data_source_id = ?`

func (a *diagnosisRepository) DeleteDiagnosisRelDataSource(projectID, diagnosisID, diagnosisDataSourceID uint32) error {
	if err := a.MasterDB.Exec(deleteDiagnosisRelDataSource, projectID, diagnosisID, diagnosisDataSourceID).Error; err != nil {
		return err
	}
	return nil
}

const selectDiagnosisDataSourceForMessage = `
select 
  ads.data_source        as data_source
  , ards.project_id      as project_id
  , a.diagnosis_account_id     as account_id
  , ards.assume_role_arn as assume_role_arn
  , ards.external_id     as external_id
from
  diagnosis_rel_data_source ards
  inner join diagnosis a using(diagnosis_id)
  inner join diagnosis_data_source ads using(diagnosis_data_source_id)
where
  ards.diagnosis_id = ?
  and ards.diagnosis_data_source_id = ?
	and ards.project_id = ? 
`

func (a *diagnosisRepository) GetDiagnosisDataSourceForMessage(diagnosisID, diagnosisDataSourceID, projectID uint32) (*message.DiagnosisQueueMessage, error) {
	data := message.DiagnosisQueueMessage{}
	if err := a.SlaveDB.Raw(selectDiagnosisDataSourceForMessage, diagnosisID, diagnosisDataSourceID, projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
