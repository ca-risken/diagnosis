package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/common"
	"github.com/ca-risken/diagnosis/pkg/message"
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			want: &finding.FindingForUpsert{
				Description:      "description",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource:     common.DataSourceNameWPScan,
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "to_s",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "readme",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource: common.DataSourceNameWPScan,
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress version num identified",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress version num identified (Insecure)",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource: common.DataSourceNameWPScan,
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
				t.Fatalf("Unexpected recommend: wantErr=%v, error=%v", c.wantErr, e)
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is closed.",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource:      common.DataSourceNameWPScan,
				WpscanSettingID: 1,
				ProjectID:       1,
				TargetURL:       "http://localhost",
				Options:         "",
				ScanOnly:        true,
			},
			finding: &finding.FindingForUpsert{
				Description:      "WordPress login page is opened.",
				DataSource:       common.DataSourceNameWPScan,
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
				DataSource: common.DataSourceNameWPScan,
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
				t.Fatalf("Unexpected recommend: wantErr=%v, error=%v", c.wantErr, e)
			}
		})
	}
}
