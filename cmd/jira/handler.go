package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/common"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/vikyd/zero"
	"go.uber.org/zap"
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
	logger.Info("Start Jira Client")
	h.findingClient = newFindingClient()
	logger.Info("Start Finding Client")
	h.alertClient = newAlertClient()
	logger.Info("Start Alert Client")
	h.diagnosisClient = newDiagnosisClient()
	logger.Info("Start Diagnosis Client")
	return h
}

func (s *sqsHandler) HandleMessage(msg *sqs.Message) error {
	msgBody := aws.StringValue(msg.Body)
	logger.Info("got message", zap.String("message", msgBody))
	// Parse message
	message, err := parseMessage(msgBody)
	if err != nil {
		logger.Error("Invalid message", zap.String("sqs_message", "message"), zap.Error(err))
		return err
	}

	// Get jira Project
	project, errMap := s.jira.getJiraProject(message.JiraKey, message.JiraID, message.IdentityField, message.IdentityValue)
	if zero.IsZeroVal(project) {
		logger.Warn("Faild to get jira project", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.String("Project", project))
		if err := s.putJiraSetting(message.JiraSettingID, message.ProjectID, false, errMap); err != nil {
			logger.Error("Faild to put jira_setting", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.Error(err))
			return nil
		}
		return nil
	}

	if zero.IsZeroVal(project) {
		if err := s.putJiraSetting(message.JiraSettingID, message.ProjectID, false, errMap); err != nil {
			logger.Error("Faild to put jira_setting", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.Error(err))
			return err
		}
	}

	// Get jira
	findings, err := s.getJira(project, message)
	if err != nil {
		logger.Error("Faild to get findngs to Diagnosis Jira", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.Uint32("ProjectID", message.ProjectID), zap.Error(err))
		return nil
	}

	// Put finding to core
	ctx := context.Background()
	if err := s.putFindings(ctx, findings); err != nil {
		logger.Error("Faild to put findngs", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.Uint32("ProjectID", message.ProjectID), zap.Error(err))
		return err
	}

	// Put JiraSetting
	if err := s.putJiraSetting(message.JiraSettingID, message.ProjectID, true, nil); err != nil {
		logger.Error("Faild to put jira_setting", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.Error(err))
		return err
	}

	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, message.ProjectID); err != nil {
		logger.Error("Faild to analyze alert.", zap.Uint32("JiraSettingID", message.JiraSettingID), zap.Error(err))
		return err
	}

	return nil

}

func (s *sqsHandler) getJira(project string, message *message.JiraQueueMessage) ([]*finding.FindingForUpsert, error) {
	putData := []*finding.FindingForUpsert{}
	issueList, err := s.jira.listIssues(project)
	if err != nil {
		logger.Error("Jira.listIssues", zap.Error(err))
		return nil, err
	}
	issues := issueList.Issues
	for _, issue := range issues {
		buf, err := json.Marshal(issue)
		if err != nil {
			logger.Error("Failed to json encoding error", zap.Error(err))
			return nil, err
		}
		if !isOpen(issue.Fields.Status.Name) {
			continue
		}
		score := getScore(issue.Fields.Priority.Name)
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
		logger.Info("Success to PutFinding", zap.Any("Finding", f))
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
		logger.Error("Failed to TagFinding.", zap.Error(err))
		return err
	}
	return nil
}

func (s *sqsHandler) putJiraSetting(jiraSettingID, projectID uint32, isSuccess bool, errMap map[string]string) error {
	ctx := context.Background()
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
	logger.Info("Success to analyze alert.")
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
