// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package opa

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/network"
	"github.com/open-policy-agent/opa/logging"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/disk"
	"github.com/open-policy-agent/opa/storage/inmem"
	"gopkg.in/yaml.v3"
)

func newFileStore(spec *v1.FileStore) (storage.Store, error) {
	store, err := newStore(spec.Directory)
	if err != nil {
		return nil, err
	}

	// store.NewTransaction does not return error
	// as for disk storage and in-memory storage.
	// 	- https://pkg.go.dev/github.com/open-policy-agent/opa/storage/disk#Store.NewTransaction
	// 	- https://pkg.go.dev/github.com/open-policy-agent/opa/storage/inmem#NewWithOpts
	txn, _ := store.NewTransaction(context.Background(), storage.WriteParams)
	ctx := context.Background()

	for k, v := range spec.Path {
		path, ok := storage.ParsePath(k)
		if !ok {
			store.Abort(ctx, txn)
			return nil, errors.New("invalid storage path " + k)
		}
		data, err := loadFile(v)
		if err != nil {
			store.Abort(ctx, txn)
			return nil, err
		}
		err = store.Write(ctx, txn, storage.AddOp, path, data)
		if err != nil {
			store.Abort(ctx, txn)
			return nil, err
		}
	}
	return store, store.Commit(ctx, txn)
}

func newHTTPStore(spec *v1.HTTPStore, rt http.RoundTripper) (storage.Store, error) {
	store, err := newStore(spec.Directory)
	if err != nil {
		return nil, err
	}

	// store.NewTransaction does not return error
	// as for disk storage and in-memory storage.
	// 	- https://pkg.go.dev/github.com/open-policy-agent/opa/storage/disk#Store.NewTransaction
	// 	- https://pkg.go.dev/github.com/open-policy-agent/opa/storage/inmem#NewWithOpts
	txn, _ := store.NewTransaction(context.Background(), storage.WriteParams)
	ctx := context.Background()

	for k, v := range spec.Endpoint {
		path, ok := storage.ParsePath(k)
		if !ok {
			store.Abort(ctx, txn)
			return nil, errors.New("invalid storage path " + k)
		}
		data, err := getDataHTTP(v, rt)
		if err != nil {
			store.Abort(ctx, txn)
			return nil, err
		}
		err = store.Write(ctx, txn, storage.AddOp, path, data)
		if err != nil {
			store.Abort(ctx, txn)
			return nil, err
		}
	}
	return store, store.Commit(ctx, txn)
}

func newStore(dir string) (storage.Store, error) {
	if dir == "" {
		return inmem.NewWithOpts(), nil
	}
	return disk.New(context.Background(), logging.NewNoOpLogger(), nil,
		disk.Options{
			Dir:        filepath.Clean(dir), // directory to store data inside of
			Partitions: nil,                 // data prefixes that enable efficient layout
		},
	)
}

func loadFile(path string) (any, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch ext := filepath.Ext(path); ext {
	case ".csv":
		r := csv.NewReader(f)
		return r.ReadAll()
	case ".json":
		dec := json.NewDecoder(f)
		dec.UseNumber()
		var val any
		return val, dec.Decode(&val)
	case ".yaml", ".yml":
		dec := yaml.NewDecoder(f)
		var val any
		return val, dec.Decode(&val)
	default:
		return nil, errors.New("unsupported file extension " + path)
	}
}

func getDataHTTP(endpoint string, rt http.RoundTripper) (any, error) {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	if rt == nil {
		rt = network.DefaultHTTPTransport
	}
	resp, err := rt.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.ReadAll(resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get policy from " + endpoint + " " + resp.Status)
	}

	mt, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch mt {
	case "application/csv", "text/csv":
		r := csv.NewReader(resp.Body)
		return r.ReadAll()
	case "application/json", "text/json":
		dec := json.NewDecoder(resp.Body)
		dec.UseNumber()
		var val any
		return val, dec.Decode(&val)
	case "application/yaml", "application/yml", "text/yaml", "text/yml":
		dec := yaml.NewDecoder(resp.Body)
		var val any
		return val, dec.Decode(&val)
	default:
		return nil, errors.New("unsupported media type " + mt)
	}
}
