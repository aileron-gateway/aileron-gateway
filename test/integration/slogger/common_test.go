// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package slogger_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
var testDataDir = "../../../test/integration/slogger/"

func testLogger(t *testing.T, lg log.Logger) {

	t.Helper()

	ctx := context.Background()
	lg.Debug(ctx, "test debug", "name", "alice", "age", 20)
	lg.Info(ctx, "test info", "name", "alice", "age", 20)
	lg.Warn(ctx, "test warn", "name", "alice", "age", 20)
	lg.Error(ctx, "test error", "name", "alice", "age", 20)

}

func TestMinimalWithoutMetadata(t *testing.T) {

	configs := []string{
		testDataDir + "config-minimal-without-metadata.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testLogger(t, lg)

}

func TestMinimalWithMetadata(t *testing.T) {

	configs := []string{
		testDataDir + "config-minimal-with-metadata.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testLogger(t, lg)

}

func TestEmptyName(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-name.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "testNamespace",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testLogger(t, lg)

}

func TestEmptyNamespace(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-namespace.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "testName",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testLogger(t, lg)

}

func TestEmptyNameNamespace(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-name-namespace.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testLogger(t, lg)

}

func TestEmptySpec(t *testing.T) {

	configs := []string{
		testDataDir + "config-empty-spec.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testLogger(t, lg)

}

func TestInvalidSpec(t *testing.T) {

	configs := []string{
		testDataDir + "config-invalid-spec.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to load configs.`)
	testutil.DiffError(t, app.ErrAppMainLoadConfigs, errPattern, err)

}
