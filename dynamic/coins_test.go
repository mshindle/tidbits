package dynamic

import "testing"

var standard = []int{1, 5, 10, 25, 100}

func TestCoins(t *testing.T) {
	type args struct {
		val    int
		denoms []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case 1",
			args: args{27, standard},
			want: 3,
		},
		{
			name: "case 2",
			args: args{7, []int{2, 4}},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Coins(tt.args.val, tt.args.denoms); got != tt.want {
				t.Errorf("Coins() = %v, want %v", got, tt.want)
			}
		})
	}
}
