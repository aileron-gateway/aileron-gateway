package idkey

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/kvs"
)

func newEnvProvider(spec *v1.IDKeyAuthnEnvProvider) (kvs.Commander[string, credential], encoder.EncodeToStringFunc, error) {
	store := &kvs.MapKVS[string, credential]{}
	store.Open(context.Background())
	if spec == nil {
		return store, nil, nil
	}

	ids := map[int]string{}
	keys := map[int]string{}

	for _, v := range os.Environ() {
		if strings.HasPrefix(v, spec.IDPrefix) {
			vv, value, found := strings.Cut(v, "=")
			if !found || value == "" {
				continue
			}
			num, err := strconv.Atoi(strings.TrimPrefix(vv, spec.IDPrefix))
			if err != nil {
				continue
			}
			ids[num] = value
			continue
		}
		if strings.HasPrefix(v, spec.KeyPrefix) {
			vv, value, found := strings.Cut(v, "=")
			if !found || value == "" {
				continue
			}
			num, err := strconv.Atoi(strings.TrimPrefix(vv, spec.KeyPrefix))
			if err != nil {
				continue
			}
			keys[num] = value
			continue
		}
	}

	var encodeFunc encoder.EncodeToStringFunc
	var decodeFunc encoder.DecodeStringFunc
	if spec.Encoding == kernel.EncodingType_EncodingTypeUnknown {
		encodeFunc = func(data []byte) string { return string(data) }
		decodeFunc = func(data string) ([]byte, error) { return []byte(data), nil }
	} else {
		encodeFunc, decodeFunc = encoder.EncoderDecoder(spec.Encoding)
	}

	for num, id := range ids {
		key, ok := keys[num]
		if !ok {
			continue
		}

		if store.Exists(context.Background(), id) {
			return nil, nil, errors.New("duplicate id :" + id)
		}
		secret, err := decodeFunc(key)
		if err != nil {
			return nil, nil, errors.New("key decode error for id :" + id)
		}

		cred := &defaultCredential{secret: secret}
		store.Set(context.Background(), id, cred)
	}

	return store, encodeFunc, nil
}

func newFileProvider(spec *v1.IDKeyAuthnFileProvider) (kvs.Commander[string, credential], encoder.EncodeToStringFunc, error) {
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

			var id, key, attrs []byte
			if i := bytes.Index(line, []byte(":")); i == -1 {
				return nil, nil, errors.New("invalid line" + string(line))
			} else {
				id = line[:i] // line is "<id>:<key>" or "<id>:<key>:<attrs>"
				line = line[i+1:]
				if j := bytes.Index(line, []byte(":")); j == -1 {
					key = line // line was "<id>:<key>"
					// Now line is empty ""
				} else {
					key = line[:j]     // line was "<id>:<key>:<attrs>"
					attrs = line[j+1:] // line was "<id>:<key>:<attrs>"
					// Now line is empty ""
				}
			}

			if store.Exists(context.Background(), string(id)) {
				return nil, nil, errors.New("duplicate id :" + string(id))
			}

			secret, err := decodeFunc(string(key))
			if err != nil {
				return nil, nil, errors.New("key decode error for id :" + string(id))
			}

			cred := &defaultCredential{secret: secret}
			store.Set(context.Background(), string(id), cred)

			if len(attrs) > 0 {
				attrMap := map[string]any{}
				if err := encoder.UnmarshalJSON(attrs, &attrMap); err != nil {
					return nil, nil, errors.New("failed to parse attributes for id :" + string(id))
				}
				cred.attrs = attrMap
			}
		}
	}

	return store, encodeFunc, nil
}
