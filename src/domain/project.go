package domain

import (
	"github.com/satori/go.uuid"
)

type Project struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
