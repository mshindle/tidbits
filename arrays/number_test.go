package arrays

import (
	"reflect"
	"testing"
)

func TestTwoNumberSum(t *testing.T) {
	type args struct {
		array  []int
		target int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "case 1",
			args: args{[]int{4, 6}, 10},
			want: [][]int{{4, 6}},
		},
		{
			name: "case 2",
			args: args{[]int{3, 5, -4, 8, 11, 1, -1, 6}, 15},
			want: [][]int{},
		},
		{
			name: "case 3",
			args: args{[]int{3, 5, -4, 8, 11, 1, -1, 6}, 11},
			want: [][]int{{3, 8}, {5, 6}},
		},
		{
			name: "case 4",
			args: args{[]int{-5, -7, -3, -1, 0, 1, 3, 5, 7, 2, -2}, -5},
			want: [][]int{{-7, 2}, {-5, 0}, {-3, -2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TwoNumberSum(tt.args.array, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwoNumberSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThreeNumberSum(t *testing.T) {
	type args struct {
		array  []int
		target int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "case 1",
			args: args{[]int{1, 2, 3}, 6},
			want: [][]int{[]int{1, 2, 3}},
		},
		{
			name: "case 2",
			args: args{[]int{1, 2, 3}, 7},
			want: [][]int{},
		},
		{
			name: "case 3",
			args: args{[]int{12, 3, 1, 2, -6, 5, 0, -8, -1}, 0},
			want: [][]int{[]int{-8, 3, 5}, []int{-6, 1, 5}, []int{-1, 0, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ThreeNumberSum(tt.args.array, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThreeNumberSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
