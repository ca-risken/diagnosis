# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [entities.proto](#entities.proto)
    - [Diagnosis](#diagnosis.Diagnosis)
    - [DiagnosisDataSource](#diagnosis.DiagnosisDataSource)
    - [DiagnosisDataSourceForUpsert](#diagnosis.DiagnosisDataSourceForUpsert)
    - [DiagnosisForUpsert](#diagnosis.DiagnosisForUpsert)
    - [RelDiagnosisDataSource](#diagnosis.RelDiagnosisDataSource)
    - [RelDiagnosisDataSourceForUpsert](#diagnosis.RelDiagnosisDataSourceForUpsert)
  
- [services.proto](#services.proto)
    - [DeleteDiagnosisDataSourceRequest](#diagnosis.DeleteDiagnosisDataSourceRequest)
    - [DeleteDiagnosisRequest](#diagnosis.DeleteDiagnosisRequest)
    - [DeleteRelDiagnosisDataSourceRequest](#diagnosis.DeleteRelDiagnosisDataSourceRequest)
    - [GetDiagnosisDataSourceRequest](#diagnosis.GetDiagnosisDataSourceRequest)
    - [GetDiagnosisDataSourceResponse](#diagnosis.GetDiagnosisDataSourceResponse)
    - [GetDiagnosisRequest](#diagnosis.GetDiagnosisRequest)
    - [GetDiagnosisResponse](#diagnosis.GetDiagnosisResponse)
    - [GetRelDiagnosisDataSourceRequest](#diagnosis.GetRelDiagnosisDataSourceRequest)
    - [GetRelDiagnosisDataSourceResponse](#diagnosis.GetRelDiagnosisDataSourceResponse)
    - [ListDiagnosisDataSourceRequest](#diagnosis.ListDiagnosisDataSourceRequest)
    - [ListDiagnosisDataSourceResponse](#diagnosis.ListDiagnosisDataSourceResponse)
    - [ListDiagnosisRequest](#diagnosis.ListDiagnosisRequest)
    - [ListDiagnosisResponse](#diagnosis.ListDiagnosisResponse)
    - [ListRelDiagnosisDataSourceRequest](#diagnosis.ListRelDiagnosisDataSourceRequest)
    - [ListRelDiagnosisDataSourceResponse](#diagnosis.ListRelDiagnosisDataSourceResponse)
    - [PutDiagnosisDataSourceRequest](#diagnosis.PutDiagnosisDataSourceRequest)
    - [PutDiagnosisDataSourceResponse](#diagnosis.PutDiagnosisDataSourceResponse)
    - [PutDiagnosisRequest](#diagnosis.PutDiagnosisRequest)
    - [PutDiagnosisResponse](#diagnosis.PutDiagnosisResponse)
    - [PutRelDiagnosisDataSourceRequest](#diagnosis.PutRelDiagnosisDataSourceRequest)
    - [PutRelDiagnosisDataSourceResponse](#diagnosis.PutRelDiagnosisDataSourceResponse)
    - [StartDiagnosisRequest](#diagnosis.StartDiagnosisRequest)
    - [StartDiagnosisResponse](#diagnosis.StartDiagnosisResponse)
  
    - [DiagnosisService](#diagnosis.DiagnosisService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="entities.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## entities.proto



<a name="diagnosis.Diagnosis"></a>

### Diagnosis



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






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






<a name="diagnosis.DiagnosisForUpsert"></a>

### DiagnosisForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="diagnosis.RelDiagnosisDataSource"></a>

### RelDiagnosisDataSource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_diagnosis_data_source_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| diagnosis_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| record_id | [string](#string) |  |  |
| jira_id | [string](#string) |  |  |
| jira_key | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="diagnosis.RelDiagnosisDataSourceForUpsert"></a>

### RelDiagnosisDataSourceForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_diagnosis_data_source_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |
| diagnosis_id | [uint32](#uint32) |  |  |
| record_id | [string](#string) |  |  |
| jira_id | [string](#string) |  |  |
| jira_key | [string](#string) |  |  |





 

 

 

 



<a name="services.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services.proto



<a name="diagnosis.DeleteDiagnosisDataSourceRequest"></a>

### DeleteDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeleteDiagnosisRequest"></a>

### DeleteDiagnosisRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_id | [uint32](#uint32) |  |  |






<a name="diagnosis.DeleteRelDiagnosisDataSourceRequest"></a>

### DeleteRelDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_diagnosis_data_source_id | [uint32](#uint32) |  |  |






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






<a name="diagnosis.GetDiagnosisRequest"></a>

### GetDiagnosisRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetDiagnosisResponse"></a>

### GetDiagnosisResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis | [Diagnosis](#diagnosis.Diagnosis) |  |  |






<a name="diagnosis.GetRelDiagnosisDataSourceRequest"></a>

### GetRelDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.GetRelDiagnosisDataSourceResponse"></a>

### GetRelDiagnosisDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_diagnosis_data_source | [RelDiagnosisDataSource](#diagnosis.RelDiagnosisDataSource) |  |  |






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






<a name="diagnosis.ListDiagnosisRequest"></a>

### ListDiagnosisRequest
Diagnosis Service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="diagnosis.ListDiagnosisResponse"></a>

### ListDiagnosisResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis | [Diagnosis](#diagnosis.Diagnosis) | repeated |  |






<a name="diagnosis.ListRelDiagnosisDataSourceRequest"></a>

### ListRelDiagnosisDataSourceRequest
RelDiagnosisDataSourceService


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis_id | [uint32](#uint32) |  |  |
| diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.ListRelDiagnosisDataSourceResponse"></a>

### ListRelDiagnosisDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_diagnosis_data_source | [RelDiagnosisDataSource](#diagnosis.RelDiagnosisDataSource) | repeated |  |






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






<a name="diagnosis.PutDiagnosisRequest"></a>

### PutDiagnosisRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| diagnosis | [DiagnosisForUpsert](#diagnosis.DiagnosisForUpsert) |  |  |






<a name="diagnosis.PutDiagnosisResponse"></a>

### PutDiagnosisResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diagnosis | [Diagnosis](#diagnosis.Diagnosis) |  |  |






<a name="diagnosis.PutRelDiagnosisDataSourceRequest"></a>

### PutRelDiagnosisDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_diagnosis_data_source | [RelDiagnosisDataSourceForUpsert](#diagnosis.RelDiagnosisDataSourceForUpsert) |  |  |






<a name="diagnosis.PutRelDiagnosisDataSourceResponse"></a>

### PutRelDiagnosisDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_diagnosis_data_source | [RelDiagnosisDataSource](#diagnosis.RelDiagnosisDataSource) |  |  |






<a name="diagnosis.StartDiagnosisRequest"></a>

### StartDiagnosisRequest
KICK Diagnosis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_diagnosis_data_source_id | [uint32](#uint32) |  |  |






<a name="diagnosis.StartDiagnosisResponse"></a>

### StartDiagnosisResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [string](#string) |  |  |





 

 

 


<a name="diagnosis.DiagnosisService"></a>

### DiagnosisService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListDiagnosis | [ListDiagnosisRequest](#diagnosis.ListDiagnosisRequest) | [ListDiagnosisResponse](#diagnosis.ListDiagnosisResponse) | Diagnosis |
| GetDiagnosis | [GetDiagnosisRequest](#diagnosis.GetDiagnosisRequest) | [GetDiagnosisResponse](#diagnosis.GetDiagnosisResponse) |  |
| PutDiagnosis | [PutDiagnosisRequest](#diagnosis.PutDiagnosisRequest) | [PutDiagnosisResponse](#diagnosis.PutDiagnosisResponse) |  |
| DeleteDiagnosis | [DeleteDiagnosisRequest](#diagnosis.DeleteDiagnosisRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListDiagnosisDataSource | [ListDiagnosisDataSourceRequest](#diagnosis.ListDiagnosisDataSourceRequest) | [ListDiagnosisDataSourceResponse](#diagnosis.ListDiagnosisDataSourceResponse) | DataSource |
| GetDiagnosisDataSource | [GetDiagnosisDataSourceRequest](#diagnosis.GetDiagnosisDataSourceRequest) | [GetDiagnosisDataSourceResponse](#diagnosis.GetDiagnosisDataSourceResponse) |  |
| PutDiagnosisDataSource | [PutDiagnosisDataSourceRequest](#diagnosis.PutDiagnosisDataSourceRequest) | [PutDiagnosisDataSourceResponse](#diagnosis.PutDiagnosisDataSourceResponse) |  |
| DeleteDiagnosisDataSource | [DeleteDiagnosisDataSourceRequest](#diagnosis.DeleteDiagnosisDataSourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListRelDiagnosisDataSource | [ListRelDiagnosisDataSourceRequest](#diagnosis.ListRelDiagnosisDataSourceRequest) | [ListRelDiagnosisDataSourceResponse](#diagnosis.ListRelDiagnosisDataSourceResponse) | RelDiagnosisDataSource |
| GetRelDiagnosisDataSource | [GetRelDiagnosisDataSourceRequest](#diagnosis.GetRelDiagnosisDataSourceRequest) | [GetRelDiagnosisDataSourceResponse](#diagnosis.GetRelDiagnosisDataSourceResponse) |  |
| PutRelDiagnosisDataSource | [PutRelDiagnosisDataSourceRequest](#diagnosis.PutRelDiagnosisDataSourceRequest) | [PutRelDiagnosisDataSourceResponse](#diagnosis.PutRelDiagnosisDataSourceResponse) |  |
| DeleteRelDiagnosisDataSource | [DeleteRelDiagnosisDataSourceRequest](#diagnosis.DeleteRelDiagnosisDataSourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
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

