package main

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/golang/protobuf/ptypes/empty"
	mock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestMakeFinding(t *testing.T) {
	cases := []struct {
		name         string
		description  string
		dataSourceID string
		score        float32
		data         *[]byte
		message      *message.WpscanQueueMessage
		want         *finding.FindingForUpsert
	}{
		{
			name:         "OK",
			description:  "description",
			dataSourceID: "dataSourceID",
			score:        1.0,
			data:         &[]byte{},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			want: &finding.FindingForUpsert{
				Description:      "description",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("dataSourceID"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    1.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := makeFinding(c.description, c.dataSourceID, c.score, c.data, c.message)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestMakeRecommend(t *testing.T) {
	cases := []struct {
		name          string
		projectID     uint32
		findingID     uint64
		recommendType string
		risk          string
		recommend     string
		want          *finding.PutRecommendRequest
	}{
		{
			name:          "OK",
			projectID:     1,
			findingID:     1,
			recommendType: "recommendType",
			risk:          "risk",
			recommend:     "recommend",
			want: &finding.PutRecommendRequest{
				ProjectId:      1,
				FindingId:      1,
				DataSource:     "diagnosis:wpscan",
				Type:           "recommendType",
				Risk:           "risk",
				Recommendation: "recommend",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := makeRecommend(c.projectID, c.findingID, c.recommendType, c.risk, c.recommend)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestGetInterestingFinding(t *testing.T) {
	cases := []struct {
		name      string
		ie        interestingFindings
		message   *message.WpscanQueueMessage
		finding   *finding.FindingForUpsert
		recommend *finding.PutRecommendRequest
		wantErr   bool
	}{
		{
			name: "Score 1.0 no recommend",
			ie: interestingFindings{
				URL:               "http://localhost",
				ToS:               "to_s",
				Type:              "type",
				InterstingEntries: []string{"ie"},
				References:        map[string]interface{}{"key": "val"},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "to_s",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("interesting_findings_to_s"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    1.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: nil,
			wantErr:   false,
		},
		{
			name: "Score 6.0 exists recommend",
			ie: interestingFindings{
				URL:               "http://localhost",
				ToS:               "readme",
				Type:              "readme",
				InterstingEntries: []string{"readme"},
				References:        map[string]interface{}{"key": "val"},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "readme",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("interesting_findings_readme"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    6.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       "readme.html",
				Risk: `Readme.html exists
	- Basic information such as version can be identified, which may provide useful information to the attacker.`,
				Recommendation: `Delete readme.html.`},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, _ := json.Marshal(map[string]interestingFindings{"data": c.ie})
			c.finding.Data = string(data)
			f, r, e := getInterestingFinding(c.ie, c.message)
			if !reflect.DeepEqual(c.finding, f) {
				t.Fatalf("Unexpected finding: want=%v, got=%v", c.finding, f)
			}
			if !reflect.DeepEqual(c.recommend, r) {
				t.Fatalf("Unexpected recommend: want=%v, got=%v", c.recommend, r)
			}
			if (c.wantErr && e == nil) || (!c.wantErr && e != nil) {
				t.Fatalf("Unexpected recommend: wantErr=%v, error=%v", c.wantErr, e)
			}
		})
	}
}

func TestGetVersionFinding(t *testing.T) {
	cases := []struct {
		name      string
		ver       version
		message   *message.WpscanQueueMessage
		finding   *finding.FindingForUpsert
		recommend *finding.PutRecommendRequest
		wantErr   bool
	}{
		{
			name: "Score 1.0 no recommend",
			ver: version{
				Number:            "num",
				Status:            "status",
				InterstingEntries: []string{"ie"},
				Vulnerabilities:   []vulnerability{{}},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress version num identified",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("version_http://localhost"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    1.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: nil,
			wantErr:   false,
		},
		{
			name: "Insecure version exists recommend",
			ver: version{
				Number:            "num",
				Status:            "insecure",
				InterstingEntries: []string{"ie"},
				Vulnerabilities:   []vulnerability{{}},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress version num identified (Insecure)",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("version_http://localhost"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    6.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typeVersionInsecure,
				Risk: `WordPress Insecure Version
	- WordPress is not up to date and not secure. 
	- Vulnerabilities may exist, and attacks can cause serious damage.`,
				Recommendation: `Update wordpress.
	- https://wordpress.org/support/article/updating-wordpress/`},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, _ := json.Marshal(map[string]version{"data": c.ver})
			c.finding.Data = string(data)
			f, r, e := getVersionFinding(c.ver, c.message)
			if !reflect.DeepEqual(c.finding, f) {
				t.Fatalf("Unexpected finding:\n want=%v,\n got=%v", c.finding, f)
			}
			if !reflect.DeepEqual(c.recommend, r) {
				t.Fatalf("Unexpected recommend:\n want=%v,\n got=%v", c.recommend, r)
			}
			if (c.wantErr && e == nil) || (!c.wantErr && e != nil) {
				t.Fatalf("Unexpected error: wantErr=%v, error=%v", c.wantErr, e)
			}
		})
	}
}

func TestGetPluginFinding(t *testing.T) {
	cases := []struct {
		name      string
		plugin    plugin
		message   *message.WpscanQueueMessage
		finding   *finding.FindingForUpsert
		recommend *finding.PutRecommendRequest
		wantErr   bool
	}{
		{
			name: "Score 1.0 no recommend",
			plugin: plugin{
				Slug:              "no-vulnerable-plugin",
				LatestVersion:     "1",
				Location:          "http://plugin/location",
				InterstingEntries: []string{"ie"},
				Vulnerabilities:   []vulnerability{},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "Plugin was found. plugin: no-vulnerable-plugin",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("plugin_no-vulnerable-plugin"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    1.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: nil,
			wantErr:   false,
		},
		{
			name: "plugin of unknown version",
			plugin: plugin{
				Slug:              "unknown-version-plugin",
				LatestVersion:     "1",
				Location:          "http://plugin/location",
				InterstingEntries: []string{"ie"},
				Vulnerabilities: []vulnerability{{
					Title:      "vulnerable-plugin",
					FixedIn:    "fixed-in",
					References: map[string]interface{}{"ref": "reference_url"},
					URL:        []string{"urls"},
				}},
				Version: version{Number: ""},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "Plugin of unknown version was found. plugin: unknown-version-plugin",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("plugin_unknown-version-plugin"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    6.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typePluginUnknownVersion,
				Risk: `Plugin of unknown version found.
	Vulnerability exists in some versions.`,
				Recommendation: `Please check the version and make sure it is not affected by the vulnerability.
	If the version is affected, please update the plugin.`},
			wantErr: false,
		},
		{
			name: "Insecure plugin",
			plugin: plugin{
				Slug:              "vulnerable-plugin",
				LatestVersion:     "1",
				Location:          "http://plugin/location",
				InterstingEntries: []string{"ie"},
				Vulnerabilities: []vulnerability{{
					Title:      "vulnerable-plugin",
					FixedIn:    "fixed-in",
					References: map[string]interface{}{"ref": "reference_url"},
					URL:        []string{"urls"},
				}},
				Version: version{Number: "1"},
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "Vulnerable plugin was found. plugin: vulnerable-plugin",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("plugin_vulnerable-plugin"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    8.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typePluginVulnerable,
				Risk: `A vulnerable plugin was found.
	See Finding for details on the vulnerability and its impact.`,
				Recommendation: `Please update your plugins.
	The version in which the vulnerability has been fixed is listed in Fixed_in of Finding.`},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, _ := json.Marshal(c.plugin)
			c.finding.Data = string(data)
			f, r, e := getPluginFinding(c.plugin, c.message)
			if !reflect.DeepEqual(c.finding, f) {
				t.Fatalf("Unexpected finding:\n want=%#v,\n got=%#v", c.finding, f)
			}
			if !reflect.DeepEqual(c.recommend, r) {
				t.Fatalf("Unexpected recommend:\n want=%v,\n got=%v", c.recommend, r)
			}
			if (c.wantErr && e == nil) || (!c.wantErr && e != nil) {
				t.Fatalf("Unexpected error: wantErr=%v, error=%v", c.wantErr, e)
			}
		})
	}
}

func TestGetAccessFinding(t *testing.T) {
	cases := []struct {
		name        string
		access      checkAccess
		isUserFound bool
		message     *message.WpscanQueueMessage
		finding     *finding.FindingForUpsert
		recommend   *finding.PutRecommendRequest
		wantErr     bool
	}{
		{
			name: "Closed no recommend",
			access: checkAccess{
				Target:   "target",
				Goal:     "goal",
				Method:   "GET",
				Type:     "Login",
				IsAccess: false,
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is closed.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    1.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: nil,
			wantErr:   false,
		},
		{
			name: "Open exists recommend",
			access: checkAccess{
				Target:   "target",
				Goal:     "goal",
				Method:   "GET",
				Type:     "Login",
				IsAccess: true,
			},
			message: &message.WpscanQueueMessage{
				DataSource:      "diagnosis:wpscan",
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is opened.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    8.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typeLoginOpened,
				Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.`,
				Recommendation: `Restrict access to the admin panel.`},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, _ := json.Marshal(map[string]interface{}{"data": map[string]string{
				"url": c.access.Target,
			}})
			c.finding.Data = string(data)
			f, r, e := getAccessFinding(c.access, c.isUserFound, c.message)
			if !reflect.DeepEqual(c.finding, f) {
				t.Fatalf("Unexpected finding:\n want=%v,\n got=%v", c.finding, f)
			}
			if !reflect.DeepEqual(c.recommend, r) {
				t.Fatalf("Unexpected recommend:\n want=%v,\n got=%v", c.recommend, r)
			}
			if (c.wantErr && e == nil) || (!c.wantErr && e != nil) {
				t.Fatalf("Unexpected error: wantErr=%v, error=%v", c.wantErr, e)
			}
		})
	}
}

func TestPutFinding(t *testing.T) {
	handler := &sqsHandler{}
	cases := []struct {
		name               string
		target             string
		finding            *finding.FindingForUpsert
		recommend          *finding.PutRecommendRequest
		mockFindingResp    *finding.PutFindingResponse
		mockFindingError   error
		mockRecommendResp  *finding.PutRecommendResponse
		mockRecommendError error
		mockTagResp        *finding.TagFindingResponse
		mockTagError       error
		wantErr            bool
	}{
		{
			name: "no recommend",
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is closed.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    1.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend:         nil,
			target:            "hoge",
			mockFindingResp:   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, Score: 0.5}},
			mockTagResp:       &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 1001}},
			mockRecommendResp: nil,
			wantErr:           false,
		},
		{
			name: "exists recommend",
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is opened.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    8.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typeLoginOpened,
				Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.`,
				Recommendation: `Restrict access to the admin panel.`},
			target:            "hoge",
			mockFindingResp:   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, Score: 0.5}},
			mockTagResp:       &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 1001}},
			mockRecommendResp: &finding.PutRecommendResponse{Recommend: &finding.Recommend{RecommendId: 1001}},
			wantErr:           false,
		},
		{
			name: "putFinding error",
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is opened.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    8.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typeLoginOpened,
				Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.`,
				Recommendation: `Restrict access to the admin panel.`},
			target:             "hoge",
			mockFindingResp:    nil,
			mockFindingError:   errors.New("putFinding error"),
			mockRecommendResp:  nil,
			mockRecommendError: nil,
			wantErr:            true,
		},
		{
			name: "tagFinding error",
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is opened.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    8.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typeLoginOpened,
				Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.`,
				Recommendation: `Restrict access to the admin panel.`},
			target:             "hoge",
			mockFindingResp:    &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, Score: 0.5}},
			mockFindingError:   nil,
			mockTagResp:        nil,
			mockTagError:       errors.New("tagFinding error"),
			mockRecommendResp:  nil,
			mockRecommendError: nil,
			wantErr:            true,
		},
		{
			name: "putRecommend error",
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is opened.",
				DataSource:       "diagnosis:wpscan",
				DataSourceId:     generateDataSourceID("Accesible_target"),
				ResourceName:     "http://localhost",
				ProjectId:        1,
				OriginalScore:    8.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			recommend: &finding.PutRecommendRequest{
				ProjectId:  1,
				FindingId:  0,
				DataSource: "diagnosis:wpscan",
				Type:       typeLoginOpened,
				Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.`,
				Recommendation: `Restrict access to the admin panel.`},
			target:             "hoge",
			mockFindingResp:    &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, Score: 0.5}},
			mockFindingError:   nil,
			mockTagResp:        &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 1001}},
			mockTagError:       nil,
			mockRecommendResp:  nil,
			mockRecommendError: errors.New("putRecommend error"),
			wantErr:            true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			findingMock := &mockFindingClient{}
			handler.findingClient = findingMock
			if c.mockFindingResp != nil || c.mockFindingError != nil {
				findingMock.On("PutFinding").Return(c.mockFindingResp, c.mockFindingError).Once()
			}
			if c.mockTagResp != nil {
				findingMock.On("TagFinding").Return(c.mockTagResp, c.mockTagError).Times(5)
			}
			if c.mockTagError != nil {
				findingMock.On("TagFinding").Return(c.mockTagResp, c.mockTagError).Once()
			}
			if c.mockRecommendResp != nil || c.mockRecommendError != nil {
				findingMock.On("PutRecommend").Return(c.mockRecommendResp, c.mockRecommendError).Once()
			}
			ctx := context.Background()
			e := handler.putFinding(ctx, c.finding, c.recommend, c.target)
			if !findingMock.AssertExpectations(t) {
				t.Fatalf("Unexpected call: , AssertExpectations=%v", findingMock.AssertExpectations(t))
			}
			if (c.wantErr && e == nil) || (!c.wantErr && e != nil) {
				t.Fatalf("Unexpected error: wantErr=%v, error=%v", c.wantErr, e)
			}
		})
	}
}

// Mock Client

type mockFindingClient struct {
	mock.Mock
}

func (m *mockFindingClient) ListFinding(context.Context, *finding.ListFindingRequest, ...grpc.CallOption) (*finding.ListFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingResponse), args.Error(1)
}

func (m *mockFindingClient) BatchListFinding(context.Context, *finding.BatchListFindingRequest, ...grpc.CallOption) (*finding.BatchListFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.BatchListFindingResponse), args.Error(1)
}
func (m *mockFindingClient) GetFinding(context.Context, *finding.GetFindingRequest, ...grpc.CallOption) (*finding.GetFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetFindingResponse), args.Error(1)
}
func (m *mockFindingClient) PutFinding(context.Context, *finding.PutFindingRequest, ...grpc.CallOption) (*finding.PutFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutFindingResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteFinding(context.Context, *finding.DeleteFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListFindingTag(context.Context, *finding.ListFindingTagRequest, ...grpc.CallOption) (*finding.ListFindingTagResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingTagResponse), args.Error(1)
}
func (m *mockFindingClient) ListFindingTagName(context.Context, *finding.ListFindingTagNameRequest, ...grpc.CallOption) (*finding.ListFindingTagNameResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingTagNameResponse), args.Error(1)
}
func (m *mockFindingClient) TagFinding(context.Context, *finding.TagFindingRequest, ...grpc.CallOption) (*finding.TagFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.TagFindingResponse), args.Error(1)
}
func (m *mockFindingClient) UntagFinding(context.Context, *finding.UntagFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListResource(context.Context, *finding.ListResourceRequest, ...grpc.CallOption) (*finding.ListResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceResponse), args.Error(1)
}
func (m *mockFindingClient) GetResource(context.Context, *finding.GetResourceRequest, ...grpc.CallOption) (*finding.GetResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetResourceResponse), args.Error(1)
}
func (m *mockFindingClient) PutResource(context.Context, *finding.PutResourceRequest, ...grpc.CallOption) (*finding.PutResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutResourceResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteResource(context.Context, *finding.DeleteResourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListResourceTag(context.Context, *finding.ListResourceTagRequest, ...grpc.CallOption) (*finding.ListResourceTagResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceTagResponse), args.Error(1)
}
func (m *mockFindingClient) ListResourceTagName(context.Context, *finding.ListResourceTagNameRequest, ...grpc.CallOption) (*finding.ListResourceTagNameResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceTagNameResponse), args.Error(1)
}
func (m *mockFindingClient) TagResource(context.Context, *finding.TagResourceRequest, ...grpc.CallOption) (*finding.TagResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.TagResourceResponse), args.Error(1)
}
func (m *mockFindingClient) UntagResource(context.Context, *finding.UntagResourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) GetPendFinding(context.Context, *finding.GetPendFindingRequest, ...grpc.CallOption) (*finding.GetPendFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetPendFindingResponse), args.Error(1)
}
func (m *mockFindingClient) PutPendFinding(context.Context, *finding.PutPendFindingRequest, ...grpc.CallOption) (*finding.PutPendFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutPendFindingResponse), args.Error(1)
}
func (m *mockFindingClient) DeletePendFinding(context.Context, *finding.DeletePendFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListFindingSetting(context.Context, *finding.ListFindingSettingRequest, ...grpc.CallOption) (*finding.ListFindingSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingSettingResponse), args.Error(1)
}
func (m *mockFindingClient) GetFindingSetting(context.Context, *finding.GetFindingSettingRequest, ...grpc.CallOption) (*finding.GetFindingSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetFindingSettingResponse), args.Error(1)
}
func (m *mockFindingClient) PutFindingSetting(context.Context, *finding.PutFindingSettingRequest, ...grpc.CallOption) (*finding.PutFindingSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutFindingSettingResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteFindingSetting(context.Context, *finding.DeleteFindingSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ClearScore(context.Context, *finding.ClearScoreRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) GetRecommend(context.Context, *finding.GetRecommendRequest, ...grpc.CallOption) (*finding.GetRecommendResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetRecommendResponse), args.Error(1)
}
func (m *mockFindingClient) PutRecommend(context.Context, *finding.PutRecommendRequest, ...grpc.CallOption) (*finding.PutRecommendResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutRecommendResponse), args.Error(1)
}
