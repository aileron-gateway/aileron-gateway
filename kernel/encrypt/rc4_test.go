package encrypt_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	_ = encrypt.EncryptFunc(encrypt.EncryptRC4) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptRC4) // Check that the function satisfies the signature.
)

func ExampleEncryptRC4() {
	key := []byte("123456789012345678901234567890")

	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptRC4(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Example Output:
	// The result is different at all times because of the iv (initial vector).
	// de5f193c0dce2b344b044b7a1bacc6e9429572714dbf02499a7f64fbc0ae0668b5e80e08ecd342eebd
}

func ExampleDecryptRC4() {
	key := []byte("123456789012345678901234567890")

	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptRC4(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptRC4(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func TestEncryptRC4(t *testing.T) {
	type condition struct {
		inputNil       bool
		input          []byte
		key            []byte
		useErrorReader bool
	}

	type action struct {
		err error
	}

	CndInputNil := "nil input"
	CndInputNonNil := "non-nil input"
	CndInputValidKey := "input valid key"
	CndInputInvalidKey := "input invalid key"
	CndIOError := "io error"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndInputNonNil, "input non-zero and non-nil bytes")
	tb.Condition(CndInputValidKey, "input valid length key")
	tb.Condition(CndInputInvalidKey, "input invalid length key")
	tb.Condition(CndIOError, "io error occurred while reading random")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the encrypted bytes does not contain the original data")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil input with valid key",
			[]string{CndInputNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: true,
				input:    nil,
				key:      []byte("this is a valid key"),
			},
			&action{
				err: nil,
			},
		),
		gen(
			"non-nil input with valid key",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    []byte("test"),
				key:      []byte("this is a valid key"),
			},
			&action{
				err: nil,
			},
		),
		gen(
			"key length 0",
			[]string{CndInputNonNil, CndInputInvalidKey},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil: false,
				input:    []byte("test"),
				key:      []byte(""),
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeRC4,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
		gen(
			"key length 1",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    []byte("test"),
				key:      []byte("1"),
			},
			&action{
				err: nil,
			},
		),
		gen(
			"key length 10",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    []byte("test"),
				key:      []byte("1234567890"),
			},
			&action{
				err: nil,
			},
		),
		gen(
			"key length 256",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    []byte("test"),
				key: []byte("12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"123456"),
			},
			&action{
				err: nil,
			},
		),
		gen(
			"key length 257",
			[]string{CndInputNonNil, CndInputInvalidKey},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil: false,
				input:    []byte("test"),
				key: []byte("12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"1234567"),
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeRC4,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
		gen(
			"io error",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil:       false,
				input:          []byte("test"),
				key:            []byte("this is a valid key"),
				useErrorReader: true,
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeRC4,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Initialize
			if tt.C().useErrorReader {
				tmp := rand.Reader
				rand.Reader = &testutil.ErrorReader{}
				defer func() {
					rand.Reader = tmp
				}()
			}

			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.EncryptRC4(tt.C().key, nil)
			} else {
				out, err = encrypt.EncryptRC4(tt.C().key, tt.C().input)
			}

			if len(tt.C().input) > 0 && bytes.Contains(out, tt.C().input) {
				t.Error("ciphertext contains original plaintext")
			}

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}

func TestDecryptRC4(t *testing.T) {
	type condition struct {
		inputNil bool
		input    string // hex encoded ciphertext
		key      []byte
	}

	type action struct {
		expect string
		err    error
	}

	CndInputNil := "nil input"
	CndInputNonNil := "non-nil input"
	CndInputValidKey := "input valid key"
	CndInputInvalidKey := "input invalid key"
	ActCheckNoError := "check no error"
	ActCheckError := "check error"
	ActCheckData := "check data"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(CndInputNil, "input nil as bytes")
	tb.Condition(CndInputNonNil, "input non-zero and non-nil bytes")
	tb.Condition(CndInputValidKey, "input valid length key")
	tb.Condition(CndInputInvalidKey, "input invalid length key")
	tb.Action(ActCheckNoError, "check that there is no error occurred")
	tb.Action(ActCheckError, "check that expected error occurred")
	tb.Action(ActCheckData, "check that the encrypted bytes does not contain the original data")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"nil input with valid key",
			[]string{CndInputNil, CndInputValidKey},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil: true,
				input:    "",
				key:      []byte("this is a valid key"),
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeRC4,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"non-nil input with valid key",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    "29eb4d3d3c499f3c85ccfad9a487c34411bd451819dc7f230ddd5acc",
				key:      []byte("this is a valid key"),
			},
			&action{
				expect: "test",
			},
		),
		gen(
			"key length 0",
			[]string{CndInputNonNil, CndInputInvalidKey},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil: false,
				input:    "29eb4d3d3c499f3c85ccfad9a487c34411bd451819dc7f230ddd5acc",
				key:      []byte(""),
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeRC4,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"key length 1",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    "97059b4023508f3e675430f6e5a3be3c62da5d8246d5afa412ddd09f",
				key:      []byte("1"),
			},
			&action{
				err:    nil,
				expect: "test",
			},
		),
		gen(
			"key length 10",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    "7d78449e091c9784d09a9b5267053a79d3233b7cd34b6ab69e7676e6",
				key:      []byte("1234567890"),
			},
			&action{
				err:    nil,
				expect: "test",
			},
		),
		gen(
			"key length 256",
			[]string{CndInputNonNil, CndInputValidKey},
			[]string{ActCheckNoError, ActCheckData},
			&condition{
				inputNil: false,
				input:    "e46b0490b2e007bd94d356991276389efd933fde55df921e9e7676e6",
				key: []byte("12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"123456"),
			},
			&action{
				err:    nil,
				expect: "test",
			},
		),
		gen(
			"key length 257",
			[]string{CndInputNonNil, CndInputInvalidKey},
			[]string{ActCheckError, ActCheckData},
			&condition{
				inputNil: false,
				input:    "29eb4d3d3c499f3c85ccfad9a487c34411bd451819dc7f230ddd5acc",
				key: []byte("12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"12345678901234567890123456789012345678901234567890" +
					"1234567"),
			},
			&action{
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeRC4,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// decode input string for input.
			in, e := hex.DecodeString(tt.C().input)
			testutil.Diff(t, nil, e)

			var out []byte
			var err error
			if tt.C().inputNil {
				out, err = encrypt.DecryptRC4(tt.C().key, nil)
			} else {
				out, err = encrypt.DecryptRC4(tt.C().key, in)
			}

			testutil.Diff(t, tt.A().expect, string(out))
			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
		})
	}
}
