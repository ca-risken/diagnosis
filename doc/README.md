# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [diagnosis/entities.proto](#diagnosis/entities.proto)
    - [ApplicationScan](#diagnosis.ApplicationScan)
    - [ApplicationScanBasicSetting](#diagnosis.ApplicationScanBasicSetting)
    - [ApplicationScanBasicSettingForUpsert](#diagnosis.ApplicationScanBasicSettingForUpsert)
    - [ApplicationScanForUpsert](#diagnosis.ApplicationScanForUpsert)
    - [DiagnosisDataSource](#diagnosis.DiagnosisDataSource)
    - [DiagnosisDataSourceForUpsert](#diagnosis.DiagnosisDataSourceForUpsert)
    - [JiraSetting](#diagnosis.JiraSetting)
    - [JiraSettingForUpsert](#diagnosis.JiraSettingForUpsert)
    - [PortscanSetting](#diagnosis.PortscanSetting)
    - [PortscanSettingForUpsert](#diagnosis.PortscanSettingForUpsert)
    - [PortscanTarget](#diagnosis.PortscanTarget)
    - [PortscanTargetForUpsert](#diagnosis.PortscanTargetForUpsert)
    - [WpscanSetting](#diagnosis.WpscanSetting)
    - [WpscanSettingForUpsert](#diagnosis.WpscanSettingForUpsert)
  
    - [ApplicationScanType](#diagnosis.ApplicationScanType)
    - [Status](#diagnosis.Status)
  
- [diagnosis/services.proto](#diagnosis/services.proto)
    - [DeleteApplicationScanBasicSettingRequest](#diagnosis.DeleteApplicationScanBasicSettingRequest)
    - [DeleteApplicationScanRequest](#diagnosis.DeleteApplicationScanRequest)
    - [DeleteDiagnosisDataSourceRequest](#diagnosis.DeleteDiagnosisDataSourceRequest)
    - [DeleteJiraSettingRequest](#diagnosis.DeleteJiraSettingRequest)
    - [DeletePortscanSettingRequest](#diagnosis.DeletePortscanSettingRequest)
    - [DeletePortscanTargetRequest](#diagnosis.DeletePortscanTargetRequest)
    - [DeleteWpscanSettingRequest](#diagnosis.DeleteWpscanSettingRequest)
    - [GetApplicationScanBasicSettingRequest](#diagnosis.GetApplicationScanBasicSettingRequest)
    - [GetApplicationScanBasicSettingResponse](#diagnosis.GetApplicationScanBasicSettingResponse)
    - [GetApplicationScanRequest](#diagnosis.GetApplicationScanRequest)
    - [GetApplicationScanResponse](#diagnosis.GetApplicationScanResponse)
    - [GetDiagnosisDataSourceRequest](#diagnosis.GetDiagnosisDataSourceRequest)
    - [GetDiagnosisDataSourceResponse](#diagnosis.GetDiagnosisDataSourceResponse)
    - [GetJiraSettingRequest](#diagnosis.GetJiraSettingRequest)
    - [GetJiraSettingResponse](#diagnosis.GetJiraSettingResponse)
    - [GetPortscanSettingRequest](#diagnosis.GetPortscanSettingRequest)
    - [GetPortscanSettingResponse](#diagnosis.GetPortscanSettingResponse)
    - [GetPortscanTargetRequest](#diagnosis.GetPortscanTargetRequest)
    - [GetPortscanTargetResponse](#diagnosis.GetPortscanTargetResponse)
    - [GetWpscanSettingRequest](#diagnosis.GetWpscanSettingRequest)
    - [GetWpscanSettingResponse](#diagnosis.GetWpscanSettingResponse)
    - [InvokeScanRequest](#diagnosis.InvokeScanRequest)
    - [InvokeScanResponse](#diagnosis.InvokeScanResponse)
    - [ListApplicationScanBasicSettingRequest](#diagnosis.ListApplicationScanBasicSettingRequest)
    - [ListApplicationScanBasicSettingResponse](#diagnosis.ListApplicationScanBasicSettingResponse)
    - [ListApplicationScanRequest](#diagnosis.ListApplicationScanRequest)
    - [ListApplicationScanResponse](#diagnosis.ListApplicationScanResponse)
    - [ListDiagnosisDataSourceRequest](#diagnosis.ListDiagnosisDataSourceRequest)
    - [ListDiagnosisDataSourceResponse](#diagnosis.ListDiagnosisDataSourceResponse)
    - [ListJiraSettingRequest](#diagnosis.ListJiraSettingRequest)
    - [ListJiraSettingResponse](#diagnosis.ListJiraSettingResponse)
    - [ListPortscanSettingRequest](#diagnosis.ListPortscanSettingRequest)
    - [ListPortscanSettingResponse](#diagnosis.ListPortscanSettingResponse)
    - [ListPortscanTargetRequest](#diagnosis.ListPortscanTargetRequest)
    - [ListPortscanTargetResponse](#diagnosis.ListPortscanTargetResponse)
    - [ListWpscanSettingRequest](#diagnosis.ListWpscanSettingRequest)
    - [ListWpscanSettingResponse](#diagnosis.ListWpscanSettingResponse)
    - [PutApplicationScanBasicSettingRequest](#diagnosis.PutApplicationScanBasicSettingRequest)
    - [PutApplicationScanBasicSettingResponse](#diagnosis.PutApplicationScanBasicSettingResponse)
    - [PutApplicationScanRequest](#diagnosis.PutApplicationScanRequest)
    - [PutApplicationScanResponse](#diagnosis.PutApplicationScanResponse)
    - [PutDiagnosisDataSourceRequest](#diagnosis.PutDiagnosisDataSourceRequest)
    - [PutDiagnosisDataSourceResponse](#diagnosis.PutDiagnosisDataSourceResponse)
    - [PutJiraSettingRequest](#diagnosis.PutJiraSettingRequest)
    - [PutJiraSettingResponse](#diagnosis.PutJiraSettingResponse)
    - [PutPortscanSettingRequest](#diagnosis.PutPortscanSettingRequest)
    - [PutPortscanSettingResponse](#diagnosis.PutPortscanSettingResponse)
    - [PutPortscanTargetRequest](#diagnosis.PutPortscanTargetRequest)
    - [PutPortscanTargetResponse](#diagnosis.PutPortscanTargetResponse)
    - [PutWpscanSettingRequest](#diagnosis.PutWpscanSettingRequest)
    - [PutWpscanSettingResponse](#diagnosis.PutWpscanSettingResponse)
  
    - [DiagnosisService](#diagnosis.DiagnosisService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="diagnosis/entities.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## diagnosis/entities.proto



<a name="diagnosis.ApplicationScan"></a>

### ApplicationScan



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |
| scan_type | [ApplicationScanType](#diagnosis.ApplicationScanType) |  |  |






<a name="diagnosis.ApplicationScanBasicSetting"></a>

### ApplicationScanBasicSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_basic_setting_id | [uint32](#uint32) |  |  |
| application_scan_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| target | [string](#string) |  |  |
| max_depth | [uint32](#uint32) |  |  |
| max_children | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.ApplicationScanBasicSettingForUpsert"></a>

### ApplicationScanBasicSettingForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_basic_setting_id | [uint32](#uint32) |  |  |
| application_scan_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| target | [string](#string) |  |  |
| max_depth | [uint32](#uint32) |  |  |
| max_children | [uint32](#uint32) |  |  |






<a name="diagnosis.ApplicationScanForUpsert"></a>

### ApplicationScanForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |
| scan_type | [ApplicationScanType](#diagnosis.ApplicationScanType) |  |  |






<a name="diagnosis.DiagnosisDataSource"></a>

### DiagnosisDataSource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| max_score | [float](#float) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.DiagnosisDataSourceForUpsert"></a>

### DiagnosisDataSourceForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| max_score | [float](#float) |  |  |






<a name="diagnosis.JiraSetting"></a>

### JiraSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jira_setting_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| identity_field | [string](#string) |  |  |
| identity_value | [string](#string) |  |  |
| jira_id | [string](#string) |  |  |
| jira_key | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.JiraSettingForUpsert"></a>

### JiraSettingForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jira_setting_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| identity_field | [string](#string) |  |  |
| identity_value | [string](#string) |  |  |
| jira_id | [string](#string) |  |  |
| jira_key | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |






<a name="diagnosis.PortscanSetting"></a>

### PortscanSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_setting_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.PortscanSettingForUpsert"></a>

### PortscanSettingForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_setting_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="diagnosis.PortscanTarget"></a>

### PortscanTarget



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_target_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| portscan_setting_id | [uint32](#uint32) |  |  |
| target | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.PortscanTargetForUpsert"></a>

### PortscanTargetForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_target_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| portscan_setting_id | [uint32](#uint32) |  |  |
| target | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |






<a name="diagnosis.WpscanSetting"></a>

### WpscanSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wpscan_setting_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| target_url | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.WpscanSettingForUpsert"></a>

### WpscanSettingForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wpscan_setting_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| target_url | [string](#string) |  |  |
| status | [Status](#diagnosis.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |





 


<a name="diagnosis.ApplicationScanType"></a>

### ApplicationScanType
Status

| Name | Number | Description |
| ---- | ------ | ----------- |
| BASIC | 0 |  |



<a name="diagnosis.Status"></a>

### Status
Status

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| OK | 1 |  |
| CONFIGURED | 2 |  |
| IN_PROGRESS | 3 |  |
| ERROR | 4 |  |


 

 

 



<a name="diagnosis/services.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## diagnosis/services.proto



<a name="diagnosis.DeleteApplicationScanBasicSettingRequest"></a>

### DeleteApplicationScanBasicSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan_basic_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeleteApplicationScanRequest"></a>

### DeleteApplicationScanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeleteDiagnosisDataSourceRequest"></a>

### DeleteDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeleteJiraSettingRequest"></a>

### DeleteJiraSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| jira_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeletePortscanSettingRequest"></a>

### DeletePortscanSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeletePortscanTargetRequest"></a>

### DeletePortscanTargetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_target_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeleteWpscanSettingRequest"></a>

### DeleteWpscanSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| wpscan_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetApplicationScanBasicSettingRequest"></a>

### GetApplicationScanBasicSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan_basic_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetApplicationScanBasicSettingResponse"></a>

### GetApplicationScanBasicSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_basic_setting | [ApplicationScanBasicSetting](#diagnosis.ApplicationScanBasicSetting) |  |  |






<a name="diagnosis.GetApplicationScanRequest"></a>

### GetApplicationScanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetApplicationScanResponse"></a>

### GetApplicationScanResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan | [ApplicationScan](#diagnosis.ApplicationScan) |  |  |






<a name="diagnosis.GetDiagnosisDataSourceRequest"></a>

### GetDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetDiagnosisDataSourceResponse"></a>

### GetDiagnosisDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis_data_source | [DiagnosisDataSource](#diagnosis.DiagnosisDataSource) |  |  |






<a name="diagnosis.GetJiraSettingRequest"></a>

### GetJiraSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| jira_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetJiraSettingResponse"></a>

### GetJiraSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jira_setting | [JiraSetting](#diagnosis.JiraSetting) |  |  |






<a name="diagnosis.GetPortscanSettingRequest"></a>

### GetPortscanSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetPortscanSettingResponse"></a>

### GetPortscanSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_setting | [PortscanSetting](#diagnosis.PortscanSetting) |  |  |






<a name="diagnosis.GetPortscanTargetRequest"></a>

### GetPortscanTargetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_target_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetPortscanTargetResponse"></a>

### GetPortscanTargetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_target | [PortscanTarget](#diagnosis.PortscanTarget) |  |  |






<a name="diagnosis.GetWpscanSettingRequest"></a>

### GetWpscanSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| wpscan_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetWpscanSettingResponse"></a>

### GetWpscanSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wpscan_setting | [WpscanSetting](#diagnosis.WpscanSetting) |  |  |






<a name="diagnosis.InvokeScanRequest"></a>

### InvokeScanRequest
KICK Diagnosis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| setting_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| scan_only | [bool](#bool) |  |  |






<a name="diagnosis.InvokeScanResponse"></a>

### InvokeScanResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [string](#string) |  |  |






<a name="diagnosis.ListApplicationScanBasicSettingRequest"></a>

### ListApplicationScanBasicSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListApplicationScanBasicSettingResponse"></a>

### ListApplicationScanBasicSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_basic_setting | [ApplicationScanBasicSetting](#diagnosis.ApplicationScanBasicSetting) | repeated |  |






<a name="diagnosis.ListApplicationScanRequest"></a>

### ListApplicationScanRequest
ApplicationScanService


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListApplicationScanResponse"></a>

### ListApplicationScanResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan | [ApplicationScan](#diagnosis.ApplicationScan) | repeated |  |






<a name="diagnosis.ListDiagnosisDataSourceRequest"></a>

### ListDiagnosisDataSourceRequest
DiagnosisDataSourceService


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="diagnosis.ListDiagnosisDataSourceResponse"></a>

### ListDiagnosisDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis_data_source | [DiagnosisDataSource](#diagnosis.DiagnosisDataSource) | repeated |  |






<a name="diagnosis.ListJiraSettingRequest"></a>

### ListJiraSettingRequest
JiraSettingService


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListJiraSettingResponse"></a>

### ListJiraSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jira_setting | [JiraSetting](#diagnosis.JiraSetting) | repeated |  |






<a name="diagnosis.ListPortscanSettingRequest"></a>

### ListPortscanSettingRequest
PortscanSettingService


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListPortscanSettingResponse"></a>

### ListPortscanSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_setting | [PortscanSetting](#diagnosis.PortscanSetting) | repeated |  |






<a name="diagnosis.ListPortscanTargetRequest"></a>

### ListPortscanTargetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_setting_id | [uint32](#uint32) |  |  |
| status | [Status](#diagnosis.Status) |  |  |






<a name="diagnosis.ListPortscanTargetResponse"></a>

### ListPortscanTargetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_target | [PortscanTarget](#diagnosis.PortscanTarget) | repeated |  |






<a name="diagnosis.ListWpscanSettingRequest"></a>

### ListWpscanSettingRequest
WpscanSettingService


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListWpscanSettingResponse"></a>

### ListWpscanSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wpscan_setting | [WpscanSetting](#diagnosis.WpscanSetting) | repeated |  |






<a name="diagnosis.PutApplicationScanBasicSettingRequest"></a>

### PutApplicationScanBasicSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan_basic_setting | [ApplicationScanBasicSettingForUpsert](#diagnosis.ApplicationScanBasicSettingForUpsert) |  |  |






<a name="diagnosis.PutApplicationScanBasicSettingResponse"></a>

### PutApplicationScanBasicSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan_basic_setting | [ApplicationScanBasicSetting](#diagnosis.ApplicationScanBasicSetting) |  |  |






<a name="diagnosis.PutApplicationScanRequest"></a>

### PutApplicationScanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| application_scan | [ApplicationScanForUpsert](#diagnosis.ApplicationScanForUpsert) |  |  |






<a name="diagnosis.PutApplicationScanResponse"></a>

### PutApplicationScanResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| application_scan | [ApplicationScan](#diagnosis.ApplicationScan) |  |  |






<a name="diagnosis.PutDiagnosisDataSourceRequest"></a>

### PutDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source | [DiagnosisDataSourceForUpsert](#diagnosis.DiagnosisDataSourceForUpsert) |  |  |






<a name="diagnosis.PutDiagnosisDataSourceResponse"></a>

### PutDiagnosisDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source | [DiagnosisDataSource](#diagnosis.DiagnosisDataSource) |  |  |






<a name="diagnosis.PutJiraSettingRequest"></a>

### PutJiraSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| jira_setting | [JiraSettingForUpsert](#diagnosis.JiraSettingForUpsert) |  |  |






<a name="diagnosis.PutJiraSettingResponse"></a>

### PutJiraSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jira_setting | [JiraSetting](#diagnosis.JiraSetting) |  |  |






<a name="diagnosis.PutPortscanSettingRequest"></a>

### PutPortscanSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_setting | [PortscanSettingForUpsert](#diagnosis.PortscanSettingForUpsert) |  |  |






<a name="diagnosis.PutPortscanSettingResponse"></a>

### PutPortscanSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_setting | [PortscanSetting](#diagnosis.PortscanSetting) |  |  |






<a name="diagnosis.PutPortscanTargetRequest"></a>

### PutPortscanTargetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| portscan_target | [PortscanTargetForUpsert](#diagnosis.PortscanTargetForUpsert) |  |  |






<a name="diagnosis.PutPortscanTargetResponse"></a>

### PutPortscanTargetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| portscan_target | [PortscanTarget](#diagnosis.PortscanTarget) |  |  |






<a name="diagnosis.PutWpscanSettingRequest"></a>

### PutWpscanSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| wpscan_setting | [WpscanSettingForUpsert](#diagnosis.WpscanSettingForUpsert) |  |  |






<a name="diagnosis.PutWpscanSettingResponse"></a>

### PutWpscanSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| wpscan_setting | [WpscanSetting](#diagnosis.WpscanSetting) |  |  |





 

 

 


<a name="diagnosis.DiagnosisService"></a>

### DiagnosisService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListDiagnosisDataSource | [ListDiagnosisDataSourceRequest](#diagnosis.ListDiagnosisDataSourceRequest) | [ListDiagnosisDataSourceResponse](#diagnosis.ListDiagnosisDataSourceResponse) | DataSource |
| GetDiagnosisDataSource | [GetDiagnosisDataSourceRequest](#diagnosis.GetDiagnosisDataSourceRequest) | [GetDiagnosisDataSourceResponse](#diagnosis.GetDiagnosisDataSourceResponse) |  |
| PutDiagnosisDataSource | [PutDiagnosisDataSourceRequest](#diagnosis.PutDiagnosisDataSourceRequest) | [PutDiagnosisDataSourceResponse](#diagnosis.PutDiagnosisDataSourceResponse) |  |
| DeleteDiagnosisDataSource | [DeleteDiagnosisDataSourceRequest](#diagnosis.DeleteDiagnosisDataSourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListJiraSetting | [ListJiraSettingRequest](#diagnosis.ListJiraSettingRequest) | [ListJiraSettingResponse](#diagnosis.ListJiraSettingResponse) | JiraSetting |
| GetJiraSetting | [GetJiraSettingRequest](#diagnosis.GetJiraSettingRequest) | [GetJiraSettingResponse](#diagnosis.GetJiraSettingResponse) |  |
| PutJiraSetting | [PutJiraSettingRequest](#diagnosis.PutJiraSettingRequest) | [PutJiraSettingResponse](#diagnosis.PutJiraSettingResponse) |  |
| DeleteJiraSetting | [DeleteJiraSettingRequest](#diagnosis.DeleteJiraSettingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListWpscanSetting | [ListWpscanSettingRequest](#diagnosis.ListWpscanSettingRequest) | [ListWpscanSettingResponse](#diagnosis.ListWpscanSettingResponse) | WpscanSetting |
| GetWpscanSetting | [GetWpscanSettingRequest](#diagnosis.GetWpscanSettingRequest) | [GetWpscanSettingResponse](#diagnosis.GetWpscanSettingResponse) |  |
| PutWpscanSetting | [PutWpscanSettingRequest](#diagnosis.PutWpscanSettingRequest) | [PutWpscanSettingResponse](#diagnosis.PutWpscanSettingResponse) |  |
| DeleteWpscanSetting | [DeleteWpscanSettingRequest](#diagnosis.DeleteWpscanSettingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListPortscanSetting | [ListPortscanSettingRequest](#diagnosis.ListPortscanSettingRequest) | [ListPortscanSettingResponse](#diagnosis.ListPortscanSettingResponse) | PortscanSetting |
| GetPortscanSetting | [GetPortscanSettingRequest](#diagnosis.GetPortscanSettingRequest) | [GetPortscanSettingResponse](#diagnosis.GetPortscanSettingResponse) |  |
| PutPortscanSetting | [PutPortscanSettingRequest](#diagnosis.PutPortscanSettingRequest) | [PutPortscanSettingResponse](#diagnosis.PutPortscanSettingResponse) |  |
| DeletePortscanSetting | [DeletePortscanSettingRequest](#diagnosis.DeletePortscanSettingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListPortscanTarget | [ListPortscanTargetRequest](#diagnosis.ListPortscanTargetRequest) | [ListPortscanTargetResponse](#diagnosis.ListPortscanTargetResponse) | PortscanTarget |
| GetPortscanTarget | [GetPortscanTargetRequest](#diagnosis.GetPortscanTargetRequest) | [GetPortscanTargetResponse](#diagnosis.GetPortscanTargetResponse) |  |
| PutPortscanTarget | [PutPortscanTargetRequest](#diagnosis.PutPortscanTargetRequest) | [PutPortscanTargetResponse](#diagnosis.PutPortscanTargetResponse) |  |
| DeletePortscanTarget | [DeletePortscanTargetRequest](#diagnosis.DeletePortscanTargetRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListApplicationScan | [ListApplicationScanRequest](#diagnosis.ListApplicationScanRequest) | [ListApplicationScanResponse](#diagnosis.ListApplicationScanResponse) | ApplicationScan |
| GetApplicationScan | [GetApplicationScanRequest](#diagnosis.GetApplicationScanRequest) | [GetApplicationScanResponse](#diagnosis.GetApplicationScanResponse) |  |
| PutApplicationScan | [PutApplicationScanRequest](#diagnosis.PutApplicationScanRequest) | [PutApplicationScanResponse](#diagnosis.PutApplicationScanResponse) |  |
| DeleteApplicationScan | [DeleteApplicationScanRequest](#diagnosis.DeleteApplicationScanRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListApplicationScanBasicSetting | [ListApplicationScanBasicSettingRequest](#diagnosis.ListApplicationScanBasicSettingRequest) | [ListApplicationScanBasicSettingResponse](#diagnosis.ListApplicationScanBasicSettingResponse) | ApplicationScanBasicSetting |
| GetApplicationScanBasicSetting | [GetApplicationScanBasicSettingRequest](#diagnosis.GetApplicationScanBasicSettingRequest) | [GetApplicationScanBasicSettingResponse](#diagnosis.GetApplicationScanBasicSettingResponse) |  |
| PutApplicationScanBasicSetting | [PutApplicationScanBasicSettingRequest](#diagnosis.PutApplicationScanBasicSettingRequest) | [PutApplicationScanBasicSettingResponse](#diagnosis.PutApplicationScanBasicSettingResponse) |  |
| DeleteApplicationScanBasicSetting | [DeleteApplicationScanBasicSettingRequest](#diagnosis.DeleteApplicationScanBasicSettingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| InvokeScan | [InvokeScanRequest](#diagnosis.InvokeScanRequest) | [InvokeScanResponse](#diagnosis.InvokeScanResponse) | KICK |
| InvokeScanAll | [.google.protobuf.Empty](#google.protobuf.Empty) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

