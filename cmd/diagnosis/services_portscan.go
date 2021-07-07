package main

import (
	"context"
	"time"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func (s *diagnosisService) ListPortscanSetting(ctx context.Context, req *diagnosis.ListPortscanSettingRequest) (*diagnosis.ListPortscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListPortscanSetting(req.ProjectId, req.DiagnosisDataSourceId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &diagnosis.ListPortscanSettingResponse{}, nil
		}
		logger.Error("Failed to List PortscanSetting", zap.Error(err))
		return nil, err
	}
	data := diagnosis.ListPortscanSettingResponse{}
	for _, d := range *list {
		data.PortscanSetting = append(data.PortscanSetting, convertPortscanSetting(&d))
	}
	return &data, nil
}

func (s *diagnosisService) GetPortscanSetting(ctx context.Context, req *diagnosis.GetPortscanSettingRequest) (*diagnosis.GetPortscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetPortscanSetting(req.ProjectId, req.PortscanSettingId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		logger.Error("Failed to Get PortscanSetting", zap.Error(err))
		return nil, err
	}

	return &diagnosis.GetPortscanSettingResponse{PortscanSetting: convertPortscanSetting(getData)}, nil
}

func (s *diagnosisService) PutPortscanSetting(ctx context.Context, req *diagnosis.PutPortscanSettingRequest) (*diagnosis.PutPortscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetPortscanSetting(req.ProjectId, req.PortscanSetting.PortscanSettingId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		logger.Error("Failed to Get PortscanSetting", zap.Error(err))
		return nil, err
	}

	var portscanSettingID uint32
	if !noRecord {
		portscanSettingID = savedData.PortscanSettingID
	}
	data := &model.PortscanSetting{
		PortscanSettingID:     portscanSettingID,
		ProjectID:             req.ProjectId,
		DiagnosisDataSourceID: req.PortscanSetting.DiagnosisDataSourceId,
		Name:                  req.PortscanSetting.Name,
		Status:                req.PortscanSetting.Status.String(),
		StatusDetail:          req.PortscanSetting.StatusDetail,
		ScanAt:                time.Unix(req.PortscanSetting.ScanAt, 0),
	}

	registerdData, err := s.repository.UpsertPortscanSetting(data)
	if err != nil {
		logger.Error("Failed to Put PortscanSetting", zap.Error(err))
		return nil, err
	}
	return &diagnosis.PutPortscanSettingResponse{PortscanSetting: convertPortscanSetting(registerdData)}, nil
}

func (s *diagnosisService) DeletePortscanSetting(ctx context.Context, req *diagnosis.DeletePortscanSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeletePortscanSetting(req.ProjectId, req.PortscanSettingId); err != nil {
		logger.Error("Failed to Delete PortscanSetting", zap.Error(err))
		return nil, err
	}
	// Delete PortscanTargetBySetting
	if err := s.repository.DeletePortscanTargetByPortscanSettingID(req.ProjectId, req.PortscanSettingId); err != nil {
		logger.Error("Failed to Delete PortscanTargetByPortscanSettingID", zap.Error(err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *diagnosisService) ListPortscanTarget(ctx context.Context, req *diagnosis.ListPortscanTargetRequest) (*diagnosis.ListPortscanTargetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListPortscanTarget(req.ProjectId, req.PortscanSettingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &diagnosis.ListPortscanTargetResponse{}, nil
		}
		logger.Error("Failed to List PortscanTarget", zap.Error(err))
		return nil, err
	}
	data := diagnosis.ListPortscanTargetResponse{}
	for _, d := range *list {
		data.PortscanTarget = append(data.PortscanTarget, convertPortscanTarget(&d))
	}
	return &data, nil
}

func (s *diagnosisService) GetPortscanTarget(ctx context.Context, req *diagnosis.GetPortscanTargetRequest) (*diagnosis.GetPortscanTargetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetPortscanTarget(req.ProjectId, req.PortscanTargetId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		logger.Error("Failed to Get PortscanTarget", zap.Error(err))
		return nil, err
	}

	return &diagnosis.GetPortscanTargetResponse{PortscanTarget: convertPortscanTarget(getData)}, nil
}

func (s *diagnosisService) PutPortscanTarget(ctx context.Context, req *diagnosis.PutPortscanTargetRequest) (*diagnosis.PutPortscanTargetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetPortscanTargetByTargetPortscanSettingID(req.ProjectId, req.PortscanTarget.PortscanSettingId, req.PortscanTarget.Target)
	noRecord := gorm.IsRecordNotFoundError(err)
	logger.Info("hoge", zap.Error(err), zap.Any("noRecord", noRecord))
	if err != nil && !noRecord {
		logger.Error("Failed to Get PortscanTarget", zap.Error(err))
		return nil, err
	}
	if !noRecord {
		return &diagnosis.PutPortscanTargetResponse{PortscanTarget: convertPortscanTarget(savedData)}, nil
	}

	var portscanSettingID uint32
	if !noRecord {
		portscanSettingID = savedData.PortscanTargetID
	}
	data := &model.PortscanTarget{
		PortscanTargetID:  portscanSettingID,
		ProjectID:         req.ProjectId,
		PortscanSettingID: req.PortscanTarget.PortscanSettingId,
		Target:            req.PortscanTarget.Target,
	}

	registerdData, err := s.repository.UpsertPortscanTarget(data)
	if err != nil {
		logger.Error("Failed to Put PortscanTarget", zap.Error(err))
		return nil, err
	}
	return &diagnosis.PutPortscanTargetResponse{PortscanTarget: convertPortscanTarget(registerdData)}, nil
}

func (s *diagnosisService) DeletePortscanTarget(ctx context.Context, req *diagnosis.DeletePortscanTargetRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeletePortscanTarget(req.ProjectId, req.PortscanTargetId); err != nil {
		logger.Error("Failed to Delete PortscanTarget", zap.Error(err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertPortscanSetting(data *model.PortscanSetting) *diagnosis.PortscanSetting {
	if data == nil {
		return &diagnosis.PortscanSetting{}
	}
	return &diagnosis.PortscanSetting{
		PortscanSettingId:     data.PortscanSettingID,
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		ProjectId:             data.ProjectID,
		Name:                  data.Name,
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
		Status:                getStatus(data.Status),
		StatusDetail:          data.StatusDetail,
		ScanAt:                data.ScanAt.Unix(),
	}
}

func convertPortscanTarget(data *model.PortscanTarget) *diagnosis.PortscanTarget {
	if data == nil {
		return &diagnosis.PortscanTarget{}
	}
	return &diagnosis.PortscanTarget{
		PortscanTargetId:  data.PortscanTargetID,
		PortscanSettingId: data.PortscanSettingID,
		ProjectId:         data.ProjectID,
		Target:            data.Target,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.CreatedAt.Unix(),
	}
}

func makePortscanMessage(projectID, settingID uint32, name string) (*message.PortscanQueueMessage, error) {
	msg := &message.PortscanQueueMessage{
		DataSource:        "diagnosis:portscan",
		PortscanSettingID: settingID,
		ProjectID:         projectID,
		Name:              name,
	}
	return msg, nil
}
