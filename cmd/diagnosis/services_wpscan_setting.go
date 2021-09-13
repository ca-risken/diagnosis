package main

import (
	"context"
	"errors"
	"time"

	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *diagnosisService) ListWpscanSetting(ctx context.Context, req *diagnosis.ListWpscanSettingRequest) (*diagnosis.ListWpscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListWpscanSetting(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListWpscanSettingResponse{}, nil
		}
		logger.Error("Failed to List WpscanSettinng", zap.Error(err))
		return nil, err
	}
	data := diagnosis.ListWpscanSettingResponse{}
	for _, d := range *list {
		data.WpscanSetting = append(data.WpscanSetting, convertWpscanSetting(&d))
	}
	return &data, nil
}

func (s *diagnosisService) GetWpscanSetting(ctx context.Context, req *diagnosis.GetWpscanSettingRequest) (*diagnosis.GetWpscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetWpscanSetting(ctx, req.ProjectId, req.WpscanSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		logger.Error("Failed to Get WpscanSettinng", zap.Error(err))
		return nil, err
	}

	return &diagnosis.GetWpscanSettingResponse{WpscanSetting: convertWpscanSetting(getData)}, nil
}

func (s *diagnosisService) PutWpscanSetting(ctx context.Context, req *diagnosis.PutWpscanSettingRequest) (*diagnosis.PutWpscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetWpscanSetting(ctx, req.ProjectId, req.WpscanSetting.WpscanSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		logger.Error("Failed to Get WpscanSetting", zap.Error(err))
		return nil, err
	}

	var jiraSettingID uint32
	if !noRecord {
		jiraSettingID = savedData.WpscanSettingID
	}
	data := &model.WpscanSetting{
		WpscanSettingID:       jiraSettingID,
		ProjectID:             req.ProjectId,
		DiagnosisDataSourceID: req.WpscanSetting.DiagnosisDataSourceId,
		TargetURL:             req.WpscanSetting.TargetUrl,
		Status:                req.WpscanSetting.Status.String(),
		StatusDetail:          req.WpscanSetting.StatusDetail,
		ScanAt:                time.Unix(req.WpscanSetting.ScanAt, 0),
	}

	registerdData, err := s.repository.UpsertWpscanSetting(ctx, data)
	if err != nil {
		logger.Error("Failed to Put WpscanSetting", zap.Error(err))
		return nil, err
	}
	return &diagnosis.PutWpscanSettingResponse{WpscanSetting: convertWpscanSetting(registerdData)}, nil
}

func (s *diagnosisService) DeleteWpscanSetting(ctx context.Context, req *diagnosis.DeleteWpscanSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteWpscanSetting(ctx, req.ProjectId, req.WpscanSettingId); err != nil {
		logger.Error("Failed to Delete WpscanSettinng", zap.Error(err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertWpscanSetting(data *model.WpscanSetting) *diagnosis.WpscanSetting {
	if data == nil {
		return &diagnosis.WpscanSetting{}
	}
	return &diagnosis.WpscanSetting{
		WpscanSettingId:       data.WpscanSettingID,
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		ProjectId:             data.ProjectID,
		TargetUrl:             data.TargetURL,
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
		Status:                getStatus(data.Status),
		StatusDetail:          data.StatusDetail,
		ScanAt:                data.ScanAt.Unix(),
	}
}

func makeWpscanMessage(ProjectID, SettingID uint32, targetURL string) (*message.WpscanQueueMessage, error) {
	msg := &message.WpscanQueueMessage{
		DataSource:      "diagnosis:wpscan",
		WpscanSettingID: SettingID,
		ProjectID:       ProjectID,
		TargetURL:       targetURL,
	}
	return msg, nil
}
