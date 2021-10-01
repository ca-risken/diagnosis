package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/vikyd/zero"
)

type sqsHandler struct {
	jira            jiraAPI
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	diagnosisClient diagnosis.DiagnosisServiceClient
}

func newHandler() *sqsHandler {
	h := &sqsHandler{}
	h.jira = newJiraClient()
	appLogger.Info("Start Jira Client")
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
	appLogger.Infof("got message: %s", msgBody)
	// Parse message
	msg, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", msgBody, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	requestID, err := logging.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: error: %v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, request_id: %v, jiraSettingID: %v", requestID, msg.JiraSettingID)

	// Get jira Project
	project, errMap := s.jira.getJiraProject(ctx, msg.JiraKey, msg.JiraID, msg.IdentityField, msg.IdentityValue)
	if zero.IsZeroVal(project) {
		appLogger.Warnf("Faild to get jira project, jiraSettingID: %v", msg.JiraSettingID)
		if err := s.putJiraSetting(ctx, msg.JiraSettingID, msg.ProjectID, false, errMap); err != nil {
			appLogger.Errorf("Faild to put jira_setting,jiraSettingID: %v, error: %v", msg.JiraSettingID, err)
			return mimosasqs.WrapNonRetryable(err)
		}
		return mimosasqs.WrapNonRetryable(err)
	}

	if zero.IsZeroVal(project) {
		if err := s.putJiraSetting(ctx, msg.JiraSettingID, msg.ProjectID, false, errMap); err != nil {
			appLogger.Errorf("Faild to put jira_setting, error: %v", err)
			return err
		}
	}

	// Get jira
	findings, err := s.getJira(ctx, project, msg)
	if err != nil {
		appLogger.Errorf("Faild to get findngs to Diagnosis Jira, error: %v", err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put finding to core
	if err := s.putFindings(ctx, findings); err != nil {
		appLogger.Errorf("Faild to put findngs, error: %v", err)
		return err
	}

	// Put JiraSetting
	if err := s.putJiraSetting(ctx, msg.JiraSettingID, msg.ProjectID, true, nil); err != nil {
		appLogger.Errorf("Faild to put jira_setting, error: %v", err)
		return err
	}
	appLogger.Infof("end Scan, request_id: %v, JiraSettingID: %v", requestID, msg.JiraSettingID)
	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Errorf("Faild to analyze alert, error: %v", err)
		return err
	}
	return nil

}

func (s *sqsHandler) getJira(ctx context.Context, project string, message *message.JiraQueueMessage) ([]*finding.FindingForUpsert, error) {
	putData := []*finding.FindingForUpsert{}
	issueList, err := s.jira.listIssues(ctx, project)
	if err != nil {
		appLogger.Errorf("Failed to list Issues, error: %v", err)
		return nil, err
	}
	issues := issueList.Issues
	for _, issue := range issues {
		buf, err := json.Marshal(issue)
		if err != nil {
			appLogger.Errorf("Failed to json encoding, error: %v", err)
			return nil, err
		}
		var score float32
		if isOpen(issue.Fields.Status.Name) {
			score = getScore(issue.Fields.Priority.Name)
		} else {
			score = 1.0
		}
		putData = append(putData, &finding.FindingForUpsert{
			Description:      issue.Fields.Summary,
			DataSource:       message.DataSource,
			DataSourceId:     issue.Key,
			ResourceName:     issue.Fields.Project.Name,
			ProjectId:        message.ProjectID,
			OriginalScore:    score,
			OriginalMaxScore: 10.0,
			Data:             string(buf),
		})
	}
	return putData, nil
}

func parseMessage(msg string) (*message.JiraQueueMessage, error) {
	message := &message.JiraQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	//	if err := message.Validate(); err != nil {
	//		return nil, err
	//	}
	return message, nil
}

func (s *sqsHandler) putFindings(ctx context.Context, findings []*finding.FindingForUpsert) error {
	for _, f := range findings {
		res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
		if err != nil {
			return err
		}
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDiagnosis)
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagJira)
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagVulnerability)
	}
	return nil
}

func (s *sqsHandler) tagFinding(ctx context.Context, projectID uint32, findingID uint64, tag string) error {

	_, err := s.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
		ProjectId: projectID,
		Tag: &finding.FindingTagForUpsert{
			FindingId: findingID,
			ProjectId: projectID,
			Tag:       tag,
		}})
	if err != nil {
		appLogger.Errorf("Failed to TagFinding, error: %v", err)
		return err
	}
	return nil
}

func (s *sqsHandler) putJiraSetting(ctx context.Context, jiraSettingID, projectID uint32, isSuccess bool, errMap map[string]string) error {
	resp, err := s.diagnosisClient.GetJiraSetting(ctx, &diagnosis.GetJiraSettingRequest{JiraSettingId: jiraSettingID, ProjectId: projectID})
	if err != nil {
		return err
	}
	jiraSetting := &diagnosis.JiraSettingForUpsert{
		JiraSettingId:         resp.JiraSetting.JiraSettingId,
		Name:                  resp.JiraSetting.Name,
		DiagnosisDataSourceId: resp.JiraSetting.DiagnosisDataSourceId,
		ProjectId:             resp.JiraSetting.ProjectId,
		IdentityField:         resp.JiraSetting.IdentityField,
		IdentityValue:         resp.JiraSetting.IdentityValue,
		JiraId:                resp.JiraSetting.JiraId,
		JiraKey:               resp.JiraSetting.JiraKey,
		ScanAt:                time.Now().Unix(),
	}
	jiraSetting.Status = getStatus(isSuccess)
	if isSuccess {
		jiraSetting.StatusDetail = ""
	} else {
		errDetail, err := json.Marshal(errMap)
		if err != nil {
			return err
		}
		jiraSetting.StatusDetail = string(errDetail)
	}
	_, err = s.diagnosisClient.PutJiraSetting(ctx, &diagnosis.PutJiraSettingRequest{ProjectId: resp.JiraSetting.ProjectId, JiraSetting: jiraSetting})
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
	MaxScore             = 10.0
	ScoreHigh            = 10.0
	ScoreMiddle          = 6.0
	ScoreLow             = 3.0
	ScoreInformation     = 1.0
	ScoreOther           = 0.1
	TypeScoreHigh        = "HIGH"
	TypeScoreMiddle      = "MIDDLE"
	TypeScoreLow         = "LOW"
	TypeScoreInformation = "INFORMATION"
	StatusClosed         = "クローズ"
)

func isOpen(status string) bool {
	if strings.Index(status, StatusClosed) > -1 {
		return false
	}
	return true
}

func getStatus(isSuccess bool) diagnosis.Status {
	if isSuccess {
		return diagnosis.Status_OK
	}
	return diagnosis.Status_ERROR
}

func getScore(name string) float32 {
	switch strings.ToUpper(name) {
	case TypeScoreHigh:
		return ScoreHigh
	case TypeScoreMiddle:
		return ScoreMiddle
	case TypeScoreLow:
		return ScoreLow
	case TypeScoreInformation:
		return ScoreInformation
	default:
		return ScoreOther
	}
}
