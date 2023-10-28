package validators

import (
	"errors"
	"unicode"
)

func ValidateCode(value string) error {
	if len(value) != 6 {
		return errors.New("code must be 6 digits")
	}
	for _, char := range value {
		if !unicode.IsDigit(char) {
			return errors.New("code must be digits")
		}
	}
	return nil
}
