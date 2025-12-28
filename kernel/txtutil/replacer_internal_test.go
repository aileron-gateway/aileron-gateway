// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewStringReplacers(t *testing.T) {
	type condition struct {
		specs []*k.ReplacerSpec
	}

	type action struct {
		inout map[string]string
		err   error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil specs", &condition{
				specs: nil,
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"nil replacer", &condition{
				specs: []*k.ReplacerSpec{
					{Replacers: nil},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"multiple specs", &condition{
				specs: []*k.ReplacerSpec{
					{
						Replacers: &k.ReplacerSpec_Value{
							Value: &k.ValueReplacer{
								FromTo: map[string]string{
									"foo": "***",
								},
							},
						},
					},
					{
						Replacers: &k.ReplacerSpec_Value{
							Value: &k.ValueReplacer{
								FromTo: map[string]string{
									"alice": "###",
								},
							},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "***=bar",
					"alice,bob":    "###,bob",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			reps, err := NewStringReplacers(tt.C.specs...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for k, v := range tt.A.inout {
				t.Log("Replace:", k, "Expect:", v)
				out := k
				for _, rep := range reps {
					out = rep.Replace(out)
				}
				testutil.Diff(t, v, out)
			}
		})
	}
}

func TestNewStringReplacer(t *testing.T) {
	type condition struct {
		spec *k.ReplacerSpec
	}

	type action struct {
		inout map[string]string
		err   error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec", &condition{
				spec: nil,
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscNil,
				},
			},
		),
		gen(
			"nil replacer", &condition{
				spec: &k.ReplacerSpec{
					Replacers: nil,
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"fixed", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Fixed{
						Fixed: &k.FixedReplacer{
							Value: "***",
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "***",
					"123-456-7890": "***",
					"foo=bar":      "***",
					"alice,bob":    "***",
				},
			},
		),
		gen(
			"value/no value", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Value{
						Value: &k.ValueReplacer{
							FromTo: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"value/with value", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Value{
						Value: &k.ValueReplacer{
							FromTo: map[string]string{
								"foo":   "***",
								"alice": "###",
							},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "***=bar",
					"alice,bob":    "###,bob",
				},
			},
		),
		gen(
			"left/empty char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Left{
						Left: &k.LeftReplacer{
							Char:   "",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "-456-7890",
					"foo=bar":      "=bar",
					"alice,bob":    "ce,bob",
				},
			},
		),
		gen(
			"left/with char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Left{
						Left: &k.LeftReplacer{
							Char:   "*",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "***-456-7890",
					"foo=bar":      "***=bar",
					"alice,bob":    "***ce,bob",
				},
			},
		),
		gen(
			"right/empty char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Right{
						Right: &k.RightReplacer{
							Char:   "",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7",
					"foo=bar":      "foo=",
					"alice,bob":    "alice,",
				},
			},
		),
		gen(
			"right/with char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Right{
						Right: &k.RightReplacer{
							Char:   "*",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7***",
					"foo=bar":      "foo=***",
					"alice,bob":    "alice,***",
				},
			},
		),
		gen(
			"trim/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Trim{
						Trim: &k.TrimReplacer{
							CutSets: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trim/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Trim{
						Trim: &k.TrimReplacer{
							CutSets: []string{"abc012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-789", // "0", "1" and "2" trimmed.
					"foo=bar":      "foo=bar",   // nothing trimmed.
					"alice,bob":    "lice,bo",   // "a" and "b" trimmed.
					"abc012":       "",          // all chars trimmed.
				},
			},
		),
		gen(
			"trim/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Trim{
						Trim: &k.TrimReplacer{
							CutSets: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-789", // "0", "1" and "2" trimmed.
					"foo=bar":      "foo=bar",   // nothing trimmed.
					"alice,bob":    "lice,bo",   // "a" and "b" trimmed.
					"abc012":       "012",       // First matched cutSets "abc" trimmed.
				},
			},
		),
		gen(
			"trimLeft/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimLeft{
						TrimLeft: &k.TrimLeftReplacer{
							CutSets: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimLeft/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimLeft{
						TrimLeft: &k.TrimLeftReplacer{
							CutSets: []string{"abc012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-7890", //  "1" and "2" trimmed.
					"foo=bar":      "foo=bar",    // nothing trimmed.
					"alice,bob":    "lice,bob",   // "a"  trimmed.
					"abc012":       "",           // all chars trimmed.
				},
			},
		),
		gen(
			"trimLeft/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimLeft{
						TrimLeft: &k.TrimLeftReplacer{
							CutSets: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-7890", //  "1" and "2" trimmed.
					"foo=bar":      "foo=bar",    // nothing trimmed.
					"alice,bob":    "lice,bob",   // "a"  trimmed.
					"abc012":       "012",        // First matched cutSets "abc" trimmed.
				},
			},
		),
		gen(
			"trimRight/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimRight{
						TrimRight: &k.TrimRightReplacer{
							CutSets: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimRight/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimRight{
						TrimRight: &k.TrimRightReplacer{
							CutSets: []string{"abc012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-789", //  "0" trimmed.
					"foo=bar":      "foo=bar",     // nothing trimmed.
					"alice,bob":    "alice,bo",    // "b"  trimmed.
					"abc012":       "",            // all chars trimmed.
				},
			},
		),
		gen(
			"trimRight/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimRight{
						TrimRight: &k.TrimRightReplacer{
							CutSets: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-789", //  "0" trimmed.
					"foo=bar":      "foo=bar",     // nothing trimmed.
					"alice,bob":    "alice,bo",    // "b"  trimmed.
					"abc012":       "abc",         // First matched cutSets "012" trimmed.
				},
			},
		),
		gen(
			"trimPrefix/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimPrefix{
						TrimPrefix: &k.TrimPrefixReplacer{
							Prefixes: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimPrefix/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimPrefix{
						TrimPrefix: &k.TrimPrefixReplacer{
							Prefixes: []string{"abc"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "012",
				},
			},
		),
		gen(
			"trimPrefix/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimPrefix{
						TrimPrefix: &k.TrimPrefixReplacer{
							Prefixes: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "012",
				},
			},
		),

		gen(
			"trimSuffix/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimSuffix{
						TrimSuffix: &k.TrimSuffixReplacer{
							Suffixes: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimSuffix/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimSuffix{
						TrimSuffix: &k.TrimSuffixReplacer{
							Suffixes: []string{"abc"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "abc012",
				},
			},
		),
		gen(
			"trimSuffix/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimSuffix{
						TrimSuffix: &k.TrimSuffixReplacer{
							Suffixes: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "abc",
				},
			},
		),
		gen(
			"encode/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encode{
						Encode: &k.EncodeReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encode/without pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encode{
						Encode: &k.EncodeReplacer{
							Encoding: k.EncodingType_Base64,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "MTIzLTQ1Ni03ODkw",
					"foo=bar":      "Zm9vPWJhcg==",
					"alice,bob":    "YWxpY2UsYm9i",
				},
			},
		),
		gen(
			"encode/with pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encode{
						Encode: &k.EncodeReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Encoding: k.EncodingType_Base64,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-NDU2-7890",
					"foo=bar":      "Zm9v=bar",
					"alice,bob":    "YWxpY2U=,bob",
				},
			},
		),
		gen(
			"hash/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hash/no hash", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hash/no encoding", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Alg: k.HashAlg_SHA256,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hash/without pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"123-456-7890": "29ec0a06044bedff0bec4d81516da5843c93b7b4cb8b3ec0bf379a78085747bd",
					"foo=bar":      "3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a",
					"alice,bob":    "f0e50e8fced90be1bda596553cfabf29ecee566958f394452b6d401b192f6dce",
				},
			},
		),
		gen(
			"hash/with pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-b3a8e0e1f9ab1bfe3a36f231f676f78bb30a519d2b21e6c530c0eee8ebb4a5d0-7890",
					"foo=bar":      "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae=bar",
					"alice,bob":    "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90,bob",
				},
			},
		),
		gen(
			"regexp/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regexp/invalid regexp", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `[0-9a-`,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regexp/invalid POSIX regexp", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `[0-9a-`,
							POSIX:   true,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regexp/posix=false,literal=false", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `([0-9]{3}|foo|alice)`,
							Replace: `*`,
							POSIX:   false,
							Literal: false,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "*-*-*0",
					"foo=bar":      "*=bar",
					"alice,bob":    "*,bob",
				},
			},
		),
		gen(
			"regexp/posix=true", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `([0-9]{3}|foo|alice)`,
							Replace: `*`,
							POSIX:   true,
							Literal: false,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "*-*-*0",
					"foo=bar":      "*=bar",
					"alice,bob":    "*,bob",
				},
			},
		),
		gen(
			"regexp/literal=true", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `([0-9]{3}|foo|alice)`,
							Replace: `$1`,
							POSIX:   false,
							Literal: true,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "$1-$1-$10",
					"foo=bar":      "$1=bar",
					"alice,bob":    "$1,bob",
				},
			},
		),
		gen(
			"expand/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"expand/invalid regexp", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern: `[0-9a-`,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"expand/invalid POSIX regexp", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern: `[0-9a-`,
							POSIX:   true,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"expand/posix=false,literal=false", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern:  `([0-9]{3}|foo|alice)`,
							Template: `*`,
							POSIX:    false,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "***",
					"foo=bar":      "*",
					"alice,bob":    "*",
				},
			},
		),
		gen(
			"expand/posix=true", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern:  `([0-9]{3}|foo|alice)`,
							Template: `*`,
							POSIX:    true,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "***",
					"foo=bar":      "*",
					"alice,bob":    "*",
				},
			},
		),
		gen(
			"encrypt/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encrypt/no alg", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encrypt/no encoding", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encrypt/invalid password", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("invalid length")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "!ERROR[crypto/aes: invalid key size 14]",
					"123-456-7890": "!ERROR[crypto/aes: invalid key size 14]",
					"foo=bar":      "!ERROR[crypto/aes: invalid key size 14]",
					"alice,bob":    "!ERROR[crypto/aes: invalid key size 14]",
				},
			},
		),
		gen(
			"encrypt/invalid password with pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("invalid length")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-!ERROR[crypto/aes: invalid key size 14]-7890",
					"foo=bar":      "!ERROR[crypto/aes: invalid key size 14]=bar",
					"alice,bob":    "!ERROR[crypto/aes: invalid key size 14],bob",
				},
			},
		),
		gen(
			"encrypt/invalid hex password", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: "INVALID_Hex",
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"encrypt/without pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "313233343536373839303132333435365ecacef641fb482c85c05d0790b2f714",
					"123-456-7890": "31323334353637383930313233343536e6983166eab85d52e24ae2b01ba6e373",
					"foo=bar":      "313233343536373839303132333435362aa0145db4669f5ea7b91aee38b73fed",
					"alice,bob":    "313233343536373839303132333435364babde52fde974ad663986822d9d6f7f",
				},
			},
		),
		gen(
			"encrypt/with pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-313233343536373839303132333435367caf7f92b33d98fca4916aabb75d1142-7890",
					"foo=bar":      "313233343536373839303132333435360ec0570462f1fb9bfe98ec20df31eff8=bar",
					"alice,bob":    "31323334353637383930313233343536b670297f836facf2358c2cbbc9590f94,bob",
				},
			},
		),
		gen(
			"hmac/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hmac/no hash", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hmac/no encoding", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Alg: k.HashAlg_SHA256,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hmac/invalid key", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base16,
							Key:      "Invalid Hex Key",
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"hmac/without pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base64,
							Key:      hex.EncodeToString([]byte("key")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "XV0TlWPJW1lnub2ajJsjOp3ttFByeUzSMtwbdIMmB9A=",
					"123-456-7890": "/nJh2RdsJUhOliz5mxMLH7m74HHNlg3CZufnEFsOCf4=",
					"foo=bar":      "Yiah6setCP3dBTAArIlpjf7aoys6dnOXJMcjRELn9HE=",
					"alice,bob":    "a2gqD0ndiZdbGCMk37kff3Hza5KoPDIg2iVW5BKW9DA=",
				},
			},
		),
		gen(
			"hmac/with pattern", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base64,
							Key:      hex.EncodeToString([]byte("key")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-1fPIVCyIOletgTeQXfNzfszl9l5kKG1LAL0hlb7ST30=-7890",
					"foo=bar":      "bqHZ9ek6jzreAmJh/+XXKhyQgE7ZRASmmJKhY7ijVJc==bar",
					"alice,bob":    "dvtV6SnAa5ewHDWVDuX3L+QVsV7TpzVsOecJkG27XEU=,bob",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			// Fix the random value for reducibility.
			// Random value is used for initial vectors when encryption.
			tmp := rand.Reader
			rand.Reader = bytes.NewReader(bytes.Repeat([]byte("1234567890123456"), 10))
			defer func() {
				rand.Reader = tmp
			}()

			rep, err := NewStringReplacer(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for k, v := range tt.A.inout {
				t.Log("Replace:", k, "Expect:", v)
				out := rep.Replace(k)
				testutil.Diff(t, v, out)
			}
		})
	}
}

