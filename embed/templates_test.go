package embed

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0 mins"},
		{1 * time.Minute, "1 mins"},
		{30 * time.Minute, "30 mins"},
		{59 * time.Minute, "59 mins"},
		{60 * time.Minute, "60 mins"},
		{61 * time.Minute, "1:00 hours"},
		{67 * time.Minute, "1:00 hours"},
		{68 * time.Minute, "1:15 hours"},
		{75 * time.Minute, "1:15 hours"},
		{1 * time.Hour, "60 mins"},
		{2 * time.Hour, "2:00 hours"},
		{119 * time.Minute, "2:00 hours"},
		{120 * time.Minute, "2:00 hours"},
		{121 * time.Minute, "2:00 hours"},
		{127 * time.Minute, "2:00 hours"},
		{128 * time.Minute, "2:15 hours"},
		{510 * time.Minute, "8:30 hours"},
	}

	for _, tt := range tests {
		t.Run(tt.duration.String(), func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.expected {
				t.Errorf("formatDuration(%v) = %q; want %q", tt.duration, got, tt.expected)
			}
		})
	}
}
