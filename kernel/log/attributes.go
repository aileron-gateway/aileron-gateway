package log

import (
	"runtime"
	"strings"
	"time"
)

const (
	keyFile = "file"
	keyLine = "line"
	keyFunc = "func"

	keyDate = "date"
	keyTime = "time"
	keyZone = "zone"

	keyDatetime = "datetime"
	keyLocation = "location"
)

// NewLocationAttrs returns a new instance of location log attributes.
func NewLocationAttrs(skip int) *LocationAttrs {
	fn := ""
	ptr, file, line, ok := runtime.Caller(skip)
	if ok {
		rpc := make([]uintptr, 1)
		rpc[0] = ptr
		f, _ := runtime.CallersFrames(rpc).Next()
		if i := strings.LastIndexByte(f.Function, '/'); i > 0 {
			fn = f.Function[i+1:]
		}
	}

	// Not to show too long path in the logs, trim some prefixes.
	if i := strings.LastIndexByte(file, '/'); i > 0 {
		if j := strings.LastIndexByte(file[:i], '/'); j > 0 {
			file = file[j+1:]
		}
	}

	return &LocationAttrs{
		name: keyLocation,
		file: file,
		line: line,
		fn:   fn,
	}
}

// LocationAttrs is the log attributes of caller's location.
// This implements log.Attributes interface.
type LocationAttrs struct {
	name string
	file string
	line int
	fn   string
}

func (a *LocationAttrs) Name() string {
	return a.name
}

func (a *LocationAttrs) Map() map[string]any {
	m := make(map[string]any, 3)
	m[keyFile] = a.file
	m[keyLine] = a.line
	m[keyFunc] = a.fn
	return m
}

func (a *LocationAttrs) KeyValues() []any {
	return []any{
		keyFile, a.file,
		keyLine, a.line,
		keyFunc, a.fn,
	}
}

// NewDatetimeAttrs returns a new datetime log attributes.
// Locale is parsed using time.LoadLocation.
// Local will be used if an invalid locale was given.
// See https://pkg.go.dev/time#LoadLocation
//
//	dfmt = "2006-01-02"
//	tfmt = "15:04:05"
//	loc  = "UTC"
func NewDatetimeAttrs(dfmt, tfmt string, loc *time.Location) *DatetimeAttrs {
	if loc == nil {
		loc = time.Local
	}
	t := time.Now().In(loc)
	return &DatetimeAttrs{
		name: keyDatetime,
		date: t.Format(dfmt),
		time: t.Format(tfmt),
		zone: loc.String(),
	}
}

// DatetimeAttrs is the log attributes of date and time.
// This implements log.Attributes interface.
type DatetimeAttrs struct {
	name string
	date string
	time string
	zone string
}

func (a *DatetimeAttrs) Name() string {
	return a.name
}

func (a *DatetimeAttrs) Map() map[string]any {
	m := make(map[string]any, 3)
	m[keyDate] = a.date
	m[keyTime] = a.time
	m[keyZone] = a.zone
	return m
}

func (a *DatetimeAttrs) KeyValues() []any {
	return []any{
		keyDate, a.date,
		keyTime, a.time,
		keyZone, a.zone,
	}
}

// NewCustomAttrs returns a new instance of custom attributes.
func NewCustomAttrs(name string, m map[string]any) *CustomAttrs {
	if m == nil {
		m = map[string]any{}
	}
	return &CustomAttrs{
		name: name,
		m:    m,
	}
}

// CustomAttrs is the customizable log attributes.
// This is the implements of log.Attributes interface.
type CustomAttrs struct {
	name string
	m    map[string]any
}

func (a *CustomAttrs) Name() string {
	return a.name
}

func (a *CustomAttrs) Map() map[string]any {
	return a.m
}

func (a *CustomAttrs) KeyValues() []any {
	arr := make([]any, 0, 2*len(a.m))
	for k, v := range a.m {
		arr = append(arr, k, v)
	}
	return arr
}
