package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type TargetController struct {
	ReadRepo  domain.TargetReadRepo
	WriteRepo domain.TargetRepo
}

func (controller *TargetController) AllForProject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rawProjectId := chi.URLParam(r, "projectId")
		w.Header().Set("Content-Type", "application/json")

		if rawProjectId == "" {
			w.WriteHeader(404) // Not found for somehow missing ID
			return
		}

		projectId, err := uuid.FromString(rawProjectId)

		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		targets, err := controller.ReadRepo.AllForProject(projectId)

		if err != nil {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)
		jsonBody, err := json.Marshal(targets)


		if err != nil {
			panic(err) // Panic whilst in development
		}

		_, err = w.Write(jsonBody)

		if err != nil {
			panic(err) // Panic whilst in development
		}
	}
}

func (controller *TargetController) ById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		w.Header().Set("Content-Type", "application/json")

		if id == "" {
			w.WriteHeader(404) // Not found for somehow missing ID
			return
		}

		targetId, err := uuid.FromString(id)

		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		target, err := controller.ReadRepo.Find(targetId)

		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}

		jsonBody, err := json.Marshal(target)

		if err != nil {
			panic(err) // Panic whilst in development
		}

		w.WriteHeader(200)
		_, err = w.Write(jsonBody)

		if err != nil {
			panic(err) // Panic whilst in development
		}
	}
}

func (controller *TargetController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		target := domain.Target{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&target)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		target.Id = uuid.NewV4()

		controller.WriteRepo.Save(target)
		w.WriteHeader(201)
	}
}
