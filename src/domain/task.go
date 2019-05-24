package domain

import (
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
)

type TaskInstance struct {
	Id        uuid.UUID `json:id`
	ProjectId uuid.UUID `json:projectId`
	Name      string    `json:name`
}

type task interface {
	analyse(target *Target) interface{}
	name() string
}

type fileCount struct {
}

func (task *fileCount) name() string {
	return "File count"
}

func (task *fileCount) analyse(target *Target) interface{} {
	count := 0

	_ = filepath.Walk(target.directory(), func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			count++
		}
		return nil
	})

	type Result struct {
		Count int `json:"count"`
	}

	return Result{
		Count: count,
	}
}
