package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelop map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelop) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil

}

func ReadIdParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid ID parameter")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid ID parameter: " + err.Error())
	}
	return id, nil
}
