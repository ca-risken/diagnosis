#!/bin/bash -e

# github
export GITHUB_USER="your-name"
export GITHUB_TOKEN="your_token"

# GO
export GOPRIVATE="github.com/CyberAgent/*"

# DB
export DB_MASTER_HOST="db"
export DB_MASTER_USER="hoge"
export DB_MASTER_PASSWORD="moge"
export DB_SLAVE_HOST="db"
export DB_SLAVE_USER="hoge"
export DB_SLAVE_PASSWORD="moge"
export DB_LOG_MODE="false"

# AWS
export AWS_REGION="ap-northeast-1"
export AWS_ACCESS_KEY_ID="testkey"
export AWS_SECRET_ACCESS_KEY="testsecretkey"
export SQS_ENDPOINT="http://sqs:9324"
export DIAGNOSIS_JIRA_QUEUE_NAME="diagnosis-jira"
export DIAGNOSIS_JIRA_QUEUE_URL="http://sqs:9324/queue/diagnosis-jira"
export DIAGNOSIS_JIRA_URL="https://ca-security.atlassian.net/"
export DIAGNOSIS_JIRA_USER_ID="jira-id"
export DIAGNOSIS_JIRA_USER_PASSWORD="jira-password"

# mimosa
export FINDING_SVC_ADDR="finding:8001"
export ALERT_SVC_ADDR="alert:8004"
export DIAGNOSIS_SVC_ADDR="diagnosis:19001"