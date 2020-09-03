package diagnosis

import (
	"testing"
)

// Diagnosis
func TestValidate_ListDiagnosisRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListDiagnosisRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListDiagnosisRequest{ProjectId: 111},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListDiagnosisRequest{},
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

func TestValidate_GetDiagnosisRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetDiagnosisRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetDiagnosisRequest{ProjectId: 111, DiagnosisId: 222},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetDiagnosisRequest{DiagnosisId: 222},
			wantErr: true,
		},
		{
			name:    "NG required(diagnosis_id)",
			input:   &GetDiagnosisRequest{ProjectId: 111},
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

func TestValidate_PutDiagnosisRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutDiagnosisRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutDiagnosisRequest{ProjectId: 1, Diagnosis: &DiagnosisForUpsert{Name: "name"}},
			wantErr: false,
		},
		{
			name:    "NG Required(Diagnosis)",
			input:   &PutDiagnosisRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutDiagnosisRequest{Diagnosis: &DiagnosisForUpsert{Name: "name"}},
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

func TestValidate_DeleteDiagnosisRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteDiagnosisRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteDiagnosisRequest{ProjectId: 1, DiagnosisId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteDiagnosisRequest{DiagnosisId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_id)",
			input:   &DeleteDiagnosisRequest{ProjectId: 1},
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

//RelDiagnosisDataSource DataSource

func TestValidate_ListRelDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListRelDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListRelDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListRelDiagnosisDataSourceRequest{},
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

func TestValidate_GetRelDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetRelDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetRelDiagnosisDataSourceRequest{ProjectId: 1, RelDiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetRelDiagnosisDataSourceRequest{RelDiagnosisDataSourceId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(diagnosis_result_id)",
			input:   &GetRelDiagnosisDataSourceRequest{ProjectId: 1},
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

func TestValidate_PutRelDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutRelDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutRelDiagnosisDataSourceRequest{ProjectId: 1, RelDiagnosisDataSource: &RelDiagnosisDataSourceForUpsert{DiagnosisId: 1, DiagnosisDataSourceId: 2, RecordId: "record_id", JiraId: "jira_id", JiraKey: "jira_key"}},
			wantErr: false,
		},
		{
			name:    "NG Required(RelDiagnosisDataSource)",
			input:   &PutRelDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutRelDiagnosisDataSourceRequest{RelDiagnosisDataSource: &RelDiagnosisDataSourceForUpsert{DiagnosisId: 1, DiagnosisDataSourceId: 2, RecordId: "record_id", JiraId: "jira_id", JiraKey: "jira_key"}},
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

func TestValidate_DeleteRelDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteRelDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteRelDiagnosisDataSourceRequest{ProjectId: 1, RelDiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteRelDiagnosisDataSourceRequest{RelDiagnosisDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeleteRelDiagnosisDataSourceRequest{ProjectId: 1},
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
			input:   &StartDiagnosisRequest{ProjectId: 1, RelDiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &StartDiagnosisRequest{RelDiagnosisDataSourceId: 2},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_result_id)",
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

func TestValidate_DiagnosisForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *DiagnosisForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DiagnosisForUpsert{Name: "name"},
			wantErr: false,
		},
		{
			name:    "NG Length(name)",
			input:   &DiagnosisForUpsert{Name: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &DiagnosisForUpsert{},
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

func TestValidate_RelDiagnosisDataSourceForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *RelDiagnosisDataSourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &RelDiagnosisDataSourceForUpsert{DiagnosisId: 1, DiagnosisDataSourceId: 2, RecordId: "record_id", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: false,
		},
		{
			name:    "NG Required(diagnosis_id)",
			input:   &RelDiagnosisDataSourceForUpsert{DiagnosisDataSourceId: 2, RecordId: "record_id", JiraId: "jira_id", JiraKey: "jira_key"},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &RelDiagnosisDataSourceForUpsert{DiagnosisId: 1, RecordId: "record_id", JiraId: "jira_id", JiraKey: "jira_key"},
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
