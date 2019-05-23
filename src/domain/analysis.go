package domain

import uuid "github.com/satori/go.uuid"

type Analysis struct {
	Id    uuid.UUID `json:"id"`
	JobId uuid.UUID `json:"jobId"`
}
