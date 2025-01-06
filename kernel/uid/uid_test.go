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

func TestIDValidator(t *testing.T) {
	type condition struct {
		typ kernel.EncodingType
	}

	type action struct {
		shouldNil bool
		valid     []string
		invalid   []string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndValidType := tb.Condition("valid type", "input a valid encode type")
	cndUnknown := tb.Condition("unknown type", "input a invalid encode type")
	actCheckType := tb.Action("check type", "check the returned encode type")
	table := tb.Build()

	must := func(v []byte, err error) []byte {
		if err != nil {
			panic(err)
		}
		return v
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Base16",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base16,
			},
			&action{
				valid: []string{
					hex.EncodeToString(must(uid.NewID())),
					hex.EncodeToString(must(uid.NewHostedID())),
					"1234567890abcdefABCDEF0123456789abcdefABCDEF0123456789abcdef", // 60 chars.
				},
				invalid: []string{
					"12345678901234567890123456789012345678901234567890123456789",   // 59 chars.
					"1234567890123456789012345678901234567890123456789012345678901", // 61 chars.
					"12345678901234567890123456789012345678901234567890123456789$",  // 60 chars with invalid char.
				},
			},
		),
		gen(
			"Base32",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base32,
			},
			&action{
				valid: []string{
					base32.StdEncoding.EncodeToString(must(uid.NewID())),
					base32.StdEncoding.EncodeToString(must(uid.NewHostedID())),
					"234567ABCDEFGHIJKLMNOPQRSTUVWXYZ234567ABCDEFGHIJ", // 48 chars.
				},
				invalid: []string{
					"234567ABCDEFGHIJKLMNOPQRSTUVWXYZ234567ABCDEFGHI",   // 47 chars.
					"234567ABCDEFGHIJKLMNOPQRSTUVWXYZ234567ABCDEFGHIJK", // 49 chars.
					"234567ABCDEFGHIJKLMNOPQRSTUVWXYZ234567ABCDEFGHI$",  // 48 chars with invalid char.
				},
			},
		),
		gen(
			"Base32Hex",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base32Hex,
			},
			&action{
				valid: []string{
					base32.HexEncoding.EncodeToString(must(uid.NewID())),
					base32.HexEncoding.EncodeToString(must(uid.NewHostedID())),
					"1234567890ABCDEFGHIJKLMNOPQRSTUV1234567890ABCDEF", // 48 chars.
				},
				invalid: []string{
					"1234567890ABCDEFGHIJKLMNOPQRSTUV1234567890ABCDE",   // 47 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUV1234567890ABCDEFG", // 49 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUV1234567890ABCDE$",  // 48 chars with invalid char.
				},
			},
		),
		gen(
			"Base32Escaped",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base32Escaped,
			},
			&action{
				valid: []string{
					"1234567890BCDFGHJKLMNPQRSTUVWXYZ1234567890BCDFGH", // 48 chars.
				},
				invalid: []string{
					"1234567890BCDFGHJKLMNPQRSTUVWXYZ1234567890BCDFG",   // 47 chars.
					"1234567890BCDFGHJKLMNPQRSTUVWXYZ1234567890BCDFGHJ", // 49 chars.
					"1234567890BCDFGHJKLMNPQRSTUVWXYZ1234567890BCDFG$",  // 48 chars with invalid char.
				},
			},
		),
		gen(
			"Base32HexEscaped",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base32HexEscaped,
			},
			&action{
				valid: []string{
					"1234567890BCDFGHJKLMNPQRSTUV1234567890BCDFGHJKLM", // 48 chars.
				},
				invalid: []string{
					"1234567890BCDFGHJKLMNPQRSTUV1234567890BCDFGHJKL",   // 47 chars.
					"1234567890BCDFGHJKLMNPQRSTUV1234567890BCDFGHJKLMN", // 49 chars.
					"1234567890BCDFGHJKLMNPQRSTUV1234567890BCDFGHJKL$",  // 48 chars with invalid char.
				},
			},
		),
		gen(
			"Base64",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base64,
			},
			&action{
				valid: []string{
					base64.StdEncoding.EncodeToString(must(uid.NewID())),
					base64.StdEncoding.EncodeToString(must(uid.NewHostedID())),
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/00", // 40 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz+/00", // 40 chars.
				},
				invalid: []string{
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/000", // 39 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz+/000", // 39 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/000", // 41 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz+/000", // 41 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/0$",  // 40 chars with invalid char.
					"1234567890abcdefghijklmnopqrstuvwxyz+/0$",  // 40 chars with invalid char.
				},
			},
		),
		gen(
			"Base64Raw",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base64Raw,
			},
			&action{
				valid: []string{
					base64.RawStdEncoding.EncodeToString(must(uid.NewID())),
					base64.RawStdEncoding.EncodeToString(must(uid.NewHostedID())),
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/00", // 40 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz+/00", // 40 chars.
				},
				invalid: []string{
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/000", // 39 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz+/000", // 39 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/000", // 41 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz+/000", // 41 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ+/0$",  // 40 chars with invalid char.
					"1234567890abcdefghijklmnopqrstuvwxyz+/0$",  // 40 chars with invalid char.
				},
			},
		),
		gen(
			"Base64URL",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base64URL,
			},
			&action{
				valid: []string{
					base64.URLEncoding.EncodeToString(must(uid.NewID())),
					base64.URLEncoding.EncodeToString(must(uid.NewHostedID())),
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_00", // 40 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz-_00", // 40 chars.
				},
				invalid: []string{
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_000", // 39 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz-_000", // 39 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_000", // 41 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz-_000", // 41 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_0$",  // 40 chars with invalid char.
					"1234567890abcdefghijklmnopqrstuvwxyz-_0$",  // 40 chars with invalid char.
				},
			},
		),
		gen(
			"Base64RawURL",
			[]string{cndValidType},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType_Base64RawURL,
			},
			&action{
				valid: []string{
					base64.RawURLEncoding.EncodeToString(must(uid.NewID())),
					base64.RawURLEncoding.EncodeToString(must(uid.NewHostedID())),
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_00", // 40 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz-_00", // 40 chars.
				},
				invalid: []string{
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_000", // 39 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz-_000", // 39 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_000", // 41 chars.
					"1234567890abcdefghijklmnopqrstuvwxyz-_000", // 41 chars.
					"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ-_0$",  // 40 chars with invalid char.
					"1234567890abcdefghijklmnopqrstuvwxyz-_0$",  // 40 chars with invalid char.
				},
			},
		),
		gen(
			"Unknown",
			[]string{cndUnknown},
			[]string{actCheckType},
			&condition{
				typ: kernel.EncodingType(999),
			},
			&action{
				shouldNil: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			reg := uid.Validator(tt.C().typ)
			if tt.A().shouldNil {
				testutil.Diff(t, (*regexp.Regexp)(nil), reg)
				return
			}
			for _, v := range tt.A().valid {
				t.Log("expect true", v)
				testutil.Diff(t, true, reg.MatchString(v))
			}
			for _, v := range tt.A().invalid {
				t.Log("expect false", v)
				testutil.Diff(t, false, reg.MatchString(v))
			}
		})
	}
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
			pattern := uid.Validator(tt.C().typ)
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
			pattern := uid.Validator(tt.C().typ)
			testutil.Diff(t, true, pattern.MatchString(eid))
		})
	}
}
