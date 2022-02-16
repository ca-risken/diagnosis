package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/message"
	diagnosisClient "github.com/ca-risken/diagnosis/proto/diagnosis"
)

type sqsHandler struct {
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosisClient.DiagnosisServiceClient
}

func (s *sqsHandler) HandleMessage(ctx context.Context, sqsMsg *sqs.Message) error {
	msgBody := aws.StringValue(sqsMsg.Body)
	appLogger.Infof("got message: %s", msgBody)
	// Parse message
	msg, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message: SQS_msg=%+v, err=%+v", msgBody, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}
	requestID, err := appLogger.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, RequestID=%s", requestID)

	// Get portscan
	portscan, err := newPortscanClient()
	if err != nil {
		appLogger.Errorf("Failed to create Portscan session: err=%+v", err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	portscan.makeTargets(msg.Target)

	xctx, segment := xray.BeginSubsegment(ctx, "getResult")
	nmapResults, err := portscan.getResult(xctx, msg)
	segment.Close(err)
	if err != nil {
		appLogger.Warnf("Failed to get findings to Diagnosis Portscan: PortscanSettingID=%+v, Target=%+v, err=%+v", msg.PortscanSettingID, msg.Target, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{msg.Target},
	}); err != nil {
		appLogger.Errorf("Failed to clear finding score. PortscanSettingID: %v, error: %v", msg.PortscanSettingID, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	// Put finding to core
	if err := s.putFindings(ctx, nmapResults, msg); err != nil {
		appLogger.Errorf("Failed to put findings: PortscanSettingID=%+v, Target=%+v, err=%+v", msg.PortscanSettingID, msg.Target, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	if err := s.putPortscanTarget(ctx, msg.PortscanTargetID, msg.ProjectID, true, ""); err != nil {
		appLogger.Errorf("Failed to put portscanTarget: PortscanSettingID=%+v, Target=%+v, err=%+v", msg.PortscanSettingID, msg.Target, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	appLogger.Infof("Scan finished. ProjectID: %v, PortscanSettingID: %v, Target: %v, RequestID: %s", msg.ProjectID, msg.PortscanSettingID, msg.Target, requestID)

	if msg.ScanOnly {
		return nil
	}
	if err := s.analyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Notifyf(logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil
}

func (s *sqsHandler) handleErrorWithUpdateStatus(ctx context.Context, msg *message.PortscanQueueMessage, err error) error {
	if updateErr := s.putPortscanTarget(ctx, msg.PortscanTargetID, msg.ProjectID, false, err.Error()); updateErr != nil {
		appLogger.Warnf("Failed to update scan status error: err=%+v", updateErr)
	}
	return mimosasqs.WrapNonRetryable(err)
}

func (s *sqsHandler) putPortscanTarget(ctx context.Context, portscanTargetID, projectID uint32, isSuccess bool, errDetail string) error {
	resp, err := s.diagnosisClient.GetPortscanTarget(ctx, &diagnosisClient.GetPortscanTargetRequest{PortscanTargetId: portscanTargetID, ProjectId: projectID})
	if err != nil {
		return err
	}

	portscanTarget := &diagnosisClient.PortscanTargetForUpsert{
		PortscanTargetId:  resp.PortscanTarget.PortscanTargetId,
		PortscanSettingId: resp.PortscanTarget.PortscanSettingId,
		ProjectId:         resp.PortscanTarget.ProjectId,
		Target:            resp.PortscanTarget.Target,
		ScanAt:            time.Now().Unix(),
	}

	if isSuccess {
		portscanTarget.Status = diagnosisClient.Status_OK
		portscanTarget.StatusDetail = ""
	} else {
		portscanTarget.Status = diagnosisClient.Status_ERROR
		portscanTarget.StatusDetail = errDetail
	}
	_, err = s.diagnosisClient.PutPortscanTarget(ctx, &diagnosisClient.PutPortscanTargetRequest{ProjectId: resp.PortscanTarget.ProjectId, PortscanTarget: portscanTarget})
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
