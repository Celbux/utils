package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (u UtilsService) Err(w http.ResponseWriter, err error) bool {
	if err != nil {
		return false
	}

	httpError, ok := err.(HTTPError)
	if ok {
		w.WriteHeader(httpError.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = u.respond(w, err)
	return false
}

func (u UtilsService) respond(w http.ResponseWriter, obj interface{}) error {
	// Writes the encoded marshalled json into the http writer mainly for the purpose of a response
	(w).Header().Set("Content-Type", "application/json; charset=utf-8")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
	    u.NewError(errors.New(fmt.Sprintf(": %v", err), ), "", http.StatusBadRequest)
		return fmt.Errorf(": %v", err)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u UtilsService) NewError(err error, detail string, status int) error {
	return &HTTPError{
		Cause:      err,
		Detail:     detail,
		StatusCode: status,
	}
}
