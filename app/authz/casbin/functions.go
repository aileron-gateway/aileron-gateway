package casbin

import (
	"errors"
	"net/http"
	"net/url"
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errKeyNotString = errors.New("authz/casbin: key is not string")
	errKeyNotFound  = errors.New("authz/casbin: key does not exists")

	errAuthValueInvalidArgLength = errors.New("authz/casbin: authValue requires at least 2 arguments")
	errAuthValueInvalidType      = errors.New("authz/casbin: authValue requires jwt.MapClaims or map[string]any")

	errQueryValueInvalidArgLength = errors.New("authz/casbin: queryValue or queryValues requires exactly 2 arguments")
	errQueryValueInvalidType      = errors.New("authz/casbin: queryValue or queryValues requires url.Values or map[string][]string")

	errHeaderValueInvalidArgLength = errors.New("authz/casbin: headerValue or queryValues requires exactly 2 arguments")
	errHeaderValueInvalidType      = errors.New("authz/casbin: headerValue or headerValues requires http.Header or map[string][]string")

	errContainsInvalidArgLength = errors.New("authz/casbin: contains functions requires exactly 2 arguments")
	errContainsInvalidSlice     = errors.New("authz/casbin: contains function got invalid type slice at 1st argument")
	errContainsInvalidValue     = errors.New("authz/casbin: contains function got invalid type data at 2nd argument")

	errAsSliceInvalidType = errors.New("authz/casbin: asSlice function got invalid type data")
)

// mapValue returns a value in a map.
func mapValue(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errAuthValueInvalidArgLength
	}

	// The first argument must be jwt.MapClaims or map[string]any.
	// Use 'r.Sub' for specific.
	sub, ok := args[0].(jwt.MapClaims)
	if !ok {
		sub, ok = args[0].(map[string]any)
		if !ok {
			return nil, errAuthValueInvalidType
		}
	}

	tmp := sub
	var value any
	for _, key := range args[1:] {
		// Keys (after the second argument) must be a string.
		key, ok := key.(string)
		if !ok {
			return nil, errKeyNotString
		}

		// Error if map has no value.
		if value, ok = tmp[key]; !ok {
			return nil, errKeyNotFound
		}

		// If the value is a map (map[string]any), update tmp and continue loop.
		// If it is not a map, the loop ends.
		if m, ok := value.(map[string]any); ok {
			tmp = m
		} else {
			break
		}
	}

	// Casbin does not support int type.
	// If the value type is int, the value is casted to float.
	if v, ok := toFloat64(value); ok {
		return v, nil
	}

	return value, nil
}

// toFloat64 converts all integer and float to float64 except for all unsigned int.
// This is because the casbin use all numeric values as float64.
func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}

// queryValue returns a query value of given key.
// The type of the returned value is (string, error)
// The given argument must be (url.Value, string) or (map[string][]string, string).
// Only 1 value is returned even multiple values were found.
func queryValue(args ...any) (any, error) {
	vs, err := queryValues(args...)
	v, _ := vs.([]string)
	if len(v) > 0 {
		return v[0], err // string, error
	}

	return "", err // string, error
}

// queryValues returns query values of given key.
// The type of the returned value is ([]string, error)
// The given argument must be (url.Value, string) or (map[string][]string, string).
func queryValues(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errQueryValueInvalidArgLength
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, errKeyNotString
	}

	// The first argument must be url.Values or map[string][]string.
	// Use 'r.Sub.Query' for specific.

	sub, ok := args[0].(url.Values)
	if ok {
		return sub[key], nil // string, error
	}

	sub, ok = args[0].(map[string][]string)
	if ok {
		return sub[key], nil // string, error
	}

	return nil, errQueryValueInvalidType
}

// headerValue returns a header value of given key.
// The type of the returned value is (string, error)
// The given argument must be (http.Header, string) or (map[string][]string, string).
// Only 1 value is returned even multiple values were found.
func headerValue(args ...any) (any, error) {
	vs, err := headerValues(args...)
	v, _ := vs.([]string)
	if len(v) > 0 {
		return v[0], err // string, error
	}

	return "", err // string, error
}

// headerValues returns header values of given key.
// The type of the returned value is ([]string, error)
// The given argument must be (http.Header, string) or (map[string][]string, string).
func headerValues(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errHeaderValueInvalidArgLength
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, errKeyNotString
	}

	// The first argument must be http.Header or map[string][]string.
	// Use 'r.Sub.Header' for specific.

	sub, ok := args[0].(http.Header)
	if ok {
		// CanonicalMIMEHeaderKey is used.
		return sub.Values(key), nil // []string, error
	}

	sub, ok = args[0].(map[string][]string)
	if ok {
		return sub[key], nil // []string, error
	}

	return nil, errHeaderValueInvalidType
}

func contains[T comparable](args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errContainsInvalidArgLength
	}

	// The first argument must be a slice.
	list, ok := args[0].([]T)
	if !ok {
		return nil, errContainsInvalidSlice
	}

	var value T
	switch v := args[1].(type) {
	case T:
		value = v
	default:
		return nil, errContainsInvalidValue
	}

	return slices.Contains(list, value), nil
}

func containsNumber[T Integer | Float](args ...any) (any, error) {
	if len(args) != 2 {
		return nil, errContainsInvalidArgLength
	}

	// The first argument must be a slice.
	list, ok := args[0].([]T)
	if !ok {
		return nil, errContainsInvalidSlice
	}

	// The second argument must be a value.
	var value T
	switch v := args[1].(type) {
	case T:
		value = v
	case float64:
		value = T(v)
	default:
		return nil, errContainsInvalidValue
	}

	return slices.Contains(list, value), nil
}

func asSlice[T any](args ...any) (any, error) {
	values := make([]T, 0, len(args))

	for i := range args {
		v, ok := args[i].(T)
		if !ok {
			return nil, errAsSliceInvalidType
		}
		values = append(values, v)
	}

	return values, nil
}

func asSliceNumber[T Integer | Float](args ...any) (any, error) {
	values := make([]T, 0, len(args))

	for i := range args {
		switch v := args[i].(type) {
		case T:
			values = append(values, v)
		case float64:
			values = append(values, T(v))
		default:
			return nil, errAsSliceInvalidType
		}
	}

	return values, nil
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}
