// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

//go:build integration

package oauth_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	apps "github.com/aileron-gateway/aileron-gateway/app"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

var reAuthenticationRequired = errors.New("authn/oauth: authentication required")

func checkAuthSuccess(t *testing.T, h apps.AuthenticationHandler, token string) {

	t.Helper()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/get", nil)
	r1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	_, result1, shouldReturn1, err1 := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthSucceeded, result1)
	testutil.Diff(t, false, shouldReturn1)
	testutil.Diff(t, nil, err1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b1))                        // Body not written in the AuthenticationHandler.

	r2 := httptest.NewRequest(http.MethodPost, "http://test.com/post", bytes.NewReader([]byte("dummy body")))
	r2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	_, result2, shouldReturn2, err2 := h.ServeAuthn(w2, r2)
	testutil.Diff(t, apps.AuthSucceeded, result2)
	testutil.Diff(t, false, shouldReturn2)
	testutil.Diff(t, nil, err2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b2))                        // Body not written in the AuthenticationHandler.

}

func checkAuthFailure(t *testing.T, h apps.AuthenticationHandler, token string) {

	t.Helper()

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/get", nil)
	r1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	_, result1, shouldReturn1, err1 := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthFailed, result1)
	testutil.Diff(t, true, shouldReturn1)
	testutil.Diff(t, "authn/oauth: authentication required", err1.Error())
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b1))                        // Body not written in the AuthenticationHandler.

	r2 := httptest.NewRequest(http.MethodPost, "http://test.com/post", bytes.NewReader([]byte("dummy body")))
	r2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	_, result2, shouldReturn2, err2 := h.ServeAuthn(w2, r2)
	testutil.Diff(t, apps.AuthFailed, result2)
	testutil.Diff(t, true, shouldReturn2)
	testutil.Diff(t, "authn/oauth: authentication required", err2.Error())
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b2))                        // Body not written in the AuthenticationHandler.

}

func TestAlgES(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-alg-ES.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "ESXXX",
			"typ": "JWT",
			"kid": "test-key-XXX"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtES256 := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTI1NiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.ykFIi72ctKHIlKfd9o7NY2wgPGix6D0PcCMc5n-kI_Sbm_PiJSXPYW8g0kfY4ym-B8D6wkCB7BuHNrwQ3mNEww"
	jwtES384 := "eyJhbGciOiJFUzM4NCIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTM4NCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ._5Px5nMJV5ss7l04OSgqcyQroQWywMQsQT_xpCko8jFxqPE8IeZKPDribEVJylgYGpVyAMjKpP41gc9E4W87wcMRe6ND2bl1LwTs8IZCZRD6LFkcHR_TuKgovcxJ2sxJ"
	jwtES512 := "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTUxMiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.AY-4W9NYDzPO-hFYrvnsLZ9lANLKnvLkrf9JSuft3EeoswmCxGseJGn7gxQRgM90RnhGZ6fdBDHMXnq1yw5mFXZwAD2cbSDsaVLw3ZCOsSR_g0t8-dP0UzX82x6yH8ecW4dd-fXbv5bbb0uazHviSyp61iIW1_5KVEVJjqqU8OLQMHL-"

	// alg is rewritten to "none".
	jwtNone256 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMjU2In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.ykFIi72ctKHIlKfd9o7NY2wgPGix6D0PcCMc5n-kI_Sbm_PiJSXPYW8g0kfY4ym-B8D6wkCB7BuHNrwQ3mNEww"
	jwtNone384 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMzg0In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ._5Px5nMJV5ss7l04OSgqcyQroQWywMQsQT_xpCko8jFxqPE8IeZKPDribEVJylgYGpVyAMjKpP41gc9E4W87wcMRe6ND2bl1LwTs8IZCZRD6LFkcHR_TuKgovcxJ2sxJ"
	jwtNone512 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktNTEyIn0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.AY-4W9NYDzPO-hFYrvnsLZ9lANLKnvLkrf9JSuft3EeoswmCxGseJGn7gxQRgM90RnhGZ6fdBDHMXnq1yw5mFXZwAD2cbSDsaVLw3ZCOsSR_g0t8-dP0UzX82x6yH8ecW4dd-fXbv5bbb0uazHviSyp61iIW1_5KVEVJjqqU8OLQMHL-"

	checkAuthSuccess(t, h, jwtES256)
	checkAuthSuccess(t, h, jwtES384)
	checkAuthSuccess(t, h, jwtES512)
	checkAuthFailure(t, h, jwtNone256)
	checkAuthFailure(t, h, jwtNone384)
	checkAuthFailure(t, h, jwtNone512)

}

func TestAlgHS(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-alg-HS.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "HSXXX",
			"typ": "JWT",
			"kid": "test-key-XXX"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtHS256 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTI1NiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.OT3tLMGKLdYc3c7kxjSRBtJoE_QnxPDWzGSz42ArqSY"
	jwtHS384 := "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTM4NCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.ELp046guA_MPtzPgJt3SflUpuckr92YCTBZnMheZNEfapuA8AkLlp-P3O8O0L_RU"
	jwtHS512 := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTUxMiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.sZmeQyAjIRpYJee1VHwXz5kLI_QKGf6HQKKvQ4TfEd4lmF33FV7zkPgxBLs_xb1kyFOT1agQ5vnzMDUubmWZCA"

	// alg is rewritten to "none".
	jwtNone256 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMjU2In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.OT3tLMGKLdYc3c7kxjSRBtJoE_QnxPDWzGSz42ArqSY"
	jwtNone384 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMzg0In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.ELp046guA_MPtzPgJt3SflUpuckr92YCTBZnMheZNEfapuA8AkLlp-P3O8O0L_RU"
	jwtNone512 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktNTEyIn0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.sZmeQyAjIRpYJee1VHwXz5kLI_QKGf6HQKKvQ4TfEd4lmF33FV7zkPgxBLs_xb1kyFOT1agQ5vnzMDUubmWZCA"

	checkAuthSuccess(t, h, jwtHS256)
	checkAuthSuccess(t, h, jwtHS384)
	checkAuthSuccess(t, h, jwtHS512)
	checkAuthFailure(t, h, jwtNone256)
	checkAuthFailure(t, h, jwtNone384)
	checkAuthFailure(t, h, jwtNone512)

}

func TestAlgPS(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-alg-PS.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "PSXXX",
			"typ": "JWT",
			"kid": "test-key-XXX"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtPS256 := "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTI1NiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.o4ZaT1EvwvtbWvvPz2Ve7dZgUA_BthD3X7JKpH5EEbV2-9nH5wqvqlFM6SwOJTthIT1UMzkrBpPuWumT7yyKFQTf0WbR_zf2oiioXaupc4v1ZjQqrWEPTimEPiZgjrXH9aU0yTCwy2md16-os6_36kE8FIAA-2F6uA7kufaAYGpSIaY001kwz8ChD7BKgujEequ5Z3G2_vjNy3LVvVfx5WN9yaEmb17mpqvKgNyqR_ajKs9D217Oz98TNNXQFxlyGJ0UspSuu845tgl2mVokgeD4VKjp16Asfkg9Kcew461aJcH4oLB7YCeFhNT-iuOhm3zYqJEVLSd-JdvyD1s7Dw"
	jwtPS384 := "eyJhbGciOiJQUzM4NCIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTM4NCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.fSgZrTEO8NxGGPZ_eFNNlSjQ2f6Od3eDJW2YDzr7UvQCcGVTjh7SHpKrykZqdr45Q7gARSd1K3eVHSbSZmN7u5rvK4YznwR4XDkXsD9NNzaueYN2IoihpeUyXrKwEvBAhwCkpGmnLJY5Nlf5brhPzfaoA9cWPV05JTBNB1nbB_bczore71PKUpbEXdyGZ6eiG4gaqaYl8wXhZxfkjxsofNe8RMZY5vbBGNnEl0ogvz_ro6J6O-bBfiANXOMSIg5Y1Apar-mDJJRxLX43O8NxYXsl2ABha063eFcHmuft4P89orJG2rjblqi3t2FSNBHNQuqoWrrO30iICII1qkCAxg"
	jwtPS512 := "eyJhbGciOiJQUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTUxMiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.RkaZEWtrX-3TyCjJZHZxq76jgUtUi2JWxwWsByQXt7kMMplt4s8SeiIAQHdUUbgB234GhHQjkkoGuMiYsbs7i0HJPnc9ypE0Pn3gU4d1i47LPgyxJLj82scSX1lRHKHO-F8Q0sk407ioLEumrFmT7pDK262D4DXbWA5H1wYhZgTn1_0fmyXmR6qrJ3A97kzIQByRLIVdWMtoGNBa-63gXltO_Owlh11mlpsdIIxFOKNXT2LMtLEWBVwI7_cbQgO3W_V0Te3jKWclQ1jSedKwQNNilHqoCf5jhVfXOml1OLlkbjYx-iJrFBHGJBmlGeLNbtq8308tjhf07q1un9V_PQ"

	// alg is rewritten to "none".
	jwtNone256 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMjU2In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.o4ZaT1EvwvtbWvvPz2Ve7dZgUA_BthD3X7JKpH5EEbV2-9nH5wqvqlFM6SwOJTthIT1UMzkrBpPuWumT7yyKFQTf0WbR_zf2oiioXaupc4v1ZjQqrWEPTimEPiZgjrXH9aU0yTCwy2md16-os6_36kE8FIAA-2F6uA7kufaAYGpSIaY001kwz8ChD7BKgujEequ5Z3G2_vjNy3LVvVfx5WN9yaEmb17mpqvKgNyqR_ajKs9D217Oz98TNNXQFxlyGJ0UspSuu845tgl2mVokgeD4VKjp16Asfkg9Kcew461aJcH4oLB7YCeFhNT-iuOhm3zYqJEVLSd-JdvyD1s7Dw"
	jwtNone384 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMzg0In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.fSgZrTEO8NxGGPZ_eFNNlSjQ2f6Od3eDJW2YDzr7UvQCcGVTjh7SHpKrykZqdr45Q7gARSd1K3eVHSbSZmN7u5rvK4YznwR4XDkXsD9NNzaueYN2IoihpeUyXrKwEvBAhwCkpGmnLJY5Nlf5brhPzfaoA9cWPV05JTBNB1nbB_bczore71PKUpbEXdyGZ6eiG4gaqaYl8wXhZxfkjxsofNe8RMZY5vbBGNnEl0ogvz_ro6J6O-bBfiANXOMSIg5Y1Apar-mDJJRxLX43O8NxYXsl2ABha063eFcHmuft4P89orJG2rjblqi3t2FSNBHNQuqoWrrO30iICII1qkCAxg"
	jwtNone512 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktNTEyIn0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.RkaZEWtrX-3TyCjJZHZxq76jgUtUi2JWxwWsByQXt7kMMplt4s8SeiIAQHdUUbgB234GhHQjkkoGuMiYsbs7i0HJPnc9ypE0Pn3gU4d1i47LPgyxJLj82scSX1lRHKHO-F8Q0sk407ioLEumrFmT7pDK262D4DXbWA5H1wYhZgTn1_0fmyXmR6qrJ3A97kzIQByRLIVdWMtoGNBa-63gXltO_Owlh11mlpsdIIxFOKNXT2LMtLEWBVwI7_cbQgO3W_V0Te3jKWclQ1jSedKwQNNilHqoCf5jhVfXOml1OLlkbjYx-iJrFBHGJBmlGeLNbtq8308tjhf07q1un9V_PQ"

	checkAuthSuccess(t, h, jwtPS256)
	checkAuthSuccess(t, h, jwtPS384)
	checkAuthSuccess(t, h, jwtPS512)
	checkAuthFailure(t, h, jwtNone256)
	checkAuthFailure(t, h, jwtNone384)
	checkAuthFailure(t, h, jwtNone512)

}

