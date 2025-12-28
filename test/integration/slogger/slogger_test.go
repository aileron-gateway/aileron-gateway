// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package slogger_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLevel(t *testing.T) {

	configs := []string{
		testDataDir + "config-level.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testutil.Diff(t, false, lg.Enabled(log.LvTrace))
	testutil.Diff(t, false, lg.Enabled(log.LvDebug))
	testutil.Diff(t, false, lg.Enabled(log.LvInfo))
	testutil.Diff(t, true, lg.Enabled(log.LvWarn))
	testutil.Diff(t, true, lg.Enabled(log.LvError))
	testutil.Diff(t, true, lg.Enabled(log.LvFatal))

	ctx := context.Background()
	lg.Debug(ctx, "test debug", "name", "alice")
	lg.Info(ctx, "test info", "name", "alice")
	lg.Warn(ctx, "test warn", "name", "alice")
	lg.Error(ctx, "test error", "name", "alice")

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	testutil.Diff(t, false, strings.Contains(buf.String(), `"level":"DEBUG"`))
	testutil.Diff(t, false, strings.Contains(buf.String(), `"level":"INFO"`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `"level":"WARN"`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `"level":"ERROR"`))

}

func TestUnstructured(t *testing.T) {

	configs := []string{
		testDataDir + "config-unstructured.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stdout = tmp }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testutil.Diff(t, false, lg.Enabled(log.LvTrace))
	testutil.Diff(t, false, lg.Enabled(log.LvDebug))
	testutil.Diff(t, true, lg.Enabled(log.LvInfo))
	testutil.Diff(t, true, lg.Enabled(log.LvWarn))
	testutil.Diff(t, true, lg.Enabled(log.LvError))
	testutil.Diff(t, true, lg.Enabled(log.LvFatal))

	ctx := context.Background()
	lg.Debug(ctx, "test debug", "name", "alice")
	lg.Info(ctx, "test info", "name", "alice")
	lg.Warn(ctx, "test warn", "name", "alice")
	lg.Error(ctx, "test error", "name", "alice")

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	testutil.Diff(t, false, strings.Contains(buf.String(), `level=DEBUG`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `level=INFO`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `level=WARN`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `level=ERROR`))

}

func TestOutputStderr(t *testing.T) {

	configs := []string{
		testDataDir + "config-output-stderr.yaml",
	}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	tmp := os.Stdout
	defer func() { os.Stderr = tmp }()
	r, w, _ := os.Pipe()
	os.Stderr = w

	ref := &kernel.Reference{
		APIVersion: "core/v1",
		Kind:       "SLogger",
		Name:       "default",
		Namespace:  "",
	}
	lg, err := api.ReferTypedObject[log.Logger](server, ref)
	testutil.DiffError(t, nil, nil, err)

	testutil.Diff(t, false, lg.Enabled(log.LvTrace))
	testutil.Diff(t, false, lg.Enabled(log.LvDebug))
	testutil.Diff(t, true, lg.Enabled(log.LvInfo))
	testutil.Diff(t, true, lg.Enabled(log.LvWarn))
	testutil.Diff(t, true, lg.Enabled(log.LvError))
	testutil.Diff(t, true, lg.Enabled(log.LvFatal))

	ctx := context.Background()
	lg.Debug(ctx, "test debug", "name", "alice")
	lg.Info(ctx, "test info", "name", "alice")
	lg.Warn(ctx, "test warn", "name", "alice")
	lg.Error(ctx, "test error", "name", "alice")

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	testutil.Diff(t, false, strings.Contains(buf.String(), `"level":"DEBUG"`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `"level":"INFO"`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `"level":"WARN"`))
	testutil.Diff(t, true, strings.Contains(buf.String(), `"level":"ERROR"`))

}

func TestOutputFileRotate(t *testing.T) {

	configs := []string{
		testDataDir + "config-output-file-rotate.yaml",
	}

	os.Mkdir("./logDir", os.ModePerm)
	os.Mkdir("./backupDir", os.ModePerm)
	defer func() {
		os.RemoveAll("./logDir")
		os.RemoveAll("./backupDir")
	}()

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

	// rotate log file once.
	ctx := context.Background()
	for i := 0; i < 5_000; i++ {
		lg.Error(ctx, "test error", "name", "alice", "age", 20) // 260 bytes per 1 line.
	}
	lg.(core.Finalizer).Finalize()

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	var logDirFiles []string
	files, err := os.ReadDir("./logDir")
	for _, f := range files {
		logDirFiles = append(logDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{}, logDirFiles, opts...)

	var backupDirFiles []string
	files, err = os.ReadDir("./backupDir")
	for _, f := range files {
		backupDirFiles = append(backupDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{"application.1.log", "application.2.log"}, backupDirFiles, opts...)
}

func TestOutputFileMaxBackup(t *testing.T) {
	configs := []string{
		testDataDir + "config-output-file-max-backup.yaml",
	}

	os.Mkdir("./logDir", os.ModePerm)
	os.Mkdir("./backupDir", os.ModePerm)
	defer func() {
		os.RemoveAll("./logDir")
		os.RemoveAll("./backupDir")
	}()

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

	// rotate log file twice.
	ctx := context.Background()
	for i := 0; i < 13_000; i++ {
		lg.Error(ctx, "test error", "name", "alice", "age", 20) // 263 bytes per 1 line.
	}
	lg.(core.Finalizer).Finalize()

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	var logDirFiles []string
	files, err := os.ReadDir("./logDir")
	for _, f := range files {
		logDirFiles = append(logDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{}, logDirFiles, opts...)

	var backupDirFiles []string
	files, err = os.ReadDir("./backupDir")
	for _, f := range files {
		backupDirFiles = append(backupDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	// "application.1.log" and "application.2.log" should be removed because of the maxBackup=2.
	testutil.Diff(t, []string{"application.3.log", "application.4.log"}, backupDirFiles, opts...)
}

func TestOutputFileMaxTotal(t *testing.T) {

	configs := []string{
		testDataDir + "config-output-file-max-total.yaml",
	}

	os.Mkdir("./logDir", os.ModePerm)
	os.Mkdir("./backupDir", os.ModePerm)
	defer func() {
		os.RemoveAll("./logDir")
		os.RemoveAll("./backupDir")
	}()

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

	// rotate log file twice.
	ctx := context.Background()
	for i := 0; i < 16_000; i++ {
		lg.Error(ctx, "test error", "name", "alice", "age", 20) // 262 bytes per 1 line.
	}
	lg.(core.Finalizer).Finalize()

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	var logDirFiles []string
	files, err := os.ReadDir("./logDir")
	for _, f := range files {
		logDirFiles = append(logDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{}, logDirFiles, opts...)

	var backupDirFiles []string
	files, err = os.ReadDir("./backupDir")
	for _, f := range files {
		backupDirFiles = append(backupDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	// "application.1.log" - "application.2.log" exceed the max total size 2MiB.
	// The size does not include the last active log "application.3.log"
	testutil.Diff(t, []string{"application.3.log", "application.4.log"}, backupDirFiles, opts...)

}

func TestOutputFileNoDir_logDir(t *testing.T) {

	configs := []string{
		testDataDir + "config-output-file-nodir.yaml",
	}

	// os.Mkdir("./logDir", os.ModePerm) // Do not create the directory.
	os.Mkdir("./backupDir", os.ModePerm)
	defer func() {
		os.RemoveAll("./logDir")
		os.RemoveAll("./backupDir")
	}()

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

	// rotate log file once.
	ctx := context.Background()
	for i := 0; i < 6_000; i++ {
		lg.Error(ctx, "test error", "name", "alice", "age", 20) // 262 bytes per 1 line.
	}
	lg.(core.Finalizer).Finalize()

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	var logDirFiles []string
	files, err := os.ReadDir("./logDir") // output dir will be created.
	for _, f := range files {
		logDirFiles = append(logDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{}, logDirFiles, opts...)

	var backupDirFiles []string
	files, err = os.ReadDir("./backupDir")
	for _, f := range files {
		backupDirFiles = append(backupDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{"application.1.log", "application.2.log"}, backupDirFiles, opts...)

}

func TestOutputFileNoDir_backupDir(t *testing.T) {

	configs := []string{
		testDataDir + "config-output-file-nodir.yaml",
	}

	os.Mkdir("./logDir", os.ModePerm)
	// os.Mkdir("./backupDir", os.ModePerm)  // Do not create the directory.
	defer func() {
		os.RemoveAll("./logDir")
		os.RemoveAll("./backupDir")
	}()

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

	// rotate log file once.
	ctx := context.Background()
	for i := 0; i < 6_000; i++ {
		lg.Error(ctx, "test error", "name", "alice", "age", 20) // 262 bytes per 1 line.
	}
	lg.(core.Finalizer).Finalize()

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	var logDirFiles []string
	files, err := os.ReadDir("./logDir")
	for _, f := range files {
		logDirFiles = append(logDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{}, logDirFiles, opts...)

	var backupDirFiles []string
	files, err = os.ReadDir("./backupDir") // backup dir will be created.
	for _, f := range files {
		backupDirFiles = append(backupDirFiles, f.Name())
	}
	testutil.Diff(t, nil, err)
	testutil.Diff(t, []string{"application.1.log", "application.2.log"}, backupDirFiles, opts...)

}
