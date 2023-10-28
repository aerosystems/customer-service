package validators

import (
	"errors"
	"regexp"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password should be of 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password should not be of more than 128 characters long")
	}
	done, err := regexp.MatchString("([a-z])+", password)
	if err != nil {
		return err
	}
	if !done {
		return errors.New("password should contain at least one lower case character")
	}
	done, err = regexp.MatchString("([A-Z])+", password)
	if err != nil {
		return err
	}
	if !done {
		return errors.New("password should contain at least one upper case character")
	}
	done, err = regexp.MatchString("([0-9])+", password)
	if err != nil {
		return err
	}
	if !done {
		return errors.New("password should contain at least one digit")
	}
	done, err = regexp.MatchString("([!@#$%^&*.?-])+", password)
	if err != nil {
		return err
	}
	if !done {
		return errors.New("password should contain at least one special character")
	}
	return nil
}
