package errorh

import (
	"fmt"
	"strings"
)

// ErrorCollection aggregates multiple errors
type ErrorCollection struct {
	Errors []error
}

func (ec *ErrorCollection) Add(err error) {
	if err != nil {
		ec.Errors = append(ec.Errors, err)
	}
}

func (ec *ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return ""
	}

	messages := make([]string, len(ec.Errors))
	for i, err := range ec.Errors {
		messages[i] = err.Error()
	}

	return fmt.Sprintf("%d errors occurred: %s",
		len(ec.Errors), strings.Join(messages, "; "))
}

func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.Errors) > 0
}
