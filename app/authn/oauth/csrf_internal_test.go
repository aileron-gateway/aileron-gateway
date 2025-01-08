package oauth

import (
	"crypto/rand"
	"errors"
	"io"
	"net/url"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
)

type testReader struct {
	b   []byte
	err error
}

func (r *testReader) Read(p []byte) (int, error) {
	return copy(p, r.b), r.err
}

func TestCSRFStateGenerator_new(t *testing.T) {
	type condition struct {
		reader io.Reader
		gen    *csrfStateGenerator
	}

	type action struct {
		state      *csrfStates
		err        any // error or errorutil.Kind
		errPattern *regexp.Regexp
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndGenState := tb.Condition("state", "generate state value")
	cndGenNonce := tb.Condition("nonce", "generate nonce value")
	cndGenPKCE := tb.Condition("pkce", "generate pkce values")
	cndPKCEPlain := tb.Condition("pkce plain", "use plain method for pkce")
	cndPKCES256 := tb.Condition("pkce S256", "use S256 method for pkce")
	actCheckError := tb.Action("error", "check that an error was returned")
	actCheckNoError := tb.Action("no error", "check that there is no error")

	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success",
			[]string{cndGenState, cndGenNonce, cndGenPKCE, cndPKCES256},
			[]string{actCheckNoError},
			&condition{
				reader: &testReader{
					b: []byte(
						"1234567890123456789012345678901234567890" +
							"1234567890123456789012345678901234567890" +
							"1234"),
				},
				gen: &csrfStateGenerator{},
			},
			&action{
				state: &csrfStates{
					State:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Nonce:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Verifier:  "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					challenge: "he36zaXySaNbq2hajdl9u49JDQeeBJU427sUdDa_uCA",
					method:    "S256",
				},
				err: nil,
			},
		),
		gen(
			"S256",
			[]string{cndGenState, cndGenNonce, cndGenPKCE, cndPKCES256},
			[]string{actCheckNoError},
			&condition{
				reader: &testReader{
					b: []byte(
						"1234567890123456789012345678901234567890" +
							"1234567890123456789012345678901234567890" +
							"1234"),
				},
				gen: &csrfStateGenerator{
					method: "S256",
				},
			},
			&action{
				state: &csrfStates{
					State:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Nonce:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Verifier:  "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					challenge: "he36zaXySaNbq2hajdl9u49JDQeeBJU427sUdDa_uCA",
					method:    "S256",
				},
				err: nil,
			},
		),
		gen(
			"plain",
			[]string{cndGenState, cndGenNonce, cndGenPKCE, cndPKCEPlain},
			[]string{actCheckNoError},
			&condition{
				reader: &testReader{
					b: []byte(
						"1234567890123456789012345678901234567890" +
							"1234567890123456789012345678901234567890" +
							"1234"),
				},
				gen: &csrfStateGenerator{
					method: "plain",
				},
			},
			&action{
				state: &csrfStates{
					State:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Nonce:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Verifier:  "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					challenge: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					method:    "plain",
				},
				err: nil,
			},
		),
		gen(
			"disable state",
			[]string{cndGenNonce, cndGenPKCE, cndPKCES256},
			[]string{actCheckNoError},
			&condition{
				reader: &testReader{
					b: []byte(
						"1234567890123456789012345678901234567890" +
							"1234567890123456789012345678901234567890" +
							"1234"),
				},
				gen: &csrfStateGenerator{
					stateDisabled: true,
				},
			},
			&action{
				state: &csrfStates{
					State:     "",
					Nonce:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Verifier:  "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					challenge: "he36zaXySaNbq2hajdl9u49JDQeeBJU427sUdDa_uCA",
					method:    "S256",
				},
				err: nil,
			},
		),
		gen(
			"disable nonce",
			[]string{cndGenState, cndGenPKCE, cndPKCES256},
			[]string{actCheckNoError},
			&condition{
				reader: &testReader{
					b: []byte(
						"1234567890123456789012345678901234567890" +
							"1234567890123456789012345678901234567890" +
							"1234"),
				},
				gen: &csrfStateGenerator{
					nonceDisabled: true,
				},
			},
			&action{
				state: &csrfStates{
					State:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Nonce:     "",
					Verifier:  "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					challenge: "he36zaXySaNbq2hajdl9u49JDQeeBJU427sUdDa_uCA",
					method:    "S256",
				},
				err: nil,
			},
		),
		gen(
			"disable pkce",
			[]string{cndGenState, cndGenNonce, cndPKCES256},
			[]string{actCheckNoError},
			&condition{
				reader: &testReader{
					b: []byte(
						"1234567890123456789012345678901234567890" +
							"1234567890123456789012345678901234567890" +
							"1234"),
				},
				gen: &csrfStateGenerator{
					pkceDisabled: true,
				},
			},
			&action{
				state: &csrfStates{
					State:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Nonce:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0",
					Verifier:  "",
					challenge: "",
					method:    "",
				},
				err: nil,
			},
		),
		gen(
			"state generate error",
			[]string{cndGenState},
			[]string{actCheckError},
			&condition{
				reader: &testReader{
					err: errors.New("test error"),
				},
				gen: &csrfStateGenerator{
					stateDisabled: false,
					nonceDisabled: true,
					pkceDisabled:  true,
				},
			},
			&action{
				state:      nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. random bytes read error`),
			},
		),
		gen(
			"nonce generate error",
			[]string{cndGenNonce},
			[]string{actCheckError},
			&condition{
				reader: &testReader{
					err: errors.New("test error"),
				},
				gen: &csrfStateGenerator{
					stateDisabled: true,
					nonceDisabled: false,
					pkceDisabled:  true,
				},
			},
			&action{
				state:      nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. random bytes read error`),
			},
		),
		gen(
			"pkce generate error",
			[]string{cndGenPKCE, cndPKCES256},
			[]string{actCheckError},
			&condition{
				reader: &testReader{
					err: errors.New("test error"),
				},
				gen: &csrfStateGenerator{
					stateDisabled: true,
					nonceDisabled: true,
					pkceDisabled:  false,
				},
			},
			&action{
				state:      nil,
				err:        core.ErrCoreGenCreateComponent,
				errPattern: regexp.MustCompile(core.ErrPrefix + `failed to create component. random bytes read error`),
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

			state, err := tt.C().gen.new()
			testutil.DiffError(t, tt.A().err, tt.A().errPattern, err)
			testutil.Diff(t, tt.A().state, state, cmp.AllowUnexported(csrfStates{}))
		})
	}
}

