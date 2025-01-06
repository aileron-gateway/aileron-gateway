package app

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/pflag"
)

// Version is the version of the binary.
var Version = "UNSET"

func ParseArgs(args []string, custom ...*pflag.FlagSet) *Options {
	opts := &Options{
		Metadata: &MetadataOptions{},
		Basic:    &BasicOptions{},
	}

	root := pflag.NewFlagSet("root", pflag.ContinueOnError)
	root.AddFlagSet(opts.Metadata.FlagSet())
	root.AddFlagSet(opts.Basic.FlagSet())
	for _, c := range custom {
		root.AddFlagSet(c) // Add custom flags.
	}

	// Parse arguments.
	if err := root.Parse(args); err != nil {
		fmt.Println(err.Error())
		fmt.Println("")
		fmt.Println("Options :")
		fmt.Println(root.FlagUsages())
		Exit(2)
		return nil
	}

	// If non-flag arguments are remained, show them and exit.
	if len(root.Args()) > 0 {
		fmt.Println("invalid arguments: ", root.Args())
		fmt.Println("")
		fmt.Println("Options :")
		fmt.Println(root.FlagUsages())
		Exit(2)
		return nil
	}

	if opts.Metadata.Version {
		fmt.Println(Version)
		Exit(0)
		return nil
	}

	if opts.Metadata.BuildInfo {
		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(info.String())
		}
		Exit(0)
		return nil
	}

	if opts.Metadata.Help {
		fmt.Println("Options :")
		fmt.Println(root.FlagUsages())
		Exit(0)
		return nil
	}

	return opts
}

type Options struct {
	Metadata *MetadataOptions
	Basic    *BasicOptions
}

type MetadataOptions struct {
	Version   bool
	BuildInfo bool
	Help      bool
}

func (o *MetadataOptions) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("metadata", pflag.ContinueOnError)
	fs.BoolVarP(&o.Version, "version", "v", false, "show version")
	fs.BoolVarP(&o.BuildInfo, "info", "i", false, "show build information")
	fs.BoolVarP(&o.Help, "help", "h", false, "show help message")
	return fs
}

type BasicOptions struct {
	Configs  []string
	Envs     []string
	Template string
	Out      string
}

func (o *BasicOptions) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("basic", pflag.ContinueOnError)
	fs.StringArrayVarP(&o.Configs, "file", "f", []string{}, "config file or directory path. absolute or relative")
	fs.StringArrayVarP(&o.Envs, "env", "e", []string{}, "env file path. each line be 'KEY=VALUE'")
	fs.StringVarP(&o.Template, "template", "t", "", "show template config. value format be 'Group/Version/Kind(/Namespace/Name)'")
	fs.StringVarP(&o.Out, "out", "o", "yaml", "template output format. yaml or json")
	return fs
}