func TestAlgRS(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-alg-RS.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RSXXX",
			"typ": "JWT",
			"kid": "test-key-XXX"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtRS256 := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTI1NiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.RGzT5l8OYIyWBUCsbTgO2vC2qf1tWEBgLiPTBNylrAqIN4_G-L4MHkjpqVFl8PGIzOKWLXmia7T1mj9M5V4vDuxHq-K1S848_xgmoUU7sJKCMH_8btIM872gbocWwP4kaV4p66allBJoe3yBvuy4qLPfn2RtZ_24WFm3FaAYyIOZRqVW1JoNJs0l1SjSCmA8-0WeFEkk_WQAN0YUafhOMx1S7sFnrqBpn3mnxI6gCz5_oa7ML4wRXcXeeRpWROLKhrreyuYPJ-n2BOfepVaLChUHA7EeOkCzG3wCMTiABpJ2S_FOcj3VDd7DyRP9fkWH_1CuCoCU99yjim7wQB08Rg"
	jwtRS384 := "eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTM4NCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.iVj62-Q5jpr9Hwlit7ZzbIrcvLlBQsdrc35e5txHDpcbV0UH67B-sQt0R5fXEnyXvi0Ja5wfZB7B254ISFFLo2C98p1sXgpYKxehWpaA8TNeCHJqX05dLLTkbYdoUFBmMMZ-YUXoItUKCWKfIJF8H7mYWHtYB99EURP2n7_t--J0rKSTPYsY-n8mJ33tOE44nUWpSZx1gTCBw0X4fp0nWoBGbdmHPmc3A8YXxm_icLuY2fBktO5KZFVfWSLAxlpWljUwgU3i3fJDpBanDMLirq1AsjRS5bEy9r4naSQjfGqRFWVHjP7nJJB85Wj7aC2iRW6vj4XtvyqAIhwxneZACg"
	jwtRS512 := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5LTUxMiJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.EHcxWWk-uNSzCECHe9uLt6o4S8lWHb0ZOC_TzXsShT4KUDf79K3xb_tNrI3f2c8uJak90IKjN7g0anI6yJBGd3L7Kq8LkmbZbjFr6E3MoIe6H295MRxfwhn8VFv39qVZHxH5YBYS6TlmiltnHx-8xq-0E6j7YCcT9R1ClKJK9wMvidd05eeUOUmJCaLAKceNwrhGEE24lI9XR7ft4EnBugPKVYFs3dXnYBfQHJdKlg8gnTuPRe9nlt4R3RmFBfGbgj1kxtyzYG6our92LZwRQcerZ_yzBlVOJyyGN4LF0fLQYixS48hNYnhddYgqyYJ1lYhYEXhxSjVIvEg90GJSpw"

	// alg is rewritten to "none".
	jwtNone256 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMjU2In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.RGzT5l8OYIyWBUCsbTgO2vC2qf1tWEBgLiPTBNylrAqIN4_G-L4MHkjpqVFl8PGIzOKWLXmia7T1mj9M5V4vDuxHq-K1S848_xgmoUU7sJKCMH_8btIM872gbocWwP4kaV4p66allBJoe3yBvuy4qLPfn2RtZ_24WFm3FaAYyIOZRqVW1JoNJs0l1SjSCmA8-0WeFEkk_WQAN0YUafhOMx1S7sFnrqBpn3mnxI6gCz5_oa7ML4wRXcXeeRpWROLKhrreyuYPJ-n2BOfepVaLChUHA7EeOkCzG3wCMTiABpJ2S_FOcj3VDd7DyRP9fkWH_1CuCoCU99yjim7wQB08Rg"
	jwtNone384 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktMzg0In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.iVj62-Q5jpr9Hwlit7ZzbIrcvLlBQsdrc35e5txHDpcbV0UH67B-sQt0R5fXEnyXvi0Ja5wfZB7B254ISFFLo2C98p1sXgpYKxehWpaA8TNeCHJqX05dLLTkbYdoUFBmMMZ-YUXoItUKCWKfIJF8H7mYWHtYB99EURP2n7_t--J0rKSTPYsY-n8mJ33tOE44nUWpSZx1gTCBw0X4fp0nWoBGbdmHPmc3A8YXxm_icLuY2fBktO5KZFVfWSLAxlpWljUwgU3i3fJDpBanDMLirq1AsjRS5bEy9r4naSQjfGqRFWVHjP7nJJB85Wj7aC2iRW6vj4XtvyqAIhwxneZACg"
	jwtNone512 := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIiwia2lkIjoidGVzdC1rZXktNTEyIn0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.EHcxWWk-uNSzCECHe9uLt6o4S8lWHb0ZOC_TzXsShT4KUDf79K3xb_tNrI3f2c8uJak90IKjN7g0anI6yJBGd3L7Kq8LkmbZbjFr6E3MoIe6H295MRxfwhn8VFv39qVZHxH5YBYS6TlmiltnHx-8xq-0E6j7YCcT9R1ClKJK9wMvidd05eeUOUmJCaLAKceNwrhGEE24lI9XR7ft4EnBugPKVYFs3dXnYBfQHJdKlg8gnTuPRe9nlt4R3RmFBfGbgj1kxtyzYG6our92LZwRQcerZ_yzBlVOJyyGN4LF0fLQYixS48hNYnhddYgqyYJ1lYhYEXhxSjVIvEg90GJSpw"

	checkAuthSuccess(t, h, jwtRS256)
	checkAuthSuccess(t, h, jwtRS384)
	checkAuthSuccess(t, h, jwtRS512)
	checkAuthFailure(t, h, jwtNone256)
	checkAuthFailure(t, h, jwtNone384)
	checkAuthFailure(t, h, jwtNone512)

}

func TestContextQuery(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-context-query.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// Context found for query foo=bar.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/?foo=bar", nil)
	r1.Header.Set("Authorization", "Bearer "+jwtToken)
	w1 := httptest.NewRecorder()
	_, result1, shouldReturn1, err1 := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthSucceeded, result1)
	testutil.Diff(t, false, shouldReturn1)
	testutil.Diff(t, nil, err1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b1))                        // Body not written in the AuthenticationHandler.

	// Context not found for header foo=bar.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r2.Header.Set("Authorization", "Bearer "+jwtToken)
	r2.Header.Set("foo", "bar")
	w2 := httptest.NewRecorder()
	_, result2, shouldReturn2, err2 := h.ServeAuthn(w2, r2)
	testutil.Diff(t, apps.AuthContinue, result2) // Try next authentication handler.
	testutil.Diff(t, false, shouldReturn2)       // Try next authentication handler.
	testutil.Diff(t, nil, err2)                  // Try next authentication handler.
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b2))                        // Body not written in the AuthenticationHandler.

	// Context not found without query and header.
	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r3.Header.Set("Authorization", "Bearer "+jwtToken)
	w3 := httptest.NewRecorder()
	_, result3, shouldReturn3, err3 := h.ServeAuthn(w3, r3)
	testutil.Diff(t, apps.AuthContinue, result3) // Try next authentication handler.
	testutil.Diff(t, false, shouldReturn3)       // Try next authentication handler.
	testutil.Diff(t, nil, err3)                  // Try next authentication handler.
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b3))                        // Body not written in the AuthenticationHandler.

}

