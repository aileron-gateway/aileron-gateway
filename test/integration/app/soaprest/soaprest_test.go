// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration
// +build integration

package soaprest_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

func TestCorrectPath(t *testing.T) {
	configs := []string{"./config-path.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		respJSON, _ := os.ReadFile("./testdata/simple/ok_simple.json")
		if !equalJSON(t, []byte(respJSON), b) {
			t.Error("result not match (xml to json)")
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	xmlBytes, _ := os.ReadFile("./testdata/simple/ok_simple.xml")
	r := httptest.NewRequest(http.MethodPost, "http://test.com/soap/", strings.NewReader(string(xmlBytes)))
	r.Header.Set("Content-Type", "text/xml")
	r.Header.Set("SOAPAction", "http://example.com/")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	if !equalXML(t, []byte(xmlBytes), b) {
		t.Error("result not match (json to xml)")
	}
}

func TestWrongPath(t *testing.T) {
	configs := []string{"./config-path.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	reqXML, _ := os.ReadFile("./testdata/simple/ok_simple.xml")
	respJSON, _ := os.ReadFile("./testdata/simple/ok_simple.json")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// since requests are being sent to different paths,
		// the conversion from SOAP/XML to REST/JSON will not take place
		if !equalXML(t, []byte(reqXML), b) {
			t.Error("result not match")
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	r := httptest.NewRequest(http.MethodPost, "http://soaprest.com/wrong", strings.NewReader(string(reqXML)))
	r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)

	// similarly to the request,
	// the response will not undergo conversion from REST/JSON to SOAP/XML
	if !equalJSON(t, []byte(respJSON), b) {
		t.Error("result not match")
	}
}

func TestRayfishEmptyStringValue(t *testing.T) {
	configs := []string{"./config-rayfish-string.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		respJSON, _ := os.ReadFile("./testdata/rayfish/ok_rayfish.json")
		if !equalJSON(t, []byte(respJSON), b) {
			t.Error("result not match (xml to json)")
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	xmlBytes, _ := os.ReadFile("./testdata/rayfish/ok_rayfish.xml")
	r := httptest.NewRequest(http.MethodPost, "http://test.com/", strings.NewReader(string(xmlBytes)))
	r.Header.Set("Content-Type", "text/xml")
	r.Header.Set("SOAPAction", "http://example.com/")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	if !equalXML(t, []byte(xmlBytes), b) {
		t.Error("result not match (json to xml)")
	}
}

func TestBadgerfishEmptyStringValue(t *testing.T) {
	configs := []string{"./config-badgerfish-string.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		respJSON, _ := os.ReadFile("./testdata/badgerfish/ok_badgerfish.json")
		if !equalJSON(t, []byte(respJSON), b) {
			t.Error("result not match (xml to json)")
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	xmlBytes, _ := os.ReadFile("./testdata/badgerfish/ok_badgerfish.xml")
	r := httptest.NewRequest(http.MethodPost, "http://test.com/", strings.NewReader(string(xmlBytes)))
	r.Header.Set("Content-Type", "text/xml")
	r.Header.Set("SOAPAction", "http://example.com/")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	if !equalXML(t, []byte(xmlBytes), b) {
		t.Error("result not match (json to xml)")
	}
}

func TestJSONEncodeDecodeOptions(t *testing.T) {
	configs := []string{"./config-simple-string.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "SOAPRESTMiddleware",
		Name:       "test",
		Namespace:  "testNamespace",
	}

	m, err := api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, nil, nil, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		respJSON, _ := os.ReadFile("./testdata/simple/ok_simple_option.json")
		if !equalJSON(t, []byte(respJSON), b) {
			t.Error("result not match (xml to json)")
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(respJSON))
	})

	h := m.Middleware(handler)

	xmlBytes, _ := os.ReadFile("./testdata/simple/ok_simple_option.xml")
	r := httptest.NewRequest(http.MethodPost, "http://test.com/", strings.NewReader(string(xmlBytes)))
	r.Header.Set("Content-Type", "text/xml")
	r.Header.Set("SOAPAction", "http://example.com/")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	testutil.Diff(t, http.StatusOK, w.Result().StatusCode)
	if !equalXML(t, []byte(xmlBytes), b) {
		t.Error("result not match (json to xml)")
	}
}

func xmlTokens(decoder *xml.Decoder, end xml.EndElement) (map[string]any, error) {
	key := ""
	m := map[string]any{}
	var tokens []xml.Token
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return m, nil
			}
			return m, err
		}
		switch t := token.(type) {
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		case xml.StartElement:
			slices.SortFunc(t.Attr, func(a, b xml.Attr) int {
				if a.Name.Space+a.Name.Local > b.Name.Space+b.Name.Local {
					return 1
				}
				return -1
			})
			key = t.Name.Space + ":" + t.Name.Local
			children, err := xmlTokens(decoder, t.End())
			if err != nil {
				return nil, err
			}
			m[key] = children
		case xml.CharData:
			t = bytes.TrimSpace([]byte(t))
			if len(t) == 0 {
				continue
			}
			token = t
		case xml.EndElement:
			if t == end {
				return m, nil
			}
			v, ok := m[key]
			if ok {
				vv := v.([]xml.Token)
				m[key] = append(vv, tokens...)
			} else {
				m[key] = tokens
			}
			tokens = nil
		}
	}
}

func equalXML(t *testing.T, a, b []byte) bool {
	tokens1, err := xmlTokens(xml.NewDecoder(bytes.NewReader(a)), xml.EndElement{})
	if err != nil {
		panic(err)
	}
	tokens2, err := xmlTokens(xml.NewDecoder(bytes.NewReader(b)), xml.EndElement{})
	if err != nil {
		panic(err)
	}
	if equal := reflect.DeepEqual(tokens1, tokens2); equal {
		return true
	}
	t.Logf("XML-1: %#v\n", tokens1)
	t.Logf("XML-2: %#v\n", tokens2)
	return false
}

func equalJSON(t *testing.T, a, b []byte) bool {
	var obj1, obj2 interface{}
	if err := json.Unmarshal(a, &obj1); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &obj2); err != nil {
		panic(err)
	}
	if equal := reflect.DeepEqual(obj1, obj2); equal {
		return true
	}
	t.Logf("JSON-1: %#v\n", obj1)
	t.Logf("JSON-2: %#v\n", obj2)
	return false
}
