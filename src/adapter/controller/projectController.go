package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type ProjectController struct {
	ReadRepo  domain.ProjectReadRepo
	WriteRepo domain.ProjectRepo
}

func (controller *ProjectController) All() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		err := marshalJSONResponse(w, controller.ReadRepo.All())
		if err != nil {
			panic(err) // Panic whilst in development
		}
	}
}

func (controller *ProjectController) ById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "id")

		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		project, err := controller.ReadRepo.Find(projectUuid)
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)

		err = marshalJSONResponse(w, project)
		if err != nil {
			panic(err) // Panic whilst in development
		}
	}
}

func (controller *ProjectController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		project := domain.Project{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&project)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		project.Id = uuid.NewV4()

		controller.WriteRepo.Save(project)

		w.WriteHeader(201)

		err = marshalJSONResponse(w, project)
		if err != nil {
			panic(err) // Panic whilst in development
		}
	}
}
