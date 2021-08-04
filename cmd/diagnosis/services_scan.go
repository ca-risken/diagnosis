package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/model"
	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (d *diagnosisService) InvokeScan(ctx context.Context, req *diagnosis.InvokeScanRequest) (*diagnosis.InvokeScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	dataSource, err := d.repository.GetDiagnosisDataSource(ctx, 0, req.DiagnosisDataSourceId)
	if err != nil {
		return nil, err
	}
	var resp *sqs.SendMessageOutput
	switch dataSource.Name {
	case "diagnosis:jira":
		data, err := d.repository.GetJiraSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		msg, err := makeJiraMessage(req.ProjectId, req.SettingId, data)
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.send(ctx, msg)
		if err != nil {
			return nil, err
		}
		if _, err = d.repository.UpsertJiraSetting(ctx, &model.JiraSetting{
			JiraSettingID:         data.JiraSettingID,
			Name:                  data.Name,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			IdentityField:         data.IdentityField,
			IdentityValue:         data.IdentityValue,
			JiraID:                data.JiraID,
			JiraKey:               data.JiraKey,
			Status:                diagnosis.Status_IN_PROGRESS.String(),
			StatusDetail:          fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
			ScanAt:                data.ScanAt,
		}); err != nil {
			return nil, err
		}
	case "diagnosis:wpscan":
		data, err := d.repository.GetWpscanSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		msg, err := makeWpscanMessage(req.ProjectId, req.SettingId, data.TargetURL)
		if err != nil {
			return nil, err
		}
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.sendWpscanMessage(ctx, msg)
		if err != nil {
			return nil, err
		}
		if _, err = d.repository.UpsertWpscanSetting(ctx, &model.WpscanSetting{
			WpscanSettingID:       data.WpscanSettingID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			TargetURL:             data.TargetURL,
			Status:                diagnosis.Status_IN_PROGRESS.String(),
			StatusDetail:          fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
			ScanAt:                data.ScanAt,
		}); err != nil {
			return nil, err
		}
	case "diagnosis:portscan":
		data, err := d.repository.GetPortscanSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			logger.Error("Error occured when getting PortscanSetting", zap.Error(err))
			return nil, err
		}
		portscanTargets, err := d.repository.ListPortscanTarget(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			logger.Error("Error occured when getting PortscanTargets", zap.Error(err))
			return nil, err
		}
		for _, target := range *portscanTargets {
			msg, err := makePortscanMessage(data.ProjectID, data.PortscanSettingID, target.PortscanTargetID, target.Target)
			if err != nil {
				logger.Error("Error occured when making Portscan message", zap.Error(err))
				continue
			}
			msg.ScanOnly = req.ScanOnly
			resp, err = d.sqs.sendPortscanMessage(ctx, msg)
			if err != nil {
				logger.Error("Error occured when sending Portscan message", zap.Error(err))
				continue
			}
			var scanAt time.Time
			if !zero.IsZeroVal(target.ScanAt) {
				scanAt = target.ScanAt
			}
			if _, err = d.repository.UpsertPortscanTarget(ctx, &model.PortscanTarget{
				PortscanTargetID:  target.PortscanTargetID,
				PortscanSettingID: target.PortscanSettingID,
				ProjectID:         target.ProjectID,
				Target:            target.Target,
				Status:            diagnosis.Status_IN_PROGRESS.String(),
				StatusDetail:      fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
				ScanAt:            scanAt,
			}); err != nil {
				logger.Error("Error occured when upsert Portscan target", zap.Error(err))
				return nil, err
			}
		}
		if _, err = d.repository.UpsertPortscanSetting(ctx, &model.PortscanSetting{
			PortscanSettingID:     data.PortscanSettingID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			Name:                  data.Name,
		}); err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	logger.Info("Invoke scanned.", zap.String("MessageId", *resp.MessageId))
	return &diagnosis.InvokeScanResponse{Message: "Start Diagnosis."}, nil
}

func (s *diagnosisService) InvokeScanAll(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {

	listjiraSetting, err := s.repository.ListAllJiraSetting(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
			ScanOnly:              true,
		}); err != nil {
			// errorが出ても続行
			logger.Error("InvokeScanAll error", zap.Error(err))
		}
	}

	listWpscanSetting, err := s.repository.ListAllWpscanSetting(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
			ScanOnly:              true,
		}); err != nil {
			// errorが出ても続行
			logger.Error("InvokeScanAll error", zap.Error(err))
		}
	}

	return &empty.Empty{}, nil
}
