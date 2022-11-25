package applicationscan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/common"
)

func (s *SqsHandler) putFindings(ctx context.Context, zapResult *zapResult, msg *message.ApplicationScanQueueMessage, target string) error {
	for _, site := range zapResult.Site {
		for _, alert := range site.Alerts {
			data, err := json.Marshal(map[string]zapResultAlert{"data": alert})
			if err != nil {
				return fmt.Errorf("failed to marshal zapResult for makeFinding: err=%w", err)
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
				return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagDiagnosis, err)
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagURL); err != nil {
				return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagURL, err)
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagApplicationScan); err != nil {
				return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagApplicationScan, err)
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagVulnerability); err != nil {
				return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagVulnerability, err)
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, fmt.Sprintf("application_scan_id:%v", msg.ApplicationScanID)); err != nil {
				return fmt.Errorf("tag finding error: tag=%s, err=%w", fmt.Sprintf("application_scan_id:%v", msg.ApplicationScanID), err)
			}
			if err = s.putRecommend(ctx, res.Finding.ProjectId, res.Finding.FindingId, msg.DataSource, &alert); err != nil {
				return fmt.Errorf("put recommend error: alert=%v, err=%w", alert, err)
			}
		}
	}
	return nil
}

func (s *SqsHandler) tagFinding(ctx context.Context, projectID uint32, findingID uint64, tag string) error {
	_, err := s.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
		ProjectId: projectID,
		Tag: &finding.FindingTagForUpsert{
			FindingId: findingID,
			ProjectId: projectID,
			Tag:       tag,
		}})
	if err != nil {
		return err
	}
	return nil
}

func (s SqsHandler) putRecommend(ctx context.Context, projectID uint32, findingID uint64, dataSource string, alert *zapResultAlert) error {
	r := s.getRecommend(ctx, alert)
	if r.Risk == "" && r.Recommendation == "" {
		s.logger.Warnf(ctx, "Failed to get Recommendation, zapResultAlert=%+v", alert)
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
