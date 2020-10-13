.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-service
list-service:
	grpcurl -plaintext localhost:19001 list diagnosis.DiagnosisService

.PHONY: list-diagnosis_datasource
list-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:19001 diagnosis.DiagnosisService.ListDiagnosisDataSource

.PHONY: get-diagnosis_datasource
get-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "diagnosis_data_source_id":1001}' \
		localhost:19001 diagnosis.DiagnosisService.GetDiagnosisDataSource

.PHONY: put-diagnosis_datasource
put-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001 ,"diagnosis_data_source":{"name":"test_ds","description":"for_test","max_score":10}}' \
		localhost:19001 diagnosis.DiagnosisService.PutDiagnosisDataSource

.PHONY: delete-diagnosis_datasource
delete-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "diagnosis_data_source_id":1002}' \
		localhost:19001 diagnosis.DiagnosisService.DeleteDiagnosisDataSource

.PHONY: list-jira_setting
list-jira_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:19001 diagnosis.DiagnosisService.ListJiraSetting

.PHONY: get-jira_setting
get-jira_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "jira_setting_id":1001}' \
		localhost:19001 diagnosis.DiagnosisService.GetJiraSetting

.PHONY: put-jira_setting
put-jira_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "jira_setting": {"name":"hoge_name", "diagnosis_data_source_id":1001, "project_id":1001, "identity_field":"field_hoge", "identity_value":"value_hoge","jira_id":"hoge_id","jira_key":"hoge_key", "status":"CONFIGURED"}}' \
		localhost:19001 diagnosis.DiagnosisService.PutJiraSetting

.PHONY: delete-jira_setting
delete-jira_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "jira_setting_id":1001}' \
		localhost:19001 diagnosis.DiagnosisService.DeleteJiraSetting

.PHONY: start-diagnosis
start-diagnosis:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "jira_setting_id":1003}' \
		localhost:19001 diagnosis.DiagnosisService.StartDiagnosis
