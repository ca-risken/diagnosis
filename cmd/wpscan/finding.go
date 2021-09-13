package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/common"
)

func (s *sqsHandler) putFindings(ctx context.Context, findings []*finding.FindingForUpsert) error {
	for _, f := range findings {
		res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
		if err != nil {
			return err
		}
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDiagnosis)
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagWordPress)
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagVulnerability)

		appLogger.Infof("Success to PutFinding. finding: %v", f)
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
