// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/url"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

var (
	_ resilience.Entry = &noopUpstream{}
	_ resilience.Entry = &lbUpstream{}
)

func TestLBUpstream_url(t *testing.T) {
	type condition struct {
		upstream *lbUpstream
	}

	type action struct {
		rawURL    string
		parsedURL *url.URL
	}

	actCheckURL := "check url string"

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	tb.Action(actCheckURL, "check the returned string")
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single upstream",
			[]string{},
			[]string{actCheckURL},
			&condition{
				upstream: &lbUpstream{
					rawURL: "http://test.com",
					parsedURL: &url.URL{
						Scheme: "http",
						Host:   "test.com",
					},
				},
			},
			&action{
				rawURL: "http://test.com",
				parsedURL: &url.URL{
					Scheme: "http",
					Host:   "test.com",
				},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			rawURL := tt.C().upstream.ID()
			testutil.Diff(t, tt.A().rawURL, rawURL)

			parsedURL := tt.C().upstream.url()
			testutil.Diff(t, tt.A().parsedURL, parsedURL)
		})
	}
}
