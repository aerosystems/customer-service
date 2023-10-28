package validators

import (
	"errors"
	"os"
	"strings"
)

func ValidateRole(role string) error {
	trustRoles := strings.Split(os.Getenv("TRUST_ROLES"), ",")
	if !Contains(trustRoles, role) {
		return errors.New("role exists in trusted roles")
	}
	return nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
