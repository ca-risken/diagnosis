package applicationscan

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type SqsHandler struct {
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosis.DiagnosisServiceClient
	appClient       applicationScanAPI
	logger          logging.Logger
}

func NewSqsHandler(
	fc finding.FindingServiceClient,
	ac alert.AlertServiceClient,
	dc diagnosis.DiagnosisServiceClient,
	appc applicationScanAPI,
	l logging.Logger,
) *SqsHandler {
	return &SqsHandler{
		findingClient:   fc,
		alertClient:     ac,
		diagnosisClient: dc,
		appClient:       appc,
		logger:          l,
	}
}

func (s *SqsHandler) HandleMessage(ctx context.Context, sqsMsg *types.Message) error {
	msgBody := aws.ToString(sqsMsg.Body)
	s.logger.Infof(ctx, "got message. message: %v", msgBody)
	// Parse message
	msg, err := message.ParseApplicationScanMessage(msgBody)
	if err != nil {
		s.logger.Errorf(ctx, "Invalid message. message: %v, error: %v", msgBody, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	beforeScanAt := time.Now()
	requestID, err := s.logger.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		s.logger.Warnf(ctx, "Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	s.logger.Infof(ctx, "start Scan, RequestID=%s", requestID)

	// basic setting must be used anytime.
	setting, err := s.GetBasicScanSetting(ctx, msg.ProjectID, msg.ApplicationScanID)
	if err != nil {
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// check scanner can access target for confirming to scan correctly
	if err = checkAccessibleTarget(setting.Target); err != nil {
		s.logger.Warnf(ctx, "Failed to access target, target: %v error: %v", setting.Target, err)
		err = fmt.Errorf("failed to access target. Target seems to be down or unaccessible from scanner. err=%w", err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Run ApplicationScan
	apiKey, err := generateAPIKey()
	if err != nil {
		s.logger.Errorf(ctx, "Failed to create apiKey, error: %v", err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	s.appClient.setApiKey(apiKey)

	tspan, _ := tracer.StartSpanFromContext(ctx, "runApplicationScan")
	scanResult, err := s.runApplicationScan(ctx, msg, setting, apiKey)
	tspan.Finish(tracer.WithError(err))
	if err != nil {
		s.logger.Errorf(ctx, "Failed to run application scan, error: %v", err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put Finding and Tag Finding
	if err := s.putFindings(ctx, scanResult, msg, setting.Target); err != nil {
		s.logger.Errorf(ctx, "Failed to put findings. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear score for inactive findings
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{fmt.Sprintf("application_scan_id:%v", msg.ApplicationScanID)},
		BeforeAt:   beforeScanAt.Unix(),
	}); err != nil {
		s.logger.Errorf(ctx, "Failed to clear finding score. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put ApplicationScan
	if err := s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, true, ""); err != nil {
		s.logger.Errorf(ctx, "Failed to put applicationscan. ApplicationScanID: %v, error: %v", msg.ApplicationScanID, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	s.logger.Infof(ctx, "end Scan, RequestID=%s", requestID)
	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		s.logger.Notifyf(ctx, logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil
}

func (s *SqsHandler) runApplicationScan(ctx context.Context, msg *message.ApplicationScanQueueMessage, setting *diagnosis.ApplicationScanBasicSetting, apiKey string) (*zapResult, error) {
	if strings.ToUpper(msg.ApplicationScanType) != "BASIC" {
		return nil, errors.New("ScanType is not configured")
	}

	pID, err := s.appClient.executeZap(ctx, apiKey)
	if err != nil {
		s.logger.Errorf(ctx, "failed to execute ZAP, error: %v", err)
		return nil, err
	}
	defer func(int) {
		err = s.appClient.terminateZap(pID)
		if err != nil {
			s.logger.Warnf(ctx, "failed to terminate Zap, error: %v", err)
		}
	}(pID)

	var scanResult *zapResult
	scanResult, err = s.appClient.handleBasicScan(setting, msg.ApplicationScanID, msg.ProjectID, msg.Name)
	if err != nil {
		s.logger.Errorf(ctx, "failed to exec basicScan, error: %v", err)
		return nil, err
	}
	return scanResult, nil
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

func (s *SqsHandler) updateStatusToError(ctx context.Context, msg *message.ApplicationScanQueueMessage, err error) {
	if updateErr := s.putApplicationScan(ctx, msg.ApplicationScanID, msg.ProjectID, false, err.Error()); updateErr != nil {
		s.logger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
	}
}

func (s *SqsHandler) putApplicationScan(ctx context.Context, applicationScanID, projectID uint32, isSuccess bool, errDetail string) error {

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

func (s *SqsHandler) GetBasicScanSetting(ctx context.Context, projectID, applicationScanID uint32) (*diagnosis.ApplicationScanBasicSetting, error) {
	resp, err := s.diagnosisClient.GetApplicationScanBasicSetting(ctx, &diagnosis.GetApplicationScanBasicSettingRequest{ApplicationScanId: applicationScanID, ProjectId: projectID})
	if err != nil {
		return nil, err
	}
	return resp.ApplicationScanBasicSetting, nil
}

func (s *SqsHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	s.logger.Info(ctx, "Success to analyze alert.")
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
