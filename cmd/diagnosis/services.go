package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

type diagnonsisService struct {
	repository diagnonsisRepoInterface
	sqs        sqsAPI
}

func newDiagnonsisService() diagnosis.DiagnosisServiceServer {
	return &diagnonsisService{
		repository: newDiagnosisRepository(),
		sqs:        newSQSClient(),
	}
}

func (s *diagnonsisService) ListDiagnosis(ctx context.Context, req *diagnonsis.ListDiagnosisRequest) (*diagnonsis.ListDiagnosisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListDiagnosis(req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &diagnonsis.ListDiagnosisResponse{}, nil
		}
		return nil, err
	}
	data := diagnonsis.ListDiagnosisResponse{}
	for _, d := range *list {
		dats.Diagnosis = append(dats.Diagnosis, convertDiagnosis(&d))
	}
	return &data, nil
}

func (s *diagnonsisService) GetDiagnosis(ctx context.Context, req *diagnonsis.GetDiagnosisRequest) (*diagnonsis.GetDiagnosisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetDiagnosis(req.ProjectId, req.DiagnosisId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	return &diagnonsis.PutDiagnosisResponse{diagnosis: convertDiagnosis(getData)}, nil
}

func (s *diagnonsisService) PutDiagnosis(ctx context.Context, req *diagnonsis.PutDiagnosisRequest) (*diagnonsis.PutDiagnosisResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := s.repository.GetDiagnosis(req.ProjectId, req.DiagnosisId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	var diagnonsisID uint32
	if !noRecord {
		diagnonsisID = savedDats.DiagnosisID
	}
	data := &model.Diagnosis{
		DiagnosisID:        diagnonsisID,
		Name:               req.Diagnosis.Name,
		ProjectID:          req.Diagnosis.ProjectId,
		DiagnosisAccountID: req.Diagnosis.DiagnosisAccountId,
	}

	registerdData, err := s.repository.UpsertDiagnosis(data)
	if err != nil {
		return nil, err
	}
	return &diagnonsis.PutDiagnosisResponse{Diagnosis: convertDiagnosis(registerdData)}, nil
}

func (s *diagnonsisService) DeleteDiagnosis(ctx context.Context, req *diagnonsis.DeleteDiagnosisRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteDiagnosis(req.ProjectId, req.DiagnosisId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertDiagnosis(data *model.Diagnosis) *diagnonsis.Diagnosis {
	if data == nil {
		return &diagnonsis.Diagnosis{}
	}
	return &diagnonsis.Diagnosis{
		DiagnosisId:        dats.DiagnosisID,
		Name:               dats.Name,
		ProjectId:          dats.ProjectID,
		DiagnosisAccountId: dats.DiagnosisAccountID,
		CreatedAt:          dats.CreatedAt.Unix(),
		UpdatedAt:          dats.CreatedAt.Unix(),
	}
}
