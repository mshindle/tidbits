package toy

import "testing"

func TestWhisper(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{1},
			want: 1,
		},
		{
			name: "test100k",
			args: args{100000},
			want: 100000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Whisper(tt.args.n); got != tt.want {
				t.Errorf("Whisper() = %v, want %v", got, tt.want)
			}
		})
	}
}
