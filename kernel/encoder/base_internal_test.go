// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encoder

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestEncoderDecoder(t *testing.T) {
	type condition struct {
		typ k.EncodingType
	}
	type action struct {
		enc     EncodeToStringFunc
		dec     DecodeStringFunc
		pattern *regexp.Regexp
	}
	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"Base32",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base32,
			},
			&action{
				enc:     base32.StdEncoding.EncodeToString,
				dec:     base32.StdEncoding.DecodeString,
				pattern: regexp.MustCompile("^[A-Z2-7=]+$"),
			},
		),
		gen(
			"Base32Hex",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base32Hex,
			},
			&action{
				enc:     base32.HexEncoding.EncodeToString,
				dec:     base32.HexEncoding.DecodeString,
				pattern: regexp.MustCompile("^[0-9A-V=]+$"),
			},
		),
		gen(
			"Base32Escaped",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base32Escaped,
			},
			&action{
				enc:     Base32StdEscapedEncoding.EncodeToString,
				dec:     Base32StdEscapedEncoding.DecodeString,
				pattern: regexp.MustCompile("^[BCDFGHJKLMNPQRSTUVWXYZ0-9=]+$"),
			},
		),
		gen(
			"Base32HexEscaped",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base32HexEscaped,
			},
			&action{
				enc:     Base32HexEscapedEncoding.EncodeToString,
				dec:     Base32HexEscapedEncoding.DecodeString,
				pattern: regexp.MustCompile("^[0-9BCDFGHJKLMNPQRSTUVWXYZ=]+$"),
			},
		),
		gen(
			"Base64",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base64,
			},
			&action{
				enc:     base64.StdEncoding.EncodeToString,
				dec:     base64.StdEncoding.DecodeString,
				pattern: regexp.MustCompile("^[A-Za-z0-9+/=]+$"),
			},
		),
		gen(
			"Base64Raw",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base64Raw,
			},
			&action{
				enc:     base64.RawStdEncoding.EncodeToString,
				dec:     base64.RawStdEncoding.DecodeString,
				pattern: regexp.MustCompile("^[A-Za-z0-9+/]+$"),
			},
		),
		gen(
			"Base64URL",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base64URL,
			},
			&action{
				enc:     base64.URLEncoding.EncodeToString,
				dec:     base64.URLEncoding.DecodeString,
				pattern: regexp.MustCompile("^[A-Za-z0-9-_=]+$"),
			},
		),
		gen(
			"Base64RawURL",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base64RawURL,
			},
			&action{
				enc:     base64.RawURLEncoding.EncodeToString,
				dec:     base64.RawURLEncoding.DecodeString,
				pattern: regexp.MustCompile("^[A-Za-z0-9-_]+$"),
			},
		),
		gen(
			"Base16",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_Base16,
			},
			&action{
				enc:     hex.EncodeToString,
				dec:     hex.DecodeString,
				pattern: regexp.MustCompile("^[0-9a-f]+$"),
			},
		),
		gen(
			"Unknown will be Base32HexEscaped",
			[]string{}, []string{},
			&condition{
				typ: k.EncodingType_EncodingTypeUnknown,
			},
			&action{
				enc:     Base32StdEscapedEncoding.EncodeToString,
				dec:     Base32StdEscapedEncoding.DecodeString,
				pattern: regexp.MustCompile("^[0-9BCDFGHJKLMNPQRSTUVWXYZ=]+$"),
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			enc, dec := EncoderDecoder(tt.C().typ)
			testutil.Diff(t, tt.A().enc, enc, cmp.Comparer(testutil.ComparePointer[EncodeToStringFunc]))
			testutil.Diff(t, tt.A().dec, dec, cmp.Comparer(testutil.ComparePointer[DecodeStringFunc]))

			data := "abcdefghijklmnopqrstuvwxyz1234 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
			e := enc([]byte(data))
			fmt.Println(tt.Name(), e)
			testutil.Diff(t, true, tt.A().pattern.MatchString(e))
			d, err := dec(e)
			testutil.Diff(t, nil, err)
			testutil.Diff(t, data, string(d))
		})
	}
}