func TestContextHeader(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-context-header.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// Context found for header foo=bar.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r1.Header.Set("Authorization", "Bearer "+jwtToken)
	r1.Header.Set("foo", "bar")
	w1 := httptest.NewRecorder()
	_, result1, shouldReturn1, err1 := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthSucceeded, result1)
	testutil.Diff(t, false, shouldReturn1)
	testutil.Diff(t, nil, err1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b1))                        // Body not written in the AuthenticationHandler.

	// Context not found for query foo=bar.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/?foo=bar", nil)
	r2.Header.Set("Authorization", "Bearer "+jwtToken)
	w2 := httptest.NewRecorder()
	_, result2, shouldReturn2, err2 := h.ServeAuthn(w2, r2)
	testutil.Diff(t, apps.AuthContinue, result2) // Try next authentication handler.
	testutil.Diff(t, false, shouldReturn2)       // Try next authentication handler.
	testutil.Diff(t, nil, err2)                  // Try next authentication handler.
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b2))                        // Body not written in the AuthenticationHandler.

	// Context not found without query and header.
	r3 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r3.Header.Set("Authorization", "Bearer "+jwtToken)
	w3 := httptest.NewRecorder()
	_, result3, shouldReturn3, err3 := h.ServeAuthn(w3, r3)
	testutil.Diff(t, apps.AuthContinue, result3) // Try next authentication handler.
	testutil.Diff(t, false, shouldReturn3)       // Try next authentication handler.
	testutil.Diff(t, nil, err3)                  // Try next authentication handler.
	b3, _ := io.ReadAll(w3.Result().Body)
	testutil.Diff(t, http.StatusOK, w3.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b3))                        // Body not written in the AuthenticationHandler.

}

func TestIntrospection(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-introspection.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	introspectionServer := &http.Server{
		Addr: ":12525",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.PostFormValue("token") {
			case "test-opaque-token-ok":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"active":true}`))
			default:
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid","error_description":"invalid"}`))
			}
		}),
	}
	go func() { introspectionServer.ListenAndServe() }()
	time.Sleep(time.Second) // Wait a little until server starts.
	defer introspectionServer.Close()

	// Use opaque token rather than JWT.
	opaqueTokenOk := "test-opaque-token-ok"
	opaqueTokenNg := "test-opaque-token-ng"

	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r1.Header.Set("Authorization", "Bearer "+opaqueTokenOk)
	r1.Header.Set("foo", "bar")
	w1 := httptest.NewRecorder()
	_, result1, shouldReturn1, err1 := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthSucceeded, result1)
	testutil.Diff(t, false, shouldReturn1)
	testutil.Diff(t, nil, err1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b1))                        // Body not written in the AuthenticationHandler.

	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r2.Header.Set("Authorization", "Bearer "+opaqueTokenNg)
	w2 := httptest.NewRecorder()
	_, result2, shouldReturn2, err2 := h.ServeAuthn(w2, r2)
	testutil.Diff(t, apps.AuthFailed, result2)
	testutil.Diff(t, true, shouldReturn2)
	errPattern := regexp.MustCompile(core.ErrPrefix + `failed to token introspection`)
	testutil.DiffError(t, apps.ErrAppAuthnIntrospection, errPattern, err2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b2))                        // Body not written in the AuthenticationHandler.

}

func TestHeaderKey(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-header-key.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// Token in "X-Access-Token" is allowed.
	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r1.Header.Set("X-Access-Token", jwtToken)
	w1 := httptest.NewRecorder()
	_, result1, shouldReturn1, err1 := h.ServeAuthn(w1, r1)
	testutil.Diff(t, apps.AuthSucceeded, result1)
	testutil.Diff(t, false, shouldReturn1)
	testutil.Diff(t, nil, err1)
	b1, _ := io.ReadAll(w1.Result().Body)
	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b1))                        // Body not written in the AuthenticationHandler.

	// Token in "Authorization" is denied.
	r2 := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
	r2.Header.Set("Authorization", "Bearer "+jwtToken)
	w2 := httptest.NewRecorder()
	_, result2, shouldReturn2, err2 := h.ServeAuthn(w2, r2)
	testutil.Diff(t, apps.AuthContinue, result2)
	testutil.Diff(t, false, shouldReturn2)
	testutil.Diff(t, nil, err2)
	b2, _ := io.ReadAll(w2.Result().Body)
	testutil.Diff(t, http.StatusOK, w2.Result().StatusCode) // Status not written in the AuthenticationHandler.
	testutil.Diff(t, "", string(b2))                        // Body not written in the AuthenticationHandler.

}

