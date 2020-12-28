package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Main
func TestMain(m *testing.M) {
	m.Run()
	os.Exit(0)
}

// Global Spy
var Spy = make(map[string]int)

// Tests
func Test_Err(t *testing.T) {
	t.Run("[Mock] Call Err() once", func(t *testing.T) {
		var w httptest.ResponseRecorder

		service := SuccessRepo{}
		ok := service.Err(&w, fixtureHTTPResponseSuccess)

		assertBool(t, false, ok)
		if Spy["SuccessRepo-Err-CallCount"] != 1 {
			t.Errorf("wanted call count of 1 but got %v", Spy["SuccessRepo-Err-CallCount"])
		}
	})

	t.Run("[Integration] Error is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		server := Server{
			Service: BlueprintRepo{
				Repo: Repo{},
			},
		}

		ok := server.Err(w, fixtureHTTPResponseSuccess)

		assertBool(t, false, ok)
		assertError(t, nil, w)
		assertStatus(t, http.StatusOK, w.Code)
	})

	t.Run("[Integration] Error is not nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		server := Server{
			Service: BlueprintRepo{
				Repo: Repo{},
			},
		}

		ok := server.Err(w, fixtureHTTPResponseErr)

		assertBool(t, true, ok)
		assertError(t, fixtureHTTPResponseErr, w)
		assertStatus(t, http.StatusBadRequest, w.Code)
	})
}

// Helpers
func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Got %v but wanted nil", err)
	}
}

func assertBool(t *testing.T, want bool, got bool) {
	if want != got {
		t.Errorf("wanted bool %v, but got %v", want, got)
	}
}

func assertError(t *testing.T, want error, w *httptest.ResponseRecorder) {
	var res *HTTPResponse

	if w.Body.Len() > 0 {
		err := json.NewDecoder(w.Body).Decode(&res)
		assertNil(t, err)

		if want == nil {
			want = errors.New("")
		}

		if res.Error() != want.Error() {
			t.Errorf("wanted error %v but got %v", want, res.Error())
		}
	} else if want != nil {
		t.Errorf("wanted error %v but w.Body is empty", want)
	}
}

func assertStatus(t *testing.T, want int, got int) {
	if want != got {
		t.Errorf("wanted http status code of %v but got %v", want, got)
	}
}

// Mocks, Stubs & Fixtures
type SuccessRepo struct {}

var fixtureHTTPResponseSuccess = HTTPResponse{
	Details: "",
	StatusCode: http.StatusOK,
}

var fixtureHTTPResponseErr = HTTPResponse{
	Details: "test error",
	StatusCode: http.StatusBadRequest,
}

func (u SuccessRepo) NewErr(details string, status int) HTTPResponse {
	Spy["SuccessRepo-NewErr-CallCount"]++
	return fixtureHTTPResponseSuccess
}

func (u SuccessRepo) Err(w http.ResponseWriter, err HTTPResponse) bool {
	Spy["SuccessRepo-Err-CallCount"]++
	return false
}

func (u SuccessRepo) encode(w http.ResponseWriter, obj interface{}) error {
	Spy["SuccessRepo-encode-CallCount"]++
	return nil
}