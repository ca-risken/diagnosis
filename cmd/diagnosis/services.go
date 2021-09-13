package main

import (
	"github.com/ca-risken/diagnosis/proto/diagnosis"
)

type diagnosisService struct {
	repository diagnosisRepoInterface
	sqs        sqsAPI
}

func newDiagnosisService(db diagnosisRepoInterface, s sqsAPI) diagnosis.DiagnosisServiceServer {
	return &diagnosisService{
		repository: db,
		sqs:        s}
}
