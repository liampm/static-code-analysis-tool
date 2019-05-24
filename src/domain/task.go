package domain

import (
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"path/filepath"
)

type TaskInstance struct {
	Id        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	Name      string    `json:"name"`
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

	directory, err := target.directory()

	if err != nil {
		panic(err)
	}

	_ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			count++
		}
		return nil
	})

	err = os.RemoveAll(directory)

	if err != nil {
		log.Printf("Failed to clear up output directory '%s'", directory)
	}

	type Result struct {
		Count int `json:"count"`
	}

	return Result{
		Count: count,
	}
}
