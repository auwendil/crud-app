package main

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func handleSuccessfulJSON(w http.ResponseWriter, msg string, payload any, statusCode int, headers ...http.Header) error {
	return writeJSON(w, false, msg, payload, statusCode, headers...)
}

func handleErrorJSON(w http.ResponseWriter, err error, statusCode int, headers ...http.Header) error {
	return writeJSON(w, true, err.Error(), nil, statusCode, headers...)
}

func writeJSON(w http.ResponseWriter, hasError bool, msg string, payload interface{}, statusCode int, headers ...http.Header) error {
	out, err := json.Marshal(JSONResponse{Error: hasError, Message: msg, Data: payload})
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}
