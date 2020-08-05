package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func (r *diagnosisRepository) ListDiagnosis(projectID uint32, name string) (*[]Diagnosis, error) {
	var data []Diagnosis
	paramName := "%" + name + "%"
	if err := r.SlaveDB.Where("project_id = ? and name like ?", projectID, paramName).Find(&data).Error; err != nil {
		logger.Error("Failed to List Diagnosis", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetDiagnosis(projectID uint32, diagnosisID uint32) (*Diagnosis, error) {
	var data Diagnosis
	if err := r.SlaveDB.Where("project_id = ? AND diagnosis_id = ?", projectID, diagnosisID).First(&data).Error; err != nil {
		logger.Error("Failed to Get Diagnosis", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertDiagnosis(input *Diagnosis) (*Diagnosis, error) {
	var data Diagnosis
	if err := r.MasterDB.Where("project_id = ? AND diagnosis_id = ?", input.ProjectID, input.DiagnosisID).Assign(Diagnosis{ProjectID: input.ProjectID, Name: input.Name}).FirstOrCreate(&data).Error; err != nil {
		logger.Error("Failed to Upsert Diagnosis", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) DeleteDiagnosis(projectID uint32, diagnosisID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND diagnosis_id = ?", projectID, diagnosisID).Delete(Diagnosis{}).Error; err != nil {
		logger.Error("Failed to Delete Diagnosis", zap.Error(err))
		return err
	}
	return nil
}
