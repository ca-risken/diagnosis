package main

import (
	"context"

	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/vikyd/zero"
)

func (r *diagnosisRepository) ListApplicationScan(ctx context.Context, projectID, diagnoosisDataSourceID uint32) (*[]model.ApplicationScan, error) {
	query := `select * from application_scan where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.ApplicationScan
	if err := r.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetApplicationScan(ctx context.Context, projectID uint32, applicationScanID uint32) (*model.ApplicationScan, error) {
	var data model.ApplicationScan
	if err := r.SlaveDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", projectID, applicationScanID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertApplicationScan(ctx context.Context, data *model.ApplicationScan) (*model.ApplicationScan, error) {
	var savedData model.ApplicationScan
	update := applicationScanToMap(data)
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", data.ProjectID, data.ApplicationScanID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeleteApplicationScan(ctx context.Context, projectID uint32, applicationScanID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", projectID, applicationScanID).Delete(model.ApplicationScan{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) ListApplicationScanBasicSetting(ctx context.Context, projectID, applicationScanID uint32) (*[]model.ApplicationScanBasicSetting, error) {
	query := `select * from application_scan_basic_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(applicationScanID) {
		query += " and application_scan_id = ?"
		params = append(params, applicationScanID)
	}
	var data []model.ApplicationScanBasicSetting
	if err := r.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetApplicationScanBasicSetting(ctx context.Context, projectID uint32, applicationScanID uint32) (*model.ApplicationScanBasicSetting, error) {
	var data model.ApplicationScanBasicSetting
	if err := r.SlaveDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", projectID, applicationScanID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertApplicationScanBasicSetting(ctx context.Context, data *model.ApplicationScanBasicSetting) (*model.ApplicationScanBasicSetting, error) {
	var savedData model.ApplicationScanBasicSetting
	update := applicationScanBasicSettingToMap(data)
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_basic_setting_id = ?", data.ProjectID, data.ApplicationScanBasicSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeleteApplicationScanBasicSetting(ctx context.Context, projectID uint32, applicationScanBasicSettingID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_basic_setting_id = ?", projectID, applicationScanBasicSettingID).Delete(model.ApplicationScanBasicSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) ListAllApplicationScan(ctx context.Context) (*[]model.ApplicationScan, error) {
	var data []model.ApplicationScan
	if err := r.SlaveDB.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func applicationScanToMap(applicationScan *model.ApplicationScan) map[string]interface{} {
	settingMap := map[string]interface{}{
		"application_scan_id":      applicationScan.ApplicationScanID,
		"diagnosis_data_source_id": applicationScan.DiagnosisDataSourceID,
		"project_id":               applicationScan.ProjectID,
		"name":                     applicationScan.Name,
		"scan_type":                applicationScan.ScanType,
		"status":                   applicationScan.Status,
		"status_detail":            applicationScan.StatusDetail,
	}
	if !zero.IsZeroVal(applicationScan.ScanAt) {
		settingMap["scan_at"] = applicationScan.ScanAt
	}
	return settingMap
}

func applicationScanBasicSettingToMap(applicationScanBasicSetting *model.ApplicationScanBasicSetting) map[string]interface{} {
	settingMap := map[string]interface{}{
		"application_scan_basic_setting_id": applicationScanBasicSetting.ApplicationScanBasicSettingID,
		"application_scan_id":               applicationScanBasicSetting.ApplicationScanID,
		"project_id":                        applicationScanBasicSetting.ProjectID,
		"target":                            applicationScanBasicSetting.Target,
		"max_depth":                         applicationScanBasicSetting.MaxDepth,
		"max_children":                      applicationScanBasicSetting.MaxChildren,
	}
	return settingMap
}
