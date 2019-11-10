package rect

import "github.com/sirupsen/logrus"

type Point struct {
	X int
	Y int
}

type Grid []Point

type pair struct {
	a int
	b int
}

// Calculate the number of rectangles given a 2-dimensional grid of points
// We calculate the number of rectangles by counting the number of vertical heights
// we have that are the same size
func CalcNumRectangles(grid Grid) int {
	var ans int
	pairCount := make(map[pair]int)
	for _, pt := range grid {
		for _, ptAbove := range grid {
			if pt.X == ptAbove.X && pt.Y < ptAbove.Y {
				p := pair{pt.Y, ptAbove.Y}
				ans = ans + pairCount[p]
				pairCount[p]++
				logrus.
					WithFields(logrus.Fields{"pt": pt, "p": p, "ans": ans, "pairCount": pairCount[p]}).
					Debug("current state")
			}
		}
	}
	return ans
}
