# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [diagnosis/entities.proto](#diagnosis/entities.proto)
    - [DiagnosisDataSource](#diagnosis.DiagnosisDataSource)
    - [DiagnosisDataSourceForUpsert](#diagnosis.DiagnosisDataSourceForUpsert)
    - [JiraSetting](#diagnosis.JiraSetting)
    - [JiraSettingForUpsert](#diagnosis.JiraSettingForUpsert)
  
    - [Status](#diagnosis.Status)
  
- [diagnosis/services.proto](#diagnosis/services.proto)
    - [DeleteDiagnosisDataSourceRequest](#diagnosis.DeleteDiagnosisDataSourceRequest)
    - [DeleteJiraSettingRequest](#diagnosis.DeleteJiraSettingRequest)
    - [GetDiagnosisDataSourceRequest](#diagnosis.GetDiagnosisDataSourceRequest)
    - [GetDiagnosisDataSourceResponse](#diagnosis.GetDiagnosisDataSourceResponse)
    - [GetJiraSettingRequest](#diagnosis.GetJiraSettingRequest)
    - [GetJiraSettingResponse](#diagnosis.GetJiraSettingResponse)
    - [ListDiagnosisDataSourceRequest](#diagnosis.ListDiagnosisDataSourceRequest)
    - [ListDiagnosisDataSourceResponse](#diagnosis.ListDiagnosisDataSourceResponse)
    - [ListJiraSettingRequest](#diagnosis.ListJiraSettingRequest)
    - [ListJiraSettingResponse](#diagnosis.ListJiraSettingResponse)
    - [PutDiagnosisDataSourceRequest](#diagnosis.PutDiagnosisDataSourceRequest)
    - [PutDiagnosisDataSourceResponse](#diagnosis.PutDiagnosisDataSourceResponse)
    - [PutJiraSettingRequest](#diagnosis.PutJiraSettingRequest)
    - [PutJiraSettingResponse](#diagnosis.PutJiraSettingResponse)
    - [StartDiagnosisRequest](#diagnosis.StartDiagnosisRequest)
    - [StartDiagnosisResponse](#diagnosis.StartDiagnosisResponse)
  
    - [DiagnosisService](#diagnosis.DiagnosisService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="diagnosis/entities.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## diagnosis/entities.proto



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
| created_at | [int64](#int64) |  |  |





 


<a name="diagnosis.Status"></a>

### Status
Status

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| OK | 1 |  |
| CONFIGURED | 2 |  |
| NOT_CONFIGURED | 3 |  |
| ERROR | 4 |  |


 

 

 



<a name="diagnosis/services.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## diagnosis/services.proto



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
| jira_setting_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListJiraSettingResponse"></a>

### ListJiraSettingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jira_setting | [JiraSetting](#diagnosis.JiraSetting) | repeated |  |






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






<a name="diagnosis.StartDiagnosisRequest"></a>

### StartDiagnosisRequest
KICK Diagnosis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| jira_setting_id | [uint32](#uint32) |  |  |






<a name="diagnosis.StartDiagnosisResponse"></a>

### StartDiagnosisResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [string](#string) |  |  |





 

 

 


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
| StartDiagnosis | [StartDiagnosisRequest](#diagnosis.StartDiagnosisRequest) | [StartDiagnosisResponse](#diagnosis.StartDiagnosisResponse) | KICK |

 



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

