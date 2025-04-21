// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package csrf

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

type mockReader struct {
	seed []byte // Predefined seed data to be used in place of random values.
	pos  int    // Current read position within the seed data.
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.seed) {
		return 0, io.EOF // Return EOF when all seed data has been used
	}

	n = copy(p, m.seed[m.pos:]) // Copy seed data into the provided buffer `p`
	m.pos += n                  // Update the read position
	return n, nil
}

func errReader(n int64) io.ReadCloser {
	return io.NopCloser(errReaderStruct{n})
}

type errReaderStruct struct {
	n int64
}

func (r errReaderStruct) Read([]byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestCsrfToken_New(t *testing.T) {
	type condition struct {
		useMockReader bool
		seedSize      int
		hashSize      int
		mockSeed      []byte
	}

	type action struct {
		expectError   bool
		expectedToken string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"success case with mock reader",
			[]string{},
			[]string{},
			&condition{
				useMockReader: true,
				seedSize:      4,                              // Small size for testing purposes
				hashSize:      32,                             // SHA256 hash size
				mockSeed:      []byte{0x01, 0x02, 0x03, 0x04}, // Fixed seed for reproducibility
			},
			&action{
				expectError:   false,
				expectedToken: "01020304512d870cf0f656ef3e22bf03836507c790a3d69e0cb1de14281008a9ae139ed1",
			},
		),
		gen(
			"error case with standard random reader",
			[]string{},
			[]string{},
			&condition{
				useMockReader: false,
				seedSize:      32,
				hashSize:      32,
			},
			&action{
				expectError: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Set up and reset rand.Reader only if using mock reader
			if tt.C().useMockReader {
				originalReader := rand.Reader
				defer func() { rand.Reader = originalReader }()  // Restore rand.Reader after test
				rand.Reader = &mockReader{seed: tt.C().mockSeed} // Use mock reader for predictable seed
			}

			csrfToken := &csrfToken{
				secret:   []byte("some-secret-key"),
				seedSize: tt.C().seedSize,
				hashSize: tt.C().hashSize,
				hmac:     mac.FromHashAlg(kernel.HashAlg_SHA256),
			}

			token, err := csrfToken.new()
			testutil.Diff(t, tt.A().expectError, err != nil)

			if !tt.A().expectError && tt.C().useMockReader {
				testutil.Diff(t, tt.A().expectedToken, token) // Validate the generated token matches expected value
			} else {
				testutil.Diff(t, token == "", err != nil)
			}
		})
	}
}

func TestCsrfToken_Verify(t *testing.T) {
	type condition struct {
		token             string
		seedSize          int
		hashSize          int
		useGeneratedToken bool
		mockSeed          []byte
	}

	type action struct {
		expectValid bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"valid token with generated seed",
			[]string{},
			[]string{},
			&condition{
				useGeneratedToken: true,
				seedSize:          4,                              // Small size for testing purposes
				hashSize:          32,                             // SHA256 hash size
				mockSeed:          []byte{0x01, 0x02, 0x03, 0x04}, // Fixed seed for reproducibility
			},
			&action{
				expectValid: true,
			},
		),
		gen(
			"short token length",
			[]string{},
			[]string{},
			&condition{
				token:    "short-token",
				seedSize: 32,
				hashSize: 32,
			},
			&action{
				expectValid: false,
			},
		),
		gen(
			"long token length",
			[]string{},
			[]string{},
			&condition{
				token:    "long-token-value-exceeding-size",
				seedSize: 32,
				hashSize: 32,
			},
			&action{
				expectValid: false,
			},
		),
		gen(
			"incorrect token length",
			[]string{},
			[]string{},
			&condition{
				token:    "abcd", // Token shorter than required seedSize + hashSize
				seedSize: 16,
				hashSize: 16,
			},
			&action{
				expectValid: false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			// Set up and reset rand.Reader only if using a generated token with mock seed
			if tt.C().useGeneratedToken {
				originalReader := rand.Reader
				defer func() { rand.Reader = originalReader }()  // Restore rand.Reader after test
				rand.Reader = &mockReader{seed: tt.C().mockSeed} // Use mock reader for predictable seed
			}

			csrfToken := &csrfToken{
				secret:   []byte("some-secret-key"),
				seedSize: tt.C().seedSize,
				hashSize: tt.C().hashSize,
				hmac:     mac.FromHashAlg(kernel.HashAlg_SHA256),
			}

			var token string
			if tt.C().useGeneratedToken {
				// Generate a valid token with the mock seed for reproducible testing
				var err error
				token, err = csrfToken.new()
				if err != nil {
					t.Fatalf("unexpected error generating token: %v", err)
				}
			} else {
				// Use provided token from the condition
				token = tt.C().token
			}

			// Verify the token
			valid := csrfToken.verify(token)

			// Check the result
			testutil.Diff(t, tt.A().expectValid, valid)
		})
	}
}

