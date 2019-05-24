package domain

import (
	"github.com/satori/go.uuid"
	"path/filepath"
)

type Target struct {
	Id        uuid.UUID    `json:"id"`
	ProjectId uuid.UUID    `json:"projectId"`
	Name      string       `json:"name"`
	Config    TargetConfiguration `json:"config"`
}

type TargetType int

const (
	GIT_REPO TargetType = iota + 1
)


type TargetConfiguration struct {
	Type TargetType `json:"type"`
	Details interface{} `json:"details"`
}

type RepoDetails struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (target *Target) directory(added string) string {
	return filepath.Join("/tmp", target.Id.String(), added)
}
