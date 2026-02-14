package utils

import (
	"errors"
	"slices"
)

type ContextKey string

func AzuthorizeUser(userRole string, alowedRoles ...string) (bool, error) {
	hasRole := slices.Contains(alowedRoles, userRole)
	if hasRole {
		return true, nil
	}
	return false, errors.New("user not authorized")
}
