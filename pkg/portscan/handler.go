package portscan

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	diagnosisClient "github.com/ca-risken/datasource-api/proto/diagnosis"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type SqsHandler struct {
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosisClient.DiagnosisServiceClient
	logger          logging.Logger
}

func NewSqsHandler(
	fc finding.FindingServiceClient,
	ac alert.AlertServiceClient,
	dc diagnosisClient.DiagnosisServiceClient,
	l logging.Logger,
) *SqsHandler {
	return &SqsHandler{
		findingClient:   fc,
		alertClient:     ac,
		diagnosisClient: dc,
		logger:          l,
	}
}

func (s *SqsHandler) HandleMessage(ctx context.Context, sqsMsg *types.Message) error {
	msgBody := aws.ToString(sqsMsg.Body)
	s.logger.Infof(ctx, "got message: %s", msgBody)
	// Parse message
	msg, err := message.ParsePortscanMessage(msgBody)
	if err != nil {
		s.logger.Errorf(ctx, "Invalid message: SQS_msg=%+v, err=%+v", msgBody, err)
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

	// Get portscan
	portscan, err := newPortscanClient()
	if err != nil {
		s.logger.Errorf(ctx, "Failed to create Portscan session: err=%+v", err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	portscan.makeTargets(msg.Target)

	tspan, tctx := tracer.StartSpanFromContext(ctx, "getResult")
	nmapResults, err := portscan.getResult(tctx, msg)
	tspan.Finish(tracer.WithError(err))
	if err != nil {
		s.logger.Warnf(ctx, "Failed to get findings to Diagnosis Portscan: PortscanSettingID=%+v, Target=%+v, err=%+v", msg.PortscanSettingID, msg.Target, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put finding to core
	if err := s.putFindings(ctx, nmapResults, msg); err != nil {
		s.logger.Errorf(ctx, "Failed to put findings: PortscanSettingID=%+v, Target=%+v, err=%+v", msg.PortscanSettingID, msg.Target, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear score for inactive findings
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{msg.Target},
		BeforeAt:   beforeScanAt.Unix(),
	}); err != nil {
		s.logger.Errorf(ctx, "Failed to clear finding score. PortscanSettingID: %v, error: %v", msg.PortscanSettingID, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	if err := s.putPortscanTarget(ctx, msg.PortscanTargetID, msg.ProjectID, true, ""); err != nil {
		s.logger.Errorf(ctx, "Failed to put portscanTarget: PortscanSettingID=%+v, Target=%+v, err=%+v", msg.PortscanSettingID, msg.Target, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	s.logger.Infof(ctx, "Scan finished. ProjectID: %v, PortscanSettingID: %v, Target: %v, RequestID: %s", msg.ProjectID, msg.PortscanSettingID, msg.Target, requestID)

	if msg.ScanOnly {
		return nil
	}
	if err := s.analyzeAlert(ctx, msg.ProjectID); err != nil {
		s.logger.Notifyf(ctx, logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil
}

func (s *SqsHandler) updateStatusToError(ctx context.Context, msg *message.PortscanQueueMessage, err error) {
	if updateErr := s.putPortscanTarget(ctx, msg.PortscanTargetID, msg.ProjectID, false, err.Error()); updateErr != nil {
		s.logger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
	}
}

func (s *SqsHandler) putPortscanTarget(ctx context.Context, portscanTargetID, projectID uint32, isSuccess bool, errDetail string) error {
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

func (s *SqsHandler) analyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{
		ProjectId: projectID,
	})
	return err
}
