package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func uuidFromParam(r *http.Request, paramName string) (uuid.UUID, error) {
	paramValue := chi.URLParam(r, paramName)

	if paramValue == "" {
		return uuid.UUID{}, errors.New("no param")
	}

	uuidValue, err := uuid.FromString(paramValue)

	if err != nil {
		return uuid.UUID{}, err
	}

	return uuidValue, nil
}

func marshalJSONResponse(w http.ResponseWriter, body interface{}) error {
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return err
	}

	_, err = w.Write(jsonBody)

	if err != nil {
		return err
	}

	return nil
}
