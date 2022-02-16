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
	"github.com/ca-risken/diagnosis/proto/diagnosis"
)

type sqsHandler struct {
	wpscanConfig    WpscanConfig
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosis.DiagnosisServiceClient
}

func (s *sqsHandler) HandleMessage(ctx context.Context, sqsMsg *sqs.Message) error {
	msgBody := aws.StringValue(sqsMsg.Body)
	appLogger.Infof("got message. message: %v", msgBody)
	// Parse message
	msg, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", msgBody, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}
	var options wpscanOptions
	if err := json.Unmarshal([]byte(msg.Options), &options); err != nil {
		appLogger.Errorf("Failed to Unmarshal options. message: %v, error: %v", msgBody, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}
	requestID, err := appLogger.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, RequestID=%s", requestID)

	// Run WPScan
	_, segment := xray.BeginSubsegment(ctx, "runWPScan")
	wpscanResult, err := s.wpscanConfig.run(msg.TargetURL, msg.WpscanSettingID, options)
	segment.Close(err)
	if err != nil {
		appLogger.Errorf("Failed exec WPScan, error: %v", err)
		_ = s.putWpscanSetting(ctx, msg.WpscanSettingID, msg.ProjectID, false, "Failed exec WPScan Ask the system administrator. ")
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}
	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{msg.TargetURL},
	}); err != nil {
		appLogger.Errorf("Failed to clear finding score. WpscanSettingID: %v, error: %v", msg.WpscanSettingID, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}
	err = s.putFindings(ctx, wpscanResult, msg)
	if err != nil {
		appLogger.Errorf("Failed put Findings, error: %v", err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	// Put WpscanSetting
	if err := s.putWpscanSetting(ctx, msg.WpscanSettingID, msg.ProjectID, true, ""); err != nil {
		appLogger.Errorf("Faild to put rel_osint_data_source. WpscanSettingID: %v, error: %v", msg.WpscanSettingID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	appLogger.Infof("end Scan, RequestID=%s", requestID)
	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Notifyf(logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil

}

func parseMessage(msg string) (*message.WpscanQueueMessage, error) {
	message := &message.WpscanQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *sqsHandler) handleErrorWithUpdateStatus(ctx context.Context, msg *message.WpscanQueueMessage, err error) error {
	if updateErr := s.putWpscanSetting(ctx, msg.WpscanSettingID, msg.ProjectID, false, err.Error()); updateErr != nil {
		appLogger.Warnf("Failed to update scan status error: err=%+v", updateErr)
	}
	return mimosasqs.WrapNonRetryable(err)
}

func (s *sqsHandler) putWpscanSetting(ctx context.Context, wpscanSettingID, projectID uint32, isSuccess bool, errDetail string) error {
	resp, err := s.diagnosisClient.GetWpscanSetting(ctx, &diagnosis.GetWpscanSettingRequest{WpscanSettingId: wpscanSettingID, ProjectId: projectID})
	if err != nil {
		return err
	}
	wpscanSetting := &diagnosis.WpscanSettingForUpsert{
		WpscanSettingId:       resp.WpscanSetting.WpscanSettingId,
		DiagnosisDataSourceId: resp.WpscanSetting.DiagnosisDataSourceId,
		ProjectId:             resp.WpscanSetting.ProjectId,
		TargetUrl:             resp.WpscanSetting.TargetUrl,
		Options:               resp.WpscanSetting.Options,
		ScanAt:                time.Now().Unix(),
	}
	wpscanSetting.Status = getStatus(isSuccess)
	if isSuccess {
		wpscanSetting.StatusDetail = ""
	} else {
		wpscanSetting.StatusDetail = string(errDetail)
	}
	_, err = s.diagnosisClient.PutWpscanSetting(ctx, &diagnosis.PutWpscanSettingRequest{ProjectId: resp.WpscanSetting.ProjectId, WpscanSetting: wpscanSetting})
	if err != nil {
		return err
	}

	return nil
}

func (s *sqsHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	appLogger.Info("Success to analyze alert.")
	return nil
}

const (
	// PriorityScore
	MaxScore = 10.0
)

func getStatus(isSuccess bool) diagnosis.Status {
	if isSuccess {
		return diagnosis.Status_OK
	}
	return diagnosis.Status_ERROR
}
