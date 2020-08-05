package main

import (
	"context"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

type diagnosisDataSourceService struct {
	repository diagnosisRepoInterface
	//	sqs sqsAPI
}

func newDiagnosisDataSourceService(db diagnosisRepoInterface) diagnosis.DiagnosisDataSourceServiceServer {
	return &diagnosisDataSourceService{
		repository: db,
		//		sqs:        newSQSClient(),
	}
}

func (s *diagnosisDataSourceService) ListDiagnosisDataSource(ctx context.Context, req *diagnosis.ListDiagnosisDataSourceRequest) (*diagnosis.ListDiagnosisDataSourceResponse, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	list, err := s.repository.ListDiagnosisDataSource(req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &diagnosis.ListDiagnosisDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := diagnosis.ListDiagnosisDataSourceResponse{}
	for _, d := range *list {
		data.DiagnosisDataSource = append(data.DiagnosisDataSource, convertDiagnosisDataSource(&d))
	}
	return &data, nil
}

func (s *diagnosisDataSourceService) GetDiagnosisDataSource(ctx context.Context, req *diagnosis.GetDiagnosisDataSourceRequest) (*diagnosis.GetDiagnosisDataSourceResponse, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	getData, err := s.repository.GetDiagnosisDataSource(req.ProjectId, req.DiagnosisDataSourceId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	return &diagnosis.GetDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(getData)}, nil
}

func (s *diagnosisDataSourceService) PutDiagnosisDataSource(ctx context.Context, req *diagnosis.PutDiagnosisDataSourceRequest) (*diagnosis.PutDiagnosisDataSourceResponse, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	savedData, err := s.repository.GetDiagnosisDataSource(req.ProjectId, req.DiagnosisDataSource.DiagnosisDataSourceId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	var diagnosisDataSourceID uint32
	if !noRecord {
		diagnosisDataSourceID = savedData.DiagnosisDataSourceID
	}
	data := &DiagnosisDataSource{
		DiagnosisDataSourceID: diagnosisDataSourceID,
		Name:                  req.DiagnosisDataSource.Name,
		Description:           req.DiagnosisDataSource.Description,
		MaxScore:              req.DiagnosisDataSource.MaxScore,
	}

	registerdData, err := s.repository.UpsertDiagnosisDataSource(data)
	if err != nil {
		return nil, err
	}
	return &diagnosis.PutDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(registerdData)}, nil
}

func (s *diagnosisDataSourceService) DeleteDiagnosisDataSource(ctx context.Context, req *diagnosis.DeleteDiagnosisDataSourceRequest) (*empty.Empty, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	if err := s.repository.DeleteDiagnosisDataSource(req.ProjectId, req.DiagnosisDataSourceId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertDiagnosisDataSource(data *DiagnosisDataSource) *diagnosis.DiagnosisDataSource {
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
