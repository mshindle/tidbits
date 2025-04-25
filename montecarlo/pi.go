package montecarlo

import (
	"math"
	"math/rand"
	"sync"

	"github.com/apex/log"
)

const radius = 1.0

type point struct {
	x float64
	y float64
}

type tally struct {
	inCircle int64
	total    int64
}

type PI struct {
	Points      int64
	Value       float64
	InCircle    int64
	TotalPoints int64
	numWorkers  int
}

func NewPI(points int64, numWorkers int) *PI {
	p := &PI{
		Points:     points,
		numWorkers: numWorkers,
	}
	return p
}

func (p *PI) Compute() error {
	tchans := make([]<-chan tally, 0, p.numWorkers)
	pchan := generate(p.Points)
	for i := 0; i < p.numWorkers; i++ {
		tchans = append(tchans, worker(pchan))
	}
	results := merge(tchans...)
	for t := range results {
		p.InCircle += t.inCircle
		p.TotalPoints += t.total
	}
	p.Value = 4.0 * (float64(p.InCircle) / float64(p.TotalPoints))

	log.WithFields(log.Fields{
		"inCircle": p.InCircle,
		"total":    p.TotalPoints,
		"pi":       p.Value,
	}).Info("pi calculated")
	return nil
}

func generate(points int64) <-chan point {
	out := make(chan point)
	go func(n int64) {
		defer close(out)
		for n > 0 {
			pt := point{
				x: (rand.Float64() * 2) - 1,
				y: (rand.Float64() * 2) - 1,
			}
			out <- pt
			n--
		}
	}(points)
	return out
}

func worker(points <-chan point) <-chan tally {
	out := make(chan tally)

	go func() {
		var t tally
		defer close(out)
		for pt := range points {
			t.total++
			distance := math.Sqrt(pt.x*pt.x + pt.y*pt.y)
			if distance <= radius {
				t.inCircle++
			}
		}
		out <- t
	}()
	return out
}

// merge all the worker output into one stream
func merge(tallies ...<-chan tally) <-chan tally {
	var wg sync.WaitGroup
	out := make(chan tally)

	output := func(tchan <-chan tally) {
		defer wg.Done()
		for t := range tchan {
			out <- t
		}
	}
	wg.Add(len(tallies))
	for _, tchan := range tallies {
		go output(tchan)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
