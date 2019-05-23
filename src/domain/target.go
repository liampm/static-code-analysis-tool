package domain

import (
	"github.com/satori/go.uuid"
)

type Target struct {
	Id        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	Name      string    `json:"name"`
}
