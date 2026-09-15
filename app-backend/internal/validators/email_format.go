package validators

import (
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//ValidateEmailFormat validates an email
func ValidateEmailFormat(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	var r bool = emailRegex.MatchString(email)

	return r

}
