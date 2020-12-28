package helpers

import "net/http"

type IUtilsRepo interface {
	FinalErr(w http.ResponseWriter, err error) error
	encode(w http.ResponseWriter, obj interface{}) error
}

type UtilsService struct {
	Error 	HTTPError
}

type HTTPError struct {
	Cause      error
	Detail     string
	StatusCode int
}

func (e HTTPError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}
