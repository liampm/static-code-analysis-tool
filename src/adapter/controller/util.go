package controller

import (
	"errors"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func uuidFromParam (r *http.Request, paramName string) (uuid.UUID, error) {
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
