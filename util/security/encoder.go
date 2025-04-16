// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package security

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
	"github.com/aileron-gateway/aileron-gateway/kernel/io"
	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
)

func NewSecureEncoder(spec *v1.SecureEncoderSpec) (*SecureEncoder, error) {
	if spec == nil {
		return nil, nil
	}

	hmac := mac.FromHashAlg(spec.HashAlg)
	size := mac.HashSize[spec.HashAlg]
	enc := encrypt.EncrypterFromType(spec.CommonKeyCryptType)
	dec := encrypt.DecrypterFromType(spec.CommonKeyCryptType)

	key, err := base64.StdEncoding.DecodeString(spec.HMACSecret)
	if err != nil {
		return nil, err
	}

	secret, err := base64.StdEncoding.DecodeString(spec.CryptSecret)
	if err != nil {
		return nil, err
	}

	s := &SecureEncoder{
		// Settings for HMAC
		hmac: hmac,
		size: size,
		key:  key,

		// Settings for common key encryption
		secret: secret,
		enc:    enc,
		dec:    dec,

		enableCompression: spec.EnableCompression,
		disableEncryption: spec.DisableEncryption,
		disableHMAC:       spec.DisableHMAC,
	}

	return s, nil
}

type SecureEncoder struct {
	hmac mac.HMACFunc
	size int
	key  []byte

	enc    encrypt.EncryptFunc
	dec    encrypt.DecryptFunc
	secret []byte

	enableCompression bool
	disableHMAC       bool
	disableEncryption bool
}

// Encode return encoded byte data.
//
//	Plaintext := Compress( data )
//	Ciphertext := Encrypt(data) HMAC(Encrypt(data))
//
// The order of Encryption and HMAC is important. See the link below.
// https://crypto.stackexchange.com/questions/202/should-we-mac-then-encrypt-or-encrypt-then-mac
func (e *SecureEncoder) Encode(b []byte) ([]byte, error) {
	data := b

	if e.enableCompression {
		var buf bytes.Buffer
		w, _ := gzip.NewWriterLevel(&buf, 6)
		if _, err := w.Write(data); err != nil {
			w.Close()
			return nil, errors.New("util/security: failed to compress data")
		}
		w.Close()
		data = buf.Bytes()
	}

	if !e.disableEncryption {
		ciphertext, err := e.enc(e.secret, data)
		if err != nil {
			return nil, errors.New("util/security: failed to encrypt data")
		}
		data = ciphertext
	}

	if !e.disableHMAC {
		digest := e.hmac(data, e.key)
		data = append(data, digest...)
	}

	return data, nil
}

func (e *SecureEncoder) Decode(b []byte) ([]byte, error) {
	data := b

	if !e.disableHMAC {
		n := len(b) - e.size
		if n < 0 {
			return nil, errors.New("util/security: invalid hash")
		}

		text, digest := b[:n], b[n:]

		d := e.hmac(text, e.key)

		if !bytes.Equal(d, digest) {
			return nil, errors.New("util/security: hashes are not matched")
		}

		data = text
	}

	if !e.disableEncryption {
		plaintext, err := e.dec(e.secret, data)
		if err != nil {
			return nil, errors.New("util/security: invalid data")
		}

		data = plaintext
	}

	if e.enableCompression {
		r, _ := gzip.NewReader(bytes.NewReader(data))
		var buf bytes.Buffer
		_, err := io.CopyBuffer(&buf, r)
		if err != nil {
			return nil, errors.New("util/security: failed to extract compressed data")
		}
		data = buf.Bytes()
	}

	return data, nil
}
