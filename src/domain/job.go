package domain

import uuid "github.com/satori/go.uuid"

type JobReference struct {
	Id   	  uuid.UUID           `json:"id"`
	ProjectId uuid.UUID           `json:"projectId"`
	Analyses  []AnalysisReference `json:"analyses"`
}

type Job struct {
	Id        uuid.UUID  `json:"id"`
	ProjectId uuid.UUID  `json:"projectId"`
	Analyses  []Analysis `json:"analyses"`
}
