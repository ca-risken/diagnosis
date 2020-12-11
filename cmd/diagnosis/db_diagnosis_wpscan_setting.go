package main

import (
	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (r *diagnosisRepository) ListWpscanSetting(projectID, diagnoosisDataSourceID uint32) (*[]model.WpscanSetting, error) {
	query := `select * from wpscan_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.WpscanSetting
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetWpscanSetting(projectID uint32, wpscanSettingID uint32) (*model.WpscanSetting, error) {
	var data model.WpscanSetting
	if err := r.SlaveDB.Where("project_id = ? AND wpscan_setting_id = ?", projectID, wpscanSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertWpscanSetting(data *model.WpscanSetting) (*model.WpscanSetting, error) {
	var savedData model.WpscanSetting
	update := wpscanSettingToMap(data)
	if err := r.MasterDB.Where("project_id = ? AND wpscan_setting_id = ?", data.ProjectID, data.WpscanSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeleteWpscanSetting(projectID uint32, wpscanSettingID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND wpscan_setting_id = ?", projectID, wpscanSettingID).Delete(model.WpscanSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) ListAllWpscanSetting() (*[]model.WpscanSetting, error) {
	var data []model.WpscanSetting
	if err := r.SlaveDB.Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func wpscanSettingToMap(wpscanSetting *model.WpscanSetting) map[string]interface{} {
	return map[string]interface{}{
		"wpscan_setting_id":        wpscanSetting.WpscanSettingID,
		"diagnosis_data_source_id": wpscanSetting.DiagnosisDataSourceID,
		"project_id":               wpscanSetting.ProjectID,
		"target_url":               wpscanSetting.TargetURL,
		"status":                   wpscanSetting.Status,
		"status_detail":            wpscanSetting.StatusDetail,
		"scan_at":                  wpscanSetting.ScanAt,
	}
}
