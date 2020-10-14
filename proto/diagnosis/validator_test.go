package diagnosis

import (
	"testing"
)

//DiagnosisDataSource DataSource

func TestValidate_ListDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListDiagnosisDataSourceRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetDiagnosisDataSourceRequest{ProjectId: 1, DiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetDiagnosisDataSourceRequest{DiagnosisDataSourceId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(diagnosis_data_source_id)",
			input:   &GetDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutDiagnosisDataSourceRequest{ProjectId: 1, DiagnosisDataSource: &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 10.0}},
			wantErr: false,
		},
		{
			name:    "NG Required(DiagnosisDataSource)",
			input:   &PutDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutDiagnosisDataSourceRequest{DiagnosisDataSource: &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 10.0}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteDiagnosisDataSourceRequest{ProjectId: 1, DiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteDiagnosisDataSourceRequest{DiagnosisDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeleteDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

//JiraSetting DataSource

func TestValidate_ListJiraSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListJiraSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListJiraSettingRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListJiraSettingRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetJiraSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetJiraSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetJiraSettingRequest{ProjectId: 1, JiraSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetJiraSettingRequest{JiraSettingId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(jira_setting_id)",
			input:   &GetJiraSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutJiraSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutJiraSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutJiraSettingRequest{ProjectId: 1001, JiraSetting: &JiraSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_name", IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"}},
			wantErr: false,
		},
		{
			name:    "NG Required(JiraSetting)",
			input:   &PutJiraSettingRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != jira_setting.project_id)",
			input:   &PutJiraSettingRequest{ProjectId: 1002, JiraSetting: &JiraSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_name", IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutJiraSettingRequest{JiraSetting: &JiraSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_name", IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteJiraSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteJiraSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteJiraSettingRequest{ProjectId: 1, JiraSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteJiraSettingRequest{JiraSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeleteJiraSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_StartDiagnosisRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *StartDiagnosisRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &StartDiagnosisRequest{ProjectId: 1, JiraSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &StartDiagnosisRequest{JiraSettingId: 2},
			wantErr: true,
		},
		{
			name:    "NG Required(jira_setting_id)",
			input:   &StartDiagnosisRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DiagnosisDataSourceForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *DiagnosisDataSourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 100},
			wantErr: false,
		},
		{
			name:    "NG Required(name)",
			input:   &DiagnosisDataSourceForUpsert{Description: "description", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Length(description)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Required(description)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Num Over(Max Score)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 100000},
			wantErr: true,
		},
		{
			name:    "NG Num Under(Max Score)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: -1.0},
			wantErr: true,
		},
		{
			name:    "NG Required(Max Score)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_JiraSettingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *JiraSettingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: false,
		},
		{
			name:    "NG Required(name)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "Too long(name)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "123456789012345678901234567890123456789012345678901", DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &JiraSettingForUpsert{DiagnosisDataSourceId: 1, Name: "hoge_name", IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "Too long(identity_field)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "123456789012345678901234567890123456789012345678901", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "NG Required(identity_value when identity_field not blank)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "Too long(identity_value)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "123456789012345678901234567890123456789012345678901", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "NG Required(identity_value,jira_id,jira_key)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "", IdentityValue: "", JiraId: "", JiraKey: ""},
			wantErr: true,
		},
		{
			name:    "Too long(jira_id)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "", IdentityValue: "", JiraId: "123456789012345678901234567890123456789012345678901", JiraKey: ""},
			wantErr: true,
		},
		{
			name:    "Too long(jira_key)",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "", IdentityValue: "", JiraId: "", JiraKey: "123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
		{
			name:    "NG Too small scan_at",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key", ScanAt: -1},
			wantErr: true,
		},
		{
			name:    "NG Too large scan_at",
			input:   &JiraSettingForUpsert{ProjectId: 1001, Name: "hoge_name", DiagnosisDataSourceId: 1, IdentityField: "hoge_field", IdentityValue: "hoge_value", JiraId: "jira_id", JiraKey: "jira_key", ScanAt: 253402268400},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
