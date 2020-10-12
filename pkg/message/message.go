package message

// DiagnosisQueueMessage is the message for SQS queue
type DiagnosisQueueMessage struct {
	DataSource    string `json:"data_source"`
	JiraSettingID uint32 `json:"jira_setting_id"`
	ProjectID     uint32 `json:"project_id"`
	IdentityField string `json:"identity_field"`
	IdentityValue string `json:"identity_value"`
	JiraID        string `json:"jira_id"`
	JiraKey       string `json:"jira_key"`
}
