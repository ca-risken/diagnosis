package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	diagnosisClient "github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsHandler struct {
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosisClient.DiagnosisServiceClient
}

func newHandler() *sqsHandler {
	return &sqsHandler{
		findingClient:   newFindingClient(),
		alertClient:     newAlertClient(),
		diagnosisClient: newDiagnosisClient(),
	}
}

func (s *sqsHandler) HandleMessage(msg *sqs.Message) error {
	msgBody := aws.StringValue(msg.Body)
	appLogger.Infof("got message: %s", msgBody)
	// Parse message
	message, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message: SQS_msg=%+v, err=%+v", msg, err)
		return err
	}
	ctx := context.Background()
	// Get portscan
	portscan, err := newPortscanClient()
	if err != nil {
		appLogger.Errorf("Failed to create Portscan session: err=%+v", err)
		return s.putPortscanSetting(message.PortscanSettingID, message.ProjectID, false, err.Error())
	}
	statusDetail := ""

	portscan.target = makeTargets(message.Target)

	findings, err := portscan.getResult(message)
	if err != nil {
		appLogger.Warnf("Failed to get findings to Diagnosis Portscan: PortscanSettingID=%+v, err=%+v", message.PortscanSettingID, err)
	}
	// Put finding to core
	if err := s.putFindings(ctx, findings); err != nil {
		appLogger.Errorf("Failed to put findings: PortscanSettingID=%+v, err=%+v", message.PortscanSettingID, err)
		statusDetail = fmt.Sprintf("%v%v", statusDetail, err.Error())
		return s.putPortscanSetting(message.PortscanSettingID, message.ProjectID, false, statusDetail)
	}

	if err := s.putPortscanSetting(message.PortscanSettingID, message.ProjectID, true, ""); err != nil {
		return err
	}
	if err := s.analyzeAlert(ctx, message.ProjectID); err != nil {
		appLogger.Errorf("Failed to analyze alert: PortscanSettingID=%+v, err=%+v", message.PortscanSettingID, err)
		return err
	}

	appLogger.Infof("Scan finished. ProjectID: %v, PortscanSettingID: %v, Target: %v", message.ProjectID, message.PortscanSettingID, message.Target)

	return nil
}

func (s *sqsHandler) putPortscanSetting(portscanSettingID, projectID uint32, isSuccess bool, errDetail string) error {
	ctx := context.Background()
	resp, err := s.diagnosisClient.GetPortscanSetting(ctx, &diagnosisClient.GetPortscanSettingRequest{PortscanSettingId: portscanSettingID, ProjectId: projectID})
	if err != nil {
		return err
	}

	portscanSetting := &diagnosisClient.PortscanSettingForUpsert{
		PortscanSettingId:     resp.PortscanSetting.PortscanSettingId,
		DiagnosisDataSourceId: resp.PortscanSetting.DiagnosisDataSourceId,
		ProjectId:             resp.PortscanSetting.ProjectId,
		Name:                  resp.PortscanSetting.Name,
		ScanAt:                time.Now().Unix(),
	}

	if isSuccess {
		portscanSetting.Status = diagnosisClient.Status_OK
		portscanSetting.StatusDetail = ""
	} else {
		portscanSetting.Status = diagnosisClient.Status_ERROR
		portscanSetting.StatusDetail = string(errDetail)
	}
	_, err = s.diagnosisClient.PutPortscanSetting(ctx, &diagnosisClient.PutPortscanSettingRequest{ProjectId: resp.PortscanSetting.ProjectId, PortscanSetting: portscanSetting})
	if err != nil {
		return err
	}

	return nil
}

func parseMessage(msg string) (*message.PortscanQueueMessage, error) {
	message := &message.PortscanQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	//	if err := message.Validate(); err != nil {
	//		return nil, err
	//	}
	return message, nil
}

func (s *sqsHandler) analyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{
		ProjectId: projectID,
	})
	return err
}
