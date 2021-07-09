package model

import (
	"time"
)

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

// WpscanSetting Entity
type WpscanSetting struct {
	WpscanSettingID       uint32 `gorm:"primary_key"`
	DiagnosisDataSourceID uint32
	ProjectID             uint32
	TargetURL             string
	Status                string
	StatusDetail          string
	ScanAt                time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// PortscanSetting Entity
type PortscanSetting struct {
	PortscanSettingID     uint32 `gorm:"primary_key"`
	DiagnosisDataSourceID uint32
	ProjectID             uint32
	Name                  string
	Status                string
	StatusDetail          string
	ScanAt                time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// PortscanTarget Entity
type PortscanTarget struct {
	PortscanTargetID  uint32 `gorm:"primary_key"`
	PortscanSettingID uint32
	ProjectID         uint32
	Target            string
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
