package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestNewAES(t *testing.T) {
	type condition struct {
		key    []byte
		reader io.Reader
	}

	type action struct {
		cipher cipher.Block
		iv     []byte
		err    error
	}

	cndValidKeyLength := "valid key length"
	cndErrReader := "use error reader"
	actCheckNoError := "check no error"
	actCheckError := "check error"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Condition(cndValidKeyLength, "input non-zero or non-nil password")
	tb.Condition(cndErrReader, "error reading random bytes")
	tb.Action(actCheckNoError, "check that there is no error")
	tb.Action(actCheckError, "check that there the expected error")
	table := tb.Build()

	validCipher, _ := aes.NewCipher([]byte("16bytes_test_key"))
	reader := bytes.NewReader([]byte("1234567890123456"))

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success to generate a new cipher",
			[]string{cndValidKeyLength},
			[]string{actCheckNoError},
			&condition{
				key:    []byte("16bytes_test_key"),
				reader: reader,
			},
			&action{
				cipher: validCipher,
				iv:     []byte("1234567890123456"),
				err:    nil,
			},
		),
		gen(
			"invalid key length",
			[]string{},
			[]string{actCheckError},
			&condition{
				key:    []byte("invalid_length_key"),
				reader: reader,
			},
			&action{
				cipher: nil,
				iv:     nil,
				err:    errors.New("crypto/aes: invalid key size 18"),
			},
		),
		gen(
			"err generate iv",
			[]string{cndValidKeyLength, cndErrReader},
			[]string{actCheckError},
			&condition{
				key:    []byte("16bytes_test_key"),
				reader: &testutil.ErrorReader{},
			},
			&action{
				cipher: nil,
				iv:     nil,
				err:    errors.New("rand read error"),
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

			cipher, iv, err := newAES(tt.C().key)
			if tt.A().err != nil {
				testutil.Diff(t, tt.A().err.Error(), err.Error())
			} else {
				testutil.Diff(t, nil, err)
			}

			testutil.Diff(t, true, reflect.DeepEqual(tt.A().cipher, cipher))
			testutil.Diff(t, tt.A().iv, iv)
		})
	}
}
