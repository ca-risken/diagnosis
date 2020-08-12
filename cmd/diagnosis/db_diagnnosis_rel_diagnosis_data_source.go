package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
	"go.uber.org/zap"
)

func (r *diagnosisRepository) ListRelDiagnosisDataSource(projectID, diagnosisID, diagnosisDataSourceID uint32) (*[]RelDiagnosisDataSource, error) {
	query := `select * from rel_diagnosis_data_source where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnosisID) {
		query += " and diagnosis_id = ?"
		params = append(params, diagnosisID)
	}
	if !zero.IsZeroVal(diagnosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnosisDataSourceID)
	}

	var data []RelDiagnosisDataSource
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		logger.Error("Failed to List RelDiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetRelDiagnosisDataSource(projectID uint32, rel_diagnosis_data_sourceID uint32) (*RelDiagnosisDataSource, error) {
	var data RelDiagnosisDataSource
	if err := r.SlaveDB.Where("project_id = ? AND rel_diagnosis_data_source_id = ?", projectID, rel_diagnosis_data_sourceID).First(&data).Error; err != nil {
		logger.Error("Failed to Get RelDiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertRelDiagnosisDataSource(input *RelDiagnosisDataSource) (*RelDiagnosisDataSource, error) {
	var data RelDiagnosisDataSource
	putData := RelDiagnosisDataSource{ProjectID: input.ProjectID, DiagnosisID: input.DiagnosisID, DiagnosisDataSourceID: input.DiagnosisDataSourceID, RecordID: input.RecordID, JiraID: input.JiraID, JiraKey: input.JiraKey}
	if err := r.MasterDB.Where("project_id = ? AND rel_diagnosis_data_source_id = ?", input.ProjectID, input.RelDiagnosisDataSourceID).Assign(putData).FirstOrCreate(&data).Error; err != nil {
		logger.Error("Failed to Upsert RelDiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) DeleteRelDiagnosisDataSource(projectID uint32, rel_diagnosis_data_sourceID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND rel_diagnosis_data_source_id = ?", projectID, rel_diagnosis_data_sourceID).Delete(RelDiagnosisDataSource{}).Error; err != nil {
		logger.Error("Failed to Delete RelDiagnosisDataSource", zap.Error(err))
		return err
	}
	return nil
}
