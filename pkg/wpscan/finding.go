package wpscan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/vikyd/zero"
)

func (s *SqsHandler) putFindings(ctx context.Context, wpscanResult *wpscanResult, message *message.WpscanQueueMessage) error {
	for _, interstingFinding := range wpscanResult.InterestingFindings {
		findingIntersting, err := interstingFinding.getFinding(message)
		if err != nil {
			return err
		}
		recommendInteresting, err := interstingFinding.getRecommend(message)
		if err != nil {
			return err
		}
		err = s.putFinding(ctx, findingIntersting, recommendInteresting, message.TargetURL)
		if err != nil {
			return err
		}
	}
	if wpscanResult.Version != nil {
		findingVersion, err := wpscanResult.Version.getFinding(ctx, message)
		if err != nil {
			return err
		}
		recommendVersion, err := wpscanResult.Version.getRecommend(ctx, message)
		if err != nil {
			return err
		}
		err = s.putFinding(ctx, findingVersion, recommendVersion, message.TargetURL)
		if err != nil {
			return err
		}
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
	if wpscanResult.CheckAccess != nil {
		wpscanResult.CheckAccess.isUserFound = isUserFound
		findingAccess, err := wpscanResult.CheckAccess.getFinding(ctx, message)
		if err != nil {
			return err
		}
		recommendAccess, err := wpscanResult.CheckAccess.getRecommend(ctx, message)
		if err != nil {
			return err
		}
		err = s.putFinding(ctx, findingAccess, recommendAccess, message.TargetURL)
		if err != nil {
			return err
		}
	}

	for _, p := range wpscanResult.Plugins {
		findingPlugin, recommendPlugin, err := getPluginFinding(ctx, p, message)
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
		DataSource:     message.DataSourceNameWPScan,
		Type:           recommendType,
		Risk:           risk,
		Recommendation: recommendation,
	}
}

func (i *interestingFindings) getFinding(message *message.WpscanQueueMessage) (*finding.FindingForUpsert, error) {
	data, err := json.Marshal(map[string]interestingFindings{"data": *i})
	if err != nil {
		return nil, err
	}
	findingInf, ok := wpscanFindingMap[i.Type]
	var desc string
	var score float32
	if !ok {
		desc = i.ToS
		score = 1.0
		f := makeFinding(desc, fmt.Sprintf("interesting_findings_%v", i.ToS), score, &data, message)
		return f, nil
	}
	desc = findingInf.Description
	score = findingInf.Score
	if zero.IsZeroVal(desc) {
		desc = i.ToS
	}
	f := makeFinding(desc, fmt.Sprintf("interesting_findings_%v", i.ToS), score, &data, message)
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return f, nil
	}

	return f, nil
}

func (i *interestingFindings) getRecommend(message *message.WpscanQueueMessage) (*finding.PutRecommendRequest, error) {
	findingInf, ok := wpscanFindingMap[i.Type]
	if !ok {
		return nil, nil
	}
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)

	return r, nil
}

func (v *version) getFinding(ctx context.Context, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, error) {
	if zero.IsZeroVal(v.Number) {
		return nil, nil
	}
	findingType := typeVersion
	if v.Status == "insecure" {
		findingType = typeVersionInsecure
	}
	findingInf, ok := wpscanFindingMap[findingType]
	if !ok {
		return nil, fmt.Errorf("unknown finding type: findingType=%s", findingType)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	f := makeFinding(fmt.Sprintf(findingInf.Description, v.Number), fmt.Sprintf("version_%v", message.TargetURL), findingInf.Score, &data, message)

	return f, nil
}

func (v *version) getRecommend(ctx context.Context, message *message.WpscanQueueMessage) (*finding.PutRecommendRequest, error) {
	if zero.IsZeroVal(v.Number) {
		return nil, nil
	}
	findingType := typeVersion
	if v.Status == "insecure" {
		findingType = typeVersionInsecure
	}
	findingInf, ok := wpscanFindingMap[findingType]
	if !ok {
		return nil, fmt.Errorf("unknown finding type: findingType=%s", findingType)
	}
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)

	return r, nil
}

func (c *checkAccess) getFinding(ctx context.Context, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var findingInf wpscanFindingInformation
	var ok bool

	if c.isFoundAccesibleURL {
		if c.isUserFound {
			findingInf, ok = wpscanFindingMap[typeLoginOpenedUserFound]
		} else {
			findingInf, ok = wpscanFindingMap[typeLoginOpened]
		}
	} else {
		findingInf, ok = wpscanFindingMap[typeLoginClosed]
	}

	if !ok {
		return nil, fmt.Errorf("failed to get access information. isFoundAccesibleURL=%v, isUserFound=%v", c.isFoundAccesibleURL, c.isUserFound)
	}
	f := makeFinding(findingInf.Description, fmt.Sprintf("Accesible_%v", message.TargetURL), findingInf.Score, &data, message)
	return f, nil
}

func (c *checkAccess) getRecommend(ctx context.Context, message *message.WpscanQueueMessage) (*finding.PutRecommendRequest, error) {
	var findingInf wpscanFindingInformation
	var ok bool

	if c.isFoundAccesibleURL {
		if c.isUserFound {
			findingInf, ok = wpscanFindingMap[typeLoginOpenedUserFound]
		} else {
			findingInf, ok = wpscanFindingMap[typeLoginOpened]
		}
	} else {
		findingInf, ok = wpscanFindingMap[typeLoginClosed]
	}

	if !ok {
		return nil, fmt.Errorf("failed to get access information. isFoundAccesibleURL=%v, isUserFound=%v", c.isFoundAccesibleURL, c.isUserFound)
	}
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)

	return r, nil
}

func getPluginFinding(ctx context.Context, plugin plugin, message *message.WpscanQueueMessage) (*finding.FindingForUpsert, *finding.PutRecommendRequest, error) {
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
		return nil, nil, fmt.Errorf("failed to get plugin information. plugin_name=%v", plugin.Slug)
	}
	f := makeFinding(fmt.Sprintf(findingInf.Description, plugin.Slug), fmt.Sprintf("plugin_%v", plugin.Slug), findingInf.Score, &data, message)
	if zero.IsZeroVal(findingInf.Risk) || zero.IsZeroVal(findingInf.Recommendation) {
		return f, nil, nil
	}
	r := makeRecommend(message.ProjectID, 0, findingInf.RecommendType, findingInf.Risk, findingInf.Recommendation)
	return f, r, nil
}

func (s *SqsHandler) putFinding(ctx context.Context, f *finding.FindingForUpsert, r *finding.PutRecommendRequest, target string) error {
	if f == nil {
		return nil
	}
	res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
	if err != nil {
		return err
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDiagnosis); err != nil {
		return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagDiagnosis, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagURL); err != nil {
		return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagURL, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagWordPress); err != nil {
		return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagWordPress, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagVulnerability); err != nil {
		return fmt.Errorf("tag finding error: tag=%s, err=%w", common.TagVulnerability, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, target); err != nil {
		return fmt.Errorf("tag finding error: tag=%s, err=%w", target, err)
	}

	if r == nil {
		return nil
	}
	r.FindingId = res.Finding.FindingId
	return s.putRecommend(ctx, r)
}

func (s *SqsHandler) putRecommend(ctx context.Context, r *finding.PutRecommendRequest) error {
	if _, err := s.findingClient.PutRecommend(ctx, r); err != nil {
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

func generateDataSourceID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
