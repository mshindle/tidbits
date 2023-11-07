package structures

import (
	"fmt"
	"testing"
)

func TestSame(t *testing.T) {
	type args struct {
		t1 *IntTree
		t2 *IntTree
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "two trees should be the same",
			args: args{
				t1: NewIntTree(1),
				t2: NewIntTree(1),
			},
			want: true,
		},
		{
			name: "two trees should be different",
			args: args{
				t1: NewIntTree(2),
				t2: NewIntTree(5),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Same(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("Same() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvert(t *testing.T) {
	type args struct {
		t *IntTree
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "invert a small tree",
			args: args{
				t: NewIntTree(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.t.String())
			Invert(tt.args.t)
			fmt.Println(tt.args.t.String())
		})
	}
}
