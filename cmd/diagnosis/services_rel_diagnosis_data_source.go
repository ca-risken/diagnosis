package main

import (
	"context"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type relDiagnosisDataSourceService struct {
	repository diagnosisRepoInterface
	sqs        sqsAPI
}

func newRelDiagnosisDataSourceService(db diagnosisRepoInterface, s sqsAPI) diagnosis.RelDiagnosisDataSourceServiceServer {
	return &relDiagnosisDataSourceService{
		repository: db,
		sqs:        s,
	}
}

func (s *relDiagnosisDataSourceService) ListRelDiagnosisDataSource(ctx context.Context, req *diagnosis.ListRelDiagnosisDataSourceRequest) (*diagnosis.ListRelDiagnosisDataSourceResponse, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	list, err := s.repository.ListRelDiagnosisDataSource(req.ProjectId, req.DiagnosisId, req.DiagnosisDataSourceId, req.RecordId, req.JiraId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &diagnosis.ListRelDiagnosisDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := diagnosis.ListRelDiagnosisDataSourceResponse{}
	for _, d := range *list {
		data.RelDiagnosisDataSource = append(data.RelDiagnosisDataSource, convertRelDiagnosisDataSource(&d))
	}
	return &data, nil
}

func (s *relDiagnosisDataSourceService) GetRelDiagnosisDataSource(ctx context.Context, req *diagnosis.GetRelDiagnosisDataSourceRequest) (*diagnosis.GetRelDiagnosisDataSourceResponse, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	getData, err := s.repository.GetRelDiagnosisDataSource(req.ProjectId, req.RelDiagnosisDataSourceId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	return &diagnosis.GetRelDiagnosisDataSourceResponse{RelDiagnosisDataSource: convertRelDiagnosisDataSource(getData)}, nil
}

func (s *relDiagnosisDataSourceService) PutRelDiagnosisDataSource(ctx context.Context, req *diagnosis.PutRelDiagnosisDataSourceRequest) (*diagnosis.PutRelDiagnosisDataSourceResponse, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	savedData, err := s.repository.GetRelDiagnosisDataSource(req.ProjectId, req.RelDiagnosisDataSource.RelDiagnosisDataSourceId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	var relDiagnosisDataSourceID uint32
	if !noRecord {
		relDiagnosisDataSourceID = savedData.RelDiagnosisDataSourceID
	}
	data := &RelDiagnosisDataSource{
		RelDiagnosisDataSourceID: relDiagnosisDataSourceID,
		ProjectID:                req.ProjectId,
		DiagnosisDataSourceID:    req.RelDiagnosisDataSource.DiagnosisDataSourceId,
		DiagnosisID:              req.RelDiagnosisDataSource.DiagnosisId,
		RecordID:                 req.RelDiagnosisDataSource.RecordId,
		JiraID:                   req.RelDiagnosisDataSource.JiraId,
	}

	registerdData, err := s.repository.UpsertRelDiagnosisDataSource(data)
	if err != nil {
		return nil, err
	}
	return &diagnosis.PutRelDiagnosisDataSourceResponse{RelDiagnosisDataSource: convertRelDiagnosisDataSource(registerdData)}, nil
}

func (s *relDiagnosisDataSourceService) DeleteRelDiagnosisDataSource(ctx context.Context, req *diagnosis.DeleteRelDiagnosisDataSourceRequest) (*empty.Empty, error) {
	//	if err := req.Validate(); err != nil {
	//		return nil, err
	//	}
	if err := s.repository.DeleteRelDiagnosisDataSource(req.ProjectId, req.RelDiagnosisDataSourceId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertRelDiagnosisDataSource(data *RelDiagnosisDataSource) *diagnosis.RelDiagnosisDataSource {
	if data == nil {
		return &diagnosis.RelDiagnosisDataSource{}
	}
	return &diagnosis.RelDiagnosisDataSource{
		RelDiagnosisDataSourceId: data.RelDiagnosisDataSourceID,
		DiagnosisDataSourceId:    data.DiagnosisDataSourceID,
		DiagnosisId:              data.DiagnosisID,
		ProjectId:                data.ProjectID,
		RecordId:                 data.RecordID,
		JiraId:                   data.JiraID,
		CreatedAt:                data.CreatedAt.Unix(),
		UpdatedAt:                data.CreatedAt.Unix(),
	}
}

func (s *relDiagnosisDataSourceService) StartDiagnosis(ctx context.Context, req *diagnosis.StartDiagnosisRequest) (*diagnosis.StartDiagnosisResponse, error) {
	data, err := s.repository.GetRelDiagnosisDataSource(req.ProjectId, req.RelDiagnosisDataSourceId)
	if err != nil {
		return nil, err
	}
	msg := &message.DiagnosisQueueMessage{
		DataSource: "diagnosis:jira",
		ProjectID:  req.ProjectId,
		RecordID:   data.RecordID,
		JiraID:     data.JiraID,
	}
	resp, err := s.sqs.send(msg)
	if err != nil {
		return nil, err
	}
	logger.Info("Invoke scanned.", zap.String("MessageId", *resp.MessageId))
	return &diagnosis.StartDiagnosisResponse{Message: "Start Diagnosis."}, nil
}
