package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *DiagnosisService) ListJiraSetting(ctx context.Context, req *diagnosis.ListJiraSettingRequest) (*diagnosis.ListJiraSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListJiraSetting(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListJiraSettingResponse{}, nil
		}
		appLogger.Errorf("Failed to List JiraSettinng, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListJiraSettingResponse{}
	for _, d := range *list {
		data.JiraSetting = append(data.JiraSetting, convertJiraSetting(&d))
	}
	return &data, nil
}

func (s *DiagnosisService) GetJiraSetting(ctx context.Context, req *diagnosis.GetJiraSettingRequest) (*diagnosis.GetJiraSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetJiraSetting(ctx, req.ProjectId, req.JiraSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf("Failed to Get JiraSettinng, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetJiraSettingResponse{JiraSetting: convertJiraSetting(getData)}, nil
}

func (s *DiagnosisService) PutJiraSetting(ctx context.Context, req *diagnosis.PutJiraSettingRequest) (*diagnosis.PutJiraSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetJiraSetting(ctx, req.ProjectId, req.JiraSetting.JiraSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf("Failed to Get JiraSetting, error: %v", err)
		return nil, err
	}

	var jiraSettingID uint32
	if !noRecord {
		jiraSettingID = savedData.JiraSettingID
	}
	data := &model.JiraSetting{
		JiraSettingID:         jiraSettingID,
		ProjectID:             req.ProjectId,
		Name:                  req.JiraSetting.Name,
		DiagnosisDataSourceID: req.JiraSetting.DiagnosisDataSourceId,
		IdentityField:         req.JiraSetting.IdentityField,
		IdentityValue:         req.JiraSetting.IdentityValue,
		JiraID:                req.JiraSetting.JiraId,
		JiraKey:               req.JiraSetting.JiraKey,
		Status:                req.JiraSetting.Status.String(),
		StatusDetail:          req.JiraSetting.StatusDetail,
		ScanAt:                time.Unix(req.JiraSetting.ScanAt, 0),
	}

	registerdData, err := s.repository.UpsertJiraSetting(ctx, data)
	if err != nil {
		appLogger.Errorf("Failed to Put JiraSettinng, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutJiraSettingResponse{JiraSetting: convertJiraSetting(registerdData)}, nil
}

func (s *DiagnosisService) DeleteJiraSetting(ctx context.Context, req *diagnosis.DeleteJiraSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteJiraSetting(ctx, req.ProjectId, req.JiraSettingId); err != nil {
		appLogger.Errorf("Failed to Delete JiraSettinng, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertJiraSetting(data *model.JiraSetting) *diagnosis.JiraSetting {
	if data == nil {
		return &diagnosis.JiraSetting{}
	}
	return &diagnosis.JiraSetting{
		JiraSettingId:         data.JiraSettingID,
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		Name:                  data.Name,
		ProjectId:             data.ProjectID,
		IdentityField:         data.IdentityField,
		IdentityValue:         data.IdentityValue,
		JiraId:                data.JiraID,
		JiraKey:               data.JiraKey,
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
		Status:                getStatus(data.Status),
		StatusDetail:          data.StatusDetail,
		ScanAt:                data.ScanAt.Unix(),
	}
}

func makeJiraMessage(ProjectID, SettingID uint32, data *model.JiraSetting) (*message.JiraQueueMessage, error) {

	msg := &message.JiraQueueMessage{
		DataSource:    "diagnosis:jira",
		JiraSettingID: SettingID,
		ProjectID:     ProjectID,
		IdentityField: data.IdentityField,
		IdentityValue: data.IdentityValue,
		JiraID:        data.JiraID,
		JiraKey:       data.JiraKey,
	}
	return msg, nil
}

func getStatus(s string) diagnosis.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := diagnosis.Status_value[statusKey]; !ok {
		return diagnosis.Status_UNKNOWN
	}
	switch statusKey {
	case diagnosis.Status_OK.String():
		return diagnosis.Status_OK
	case diagnosis.Status_CONFIGURED.String():
		return diagnosis.Status_CONFIGURED
	case diagnosis.Status_IN_PROGRESS.String():
		return diagnosis.Status_IN_PROGRESS
	case diagnosis.Status_ERROR.String():
		return diagnosis.Status_ERROR
	default:
		return diagnosis.Status_UNKNOWN
	}
}
