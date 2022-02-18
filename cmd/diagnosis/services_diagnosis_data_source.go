package main

import (
	"context"
	"errors"

	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *DiagnosisService) ListDiagnosisDataSource(ctx context.Context, req *diagnosis.ListDiagnosisDataSourceRequest) (*diagnosis.ListDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListDiagnosisDataSource(ctx, req.ProjectId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListDiagnosisDataSourceResponse{}, nil
		}
		appLogger.Errorf("Failed to List DiagnosisDataSource, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListDiagnosisDataSourceResponse{}
	for _, d := range *list {
		data.DiagnosisDataSource = append(data.DiagnosisDataSource, convertDiagnosisDataSource(&d))
	}
	return &data, nil
}

func (s *DiagnosisService) GetDiagnosisDataSource(ctx context.Context, req *diagnosis.GetDiagnosisDataSourceRequest) (*diagnosis.GetDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetDiagnosisDataSource(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf("Failed to Get DiagnosisDataSource, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(getData)}, nil
}

func (s *DiagnosisService) PutDiagnosisDataSource(ctx context.Context, req *diagnosis.PutDiagnosisDataSourceRequest) (*diagnosis.PutDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetDiagnosisDataSource(ctx, req.ProjectId, req.DiagnosisDataSource.DiagnosisDataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf("Failed to Get DiagnosisDataSource, error: %v", err)
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

	registerdData, err := s.repository.UpsertDiagnosisDataSource(ctx, data)
	if err != nil {
		appLogger.Errorf("Failed to Put DiagnosisDataSource, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(registerdData)}, nil
}

func (s *DiagnosisService) DeleteDiagnosisDataSource(ctx context.Context, req *diagnosis.DeleteDiagnosisDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteDiagnosisDataSource(ctx, req.ProjectId, req.DiagnosisDataSourceId); err != nil {
		appLogger.Errorf("Failed to Delete DiagnosisDataSource, error: %v", err)
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
