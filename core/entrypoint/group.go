// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package entrypoint

import (
	"context"
	"errors"
	"sync"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
)

// channelGroup runs multiple runner in a new goroutine.
// This implements core.Runner interface.
type channelGroup struct {
	lg log.Logger

	runners      []core.Runner
	initializers []core.Initializer
	finalizers   []core.Finalizer
}

// Run run the registered runner in a new goroutine.
// This is the implements of core.Runner.Run method.
// This returns an error or nil when at least one runner is done.
func (g *channelGroup) Run(ctx context.Context) (err error) {
	if len(g.runners) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err = g.initialize(); err != nil {
		return err
	}
	defer func() {
		e := g.finalize()
		if err == nil && e != nil {
			err = e
		}
	}()

	errChan := make(chan error)

	for _, runner := range g.runners {
		r := runner

		waitGoroutine := make(chan struct{})
		go func() {
			close(waitGoroutine)
			if err := r.Run(ctx); err != nil {
				errChan <- err
				return
			}
			errChan <- nil
		}()
		<-waitGoroutine
	}

	// Return an error if any.
	if err := <-errChan; err != nil {
		err := core.ErrCoreEntrypointRun.WithStack(err, nil)
		g.lg.Error(ctx, "error returned from runner", err.Name(), err.Map())
		return err
	}

	return nil
}

func (g *channelGroup) initialize() error {
	for _, f := range g.initializers {
		if f == nil {
			continue
		}
		if err := f.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

func (g *channelGroup) finalize() error {
	for _, f := range g.finalizers {
		if f == nil {
			continue
		}
		if err := f.Finalize(); err != nil {
			return err
		}
	}
	return nil
}

// waitGroup runs multiple runner in a new goroutine.
// This implements core.Runner interface.
type waitGroup struct {
	lg log.Logger

	// wg is the wait group that runs multiple runners.
	wg sync.WaitGroup

	// runners are the runner objects which are run
	// as a member of the waitgroup above.
	runners      []core.Runner
	initializers []core.Initializer
	finalizers   []core.Finalizer
}

// Run run the registered runner in a new goroutine.
// This is the implementation of core.Runner.Run method.
// This returns an error or nil when all runners are done.
func (g *waitGroup) Run(ctx context.Context) (err error) {
	if len(g.runners) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err = g.initialize(); err != nil {
		return err
	}
	defer func() {
		e := g.finalize()
		if err == nil && e != nil {
			err = e
		}
	}()

	// errs collects all errors returned from the runners.
	var errs []error
	// mu protects errs.
	var mu sync.Mutex

	for _, runner := range g.runners {
		r := runner
		g.wg.Add(1)

		waitGoroutine := make(chan struct{})

		go func() {
			close(waitGoroutine)
			defer g.wg.Done()
			if err := r.Run(ctx); err != nil {
				mu.Lock()
				defer mu.Unlock()
				errs = append(errs, err)
			}
		}()

		<-waitGoroutine
	}

	// Wait until all goroutines are done.
	g.wg.Wait()

	// Return errors if any.
	if len(errs) > 0 {
		err := core.ErrCoreEntrypointRun.WithStack(errors.Join(errs...), nil)
		g.lg.Error(ctx, "error returned from runner", err.Name(), err.Map())
		return err
	}

	return nil
}

func (g *waitGroup) initialize() error {
	for _, f := range g.initializers {
		if f == nil {
			continue
		}
		if err := f.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

func (g *waitGroup) finalize() error {
	for _, f := range g.finalizers {
		if f == nil {
			continue
		}
		if err := f.Finalize(); err != nil {
			return err
		}
	}
	return nil
}
