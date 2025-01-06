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
	_ = encrypt.EncryptFunc(encrypt.EncryptAESGCM) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptAESCBC) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptAESCFB) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptAESCTR) // Check that the function satisfies the signature.
	_ = encrypt.EncryptFunc(encrypt.EncryptAESOFB) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptAESGCM) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptAESCBC) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptAESCFB) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptAESCTR) // Check that the function satisfies the signature.
	_ = encrypt.DecryptFunc(encrypt.DecryptAESOFB) // Check that the function satisfies the signature.
)

func ExampleEncryptAESGCM() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESGCM(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 66697865642076616c756520ad6e5f0773705d3964ddc5173946355bc680d98834bef3351f2b4eb25d36fae4f2
}

func ExampleDecryptAESGCM() {
	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESGCM(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptAESGCM(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptAESCBC() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESCBC(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 66697865642076616c75652077696c6c152fece4213897d9a86bca729e6f6be4cce941d02efc6e0efa677c6ec2ba4b6a
}

func ExampleDecryptAESCBC() {
	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESCBC(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptAESCBC(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptAESCFB() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESCFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae869
}

func ExampleDecryptAESCFB() {
	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESCFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptAESCFB(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptAESCTR() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESCTR(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae85e
}

func ExampleDecryptAESCTR() {
	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESCTR(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptAESCTR(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

func ExampleEncryptAESOFB() {
	// Replace random value source for reproducibility.
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("fixed value will be returned for reproducibility"))
	defer func() {
		rand.Reader = tmp
	}()

	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESOFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(hex.EncodeToString(ciphertext))
	// Output:
	// 66697865642076616c75652077696c6c4c6e78d7a206d2db4dde7b34618dfae87e
}

func ExampleDecryptAESOFB() {
	key := []byte("1234567890123456")
	plaintext := []byte("plaintext message")

	ciphertext, err := encrypt.EncryptAESOFB(key, plaintext)
	if err != nil {
		panic("handle error here")
	}

	decrypted, err := encrypt.DecryptAESOFB(key, ciphertext)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decrypted))
	// Output:
	// plaintext message
}

type testReader struct {
	io.Reader
	b []byte
}

func (r *testReader) Read(p []byte) (n int, err error) {
	copy(p, r.b)
	if len(p) > len(r.b) {
		return len(r.b), nil
	}
	return len(p), nil
}

func TestEncryptAESGCM(t *testing.T) {
	type condition struct {
		key       []byte
		plaintext []byte
		reader    io.Reader
	}

	type action struct {
		nonce []byte
		err   error
	}

	cndValidKeyLength := "valid key length"
	cndErrReader := "use error reader"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Condition(cndErrReader, "error reading random bytes")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to encrypt/16 byte",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				nonce: []byte("123456789012"),
				err:   nil,
			},
		),
		gen(
			"success to encrypt/24 byte",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				nonce: []byte("123456789012"),
				err:   nil,
			},
		),
		gen(
			"success to encrypt/32 byte",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				nonce: []byte("123456789012"),
				err:   nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:       []byte("invalid_length_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				nonce: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
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

			ciphertext, err := encrypt.EncryptAESGCM(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().nonce))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptAESGCM(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext []byte
	}

	type action struct {
		plaintext []byte
		err       error
	}

	cndValidKeyLength := "valid key length"
	cndErrReader := "use error reader"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input valid length key")
	tb.Condition(cndErrReader, "error reading random bytes")
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
				key:        []byte("16bytes_test_key"),
				ciphertext: []byte("3132333435363738393031327911ca6ae761635f0891950b51f47ac6eafa35e6"),
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
				ciphertext: []byte("3132333435363738393031327911ca6ae761635f0891950b51f47ac6eafa35e6"),
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"ciphertext is shorter than nonce",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:        []byte("16bytes_test_key"),
				ciphertext: []byte("31323334353637383930"),
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"invalid nonce",
			[]string{},
			[]string{actCheckError},
			&condition{
				key: []byte("16bytes_test_key"),
				// invalid nonce is "223456789012" instead of "123456789012"
				ciphertext: []byte("3232333435363738393031327911ca6ae761635f0891950b51f47ac6eafa35e6"),
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertext, _ := hex.DecodeString(string(tt.C().ciphertext))
			plaintext, err := encrypt.DecryptAESGCM(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptAESGCM(t *testing.T) {
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
				key:       []byte("16bytes_test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"16 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"24 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"32 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
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
				ciphertext, err := encrypt.EncryptAESGCM(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptAESGCM(tt.C().key, ciphertext)
				testutil.Diff(t, nil, err)
				testutil.Diff(t, tt.C().plaintext, plaintext, cmpopts.EquateEmpty())

				hc := hex.EncodeToString(ciphertext)
				_, ok := ciphertexts[hc]
				testutil.Diff(t, false, ok)
				ciphertexts[hc] = struct{}{}
			}
		})
	}
}

func TestEncryptAESCBC(t *testing.T) {
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
			"success to encrypt/16 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/24 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/32 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
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
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
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

			ciphertext, err := encrypt.EncryptAESCBC(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptAESCBC(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
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
				key:        []byte("16bytes_test_key"),
				ciphertext: "3132333435363738393031323334353620eadd0f53e2bffad853c42472efc5f3",
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
				ciphertext: "3132333435363738393031323334353620eadd0f53e2bffad853c42472efc5f3",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("16bytes_test_key"),
				ciphertext: "31323334353637383930313233343536",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"decrypt with wrong key",
			[]string{cndValidKeyLength, cndWrongKey},
			[]string{actCheckError},
			&condition{
				key:        []byte("16bytes_xxxx_key"),
				ciphertext: "3132333435363738393031323334353620eadd0f53e2bffad853c42472efc5f3",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptAESCBC(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptAESCBC(t *testing.T) {
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
				key:       []byte("16bytes_test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"16 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"24 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"32 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
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
				ciphertext, err := encrypt.EncryptAESCBC(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptAESCBC(tt.C().key, ciphertext)
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

func TestEncryptAESCFB(t *testing.T) {
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
			"success to encrypt/16 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/24 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/32 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
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
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
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

			ciphertext, err := encrypt.EncryptAESCFB(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptAESCFB(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
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
				key:        []byte("16bytes_test_key"),
				ciphertext: "3132333435363738393031323334353611c2885d",
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
				ciphertext: "3132333435363738393031323334353611c2885d",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("16bytes_test_key"),
				ciphertext: "31323334353637383930",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptAESCFB(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptAESCFB(t *testing.T) {
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
				key:       []byte("16bytes_test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"16 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"24 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"32 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
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
				ciphertext, err := encrypt.EncryptAESCFB(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptAESCFB(tt.C().key, ciphertext)
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

func TestEncryptAESCTR(t *testing.T) {
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
			"success to encrypt/16 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/24 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/32 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
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
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
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

			ciphertext, err := encrypt.EncryptAESCTR(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptAESCTR(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
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
				key:        []byte("16bytes_test_key"),
				ciphertext: "3132333435363738393031323334353611c2885d",
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
				ciphertext: "3132333435363738393031323334353611c2885d",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("16bytes_test_key"),
				ciphertext: "3132333435363738393",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptAESCTR(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptAESCTR(t *testing.T) {
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
				key:       []byte("16bytes_test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"16 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"24 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"32 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
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
				ciphertext, err := encrypt.EncryptAESCTR(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptAESCTR(tt.C().key, ciphertext)
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

func TestEncryptAESOFB(t *testing.T) {
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
			"success to encrypt/16 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/24 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
				err: nil,
			},
		),
		gen(
			"success to encrypt/48 bytes",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("48bytes_test_key_000000000000000"),
				plaintext: []byte("test"),
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv:  []byte("1234567890123456"),
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
				reader:    &testReader{b: []byte("1234567890123456")},
			},
			&action{
				iv: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
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

			ciphertext, err := encrypt.EncryptAESOFB(tt.C().key, tt.C().plaintext)
			t.Log(hex.EncodeToString(ciphertext))

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, true, bytes.HasPrefix(ciphertext, tt.A().iv))
			testutil.Diff(t, false, bytes.Contains(ciphertext, tt.C().plaintext))
		})
	}
}

func TestDecryptAESOFB(t *testing.T) {
	type condition struct {
		key        []byte
		ciphertext string
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
				key:        []byte("16bytes_test_key"),
				ciphertext: "3132333435363738393031323334353611c2885d",
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
				ciphertext: "3132333435363738393031323334353611c2885d",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
		gen(
			"too short ciphertext",
			[]string{cndValidKeyLength},
			[]string{actCheckError},
			&condition{
				key:        []byte("16bytes_test_key"),
				ciphertext: "31323334353637383930",
			},
			&action{
				plaintext: nil,
				err: &er.Error{
					Package:     encrypt.ErrPkg,
					Type:        encrypt.ErrTypeAES,
					Description: encrypt.ErrDscDecrypt,
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ciphertext, _ := hex.DecodeString(tt.C().ciphertext)
			plaintext, err := encrypt.DecryptAESOFB(tt.C().key, ciphertext)

			testutil.Diff(t, tt.A().err, err, cmpopts.EquateErrors())
			testutil.Diff(t, tt.A().plaintext, plaintext)
		})
	}
}

func TestEncryptDecryptAESOFB(t *testing.T) {
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
				key:       []byte("16bytes_test_key"),
				plaintext: []byte(""),
			},
			&action{},
		),
		gen(
			"16 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("16bytes_test_key"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"24 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("24bytes_test_key_0000000"),
				plaintext: []byte("test"),
			},
			&action{},
		),
		gen(
			"32 bytes key",
			[]string{},
			[]string{actCheckNoError},
			&condition{
				key:       []byte("32bytes_test_key_000000000000000"),
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
				ciphertext, err := encrypt.EncryptAESOFB(tt.C().key, tt.C().plaintext)
				testutil.Diff(t, nil, err)

				plaintext, err := encrypt.DecryptAESOFB(tt.C().key, ciphertext)
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
