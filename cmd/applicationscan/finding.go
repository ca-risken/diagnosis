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

func (s *sqsHandler) putFindings(ctx context.Context, zapResult *zapResult, msg *message.ApplicationScanQueueMessage, target string) error {
	for _, site := range zapResult.Site {
		for _, alert := range site.Alerts {
			data, err := json.Marshal(map[string]zapResultAlert{"data": alert})
			if err != nil {
				appLogger.Errorf("Failed to marshal zapResult for makeFinding. err: %v", err)
			}
			res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{
				Description:      getDescription(&alert, target),
				DataSource:       msg.DataSource,
				DataSourceId:     generateDataSourceID(fmt.Sprintf("%v_%v", target, alert.Alert)),
				ResourceName:     target,
				ProjectId:        msg.ProjectID,
				OriginalScore:    getScore(&alert),
				OriginalMaxScore: MaxScore,
				Data:             string(data),
			}})
			if err != nil {
				return err
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDiagnosis); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagDiagnosis, err)
				return err
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagURL); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagURL, err)
				return err
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagApplicationScan); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagApplicationScan, err)
				return err
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagVulnerability); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagVulnerability, err)
				return err
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, target); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", target, err)
				return err
			}
			if err = s.putRecommend(ctx, res.Finding.ProjectId, res.Finding.FindingId, msg.DataSource, &alert); err != nil {
				appLogger.Errorf("Failed to put Recommend. alert: %v, error: %v", alert, err)
				return err
			}
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

func (s *sqsHandler) putRecommend(ctx context.Context, projectID uint32, findingID uint64, dataSource string, alert *zapResultAlert) error {
	r := getRecommend(alert)
	if r.Risk == "" && r.Recommendation == "" {
		appLogger.Warnf("Failed to get Recommendation, zapResultAlert=%+v", alert)
		return nil
	}
	_, err := s.findingClient.PutRecommend(ctx, &finding.PutRecommendRequest{
		ProjectId:      projectID,
		FindingId:      findingID,
		DataSource:     dataSource,
		Type:           alert.Alert,
		Risk:           r.Risk,
		Recommendation: r.Recommendation,
	})
	if err != nil {
		appLogger.Errorf("Failed to PutRecommend. projectID: %d, findingID: %d, error: %v", projectID, findingID, err)
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