func TestValidationDefault(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-validation-default.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	// jwtToken is the valid token with the header and payload above.
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// jwtWrongAlgToken is the same as jwtToken except for alg. signed by HS256 with "test-secret".
	jwtWrongAlgToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.b4ZFXTKPmcLLT2zbnolCjHGvjC4g1plasL5R0N23K1c"
	// jwtWrongKidToken is the same as jwtToken except for kid. kid="wong-key".
	jwtWrongKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Indyb25nLWtleSJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.blkzTXsUyEOZId3Mq2T-M8vBx8X8-MfxT_R5cAmS8PjzXVfQgQLuzT4eXwQoT0tMXiOp6bSEYPwk-7uVBZzQ9R8zH2wmgR5L69EQp2xnSII6yEeCNCXNrBF86tHbc1tv3_W7blQPSKFs2lEU8gBPqwHsC3aftClcOCtUtME3FIEhSfpu3XvxW2iG3YBOIhvS9vEeO-dmaCpI52UYPrpFIOYCa-RkykQ_9QGCGU-GbOKfSZLX2dWgbvCBivkqH11JI2UuD_9r91iIsQeGb0EH8D8HN73SZ-p3xS8qWC9zDiCk3g_8ml1Qe5Es6L_vH1zp-QfqKDauaZYpHpxX4s4_iA"
	// jwtWrongExpiredToken is the same as jwtToken except for exp. exp=987654321.
	jwtWrongExpiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMSwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.Idotjnq-hw8yOLEfegpHeGPBcrDqeeuESwvBX7OMq8DR44Wcm4h5mg5DWf582Q2mQTVgHMHPuwct6HSNdCxc5Gx4HHqE4CVOq54dv99duYVluQNReDD31G9Uc6Lgnb9zIW85bsg8Q-TnpWO9aO_RqePRVrUJ7wdRdozGipNYSwtBQfLQzORjdLnz9u6JnbLa9eQiEyihL6rAaJz3DqX5Epft2UaqO2b0s31_y6ar0IHV3nL5lR_FkaSrjOvFxA8SRFJZk7BvnA0v2xR4Jqu0MFN7ZGdo12PfPisG2n3oIyaFEtDbETswCYl4zL6lKpWN7RQZbBkg0-4DCSiXJqvckQ"
	// jwtWrongIssToken is the same as jwtToken except for iss. iss="http://wrong.provider.com/".
	jwtWrongIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly93cm9uZy5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.NQXXYsJ7XYvEOLV8dqVq33jPvyBw5OALBELObO9F0Z8YsJB0hUF5lwqoezn745ue_MqrTW4xlFBc2ctGIa49FZdgl9ptquTolYjDP4rv8m1B0-2zu4OBx_dxG_cICFGBmYGJtK9_wInyKKZTr7ilKQW74yM01dELBzPlVTcZsfOLDwnQfEPFae3G958feporeAt5ME7-wnvzfM2AdhV9HozB7Cmvy7z3eUcO5LGgPp5RrKKJupq36HvFZuFKnc67RmHnwaF23FZf_Fv7EJSzeNeyrpd9lBo3hJsZTokPKPAv_ebUkRNpNADfr89Yxq_oww_DaccIUAN75wPBrR_vYA"
	// jwtWrongAudToken is the same as jwtToken except for aud. aud=["wrong-audience"].
	jwtWrongAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsid3JvbmctYXVkaWVuY2UiXX0.Tny_Xsj9r8TED_Xooa6EYwjCm4m_KX-6MKKU9Ix96wMc3svCUCpKNmMGsZ5qzeW-37SCH-Hx1_ZMm1MnaDp0DiBIRKgPRxfNAKzXC58DhZTEkNZEDRRzN5SpzK_nwTdjkYcE02qewb2wLHzkZ7Bhu1Q06M3JYd42krLjEScWu9T5rr8ErWotEHtYUaY0mF5ez79byUTtlJ0bMUYoVMnbwvwcqMtRA5xEiLCWF7ad39xgSvY2aAGsZuwZrKPeUz2yrahZiZQstR18gqClIR5Kk63yUbtX86iBQtpL4kE0SrtblhXG0yKwjzrBIofdW5qOoIrZQnXt1FwmHgN2jhzNbQ"

	// jwtNoKidToken is the same as jwtToken except for kid. kid not exists.
	jwtNoKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.J86YrOgvnvB5zv-ngOI4AP_F7vBTMEp3231oixD14BmNWU9S54SIueinfSxFkARYp-4kAHeIv-368FI4__gHphOYBnhSUPHvr5TR-QGblqtjOnfjmgJ3NvUkfNdaAKaTnpRy9telAAzkPUEzAxJKig8gYVyv8fK_guaoucX9v366hvld6jLx-8QSE8czqpriNhV3i2vRJPkF2hiRpUE5Ci6UtGRh3nzicao_4uktm7Kg1KcawQjtLDIKnThyMtSkKYbHv99mnFPF11cR857L0H5RSx9gHCgPtyyoUXGTiW0lQXTPsHK6WhGEzOr32969hGFR2eJyBoZxOZo9w1666A"
	// jwtNoExpToken is the same as jwtToken except for exp. exp not exists.
	jwtNoExpToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.Ao7bMa_sLsOn5LscGTRiRBprA-IuQj3RYRiQS2IOYNOLeY3tdx0YGs1vgW92hE7GyuljpHepY7Ppr23PNFOjw9bX0DnjiBycAINEKE7lna-p4YGtIi4GDpSATHR5kUHo6ck38eeu_OXJWzrn2ahaJ3lW-ZxHy0Th09tuUHLNuLjOEJGVvQTOYjnxYOawp5g5mC1b6SvtJoRd9bZsVLpWFKjLVI1xcDfuaBQLTFCD7sbHiv-WkNGtwjlOHV5h7HqZ0zHK8vuBznZtog5dkP8jOGXYoprtIbUjpOsXsKB85LFkwVdhGGI8eOahz_g_FyxYXhcyIt7TemWlRDqw7cuRww"
	// jwtNoIssToken is the same as jwtToken except for iss. iss not exists.
	jwtNoIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.GqD9K2v74Dp48epBthaRTnSG7cRbmS2qbHR4v69H3zV_kptoWn4IYqPLYhhN2t1FiPiASOax1ZlFGkJaD5WcM4S15UzJ0lz1KpV1o2zL00mKW38qz6XA75mGG64VtBSjFnyk2NjSAC-cV6Cy3BCwnMBW798FHwmFI1QpgMVbs-LM7hwROwENIH5YzOvO2oN-7kDamtOLj8g5hXLN_LKwxDYlWqd2cz5FWDSlTPaaV1OtlKIIsYX4UySx2m-5UgbdROYmxNRdWxVVpGPPad0qbEybbNHPea2NA7BHB9b1I7iNkDLUwysIwnTl0V3wzDn49rc4rRa64BMQ9VOWotvzNA"
	// jwtNoAudToken is the same as jwtToken except for aud. aud not exists.
	jwtNoAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8ifQ.IKWx1IvYy27SE7iy1FLnqN-y5CX860XQLRkLhhlvPFqOyK7bsbBbWZkKBuhcfFkJmqHErsSi90R0oJ0Obdujp6EUbSgWepVz9T1oiGLMTzkewAzZy7tmIBhk9K9Bbs0VhPEB0gyZ5CC167flGdEq-fuSr7o-3l1D5fpXUjtRZXidb3vuLwmMiW1jbpTymub4Y0J1B1z58dt_HMis98AvW7ToOSND39JF3B716HJm-mGOkGP7722X-nXVcBGEygePN5xM8jb5Joc72PVOPip_uZLvj9hB64pUMDXYYAMnihq8R7BXW85sBoU-VLpuBOBE01_h5vLp8vohEkcQA_TeAw"

	// jwtValidNbf is the same as jwtToken except for nbf addition. nbf=987654321 (valid).
	jwtValidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.lNwivTGxUHl1Gzie-YnzMGiByxu9NpITQEMoUohPI1t19dJe1tFm2XeT3xC5K3Lm7UcI_5uQ7ysLQxCyqSabT027HIZjrrvdPbsM-kbxrRjld_8IPKJUwIilKqN4562rtn-dMbquHlzNZQX4MxfZNXD1Dr2ViXhBd6TklhDOrI8S9AF_mEaNZTV_hOjcxw5WAzAXA9EW-QXGNaIOFvCdk5hqZyo1QWMU2SOSNDALnRhtamdMYmt6wk0UlCpwIeJeL1zjp0OgJUOnPTOjdPTHFs6yjDJHO2zrwtx02YB4lT928SdMRDaqMzBdvI2nHxfYwR-fhaKysmOnS2x98S_nVg"
	// jwtInvalidNbf is the same as jwtToken except for nbf addition. nbf=9876543210 (invalid).
	jwtInvalidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.HimYFFhJbg1QVw0Y3uHM2qFmFqSFGOLixB-R4NzGMWIsXHzFT5XOm2PLRlpWaGLDvAPlVMs50A-7qD6BPwVql1N6uHtRmrk4Cd1FBQjIi5sCiyIO_RZ4dJaP9Gcn3hgimHO-e5MQ9E2vOLIdAZizyAV8-_4gNpCEq_1tt0vMPKd1nXdClRVL4twCY2nYAuC4DZOVfTU30LRlM2SBB6bHXGGHHkN8hAcV8G9g4jaEelmW6-sWN785hRujw9a5vVnCVbgC20X0JoH6bBhUN51zTQ1TJvxDmmXgVEWh6ig6elZHbTypM441CHM6kkAnQ_8f97uwQgKAJn6R32Ghdm4Qzw"
	// jwtValidIat is the same as jwtToken except for iat addition. iat=987654321 (valid).
	jwtValidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.eubgpgd9ll7dtgSSE3RVnYZ_4VkSliGhGVXTaNVToCBmF-pABR_JXXrA49aQA_w3gaRmPA5IvCFx6wRBam9blX63JfM3-CZ-NVpnMO44nhNfiA4vE9lJo_uPg-ud3cWteC5QHTVjQOMBj2KUNqjmJTZoSvrcOhbd_FHtNdocu_jK5VH5_7QObgxE1i4fkbnhBaNRUQ0AFVyJs8ui0zUYhty1Vch_87ipCBsJ0VgT03Rb7xOhX2oRb5Zcn4BJH2uQ5OoFrRy63u_43js8NqF-M59vm0wrNXwWgHb3KiVCqzjDvSI8qPL4rZWp3wKxEzeg75Wqn4zsk9L0PWehMORk9A"
	// jwtInvalidIat is the same as jwtToken except for iat addition. iat=9876543210 (invalid).
	jwtInvalidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.mx9coQPAFsgL12zVjCFcIIZ4oVPFlG7vWG63pL1doC4KiWVY-h_Ingxlf19us-P5rbQqH2_ylKtKhxpQNNIE2MKb-Fn3b4PZUSMNb7SKEFogENy5Lcf4puLd7_fqHkwv6cBxH_WcvD3wzWZIEkR1TlCDV1xdezrPFJgBZUyKEGf4CJ2uF_8eyn5knCVqok3TRZZk_sXPo2T1VQBibT7Sy8jSYmtKB1NR1AxAtK96G_oRnFQjjNEpRez6bZONykLzKLp73jnyDIzIcbcMrrcoG7QfGTbgGujGZqEq85_hgZjAwRcaA9k16nkzhc4hFEMMu0V4lUww9eOqoUMQAGaqFA"

	checkAuthSuccess(t, h, jwtToken)

	checkAuthFailure(t, h, jwtWrongAlgToken)
	checkAuthFailure(t, h, jwtWrongKidToken)
	checkAuthFailure(t, h, jwtWrongExpiredToken)
	checkAuthFailure(t, h, jwtWrongIssToken)
	checkAuthFailure(t, h, jwtWrongAudToken)

	checkAuthFailure(t, h, jwtNoKidToken)
	checkAuthFailure(t, h, jwtNoExpToken)
	checkAuthFailure(t, h, jwtNoIssToken)
	checkAuthFailure(t, h, jwtNoAudToken)

	checkAuthSuccess(t, h, jwtValidNbf)
	checkAuthFailure(t, h, jwtInvalidNbf)
	checkAuthSuccess(t, h, jwtValidIat)
	checkAuthFailure(t, h, jwtInvalidIat)

}

