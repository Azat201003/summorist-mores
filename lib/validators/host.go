package validators

import (
	"errors"
)

var emptyHostError = errors.New("Host is empty")

func ValidateHost(host string) (err error) {
	if len(host) == 0 {
		return emptyHostError
	}
	return nil
}
