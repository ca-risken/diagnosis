package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *DiagnosisService) ListApplicationScan(ctx context.Context, req *diagnosis.ListApplicationScanRequest) (*diagnosis.ListApplicationScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListApplicationScan(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListApplicationScanResponse{}, nil
		}
		appLogger.Errorf(ctx, "Failed to List ApplicationScan, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListApplicationScanResponse{}
	for _, d := range *list {
		data.ApplicationScan = append(data.ApplicationScan, convertApplicationScan(&d))
	}
	return &data, nil
}

func (s *DiagnosisService) GetApplicationScan(ctx context.Context, req *diagnosis.GetApplicationScanRequest) (*diagnosis.GetApplicationScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetApplicationScan(ctx, req.ProjectId, req.ApplicationScanId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf(ctx, "Failed to Get ApplicationScan, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetApplicationScanResponse{ApplicationScan: convertApplicationScan(getData)}, nil
}

func (s *DiagnosisService) PutApplicationScan(ctx context.Context, req *diagnosis.PutApplicationScanRequest) (*diagnosis.PutApplicationScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetApplicationScan(ctx, req.ProjectId, req.ApplicationScan.ApplicationScanId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf(ctx, "Failed to Get ApplicationScan, error: %v", err)
		return nil, err
	}

	var applicationScanID uint32
	if !noRecord {
		applicationScanID = savedData.ApplicationScanID
	}
	data := &model.ApplicationScan{
		ApplicationScanID:     applicationScanID,
		ProjectID:             req.ProjectId,
		DiagnosisDataSourceID: req.ApplicationScan.DiagnosisDataSourceId,
		Name:                  req.ApplicationScan.Name,
		ScanType:              req.ApplicationScan.ScanType.String(),
		Status:                req.ApplicationScan.Status.String(),
		StatusDetail:          req.ApplicationScan.StatusDetail,
		ScanAt:                time.Unix(req.ApplicationScan.ScanAt, 0),
	}

	registerdData, err := s.repository.UpsertApplicationScan(ctx, data)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to Put ApplicationScan, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutApplicationScanResponse{ApplicationScan: convertApplicationScan(registerdData)}, nil
}

func (s *DiagnosisService) DeleteApplicationScan(ctx context.Context, req *diagnosis.DeleteApplicationScanRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteApplicationScan(ctx, req.ProjectId, req.ApplicationScanId); err != nil {
		appLogger.Errorf(ctx, "Failed to Delete ApplicationScan, error: %v", err)
		return nil, err
	}
	// Delete ApplicationScanBasicSetting
	//	if err := s.repository.DeleteApplicationScanBasicSettingByApplicationScanID(ctx, req.ProjectId, req.ApplicationScanId); err != nil {
	//		appLogger.Errorf("Failed to Delete ApplicationScanBasicSettingByApplicationScanID, error: %v", err)
	//		return nil, err
	//	}
	return &empty.Empty{}, nil
}

func (s *DiagnosisService) ListApplicationScanBasicSetting(ctx context.Context, req *diagnosis.ListApplicationScanBasicSettingRequest) (*diagnosis.ListApplicationScanBasicSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListApplicationScanBasicSetting(ctx, req.ProjectId, req.ApplicationScanId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListApplicationScanBasicSettingResponse{}, nil
		}
		appLogger.Errorf(ctx, "Failed to List ApplicationScanBasicSetting, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListApplicationScanBasicSettingResponse{}
	for _, d := range *list {
		data.ApplicationScanBasicSetting = append(data.ApplicationScanBasicSetting, convertApplicationScanBasicSetting(&d))
	}
	return &data, nil
}

func (s *DiagnosisService) GetApplicationScanBasicSetting(ctx context.Context, req *diagnosis.GetApplicationScanBasicSettingRequest) (*diagnosis.GetApplicationScanBasicSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetApplicationScanBasicSetting(ctx, req.ProjectId, req.ApplicationScanId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf(ctx, "Failed to Get ApplicationScanBasicSetting, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetApplicationScanBasicSettingResponse{ApplicationScanBasicSetting: convertApplicationScanBasicSetting(getData)}, nil
}

func (s *DiagnosisService) PutApplicationScanBasicSetting(ctx context.Context, req *diagnosis.PutApplicationScanBasicSettingRequest) (*diagnosis.PutApplicationScanBasicSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.ApplicationScanBasicSetting{
		ApplicationScanBasicSettingID: req.ApplicationScanBasicSetting.ApplicationScanBasicSettingId,
		ProjectID:                     req.ProjectId,
		ApplicationScanID:             req.ApplicationScanBasicSetting.ApplicationScanId,
		Target:                        req.ApplicationScanBasicSetting.Target,
		MaxDepth:                      req.ApplicationScanBasicSetting.MaxDepth,
		MaxChildren:                   req.ApplicationScanBasicSetting.MaxChildren,
	}

	registerdData, err := s.repository.UpsertApplicationScanBasicSetting(ctx, data)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to Put ApplicationScanBasicSetting, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutApplicationScanBasicSettingResponse{ApplicationScanBasicSetting: convertApplicationScanBasicSetting(registerdData)}, nil
}

func (s *DiagnosisService) DeleteApplicationScanBasicSetting(ctx context.Context, req *diagnosis.DeleteApplicationScanBasicSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteApplicationScanBasicSetting(ctx, req.ProjectId, req.ApplicationScanBasicSettingId); err != nil {
		appLogger.Errorf(ctx, "Failed to Delete ApplicationScanBasicSetting, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertApplicationScan(data *model.ApplicationScan) *diagnosis.ApplicationScan {
	if data == nil {
		return &diagnosis.ApplicationScan{}
	}
	return &diagnosis.ApplicationScan{
		ApplicationScanId:     data.ApplicationScanID,
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		ProjectId:             data.ProjectID,
		Name:                  data.Name,
		ScanType:              getScanType(data.ScanType),
		Status:                getStatus(data.Status),
		StatusDetail:          data.StatusDetail,
		ScanAt:                data.ScanAt.Unix(),
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
	}
}

func convertApplicationScanBasicSetting(data *model.ApplicationScanBasicSetting) *diagnosis.ApplicationScanBasicSetting {
	if data == nil {
		return &diagnosis.ApplicationScanBasicSetting{}
	}
	return &diagnosis.ApplicationScanBasicSetting{
		ApplicationScanBasicSettingId: data.ApplicationScanBasicSettingID,
		ApplicationScanId:             data.ApplicationScanID,
		ProjectId:                     data.ProjectID,
		Target:                        data.Target,
		MaxDepth:                      data.MaxDepth,
		MaxChildren:                   data.MaxChildren,
		CreatedAt:                     data.CreatedAt.Unix(),
		UpdatedAt:                     data.CreatedAt.Unix(),
	}
}

func getScanType(s string) diagnosis.ApplicationScanType {
	statusKey := strings.ToUpper(s)
	if _, ok := diagnosis.ApplicationScanType_value[statusKey]; !ok {
		return diagnosis.ApplicationScanType_NOT_CONFIGURED
	}
	switch statusKey {
	case diagnosis.ApplicationScanType_BASIC.String():
		return diagnosis.ApplicationScanType_BASIC
	default:
		return diagnosis.ApplicationScanType_NOT_CONFIGURED
	}
}

func makeApplicationScanMessage(projectID, applicationScanID uint32, name, scanType string) (*message.ApplicationScanQueueMessage, error) {
	msg := &message.ApplicationScanQueueMessage{
		DataSource:          common.DataSourceNameApplicationScan,
		ApplicationScanID:   applicationScanID,
		ProjectID:           projectID,
		Name:                name,
		ApplicationScanType: scanType,
	}
	return msg, nil
}