func TestHeaderExtractor(t *testing.T) {
	type condition struct {
		header     map[string]string // Define headers as key-value pairs for direct setting
		headerName string            // Name of the header to extract
	}

	type action struct {
		expectedToken string
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"extract from header",
			[]string{},
			[]string{},
			&condition{
				header:     map[string]string{"X-CSRF-TOKEN": "header-token"}, // Define header as key-value
				headerName: "X-CSRF-TOKEN",
			},
			&action{
				expectedToken: "header-token",
			},
		),
		gen(
			"empty header",
			[]string{},
			[]string{},
			&condition{
				header:     map[string]string{}, // Empty header case
				headerName: "X-CSRF-TOKEN",
			},
			&action{
				expectedToken: "",
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Set headers directly using key-value pairs in condition
			for key, value := range tt.C().header {
				req.Header.Set(key, value)
			}

			// Debug output to confirm headers are set correctly
			fmt.Printf("Testing with headerName: %s and request headers: %v", tt.C().headerName, req.Header)

			extractor := &headerExtractor{headerName: tt.C().headerName}
			token, err := extractor.extract(req)

			// Additional debug information to confirm expected behavior
			fmt.Printf("Extracted token: %s, Expected token: %s", token, tt.A().expectedToken)

			testutil.Diff(t, nil, err)
			testutil.Diff(t, tt.A().expectedToken, token)
		})
	}
}

func TestFormExtractor(t *testing.T) {
	type condition struct {
		paramName   string
		contentType string
		body        io.Reader // Specify the request body directly
	}

	type action struct {
		expectedToken string
		expectError   bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"invalid content type",
			[]string{},
			[]string{},
			&condition{
				paramName:   "csrf_token",
				contentType: "application/json",
				body:        strings.NewReader("csrf_token=form-token"), // Form data in JSON content-type
			},
			&action{
				expectedToken: "",
				expectError:   true,
			},
		),
		gen(
			"body read error",
			[]string{},
			[]string{},
			&condition{
				paramName:   "csrf_token",
				contentType: "application/x-www-form-urlencoded",
				body:        errReader(0), // Simulate a read error
			},
			&action{
				expectedToken: "",
				expectError:   true,
			},
		),
		gen(
			"missing token in form",
			[]string{},
			[]string{},
			&condition{
				paramName:   "csrf_token",
				contentType: "application/x-www-form-urlencoded",
				body:        strings.NewReader("csrf_token="), // Empty form value
			},
			&action{
				expectedToken: "",
				expectError:   false,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", tt.C().body) // Use body directly from condition
			req.Header.Set("Content-Type", tt.C().contentType)

			extractor := &formExtractor{paramName: tt.C().paramName}
			token, err := extractor.extract(req)

			testutil.Diff(t, tt.A().expectError, err != nil)
			testutil.Diff(t, tt.A().expectedToken, token)
		})
	}
}

func TestJsonExtractor(t *testing.T) {
	type condition struct {
		jsonPath    string
		jsonContent string
		contentType string
		bodyError   bool
	}

	type action struct {
		expectedToken string
		expectError   bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"extract from json",
			[]string{},
			[]string{},
			&condition{
				jsonPath:    "csrf.token",
				jsonContent: `{"csrf": {"token": "json-token"}}`,
				contentType: "application/json",
			},
			&action{
				expectedToken: "json-token",
				expectError:   false,
			},
		),
		gen(
			"invalid content type",
			[]string{},
			[]string{},
			&condition{
				jsonPath:    "csrf.token",
				jsonContent: `{"csrf": {"token": "json-token"}}`,
				contentType: "application/x-www-form-urlencoded",
			},
			&action{
				expectedToken: "",
				expectError:   true,
			},
		),
		gen(
			"body read error",
			[]string{},
			[]string{},
			&condition{
				jsonPath:    "csrf.token",
				contentType: "application/json",
				bodyError:   true,
			},
			&action{
				expectedToken: "",
				expectError:   true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {

			var req *http.Request
			if tt.C().bodyError {
				req = httptest.NewRequest(http.MethodPost, "/", errReader(0))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.C().jsonContent))
			}

			req.Header.Set("Content-Type", tt.C().contentType)

			extractor := &jsonExtractor{jsonPath: tt.C().jsonPath}

			token, err := extractor.extract(req)
			testutil.Diff(t, tt.A().expectError, err != nil)
			testutil.Diff(t, tt.A().expectedToken, token)
		})
	}
}
