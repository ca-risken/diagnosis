package diagnosis

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// DiagnosisDataSource

// Validate ListDiagnosisDataSourceRequest
func (r *ListDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetDiagnosisDataSourceRequest
func (r *GetDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
	)
}

// Validate PutDiagnosisDataSourceRequest
func (r *PutDiagnosisDataSourceRequest) Validate() error {
	if r.DiagnosisDataSource == nil {
		return errors.New("Required DiagnosisDataSource")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	); err != nil {
		return err
	}
	return r.DiagnosisDataSource.Validate()
}

// Validate DeleteDataSourceRequest
func (r *DeleteDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
	)
}

// JiraSetting

// Validate ListJiraSettingRequest
func (r *ListJiraSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetJiraSettingRequest
func (r *GetJiraSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.JiraSettingId, validation.Required),
	)
}

// Validate PutJiraSettingRequest
func (r *PutJiraSettingRequest) Validate() error {
	if r.JiraSetting == nil {
		return errors.New("Required JiraSetting")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.JiraSetting.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.JiraSetting.Validate()
}

// Validate DeleteResultRequest
func (r *DeleteJiraSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.JiraSettingId, validation.Required),
	)
}

// Validate StartDiagnosisRequest
func (r *StartDiagnosisRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.JiraSettingId, validation.Required),
	)
}

/**
 * Entity
**/

// Validate DiagnosisDataSourceForUpsert
func (d *DiagnosisDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&d.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&d.MaxScore, validation.Required, validation.Min(0.0), validation.Max(99999.0)),
	)
}

// Validate JiraSettingForUpsert
func (r *JiraSettingForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&r.IdentityField, validation.Length(0, 50)),
		validation.Field(&r.IdentityValue, validation.Required.When(r.IdentityField != "")),
		validation.Field(&r.IdentityValue, validation.Required.When(r.JiraId == "" && r.JiraKey == "").Error("Any one IdentityValue, JiraId, JiraKey is required."), validation.Length(0, 50)),
		validation.Field(&r.JiraId, validation.Required.When(r.IdentityValue == "" && r.JiraKey == "").Error("Any one IdentityValue, JiraId, JiraKey is required."), validation.Length(0, 50)),
		validation.Field(&r.JiraKey, validation.Required.When(r.IdentityValue == "" && r.JiraId == "").Error("Any one IdentityValue, JiraId, JiraKey is required."), validation.Length(0, 50)),
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
	)
}
