// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package key

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
)

func newEnvProvider(spec *v1.KeyAuthnEnvProvider) (kvs.Commander[string, credential], encoder.EncodeToStringFunc, error) {
	store := &kvs.MapKVS[string, credential]{}
	store.Open(context.Background())
	if spec == nil {
		return store, nil, nil
	}

	var encodeFunc encoder.EncodeToStringFunc
	var decodeFunc encoder.DecodeStringFunc
	if spec.Encoding == kernel.EncodingType_EncodingTypeUnknown {
		encodeFunc = func(data []byte) string { return string(data) }
		decodeFunc = func(data string) ([]byte, error) { return []byte(data), nil }
	} else {
		encodeFunc, decodeFunc = encoder.EncoderDecoder(spec.Encoding)
	}

	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, spec.KeyPrefix) {
			continue
		}
		_, value, found := strings.Cut(v, "=")
		if !found || value == "" {
			continue
		}

		if store.Exists(context.Background(), value) {
			return nil, nil, errors.New("duplicate key :" + v)
		}
		secret, err := decodeFunc(value)
		if err != nil {
			return nil, nil, errors.New("key decode error :" + v)
		}
		cred := &defaultCredential{secret: secret}
		store.Set(context.Background(), value, cred)
	}

	return store, encodeFunc, nil
}

func newFileProvider(spec *v1.KeyAuthnFileProvider) (kvs.Commander[string, credential], encoder.EncodeToStringFunc, error) {
	store := &kvs.MapKVS[string, credential]{}
	store.Open(context.Background())
	if spec == nil || len(spec.Paths) == 0 {
		return store, nil, nil
	}

	var encodeFunc encoder.EncodeToStringFunc
	var decodeFunc encoder.DecodeStringFunc
	if spec.Encoding == kernel.EncodingType_EncodingTypeUnknown {
		encodeFunc = func(data []byte) string { return string(data) }
		decodeFunc = func(data string) ([]byte, error) { return []byte(data), nil }
	} else {
		encodeFunc, decodeFunc = encoder.EncoderDecoder(spec.Encoding)
	}

	bodies, err := io.ReadFiles(false, spec.Paths...)
	if err != nil {
		return nil, nil, err
	}
	for _, b := range bodies {
		lines := bytes.Split(b, []byte("\n"))
		for _, line := range lines {
			line = bytes.Trim(line, " \n\r\t\f")
			if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
				continue
			}

			var key, attrs []byte
			if i := bytes.Index(line, []byte(":")); i == -1 {
				key = line // line was "<key>"
				// Now line is empty ""
			} else {
				key = line[:i] // line was "<key>:<attrs>"
				attrs = line[i+1:]
				// Now line is empty ""
			}

			if store.Exists(context.Background(), string(key)) {
				return nil, nil, errors.New("duplicate key :" + string(key))
			}
			secret, err := decodeFunc(string(key))
			if err != nil {
				return nil, nil, errors.New("key decode error :" + string(key))
			}
			cred := &defaultCredential{secret: secret}
			store.Set(context.Background(), string(key), cred)

			if len(attrs) > 0 {
				attrMap := map[string]any{}
				if err := encoder.UnmarshalJSON(attrs, &attrMap); err != nil {
					return nil, nil, errors.New("failed to parse attributes :" + string(key))
				}
				cred.attrs = attrMap
			}
		}
	}

	return store, encodeFunc, nil
}
