package rect

import "testing"

func TestCalcNumRectangles(t *testing.T) {
	type args struct {
		grid Grid
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "single rect",
			args: args{
				grid: Grid{
					Point{0, 0},
					Point{0, 1},
					Point{1, 0},
					Point{1, 1},
				},
			},
			want: 1,
		},
		{
			name: "three rects",
			args: args{
				grid: Grid{
					Point{0, 0},
					Point{0, 1},
					Point{1, 0},
					Point{1, 1},
					Point{2, 0},
					Point{2, 1},
				},
			},
			want: 3,
		},
		{
			name: "missing point",
			args: args{
				grid: Grid{
					Point{0, 0},
					Point{0, 1},
					Point{1, 0},
					Point{2, 0},
					Point{2, 1},
				},
			},
			want: 1,
		},
		{
			name: "three rows",
			args: args{
				grid: Grid{
					Point{0, 0},
					Point{0, 1},
					Point{0, 2},
					Point{1, 0},
					Point{1, 1},
					Point{1, 2},
					Point{2, 0},
					Point{2, 1},
					Point{2, 2},
				},
			},
			want: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcNumRectangles(tt.args.grid); got != tt.want {
				t.Errorf("CalcNumRectangles() = %v, want %v", got, tt.want)
			}
		})
	}
}
