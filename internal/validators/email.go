package validators

import "net/mail"

func ValidateEmail(data string) (string, error) {
	email, err := mail.ParseAddress(data)
	if err != nil {
		return "", err
	}

	return email.Address, nil
}
