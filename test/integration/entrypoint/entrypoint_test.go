// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package entrypoint_test

import (
	"context"
	"io"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

type mocRunner struct {
	err    error
	called int
}

func (r *mocRunner) Run(ctx context.Context) error {
	r.called += 1
	return r.err
}

func TestRunner1(t *testing.T) {

	configs := []string{
		testDataDir + "config-runner-1.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test",
		Namespace:  "",
	}
	mocRunner := &mocRunner{}
	common.PostTestResource(server, mocRef, mocRunner)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "Entrypoint",
		Name:       ".entrypoint",
		Namespace:  ".entrypoint",
	}
	r, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = r.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, 1, mocRunner.called)

}

func TestRunner2(t *testing.T) {

	configs := []string{testDataDir + "config-runner-2.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	mocRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test1",
		Namespace:  "",
	}
	mocRunner1 := &mocRunner{}
	common.PostTestResource(server, mocRef1, mocRunner1)
	mocRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test2",
		Namespace:  "",
	}
	mocRunner2 := &mocRunner{}
	common.PostTestResource(server, mocRef2, mocRunner2)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "Entrypoint",
		Name:       ".entrypoint",
		Namespace:  ".entrypoint",
	}
	r, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = r.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, 1, mocRunner1.called)
	testutil.Diff(t, 1, mocRunner2.called)

}

func TestRunnerWait1(t *testing.T) {

	configs := []string{
		testDataDir + "config-runner-wait-1.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test",
		Namespace:  "",
	}
	mocRunner := &mocRunner{}
	common.PostTestResource(server, mocRef, mocRunner)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "Entrypoint",
		Name:       ".entrypoint",
		Namespace:  ".entrypoint",
	}
	r, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = r.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, 1, mocRunner.called)

}

func TestRunnerWait2(t *testing.T) {

	configs := []string{
		testDataDir + "config-runner-wait-2.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	mocRef1 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test1",
		Namespace:  "",
	}
	mocRunner1 := &mocRunner{}
	common.PostTestResource(server, mocRef1, mocRunner1)
	mocRef2 := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test2",
		Namespace:  "",
	}
	mocRunner2 := &mocRunner{}
	common.PostTestResource(server, mocRef2, mocRunner2)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "Entrypoint",
		Name:       ".entrypoint",
		Namespace:  ".entrypoint",
	}
	r, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = r.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)
	testutil.Diff(t, 1, mocRunner1.called)
	testutil.Diff(t, 1, mocRunner2.called)

}

func TestRunnerError(t *testing.T) {

	configs := []string{
		testDataDir + "config-runner-error.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	mocRef := &kernel.Reference{
		APIVersion: "container/v1",
		Kind:       "TestRunner",
		Name:       "test",
		Namespace:  "",
	}
	mocRunner := &mocRunner{err: io.EOF}
	common.PostTestResource(server, mocRef, mocRunner)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "Entrypoint",
		Name:       ".entrypoint",
		Namespace:  ".entrypoint",
	}
	r, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = r.Run(context.Background())
	errPattern := regexp.MustCompile(core.ErrPrefix + `error on running entrypoint`)
	testutil.DiffError(t, core.ErrCoreEntrypointRun, errPattern, err)
	testutil.Diff(t, 1, mocRunner.called)

}
