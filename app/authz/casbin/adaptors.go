// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package casbin

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"gopkg.in/yaml.v3"
)

var errNotImplemented = errors.New("not implemented")

type noopAdapter struct{}

func (a *noopAdapter) AddPolicies(_ string, _ string, _ [][]string) error {
	return errNotImplemented
}

func (a *noopAdapter) AddPolicy(_ string, _ string, _ []string) error {
	return errNotImplemented
}

func (a *noopAdapter) LoadPolicy(_ model.Model) error {
	return errNotImplemented
}

func (a *noopAdapter) RemoveFilteredPolicy(_ string, _ string, _ int, _ ...string) error {
	return errNotImplemented
}

func (a *noopAdapter) RemovePolicies(_ string, _ string, _ [][]string) error {
	return errNotImplemented
}

func (a *noopAdapter) RemovePolicy(_ string, _ string, _ []string) error {
	return errNotImplemented
}

func (a *noopAdapter) SavePolicy(_ model.Model) error {
	return errNotImplemented
}

func newFileAdapter(f string) (persist.Adapter, error) {
	switch ext := filepath.Ext(f); ext {
	case ".csv":
		return (&csvAdapter{
			Adapter:  &noopAdapter{},
			filePath: f,
		}), nil
	case ".json":
		return (&jsonAdapter{
			Adapter:  &noopAdapter{},
			filePath: f,
		}), nil
	case ".xml":
		return (&xmlAdapter{
			Adapter:  &noopAdapter{},
			filePath: f,
		}), nil
	case ".yaml", ".yml":
		return (&yamlAdapter{
			Adapter:  &noopAdapter{},
			filePath: f,
		}), nil
	default:
		return nil, errors.New("unsupported file extension " + f)
	}
}

type casbinRule struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"policy"`
	PType   string   `json:"pType" yaml:"pType" xml:"pType"`
	V0      string   `json:"v0" yaml:"v0" xml:"v0"`
	V1      string   `json:"v1" yaml:"v1" xml:"v1"`
	V2      string   `json:"v2" yaml:"v2" xml:"v2"`
	V3      string   `json:"v3" yaml:"v3" xml:"v3"`
	V4      string   `json:"v4" yaml:"v4" xml:"v4"`
	V5      string   `json:"v5" yaml:"v5" xml:"v5"`
	V6      string   `json:"v6" yaml:"v6" xml:"v6"`
	V7      string   `json:"v7" yaml:"v7" xml:"v7"`
	V8      string   `json:"v8" yaml:"v8" xml:"v8"`
	V9      string   `json:"v9" yaml:"v9" xml:"v9"`
}

type casbinRules struct {
	XMLName xml.Name      `xml:"policies"`
	Rules   []*casbinRule `xml:"policy"`
}

type csvAdapter struct {
	persist.Adapter
	filePath string
}

func (a *csvAdapter) LoadPolicy(model model.Model) error {
	b, err := os.ReadFile(a.filePath)
	if err != nil {
		return err
	}
	return a.loadPolicy(b, model)
}

func (a *csvAdapter) loadPolicy(b []byte, model model.Model) (err error) {
	// ClearPolicy must be called in the case for policy reloads.
	model.ClearPolicy()

	defer func() {
		err, _ = recover().(error)
	}()
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		persist.LoadPolicyLine(line, model) // Panics for invalid csv format.
	}
	return nil
}

type jsonAdapter struct {
	persist.Adapter
	filePath string
}

func (a *jsonAdapter) LoadPolicy(model model.Model) error {
	b, err := os.ReadFile(a.filePath)
	if err != nil {
		return err
	}
	return a.loadPolicy(b, model)
}

func (a *jsonAdapter) loadPolicy(b []byte, model model.Model) error {
	var rules []*casbinRule
	if err := json.Unmarshal(b, &rules); err != nil {
		return err
	}

	// ClearPolicy must be called in the case for policy reloads.
	model.ClearPolicy()

	var builder strings.Builder
	w := csv.NewWriter(&builder)
	for _, r := range rules {
		builder.Reset()
		w.Write([]string{r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5, r.V6, r.V7, r.V8, r.V9})
		w.Flush()
		persist.LoadPolicyLine(strings.TrimRight(builder.String(), ",\n"), model)
	}

	return nil
}

type yamlAdapter struct {
	persist.Adapter
	filePath string
}

func (a *yamlAdapter) LoadPolicy(model model.Model) error {
	b, err := os.ReadFile(a.filePath)
	if err != nil {
		return err
	}
	return a.loadPolicy(b, model)
}

func (a *yamlAdapter) loadPolicy(b []byte, model model.Model) error {
	var rules []*casbinRule
	if err := yaml.Unmarshal(b, &rules); err != nil {
		return err
	}

	// ClearPolicy must be called in the case for policy reloads.
	model.ClearPolicy()

	var builder strings.Builder
	w := csv.NewWriter(&builder)
	for _, r := range rules {
		builder.Reset()
		w.Write([]string{r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5, r.V6, r.V7, r.V8, r.V9})
		w.Flush()
		persist.LoadPolicyLine(strings.TrimRight(builder.String(), ",\n"), model)
	}

	return nil
}

type xmlAdapter struct {
	persist.Adapter
	filePath string
}

func (a *xmlAdapter) LoadPolicy(model model.Model) error {
	b, err := os.ReadFile(a.filePath)
	if err != nil {
		return err
	}
	return a.loadPolicy(b, model)
}

func (a *xmlAdapter) loadPolicy(b []byte, model model.Model) error {
	rules := &casbinRules{}
	if err := xml.Unmarshal(b, rules); err != nil {
		return err
	}

	// ClearPolicy must be called in the case for policy reloads.
	model.ClearPolicy()

	var builder strings.Builder
	w := csv.NewWriter(&builder)
	for _, r := range rules.Rules {
		builder.Reset()
		w.Write([]string{r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5, r.V6, r.V7, r.V8, r.V9})
		w.Flush()
		persist.LoadPolicyLine(strings.TrimRight(builder.String(), ",\n"), model)
	}

	return nil
}

type httpAdapter struct {
	persist.Adapter
	endpoint string
	rt       http.RoundTripper
}

func (a *httpAdapter) LoadPolicy(model model.Model) error {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, a.endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := a.rt.RoundTrip(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get policy from " + a.endpoint + " " + resp.Status)
	}

	mt, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch mt {
	case "application/csv", "text/csv":
		return (&csvAdapter{}).loadPolicy(b, model)
	case "application/json", "text/json":
		return (&jsonAdapter{}).loadPolicy(b, model)
	case "application/xml", "text/xml":
		return (&xmlAdapter{}).loadPolicy(b, model)
	case "application/yaml", "application/yml", "text/yaml", "text/yml":
		return (&yamlAdapter{}).loadPolicy(b, model)
	default:
		return errors.New("unsupported media type " + mt)
	}
}
