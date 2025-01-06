package io

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// Setenv is the function to set environmental variable.
// This can be replaced when testing.
var Setenv func(key string, value string) error = os.Setenv

// LoadEnv load environmental variables from definition text.
// The second argument should be bytes data of environmental variable.
// If the first argument overwrite is set to true,
// the already defined variables are overwritten.
func LoadEnv(overwrite bool, bs ...[]byte) error {
	re := regexp.MustCompile(`(\w+)\s*=\s*(.*)`)
	for _, b := range bs {
		envs := map[string]string{}

		lines := bytes.Split(b, []byte("\n"))
		for _, line := range lines {
			line = bytes.Trim(line, "\t\n\f\r ")
			if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
				continue
			}
			kv := re.FindSubmatch(line)
			if len(kv) != 3 {
				return &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeEnv,
					Description: ErrDscLoadEnv,
					Detail:      "invalid line `" + string(line) + "`",
				}
			}
			envs[string(kv[1])] = string(kv[2])
		}

		for k, v := range envs {
			_, exists := os.LookupEnv(k)
			// Set variable when not defined or specified to overwrite it.
			if !exists || overwrite {
				// Set variable because it is not defined.
				if err := Setenv(k, v); err != nil {
					return (&er.Error{
						Package:     ErrPkg,
						Type:        ErrTypeEnv,
						Description: ErrDscSetEnv,
						Detail:      "`" + k + "=" + v + "`",
					}).Wrap(err)
				}
			}
		}
	}
	return nil
}

// ResolveEnv resolves environmental variable in a text.
// Expressions are basically derived from shell parameter substitution.
//   - https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
//   - https://tldp.org/LDP/abs/html/parameter-substitution.html
//
// allowed expressions are:
//   - ${FOO:-bar} : bar is used if FOO is empty or is not defined.
//   - ${FOO-bar} : bar is used if FOO is not defined.
//   - ${FOO+bar} : bar is used if FOO is defined.
//   - ${#FOO} : the length of FOO.
//
// and can use inner variable like:
//   - ${${INNER_ENV}}
//   - ${OUTER_ENV:-${INNER_ENV}}
//   - ${${INNER_ENV}:-default}
//   - ${OUTER_ENV-${INNER_ENV}}
//   - ${${INNER_ENV}-default}
//   - ${OUTER_ENV+${INNER_ENV}}
//   - ${${INNER_ENV}+default}
//   - ${#${INNER_ENV}}
func ResolveEnv(b []byte) []byte {
	exp := regexp.MustCompile(`\$\{[0-9a-zA-Z_:\-+#]+\}`) // Expression of environmental variables.
	result := exp.ReplaceAllFunc(b, resolve)              // The first ReplaceAllFunc resolves single variable and inner variables.
	result = exp.ReplaceAllFunc(result, resolve)          // The second ReplaceAllFunc resolves outer variables.
	return result
}

// resolve resolves and expands environmental variable expressions.
// The arguments must match the regular expressions of `\$\{[0-9a-zA-Z:\-+#]+\}`.
// If invalid format of the argument is given, this function returns the string as it is.
// If nil is passed as an argument, nil is returned.
//
// expressions:
//   - ${FOO:-bar} : bar is used if FOO is empty or is not defined.
//   - ${FOO-bar} : bar is used if FOO is not defined.
//   - ${FOO+bar} : bar is used if FOO is defined.
//   - ${#FOO} : the length of FOO.
func resolve(in []byte) []byte {
	target := strings.TrimPrefix(string(in), "$")
	target = strings.TrimPrefix(target, "{")
	target = strings.TrimSuffix(target, "}")

	// This assumes ${FOO:-bar}.
	// If FOO is defined and has at least 1 character,
	// this returns FOO value, otherwise bar.
	if i := strings.Index(target, ":-"); i != -1 {
		key := target[:i]
		def := target[i+2:]
		if val := os.Getenv(key); val != "" {
			return []byte(val)
		}
		return []byte(def)
	}

	// This assumes ${FOO-bar}.
	// If FOO is defined even it is empty, this returns FOO value, otherwise bar.
	if i := strings.Index(target, "-"); i != -1 {
		key := target[:i]
		def := target[i+1:]
		if val, defined := os.LookupEnv(key); defined {
			return []byte(val)
		}
		return []byte(def)
	}

	// This assumes ${FOO+bar}.
	// If FOO is defined, this returns bar value.
	// If FOO is not defined, this returns empty.
	if i := strings.Index(target, "+"); i != -1 {
		key := target[:i]
		val := target[i+1:]
		if _, defined := os.LookupEnv(key); defined {
			return []byte(val)
		}
		return []byte("")
	}

	// This assumes ${#FOO}.
	// Returns the length of FOO.
	if strings.HasPrefix(target, "#") {
		key := target[1:]
		val := os.Getenv(key)
		return []byte(strconv.Itoa(len([]rune(val))))
	}

	// This assumes normal expression like ${FOO}.
	// This returns empty if the variable is not defined.
	val := os.Getenv(target)
	return []byte(val)
}
