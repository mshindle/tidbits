package errorh

import "testing"

func Test_validateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid01",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "missing domain",
			email:   "invalid-email@",
			wantErr: true,
		},
		{
			name:    "invalid domain",
			email:   "user@.com",
			wantErr: true,
		},
		{
			name:    "valid02",
			email:   "test@domain-with-hyphen.com",
			wantErr: false,
		},
		{
			name:    "missing tld",
			email:   "user@tld.",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEmail(tt.email); (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err != nil, tt.wantErr)
			}
		})
	}
}
