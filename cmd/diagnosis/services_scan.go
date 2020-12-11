package main

import (
	"context"

	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func (s *diagnosisService) InvokeScan(ctx context.Context, req *diagnosis.InvokeScanRequest) (*diagnosis.InvokeScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	dataSource, err := s.repository.GetDiagnosisDataSource(0, req.DiagnosisDataSourceId)
	if err != nil {
		return nil, err
	}

	var resp *sqs.SendMessageOutput
	switch dataSource.Name {
	case "diagnosis:jira":
		data, err := s.repository.GetJiraSetting(req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		msg, err := makeJiraMessage(req.ProjectId, req.SettingId, data)
		resp, err = s.sqs.send(msg)
		if err != nil {
			return nil, err
		}
	case "diagnosis:wpscan":
		data, err := s.repository.GetWpscanSetting(req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		msg, err := makeWpscanMessage(req.ProjectId, req.SettingId, data.TargetURL)
		if err != nil {
			return nil, err
		}
		resp, err = s.sqs.sendWpscanMessage(msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	logger.Info("Invoke scanned.", zap.String("MessageId", *resp.MessageId))
	return &diagnosis.InvokeScanResponse{Message: "Start Diagnosis."}, nil
}

func (s *diagnosisService) InvokeScanAll(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {

	listjiraSetting, err := s.repository.ListAllJiraSetting()
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &empty.Empty{}, nil
		}
		logger.Error("Failed to List All JiraSetting.", zap.Error(err))
		return nil, err
	}

	for _, jiraSetting := range *listjiraSetting {
		if _, err := s.InvokeScan(ctx, &diagnosis.InvokeScanRequest{
			ProjectId:             jiraSetting.ProjectID,
			SettingId:             jiraSetting.JiraSettingID,
			DiagnosisDataSourceId: jiraSetting.DiagnosisDataSourceID,
		}); err != nil {
			// errorが出ても続行
			logger.Error("InvokeScanAll error", zap.Error(err))
		}
	}

	listWpscanSetting, err := s.repository.ListAllWpscanSetting()
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &empty.Empty{}, nil
		}
		logger.Error("Failed to List All WPScanSetting.", zap.Error(err))
		return nil, err
	}

	for _, WpscanSetting := range *listWpscanSetting {
		if _, err := s.InvokeScan(ctx, &diagnosis.InvokeScanRequest{
			ProjectId:             WpscanSetting.ProjectID,
			SettingId:             WpscanSetting.WpscanSettingID,
			DiagnosisDataSourceId: WpscanSetting.DiagnosisDataSourceID,
		}); err != nil {
			// errorが出ても続行
			logger.Error("InvokeScanAll error", zap.Error(err))
		}
	}

	return &empty.Empty{}, nil
}
