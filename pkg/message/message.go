package message

const Jira = "diagnosis:jira"

// DiagnosisQueueMessage is the message for SQS queue
type DiagnosisQueueMessage struct {
	DataSource string `json:"data_source"`
	ProjectID  uint32 `json:"project_id"`
	RecordID   string `json:"record_id"`
	JiraID     string `json:"jira_id"`
}