func TestCSRFStates_set(t *testing.T) {
	type condition struct {
		states *csrfStates
		v      url.Values
	}

	type action struct {
		v url.Values
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	cndInputNil := tb.Condition("input nil", "input nil as the argument")
	cndSetState := tb.Condition("set nonce", "set non empty state value")
	cndSetNonce := tb.Condition("set challenge", "set non empty nonce value")
	cndSetChallenge := tb.Condition("set state", "set non empty challenge values")
	actCheckState := tb.Action("check state", "check the state value is set")
	actCheckNonce := tb.Action("check nonce", "check the nonce value is set")
	actCheckChallenge := tb.Action("check challenge", "check the challenge values")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"set state",
			[]string{cndSetState},
			[]string{actCheckState},
			&condition{
				states: &csrfStates{
					State: "test-state",
				},
				v: url.Values{},
			},
			&action{
				v: url.Values{
					"state": []string{"test-state"},
				},
			},
		),
		gen(
			"set nonce",
			[]string{cndSetNonce},
			[]string{actCheckNonce},
			&condition{
				states: &csrfStates{
					Nonce: "test-nonce",
				},
				v: url.Values{},
			},
			&action{
				v: url.Values{
					"nonce": []string{"test-nonce"},
				},
			},
		),
		gen(
			"set challenge",
			[]string{cndSetChallenge},
			[]string{actCheckChallenge},
			&condition{
				states: &csrfStates{
					method:    "test-method",
					challenge: "test-challenge",
				},
				v: url.Values{},
			},
			&action{
				v: url.Values{
					"code_challenge_method": []string{"test-method"},
					"code_challenge":        []string{"test-challenge"},
				},
			},
		),
		gen(
			"set all",
			[]string{cndSetState, cndSetNonce, cndSetChallenge},
			[]string{actCheckState, actCheckNonce, actCheckChallenge},
			&condition{
				states: &csrfStates{
					State:     "test-state",
					Nonce:     "test-nonce",
					Verifier:  "test-verify",
					method:    "test-method",
					challenge: "test-challenge",
				},
				v: url.Values{},
			},
			&action{
				v: url.Values{
					"state":                 []string{"test-state"},
					"nonce":                 []string{"test-nonce"},
					"code_challenge_method": []string{"test-method"},
					"code_challenge":        []string{"test-challenge"},
				},
			},
		),
		gen(
			"nil",
			[]string{cndSetState, cndSetNonce, cndSetChallenge, cndInputNil},
			[]string{},
			&condition{
				states: &csrfStates{
					State:     "test-state",
					Nonce:     "test-nonce",
					Verifier:  "test-verify",
					method:    "test-method",
					challenge: "test-challenge",
				},
				v: nil,
			},
			&action{
				v: nil,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			tt.C().states.set(tt.C().v)
			testutil.Diff(t, tt.A().v, tt.C().v)
		})
	}
}
