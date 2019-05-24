package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
)

type TargetController struct {
	ReadRepo  domain.TargetReadRepo
	WriteRepo domain.TargetRepo
}

func (controller *TargetController) AllForProject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "projectId")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		targets, err := controller.ReadRepo.AllForProject(projectUuid)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)

		err = marshalJSONResponse(w, targets)
		if err != nil {
			panic(err) // Panic whilst in development
		}

	}
}

func (controller *TargetController) ById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		targetId, err := uuidFromParam(r, "id")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}

		target, err := controller.ReadRepo.Find(targetId)
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}

		err = marshalJSONResponse(w, target)
		if err != nil {
			panic(err) // Panic whilst in development
		}

		w.WriteHeader(200)
	}
}

func (controller *TargetController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		projectUuid, err := uuidFromParam(r, "projectId")
		if err != nil {
			w.WriteHeader(404) // Not found for any invalid IDs
			return
		}
		
		type request struct {
			TargetName  string `json:"name"`
			RequestType string `json:"type"`
		}

		processedRequest := request{}

		bodyBytes, _ := ioutil.ReadAll(r.Body)
		
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

		err = decoder.Decode(&processedRequest)
		if err != nil {
			w.WriteHeader(400)
			return
		}
log.Println(processedRequest)
		var config domain.TargetConfig

		if processedRequest.RequestType == "git-repo" {

			type configRequest struct {
				Config domain.RepoConfig `json:"Config"`
			}

			processedConfig := configRequest{}
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

			err = decoder.Decode(&processedConfig)
			if err != nil {
				w.WriteHeader(400)
				writeError(w, "Failed to process Config")
				return
			}
log.Println(&processedConfig.Config)
			if processedConfig.Config.Username == "" {
				w.WriteHeader(http.StatusBadRequest)
				writeError(w, "A user name must be provided")
				return
			} else if processedConfig.Config.Url == "" {
				w.WriteHeader(http.StatusBadRequest)
				writeError(w, "A URL must be provided")
				return
			} else if processedConfig.Config.Token == "" {
				w.WriteHeader(http.StatusBadRequest)
				writeError(w, "A token must be provided")
				return
			}

			config = &processedConfig.Config
		} else {
			w.WriteHeader(400)
			writeError(w, fmt.Sprintf("Unrecognised target type '%s'", processedRequest.RequestType))
			return
		}

		target := domain.Target{
			Id:        uuid.NewV4(),
			ProjectId: projectUuid,
			Name:      processedRequest.TargetName,
			Config:    config,
		}

		controller.WriteRepo.Save(target)

		w.WriteHeader(201)
	}
}

func writeError(w http.ResponseWriter, errorMessage string) {
	error, _ := json.Marshal(map[string]string{"error": errorMessage})
	_, _ = w.Write(error)
}
