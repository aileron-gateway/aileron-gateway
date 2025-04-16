// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

func ExampleFixedReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Fixed{
			Fixed: &k.FixedReplacer{
				Value: "replaced",
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "example"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	replaced
}

func ExampleValueReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Value{
			Value: &k.ValueReplacer{
				FromTo: map[string]string{
					"foo": "alice",
					"bar": "bob",
				},
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "foo and bar"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	alice and bob
}

func ExampleLeftReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Left{
			Left: &k.LeftReplacer{
				Char:   "*",
				Length: 3,
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	***-456-789
}

func ExampleRightReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Right{
			Right: &k.RightReplacer{
				Char:   "*",
				Length: 3,
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-456-***
}

func ExampleTrimReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Trim{
			Trim: &k.TrimReplacer{
				CutSets: []string{"*#"},
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "*#*#123-456-789*#*#"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-456-789
}

func ExampleTrimLeftReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_TrimLeft{
			TrimLeft: &k.TrimLeftReplacer{
				CutSets: []string{"*#"},
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "*#*#123-456-789*#*#"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-456-789*#*#
}

func ExampleTrimRightReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_TrimRight{
			TrimRight: &k.TrimRightReplacer{
				CutSets: []string{"*#"},
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "*#*#123-456-789*#*#"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	*#*#123-456-789
}

func ExampleTrimPrefixReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_TrimPrefix{
			TrimPrefix: &k.TrimPrefixReplacer{
				Prefixes: []string{"Phone: "},
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "Phone: 123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-456-789
}

func ExampleTrimSuffixReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_TrimSuffix{
			TrimSuffix: &k.TrimSuffixReplacer{
				Suffixes: []string{" (+81)"},
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789 (+81)"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-456-789
}

func ExampleEncodeReplacer_encodeAll() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Encode{
			Encode: &k.EncodeReplacer{
				Encoding: k.EncodingType_Base64,
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	MTIzLTQ1Ni03ODk=
}

func ExampleEncodeReplacer_encodeMatch() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Encode{
			Encode: &k.EncodeReplacer{
				Pattern:  "[4-6]{3}",
				Encoding: k.EncodingType_Base64,
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-NDU2-789
}

func ExampleHashReplacer_hashAll() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Hash{
			Hash: &k.HashReplacer{
				Alg:      k.HashAlg_SHA256,
				Encoding: k.EncodingType_Base64,
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	XedIoOiHNR5HjPNnACwzdo4ug7MLk6RSC+/qSWwSBnk=
}

func ExampleHashReplacer_hashPart() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Hash{
			Hash: &k.HashReplacer{
				Pattern:  "[4-6]{3}",
				Alg:      k.HashAlg_SHA256,
				Encoding: k.EncodingType_Base64,
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-s6jg4fmrG/46NvIx9nb3i7MKUZ0rIebFMMDu6Ou0pdA=-789
}

func ExampleRegexpReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Regexp{
			Regexp: &k.RegexpReplacer{
				Pattern: "[0-9]([0-9])[0-9]",
				Replace: "$1",
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	2-5-8
}

func ExampleExpandReplacer() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Expand{
			Expand: &k.ExpandReplacer{
				Pattern:  "[0-9]([0-9])[0-9]",
				Template: "$1",
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	258
}

func ExampleEncryptReplacer() {
	// Fix the random value for reducibility.
	// Random value is used as a initial vector.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("1234567890123456"))
	defer func() {
		rand.Reader = tmp
	}()

	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_Encrypt{
			Encrypt: &k.EncryptReplacer{
				Pattern:  "[4-6]{3}",
				Alg:      k.CommonKeyCryptType_AESCBC,
				Encoding: k.EncodingType_Base16,
				Password: hex.EncodeToString([]byte("16_bytes_secrets")),
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-3132333435363738393031323334353670c048862b4ccd751b35ae4889e901f7-789
}

func ExampleHMACReplacer_hashAll() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_HMAC{
			HMAC: &k.HMACReplacer{
				Alg:      k.HashAlg_SHA256,
				Encoding: k.EncodingType_Base64,
				Key:      hex.EncodeToString([]byte("hmac-key")),
			},
		},
	}
	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// dzBuTQoT6Bpub2nIGs6vQfOba9TEYWm+gWqgo/hdpyw=
}

func ExampleHMACReplacer_hashPart() {
	spec := &k.ReplacerSpec{
		Replacers: &k.ReplacerSpec_HMAC{
			HMAC: &k.HMACReplacer{
				Pattern:  "[4-6]{3}",
				Alg:      k.HashAlg_SHA256,
				Encoding: k.EncodingType_Base64,
				Key:      hex.EncodeToString([]byte("hmac-key")),
			},
		},
	}

	replacer, err := txtutil.NewStringReplacer(spec)
	if err != nil {
		panic("handle error here")
	}

	value := "123-456-789"
	result := replacer.Replace(value)

	fmt.Println(string(result))
	// Output:
	// 	123-fqlEWXgPubzJUBErpbq1C0QvNfXNSFWOxG/ufgAg6iE=-789
}
