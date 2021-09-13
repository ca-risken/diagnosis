package main

import (
	"context"

	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/vikyd/zero"
)

func (r *diagnosisRepository) ListPortscanSetting(ctx context.Context, projectID, diagnoosisDataSourceID uint32) (*[]model.PortscanSetting, error) {
	query := `select * from portscan_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.PortscanSetting
	if err := r.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetPortscanSetting(ctx context.Context, projectID uint32, portscanSettingID uint32) (*model.PortscanSetting, error) {
	var data model.PortscanSetting
	if err := r.SlaveDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", projectID, portscanSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertPortscanSetting(ctx context.Context, data *model.PortscanSetting) (*model.PortscanSetting, error) {
	var savedData model.PortscanSetting
	update := portscanSettingToMap(data)
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", data.ProjectID, data.PortscanSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeletePortscanSetting(ctx context.Context, projectID uint32, portscanSettingID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", projectID, portscanSettingID).Delete(model.PortscanSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) ListPortscanTarget(ctx context.Context, projectID, portscanSettingID uint32) (*[]model.PortscanTarget, error) {
	query := `select * from portscan_target where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(portscanSettingID) {
		query += " and portscan_setting_id = ?"
		params = append(params, portscanSettingID)
	}
	var data []model.PortscanTarget
	if err := r.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) GetPortscanTarget(ctx context.Context, projectID uint32, portscanTargetID uint32) (*model.PortscanTarget, error) {
	var data model.PortscanTarget
	if err := r.SlaveDB.WithContext(ctx).Where("project_id = ? AND portscan_target_id = ?", projectID, portscanTargetID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *diagnosisRepository) UpsertPortscanTarget(ctx context.Context, data *model.PortscanTarget) (*model.PortscanTarget, error) {
	var savedData model.PortscanTarget
	update := portscanTargetToMap(data)
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_target_id = ?", data.ProjectID, data.PortscanTargetID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *diagnosisRepository) DeletePortscanTarget(ctx context.Context, projectID uint32, portscanTargetID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_target_id = ?", projectID, portscanTargetID).Delete(model.PortscanTarget{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) DeletePortscanTargetByPortscanSettingID(ctx context.Context, projectID uint32, portscanSettingID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", projectID, portscanSettingID).Delete(model.PortscanTarget{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *diagnosisRepository) ListAllPortscanSetting(ctx context.Context) (*[]model.PortscanSetting, error) {
	var data []model.PortscanSetting
	if err := r.SlaveDB.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func portscanSettingToMap(portscanSetting *model.PortscanSetting) map[string]interface{} {
	return map[string]interface{}{
		"portscan_setting_id":      portscanSetting.PortscanSettingID,
		"diagnosis_data_source_id": portscanSetting.DiagnosisDataSourceID,
		"project_id":               portscanSetting.ProjectID,
		"name":                     portscanSetting.Name,
	}
}

func portscanTargetToMap(portscanTarget *model.PortscanTarget) map[string]interface{} {
	settingMap := map[string]interface{}{
		"portscan_Target_id":  portscanTarget.PortscanTargetID,
		"portscan_setting_id": portscanTarget.PortscanSettingID,
		"project_id":          portscanTarget.ProjectID,
		"target":              portscanTarget.Target,
		"status":              portscanTarget.Status,
		"status_detail":       portscanTarget.StatusDetail,
	}
	if !zero.IsZeroVal(portscanTarget.ScanAt) {
		settingMap["scan_at"] = portscanTarget.ScanAt
	}
	return settingMap
}
