package main

import "time"

// Diagnosis entity
type Diagnosis struct {
	DiagnosisID uint32 `gorm:"column:diagnosis_id"`
	ProjectID   uint32
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DiagnosisDataSource entity
type DiagnosisDataSource struct {
	DiagnosisDataSourceID uint32 `gorm:"column:diagnosis_data_source_id"`
	Name                  string
	Description           string
	MaxScore              float32
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// RelDiagnosisDataSource entity
type RelDiagnosisDataSource struct {
	RelDiagnosisDataSourceID uint32 `gorm:"column:rel_diagnosis_data_source_id"`
	DiagnosisID              uint32 `gorm:"column:diagnosis_id"`
	DiagnosisDataSourceID    uint32 `gorm:"column:diagnosis_data_source_id"`
	ProjectID                uint32
	RecordID                 string
	JiraID                   string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
