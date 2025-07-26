// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package slogger

import (
	"cmp"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/log"
	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/ztime/zcron"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	apiVersion = "core/v1"
	kind       = "SLogger"
	Key        = apiVersion + "/" + kind
)

var Resource api.Resource = &API{
	BaseResource: &api.BaseResource{
		DefaultProto: &v1.SLogger{
			APIVersion: apiVersion,
			Kind:       kind,
			Metadata: &kernel.Metadata{
				Namespace: "default",
				Name:      "default",
			},
			Spec: &v1.SLoggerSpec{
				OutputTimeFormat: time.DateTime,
				DateFormat:       "2006-01-02",
				TimeFormat:       "15:04:05.000",
				Level:            v1.LogLevel_Info,
				LogOutput: &v1.LogOutputSpec{
					OutputTarget: v1.OutputTarget_Stdout,
					RotateSize:   1024, // = 1024 MiB = 1 GiB
					LogFileName:  "application.log",
					TimeZone:     "Local",
				},
			},
		},
	},
}

type API struct {
	*api.BaseResource
}

func (*API) Create(_ api.API[*api.Request, *api.Response], msg protoreflect.ProtoMessage) (any, error) {
	c := msg.(*v1.SLogger)

	timeZone, err := time.LoadLocation(c.Spec.LogOutput.TimeZone)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}

	repl, err := newReplaceFunc(c.Spec.FieldReplacers)
	if err != nil {
		return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
	}
	repl.timeFmt = c.Spec.OutputTimeFormat
	repl.timeZone = timeZone

	var w io.Writer
	outSpec := c.Spec.LogOutput
	switch outSpec.OutputTarget {
	case v1.OutputTarget_Discard:
		w = io.Discard
	case v1.OutputTarget_Stdout:
		w = os.Stdout
	case v1.OutputTarget_Stderr:
		w = os.Stderr
	case v1.OutputTarget_File:
		w, err = newFileWriter(c.Spec.LogOutput)
		if err != nil {
			return nil, core.ErrCoreGenCreateObject.WithStack(err, map[string]any{"kind": kind})
		}
	default:
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:       log.LvToSLogLevel(Level(c.Spec.Level)),
		ReplaceAttr: repl.replaceAttr,
	}

	var slg *log.SLogger
	if c.Spec.Unstructured {
		slg = log.NewTextSLogger(w, opts)
	} else {
		slg = log.NewJSONSLogger(w, opts)
	}
	slg.DateFmt = c.Spec.DateFormat
	slg.TimeFmt = c.Spec.TimeFormat
	slg.NoLocation = c.Spec.NoLocation
	slg.NoDatetime = c.Spec.NoDatetime
	slg.Location = timeZone

	closer, _ := w.(io.Closer)
	return &finalizableLogger{
		Writer: w,
		Logger: slg,
		closer: closer,
	}, nil
}

// finalizableLogger is the logger that can be
// finalized on exit of the application.
// File writer with log rotations should be finalized
// before exist of the application.
// For example, log rotation and log compression should be done.
// This implements core.Finalizer interface.
type finalizableLogger struct {
	io.Writer
	log.Logger
	closer io.Closer
}

func (l *finalizableLogger) Finalize() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

func newFileWriter(spec *v1.LogOutputSpec) (*zlog.LogicalFile, error) {
	config := &zlog.LogicalFileConfig{
		Manager: &zlog.FileManagerConfig{
			MaxAge:        time.Duration(spec.MaxAge) * time.Second,
			MaxHistory:    int(spec.MaxBackup),
			MaxTotalBytes: 1024 * 1024 * int64(spec.MaxTotalSize),
			GzipLv:        int(spec.CompressLevel),
			SrcDir:        filepath.Clean(spec.LogDir),
			DstDir:        filepath.Clean(cmp.Or(spec.BackupDir, spec.LogDir)), // Default is SrcDir.
			Pattern:       cmp.Or(spec.ArchivedFilePattern, spec.LogFileName),
		},
		RotateBytes: 1024 * 1024 * int64(spec.RotateSize), // Max size of a single file.
		FileName:    spec.LogFileName,                     // Active file name.
	}
	lf, err := zlog.NewLogicalFile(config)
	if err != nil {
		return nil, err
	}

	if spec.Cron == "" {
		return lf, nil
	}
	c := &zcron.Config{
		Crontab: spec.Cron,
		JobFunc: func(ctx context.Context) error { return lf.Swap() },
	}
	cron, err := zcron.NewCron(c)
	if err != nil {
		return nil, err
	}
	go cron.Start()
	return lf, nil
}

// Level convert v1.LogLevel to logger.LogLevel.
//
// conversions:
//   - v1.LogLevel_Trace >- logger.LvTrace
//   - v1.LogLevel_Debug >- logger.LvDebug
//   - v1.LogLevel_Info >- logger.LvInfo
//   - v1.LogLevel_Warn >- logger.LvWarn
//   - v1.LogLevel_Error >- logger.LvError
//   - v1.LogLevel_Fatal >- logger.LvFatal
func Level(level v1.LogLevel) log.LogLevel {
	switch {
	case level <= v1.LogLevel_Trace:
		return log.LvTrace
	case level <= v1.LogLevel_Debug:
		return log.LvDebug
	case level <= v1.LogLevel_Info:
		return log.LvInfo
	case level <= v1.LogLevel_Warn:
		return log.LvWarn
	case level <= v1.LogLevel_Error:
		return log.LvError
	case level <= v1.LogLevel_Fatal:
		return log.LvFatal
	default:
		return log.LvFatal
	}
}
