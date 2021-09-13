#!/bin/bash -e

# github
export GITHUB_USER="your-name"
export GITHUB_TOKEN="your-token"

# GO
export GOPRIVATE="github.com/CyberAgent/*,github.com/ca-risken/*"

# DB
export DB_MASTER_HOST="db"
export DB_MASTER_USER="hoge"
export DB_MASTER_PASSWORD="moge"
export DB_SLAVE_HOST="db"
export DB_SLAVE_USER="hoge"
export DB_SLAVE_PASSWORD="moge"
export DB_LOG_MODE="false"
export DB_SCHEMA="mimosa"
export DB_PORT="3306"

# AWS
export AWS_REGION="ap-northeast-1"
export AWS_ACCESS_KEY_ID="testkey"
export AWS_SECRET_ACCESS_KEY="testsecretkey"
export SQS_ENDPOINT="http://sqs:9324"
export DIAGNOSIS_JIRA_QUEUE_NAME="diagnosis-jira"
export DIAGNOSIS_JIRA_QUEUE_URL="http://sqs:9324/queue/diagnosis-jira"
export DIAGNOSIS_WPSCAN_QUEUE_NAME="diagnosis-wpscan"
export DIAGNOSIS_WPSCAN_QUEUE_URL="http://sqs:9324/queue/diagnosis-wpscan"
export DIAGNOSIS_PORTSCAN_QUEUE_NAME="diagnosis-portscan"
export DIAGNOSIS_PORTSCAN_QUEUE_URL="http://sqs:9324/queue/diagnosis-portscan"
export DIAGNOSIS_APPLICATION_SCAN_QUEUE_BANE="diagnosis-applicationscan"
export DIAGNOSIS_APPLICATION_SCAN_QUEUE_URL="http://sqs:9324/queue/diagnosis-applicationscan"
export DIAGNOSIS_JIRA_URL="https://ca-security.atlassian.net/"
export DIAGNOSIS_JIRA_USER_ID="jira-user"
export DIAGNOSIS_JIRA_USER_PASSWORD="jira-password"

# mimosa
export FINDING_SVC_ADDR="finding:8001"
export ALERT_SVC_ADDR="alert:8004"
export DIAGNOSIS_SVC_ADDR="diagnosis:19001"

# WPScan
export RESULT_PATH="/tmp"
export WPSCAN_VULNDB_APIKEY="wpscan-api"

# Portscan
export MAX_NUMBER_OF_MESSAGE="10"

# Application Scan
export ZAP_PORT="8080"
export ZAP_API_KEY_NAME="apikey"
export ZAP_API_KEY_HEADER="X-ZAP-API-Key"