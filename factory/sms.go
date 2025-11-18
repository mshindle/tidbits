package factory

import "github.com/apex/log"

// SMSNotifier implements the Notifier interface
type SMSNotifier struct {
	accountID string
}

// Private constructor for EmailNotifier
func newSMSNotifier(config map[string]string) *SMSNotifier {
	return &SMSNotifier{
		accountID: config["accountID"],
	}
}

func (e *SMSNotifier) Send(recipient string, message string) error {
	// Implementation details for sending email
	// ...
	log.WithField("recipient", recipient).Info("sending email")
	return nil
}
