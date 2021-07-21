package main

import (
	"context"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
)

func (r *diagnosisRepository) ListDiagnosisDataSource(ctx context.Context, projectID uint32, name string) (*[]model.DiagnosisDataSource, error) {
	var data []model.DiagnosisDataSource
	paramName := "%" + name + "%"
	if err := r.SlaveDB.WithContext(ctx).Where("name like ?", paramName).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetDiagnosisDataSource(ctx context.Context, projectID uint32, diagnosisDataSourceID uint32) (*model.DiagnosisDataSource, error) {
	var data model.DiagnosisDataSource
	if err := r.SlaveDB.WithContext(ctx).Where("diagnosis_data_source_id = ?", diagnosisDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertDiagnosisDataSource(ctx context.Context, input *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error) {
	var data model.DiagnosisDataSource
	if err := r.MasterDB.WithContext(ctx).Where("diagnosis_data_source_id = ?", input.DiagnosisDataSourceID).Assign(model.DiagnosisDataSource{Name: input.Name, Description: input.Description, MaxScore: input.MaxScore}).FirstOrCreate(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) DeleteDiagnosisDataSource(ctx context.Context, projectID uint32, diagnosisDataSourceID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("diagnosis_data_source_id =  ?", diagnosisDataSourceID).Delete(model.DiagnosisDataSource{}).Error; err != nil {
		return err
	}
	return nil
}
