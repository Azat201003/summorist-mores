package validators

import "strconv"

func ValidatePort(port string) (err error) {
	_, err = strconv.ParseUint(port, 10, 16)
	return
}
