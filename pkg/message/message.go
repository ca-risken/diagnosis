package message

// JiraQueueMessage is the message for SQS queue for Jira
type JiraQueueMessage struct {
	DataSource    string `json:"data_source"`
	JiraSettingID uint32 `json:"jira_setting_id"`
	ProjectID     uint32 `json:"project_id"`
	IdentityField string `json:"identity_field"`
	IdentityValue string `json:"identity_value"`
	JiraID        string `json:"jira_id"`
	JiraKey       string `json:"jira_key"`
	ScanOnly      bool   `json:"scan_only,string"`
}

// WpscanQueueMessage is the message for SQS queue for Wpscan
type WpscanQueueMessage struct {
	DataSource      string `json:"data_source"`
	WpscanSettingID uint32 `json:"jira_setting_id"`
	ProjectID       uint32 `json:"project_id"`
	TargetURL       string `json:"target_url"`
	ScanOnly        bool   `json:"scan_only,string"`
}

// PortscanQueueMessage is the message for SQS queue for Portscan
type PortscanQueueMessage struct {
	DataSource        string `json:"data_source"`
	PortscanSettingID uint32 `json:"portscan_setting_id"`
	PortscanTargetID  uint32 `json:"portscan_target_id"`
	ProjectID         uint32 `json:"project_id"`
	Target            string `json:"target"`
	ScanOnly          bool   `json:"scan_only,string"`
}
