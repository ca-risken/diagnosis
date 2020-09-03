package main

import (
	"context"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (s *diagnosisService) ListDiagnosis(ctx context.Context, req *diagnosis.ListDiagnosisRequest) (*diagnosis.ListDiagnosisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListDiagnosis(req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &diagnosis.ListDiagnosisResponse{}, nil
		}
		return nil, err
	}
	data := diagnosis.ListDiagnosisResponse{}
	for _, d := range *list {
		data.Diagnosis = append(data.Diagnosis, convertDiagnosis(&d))
	}
	return &data, nil
}

func (s *diagnosisService) GetDiagnosis(ctx context.Context, req *diagnosis.GetDiagnosisRequest) (*diagnosis.GetDiagnosisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetDiagnosis(req.ProjectId, req.DiagnosisId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	return &diagnosis.GetDiagnosisResponse{Diagnosis: convertDiagnosis(getData)}, nil
}

func (s *diagnosisService) PutDiagnosis(ctx context.Context, req *diagnosis.PutDiagnosisRequest) (*diagnosis.PutDiagnosisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetDiagnosis(req.ProjectId, req.Diagnosis.DiagnosisId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	var diagnosisID uint32
	if !noRecord {
		diagnosisID = savedData.DiagnosisID
	}
	data := &Diagnosis{
		DiagnosisID: diagnosisID,
		Name:        req.Diagnosis.Name,
		ProjectID:   req.ProjectId,
	}

	registerdData, err := s.repository.UpsertDiagnosis(data)
	if err != nil {
		return nil, err
	}
	return &diagnosis.PutDiagnosisResponse{Diagnosis: convertDiagnosis(registerdData)}, nil
}

func (s *diagnosisService) DeleteDiagnosis(ctx context.Context, req *diagnosis.DeleteDiagnosisRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteDiagnosis(req.ProjectId, req.DiagnosisId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertDiagnosis(data *Diagnosis) *diagnosis.Diagnosis {
	if data == nil {
		return &diagnosis.Diagnosis{}
	}
	return &diagnosis.Diagnosis{
		DiagnosisId: data.DiagnosisID,
		Name:        data.Name,
		ProjectId:   data.ProjectID,
		CreatedAt:   data.CreatedAt.Unix(),
		UpdatedAt:   data.CreatedAt.Unix(),
	}
}
