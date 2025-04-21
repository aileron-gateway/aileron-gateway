// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package entrypoint_test

import (
	"cmp"
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

// testDataDir is the path to the test data.
// This path can be changed by the environmental variable.
var testDataDir = cmp.Or(os.Getenv("TEST_DIR"), "../../../test/") + "integration/entrypoint/"

func TestMinimalWithoutMetadata(t *testing.T) {

	configs := []string{
		testDataDir + "config-minimal-without-metadata.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "Entrypoint",
		Name:       ".entrypoint", // Name cannot changed by the config.
		Namespace:  ".entrypoint", // Namespace cannot changed by the config.
	}
	entrypoint, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = entrypoint.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)

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
		Kind:       "Entrypoint",
		Name:       ".entrypoint", // Name cannot changed by the config.
		Namespace:  ".entrypoint", // Namespace cannot changed by the config.
	}
	entrypoint, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = entrypoint.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)

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
		Kind:       "Entrypoint",
		Name:       ".entrypoint", // Name cannot changed by the config.
		Namespace:  ".entrypoint", // Namespace cannot changed by the config.
	}
	entrypoint, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = entrypoint.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)
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
		Kind:       "Entrypoint",
		Name:       ".entrypoint", // Name cannot changed by the config.
		Namespace:  ".entrypoint", // Namespace cannot changed by the config.
	}
	entrypoint, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = entrypoint.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)

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
		Kind:       "Entrypoint",
		Name:       ".entrypoint", // Name cannot changed by the config.
		Namespace:  ".entrypoint", // Namespace cannot changed by the config.
	}
	entrypoint, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = entrypoint.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)

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
		Kind:       "Entrypoint",
		Name:       ".entrypoint", // Name cannot changed by the config.
		Namespace:  ".entrypoint", // Namespace cannot changed by the config.
	}
	entrypoint, err := api.ReferTypedObject[core.Runner](server, ref)
	testutil.DiffError(t, nil, nil, err)

	err = entrypoint.Run(context.Background())
	testutil.DiffError(t, nil, nil, err)

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
