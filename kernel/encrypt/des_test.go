package encrypt_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	_ = encrypt.EncryptFunc(encrypt.EncryptDESCFB) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptDESCTR) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptDESOFB) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptDESCBC) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptDESCFB) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptDESCTR) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptDESOFB) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptDESCBC) // Check that the function satisfies the signature.

	_ = encrypt.EncryptFunc(encrypt.EncryptTripleDESCFB) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptTripleDESCTR) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptTripleDESOFB) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptTripleDESCBC) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptTripleDESCFB) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptTripleDESCTR) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptTripleDESOFB) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptTripleDESCBC) // Check that the function satisfies the signature.
)

func ExampleEncryptDESCBC() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESCBC(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 6669786564207661396e874f3edea5920d97f7fecdf977476b6bb7a10d4b9932
}

func ExampleDecryptDESCBC() {
	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESCBC(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptDESCBC(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptDESCFB() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESCFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 666978656420766105f56ee1a289537533f15bc8e76421d7f7
}

func ExampleDecryptDESCFB() {
	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESCFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptDESCFB(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptDESCTR() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESCTR(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 666978656420766105f56ee1a28953757837a5f5d6e0e14661
}

func ExampleDecryptDESCTR() {
	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESCTR(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptDESCTR(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptDESOFB() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESOFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 666978656420766105f56ee1a289537557166bace6660729d0
}

func ExampleDecryptDESOFB() {
	key := []byte("12345678")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptDESOFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptDESOFB(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptTripleDESCBC() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESCBC(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 666978656420766103481229211cd8c8ecc9d22cf523a980144cc9b2ab37a0b4
}

func ExampleDecryptTripleDESCBC() {
	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESCBC(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptTripleDESCBC(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptTripleDESCFB() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESCFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 6669786564207661f4c4c29daf73d5dc3c844f1c9e324a5249
}

func ExampleDecryptTripleDESCFB() {
	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESCFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptTripleDESCFB(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptTripleDESCTR() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESCTR(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 6669786564207661f4c4c29daf73d5dcf5728ea80754ab33c4
}

func ExampleDecryptTripleDESCTR() {
	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESCTR(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptTripleDESCTR(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptTripleDESOFB() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESOFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 6669786564207661f4c4c29daf73d5dcbb3ad6464dda8c8b16
}

func ExampleDecryptTripleDESOFB() {
	key := []byte("123456789012345678901234")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptTripleDESOFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptTripleDESOFB(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func TestEncryptDESCBC(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptDESCBC(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptDESCBC(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	cndWrongKey := "wrong key"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Condition(cndWrongKey, "decrypt with a wrong key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "313233343536373851217225794069b1",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "313233343536373851217225794069b1",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"decrypt with wrong key",
			[]string{cndValidKeyLength, cndWrongKey},
			[]string{actCheckError},
			&condition{
				key:        []byte("xxxx_key"),
				ciphertext: "313233343536373851217225794069b1",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptDESCBC(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptDESCBC(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptDESCBC(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptDESCBC(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptTripleDESCBC(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptTripleDESCBC(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptTripleDESCBC(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	cndWrongKey := "wrong key"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Condition(cndWrongKey, "decrypt with a wrong key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "31323334353637387862e07b1a2c95f9",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "31323334353637387862e07b1a2c95f9",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"decrypt with wrong key",
			[]string{cndValidKeyLength, cndWrongKey},
			[]string{actCheckError},
			&condition{
				key:        []byte("xxxxkey_xxxxkey_xxxxkey_"),
				ciphertext: "31323334353637387862e07b1a2c95f9",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptTripleDESCBC(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptTripleDESCBC(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptTripleDESCBC(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptTripleDESCBC(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptDESCFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptDESCFB(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptDESCFB(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "3132333435363738267e72d0",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738267e72d0",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptDESCFB(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptDESCFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptDESCFB(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptDESCFB(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptTripleDESCFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptTripleDESCFB(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptTripleDESCFB(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "3132333435363738f7f33b00",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738f7f33b00",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptTripleDESCFB(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptTripleDESCFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptTripleDESCFB(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptTripleDESCFB(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptDESCTR(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptDESCTR(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptDESCTR(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "3132333435363738267e72d0",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738267e72d0",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptDESCTR(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptDESCTR(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptDESCTR(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptDESCTR(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptTripleDESCTR(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptTripleDESCTR(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptTripleDESCTR(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "3132333435363738f7f33b00",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738f7f33b00",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptTripleDESCTR(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptTripleDESCTR(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptTripleDESCTR(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptTripleDESCTR(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptDESOFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptDESOFB(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptDESOFB(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "3132333435363738267e72d0",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738267e72d0",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("test_key"),
				ciphertext: "333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptDESOFB(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptDESOFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptDESOFB(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptDESOFB(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptTripleDESOFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		iv  []byte
		err error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv:  []byte("12345678"),
				err: nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("12345678")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscEncrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, err := encrypt.EncryptTripleDESOFB(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptTripleDESOFB(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
		reader     io.Reader
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to decrypt",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "3132333435363738f7f33b00",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: []byte("test"),
				err:       nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("invalid_length_key"),
				ciphertext: "3132333435363738f7f33b00",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("testkey_testkey_testkey_"),
				ciphertext: "333435363738",
				reader:     &testReader{b: []byte("12345678")},
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeDES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tmp := rand.Reader
			rand.Reader = tt.C().reader
			defer func() {
				rand.Reader = tmp
			}()

			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptTripleDESOFB(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptTripleDESOFB(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
	}

	type action struct {
	}

	cndEmptyPlaintext := "empty plaintext"
	actCheckNoError := "check no error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndEmptyPlaintext, "input empty bytes as plaintext")
	tb.Action(actCheckNoError, "check that there is no error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"empty plaintext",
			[]string{cndEmptyPlaintext},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"non-empty plaintext",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("testkey_testkey_testkey_"),
				plaintext: []byte("test"),
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertexts := map[string]struct{}{}

			// Check that different ciphertexts are obtained for each encryption.
			for i := 0; i < 100; i++ {
				ciphertext, err := encrypt.EncryptTripleDESOFB(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptTripleDESOFB(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)

				testutil.Diff(t, tt.C().plaintext, plaintext)

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}
