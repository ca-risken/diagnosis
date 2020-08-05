#!/bin/bash -e

# github
export GITHUB_USER="iiiidaaa"
export GITHUB_TOKEN="069ac4aa7c3c32522f447506622f80d81a790cf2"

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
export ENDPOINT="http://sqs:9324"
export DIAGNOSIS_JIRA_QUEUE_NAME="diagnosis-jira"
export DIAGNOSIS_JIRA_QUEUE_URL="http://sqs:9324/queue/diagnosis-jira"
