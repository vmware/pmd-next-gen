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

func httpResponse(m *JSONResponseMessage, response interface{}, w http.ResponseWriter) error {
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

func JSONResponse(response interface{}, w http.ResponseWriter) error {
	m := JSONResponseMessage{
		Success: true,
		Message: response,
	}

	return httpResponse(&m , response, w)
}

func JSONResponseError(err error, w http.ResponseWriter) error {
	m := JSONResponseMessage{
		Success: false,
		Errors:  err.Error(),
	}

	return httpResponse(&m , nil, w)
}
