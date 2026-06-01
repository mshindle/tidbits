package montecarlo

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"
	"sync"
	"sync/atomic"

	"github.com/apex/log"
	"github.com/mshindle/structures/ringbuffer"
)

const radius = 1.0

type point struct {
	x float64
	y float64
}

// RenderPoint represents the sampled output payload for our data visualization
type RenderPoint struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	InCircle bool    `json:"in_circle"`
}

type PI struct {
	Points       int64
	Value        float64
	InCircle     atomic.Int64
	TotalPoints  atomic.Int64
	numWorkers   int
	SampleBuffer *ringbuffer.RingBuffer[RenderPoint]
}

func NewPI(points int64, numWorkers, sampleCapacity int) *PI {
	p := &PI{
		Points:       points,
		numWorkers:   numWorkers,
		SampleBuffer: ringbuffer.New[RenderPoint](sampleCapacity),
	}
	return p
}

func (p *PI) Compute(ctx context.Context) error {
	pChan := generate(ctx, p.Points)

	var wg sync.WaitGroup
	for i := 0; i < p.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.worker(ctx, pChan)
		}()
	}
	wg.Wait()

	in := p.InCircle.Load()
	total := p.TotalPoints.Load()
	if total == 0 {
		return errors.New("no points processed, cannot calculate PI")
	}

	// grab the value of PI.
	p.Value = 4.0 * (float64(in) / float64(total))

	log.WithFields(log.Fields{
		"inCircle": in,
		"total":    total,
		"pi":       p.Value,
	}).Info("pi calculated")
	return ctx.Err()
}

func generate(ctx context.Context, points int64) <-chan point {
	out := make(chan point)
	go func(n int64) {
		defer close(out)
		for n > 0 {
			pt := point{
				x: (rand.Float64() * 2) - 1,
				y: (rand.Float64() * 2) - 1,
			}
			select {
			case <-ctx.Done():
				return
			case out <- pt:
				n--
			}
		}
	}(points)
	return out
}

func (p *PI) worker(ctx context.Context, points <-chan point) {
	var localTotal, localIn int64

	for {
		select {
		case <-ctx.Done():
			return
		case pt, ok := <-points:
			if !ok {
				// Flush any remaining tallies before exiting
				if localTotal > 0 {
					p.TotalPoints.Add(localTotal)
					p.InCircle.Add(localIn)
				}
				return
			}
			localTotal++

			inCircle := false
			distance := math.Sqrt(pt.x*pt.x + pt.y*pt.y)
			if distance <= radius {
				localIn++
				inCircle = true
			}

			// Stratified Sampling: Only sample 5% of points into the visualization buffer
			// to protect memory bandwidth while remaining statistically accurate.
			if localTotal%50 == 0 {
				p.SampleBuffer.OverwritePush(RenderPoint{
					X:        pt.x,
					Y:        pt.y,
					InCircle: inCircle,
				})
			}

			// Batch flush to atomics to prevent cache line contention
			if localTotal == 1000 {
				p.TotalPoints.Add(localTotal)
				p.InCircle.Add(localIn)
				localTotal = 0
				localIn = 0
			}
		}
	}
}
