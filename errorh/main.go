package errorh

import (
	"errors"
	"net/mail"

	"github.com/apex/log"
)

type User struct {
	Name  string
	Email string
}

func Execute() {
	users := []User{
		{
			Name:  "John Invalid",
			Email: "invalid",
		},
		{
			Name:  "John Valid",
			Email: "jvalid@example.com",
		},
	}
	for _, user := range users {
		logger := log.WithField("name", user.Name).WithField("email", user.Email)
		logger.Info("processing user")
		err := processUser(user)
		if err != nil {
			logger.WithError(err).Error("error processing user")
			continue
		}
		logger.Info("processed user")
	}
}

// processUser just validates the user. It would normally be a more robust function.
func processUser(user User) error {
	return validateUser(user)

}

func validateUser(user User) error {
	var errs = &ErrorCollection{}

	if len(user.Name) < 2 {
		errs.Add(errors.New("name too short"))
	}

	if len(user.Email) == 0 {
		errs.Add(errors.New("email required"))
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		errs.Add(err)
	}

	if !errs.HasErrors() {
		return nil
	}
	return errs
}
