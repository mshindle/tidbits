package factory

import (
	"github.com/apex/log"
)

// Notifier defines the behavior for sending notifications
type Notifier interface {
	Send(recipient string, message string) error
}

const (
	ChannelEmail = "email"
	ChannelSMS   = "sms"
	ChannelPush  = "push"
)

// NewNotifier is factory function that returns a Notifier
func NewNotifier(channel string, config map[string]string) Notifier {
	switch channel {
	case ChannelEmail:
		return newEmailNotifier(config)
	case ChannelSMS:
		return newSMSNotifier(config)
	case ChannelPush:
		return newPushNotifier(config)
	default:
		// Default to email if channel is unknown

		log.WithField("channel", channel).Warn("unknown channel; defaulting to email")
		return newEmailNotifier(config)
	}
}
