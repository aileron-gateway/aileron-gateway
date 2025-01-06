package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// LvToSLogLevel convert LogLevel to slog.Level.
// Because slog.Level does not have the level equivalent to Fatal and Trace,
// they are converted to other levels.
//
// conversions:
//   - LvTrace -> slog.LevelDebug
//   - LvDebug -> slog.LevelDebug
//   - LvInfo  -> slog.LevelInfo
//   - LvWarn  -> slog.LevelWarn
//   - LvError -> slog.LevelError
//   - LvFatal -> slog.LevelError
func LvToSLogLevel(level LogLevel) slog.Level {
	switch {
	case level <= LvTrace:
		return slog.LevelDebug
	case level <= LvDebug:
		return slog.LevelDebug
	case level <= LvInfo:
		return slog.LevelInfo
	case level <= LvWarn:
		return slog.LevelWarn
	case level <= LvError:
		return slog.LevelError
	case level <= LvFatal:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

// LvFromSLogLevel convert slog.Level to LogLevel.
// Because slog.Level does not have the level equivalent to Fatal and Trace,
// they are converted to other levels.
//
// conversions:
//   - slog.LevelDebug -> LvDebug
//   - slog.LevelInfo -> LvInfo
//   - slog.LevelWarn -> LvWarn
//   - slog.LevelError -> LvError
func LvFromSLogLevel(level slog.Level) LogLevel {
	switch {
	case level <= slog.LevelDebug:
		return LvDebug
	case level <= slog.LevelInfo:
		return LvInfo
	case level <= slog.LevelWarn:
		return LvWarn
	case level <= slog.LevelError:
		return LvError
	default:
		return LvError
	}
}

// NewTextSLogger returns a new Logger created with slog.NewTextHandler.
// os.Stdio is used if the given io.Writer is nil.
// Zero value of slog.HandlerOptions is used when nil option was given
// which means the default log level becomes INFO.
func NewTextSLogger(w io.Writer, opts *slog.HandlerOptions) *SLogger {
	if w == nil {
		w = os.Stdout
	}
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo // 0
	}
	return &SLogger{
		w:        w,
		lg:       slog.New(slog.NewTextHandler(w, opts)),
		lv:       LvFromSLogLevel(opts.Level.Level()),
		DateFmt:  "2006-01-02",
		TimeFmt:  "15:04:05.000",
		Location: time.Local,
	}
}

// NewJSONSLogger returns a new Logger created with slog.NewJSONHandler.
// os.Stdio is used if the given io.Writer is nil.
// Zero value of slog.HandlerOptions is used when nil option was given
// which means the default log level becomes INFO.
func NewJSONSLogger(w io.Writer, opts *slog.HandlerOptions) *SLogger {
	if w == nil {
		w = os.Stdout
	}
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo // 0
	}
	return &SLogger{
		w:        w,
		lg:       slog.New(slog.NewJSONHandler(w, opts)),
		lv:       LvFromSLogLevel(opts.Level.Level()),
		DateFmt:  "2006-01-02",
		TimeFmt:  "15:04:05.000",
		Location: time.Local,
	}
}

// SLogger is a logger using slog.Logger.
// This implements logger.Logger interface.
type SLogger struct {
	w  io.Writer
	lg *slog.Logger
	lv LogLevel

	NoLocation bool
	NoDatetime bool
	DateFmt    string
	TimeFmt    string
	Location   *time.Location
}

func (l *SLogger) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}

func (l *SLogger) Enabled(level LogLevel) bool {
	return l.lv <= level
}

func (l *SLogger) Debug(ctx context.Context, msg string, v ...any) {
	if l.lv > LvDebug {
		return
	}
	l.lg.DebugContext(ctx, msg, append(l.newAttrs(ctx), v...)...)
}

func (l *SLogger) Info(ctx context.Context, msg string, v ...any) {
	if l.lv > LvInfo {
		return
	}
	l.lg.InfoContext(ctx, msg, append(l.newAttrs(ctx), v...)...)
}

func (l *SLogger) Warn(ctx context.Context, msg string, v ...any) {
	if l.lv > LvWarn {
		return
	}
	l.lg.WarnContext(ctx, msg, append(l.newAttrs(ctx), v...)...)
}

func (l *SLogger) Error(ctx context.Context, msg string, v ...any) {
	if l.lv > LvError {
		return
	}
	l.lg.ErrorContext(ctx, msg, append(l.newAttrs(ctx), v...)...)
}

func (l *SLogger) newAttrs(ctx context.Context) []any {
	attrs := make([]any, 0, 2)
	if !l.NoDatetime {
		attrs = append(attrs, keyDatetime, NewDatetimeAttrs(l.DateFmt, l.TimeFmt, l.Location).Map())
	}
	if !l.NoLocation {
		attrs = append(attrs, keyLocation, NewLocationAttrs(3).Map())
	}
	a := AttrsFromContext(ctx)
	for i := range a {
		attrs = append(attrs, a[i].Name(), a[i].Map())
	}
	return attrs
}
