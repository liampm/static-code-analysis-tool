package domain

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type ProjectReference struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
}

type Project struct {
	Id         uuid.UUID      `json:"id" gorm:"primary_key"`
	Name       string         `json:"name"`
	Tasks      []TaskInstance
	Targets    []Target
	tasksCache map[string]task
}

func (project *Project) Analyse(jobId uuid.UUID) (job Job) {
	var analyses []Analysis
	log.Printf("Analysing job %s\n", jobId)

	time.Sleep(3000 * time.Millisecond)
	job.Id = jobId
	job.ProjectId = project.Id

	for _, target := range project.Targets {
		log.Printf("Analysing target %s\n", target.Name)
		for _, task := range project.Tasks {
			log.Printf("Performing task %s\n", task.Name)
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
	} else {
		panic(errors.New("unknown task"))
	}

	project.tasksCache[name] = actualTask

	return actualTask
}
