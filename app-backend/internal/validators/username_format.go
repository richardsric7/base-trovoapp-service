package validators

import (
	"regexp"
	"trovo-wallet-api/internal/errors"
)

// ValidateUsernameFormat validates a username
func ValidateUsernameFormat(username string) error {
	var err errors.ErrorInvalidUsernameFormat

	if len(username) < 3 {
		err.Detail = "Username cannot be less than 3 characters"
		return &err
	}

	if len(username) > 25 {
		err.Detail = "Username cannot be more than 25 characters"
		return &err
	}

	//check if username contains invalid characters

	matched, _ := regexp.MatchString("^[a-z][a-z0-9]+$", username)

	if !matched {
		err.Detail = "Usernames can only start with a letter [a-z], followed by either digits [0-9] or letters"
		return &err
	}

	return nil

}
