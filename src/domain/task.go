package domain

import (
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
)

type TaskInstance struct {
	Id uuid.UUID `json:id`
	ProjectId uuid.UUID `json:projectId`
	Name string `json:name`
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
	files, _ := ioutil.ReadDir(target.Directory())

	type Result struct {
		count int
	}

	return Result{
		count: len(files),
	}
}
