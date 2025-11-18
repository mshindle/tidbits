package factory

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewNotifier(t *testing.T) {
	var smtpConfig = map[string]string{
		"smtpServer": "smtp.example.com",
		"port":       "587",
		"username":   "user",
		"password":   "pass",
	}
	type args struct {
		channel string
		config  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Notifier
	}{
		{
			name: "TestNewEmailNotifier",
			args: args{
				channel: ChannelEmail,
				config:  smtpConfig,
			},
			want: &EmailNotifier{
				smtpServer: smtpConfig["smtpServer"],
				port:       smtpConfig["port"],
				username:   smtpConfig["username"],
				password:   smtpConfig["password"],
			},
		},
		{
			name: "TestNewPushNotifier",
			args: args{
				channel: ChannelPush,
				config: map[string]string{
					"accountID": "abc123",
					"appID":     "frodo",
				},
			},
			want: &PushNotifier{accountID: "abc123", appID: "frodo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotifier(tt.args.channel, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MockNotifier for testing
type MockNotifier struct {
	RecipientReceived string
	MessageReceived   string
	ShouldFail        bool
}

func (m *MockNotifier) Send(recipient string, message string) error {
	m.RecipientReceived = recipient
	m.MessageReceived = message
	if m.ShouldFail {
		return fmt.Errorf("mock error")
	}
	return nil
}
