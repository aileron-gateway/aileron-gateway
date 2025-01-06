package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// DefaultLoggerName is the name of the logger that is registered
// in the global logger holder by default.
//
// The default logger can be obtained like below.
// logger.NopLogger is returned by default.
// The logger obtained by the name of logger.DefaultLoggerName will never be nil.
//
//	lg := logger.GlobalLogger(logger.DefaultLoggerName)
//
// The default logger can be replaced by
//
//	var lg logger.Logger
//	lg = <logger you want to use>
//	logger.SetGlobalLogger(logger.DefaultLoggerName, lg)
const DefaultLoggerName = "__default__"

// NoopLoggerName is the noop logger name.
// Default noop logger is available by
// GlobalLogger func.
const NoopLoggerName = "__noop__"

var (
	// mu protects loggers.
	mu = sync.RWMutex{}
	// loggers is the global logger set.
	loggers = map[string]Logger{
		DefaultLoggerName: &SLogger{
			w: os.Stdout,
			lg: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:       LvToSLogLevel(LevelFromText(os.Getenv("GATEWAY_LOG_LEVEL"))),
				ReplaceAttr: replaceTime,
			})),
			lv:         LevelFromText(os.Getenv("GATEWAY_LOG_LEVEL")),
			DateFmt:    "2006-01-02",
			TimeFmt:    "15:04:05.000",
			Location:   time.Local,
			NoLocation: os.Getenv("GATEWAY_LOG_LOCATION") == "false",
			NoDatetime: os.Getenv("GATEWAY_LOG_DATETIME") == "false",
		},
		NoopLoggerName: NoopLogger,
	}
)

// replaceTime replaces time format of slogger.
// This function is used only in the default logger.
// See the default logger above.
// This function is used at slog.HandlerOptions.ReplaceAttr.
func replaceTime(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key == slog.TimeKey {
		attr.Value = slog.StringValue(time.Now().Format(time.DateTime)) // Local zone.
	}
	return attr
}

// DefaultOr returns a globally used logger.
// DefaultOr returns default logger if a logger
// with given name does not exist.
func DefaultOr(name string) Logger {
	mu.RLock()
	defer mu.RUnlock()
	logger, ok := loggers[name]
	if !ok {
		return loggers[DefaultLoggerName]
	}
	return logger
}

// GlobalLogger returns a logger which is stored in the global log holder by name.
// A slogger is registered by default with the name of logger.DefaultLoggerName.
// If there is no logger, this function returns nil.
//
// When getting the default logger, use like below.
// The logger gotten by the name DefaultLoggerName won't be nil.
//
//	lg := logger.GlobalLogger(logger.DefaultLoggerName)
//
// If it's not the default logger, nil check should be taken.
//
//	lg := logger.GlobalLogger("yourLoggerName")
//	if lg == nil {
//		// use other logger.
//	}
func GlobalLogger(name string) Logger {
	mu.RLock()
	defer mu.RUnlock()
	return loggers[name]
}

// SetGlobalLogger stores the given logger in the global log holder.
// This replaces the existing logger if there have already been the same named logger.
// To delete the logger, set nil as the second argument.
// The logger named logger.DefaultLoggerName can be replaced but cannot be deleted.
//
// To delete logger:
//
//	logger.SetGlobalLogger("loggerName", nil)
//
// To replace default logger:
//
//	lg = <logger you want to use>
//	logger.SetGlobalLogger(logger.DefaultLoggerName, lg)
func SetGlobalLogger(name string, logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	if logger == nil {
		if name != DefaultLoggerName && name != NoopLoggerName {
			delete(loggers, name)
		}
		return
	}
	loggers[name] = logger
}

// NoopLogger is a no-operation logger which do nothing.
var NoopLogger Logger = &Noop{
	Writer: io.Discard,
}

// Noop is the logger that do nothing.
// This implements logger.Logger interface.
type Noop struct {
	io.Writer
}

func (l *Noop) Enabled(_ LogLevel) bool {
	return false
}

func (l *Noop) Debug(_ context.Context, _ string, _ ...any) {}

func (l *Noop) Info(_ context.Context, _ string, _ ...any) {}

func (l *Noop) Warn(_ context.Context, _ string, _ ...any) {}

func (l *Noop) Error(_ context.Context, _ string, _ ...any) {}
