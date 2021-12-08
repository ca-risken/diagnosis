package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosis.DiagnosisServiceClient
}

func newHandler() *sqsHandler {
	h := &sqsHandler{}
	h.findingClient = newFindingClient()
	appLogger.Info("Start Finding Client")
	h.alertClient = newAlertClient()
	appLogger.Info("Start Alert Client")
	h.diagnosisClient = newDiagnosisClient()
	appLogger.Info("Start Diagnosis Client")
	return h
}

func (s *sqsHandler) HandleMessage(ctx context.Context, sqsMsg *sqs.Message) error {
	msgBody := aws.StringValue(sqsMsg.Body)
	appLogger.Infof("got message. message: %v", msgBody)
	// Parse message
	msg, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", msgBody, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	requestID, err := logging.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, RequestID=%s", requestID)

	// basic setting must be used anytime.
	setting, err := s.GetBasicScanSetting(ctx, msg.ProjectID, msg.ApplicationScanID)
	if err != nil {
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, err.Error())
		return mimosasqs.WrapNonRetryable(err)
	}

	// check scanner can access target for confirming to scan correctly
	if err = checkAccessibleTarget(setting.Target); err != nil {
		appLogger.Warnf("Failed to access target, target: %v error: %v", setting.Target, err)
		errOutput := fmt.Errorf("Failed to access target. Target seems to be down or unaccessible from scanner.")
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, errOutput.Error())
		return mimosasqs.WrapNonRetryable(errOutput)
	}

	// Run ApplicationScan
	apiKey, err := generateAPIKey()
	if err != nil {
		appLogger.Errorf("Failed to create apiKey, error: %v", err)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed exec application scan Ask the system administrator. ")
		return mimosasqs.WrapNonRetryable(err)
	}
	cli, err := newApplicationScanClient(apiKey)
	if err != nil {
		appLogger.Errorf("Failed to create ApplicationScanClient, error: %v", err)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed exec application scan Ask the system administrator. ")
		return mimosasqs.WrapNonRetryable(err)
	}
	if err != nil {
		appLogger.Errorf("Failed to generate API Key, error: %v", err)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed exec application scan Ask the system administrator. ")
		return mimosasqs.WrapNonRetryable(err)
	}
	_, segment := xray.BeginSubsegment(ctx, "runApplicationScan")
	pID, err := cli.executeZap(apiKey)
	if err != nil {
		appLogger.Errorf("Failed to execute ZAP, error: %v", err)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed exec application scan Ask the system administrator. ")
		return mimosasqs.WrapNonRetryable(err)
	}
	var scanResult *zapResult
	switch strings.ToUpper(msg.ApplicationScanType) {
	case "BASIC":
		scanResult, err = cli.handleBasicScan(setting, msg.ApplicationScanID, msg.ProjectID, msg.Name)
	default:
		err = errors.New("ScanType is not configured.")
	}
	errTerminate := cli.terminateZap(pID)
	if errTerminate != nil {
		appLogger.Errorf("Failed to terminate Zap, error: %v", errTerminate)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed exec application scan Ask the system administrator. ")
		return mimosasqs.WrapNonRetryable(err)
	}
	segment.Close(err)
	if err != nil {
		appLogger.Errorf("Failed to exec basicScan, error: %v", err)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed exec application scan Ask the system administrator. ")
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{setting.Target},
	}); err != nil {
		appLogger.Errorf("Failed to clear finding score. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		_ = s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, "Failed to clear finding score. ApplicationScanID")
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put Finding and Tag Finding
	if err := s.putFindings(ctx, scanResult, msg, setting.Target); err != nil {
		appLogger.Errorf("Faild to put findings. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put ApplicationScan
	if err := s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, true, ""); err != nil {
		appLogger.Errorf("Faild to put applicationscan. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	appLogger.Infof("end Scan, RequestID=%s", requestID)
	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Errorf("Faild to analyze alert. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil

}

func parseMessage(msg string) (*message.ApplicationScanQueueMessage, error) {
	message := &message.ApplicationScanQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	//	if err := message.Validate(); err != nil {
	//		return nil, err
	//	}
	return message, nil
}

func checkAccessibleTarget(target string) error {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *sqsHandler) putApplicationScan(ctx context.Context, applicationScanID, projectID uint32, isSuccess bool, errDetail string) error {

	resp, err := s.diagnosisClient.GetApplicationScan(ctx, &diagnosis.GetApplicationScanRequest{ApplicationScanId: applicationScanID, ProjectId: projectID})
	if err != nil {
		return err
	}
	applicationScan := &diagnosis.ApplicationScanForUpsert{
		ApplicationScanId:     resp.ApplicationScan.ApplicationScanId,
		DiagnosisDataSourceId: resp.ApplicationScan.DiagnosisDataSourceId,
		ProjectId:             resp.ApplicationScan.ProjectId,
		Name:                  resp.ApplicationScan.Name,
		ScanType:              resp.ApplicationScan.ScanType,
		Status:                getStatus(isSuccess),
		ScanAt:                time.Now().Unix(),
	}
	if isSuccess {
		applicationScan.StatusDetail = ""
	} else {
		applicationScan.StatusDetail = string(errDetail)
	}
	_, err = s.diagnosisClient.PutApplicationScan(ctx, &diagnosis.PutApplicationScanRequest{ProjectId: resp.ApplicationScan.ProjectId, ApplicationScan: applicationScan})
	if err != nil {
		return err
	}

	return nil
}

func (s *sqsHandler) GetBasicScanSetting(ctx context.Context, projectID, applicationScanID uint32) (*diagnosis.ApplicationScanBasicSetting, error) {
	resp, err := s.diagnosisClient.GetApplicationScanBasicSetting(ctx, &diagnosis.GetApplicationScanBasicSettingRequest{ApplicationScanId: applicationScanID, ProjectId: projectID})
	if err != nil {
		return nil, err
	}
	return resp.ApplicationScanBasicSetting, nil
}

func (s *sqsHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	appLogger.Info("Success to analyze alert.")
	return nil
}

func getStatus(isSuccess bool) diagnosis.Status {
	if isSuccess {
		return diagnosis.Status_OK
	}
	return diagnosis.Status_ERROR
}

func generateAPIKey() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil

}
