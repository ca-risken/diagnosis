package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/vikyd/zero"
)

func (s *sqsHandler) putFindings(ctx context.Context, wpscanResult *wpscanResult, message *message.WpscanQueueMessage) error {
	for _, interstingFinding := range wpscanResult.InterestingFindings {
		findingIntersting, recommendInteresting, err := getInterestingFinding(interstingFinding, message)
		if err != nil {
			return err
		}
		err = s.putFinding(ctx, findingIntersting, recommendInteresting, message.TargetURL)
		if err != nil {
			return err
		}
	}
	findingVersion, recommendVersion, err := getVersionFinding(wpscanResult.Version, message)
	if err != nil {
		return err
	}
	err = s.putFinding(ctx, findingVersion, recommendVersion, message.TargetURL)
	if err != nil {
		return err
	}
	isUserFound := false
	for key, val := range wpscanResult.Users {
		isUserFound = true
		data, err := json.Marshal(map[string]interface{}{"data": val})
		if err != nil {
			return err
		}
		desc := fmt.Sprintf("User %v was found.", key)
		score := float32(3.0)
		findingUser := makeFinding(desc, fmt.Sprintf("username_%v", key), score, &data, message)
		err = s.putFinding(ctx, findingUser, nil, message.TargetURL)
		if err != nil {
			return err
		}
	}
	for _, access := range wpscanResult.AccessList {
		findingAccess, recommendAccess, err := getAccessFinding(access, isUserFound, message)
		if err != nil {
			return err
		}
		err = s.putFinding(ctx, findingAccess, recommendAccess, message.TargetURL)
		if err != nil {
			return err
		}

	}

	for _, p := range wpscanResult.Plugins {
		findingPlugin, recommendPlugin, err := getPluginFinding(p, message)
		if err != nil {
			return err
		}
		err = s.putFinding(ctx, findingPlugin, recommendPlugin, message.TargetURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeFinding(description, dataSourceID string, score float32, data *[]byte, message *message.WpscanQueueMessage) *finding.FindingForUpsert {
	return &finding.FindingForUpsert{
		Description:      description,
		DataSource:       message.DataSource,
		DataSourceId:     generateDataSourceID(dataSourceID),
		ResourceName:     message.TargetURL,
		ProjectId:        message.ProjectID,
		OriginalScore:    score,
		OriginalMaxScore: MaxScore,
		Data:             string(*data),
	}
}

func makeRecommend(projectID uint32, findingID uint64, recommendType, risk, recommendation string) *finding.PutRecommendRequest {
	return &finding.PutRecommendRequest{
		ProjectId:      projectID,
		FindingId:      findingID,
		DataSource:     common.DataSourceNameWPScan,
		Type:           recommendType,
		Risk:           risk,
		Recommendation: recommendation,
	}
}
func getInterestingFinding(ie interestingFindings, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, *finding.PutRecommendRequest, error) {
	data, err := json.Marshal(map[string]interestingFindings{"data": ie})
	if err != nil {
		return nil, nil, err
	}
	findingInf, ok := wpscanFindingMap[ie.Type]
	var desc string
	var score float32
	if !ok {
		desc = ie.ToS
		score = 1.0
		f := makeFinding(desc, fmt.Sprintf("interesting_findings_%v", ie.ToS), score, &data, message)
		return f, nil, nil
	}
	desc = findingInf.Description
	score = findingInf.Score
	if zero.IsZeroVal(desc) {
		desc = ie.ToS
	}
	f := makeFinding(desc, fmt.Sprintf("interesting_findings_%v", ie.ToS), score, &data, message)
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return f, nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)

	return f, r, nil
}

func getVersionFinding(wpScanVersion version, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, *finding.PutRecommendRequest, error) {
	if zero.IsZeroVal(wpScanVersion.Number) {
		return nil, nil, nil
	}
	findingType := typeVersion
	if wpScanVersion.Status == "insecure" {
		findingType = typeVersionInsecure
	}
	findingInf, ok := wpscanFindingMap[findingType]
	if !ok {
		appLogger.Warnf("Failed to get finding information, Unknown findingType=%v", findingType)
		return nil, nil, nil
	}
	data, err := json.Marshal(map[string]version{"data": wpScanVersion})
	if err != nil {
		return nil, nil, err
	}
	f := makeFinding(fmt.Sprintf(findingInf.Description, wpScanVersion.Number), fmt.Sprintf("version_%v", message.TargetURL), findingInf.Score, &data, message)
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return f, nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)

	return f, r, nil
}

func getAccessFinding(access checkAccess, isUserFound bool, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, *finding.PutRecommendRequest, error) {
	data, err := json.Marshal(map[string]interface{}{"data": map[string]string{
		"url": access.Target,
	}})
	if err != nil {
		return nil, nil, err
	}
	var findingInf wpscanFindingInformation
	var ok bool
	switch access.Type {
	case "Login":
		if !access.IsAccess {
			findingInf, ok = wpscanFindingMap[typeLoginClosed]
		} else if isUserFound {
			findingInf, ok = wpscanFindingMap[typeLoginOpenedUserFound]
		} else {
			findingInf, ok = wpscanFindingMap[typeLoginOpened]
		}
	default:
		return nil, nil, nil
	}
	if !ok {
		appLogger.Warnf("Failed to get access information, Unknown access.Type=%v", access.Type)
		return nil, nil, nil
	}
	f := makeFinding(findingInf.Description, fmt.Sprintf("Accesible_%v", access.Target), findingInf.Score, &data, message)
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return f, nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)

	return f, r, nil
}

func getPluginFinding(plugin plugin, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, *finding.PutRecommendRequest, error) {
	data, err := json.Marshal(plugin)
	if err != nil {
		return nil, nil, err
	}
	var findingInf wpscanFindingInformation
	var ok bool
	if len(plugin.Vulnerabilities) == 0 {
		findingInf, ok = wpscanFindingMap[typePlugin]
	} else if zero.IsZeroVal(plugin.Version.Number) {
		findingInf, ok = wpscanFindingMap[typePluginUnknownVersion]
	} else {
		findingInf, ok = wpscanFindingMap[typePluginVulnerable]
	}
	if !ok {
		appLogger.Errorf("Failed to get plugin information, plugin=%v", plugin)
		return nil, nil, fmt.Errorf("failed to get plugin information. plugin_name=%v", plugin.Slug)
	}
	f := makeFinding(fmt.Sprintf(findingInf.Description, plugin.Slug), fmt.Sprintf("plugin_%v", plugin.Slug), findingInf.Score, &data, message)
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return f, nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)
	return f, r, nil
}

func (s *sqsHandler) putFinding(ctx context.Context, f *finding.FindingForUpsert, r *finding.PutRecommendRequest, target string) error {
	if f == nil {
		return nil
	}
	res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
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
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagWordPress); err != nil {
		appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagWordPress, err)
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

	if r == nil {
		return nil
	}
	r.FindingId = res.Finding.FindingId
	return s.putRecommend(ctx, r)
}

func (s *sqsHandler) putRecommend(ctx context.Context, r *finding.PutRecommendRequest) error {
	if _, err := s.findingClient.PutRecommend(ctx, r); err != nil {
		return err
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
