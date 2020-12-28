package utils

import "net/http"

type Server struct {
	Service BlueprintRepo
}

type IRepo interface {
	NewErr(details string, status int) HTTPResponse
	Err(w http.ResponseWriter, err HTTPResponse) bool

	respond(w http.ResponseWriter, obj interface{}) error
}

type BlueprintRepo struct {
	Repo IRepo
}

type Repo struct {}

type HTTPResponse struct {
	Details    string
	StatusCode int
}