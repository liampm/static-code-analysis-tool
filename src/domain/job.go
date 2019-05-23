package domain

import uuid "github.com/satori/go.uuid"

type Job struct {
	Id        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}
