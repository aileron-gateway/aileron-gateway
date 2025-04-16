// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package log

import (
	"cmp"
	"strings"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
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

// Level convert v1.LogLevel to logger.LogLevel.
//
// conversions:
//   - v1.LogLevel_Trace >- logger.LvTrace
//   - v1.LogLevel_Debug >- logger.LvDebug
//   - v1.LogLevel_Info >- logger.LvInfo
//   - v1.LogLevel_Warn >- logger.LvWarn
//   - v1.LogLevel_Error >- logger.LvError
//   - v1.LogLevel_Fatal >- logger.LvFatal
func Level(level k.LogLevel) LogLevel {
	switch {
	case level <= k.LogLevel_Trace:
		return LvTrace
	case level <= k.LogLevel_Debug:
		return LvDebug
	case level <= k.LogLevel_Info:
		return LvInfo
	case level <= k.LogLevel_Warn:
		return LvWarn
	case level <= k.LogLevel_Error:
		return LvError
	case level <= k.LogLevel_Fatal:
		return LvFatal
	default:
		return LvFatal
	}
}
