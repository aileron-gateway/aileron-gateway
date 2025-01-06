package slogger

import (
	"log/slog"
	"strings"
	"time"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/core/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

// newReplaceFunc return a new instance of replacer.
// When the returned error is not nil, returned replacer is nil.
// If the nil error was returned, the returned replacer is always non-nil.
// As shown below, set timeFmt and timeZone after obtained a replacer.
//
//	repl, err := newReplaceFunc(c.Spec.FieldReplacers)
//	if err != nil {
//		// Handle error.
//	}
//	repl.timeFmt = "2006-01-02 15:04:05" // Set non empty time format.
//	repl.timeZone = time.Local // Set non-nil location. Otherwise panics.
func newReplaceFunc(specs []*v1.FieldReplacerSpec) (*replacer, error) {
	r := &replacer{
		timeFmt:  time.DateTime,
		timeZone: time.Local,
	}
	for _, spec := range specs {
		if spec.Field == "" {
			continue
		}
		var repl txtutil.ReplaceFunc[string]
		if spec.Replacer != nil {
			r, err := txtutil.NewStringReplacer(spec.Replacer)
			if err != nil {
				return nil, err
			}
			repl = r.Replace
		}
		paths := strings.Split(spec.Field, ".")
		f := replaceAttrFunc(paths[1:], repl)
		r.keys = append(r.keys, paths[0])
		r.repl = append(r.repl, f)
	}
	return r, nil
}

// replacer is the attribute replacer for slogger.
// replacer.replaceAttr is intended to be used in
// slog.HandlerOptions.ReplaceAttr.
// Use newReplaceFunc to create replacer instances.
//
//	opts := &slog.HandlerOptions{
//		ReplaceAttr: repl.replaceAttr, // Use replaceAttr here.
//	}
type replacer struct {
	timeFmt  string
	timeZone *time.Location
	keys     []string
	repl     []func(slog.Attr) (slog.Attr, bool)
}

func (r *replacer) replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		a.Value = slog.StringValue(time.Now().In(r.timeZone).Format(r.timeFmt))
		return a
	}
	for i := 0; i < len(r.keys); i++ {
		if a.Key == r.keys[i] {
			if aa, ok := r.repl[i](a); ok {
				return aa
			}
		}
	}
	return a
}

func replaceAttrFunc(keys []string, replFunc txtutil.ReplaceFunc[string]) func(slog.Attr) (slog.Attr, bool) {
	if len(keys) == 0 {
		if replFunc == nil {
			noAttr := slog.Attr{}
			return func(_ slog.Attr) (slog.Attr, bool) {
				return noAttr, true
			}
		} else {
			return func(a slog.Attr) (slog.Attr, bool) {
				v, ok := a.Value.Any().(string)
				if ok {
					a.Value = slog.StringValue(replFunc(v))
				}
				return a, ok
			}
		}
	}

	lastKey := keys[len(keys)-1]
	return func(a slog.Attr) (slog.Attr, bool) {
		v := a.Value.Any()
		var vv map[string]any
		var ok bool
		for i := 0; i < len(keys); i++ {
			vv, ok = v.(map[string]any)
			if !ok {
				return a, false
			}
			v, ok = vv[keys[i]]
			if !ok {
				return a, false
			}
		}

		if replFunc == nil {
			delete(vv, lastKey)
			return a, true
		} else {
			if s, ok := v.(string); ok {
				vv[lastKey] = replFunc(s)
				return a, true
			}
		}

		return a, false
	}
}
