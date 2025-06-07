// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-projects/go/zos"
	"github.com/spf13/pflag"
)

// Exit is the function that will be called
// when exiting the application.
// This can be replaced for testing.
var Exit func(int) = os.Exit

// New returns a new instance os app.App struct.
func New() *App {
	return &App{}
}

// App is the application.
// Use app.New() to get a new instance of this struct.
// This implements core.Runner interface.
type App struct {
	args []string
	opts *Options
}

// ParseArgs parse arguments.
// Custom or used-defined flag sets can be given at the second argument.
// Custom flag sets will be filled whe parsing args.
// When the given args were empty, default args of ["-h"] will be used.
func (a *App) ParseArgs(args []string, custom ...*pflag.FlagSet) {
	a.args = append(a.args, args...)
	if len(args) == 0 {
		a.args = append(a.args, "-h")
	}
	a.opts = ParseArgs(a.args, custom...)
}

// Run runs this app.
func (a *App) Run(server api.API[*api.Request, *api.Response]) error {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	// Load env files before loading config files.
	if err := LoadEnvFiles(a.opts.Basic.Envs); err != nil {
		return err // Return err as-is.
	}
	// Load config files.
	// Environmental variables in the configs will be resolved.
	if err := LoadConfigFiles(server, a.opts.Basic.Configs); err != nil {
		return err // Return err as-is.
	}

	// Show resource template and exit.
	ShowTemplate(server, a.opts.Basic.Template, a.opts.Basic.Out)

	// Get the entrypoint resource and run it.
	// Entrypoint must implement core.Runner interface.
	req := &api.Request{
		Method: api.MethodGet,
		Key:    "core/v1/Entrypoint",
		Format: api.FormatProtoReference,
		Content: &k.Reference{
			APIVersion: "core/v1",
			Kind:       "Entrypoint",
			Namespace:  ".entrypoint",
			Name:       ".entrypoint",
		},
	}

	res, err := server.Serve(ctx, req)
	if err != nil {
		return ErrAppMainGetEntrypoint.WithStack(err, nil)
	}
	entrypoint, ok := res.Content.(core.Runner)
	if !ok {
		err := fmt.Errorf("entrypoint type %T is not runnable because it lacks core.Runner interface", res.Content)
		return ErrAppMainGetEntrypoint.WithStack(err, nil)
	}

	// Run the entrypoint runner.
	if err := entrypoint.Run(ctx); err != nil {
		return ErrAppMainRun.WithStack(err, nil)
	}

	return nil
}

// LoadEnvFiles load environmental variables from given file paths.
func LoadEnvFiles(paths []string) error {
	envs, err := zos.ReadFiles(true, paths...)
	if err != nil {
		return ErrAppMainLoadEnv.WithStack(err, nil)
	}
	for k, v := range envs {
		if _, err := zos.LoadEnv(v); err != nil {
			return ErrAppMainLoadEnv.WithStack(err, map[string]any{"path": k})
		}
	}
	return nil
}

// LoadConfigFiles load config files in json or yaml format from given file paths.
// Currently only ".json", ".yaml" and ".yml" file extensions are supported.
// Others will be ignored.
// This function panics when the given server is nil.
func LoadConfigFiles(server api.API[*api.Request, *api.Response], paths []string) error {
	configs, err := zos.ReadFiles(false, paths...)
	if err != nil {
		return ErrAppMainLoadConfigs.WithStack(err, nil)
	}

	for path, manifest := range configs {
		var format api.Format
		var unmarshalFunc func([]byte, any) error
		if strings.HasSuffix(path, ".json") {
			format = api.FormatJSON
			unmarshalFunc = encoder.UnmarshalJSON
		} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			format = api.FormatYAML
			unmarshalFunc = encoder.UnmarshalYAML
		} else {
			continue // Ignore other formats.
		}

		manifest = bytes.ReplaceAll(manifest, []byte("\r\n"), []byte("\n"))
		bs := SplitMultiDoc(manifest, "---\n")
		for _, b := range bs {
			into := &struct {
				APIVersion string `json:"apiVersion" yaml:"apiVersion"`
				Kind       string `json:"kind" yaml:"kind"`
			}{}

			if err := unmarshalFunc(b, into); err != nil {
				return ErrAppMainLoadConfigs.WithStack(err, map[string]any{"path": path})
			}

			if into.APIVersion == "" && into.Kind == "" {
				continue //Skip
			}

			req := &api.Request{
				Method:  api.MethodPost,
				Key:     into.APIVersion + "/" + into.Kind,
				Format:  format,
				Content: b,
			}

			if _, err := server.Serve(context.Background(), req); err != nil {
				return ErrAppMainLoadConfigs.WithStack(err, map[string]any{"path": path})
			}
		}
	}

	return nil
}

// ShowTemplate shows the resource template.
// This function panics when the given server is nil.
func ShowTemplate(server api.API[*api.Request, *api.Response], tpl, out string) {
	if tpl == "" {
		return
	}

	// tpl would be like "core/v1/Entrypoint" or "core/v1/Entrypoint/foo/bar"
	arr := strings.Split(strings.Trim(tpl, " "), "/")

	var ref *k.Reference
	if len(arr) == 3 {
		ref = &k.Reference{
			APIVersion: arr[0] + "/" + arr[1],
			Kind:       arr[2],
			Namespace:  "template",
			Name:       "template",
		}
	} else if len(arr) == 5 {
		ref = &k.Reference{
			APIVersion: arr[0] + "/" + arr[1],
			Kind:       arr[2],
			Namespace:  arr[3],
			Name:       arr[4],
		}
	} else {
		fmt.Println("invalid template format: " + tpl)
		fmt.Println("Should be \"apiGroup/apiVersion/kind\" or \"apiGroup/apiVersion/kind/namespace/name\"")
		Exit(2)
		return
	}

	var accept api.Format
	switch out { // Output format.
	case "json":
		accept = api.FormatJSON
	default:
		accept = api.FormatYAML
	}

	req := &api.Request{
		Method:  api.MethodGet,
		Key:     ref.APIVersion + "/" + ref.Kind,
		Format:  api.FormatProtoReference,
		Params:  map[string]string{api.KeyAccept: string(accept)},
		Content: ref,
	}

	res, err := server.Serve(context.Background(), req)
	if err != nil {
		fmt.Println(err.Error())
		Exit(2)
		return
	}

	fmt.Println(string(res.Content.([]byte)))
	Exit(0)
}

// SplitMultiDoc splits a documents in []byte format with a given separator.
// If an empty separator is given, the default separator "---\n" is used.
// Empty contents are ignored.
func SplitMultiDoc(in []byte, sep string) [][]byte {
	if sep == "" {
		sep = "---\n"
	}

	var outArr [][]byte
	inArr := bytes.Split(in, []byte(sep))
	for _, b := range inArr {
		// exclude empty documents.
		b = bytes.Trim(b, " \n")
		if len(b) != 0 {
			outArr = append(outArr, b)
		}
	}

	return outArr
}
