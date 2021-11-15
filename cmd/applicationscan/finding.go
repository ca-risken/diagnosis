package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
)

func makeFindings(zapResult *zapResult, message *message.ApplicationScanQueueMessage, target string) ([]*finding.FindingForUpsert, error) {
	var findings []*finding.FindingForUpsert
	for _, site := range zapResult.Site {
		for _, alert := range site.Alerts {
			data, err := json.Marshal(map[string]zapResultAlert{"data": alert})
			if err != nil {
				appLogger.Errorf("Failed to marshal zapResult for makeFinding. err: %v", err)
			}
			findings = append(findings, &finding.FindingForUpsert{
				Description:      getDescription(&alert, target),
				DataSource:       message.DataSource,
				DataSourceId:     generateDataSourceID(fmt.Sprintf("%v_%v", target, alert.Alert)),
				ResourceName:     target,
				ProjectId:        message.ProjectID,
				OriginalScore:    getScore(&alert),
				OriginalMaxScore: MaxScore,
				Data:             string(data),
			})
		}
	}

	return findings, nil
}

func (s *sqsHandler) putFindings(ctx context.Context, findings []*finding.FindingForUpsert, target string) error {
	for _, f := range findings {
		res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
		if err != nil {
			return err
		}
		if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDiagnosis); err != nil {
			appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagDiagnosis, err)
		}
		if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagApplicationScan); err != nil {
			appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagApplicationScan, err)
		}
		if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagVulnerability); err != nil {
			appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagVulnerability, err)
		}
		if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, target); err != nil {
			appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", target, err)
		}
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
		appLogger.Errorf("Failed to TagFinding. error: %v", err)
		return err
	}
	return nil
}

func generateDataSourceID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func getDescription(alert *zapResultAlert, target string) string {
	return fmt.Sprintf("%v found in %v.", alert.Alert, target)
}

func getScore(alert *zapResultAlert) float32 {
	if strings.HasPrefix(strings.ToLower(alert.RiskDesc), ScorePrefixInformation) {
		return ScoreInformation
	}
	if strings.HasPrefix(strings.ToLower(alert.RiskDesc), ScorePrefixLow) {
		return ScoreLow
	}
	if strings.HasPrefix(strings.ToLower(alert.RiskDesc), ScorePrefixMedium) {
		return ScoreMedium
	}
	if strings.HasPrefix(strings.ToLower(alert.RiskDesc), ScorePrefixHigh) {
		return ScoreHigh
	}
	return ScoreOther
}

const (
	// PriorityScore
	MaxScore               = 10.0
	ScoreInformation       = 1.0
	ScoreLow               = 3.0
	ScoreMedium            = 6.0
	ScoreHigh              = 8.0
	ScoreOther             = 1.0
	ScorePrefixInformation = "information"
	ScorePrefixLow         = "low"
	ScorePrefixMedium      = "medium"
	ScorePrefixHigh        = "high"
)
