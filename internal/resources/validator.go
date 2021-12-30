package resources

import "errors"

func ValidateUser(username, password string) (bool, error) {
	if username == "" || password == "" {
		return false, errors.New("username and password should not be empty")
	}

	if len(username) < 8 || len(password) < 8 {
		return false, errors.New("username and password length must be greater than 8")
	}
	return true, nil
}
