package errorh

import (
	"errors"
	"regexp"
)

func validateEmail(email string) error {
	// Define the stricter email regex pattern (based on RFC 5322)
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("email invalid")
	}
	return nil
}
