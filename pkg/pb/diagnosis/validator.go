package diagnosis

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Diagnosis Service
// Validate ListDiagnosisRequest
func (r *ListDiagnosisRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetDiagnosisRequest
func (r *GetDiagnosisRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.DiagnosisId, validation.Required),
	)
}

// Validate PutDiagnosisRequest
func (r *PutDiagnosisRequest) Validate() error {
	if r.Diagnosis == nil {
		return errors.New("Required Diagnosis")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	); err != nil {
		return err
	}
	return r.Diagnosis.Validate()
}

// Validate DeleteDiagnosisRequest
func (r *DeleteDiagnosisRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.DiagnosisId, validation.Required),
	)
}

// DiagnosisDataSource Service
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

// RelDiagnosisDataSource Service
// Validate ListRelDiagnosisDataSourceRequest
func (r *ListRelDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetRelDiagnosisDataSourceRequest
func (r *GetRelDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelDiagnosisDataSourceId, validation.Required),
	)
}

// Validate PutRelDiagnosisDataSourceRequest
func (r *PutRelDiagnosisDataSourceRequest) Validate() error {
	if r.RelDiagnosisDataSource == nil {
		return errors.New("Required RelDiagnosisDataSource")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	); err != nil {
		return err
	}
	return r.RelDiagnosisDataSource.Validate()
}

// Validate DeleteResultRequest
func (r *DeleteRelDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelDiagnosisDataSourceId, validation.Required),
	)
}

// Validate StartDiagnosisRequest
func (r *StartDiagnosisRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelDiagnosisDataSourceId, validation.Required),
	)
}

/**
 * Entity
**/

// Validate DiagnosisForUpsert
func (o *DiagnosisForUpsert) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Name, validation.Required, validation.Length(0, 200)),
	)
}

// Validate DiagnosisDataSourceForUpsert
func (d *DiagnosisDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&d.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&d.MaxScore, validation.Required, validation.Min(0.0), validation.Max(99999.0)),
	)
}

// Validate RelDiagnosisDataSourceForUpsert
func (r *RelDiagnosisDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DiagnosisId, validation.Required),
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
	)
}
