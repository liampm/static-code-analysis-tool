package domain

import (
	"errors"
	"github.com/satori/go.uuid"
	"time"
)

type ProjectReference struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Project struct {
	Id         uuid.UUID `json:"id" gorm:"primary_key"`
	Name       string    `json:"name"`
	Tasks      []TaskInstance
	Targets    []Target
	tasksCache map[string]task
}

func (project *Project) Analyse(jobId uuid.UUID) (job Job) {
	project.tasksCache = make(map[string]task)
	var analyses []Analysis

	time.Sleep(3000 * time.Millisecond) // TODO Remove after testing

	job.Id = jobId
	job.ProjectId = project.Id

	for _, target := range project.Targets {
		for _, task := range project.Tasks {
			analyses = append(analyses, fromTask(&job, &target, project.loadTask(task.Name), task.Id))
		}
	}

	job.Analyses = analyses

	return
}

func (project *Project) loadTask(name string) task {
	var actualTask task
	if actualTask, ok := project.tasksCache[name]; ok {
		return actualTask
	}

	if name == "fileCount" {
		actualTask = &fileCount{}
	} else if name == "push" {
		actualTask = &addAFile{
			filename: "test-it.txt",
			content:  "Does this work?",
		}
	} else {
		panic(errors.New("unknown task"))
	}

	project.tasksCache[name] = actualTask

	return actualTask
}
