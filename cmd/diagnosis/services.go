package main

import (
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
)

type diagnosisService struct {
	repository    diagnosisRepoInterface
	sqs           sqsAPI
	projectClient project.ProjectServiceClient
}

func newDiagnosisService(conf *diagnosisConfig) diagnosis.DiagnosisServiceServer {
	return &diagnosisService{
		repository:    conf.DB,
		sqs:           conf.SQS,
		projectClient: conf.projectClient,
	}
}
