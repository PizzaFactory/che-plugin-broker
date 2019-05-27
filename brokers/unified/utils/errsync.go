package utils

import (
	"context"
	"sync"
)

type WaitGroup struct {
	wg sync.WaitGroup

	/* Set `err` only inside `sync.Once` for avoiding over-written `nil` by the other goroutine. */
	errOnce sync.Once
	err     error
}

// The wrapper for generating goroutine..
func (g *WaitGroup) RunInParallel(fn func() error) {
	if g.err != nil {
		// Just return if the other gorutine was failed with the error.
		return;
	}

	g.wg.Add(1) // Will be closed by `Done()` in the goroutine below.

	go func() {
		defer g.wg.Done()

		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

// Wait until finishing all parallel jobs.
// Returns the first detected error if avaliable, or `nil` if there is no error.
func (g *WaitGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

