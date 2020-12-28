package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (u Repo) NewErr(details string, status int) HTTPResponse {
	return HTTPResponse{
		Details:    details,
		StatusCode: status,
	}
}

func (u Repo) Err(w http.ResponseWriter, err HTTPResponse) bool {
	if err.Details == "" {
		return false
	}

	if err.StatusCode == 0 {
		err.StatusCode = http.StatusInternalServerError
	}

	w.WriteHeader(err.StatusCode)

	_ = u.respond(w, err)
	return true
}

// Writes the encoded marshalled json into the http writer mainly for the purpose of a response
func (u Repo) respond(w http.ResponseWriter, obj interface{}) error {
	(w).Header().Set("Content-Type", "application/json; charset=utf-8")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
	    return u.NewErr(fmt.Sprintf(": %v", err), http.StatusBadRequest)
	}

	return nil
}