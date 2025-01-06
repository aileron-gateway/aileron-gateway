package testutil

import (
	"errors"
	"io"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Diff compares two value using go-comp.
func Diff(t *testing.T, want, got any, opts ...cmp.Option) {
	t.Helper()
	// opts = append(opts, cmpopts.EquateEmpty())
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
		t.Errorf("StackTrace:\n%s", debug.Stack())
	}
}

func Equal(t *testing.T, want, got any, opts ...cmp.Option) {
	t.Helper()
	// opts = append(opts, cmpopts.EquateEmpty())
	if ok := cmp.Equal(want, got, opts...); !ok {
		t.Errorf("mismatch :\n")
		t.Errorf("StackTrace:\n%s", debug.Stack())
	}
}

// Iser is the interface which check if the given error is
// the same as this error instance.
// See https://pkg.go.dev/errors#Is
type Iser interface {
	Is(error) bool
}

// DiffError compares two error.
func DiffError(t *testing.T, want any, pattern *regexp.Regexp, got error, opts ...cmp.Option) {
	t.Helper()

	if want == nil || got == nil {
		Diff(t, want, got, opts...)
		return
	}

	if pattern != nil && !pattern.MatchString(got.Error()) {
		t.Errorf("error message mismatch :\n")
		t.Errorf(" (-want) %s\n", pattern.String())
		t.Errorf(" (+got)  %s\n", got.Error())
	}

	iser, ok := want.(Iser)
	if !ok {
		Diff(t, want, got, opts...)
		return
	}

	if !iser.Is(got) {
		t.Errorf("error mismatch :\n")
		t.Errorf(" (-want) %#v\n", want)
		t.Errorf(" (+got)  %#v\n", got)
	}
}

// ComparePointer is the function to compare two value pointers.
// This function is mainly intended to be used as an option of go-cmp.
// For example, use this option
//
//	cmp.Comparer(testutil.ComparePointer[foo.Bar])
func ComparePointer[T any](x, y T) bool {
	return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
}

// ErrorReader is an io.Reader which returns an error.
// This implements io.Reader interface.
type ErrorReader struct {
	io.Reader
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("rand read error")
}

// DeepAllowUnexported returns compare options
// like reflect.DeepEqual.
// See https://github.com/google/go-cmp/issues/40
func DeepAllowUnexported(vs ...any) cmp.Option {
	m := make(map[reflect.Type]struct{})
	for _, v := range vs {
		structTypes(reflect.ValueOf(v), m)
	}
	typs := make([]any, 0, len(m))
	for t := range m {
		typs = append(typs, reflect.New(t).Elem().Interface())
	}
	return cmp.AllowUnexported(typs...)
}

func structTypes(v reflect.Value, m map[reflect.Type]struct{}) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			structTypes(v.Elem(), m)
		}
	case reflect.Interface:
		if !v.IsNil() {
			structTypes(v.Elem(), m)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			structTypes(v.Index(i), m)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			structTypes(v.MapIndex(k), m)
		}
	case reflect.Struct:
		m[v.Type()] = struct{}{}
		for i := 0; i < v.NumField(); i++ {
			structTypes(v.Field(i), m)
		}
	}
}
