package main

import (
	"context"

	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/vikyd/zero"
)

func (r *diagnosisRepository) ListJiraSetting(ctx context.Context, projectID, diagnoosisDataSourceID uint32) (*[]model.JiraSetting, error) {
	query := `select * from jira_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.JiraSetting
	if err := r.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetJiraSetting(ctx context.Context, projectID uint32, jiraSettingID uint32) (*model.JiraSetting, error) {
	var data model.JiraSetting
	if err := r.SlaveDB.WithContext(ctx).Where("project_id = ? AND jira_setting_id = ?", projectID, jiraSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertJiraSetting(ctx context.Context, data *model.JiraSetting) (*model.JiraSetting, error) {
	var savedData model.JiraSetting
	update := jiraSettingToMap(data)
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND jira_setting_id = ?", data.ProjectID, data.JiraSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeleteJiraSetting(ctx context.Context, projectID uint32, jiraSettingID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND jira_setting_id = ?", projectID, jiraSettingID).Delete(model.JiraSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) ListAllJiraSetting(ctx context.Context) (*[]model.JiraSetting, error) {
	var data []model.JiraSetting
	if err := r.SlaveDB.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func jiraSettingToMap(jiraSetting *model.JiraSetting) map[string]interface{} {
	settingMap := map[string]interface{}{
		"jira_setting_id":          jiraSetting.JiraSettingID,
		"name":                     jiraSetting.Name,
		"diagnosis_data_source_id": jiraSetting.DiagnosisDataSourceID,
		"project_id":               jiraSetting.ProjectID,
		"identity_field":           jiraSetting.IdentityField,
		"identity_value":           jiraSetting.IdentityValue,
		"jira_id":                  jiraSetting.JiraID,
		"jira_key":                 jiraSetting.JiraKey,
		"status":                   jiraSetting.Status,
		"status_detail":            jiraSetting.StatusDetail,
	}
	if !zero.IsZeroVal(jiraSetting.ScanAt) {
		settingMap["scan_at"] = jiraSetting.ScanAt
	}
	return settingMap
}
