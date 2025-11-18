package factory

import "github.com/apex/log"

// PushNotifier implements the Notifier interface
type PushNotifier struct {
	accountID string
	appID     string
}

// Private constructor for EmailNotifier
func newPushNotifier(config map[string]string) *PushNotifier {
	return &PushNotifier{
		accountID: config["accountID"],
		appID:     config["appID"],
	}
}

func (e *PushNotifier) Send(recipient string, message string) error {
	// Implementation details for sending email
	// ...
	log.WithField("recipient", recipient).Info("sending email")
	return nil
}
