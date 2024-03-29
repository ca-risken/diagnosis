package portscan

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/common/pkg/portscan"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/common"
)

const (
	MaxScore = 10.0
)

func (s *SqsHandler) putNmapFinding(ctx context.Context, nmapResult *portscan.NmapResult, projectID uint32, dataSource, data string, target string) error {
	putFinding := &finding.FindingForUpsert{
		Description:      nmapResult.GetDescription(),
		DataSource:       dataSource,
		DataSourceId:     nmapResult.GetDataSourceID(""),
		ResourceName:     nmapResult.ResourceName,
		ProjectId:        projectID,
		OriginalScore:    nmapResult.GetScore(),
		OriginalMaxScore: MaxScore,
		Data:             data,
	}
	resFinding, err := s.putFinding(ctx, putFinding, target)
	if err != nil {
		return err
	}
	recommend := getRecommend(recommendTypeNmap)
	putRecommend := &finding.PutRecommendRequest{
		ProjectId:      projectID,
		FindingId:      resFinding.FindingId,
		DataSource:     dataSource,
		Type:           recommendTypeNmap,
		Risk:           recommend.Risk,
		Recommendation: recommend.Recommendation,
	}
	err = s.putRecommend(ctx, putRecommend)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqsHandler) putAdditionalFinding(ctx context.Context, nmapResult *portscan.NmapResult, projectID uint32, dataSource, data, target string) error {
	for key, detail := range nmapResult.ScanDetail {
		additionalCheckResult, ok := portscan.GetAdditionalCheckResult(key)
		if !ok || detail == false {
			continue
		}
		addFinding := &finding.FindingForUpsert{
			Description:      additionalCheckResult.GetDescription(nmapResult.Target, nmapResult.Port),
			DataSource:       dataSource,
			DataSourceId:     nmapResult.GetDataSourceID(additionalCheckResult.GetRecommendType()),
			ResourceName:     nmapResult.ResourceName,
			ProjectId:        projectID,
			OriginalScore:    additionalCheckResult.GetScore(),
			OriginalMaxScore: MaxScore,
			Data:             data,
		}
		resFinding, err := s.putFinding(ctx, addFinding, target)
		if err != nil {
			return err
		}
		putRecommend := &finding.PutRecommendRequest{
			ProjectId:      projectID,
			FindingId:      resFinding.FindingId,
			DataSource:     dataSource,
			Type:           additionalCheckResult.GetRecommendType(),
			Risk:           additionalCheckResult.GetRisk(),
			Recommendation: additionalCheckResult.GetRecommendation(),
		}
		err = s.putRecommend(ctx, putRecommend)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqsHandler) putFindings(ctx context.Context, results []*portscan.NmapResult, message *message.PortscanQueueMessage) error {
	for _, r := range results {
		externalLink := makeURL(r.Target, r.Port)
		data, err := json.Marshal(map[string]interface{}{"data": *r, "external_link": externalLink})
		if err != nil {
			return err
		}
		err = s.putNmapFinding(ctx, r, message.ProjectID, message.DataSource, string(data), message.Target)
		if err != nil {
			return err
		}
		err = s.putAdditionalFinding(ctx, r, message.ProjectID, message.DataSource, string(data), message.Target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqsHandler) putFinding(ctx context.Context, f *finding.FindingForUpsert, target string) (*finding.Finding, error) {
	res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
	if err != nil {
		return nil, err
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDiagnosis); err != nil {
		return nil, fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagDiagnosis, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagFQDN); err != nil {
		return nil, fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagFQDN, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagPortscan); err != nil {
		return nil, fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagPortscan, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, target); err != nil {
		return nil, fmt.Errorf("tag finding error: tag=%s, err=%w", target, err)
	}
	return res.Finding, nil
}

func (s *SqsHandler) putRecommend(ctx context.Context, recommend *finding.PutRecommendRequest) error {
	if _, err := s.findingClient.PutRecommend(ctx, recommend); err != nil {
		return err
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

func makeURL(target string, port int) string {
	switch port {
	case 443:
		return fmt.Sprintf("https://%v", target)
	case 80:
		return fmt.Sprintf("http://%v", target)
	default:
		return fmt.Sprintf("http://%v:%v", target, port)
	}
}
