package web

import (
	"encoding/json"
	"net/http"
)

type JSONResponseMessage struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Errors  string      `json:"errors"`
}

func JSONResponse(response interface{}, w http.ResponseWriter) error {
	m := JSONResponseMessage{
		Success: true,
		Message: response,
	}

	j, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)

	return nil
}

func JSONResponseError(err error, w http.ResponseWriter) error {
	http.Error(w, err.Error(), http.StatusInternalServerError)

	m := JSONResponseMessage{
		Success: false,
		Errors:  err.Error(),
	}

	j, err := json.Marshal(m)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)

	return nil
}
