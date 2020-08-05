package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func (r *diagnosisRepository) ListDiagnosisDataSource(projectID uint32, name string) (*[]DiagnosisDataSource, error) {
	var data []DiagnosisDataSource
	paramName := "%" + name + "%"
	if err := r.SlaveDB.Where("name like ?", paramName).Find(&data).Error; err != nil {
		logger.Error("Failed to List DiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetDiagnosisDataSource(projectID uint32, diagnosisDataSourceID uint32) (*DiagnosisDataSource, error) {
	var data DiagnosisDataSource
	if err := r.SlaveDB.Where("diagnosis_data_source_id = ?", diagnosisDataSourceID).First(&data).Error; err != nil {
		logger.Error("Failed to Get DiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertDiagnosisDataSource(input *DiagnosisDataSource) (*DiagnosisDataSource, error) {
	var data DiagnosisDataSource
	if err := r.MasterDB.Where("diagnosis_data_source_id = ?", input.DiagnosisDataSourceID).Assign(DiagnosisDataSource{Name: input.Name, Description: input.Description, MaxScore: input.MaxScore}).FirstOrCreate(&data).Error; err != nil {
		logger.Error("Failed to Upsert DiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) DeleteDiagnosisDataSource(projectID uint32, diagnosisDataSourceID uint32) error {
	if err := r.MasterDB.Where("diagnosis_data_source_id =  ?", diagnosisDataSourceID).Delete(DiagnosisDataSource{}).Error; err != nil {
		logger.Error("Failed to Delete DiagnosisDataSource", zap.Error(err))
		return err
	}
	return nil
}
