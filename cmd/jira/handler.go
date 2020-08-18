package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

type sqsHandler struct {
	jira          jiraAPI
	findingClient finding.FindingServiceClient
}

func newHandler() *sqsHandler {
	h := &sqsHandler{}
	h.jira = newJiraClient()
	logger.Info("Start Jira Client")
	h.findingClient = newFindingClient()
	logger.Info("Start Finding Client")
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

	// Get jira
	findings, err := s.getJira(message)

	if err != nil {
		logger.Error("Faild to get findngs to Diagnosis Jira", zap.String("RecordID", message.RecordID), zap.String("JiraID", message.JiraID), zap.Error(err))
		return err
	}

	//Put finding to core
	ctx := context.Background()
	if err := s.putFindings(ctx, findings); err != nil {
		logger.Error("Faild to put findngs", zap.String("RecordID", message.RecordID), zap.String("JiraID", message.JiraID), zap.Error(err))
		return err
	}
	return nil
}

func (s *sqsHandler) getJira(message *message.DiagnosisQueueMessage) ([]*finding.FindingForUpsert, error) {
	putData := []*finding.FindingForUpsert{}
	//	projects, err := s.jira.listProjects()
	//	for _, project := range *projects {
	//	}
	//	if err != nil {
	//		logger.Error("Jira.listProjects", zap.Error(err))
	//		return nil, err
	//	}
	issueList, err := s.jira.listIssues(message.JiraKey, message.JiraID, message.RecordID)
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
		score := getScore(issue.Fields.Priority.Name)
		resource := getResourceName(issue.Fields.Target)
		putData = append(putData, &finding.FindingForUpsert{
			Description:      issue.Fields.Summary,
			DataSource:       message.DataSource,
			DataSourceId:     issue.Key,
			ResourceName:     resource,
			ProjectId:        message.ProjectID,
			OriginalScore:    score,
			OriginalMaxScore: 10.0,
			Data:             string(buf),
		})
	}
	return putData, nil
}

func parseMessage(msg string) (*message.DiagnosisQueueMessage, error) {
	message := &message.DiagnosisQueueMessage{}
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
		_, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
		if err != nil {
			return err
		}
		logger.Info("Success to PutFinding", zap.Any("Finding", f))
	}
	return nil
}

const (
	// PriorityScore
	MaxScore             = 10.0
	ScoreHigh            = 10.0
	ScoreMiddle          = 5.0
	ScoreLow             = 3.0
	ScoreInformation     = 1.0
	ScoreOther           = 0.0
	TypeScoreHigh        = "HIGH"
	TypeScoreMiddle      = "MIDDLE"
	TypeScoreLow         = "LOW"
	TypeScoreInformation = "INFORMATION"
)

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

func getResourceName(target string) string {
	urlLines := strings.Split(target, "\r\n")
	var resources []string
	for _, s := range urlLines {
		if strings.Index(s, "http") > -1 {
			u, err := url.Parse(s)
			if err != nil {
				continue
			}
			retResource := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
			resources = append(resources, retResource)
		}
	}
	m := make(map[string]struct{})
	newList := make([]string, 0)
	for _, resource := range resources {
		if _, ok := m[resource]; !ok {
			m[resource] = struct{}{}
			newList = append(newList, resource)
		}
	}
	ret := strings.Join(newList, ",")
	ret = "target::" + ret
	if len(ret) > 255 {
		ret = ret[0:252] + "..."
	}
	return ret
}
