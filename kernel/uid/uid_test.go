// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package uid_test

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/uid"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func ExampleNewID() {
	id, err := uid.NewID()
	if err != nil {
		panic("handle error")
	}

	fmt.Println(hex.EncodeToString(id))                   // `^[0-9a-fA-F]{60}$`
	fmt.Println(base32.StdEncoding.EncodeToString(id))    // `^[2-7A-Z]{48}$`
	fmt.Println(base32.HexEncoding.EncodeToString(id))    // `^[0-9A-V]{48}$`
	fmt.Println(base64.StdEncoding.EncodeToString(id))    // `^[0-9a-zA-Z+/]{40}$`
	fmt.Println(base64.RawStdEncoding.EncodeToString(id)) // `^[0-9a-zA-Z+/]{40}$`
	fmt.Println(base64.URLEncoding.EncodeToString(id))    // `^[0-9a-zA-Z-_]{40}$`
	fmt.Println(base64.RawURLEncoding.EncodeToString(id)) // `^[0-9a-zA-Z-_]{40}$`
	// Example Output:
	// 0006221060e335a3fb86351a8a6366cedd480dbacf5c94af4d63392aa98f
	// AADCEEDA4M22H64GGUNIUY3GZ3OUQDN2Z5OJJL2NMM4SVKMP
	// 00324430SCQQ7US66KD8KOR6PREKG3DQPTE99BQDCCSILACF
	// AAYiEGDjNaP7hjUaimNmzt1IDbrPXJSvTWM5KqmP
	// AAYiEGDjNaP7hjUaimNmzt1IDbrPXJSvTWM5KqmP
	// AAYiEGDjNaP7hjUaimNmzt1IDbrPXJSvTWM5KqmP
	// AAYiEGDjNaP7hjUaimNmzt1IDbrPXJSvTWM5KqmP
}

func ExampleNewHostedID() {
	id, err := uid.NewHostedID()
	if err != nil {
		panic("handle error")
	}

	fmt.Println(hex.EncodeToString(id))                   // `^[0-9a-fA-F]{60}$`
	fmt.Println(base32.StdEncoding.EncodeToString(id))    // `^[2-7A-Z]{48}$`
	fmt.Println(base32.HexEncoding.EncodeToString(id))    // `^[0-9A-V]{48}$`
	fmt.Println(base64.StdEncoding.EncodeToString(id))    // `^[0-9a-zA-Z+/]{40}$`
	fmt.Println(base64.RawStdEncoding.EncodeToString(id)) // `^[0-9a-zA-Z+/]{40}$`
	fmt.Println(base64.URLEncoding.EncodeToString(id))    // `^[0-9a-zA-Z-_]{40}$`
	fmt.Println(base64.RawURLEncoding.EncodeToString(id)) // `^[0-9a-zA-Z-_]{40}$`
	// Example Output:
	// 00062210663d9b15c8869cf17b54f8ead187500e5e648154989224e33f35
	// AADCEEDGHWNRLSEGTTYXWVHY5LIYOUAOLZSICVEYSISOGPZV
	// 003244367MDHBI46JJONML7OTB8OEK0EBPI82L4OI8IE6FPL
	// AAYiEGY9mxXIhpzxe1T46tGHUA5eZIFUmJIk4z81
	// AAYiEGY9mxXIhpzxe1T46tGHUA5eZIFUmJIk4z81
	// AAYiEGY9mxXIhpzxe1T46tGHUA5eZIFUmJIk4z81
	// AAYiEGY9mxXIhpzxe1T46tGHUA5eZIFUmJIk4z81
}

type testErrorReader struct {
	io.Reader
	err error
}

