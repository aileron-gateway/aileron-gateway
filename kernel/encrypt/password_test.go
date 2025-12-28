// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package encrypt_test

import (
	"testing"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
)

func TestNewPasswordCrypt(t *testing.T) {
	testCases := map[string]struct {
		spec *k.PasswordCryptSpec
	}{
		"nil": {
			spec: nil,
		},
		"unknown": {
			spec: &k.PasswordCryptSpec{PasswordCrypts: nil},
		},
		"BCrypt": {
			spec: &k.PasswordCryptSpec{
				PasswordCrypts: &k.PasswordCryptSpec_BCrypt{
					BCrypt: &k.BCryptSpec{},
				},
			},
		},
		"SCrypt": {
			spec: &k.PasswordCryptSpec{
				PasswordCrypts: &k.PasswordCryptSpec_SCrypt{
					SCrypt: &k.SCryptSpec{},
				},
			},
		},
		"PBKDF2": {
			spec: &k.PasswordCryptSpec{
				PasswordCrypts: &k.PasswordCryptSpec_PBKDF2{
					PBKDF2: &k.PBKDF2Spec{},
				},
			},
		},
		"Argon2i": {
			spec: &k.PasswordCryptSpec{
				PasswordCrypts: &k.PasswordCryptSpec_Argon2I{
					Argon2I: &k.Argon2Spec{},
				},
			},
		},
		"Argon2id": {
			spec: &k.PasswordCryptSpec{
				PasswordCrypts: &k.PasswordCryptSpec_Argon2Id{
					Argon2Id: &k.Argon2Spec{},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := encrypt.NewPasswordCrypt(tc.spec)
			testutil.Diff(t, nil, err)
			if c == nil {
				return
			}
			s, err := c.Sum([]byte("password"))
			testutil.Diff(t, nil, err)
			testutil.Diff(t, true, c.Compare(s, []byte("password")) == nil)
			testutil.Diff(t, false, c.Compare(s, []byte("Password")) == nil)
		})
	}
}
