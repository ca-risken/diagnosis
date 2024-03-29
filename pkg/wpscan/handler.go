package wpscan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	wpscanConfig    *WpscanConfig
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosis.DiagnosisServiceClient
	logger          logging.Logger
}

func NewSqsHandler(
	wc *WpscanConfig,
	fc finding.FindingServiceClient,
	ac alert.AlertServiceClient,
	dc diagnosis.DiagnosisServiceClient,
	l logging.Logger,
) *SqsHandler {
	return &SqsHandler{
		wpscanConfig:    wc,
		findingClient:   fc,
		alertClient:     ac,
		diagnosisClient: dc,
		logger:          l,
	}
}

func (s *SqsHandler) HandleMessage(ctx context.Context, sqsMsg *types.Message) error {
	msgBody := aws.ToString(sqsMsg.Body)
	s.logger.Infof(ctx, "got message. message: %v", msgBody)
	// Parse message
	msg, err := message.ParseWpscanMessage(msgBody)
	if err != nil {
		s.logger.Errorf(ctx, "Invalid message. message: %v, error: %v", msgBody, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	var options wpscanOptions
	if err := json.Unmarshal([]byte(msg.Options), &options); err != nil {
		s.logger.Errorf(ctx, "Failed to Unmarshal options. message: %v, error: %v", msgBody, err)
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

	// Run WPScan
	tspan, _ := tracer.StartSpanFromContext(ctx, "runWPScan")
	wpscanResult, err := s.wpscanConfig.run(ctx, msg.TargetURL, msg.WpscanSettingID, options)
	tspan.Finish(tracer.WithError(err))
	if err != nil {
		s.logger.Errorf(ctx, "Failed exec WPScan, error: %v", err)
		// Customize error message when failed WPScan
		s.updateStatusToError(ctx, msg, errors.New("failed exec WPScan. See the document. https://docs.security-hub.jp/contact/faq/#wpscan"))
		return mimosasqs.WrapNonRetryable(err)
	}
	err = s.putFindings(ctx, wpscanResult, msg)
	if err != nil {
		s.logger.Errorf(ctx, "Failed put Findings, error: %v", err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear score for inactive findings
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{msg.TargetURL},
		BeforeAt:   beforeScanAt.Unix(),
	}); err != nil {
		s.logger.Errorf(ctx, "Failed to clear finding score. WpscanSettingID: %v, error: %v", msg.WpscanSettingID, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put WpscanSetting
	if err := s.putWpscanSetting(ctx, msg.WpscanSettingID, msg.ProjectID, true, ""); err != nil {
		s.logger.Errorf(ctx, "Faild to put rel_osint_data_source. WpscanSettingID: %v, error: %v", msg.WpscanSettingID, err)
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

func (s *SqsHandler) updateStatusToError(ctx context.Context, msg *message.WpscanQueueMessage, err error) {
	if updateErr := s.putWpscanSetting(ctx, msg.WpscanSettingID, msg.ProjectID, false, err.Error()); updateErr != nil {
		s.logger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
	}
}

func (s *SqsHandler) putWpscanSetting(ctx context.Context, wpscanSettingID, projectID uint32, isSuccess bool, errDetail string) error {
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

func (s *SqsHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	s.logger.Info(ctx, "Success to analyze alert.")
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