func (r *testErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func TestNewID(t *testing.T) {
	type condition struct {
		encodeFunc     func([]byte) string
		typ            kernel.EncodingType
		useErrorReader bool
	}

	type action struct {
		err error
	}

	CndValidID := "expected id"
	CndErrorReader := "rand.Reader returns an error"
	ActCheckValidFormat := "non-empty expected string"
	ActCheckNoDuplication := "no duplication"
	ActCheckNoError := "check that there is no error returned"
	ActCheckExpectedError := "check that an expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndValidID, "an ID in a valid format is returned")
	tb.Condition(CndErrorReader, "rand.Reader returns an error")
	tb.Action(ActCheckValidFormat, "check that ID has a valid format")
	tb.Action(ActCheckNoDuplication, "check that there is no duplication among IDs generated in a short time")
	tb.Action(ActCheckNoError, "check that there is no error returned")
	tb.Action(ActCheckExpectedError, "check that an expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Base16",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base16,
				encodeFunc: hex.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base32",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base32,
				encodeFunc: base32.StdEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base32Hex",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base32Hex,
				encodeFunc: base32.HexEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64,
				encodeFunc: base64.StdEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64Raw",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64Raw,
				encodeFunc: base64.RawStdEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64URL",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64URL,
				encodeFunc: base64.URLEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64RawURL",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64RawURL,
				encodeFunc: base64.RawURLEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"error",
			[]string{CndErrorReader},
			[]string{ActCheckExpectedError},
			&condition{
				useErrorReader: true,
			},
			&action{
				err: &er.Error{Package: uid.ErrPkg, Type: "id", Description: uid.ErrDscNew},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Replace rand.Reader for testing.
			if tt.C().useErrorReader {
				tmp := rand.Reader
				rand.Reader = &testErrorReader{err: io.EOF}
				defer func() {
					rand.Reader = tmp
				}()
			}

			id, err := uid.NewID()
			testutil.DiffError(t, tt.A().err, nil, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			eid := string(tt.C().encodeFunc(id))
			pattern := Validator(tt.C().typ)
			testutil.Diff(t, true, pattern.MatchString(eid))
		})
	}
}

func TestNewHostedID(t *testing.T) {
	type condition struct {
		encodeFunc     func([]byte) string
		typ            kernel.EncodingType
		useErrorReader bool
	}

	type action struct {
		err error
	}

	CndValidID := "expected id"
	CndErrorReader := "rand.Reader returns an error"
	ActCheckValidFormat := "non-empty expected string"
	ActCheckNoDuplication := "no duplication"
	ActCheckNoError := "check that there is no error returned"
	ActCheckExpectedError := "check that an expected error returned"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndValidID, "an ID in a valid format is returned")
	tb.Condition(CndErrorReader, "rand.Reader returns an error")
	tb.Action(ActCheckValidFormat, "check that ID has a valid format")
	tb.Action(ActCheckNoDuplication, "check that there is no duplication among IDs generated in a short time")
	tb.Action(ActCheckNoError, "check that there is no error returned")
	tb.Action(ActCheckExpectedError, "check that an expected error returned")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Base16",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base16,
				encodeFunc: hex.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base32",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base32,
				encodeFunc: base32.StdEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base32Hex",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base32Hex,
				encodeFunc: base32.HexEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64,
				encodeFunc: base64.StdEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64Raw",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64Raw,
				encodeFunc: base64.RawStdEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64URL",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64URL,
				encodeFunc: base64.URLEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"Base64RawURL",
			[]string{CndValidID},
			[]string{ActCheckValidFormat, ActCheckNoDuplication, ActCheckNoError},
			&condition{
				typ:        kernel.EncodingType_Base64RawURL,
				encodeFunc: base64.RawURLEncoding.EncodeToString,
			},
			&action{},
		),
		gen(
			"error",
			[]string{CndErrorReader},
			[]string{ActCheckExpectedError},
			&condition{
				useErrorReader: true,
			},
			&action{
				err: &er.Error{Package: uid.ErrPkg, Type: "hosted id", Description: uid.ErrDscNew},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Replace rand.Reader for testing.
			if tt.C().useErrorReader {
				tmp := rand.Reader
				rand.Reader = &testErrorReader{err: io.EOF}
				defer func() {
					rand.Reader = tmp
				}()
			}

			id, err := uid.NewHostedID()
			testutil.DiffError(t, tt.A().err, nil, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			eid := string(tt.C().encodeFunc(id))
			pattern := Validator(tt.C().typ)
			testutil.Diff(t, true, pattern.MatchString(eid))
		})
	}
}

// Validator returns a regular expression for IDs.
// Returned expression can be used for validating IDs
// created by NewID() and NewHostedID().
// nil will be returned if an unknown EncodeType was given.
// Returned regular expressions are as follow.
//   - Base16           : ^[0-9a-fA-F]{60}$
//   - Base32           : ^[2-7A-Z]{48}$
//   - Base32Hex        : ^[0-9A-V]{48}$
//   - Base32Escaped    : ^[0-9B-DF-HJ-NP-TV-Z]{48}$
//   - Base32HexEscaped : ^[0-9B-DF-HJ-NP-TV-Z]{48}$
//   - Base64           : ^[0-9a-zA-Z+/]{40}$
//   - Base64Raw        : ^[0-9a-zA-Z+/]{40}$
//   - Base64URL        : ^[0-9a-zA-Z-_]{40}$
//   - Base64RawURL     : ^[0-9a-zA-Z-_]{40}$
func Validator(t kernel.EncodingType) *regexp.Regexp {
	switch t {
	case kernel.EncodingType_Base16:
		return regexp.MustCompile(`^[0-9a-fA-F]{60}$`)
	case kernel.EncodingType_Base32:
		return regexp.MustCompile(`^[2-7A-Z]{48}$`)
	case kernel.EncodingType_Base32Escaped:
		return regexp.MustCompile(`^[0-9B-DF-HJ-NP-TV-Z]{48}$`)
	case kernel.EncodingType_Base32Hex:
		return regexp.MustCompile(`^[0-9A-V]{48}$`)
	case kernel.EncodingType_Base32HexEscaped:
		return regexp.MustCompile(`^[0-9B-DF-HJ-NP-TV-Z]{48}$`)
	case kernel.EncodingType_Base64, kernel.EncodingType_Base64Raw:
		return regexp.MustCompile(`^[0-9a-zA-Z+/]{40}$`)
	case kernel.EncodingType_Base64URL, kernel.EncodingType_Base64RawURL:
		return regexp.MustCompile(`^[0-9a-zA-Z-_]{40}$`)
	default:
		return nil
	}
}
