package concurrency

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apex/log"
)

var taskCount atomic.Int64

func watch(ctx context.Context, tasks <-chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-tasks:
			if !ok {
				return
			}
			taskCount.Add(v)
		}
	}
}

func streamData(ctx context.Context, v int64, stream chan int64) {
	defer close(stream)
	for {
		select {
		case <-ctx.Done():
			return
		case stream <- v:
		}
	}
}

func Leaking() bool {
	tasks := make(chan int64)

	// start streaming data
	ctxStream, cancelStream := context.WithCancel(context.Background())
	defer cancelStream()
	go streamData(ctxStream, 1, tasks)

	// create the consumers
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.WithField("consumer", i).Info("starting watcher")
		go watch(ctx, tasks, wg)
	}

	log.Info("sleeping for 100ms; cancel consumer context")
	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	before := taskCount.Load()
	log.WithField("taskCount", before).Info("grabbing before sleeping")
	time.Sleep(100 * time.Millisecond)
	after := taskCount.Load()
	log.WithField("taskCount", after).Info("grabbing after sleeping")
	return before == after
}

// LeakingCascade runs the same goroutines, but instead of canceling the consumers directly,
// it cancels the context of the streamer and allows the closed stream to
// terminate the consumers.
func LeakingCascade() bool {
	tasks := make(chan int64)

	// start streaming data
	ctxStream, cancelStream := context.WithCancel(context.Background())
	go streamData(ctxStream, 1, tasks)

	// create the consumers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.WithField("consumer", i).Info("starting watcher")
		go watch(ctx, tasks, wg)
	}

	log.Info("sleep for 100ms; cancel stream context")
	time.Sleep(100 * time.Millisecond)
	cancelStream()
	before := taskCount.Load()
	log.WithField("taskCount", before).Info("grabbing after stream cancelled")

	// wait for all consumers to exit
	wg.Wait()
	log.Info("all consumers exited")

	after := taskCount.Load()
	log.WithField("taskCount", after).
		Info("time for process to drain completely")

	return true
}
