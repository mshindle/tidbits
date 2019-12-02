package toy

import "testing"

func TestProductSum(t *testing.T) {
	type args struct {
		array []interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"test_1",
			args{SpecialArray{1, 2, 3, 4, 5}},
			15,
		},
		{
			"test_2",
			args{SpecialArray{1, 2, SpecialArray{3}, 4, 5}},
			18,
		},
		{
			"test_3",
			args{SpecialArray{1, 2, SpecialArray{3}, 4, 5}},
			18,
		},
		{
			"test_4",
			args{SpecialArray{1, 2, SpecialArray{3, SpecialArray{6, 7, 8, 9, 10}}, 4, 5}},
			258,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProductSum(tt.args.array); got != tt.want {
				t.Errorf("ProductSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
