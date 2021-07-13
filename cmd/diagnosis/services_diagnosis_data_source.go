package main

import (
	"context"
	"errors"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *diagnosisService) ListDiagnosisDataSource(ctx context.Context, req *diagnosis.ListDiagnosisDataSourceRequest) (*diagnosis.ListDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListDiagnosisDataSource(req.ProjectId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListDiagnosisDataSourceResponse{}, nil
		}
		logger.Error("Failed to List DiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	data := diagnosis.ListDiagnosisDataSourceResponse{}
	for _, d := range *list {
		data.DiagnosisDataSource = append(data.DiagnosisDataSource, convertDiagnosisDataSource(&d))
	}
	return &data, nil
}

func (s *diagnosisService) GetDiagnosisDataSource(ctx context.Context, req *diagnosis.GetDiagnosisDataSourceRequest) (*diagnosis.GetDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetDiagnosisDataSource(req.ProjectId, req.DiagnosisDataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		logger.Error("Failed to Get DiagnosisDataSource", zap.Error(err))
		return nil, err
	}

	return &diagnosis.GetDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(getData)}, nil
}

func (s *diagnosisService) PutDiagnosisDataSource(ctx context.Context, req *diagnosis.PutDiagnosisDataSourceRequest) (*diagnosis.PutDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetDiagnosisDataSource(req.ProjectId, req.DiagnosisDataSource.DiagnosisDataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		logger.Error("Failed to Get DiagnosisDataSource", zap.Error(err))
		return nil, err
	}

	var diagnosisDataSourceID uint32
	if !noRecord {
		diagnosisDataSourceID = savedData.DiagnosisDataSourceID
	}
	data := &model.DiagnosisDataSource{
		DiagnosisDataSourceID: diagnosisDataSourceID,
		Name:                  req.DiagnosisDataSource.Name,
		Description:           req.DiagnosisDataSource.Description,
		MaxScore:              req.DiagnosisDataSource.MaxScore,
	}

	registerdData, err := s.repository.UpsertDiagnosisDataSource(data)
	if err != nil {
		logger.Error("Failed to Put DiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &diagnosis.PutDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(registerdData)}, nil
}

func (s *diagnosisService) DeleteDiagnosisDataSource(ctx context.Context, req *diagnosis.DeleteDiagnosisDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteDiagnosisDataSource(req.ProjectId, req.DiagnosisDataSourceId); err != nil {
		logger.Error("Failed to Delete DiagnosisDataSource", zap.Error(err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertDiagnosisDataSource(data *model.DiagnosisDataSource) *diagnosis.DiagnosisDataSource {
	if data == nil {
		return &diagnosis.DiagnosisDataSource{}
	}
	return &diagnosis.DiagnosisDataSource{
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		Name:                  data.Name,
		Description:           data.Description,
		MaxScore:              data.MaxScore,
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
	}
}
