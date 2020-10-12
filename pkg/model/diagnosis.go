package main

import (
	"time"
)

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
	DiagnosisDataSourceID uint32 `gorm:"primary_key"`
	Name                  string
	Description           string
	MaxScore              float32
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// JiraSetting entity
type JiraSetting struct {
	JiraSettingID         uint32 `gorm:"primary_key"`
	Name                  string
	DiagnosisDataSourceID uint32
	ProjectID             uint32
	IdentityField         string
	IdentityValue         string
	JiraID                string
	JiraKey               string
	Status                string
	StatusDetail          string
	ScanAt                time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
