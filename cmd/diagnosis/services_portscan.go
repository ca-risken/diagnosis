package main

import (
	"context"
	"errors"
	"time"

	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *DiagnosisService) ListPortscanSetting(ctx context.Context, req *diagnosis.ListPortscanSettingRequest) (*diagnosis.ListPortscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListPortscanSetting(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListPortscanSettingResponse{}, nil
		}
		appLogger.Errorf(ctx, "Failed to List PortscanSetting, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListPortscanSettingResponse{}
	for _, d := range *list {
		data.PortscanSetting = append(data.PortscanSetting, convertPortscanSetting(&d))
	}
	return &data, nil
}

func (s *DiagnosisService) GetPortscanSetting(ctx context.Context, req *diagnosis.GetPortscanSettingRequest) (*diagnosis.GetPortscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetPortscanSetting(ctx, req.ProjectId, req.PortscanSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf(ctx, "Failed to Get PortscanSetting, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetPortscanSettingResponse{PortscanSetting: convertPortscanSetting(getData)}, nil
}

func (s *DiagnosisService) PutPortscanSetting(ctx context.Context, req *diagnosis.PutPortscanSettingRequest) (*diagnosis.PutPortscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetPortscanSetting(ctx, req.ProjectId, req.PortscanSetting.PortscanSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf(ctx, "Failed to Get PortscanSetting, error: %v", err)
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
	}

	registerdData, err := s.repository.UpsertPortscanSetting(ctx, data)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to Put PortscanSetting, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutPortscanSettingResponse{PortscanSetting: convertPortscanSetting(registerdData)}, nil
}

func (s *DiagnosisService) DeletePortscanSetting(ctx context.Context, req *diagnosis.DeletePortscanSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeletePortscanSetting(ctx, req.ProjectId, req.PortscanSettingId); err != nil {
		appLogger.Errorf(ctx, "Failed to Delete PortscanSetting, error: %v", err)
		return nil, err
	}
	// Delete PortscanTargetBySetting
	if err := s.repository.DeletePortscanTargetByPortscanSettingID(ctx, req.ProjectId, req.PortscanSettingId); err != nil {
		appLogger.Errorf(ctx, "Failed to Delete PortscanTargetByPortscanSettingID, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *DiagnosisService) ListPortscanTarget(ctx context.Context, req *diagnosis.ListPortscanTargetRequest) (*diagnosis.ListPortscanTargetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListPortscanTarget(ctx, req.ProjectId, req.PortscanSettingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListPortscanTargetResponse{}, nil
		}
		appLogger.Errorf(ctx, "Failed to List PortscanTarget, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListPortscanTargetResponse{}
	for _, d := range *list {
		data.PortscanTarget = append(data.PortscanTarget, convertPortscanTarget(&d))
	}
	return &data, nil
}

func (s *DiagnosisService) GetPortscanTarget(ctx context.Context, req *diagnosis.GetPortscanTargetRequest) (*diagnosis.GetPortscanTargetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetPortscanTarget(ctx, req.ProjectId, req.PortscanTargetId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf(ctx, "Failed to Get PortscanTarget, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetPortscanTargetResponse{PortscanTarget: convertPortscanTarget(getData)}, nil
}

func (s *DiagnosisService) PutPortscanTarget(ctx context.Context, req *diagnosis.PutPortscanTargetRequest) (*diagnosis.PutPortscanTargetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.PortscanTarget{
		PortscanTargetID:  req.PortscanTarget.PortscanTargetId,
		ProjectID:         req.ProjectId,
		PortscanSettingID: req.PortscanTarget.PortscanSettingId,
		Target:            req.PortscanTarget.Target,
		Status:            req.PortscanTarget.Status.String(),
		StatusDetail:      req.PortscanTarget.StatusDetail,
		ScanAt:            time.Unix(req.PortscanTarget.ScanAt, 0),
	}

	registerdData, err := s.repository.UpsertPortscanTarget(ctx, data)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to Put PortscanTarget, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutPortscanTargetResponse{PortscanTarget: convertPortscanTarget(registerdData)}, nil
}

func (s *DiagnosisService) DeletePortscanTarget(ctx context.Context, req *diagnosis.DeletePortscanTargetRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeletePortscanTarget(ctx, req.ProjectId, req.PortscanTargetId); err != nil {
		appLogger.Errorf(ctx, "Failed to Delete PortscanTarget, error: %v", err)
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
		Status:            getStatus(data.Status),
		StatusDetail:      data.StatusDetail,
		ScanAt:            data.ScanAt.Unix(),
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.CreatedAt.Unix(),
	}
}

func makePortscanMessage(projectID, settingID, portscanTargetID uint32, target string) (*message.PortscanQueueMessage, error) {
	msg := &message.PortscanQueueMessage{
		DataSource:        common.DataSourceNamePortScan,
		PortscanSettingID: settingID,
		PortscanTargetID:  portscanTargetID,
		ProjectID:         projectID,
		Target:            target,
	}
	return msg, nil
}
