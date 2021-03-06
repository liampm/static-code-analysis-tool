package controller

import (
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type JobController struct {
	projectRepo domain.ProjectRepo
	readRepo    domain.JobReadRepo
	writeRepo   domain.JobRepo
}

func NewJobController(projectRepo domain.ProjectRepo, readRepo domain.JobReadRepo, writeRepo domain.JobRepo) *JobController {
	return &JobController{
		projectRepo: projectRepo,
		readRepo:    readRepo,
		writeRepo:   writeRepo,
	}
}

func (controller *JobController) Initiate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "projectId")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		project, err := controller.projectRepo.Find(projectUuid)
		if err != nil {
			w.WriteHeader(404) // Not found the project
			return
		}

		job := domain.Job{
			Id:        uuid.NewV4(),
			ProjectId: project.Id,
		}

		go func() {
			controller.writeRepo.Save(project.Analyse(job.Id))
		}()

		controller.writeRepo.Save(job)
		w.WriteHeader(201)
	}
}

func (controller *JobController) ByProjectId() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "projectId")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		project, err := controller.projectRepo.Find(projectUuid)
		if err != nil {
			w.WriteHeader(404) // Not found the project
			return
		}

		err = marshalJSONResponse(w, controller.readRepo.AllForProject(project))
		if err != nil {
			panic(err) // Panic whilst in development
		}
	}
}
