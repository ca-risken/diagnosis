package main

import (
	"github.com/ca-risken/core/proto/project"
)

type DiagnosisService struct {
	repository    diagnosisRepoInterface
	sqs           sqsAPI
	projectClient project.ProjectServiceClient
}
