package domain

import (
	"github.com/satori/go.uuid"
	git "gopkg.in/src-d/go-git.v4"
)

type Target struct {
	Id        uuid.UUID    `json:"id"`
	ProjectId uuid.UUID    `json:"projectId"`
	Name      string       `json:"name"`
	Config    TargetConfig `json:"config"`
}

type targetType int

const (
	GIT_REPO targetType = iota
)

type TargetConfig interface {
	Type() targetType
}

type RepoConfig struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (config *RepoConfig) Type() targetType {
	return GIT_REPO
}

func (target *Target) directory() string {
	
	return "/go/"
}
