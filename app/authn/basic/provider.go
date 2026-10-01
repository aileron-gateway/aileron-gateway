// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package basic

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/encoder"
	"github.com/aileron-gateway/aileron-gateway/internal/kvs"
	"github.com/aileron-projects/go/zos"
)

func newEnvProvider(spec *v1.BasicAuthnEnvProvider) (kvs.Commander[string, credential], error) {
	store := &kvs.MapKVS[string, credential]{}
	store.Open(context.Background())
	if spec == nil {
		return store, nil
	}

	username := map[int]string{}
	password := map[int]string{}

	for _, v := range os.Environ() {
		if strings.HasPrefix(v, spec.UsernamePrefix) {
			vv, value, found := strings.Cut(v, "=")
			if !found || value == "" {
				continue
			}
			num, err := strconv.Atoi(strings.TrimPrefix(vv, spec.UsernamePrefix))
			if err != nil {
				continue
			}
			username[num] = value
			continue
		}
		if strings.HasPrefix(v, spec.PasswordPrefix) {
			vv, value, found := strings.Cut(v, "=")
			if !found || value == "" {
				continue
			}
			num, err := strconv.Atoi(strings.TrimPrefix(vv, spec.PasswordPrefix))
			if err != nil {
				continue
			}
			password[num] = value
			continue
		}
	}

	var decodeFunc encoder.DecodeStringFunc
	if spec.Encoding == kernel.EncodingType_EncodingTypeUnknown {
		decodeFunc = func(data string) ([]byte, error) { return []byte(data), nil }
	} else {
		_, decodeFunc = encoder.EncoderDecoder(spec.Encoding)
	}

	for num, un := range username {
		pw, ok := password[num]
		if !ok {
			continue
		}

		if store.Exists(context.Background(), un) {
			return nil, errors.New("duplicate user :" + un)
		}
		secret, err := decodeFunc(pw)
		if err != nil {
			return nil, errors.New("password decode error for user :" + un)
		}
		cred := &defaultCredential{secret: secret}
		store.Set(context.Background(), un, cred)
	}

	return store, nil
}

func newFileProvider(spec *v1.BasicAuthnFileProvider) (kvs.Commander[string, credential], error) {
	store := &kvs.MapKVS[string, credential]{}
	store.Open(context.Background())
	if spec == nil || len(spec.Paths) == 0 {
		return store, nil
	}

	var decodeFunc encoder.DecodeStringFunc
	if spec.Encoding == kernel.EncodingType_EncodingTypeUnknown {
		decodeFunc = func(data string) ([]byte, error) { return []byte(data), nil }
	} else {
		_, decodeFunc = encoder.EncoderDecoder(spec.Encoding)
	}

	bodies, err := zos.ReadFiles(false, spec.Paths...)
	if err != nil {
		return nil, err
	}
	for _, b := range bodies {
		lines := bytes.Split(b, []byte("\n"))
		for _, line := range lines {
			line = bytes.Trim(line, " \n\r\t\f")
			if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
				continue
			}

			var un, pw, attrs []byte
			if i := bytes.Index(line, []byte(":")); i == -1 {
				return nil, errors.New("invalid line" + string(line))
			} else {
				un = line[:i]     // line was "<username>:<password>" or "<username>:<password>:<attrs>"
				line = line[i+1:] // Now line is "<password>" or "<password>:<attrs>"
				if j := bytes.Index(line, []byte(":")); j == -1 {
					pw = line // line was "<username>:<password>"
					// Now line is empty ""
				} else {
					pw = line[:j]      // line was "<username>:<password>:<attrs>"
					attrs = line[j+1:] // line was "<username>:<password>:<attrs>"
					// Now line is empty ""
				}
			}

			if store.Exists(context.Background(), string(un)) {
				return nil, errors.New("duplicate user :" + string(un))
			}

			pw, err := decodeFunc(string(pw))
			if err != nil {
				return nil, errors.New("password decode error for user :" + string(un))
			}
			cred := &defaultCredential{secret: pw}
			store.Set(context.Background(), string(un), cred)

			if len(attrs) > 0 {
				attrMap := map[string]any{}
				if err := encoder.UnmarshalJSON(attrs, &attrMap); err != nil {
					return nil, errors.New("failed to parse attributes for user :" + string(un))
				}
				cred.attrs = attrMap
			}
		}
	}

	return store, nil
}