func TestValidationIgnoreAud(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-validation-ignore-aud.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	// jwtToken is the valid token with the header and payload above.
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// jwtWrongAlgToken is the same as jwtToken except for alg. signed by HS256 with "test-secret".
	jwtWrongAlgToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.b4ZFXTKPmcLLT2zbnolCjHGvjC4g1plasL5R0N23K1c"
	// jwtWrongKidToken is the same as jwtToken except for kid. kid="wong-key".
	jwtWrongKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Indyb25nLWtleSJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.blkzTXsUyEOZId3Mq2T-M8vBx8X8-MfxT_R5cAmS8PjzXVfQgQLuzT4eXwQoT0tMXiOp6bSEYPwk-7uVBZzQ9R8zH2wmgR5L69EQp2xnSII6yEeCNCXNrBF86tHbc1tv3_W7blQPSKFs2lEU8gBPqwHsC3aftClcOCtUtME3FIEhSfpu3XvxW2iG3YBOIhvS9vEeO-dmaCpI52UYPrpFIOYCa-RkykQ_9QGCGU-GbOKfSZLX2dWgbvCBivkqH11JI2UuD_9r91iIsQeGb0EH8D8HN73SZ-p3xS8qWC9zDiCk3g_8ml1Qe5Es6L_vH1zp-QfqKDauaZYpHpxX4s4_iA"
	// jwtWrongExpiredToken is the same as jwtToken except for exp. exp=987654321.
	jwtWrongExpiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMSwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.Idotjnq-hw8yOLEfegpHeGPBcrDqeeuESwvBX7OMq8DR44Wcm4h5mg5DWf582Q2mQTVgHMHPuwct6HSNdCxc5Gx4HHqE4CVOq54dv99duYVluQNReDD31G9Uc6Lgnb9zIW85bsg8Q-TnpWO9aO_RqePRVrUJ7wdRdozGipNYSwtBQfLQzORjdLnz9u6JnbLa9eQiEyihL6rAaJz3DqX5Epft2UaqO2b0s31_y6ar0IHV3nL5lR_FkaSrjOvFxA8SRFJZk7BvnA0v2xR4Jqu0MFN7ZGdo12PfPisG2n3oIyaFEtDbETswCYl4zL6lKpWN7RQZbBkg0-4DCSiXJqvckQ"
	// jwtWrongIssToken is the same as jwtToken except for iss. iss="http://wrong.provider.com/".
	jwtWrongIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly93cm9uZy5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.NQXXYsJ7XYvEOLV8dqVq33jPvyBw5OALBELObO9F0Z8YsJB0hUF5lwqoezn745ue_MqrTW4xlFBc2ctGIa49FZdgl9ptquTolYjDP4rv8m1B0-2zu4OBx_dxG_cICFGBmYGJtK9_wInyKKZTr7ilKQW74yM01dELBzPlVTcZsfOLDwnQfEPFae3G958feporeAt5ME7-wnvzfM2AdhV9HozB7Cmvy7z3eUcO5LGgPp5RrKKJupq36HvFZuFKnc67RmHnwaF23FZf_Fv7EJSzeNeyrpd9lBo3hJsZTokPKPAv_ebUkRNpNADfr89Yxq_oww_DaccIUAN75wPBrR_vYA"
	// jwtWrongAudToken is the same as jwtToken except for aud. aud=["wrong-audience"].
	jwtWrongAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsid3JvbmctYXVkaWVuY2UiXX0.Tny_Xsj9r8TED_Xooa6EYwjCm4m_KX-6MKKU9Ix96wMc3svCUCpKNmMGsZ5qzeW-37SCH-Hx1_ZMm1MnaDp0DiBIRKgPRxfNAKzXC58DhZTEkNZEDRRzN5SpzK_nwTdjkYcE02qewb2wLHzkZ7Bhu1Q06M3JYd42krLjEScWu9T5rr8ErWotEHtYUaY0mF5ez79byUTtlJ0bMUYoVMnbwvwcqMtRA5xEiLCWF7ad39xgSvY2aAGsZuwZrKPeUz2yrahZiZQstR18gqClIR5Kk63yUbtX86iBQtpL4kE0SrtblhXG0yKwjzrBIofdW5qOoIrZQnXt1FwmHgN2jhzNbQ"

	// jwtNoKidToken is the same as jwtToken except for kid. kid not exists.
	jwtNoKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.J86YrOgvnvB5zv-ngOI4AP_F7vBTMEp3231oixD14BmNWU9S54SIueinfSxFkARYp-4kAHeIv-368FI4__gHphOYBnhSUPHvr5TR-QGblqtjOnfjmgJ3NvUkfNdaAKaTnpRy9telAAzkPUEzAxJKig8gYVyv8fK_guaoucX9v366hvld6jLx-8QSE8czqpriNhV3i2vRJPkF2hiRpUE5Ci6UtGRh3nzicao_4uktm7Kg1KcawQjtLDIKnThyMtSkKYbHv99mnFPF11cR857L0H5RSx9gHCgPtyyoUXGTiW0lQXTPsHK6WhGEzOr32969hGFR2eJyBoZxOZo9w1666A"
	// jwtNoExpToken is the same as jwtToken except for exp. exp not exists.
	jwtNoExpToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.Ao7bMa_sLsOn5LscGTRiRBprA-IuQj3RYRiQS2IOYNOLeY3tdx0YGs1vgW92hE7GyuljpHepY7Ppr23PNFOjw9bX0DnjiBycAINEKE7lna-p4YGtIi4GDpSATHR5kUHo6ck38eeu_OXJWzrn2ahaJ3lW-ZxHy0Th09tuUHLNuLjOEJGVvQTOYjnxYOawp5g5mC1b6SvtJoRd9bZsVLpWFKjLVI1xcDfuaBQLTFCD7sbHiv-WkNGtwjlOHV5h7HqZ0zHK8vuBznZtog5dkP8jOGXYoprtIbUjpOsXsKB85LFkwVdhGGI8eOahz_g_FyxYXhcyIt7TemWlRDqw7cuRww"
	// jwtNoIssToken is the same as jwtToken except for iss. iss not exists.
	jwtNoIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.GqD9K2v74Dp48epBthaRTnSG7cRbmS2qbHR4v69H3zV_kptoWn4IYqPLYhhN2t1FiPiASOax1ZlFGkJaD5WcM4S15UzJ0lz1KpV1o2zL00mKW38qz6XA75mGG64VtBSjFnyk2NjSAC-cV6Cy3BCwnMBW798FHwmFI1QpgMVbs-LM7hwROwENIH5YzOvO2oN-7kDamtOLj8g5hXLN_LKwxDYlWqd2cz5FWDSlTPaaV1OtlKIIsYX4UySx2m-5UgbdROYmxNRdWxVVpGPPad0qbEybbNHPea2NA7BHB9b1I7iNkDLUwysIwnTl0V3wzDn49rc4rRa64BMQ9VOWotvzNA"
	// jwtNoAudToken is the same as jwtToken except for aud. aud not exists.
	jwtNoAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8ifQ.IKWx1IvYy27SE7iy1FLnqN-y5CX860XQLRkLhhlvPFqOyK7bsbBbWZkKBuhcfFkJmqHErsSi90R0oJ0Obdujp6EUbSgWepVz9T1oiGLMTzkewAzZy7tmIBhk9K9Bbs0VhPEB0gyZ5CC167flGdEq-fuSr7o-3l1D5fpXUjtRZXidb3vuLwmMiW1jbpTymub4Y0J1B1z58dt_HMis98AvW7ToOSND39JF3B716HJm-mGOkGP7722X-nXVcBGEygePN5xM8jb5Joc72PVOPip_uZLvj9hB64pUMDXYYAMnihq8R7BXW85sBoU-VLpuBOBE01_h5vLp8vohEkcQA_TeAw"

	// jwtValidNbf is the same as jwtToken except for nbf addition. nbf=987654321 (valid).
	jwtValidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.lNwivTGxUHl1Gzie-YnzMGiByxu9NpITQEMoUohPI1t19dJe1tFm2XeT3xC5K3Lm7UcI_5uQ7ysLQxCyqSabT027HIZjrrvdPbsM-kbxrRjld_8IPKJUwIilKqN4562rtn-dMbquHlzNZQX4MxfZNXD1Dr2ViXhBd6TklhDOrI8S9AF_mEaNZTV_hOjcxw5WAzAXA9EW-QXGNaIOFvCdk5hqZyo1QWMU2SOSNDALnRhtamdMYmt6wk0UlCpwIeJeL1zjp0OgJUOnPTOjdPTHFs6yjDJHO2zrwtx02YB4lT928SdMRDaqMzBdvI2nHxfYwR-fhaKysmOnS2x98S_nVg"
	// jwtInvalidNbf is the same as jwtToken except for nbf addition. nbf=9876543210 (invalid).
	jwtInvalidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.HimYFFhJbg1QVw0Y3uHM2qFmFqSFGOLixB-R4NzGMWIsXHzFT5XOm2PLRlpWaGLDvAPlVMs50A-7qD6BPwVql1N6uHtRmrk4Cd1FBQjIi5sCiyIO_RZ4dJaP9Gcn3hgimHO-e5MQ9E2vOLIdAZizyAV8-_4gNpCEq_1tt0vMPKd1nXdClRVL4twCY2nYAuC4DZOVfTU30LRlM2SBB6bHXGGHHkN8hAcV8G9g4jaEelmW6-sWN785hRujw9a5vVnCVbgC20X0JoH6bBhUN51zTQ1TJvxDmmXgVEWh6ig6elZHbTypM441CHM6kkAnQ_8f97uwQgKAJn6R32Ghdm4Qzw"
	// jwtValidIat is the same as jwtToken except for iat addition. iat=987654321 (valid).
	jwtValidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.eubgpgd9ll7dtgSSE3RVnYZ_4VkSliGhGVXTaNVToCBmF-pABR_JXXrA49aQA_w3gaRmPA5IvCFx6wRBam9blX63JfM3-CZ-NVpnMO44nhNfiA4vE9lJo_uPg-ud3cWteC5QHTVjQOMBj2KUNqjmJTZoSvrcOhbd_FHtNdocu_jK5VH5_7QObgxE1i4fkbnhBaNRUQ0AFVyJs8ui0zUYhty1Vch_87ipCBsJ0VgT03Rb7xOhX2oRb5Zcn4BJH2uQ5OoFrRy63u_43js8NqF-M59vm0wrNXwWgHb3KiVCqzjDvSI8qPL4rZWp3wKxEzeg75Wqn4zsk9L0PWehMORk9A"
	// jwtInvalidIat is the same as jwtToken except for iat addition. iat=9876543210 (invalid).
	jwtInvalidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.mx9coQPAFsgL12zVjCFcIIZ4oVPFlG7vWG63pL1doC4KiWVY-h_Ingxlf19us-P5rbQqH2_ylKtKhxpQNNIE2MKb-Fn3b4PZUSMNb7SKEFogENy5Lcf4puLd7_fqHkwv6cBxH_WcvD3wzWZIEkR1TlCDV1xdezrPFJgBZUyKEGf4CJ2uF_8eyn5knCVqok3TRZZk_sXPo2T1VQBibT7Sy8jSYmtKB1NR1AxAtK96G_oRnFQjjNEpRez6bZONykLzKLp73jnyDIzIcbcMrrcoG7QfGTbgGujGZqEq85_hgZjAwRcaA9k16nkzhc4hFEMMu0V4lUww9eOqoUMQAGaqFA"

	checkAuthSuccess(t, h, jwtToken)

	checkAuthFailure(t, h, jwtWrongAlgToken)
	checkAuthFailure(t, h, jwtWrongKidToken)
	checkAuthFailure(t, h, jwtWrongExpiredToken)
	checkAuthFailure(t, h, jwtWrongIssToken)
	checkAuthSuccess(t, h, jwtWrongAudToken) // Invalid aud but should success.

	checkAuthFailure(t, h, jwtNoKidToken)
	checkAuthFailure(t, h, jwtNoExpToken)
	checkAuthFailure(t, h, jwtNoIssToken)
	checkAuthSuccess(t, h, jwtNoAudToken) // No aud but should success.

	checkAuthSuccess(t, h, jwtValidNbf)
	checkAuthFailure(t, h, jwtInvalidNbf)
	checkAuthSuccess(t, h, jwtValidIat)
	checkAuthFailure(t, h, jwtInvalidIat)

}

