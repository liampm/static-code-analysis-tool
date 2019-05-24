package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type TaskController struct {
	ReadRepo  domain.TaskReadRepo
	WriteRepo domain.TaskRepo
}

func (controller *TaskController) All() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "projectId")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		err = marshalJSONResponse(w, controller.ReadRepo.AllForProject(projectUuid))
		if err != nil {
			panic(err) // Panic whilst in development
		}

		w.WriteHeader(200)
	}
}

func (controller *TaskController) ById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		uuid, err := uuidFromParam(r, "id")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		task, err := controller.ReadRepo.Find(uuid)
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}

		err = marshalJSONResponse(w, task)
		if err != nil {
			panic(err) // Panic whilst in development
		}

		w.WriteHeader(200)
	}
}

func (controller *TaskController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "projectId")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		task := domain.TaskInstance{}
		decoder := json.NewDecoder(r.Body)

		err = decoder.Decode(&task)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		task.Id = uuid.NewV4()
		task.ProjectId = projectUuid

		controller.WriteRepo.Save(task)

		w.WriteHeader(201)
	}
}