func TestNewBytesReplacers(t *testing.T) {
	type condition struct {
		specs []*k.ReplacerSpec
	}

	type action struct {
		inout map[string]string
		err   error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil specs", &condition{
				specs: nil,
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"nil replacer", &condition{
				specs: []*k.ReplacerSpec{
					{Replacers: nil},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"multiple specs", &condition{
				specs: []*k.ReplacerSpec{
					{
						Replacers: &k.ReplacerSpec_Value{
							Value: &k.ValueReplacer{
								FromTo: map[string]string{
									"foo": "***",
								},
							},
						},
					},
					{
						Replacers: &k.ReplacerSpec_Value{
							Value: &k.ValueReplacer{
								FromTo: map[string]string{
									"alice": "###",
								},
							},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "***=bar",
					"alice,bob":    "###,bob",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			reps, err := NewBytesReplacers(tt.C.specs...)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for k, v := range tt.A.inout {
				t.Log("Replace:", k, "Expect:", v)
				out := []byte(k)
				for _, rep := range reps {
					out = rep.Replace(out)
				}
				testutil.Diff(t, v, string(out))
			}
		})
	}
}

func TestNewBytesReplacer(t *testing.T) {
	type condition struct {
		spec *k.ReplacerSpec
	}

	type action struct {
		inout map[string]string
		err   error
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil spec", &condition{
				spec: nil,
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscNil,
				},
			},
		),
		gen(
			"nil replacer", &condition{
				spec: &k.ReplacerSpec{
					Replacers: nil,
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"fixed", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Fixed{
						Fixed: &k.FixedReplacer{
							Value: "***",
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "***",
					"123-456-7890": "***",
					"foo=bar":      "***",
					"alice,bob":    "***",
				},
			},
		),
		gen(
			"value/no value", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Value{
						Value: &k.ValueReplacer{
							FromTo: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"value/with value", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Value{
						Value: &k.ValueReplacer{
							FromTo: map[string]string{
								"foo":   "***",
								"alice": "###",
							},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "***=bar",
					"alice,bob":    "###,bob",
				},
			},
		),
		gen(
			"left/empty char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Left{
						Left: &k.LeftReplacer{
							Char:   "",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "-456-7890",
					"foo=bar":      "=bar",
					"alice,bob":    "ce,bob",
				},
			},
		),
		gen(
			"left/with char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Left{
						Left: &k.LeftReplacer{
							Char:   "*",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "***-456-7890",
					"foo=bar":      "***=bar",
					"alice,bob":    "***ce,bob",
				},
			},
		),
		gen(
			"right/empty char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Right{
						Right: &k.RightReplacer{
							Char:   "",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7",
					"foo=bar":      "foo=",
					"alice,bob":    "alice,",
				},
			},
		),
		gen(
			"right/with char", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Right{
						Right: &k.RightReplacer{
							Char:   "*",
							Length: 3,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7***",
					"foo=bar":      "foo=***",
					"alice,bob":    "alice,***",
				},
			},
		),
		gen(
			"trim/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Trim{
						Trim: &k.TrimReplacer{
							CutSets: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trim/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Trim{
						Trim: &k.TrimReplacer{
							CutSets: []string{"abc012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-789", // "0", "1" and "2" trimmed.
					"foo=bar":      "foo=bar",   // nothing trimmed.
					"alice,bob":    "lice,bo",   // "a" and "b" trimmed.
					"abc012":       "",          // all chars trimmed.
				},
			},
		),
		gen(
			"trim/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Trim{
						Trim: &k.TrimReplacer{
							CutSets: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-789", // "0", "1" and "2" trimmed.
					"foo=bar":      "foo=bar",   // nothing trimmed.
					"alice,bob":    "lice,bo",   // "a" and "b" trimmed.
					"abc012":       "012",       // First matched cutSets "abc" trimmed.
				},
			},
		),
		gen(
			"trimLeft/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimLeft{
						TrimLeft: &k.TrimLeftReplacer{
							CutSets: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimLeft/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimLeft{
						TrimLeft: &k.TrimLeftReplacer{
							CutSets: []string{"abc012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-7890", //  "1" and "2" trimmed.
					"foo=bar":      "foo=bar",    // nothing trimmed.
					"alice,bob":    "lice,bob",   // "a"  trimmed.
					"abc012":       "",           // all chars trimmed.
				},
			},
		),
		gen(
			"trimLeft/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimLeft{
						TrimLeft: &k.TrimLeftReplacer{
							CutSets: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "3-456-7890", //  "1" and "2" trimmed.
					"foo=bar":      "foo=bar",    // nothing trimmed.
					"alice,bob":    "lice,bob",   // "a"  trimmed.
					"abc012":       "012",        // First matched cutSets "abc" trimmed.
				},
			},
		),
		gen(
			"trimRight/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimRight{
						TrimRight: &k.TrimRightReplacer{
							CutSets: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimRight/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimRight{
						TrimRight: &k.TrimRightReplacer{
							CutSets: []string{"abc012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-789", //  "0" trimmed.
					"foo=bar":      "foo=bar",     // nothing trimmed.
					"alice,bob":    "alice,bo",    // "b"  trimmed.
					"abc012":       "",            // all chars trimmed.
				},
			},
		),
		gen(
			"trimRight/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimRight{
						TrimRight: &k.TrimRightReplacer{
							CutSets: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-789", //  "0" trimmed.
					"foo=bar":      "foo=bar",     // nothing trimmed.
					"alice,bob":    "alice,bo",    // "b"  trimmed.
					"abc012":       "abc",         // First matched cutSets "012" trimmed.
				},
			},
		),
		gen(
			"trimPrefix/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimPrefix{
						TrimPrefix: &k.TrimPrefixReplacer{
							Prefixes: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimPrefix/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimPrefix{
						TrimPrefix: &k.TrimPrefixReplacer{
							Prefixes: []string{"abc"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "012",
				},
			},
		),
		gen(
			"trimPrefix/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimPrefix{
						TrimPrefix: &k.TrimPrefixReplacer{
							Prefixes: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "012",
				},
			},
		),
		gen(
			"trimSuffix/no cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimSuffix{
						TrimSuffix: &k.TrimSuffixReplacer{
							Suffixes: nil,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
				},
			},
		),
		gen(
			"trimSuffix/single cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimSuffix{
						TrimSuffix: &k.TrimSuffixReplacer{
							Suffixes: []string{"abc"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "abc012",
				},
			},
		),
		gen(
			"trimSuffix/multiple cutSets", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_TrimSuffix{
						TrimSuffix: &k.TrimSuffixReplacer{
							Suffixes: []string{"abc", "012"},
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-456-7890",
					"foo=bar":      "foo=bar",
					"alice,bob":    "alice,bob",
					"abc012":       "abc",
				},
			},
		),
		gen(
			"encode/empty spec", &condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encode{
						Encode: &k.EncodeReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encode/without pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encode{
						Encode: &k.EncodeReplacer{
							Encoding: k.EncodingType_Base64,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "MTIzLTQ1Ni03ODkw",
					"foo=bar":      "Zm9vPWJhcg==",
					"alice,bob":    "YWxpY2UsYm9i",
				},
			},
		),
		gen(
			"encode/with pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encode{
						Encode: &k.EncodeReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Encoding: k.EncodingType_Base64,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-NDU2-7890",
					"foo=bar":      "Zm9v=bar",
					"alice,bob":    "YWxpY2U=,bob",
				},
			},
		),
		gen(
			"hash/empty spec",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hash/no hash",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hash/no encoding",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Alg: k.HashAlg_SHA256,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hash/without pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"123-456-7890": "29ec0a06044bedff0bec4d81516da5843c93b7b4cb8b3ec0bf379a78085747bd",
					"foo=bar":      "3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a",
					"alice,bob":    "f0e50e8fced90be1bda596553cfabf29ecee566958f394452b6d401b192f6dce",
				},
			},
		),
		gen(
			"hash/with pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Hash{
						Hash: &k.HashReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-b3a8e0e1f9ab1bfe3a36f231f676f78bb30a519d2b21e6c530c0eee8ebb4a5d0-7890",
					"foo=bar":      "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae=bar",
					"alice,bob":    "2bd806c97f0e00af1a1fc3328fa763a9269723c8db8fac4f93af71db186d6e90,bob",
				},
			},
		),
		gen(
			"regexp/empty spec",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regexp/invalid regexp",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `[0-9a-`,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regexp/invalid POSIX regexp",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `[0-9a-`,
							POSIX:   true,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"regexp/posix=false,literal=false",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `([0-9]{3}|foo|alice)`,
							Replace: `*`,
							POSIX:   false,
							Literal: false,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "*-*-*0",
					"foo=bar":      "*=bar",
					"alice,bob":    "*,bob",
				},
			},
		),
		gen(
			"regexp/posix=true",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `([0-9]{3}|foo|alice)`,
							Replace: `*`,
							POSIX:   true,
							Literal: false,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "*-*-*0",
					"foo=bar":      "*=bar",
					"alice,bob":    "*,bob",
				},
			},
		),
		gen(
			"regexp/literal=true",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Regexp{
						Regexp: &k.RegexpReplacer{
							Pattern: `([0-9]{3}|foo|alice)`,
							Replace: `$1`,
							POSIX:   false,
							Literal: true,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "$1-$1-$10",
					"foo=bar":      "$1=bar",
					"alice,bob":    "$1,bob",
				},
			},
		),
		gen(
			"expand/empty spec",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"expand/invalid regexp",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern: `[0-9a-`,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"expand/invalid POSIX regexp",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern: `[0-9a-`,
							POSIX:   true,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"expand/posix=false,literal=false",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern:  `([0-9]{3}|foo|alice)`,
							Template: `*`,
							POSIX:    false,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "***",
					"foo=bar":      "*",
					"alice,bob":    "*",
				},
			},
		),
		gen(
			"expand/posix=true",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Expand{
						Expand: &k.ExpandReplacer{
							Pattern:  `([0-9]{3}|foo|alice)`,
							Template: `*`,
							POSIX:    true,
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "***",
					"foo=bar":      "*",
					"alice,bob":    "*",
				},
			},
		),
		gen(
			"encrypt/empty spec",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encrypt/no alg",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encrypt/no encoding",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"encrypt/invalid password",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("invalid length")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "!ERROR[crypto/aes: invalid key size 14]",
					"123-456-7890": "!ERROR[crypto/aes: invalid key size 14]",
					"foo=bar":      "!ERROR[crypto/aes: invalid key size 14]",
					"alice,bob":    "!ERROR[crypto/aes: invalid key size 14]",
				},
			},
		),
		gen(
			"encrypt/invalid password with pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("invalid length")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-!ERROR[crypto/aes: invalid key size 14]-7890",
					"foo=bar":      "!ERROR[crypto/aes: invalid key size 14]=bar",
					"alice,bob":    "!ERROR[crypto/aes: invalid key size 14],bob",
				},
			},
		),
		gen(
			"encrypt/invalid hex password",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: "INVALID_Hex",
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"encrypt/without pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "313233343536373839303132333435365ecacef641fb482c85c05d0790b2f714",
					"123-456-7890": "31323334353637383930313233343536e6983166eab85d52e24ae2b01ba6e373",
					"foo=bar":      "313233343536373839303132333435362aa0145db4669f5ea7b91aee38b73fed",
					"alice,bob":    "313233343536373839303132333435364babde52fde974ad663986822d9d6f7f",
				},
			},
		),
		gen(
			"encrypt/with pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_Encrypt{
						Encrypt: &k.EncryptReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.CommonKeyCryptType_AESCBC,
							Encoding: k.EncodingType_Base16,
							Password: hex.EncodeToString([]byte("16_byte_password")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-313233343536373839303132333435367caf7f92b33d98fca4916aabb75d1142-7890",
					"foo=bar":      "313233343536373839303132333435360ec0570462f1fb9bfe98ec20df31eff8=bar",
					"alice,bob":    "31323334353637383930313233343536b670297f836facf2358c2cbbc9590f94,bob",
				},
			},
		),
		gen(
			"hmac/empty spec",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hmac/no hash",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Encoding: k.EncodingType_Base16,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hmac/no encoding",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Alg: k.HashAlg_SHA256,
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscUnsupported,
				},
			},
		),
		gen(
			"hmac/invalid key",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base16,
							Key:      "Invalid Hex Key",
						},
					},
				},
			},
			&action{
				err: &er.Error{
					Package:     ErrPkg,
					Type:        ErrTypeReplacer,
					Description: ErrDscPattern,
				},
			},
		),
		gen(
			"hmac/without pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base64,
							Key:      hex.EncodeToString([]byte("key")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "XV0TlWPJW1lnub2ajJsjOp3ttFByeUzSMtwbdIMmB9A=",
					"123-456-7890": "/nJh2RdsJUhOliz5mxMLH7m74HHNlg3CZufnEFsOCf4=",
					"foo=bar":      "Yiah6setCP3dBTAArIlpjf7aoys6dnOXJMcjRELn9HE=",
					"alice,bob":    "a2gqD0ndiZdbGCMk37kff3Hza5KoPDIg2iVW5BKW9DA=",
				},
			},
		),
		gen(
			"hmac/with pattern",
			&condition{
				spec: &k.ReplacerSpec{
					Replacers: &k.ReplacerSpec_HMAC{
						HMAC: &k.HMACReplacer{
							Pattern:  `([4-6]{3}|foo|alice)`,
							Alg:      k.HashAlg_SHA256,
							Encoding: k.EncodingType_Base64,
							Key:      hex.EncodeToString([]byte("key")),
						},
					},
				},
			},
			&action{
				inout: map[string]string{
					"":             "",
					"123-456-7890": "123-1fPIVCyIOletgTeQXfNzfszl9l5kKG1LAL0hlb7ST30=-7890",
					"foo=bar":      "bqHZ9ek6jzreAmJh/+XXKhyQgE7ZRASmmJKhY7ijVJc==bar",
					"alice,bob":    "dvtV6SnAa5ewHDWVDuX3L+QVsV7TpzVsOecJkG27XEU=,bob",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			// Fix the random value for reducibility.
			// Random value is used for initial vectors when encryption.
			tmp := rand.Reader
			rand.Reader = bytes.NewReader(bytes.Repeat([]byte("1234567890123456"), 20))
			defer func() {
				rand.Reader = tmp
			}()

			rep, err := NewBytesReplacer(tt.C.spec)
			testutil.Diff(t, tt.A.err, err, cmpopts.EquateErrors())
			if err != nil {
				return
			}

			for k, v := range tt.A.inout {
				t.Log("Replace:", k, "Expect:", v)
				out := rep.Replace([]byte(k))
				testutil.Diff(t, v, string(out))
			}
		})
	}
}
