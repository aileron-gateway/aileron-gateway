// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package httpproxy

import (
	"net/http"
	"testing"

	"github.com/cespare/xxhash/v2"
)

func TestClientAddrHasher(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	r.RemoteAddr = "127.0.0.1:12345"

	testCases := map[string]struct {
		name  string
		value uint64
	}{
		"case01": {"", xxhash.Sum64String("127.0.0.1:12345")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := clientAddrHasher(tc.name)
			v := h.Hash(r)
			if v != tc.value {
				t.Error("hash value not match.", "want:", tc.value, "got:", v)
			}
		})
	}
}

func TestHeaderHasher(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	r.Header.Add("foo", "FOO")
	r.Header.Add("bar", "BAR")
	r.Header.Add("baz", "BAZ")

	testCases := map[string]struct {
		name  string
		value uint64
	}{
		"case01": {"foo", xxhash.Sum64String("FOO")},
		"case02": {"bar", xxhash.Sum64String("BAR")},
		"case03": {"baz", xxhash.Sum64String("BAZ")},
		"case04": {"alice", xxhash.Sum64String("")},
		"case05": {"bob", xxhash.Sum64String("")},
		"case06": {"", xxhash.Sum64String("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := headerHasher(tc.name)
			v := h.Hash(r)
			if v != tc.value {
				t.Error("hash value not match.", "want:", tc.value, "got:", v)
			}
		})
	}
}

func TestCookieHasher(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	r.Header.Add("Cookie", (&http.Cookie{Name: "foo", Value: "FOO"}).String())
	r.Header.Add("Cookie", (&http.Cookie{Name: "bar", Value: "BAR"}).String())
	r.Header.Add("Cookie", (&http.Cookie{Name: "baz", Value: "BAZ"}).String())

	testCases := map[string]struct {
		name  string
		value uint64
	}{
		"case01": {"foo", xxhash.Sum64String("FOO")},
		"case02": {"bar", xxhash.Sum64String("BAR")},
		"case03": {"baz", xxhash.Sum64String("BAZ")},
		"case04": {"alice", xxhash.Sum64String("")},
		"case05": {"bob", xxhash.Sum64String("")},
		"case06": {"", xxhash.Sum64String("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := cookieHasher(tc.name)
			v := h.Hash(r)
			if v != tc.value {
				t.Error("hash value not match.", "want:", tc.value, "got:", v)
			}
		})
	}
}

func TestQueryHasher(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest(http.MethodGet, "http://test.com?foo=FOO&bar=BAR&baz=BAZ", nil)

	testCases := map[string]struct {
		name  string
		value uint64
	}{
		"case01": {"foo", xxhash.Sum64String("FOO")},
		"case02": {"bar", xxhash.Sum64String("BAR")},
		"case03": {"baz", xxhash.Sum64String("BAZ")},
		"case04": {"alice", xxhash.Sum64String("")},
		"case05": {"bob", xxhash.Sum64String("")},
		"case06": {"", xxhash.Sum64String("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := queryHasher(tc.name)
			v := h.Hash(r)
			if v != tc.value {
				t.Error("hash value not match.", "want:", tc.value, "got:", v)
			}
		})
	}
}

func TestPathParamHasher(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	r.SetPathValue("foo", "FOO")
	r.SetPathValue("bar", "BAR")
	r.SetPathValue("baz", "BAZ")

	testCases := map[string]struct {
		name  string
		value uint64
	}{
		"case01": {"foo", xxhash.Sum64String("FOO")},
		"case02": {"bar", xxhash.Sum64String("BAR")},
		"case03": {"baz", xxhash.Sum64String("BAZ")},
		"case04": {"alice", xxhash.Sum64String("")},
		"case05": {"bob", xxhash.Sum64String("")},
		"case06": {"", xxhash.Sum64String("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := pathParamHasher(tc.name)
			v := h.Hash(r)
			if v != tc.value {
				t.Error("hash value not match.", "want:", tc.value, "got:", v)
			}
		})
	}
}
