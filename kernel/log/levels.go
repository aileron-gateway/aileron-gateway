// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"cmp"
	"strings"
)

type LogLevel int

// Log levels are defined based on the severity of OpenTelemetry.
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/data-model.md
const (
	LvTrace LogLevel = 4       //  -4
	LvDebug LogLevel = 8       // 5-8
	LvInfo  LogLevel = 12      // 9-12
	LvWarn  LogLevel = 16      // 13-16
	LvError LogLevel = 20      // 17-20
	LvFatal LogLevel = 24      // 21-
	Trace   string   = "TRACE" //  -4
	Debug   string   = "DEBUG" // 5-8
	Info    string   = "INFO"  // 9-12
	Warn    string   = "WARN"  // 13-16
	Error   string   = "ERROR" // 17-20
	Fatal   string   = "FATAL" // 21-
)

// LevelFromText return log level from string.
// LevelFromText return LvInfo if the unknown
// level string was given.
//
// conversions:
//   - "TRACE" -> LvTrace ( = 4 )
//   - "DEBUG" -> LvDebug ( = 8 )
//   - "INFO"  -> LvInfo  ( = 12 )
//   - "WARN"  -> LvWarn  ( = 16 )
//   - "ERROR" -> LvError ( = 20 )
//   - "FATAL" -> LvFatal ( = 24 )
func LevelFromText(lv string) LogLevel {
	lv = strings.ToUpper(lv)
	return cmp.Or(map[string]LogLevel{
		Trace: LvTrace,
		Debug: LvDebug,
		Info:  LvInfo,
		Warn:  LvWarn,
		Error: LvError,
		Fatal: LvFatal,
	}[lv], LvInfo)
}