func TestValidationIgnoreIss(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-validation-ignore-iss.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	// jwtToken is the valid token with the header and payload above.
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// jwtWrongAlgToken is the same as jwtToken except for alg. signed by HS256 with "test-secret".
	jwtWrongAlgToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.b4ZFXTKPmcLLT2zbnolCjHGvjC4g1plasL5R0N23K1c"
	// jwtWrongKidToken is the same as jwtToken except for kid. kid="wong-key".
	jwtWrongKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Indyb25nLWtleSJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.blkzTXsUyEOZId3Mq2T-M8vBx8X8-MfxT_R5cAmS8PjzXVfQgQLuzT4eXwQoT0tMXiOp6bSEYPwk-7uVBZzQ9R8zH2wmgR5L69EQp2xnSII6yEeCNCXNrBF86tHbc1tv3_W7blQPSKFs2lEU8gBPqwHsC3aftClcOCtUtME3FIEhSfpu3XvxW2iG3YBOIhvS9vEeO-dmaCpI52UYPrpFIOYCa-RkykQ_9QGCGU-GbOKfSZLX2dWgbvCBivkqH11JI2UuD_9r91iIsQeGb0EH8D8HN73SZ-p3xS8qWC9zDiCk3g_8ml1Qe5Es6L_vH1zp-QfqKDauaZYpHpxX4s4_iA"
	// jwtWrongExpiredToken is the same as jwtToken except for exp. exp=987654321.
	jwtWrongExpiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMSwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.Idotjnq-hw8yOLEfegpHeGPBcrDqeeuESwvBX7OMq8DR44Wcm4h5mg5DWf582Q2mQTVgHMHPuwct6HSNdCxc5Gx4HHqE4CVOq54dv99duYVluQNReDD31G9Uc6Lgnb9zIW85bsg8Q-TnpWO9aO_RqePRVrUJ7wdRdozGipNYSwtBQfLQzORjdLnz9u6JnbLa9eQiEyihL6rAaJz3DqX5Epft2UaqO2b0s31_y6ar0IHV3nL5lR_FkaSrjOvFxA8SRFJZk7BvnA0v2xR4Jqu0MFN7ZGdo12PfPisG2n3oIyaFEtDbETswCYl4zL6lKpWN7RQZbBkg0-4DCSiXJqvckQ"
	// jwtWrongIssToken is the same as jwtToken except for iss. iss="http://wrong.provider.com/".
	jwtWrongIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly93cm9uZy5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.NQXXYsJ7XYvEOLV8dqVq33jPvyBw5OALBELObO9F0Z8YsJB0hUF5lwqoezn745ue_MqrTW4xlFBc2ctGIa49FZdgl9ptquTolYjDP4rv8m1B0-2zu4OBx_dxG_cICFGBmYGJtK9_wInyKKZTr7ilKQW74yM01dELBzPlVTcZsfOLDwnQfEPFae3G958feporeAt5ME7-wnvzfM2AdhV9HozB7Cmvy7z3eUcO5LGgPp5RrKKJupq36HvFZuFKnc67RmHnwaF23FZf_Fv7EJSzeNeyrpd9lBo3hJsZTokPKPAv_ebUkRNpNADfr89Yxq_oww_DaccIUAN75wPBrR_vYA"
	// jwtWrongAudToken is the same as jwtToken except for aud. aud=["wrong-audience"].
	jwtWrongAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsid3JvbmctYXVkaWVuY2UiXX0.Tny_Xsj9r8TED_Xooa6EYwjCm4m_KX-6MKKU9Ix96wMc3svCUCpKNmMGsZ5qzeW-37SCH-Hx1_ZMm1MnaDp0DiBIRKgPRxfNAKzXC58DhZTEkNZEDRRzN5SpzK_nwTdjkYcE02qewb2wLHzkZ7Bhu1Q06M3JYd42krLjEScWu9T5rr8ErWotEHtYUaY0mF5ez79byUTtlJ0bMUYoVMnbwvwcqMtRA5xEiLCWF7ad39xgSvY2aAGsZuwZrKPeUz2yrahZiZQstR18gqClIR5Kk63yUbtX86iBQtpL4kE0SrtblhXG0yKwjzrBIofdW5qOoIrZQnXt1FwmHgN2jhzNbQ"

	// jwtNoKidToken is the same as jwtToken except for kid. kid not exists.
	jwtNoKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.J86YrOgvnvB5zv-ngOI4AP_F7vBTMEp3231oixD14BmNWU9S54SIueinfSxFkARYp-4kAHeIv-368FI4__gHphOYBnhSUPHvr5TR-QGblqtjOnfjmgJ3NvUkfNdaAKaTnpRy9telAAzkPUEzAxJKig8gYVyv8fK_guaoucX9v366hvld6jLx-8QSE8czqpriNhV3i2vRJPkF2hiRpUE5Ci6UtGRh3nzicao_4uktm7Kg1KcawQjtLDIKnThyMtSkKYbHv99mnFPF11cR857L0H5RSx9gHCgPtyyoUXGTiW0lQXTPsHK6WhGEzOr32969hGFR2eJyBoZxOZo9w1666A"
	// jwtNoExpToken is the same as jwtToken except for exp. exp not exists.
	jwtNoExpToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.Ao7bMa_sLsOn5LscGTRiRBprA-IuQj3RYRiQS2IOYNOLeY3tdx0YGs1vgW92hE7GyuljpHepY7Ppr23PNFOjw9bX0DnjiBycAINEKE7lna-p4YGtIi4GDpSATHR5kUHo6ck38eeu_OXJWzrn2ahaJ3lW-ZxHy0Th09tuUHLNuLjOEJGVvQTOYjnxYOawp5g5mC1b6SvtJoRd9bZsVLpWFKjLVI1xcDfuaBQLTFCD7sbHiv-WkNGtwjlOHV5h7HqZ0zHK8vuBznZtog5dkP8jOGXYoprtIbUjpOsXsKB85LFkwVdhGGI8eOahz_g_FyxYXhcyIt7TemWlRDqw7cuRww"
	// jwtNoIssToken is the same as jwtToken except for iss. iss not exists.
	jwtNoIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.GqD9K2v74Dp48epBthaRTnSG7cRbmS2qbHR4v69H3zV_kptoWn4IYqPLYhhN2t1FiPiASOax1ZlFGkJaD5WcM4S15UzJ0lz1KpV1o2zL00mKW38qz6XA75mGG64VtBSjFnyk2NjSAC-cV6Cy3BCwnMBW798FHwmFI1QpgMVbs-LM7hwROwENIH5YzOvO2oN-7kDamtOLj8g5hXLN_LKwxDYlWqd2cz5FWDSlTPaaV1OtlKIIsYX4UySx2m-5UgbdROYmxNRdWxVVpGPPad0qbEybbNHPea2NA7BHB9b1I7iNkDLUwysIwnTl0V3wzDn49rc4rRa64BMQ9VOWotvzNA"
	// jwtNoAudToken is the same as jwtToken except for aud. aud not exists.
	jwtNoAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8ifQ.IKWx1IvYy27SE7iy1FLnqN-y5CX860XQLRkLhhlvPFqOyK7bsbBbWZkKBuhcfFkJmqHErsSi90R0oJ0Obdujp6EUbSgWepVz9T1oiGLMTzkewAzZy7tmIBhk9K9Bbs0VhPEB0gyZ5CC167flGdEq-fuSr7o-3l1D5fpXUjtRZXidb3vuLwmMiW1jbpTymub4Y0J1B1z58dt_HMis98AvW7ToOSND39JF3B716HJm-mGOkGP7722X-nXVcBGEygePN5xM8jb5Joc72PVOPip_uZLvj9hB64pUMDXYYAMnihq8R7BXW85sBoU-VLpuBOBE01_h5vLp8vohEkcQA_TeAw"

	// jwtValidNbf is the same as jwtToken except for nbf addition. nbf=987654321 (valid).
	jwtValidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.lNwivTGxUHl1Gzie-YnzMGiByxu9NpITQEMoUohPI1t19dJe1tFm2XeT3xC5K3Lm7UcI_5uQ7ysLQxCyqSabT027HIZjrrvdPbsM-kbxrRjld_8IPKJUwIilKqN4562rtn-dMbquHlzNZQX4MxfZNXD1Dr2ViXhBd6TklhDOrI8S9AF_mEaNZTV_hOjcxw5WAzAXA9EW-QXGNaIOFvCdk5hqZyo1QWMU2SOSNDALnRhtamdMYmt6wk0UlCpwIeJeL1zjp0OgJUOnPTOjdPTHFs6yjDJHO2zrwtx02YB4lT928SdMRDaqMzBdvI2nHxfYwR-fhaKysmOnS2x98S_nVg"
	// jwtInvalidNbf is the same as jwtToken except for nbf addition. nbf=9876543210 (invalid).
	jwtInvalidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.HimYFFhJbg1QVw0Y3uHM2qFmFqSFGOLixB-R4NzGMWIsXHzFT5XOm2PLRlpWaGLDvAPlVMs50A-7qD6BPwVql1N6uHtRmrk4Cd1FBQjIi5sCiyIO_RZ4dJaP9Gcn3hgimHO-e5MQ9E2vOLIdAZizyAV8-_4gNpCEq_1tt0vMPKd1nXdClRVL4twCY2nYAuC4DZOVfTU30LRlM2SBB6bHXGGHHkN8hAcV8G9g4jaEelmW6-sWN785hRujw9a5vVnCVbgC20X0JoH6bBhUN51zTQ1TJvxDmmXgVEWh6ig6elZHbTypM441CHM6kkAnQ_8f97uwQgKAJn6R32Ghdm4Qzw"
	// jwtValidIat is the same as jwtToken except for iat addition. iat=987654321 (valid).
	jwtValidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.eubgpgd9ll7dtgSSE3RVnYZ_4VkSliGhGVXTaNVToCBmF-pABR_JXXrA49aQA_w3gaRmPA5IvCFx6wRBam9blX63JfM3-CZ-NVpnMO44nhNfiA4vE9lJo_uPg-ud3cWteC5QHTVjQOMBj2KUNqjmJTZoSvrcOhbd_FHtNdocu_jK5VH5_7QObgxE1i4fkbnhBaNRUQ0AFVyJs8ui0zUYhty1Vch_87ipCBsJ0VgT03Rb7xOhX2oRb5Zcn4BJH2uQ5OoFrRy63u_43js8NqF-M59vm0wrNXwWgHb3KiVCqzjDvSI8qPL4rZWp3wKxEzeg75Wqn4zsk9L0PWehMORk9A"
	// jwtInvalidIat is the same as jwtToken except for iat addition. iat=9876543210 (invalid).
	jwtInvalidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.mx9coQPAFsgL12zVjCFcIIZ4oVPFlG7vWG63pL1doC4KiWVY-h_Ingxlf19us-P5rbQqH2_ylKtKhxpQNNIE2MKb-Fn3b4PZUSMNb7SKEFogENy5Lcf4puLd7_fqHkwv6cBxH_WcvD3wzWZIEkR1TlCDV1xdezrPFJgBZUyKEGf4CJ2uF_8eyn5knCVqok3TRZZk_sXPo2T1VQBibT7Sy8jSYmtKB1NR1AxAtK96G_oRnFQjjNEpRez6bZONykLzKLp73jnyDIzIcbcMrrcoG7QfGTbgGujGZqEq85_hgZjAwRcaA9k16nkzhc4hFEMMu0V4lUww9eOqoUMQAGaqFA"

	checkAuthSuccess(t, h, jwtToken)

	checkAuthFailure(t, h, jwtWrongAlgToken)
	checkAuthFailure(t, h, jwtWrongKidToken)
	checkAuthFailure(t, h, jwtWrongExpiredToken)
	checkAuthSuccess(t, h, jwtWrongIssToken) // Invalid iss but should success.
	checkAuthFailure(t, h, jwtWrongAudToken)

	checkAuthFailure(t, h, jwtNoKidToken)
	checkAuthFailure(t, h, jwtNoExpToken)
	checkAuthSuccess(t, h, jwtNoIssToken) // No iss but should success.
	checkAuthFailure(t, h, jwtNoAudToken)

	checkAuthSuccess(t, h, jwtValidNbf)
	checkAuthFailure(t, h, jwtInvalidNbf)
	checkAuthSuccess(t, h, jwtValidIat)
	checkAuthFailure(t, h, jwtInvalidIat)

}

