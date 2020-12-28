package utils

import "net/http"

func (s Server) NewErr(details string, status int) HTTPResponse {
	return s.Service.Repo.NewErr(details, status)
}

func (s Server) Err(w http.ResponseWriter, err HTTPResponse) bool {
	return s.Service.Repo.Err(w, err)
}

func (r HTTPResponse) Error() string {
	return r.Details
}