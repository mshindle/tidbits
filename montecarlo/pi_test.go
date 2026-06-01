package montecarlo

import (
	"context"
	"testing"
	"time"
)

// TestCompute_Success ensures that the Monte Carlo computation finishes
// normally and yields a mathematically reasonable estimation of PI.
func TestCompute_Success(t *testing.T) {
	ctx := context.Background()

	// Create PI engine with 100k points, 4 workers, a no-op tracer, and a UI sample capacity of 2500
	piCalculator := NewPI(100000, 4, 2500)

	err := piCalculator.Compute(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Validate using the new atomic.Int64 Load() method
	if piCalculator.TotalPoints.Load() != 100000 {
		t.Errorf("expected 100000 total points processed, got %d", piCalculator.TotalPoints.Load())
	}

	// PI should logically hover around 3.14 given enough samples
	if piCalculator.Value < 2.5 || piCalculator.Value > 3.8 {
		t.Errorf("calculated PI value out of expected rough bounds: %f", piCalculator.Value)
	}

	// Validate the new visualizer pipeline collected sampled points
	if piCalculator.SampleBuffer.Len() == 0 {
		t.Error("expected visualization sample buffer to contain points, but it was empty")
	}
}

// TestCompute_ZeroPoints ensures that entering an empty workload
// doesn't trigger a division-by-zero panic or block indefinitely.
func TestCompute_ZeroPoints(t *testing.T) {
	ctx := context.Background()
	piCalculator := NewPI(0, 2, 100)

	err := piCalculator.Compute(ctx)
	if err == nil {
		t.Fatalf("expected error for 0 points, got: %v", err)
	}

	if piCalculator.TotalPoints.Load() != 0 {
		t.Errorf("expected 0 total points, got %d", piCalculator.TotalPoints.Load())
	}
}

// TestCompute_ContextCancellation guarantees that when the parent context
// is canceled early, Compute() stops instantly and doesn't leak blocked workers.
func TestCompute_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Set up a massive simulation designed to run a long time if unchecked
	piCalculator := NewPI(1000000000, 4, 2500)

	// Trigger cancellation asynchronously after a fraction of a second
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// Channel to monitor completion
	done := make(chan error, 1)
	go func() {
		done <- piCalculator.Compute(ctx)
	}()

	// Enforce a strict test timeout to catch channel deadlocks
	select {
	case err := <-done:
		if err == nil {
			t.Error("expected context cancellation error, got nil")
		}
		if ctx.Err() == nil {
			t.Error("expected parent context status to register as canceled")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("test deadlocked: Compute() failed to unwind workers after context cancellation")
	}
}

// TestGenerate_ContextCancellation specifically verifies that the data generator
// unblocks its thread immediately if the pipeline down-stream listener abandons it.
func TestGenerate_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a loop that generates an infinite flow of points
	out := generate(ctx, 999999999)

	// Consume a single point to prove the generator works
	_, ok := <-out
	if !ok {
		t.Fatal("expected channel to stream a point, but it was prematurely closed")
	}

	// Abruptly cancel context and abandon reading the stream
	cancel()

	// Verify the background thread detected the cancellation and closed the outbox
	done := make(chan struct{})
	go func() {
		for range out {
		}
		close(done)
	}()

	select {
	case <-done:
		// Success: the generator goroutine broke its loop and safely closed the channel
	case <-time.After(200 * time.Millisecond):
		t.Fatal("generator goroutine leaked: failed to close channel after context abort")
	}
}
