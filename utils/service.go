package helpers

import "net/http"

func (s Utils) FinalError(w http.ResponseWriter, err error) error {
	err := s.RegRepo.RegisterMerchant(request)
	if err != nil {
		return ErrRegisterMerchant
	}
	return nil
}
