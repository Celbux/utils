package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Global Spy
var Spy = make(map[string]interface{})

// Tests
func TestUtils_FinalError_Successful(t *testing.T) {
	t.Run("successful RegisterMerchant", func(t *testing.T) {
		Utils := UtilsService{}
		var w httptest.ResponseRecorder
		response := Response
		Utils.Err(w, err)

		assertStatus(t, res.Code, http.StatusOK)
		assertError(t, &res, nil)

		if Spy["SuccessRegisterMerchantRepo-DecodeNewEmail-CallCount"] != 1 {
			t.Errorf("expected call count of %v but got %v", 1, Spy["SuccessRegisterMerchantRepo-DecodeNewEmail-CallCount"])
		}

		if Spy["SuccessRegisterMerchantRepo-RegisterMerchant-CallCount"] != 1 {
			t.Errorf("expected call count of %v but got %v", 1, Spy["SuccessRegisterMerchantRepo-RegisterMerchant-CallCount"])
		}
	})
}

// Helpers
func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Got %v but wanted nil", err)
	}
}

func assertError(t *testing.T, w *httptest.ResponseRecorder, expected error) {
	var res *Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assertNil(t, err)

	if expected == nil {
		expected = errors.New("")
	}

	if res.Error != expected.Error() {
		t.Errorf("expected error %v but got %v", res.Error, expected)
	}
}

func assertStatus(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("expected http status code of %v but got %v", expected, actual)
	}
}

func setHeaders(r *http.Request)  {
	r.Header.Set("Type", fieldsFixture.Type)
	r.Header.Set("FileID", fieldsFixture.FileID)
	r.Header.Set("EmailMain", fieldsFixture.EmailMain)
	r.Header.Set("NameMain", fieldsFixture.NameMain)
	r.Header.Set("IdNumberMain", fieldsFixture.IdNumberMain)
	r.Header.Set("EmailAdmin", fieldsFixture.EmailAdmin)
	r.Header.Set("NameAdmin", fieldsFixture.NameAdmin)
	r.Header.Set("ContactAdmin", fieldsFixture.ContactAdmin)
	r.Header.Set("EmailAlt", fieldsFixture.EmailAlt)
	r.Header.Set("NameAlt", fieldsFixture.NameAlt)
	r.Header.Set("ContactAlt", fieldsFixture.ContactAlt)
	r.Header.Set("AddressLine1", fieldsFixture.AddressLine1)
	r.Header.Set("AddressLine2", fieldsFixture.AddressLine2)
	r.Header.Set("AddressLine3", fieldsFixture.AddressLine3)
	r.Header.Set("AddressLine4", fieldsFixture.AddressLine4)
}

func setupRegisterMerchantRequest(t *testing.T) (w httptest.ResponseRecorder, r *http.Request) {
	dts, err := json.Marshal(&fileFixture)
	assertNil(t, err)

	req, err := http.NewRequest("POST", "/RegisterMerchant", bytes.NewBuffer(dts))
	assertNil(t, err)
	setHeaders(req)
	return *httptest.NewRecorder(), req
}

// Mocks, Stubs & Fixtures
type SuccessRegisterMerchantRepo struct {}

func (s SuccessRegisterMerchantRepo) DecodeNewEmail(r http.Request) (NewRegisterMerchantRequest, error) {
	Spy["SuccessRegisterMerchantRepo-DecodeNewEmail-CallCount"] = 1
	return NewRegisterMerchantRequest{}, nil
}

func (s SuccessRegisterMerchantRepo) RegisterMerchant(request NewRegisterMerchantRequest) error {
	Spy["SuccessRegisterMerchantRepo-RegisterMerchant-CallCount"] = 1
	return nil
}

type FailureDecodeNewEmailRepo struct {}

func (s *FailureDecodeNewEmailRepo) DecodeNewEmail(r http.Request) (NewRegisterMerchantRequest, error) {
	Spy["FailureDecodeNewEmailRepo-DecodeNewEmail-CallCount"] = 1
	return NewRegisterMerchantRequest{}, errors.New("some error")
}

func (s *FailureDecodeNewEmailRepo) RegisterMerchant(request NewRegisterMerchantRequest) error {
	return nil
}

type FailureRegisterMerchantRepo struct {}

func (s *FailureRegisterMerchantRepo) DecodeNewEmail(r http.Request) (NewRegisterMerchantRequest, error) {
	Spy["FailureRegisterMerchantRepo-DecodeNewEmail-CallCount"] = 1
	return NewRegisterMerchantRequest{}, nil
}

func (s *FailureRegisterMerchantRepo) RegisterMerchant(request NewRegisterMerchantRequest) error {
	Spy["FailureRegisterMerchantRepo-RegisterMerchant-CallCount"] = 1
	return errors.New("some error")
}

type SuccessPayRepo struct {}

func (s SuccessPayRepo) Pay(request NewPayRequest) error {
	Spy["SuccessPayRepo-Pay-CallCount"] = 1
	return nil
}
