package factory

import "github.com/apex/log"

// EmailNotifier implements the Notifier interface
type EmailNotifier struct {
	smtpServer string
	port       string
	username   string
	password   string
}

// Private constructor for EmailNotifier
func newEmailNotifier(config map[string]string) *EmailNotifier {
	return &EmailNotifier{
		smtpServer: config["smtpServer"],
		port:       config["port"],
		username:   config["username"],
		password:   config["password"],
	}
}

func (e *EmailNotifier) Send(recipient string, message string) error {
	// Implementation details for sending email
	// ...
	log.WithField("recipient", recipient).Info("sending email")
	return nil
}