func TestValidationOptionalExp(t *testing.T) {

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, []string{"./config-validation-optional-exp.yaml"})
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "OAuthAuthenticationHandler",
	}
	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
	testutil.DiffError(t, nil, nil, err)

	// Basic JWT header and payload used for testing.
	// Test JWTs are created with  https://jwt.io/ using "./private.key" and "./public.key"
	/*
		{
			"alg": "RS256",
			"typ": "JWT",
			"kid": "test-key"
		}
		{
			"exp": 9876543210,
			"iss": "http://test.provider.com/",
			"aud": ["test-audience"]
		}
	*/

	// jwtToken is the valid token with the header and payload above.
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.DWRz9fYHavoHkpPMpyuouulEWJg4bcM4n6Ic-QpjJTf-sOxoJlKpvhwlCJJaD4O2Y63i-JO6tr4a002OicXW8FrsiqfXUHxxLBrvO34caGOJ3dFtF4-gT0JIrUdAlFyCy1N_ZUy5nlbrSxcYC5qfo62zzx_qgaUg7kXib9Nzhjrfs57kLoGzEsUSrB3uJ7XdKkk8_gUYowzEtxjDxQUhelxGbT_9TkZae0NIKHmrPOMJ-Co1njX3b9Z1aWNg7PzjMgQfTembJX1BV_qp-DYYu4BLTVCgIOjExkwXtg7exWyBupXWaW8G762L6q9ux2iKVbyIElVQdVtWrWsuer01sQ"

	// jwtWrongAlgToken is the same as jwtToken except for alg. signed by HS256 with "test-secret".
	jwtWrongAlgToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.b4ZFXTKPmcLLT2zbnolCjHGvjC4g1plasL5R0N23K1c"
	// jwtWrongKidToken is the same as jwtToken except for kid. kid="wong-key".
	jwtWrongKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Indyb25nLWtleSJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.blkzTXsUyEOZId3Mq2T-M8vBx8X8-MfxT_R5cAmS8PjzXVfQgQLuzT4eXwQoT0tMXiOp6bSEYPwk-7uVBZzQ9R8zH2wmgR5L69EQp2xnSII6yEeCNCXNrBF86tHbc1tv3_W7blQPSKFs2lEU8gBPqwHsC3aftClcOCtUtME3FIEhSfpu3XvxW2iG3YBOIhvS9vEeO-dmaCpI52UYPrpFIOYCa-RkykQ_9QGCGU-GbOKfSZLX2dWgbvCBivkqH11JI2UuD_9r91iIsQeGb0EH8D8HN73SZ-p3xS8qWC9zDiCk3g_8ml1Qe5Es6L_vH1zp-QfqKDauaZYpHpxX4s4_iA"
	// jwtWrongExpiredToken is the same as jwtToken except for exp. exp=987654321.
	jwtWrongExpiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMSwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.Idotjnq-hw8yOLEfegpHeGPBcrDqeeuESwvBX7OMq8DR44Wcm4h5mg5DWf582Q2mQTVgHMHPuwct6HSNdCxc5Gx4HHqE4CVOq54dv99duYVluQNReDD31G9Uc6Lgnb9zIW85bsg8Q-TnpWO9aO_RqePRVrUJ7wdRdozGipNYSwtBQfLQzORjdLnz9u6JnbLa9eQiEyihL6rAaJz3DqX5Epft2UaqO2b0s31_y6ar0IHV3nL5lR_FkaSrjOvFxA8SRFJZk7BvnA0v2xR4Jqu0MFN7ZGdo12PfPisG2n3oIyaFEtDbETswCYl4zL6lKpWN7RQZbBkg0-4DCSiXJqvckQ"
	// jwtWrongIssToken is the same as jwtToken except for iss. iss="http://wrong.provider.com/".
	jwtWrongIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly93cm9uZy5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.NQXXYsJ7XYvEOLV8dqVq33jPvyBw5OALBELObO9F0Z8YsJB0hUF5lwqoezn745ue_MqrTW4xlFBc2ctGIa49FZdgl9ptquTolYjDP4rv8m1B0-2zu4OBx_dxG_cICFGBmYGJtK9_wInyKKZTr7ilKQW74yM01dELBzPlVTcZsfOLDwnQfEPFae3G958feporeAt5ME7-wnvzfM2AdhV9HozB7Cmvy7z3eUcO5LGgPp5RrKKJupq36HvFZuFKnc67RmHnwaF23FZf_Fv7EJSzeNeyrpd9lBo3hJsZTokPKPAv_ebUkRNpNADfr89Yxq_oww_DaccIUAN75wPBrR_vYA"
	// jwtWrongAudToken is the same as jwtToken except for aud. aud=["wrong-audience"].
	jwtWrongAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsid3JvbmctYXVkaWVuY2UiXX0.Tny_Xsj9r8TED_Xooa6EYwjCm4m_KX-6MKKU9Ix96wMc3svCUCpKNmMGsZ5qzeW-37SCH-Hx1_ZMm1MnaDp0DiBIRKgPRxfNAKzXC58DhZTEkNZEDRRzN5SpzK_nwTdjkYcE02qewb2wLHzkZ7Bhu1Q06M3JYd42krLjEScWu9T5rr8ErWotEHtYUaY0mF5ez79byUTtlJ0bMUYoVMnbwvwcqMtRA5xEiLCWF7ad39xgSvY2aAGsZuwZrKPeUz2yrahZiZQstR18gqClIR5Kk63yUbtX86iBQtpL4kE0SrtblhXG0yKwjzrBIofdW5qOoIrZQnXt1FwmHgN2jhzNbQ"

	// jwtNoKidToken is the same as jwtToken except for kid. kid not exists.
	jwtNoKidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8iLCJhdWQiOlsidGVzdC1hdWRpZW5jZSJdfQ.J86YrOgvnvB5zv-ngOI4AP_F7vBTMEp3231oixD14BmNWU9S54SIueinfSxFkARYp-4kAHeIv-368FI4__gHphOYBnhSUPHvr5TR-QGblqtjOnfjmgJ3NvUkfNdaAKaTnpRy9telAAzkPUEzAxJKig8gYVyv8fK_guaoucX9v366hvld6jLx-8QSE8czqpriNhV3i2vRJPkF2hiRpUE5Ci6UtGRh3nzicao_4uktm7Kg1KcawQjtLDIKnThyMtSkKYbHv99mnFPF11cR857L0H5RSx9gHCgPtyyoUXGTiW0lQXTPsHK6WhGEzOr32969hGFR2eJyBoZxOZo9w1666A"
	// jwtNoExpToken is the same as jwtToken except for exp. exp not exists.
	jwtNoExpToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.Ao7bMa_sLsOn5LscGTRiRBprA-IuQj3RYRiQS2IOYNOLeY3tdx0YGs1vgW92hE7GyuljpHepY7Ppr23PNFOjw9bX0DnjiBycAINEKE7lna-p4YGtIi4GDpSATHR5kUHo6ck38eeu_OXJWzrn2ahaJ3lW-ZxHy0Th09tuUHLNuLjOEJGVvQTOYjnxYOawp5g5mC1b6SvtJoRd9bZsVLpWFKjLVI1xcDfuaBQLTFCD7sbHiv-WkNGtwjlOHV5h7HqZ0zHK8vuBznZtog5dkP8jOGXYoprtIbUjpOsXsKB85LFkwVdhGGI8eOahz_g_FyxYXhcyIt7TemWlRDqw7cuRww"
	// jwtNoIssToken is the same as jwtToken except for iss. iss not exists.
	jwtNoIssToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.GqD9K2v74Dp48epBthaRTnSG7cRbmS2qbHR4v69H3zV_kptoWn4IYqPLYhhN2t1FiPiASOax1ZlFGkJaD5WcM4S15UzJ0lz1KpV1o2zL00mKW38qz6XA75mGG64VtBSjFnyk2NjSAC-cV6Cy3BCwnMBW798FHwmFI1QpgMVbs-LM7hwROwENIH5YzOvO2oN-7kDamtOLj8g5hXLN_LKwxDYlWqd2cz5FWDSlTPaaV1OtlKIIsYX4UySx2m-5UgbdROYmxNRdWxVVpGPPad0qbEybbNHPea2NA7BHB9b1I7iNkDLUwysIwnTl0V3wzDn49rc4rRa64BMQ9VOWotvzNA"
	// jwtNoAudToken is the same as jwtToken except for aud. aud not exists.
	jwtNoAudToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJleHAiOjk4NzY1NDMyMTAsImlzcyI6Imh0dHA6Ly90ZXN0LnByb3ZpZGVyLmNvbS8ifQ.IKWx1IvYy27SE7iy1FLnqN-y5CX860XQLRkLhhlvPFqOyK7bsbBbWZkKBuhcfFkJmqHErsSi90R0oJ0Obdujp6EUbSgWepVz9T1oiGLMTzkewAzZy7tmIBhk9K9Bbs0VhPEB0gyZ5CC167flGdEq-fuSr7o-3l1D5fpXUjtRZXidb3vuLwmMiW1jbpTymub4Y0J1B1z58dt_HMis98AvW7ToOSND39JF3B716HJm-mGOkGP7722X-nXVcBGEygePN5xM8jb5Joc72PVOPip_uZLvj9hB64pUMDXYYAMnihq8R7BXW85sBoU-VLpuBOBE01_h5vLp8vohEkcQA_TeAw"

	// jwtValidNbf is the same as jwtToken except for nbf addition. nbf=987654321 (valid).
	jwtValidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.lNwivTGxUHl1Gzie-YnzMGiByxu9NpITQEMoUohPI1t19dJe1tFm2XeT3xC5K3Lm7UcI_5uQ7ysLQxCyqSabT027HIZjrrvdPbsM-kbxrRjld_8IPKJUwIilKqN4562rtn-dMbquHlzNZQX4MxfZNXD1Dr2ViXhBd6TklhDOrI8S9AF_mEaNZTV_hOjcxw5WAzAXA9EW-QXGNaIOFvCdk5hqZyo1QWMU2SOSNDALnRhtamdMYmt6wk0UlCpwIeJeL1zjp0OgJUOnPTOjdPTHFs6yjDJHO2zrwtx02YB4lT928SdMRDaqMzBdvI2nHxfYwR-fhaKysmOnS2x98S_nVg"
	// jwtInvalidNbf is the same as jwtToken except for nbf addition. nbf=9876543210 (invalid).
	jwtInvalidNbf := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJuYmYiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.HimYFFhJbg1QVw0Y3uHM2qFmFqSFGOLixB-R4NzGMWIsXHzFT5XOm2PLRlpWaGLDvAPlVMs50A-7qD6BPwVql1N6uHtRmrk4Cd1FBQjIi5sCiyIO_RZ4dJaP9Gcn3hgimHO-e5MQ9E2vOLIdAZizyAV8-_4gNpCEq_1tt0vMPKd1nXdClRVL4twCY2nYAuC4DZOVfTU30LRlM2SBB6bHXGGHHkN8hAcV8G9g4jaEelmW6-sWN785hRujw9a5vVnCVbgC20X0JoH6bBhUN51zTQ1TJvxDmmXgVEWh6ig6elZHbTypM441CHM6kkAnQ_8f97uwQgKAJn6R32Ghdm4Qzw"
	// jwtValidIat is the same as jwtToken except for iat addition. iat=987654321 (valid).
	jwtValidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMSwiZXhwIjo5ODc2NTQzMjEwLCJpc3MiOiJodHRwOi8vdGVzdC5wcm92aWRlci5jb20vIiwiYXVkIjpbInRlc3QtYXVkaWVuY2UiXX0.eubgpgd9ll7dtgSSE3RVnYZ_4VkSliGhGVXTaNVToCBmF-pABR_JXXrA49aQA_w3gaRmPA5IvCFx6wRBam9blX63JfM3-CZ-NVpnMO44nhNfiA4vE9lJo_uPg-ud3cWteC5QHTVjQOMBj2KUNqjmJTZoSvrcOhbd_FHtNdocu_jK5VH5_7QObgxE1i4fkbnhBaNRUQ0AFVyJs8ui0zUYhty1Vch_87ipCBsJ0VgT03Rb7xOhX2oRb5Zcn4BJH2uQ5OoFrRy63u_43js8NqF-M59vm0wrNXwWgHb3KiVCqzjDvSI8qPL4rZWp3wKxEzeg75Wqn4zsk9L0PWehMORk9A"
	// jwtInvalidIat is the same as jwtToken except for iat addition. iat=9876543210 (invalid).
	jwtInvalidIat := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJpYXQiOjk4NzY1NDMyMTAsImV4cCI6OTg3NjU0MzIxMCwiaXNzIjoiaHR0cDovL3Rlc3QucHJvdmlkZXIuY29tLyIsImF1ZCI6WyJ0ZXN0LWF1ZGllbmNlIl19.mx9coQPAFsgL12zVjCFcIIZ4oVPFlG7vWG63pL1doC4KiWVY-h_Ingxlf19us-P5rbQqH2_ylKtKhxpQNNIE2MKb-Fn3b4PZUSMNb7SKEFogENy5Lcf4puLd7_fqHkwv6cBxH_WcvD3wzWZIEkR1TlCDV1xdezrPFJgBZUyKEGf4CJ2uF_8eyn5knCVqok3TRZZk_sXPo2T1VQBibT7Sy8jSYmtKB1NR1AxAtK96G_oRnFQjjNEpRez6bZONykLzKLp73jnyDIzIcbcMrrcoG7QfGTbgGujGZqEq85_hgZjAwRcaA9k16nkzhc4hFEMMu0V4lUww9eOqoUMQAGaqFA"

	checkAuthSuccess(t, h, jwtToken)

	checkAuthFailure(t, h, jwtWrongAlgToken)
	checkAuthFailure(t, h, jwtWrongKidToken)
	checkAuthFailure(t, h, jwtWrongExpiredToken)
	checkAuthFailure(t, h, jwtWrongIssToken)
	checkAuthFailure(t, h, jwtWrongAudToken)

	checkAuthFailure(t, h, jwtNoKidToken)
	checkAuthSuccess(t, h, jwtNoExpToken) // No exp but should success.
	checkAuthFailure(t, h, jwtNoIssToken)
	checkAuthFailure(t, h, jwtNoAudToken)

	checkAuthSuccess(t, h, jwtValidNbf)
	checkAuthFailure(t, h, jwtInvalidNbf)
	checkAuthSuccess(t, h, jwtValidIat)
	checkAuthFailure(t, h, jwtInvalidIat)

}
