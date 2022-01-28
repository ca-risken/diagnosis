package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/diagnosis/pkg/model"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
)

const (
	InvokeScanAllDataSource       = "all"
	InvokeScanJira                = "jira"
	InvokeScanWPScan              = "wpscan"
	InvokeScanPortScan            = "portscan"
	InvokeScanApplicationScan     = "applicationscan"
	DataSourceNameJira            = "diagnosis:jira"
	DataSourceNameWPScan          = "diagnosis:wpscan"
	DataSourceNamePortScan        = "diagnosis:portscan"
	DataSourceNameApplicationScan = "diagnosis:application-scan"
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
	case DataSourceNameJira:
		data, err := d.repository.GetJiraSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		msg, err := makeJiraMessage(req.ProjectId, req.SettingId, data)
		if err != nil {
			appLogger.Errorf("Error occured when making Jira message, error: %v", err)
			return nil, err
		}
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.sendJiraMessage(ctx, msg)
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
	case DataSourceNameWPScan:
		data, err := d.repository.GetWpscanSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		options := data.Options
		if zero.IsZeroVal(options) {
			options = "{}"
		}
		msg, err := makeWpscanMessage(req.ProjectId, req.SettingId, data.TargetURL, options)
		if err != nil {
			appLogger.Errorf("Error occured when making WPScan message, error: %v", err)
			return nil, err
		}
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.sendWpscanMessage(ctx, msg)
		if err != nil {
			appLogger.Errorf("Error occured when sending WPScan message, error: %v", err)
			return nil, err
		}
		var scanAt time.Time
		if !zero.IsZeroVal(data.ScanAt) {
			scanAt = data.ScanAt
		}
		if _, err = d.repository.UpsertWpscanSetting(ctx, &model.WpscanSetting{
			WpscanSettingID:       data.WpscanSettingID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			TargetURL:             data.TargetURL,
			Options:               options,
			Status:                diagnosis.Status_IN_PROGRESS.String(),
			StatusDetail:          fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
			ScanAt:                scanAt,
		}); err != nil {
			appLogger.Errorf("Error occured when upsert WPScanSetting, error: %v", err)
			return nil, err
		}
	case DataSourceNamePortScan:
		data, err := d.repository.GetPortscanSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			appLogger.Errorf("Error occured when getting PortscanSetting, error: %v", err)
			return nil, err
		}
		portscanTargets, err := d.repository.ListPortscanTarget(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			appLogger.Errorf("Error occured when getting PortscanTargets, error: %v", err)
			return nil, err
		}
		for _, target := range *portscanTargets {
			msg, err := makePortscanMessage(data.ProjectID, data.PortscanSettingID, target.PortscanTargetID, target.Target)
			if err != nil {
				appLogger.Errorf("Error occured when making Portscan message, error: %v", err)
				continue
			}
			msg.ScanOnly = req.ScanOnly
			resp, err = d.sqs.sendPortscanMessage(ctx, msg)
			if err != nil {
				appLogger.Errorf("Error occured when sending Portscan message, error: %v", err)
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
				appLogger.Errorf("Error occured when upsert Portscan target, error: %v", err)
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
	case DataSourceNameApplicationScan:
		data, err := d.repository.GetApplicationScan(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			appLogger.Errorf("Error occured when getting PortscanSetting, error: %v", err)
			return nil, err
		}
		msg, err := makeApplicationScanMessage(req.ProjectId, req.SettingId, data.Name, data.ScanType)
		if err != nil {
			return nil, err
		}
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.sendApplicationScanMessage(ctx, msg)
		if err != nil {
			return nil, err
		}
		var scanAt time.Time
		if !zero.IsZeroVal(data.ScanAt) {
			scanAt = data.ScanAt
		}
		if _, err = d.repository.UpsertApplicationScan(ctx, &model.ApplicationScan{
			ApplicationScanID:     data.ApplicationScanID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			Name:                  data.Name,
			ScanType:              data.ScanType,
			Status:                diagnosis.Status_IN_PROGRESS.String(),
			StatusDetail:          fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
			ScanAt:                scanAt,
		}); err != nil {
			appLogger.Errorf("Error occured when upsert Application scan, error: %v", err)
			return nil, err
		}
	default:
		return nil, nil
	}

	appLogger.Infof("Invoke scanned, MessageID: %v", *resp.MessageId)
	return &diagnosis.InvokeScanResponse{Message: "Start Diagnosis."}, nil
}

func (s *diagnosisService) InvokeScanAll(ctx context.Context, req *diagnosis.InvokeScanAllRequest) (*empty.Empty, error) {

	scanDataSource := InvokeScanAllDataSource
	if !zero.IsZeroVal(req.DiagnosisDataSourceId) {
		dataSource, err := s.repository.GetDiagnosisDataSource(ctx, 0, req.DiagnosisDataSourceId)
		if err != nil {
			return nil, err
		}
		switch dataSource.Name {
		case DataSourceNameJira:
			scanDataSource = InvokeScanJira
		case DataSourceNameWPScan:
			scanDataSource = InvokeScanWPScan
		case DataSourceNamePortScan:
			scanDataSource = InvokeScanPortScan
		case DataSourceNameApplicationScan:
			scanDataSource = InvokeScanApplicationScan
		}
	}
	if isInvokeScan(scanDataSource, InvokeScanJira) {
		listjiraSetting, err := s.repository.ListAllJiraSetting(ctx)
		if err != nil {
			appLogger.Errorf("Failed to List All JiraSetting., error: %v", err)
			return nil, err
		}

		for _, jiraSetting := range *listjiraSetting {
			if s.skipProject(ctx, jiraSetting.ProjectID) {
				continue
			}
			if _, err := s.InvokeScan(ctx, &diagnosis.InvokeScanRequest{
				ProjectId:             jiraSetting.ProjectID,
				SettingId:             jiraSetting.JiraSettingID,
				DiagnosisDataSourceId: jiraSetting.DiagnosisDataSourceID,
				ScanOnly:              true,
			}); err != nil {
				// errorが出ても続行
				appLogger.Errorf("InvokeScanAll error, error: %v", err)
			}
		}
	}
	if isInvokeScan(scanDataSource, InvokeScanWPScan) {
		listWpscanSetting, err := s.repository.ListAllWpscanSetting(ctx)
		if err != nil {
			appLogger.Errorf("Failed to List All WPScanSetting., error: %v", err)
			return nil, err
		}

		for _, WpscanSetting := range *listWpscanSetting {
			if s.skipProject(ctx, WpscanSetting.ProjectID) {
				continue
			}
			if _, err := s.InvokeScan(ctx, &diagnosis.InvokeScanRequest{
				ProjectId:             WpscanSetting.ProjectID,
				SettingId:             WpscanSetting.WpscanSettingID,
				DiagnosisDataSourceId: WpscanSetting.DiagnosisDataSourceID,
				ScanOnly:              true,
			}); err != nil {
				// errorが出ても続行
				appLogger.Errorf("InvokeScanAll error, error: %v", err)
			}
		}
	}

	return &empty.Empty{}, nil
}

func isInvokeScan(scanDataSource, targetDataSource string) bool {
	if scanDataSource == InvokeScanAllDataSource {
		return true
	}
	return scanDataSource == targetDataSource
}

func (s *diagnosisService) skipProject(ctx context.Context, projectID uint32) bool {
	if resp, err := s.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: projectID}); err != nil {
		appLogger.Errorf("Failed to project.IsActive API, err=%+v", err)
		return true
	} else if !resp.Active {
		appLogger.Infof("Skip deactive project, project_id=%d", projectID)
		return true
	}
	return false
}
