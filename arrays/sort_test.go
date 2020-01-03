package arrays

import (
	"reflect"
	"testing"
)

func TestSubarraySort(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case 1",
			args: args{arr: []int{1, 2}},
			want: []int{-1, -1},
		},
		{
			name: "case 2",
			args: args{arr: []int{2, 1}},
			want: []int{0, 1},
		},
		{
			name: "case 3",
			args: args{arr: []int{1, 2, 4, 7, 10, 11, 7, 12, 6, 7, 16, 18, 19}},
			want: []int{3, 9},
		},
		{
			name: "case 4",
			args: args{arr: []int{1, 2, 4, 7, 10, 11, 7, 12, 7, 7, 16, 18, 19}},
			want: []int{4, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubarraySort(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubarraySort() = %v, want %v", got, tt.want)
			}
		})
	}
}
