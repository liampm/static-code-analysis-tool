package domain

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"log"
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

func (target *Target) directory() (string, error) {
	details := target.Config.Details
	log.Println(target)
	switch configType := details.(type) {
	case RepoDetails:
		log.Println(configType)
		return getRepoDirectory(details.(RepoDetails), target.Id)
	default:
		return "", errors.New("unknown config type")
	}
}

func getRepoDirectory(config RepoDetails, targetId uuid.UUID) (string, error) {
	outputDir := filepath.Join("/tmp", targetId.String())

	log.Printf("Loading repo into %s", outputDir)
	_, err := git.PlainClone(outputDir, true, &git.CloneOptions{
		URL:config.Url,
		Auth: &http.BasicAuth{
			Username: config.Username,
			Password: config.Token,
		},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return "", err
	}

	return outputDir, nil
}