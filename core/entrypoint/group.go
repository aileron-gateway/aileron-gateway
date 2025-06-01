// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package entrypoint

import (
	"context"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
)

// channelGroup runs multiple runner in a new goroutine.
// This implements core.Runner interface.
type channelGroup struct {
	lg           log.Logger
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
