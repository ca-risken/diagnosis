package main

import (
	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (r *diagnosisRepository) ListJiraSetting(projectID, jiraSettingID uint32) (*[]model.JiraSetting, error) {
	query := `select * from jira_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(jiraSettingID) {
		query += " and jira_setting_id = ?"
		params = append(params, jiraSettingID)
	}
	var data []model.JiraSetting
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetJiraSetting(projectID uint32, jiraSettingID uint32) (*model.JiraSetting, error) {
	var data model.JiraSetting
	if err := r.SlaveDB.Where("project_id = ? AND jira_setting_id = ?", projectID, jiraSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertJiraSetting(data *model.JiraSetting) (*model.JiraSetting, error) {
	var savedData model.JiraSetting
	update := jiraSettingToMap(data)
	if err := r.MasterDB.Where("project_id = ? AND jira_setting_id = ?", data.ProjectID, data.JiraSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeleteJiraSetting(projectID uint32, jiraSettingID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND jira_setting_id = ?", projectID, jiraSettingID).Delete(JiraSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func jiraSettingToMap(jiraSetting *model.JiraSetting) map[string]interface{} {
	return map[string]interface{}{
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
		"scan_at":                  jiraSetting.ScanAt,
	}
}
