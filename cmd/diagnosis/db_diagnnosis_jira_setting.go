package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (r *diagnosisRepository) ListJiraSetting(projectID, jiraSettingID uint32) (*[]JiraSetting, error) {
	query := `select * from jira_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(jiraSettingID) {
		query += " and jira_setting_id = ?"
		params = append(params, jiraSettingID)
	}
	var data []JiraSetting
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetJiraSetting(projectID uint32, jiraSettingID uint32) (*JiraSetting, error) {
	var data JiraSetting
	if err := r.SlaveDB.Where("project_id = ? AND jira_setting_id = ?", projectID, jiraSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertJiraSetting(data *JiraSetting) (*JiraSetting, error) {
	var savedData JiraSetting
	if err := r.MasterDB.Where("project_id = ? AND jira_setting_id = ?", data.ProjectID, data.JiraSettingID).Assign(data).FirstOrCreate(&savedData).Error; err != nil {
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
