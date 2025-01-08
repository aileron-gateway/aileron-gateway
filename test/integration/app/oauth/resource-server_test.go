//go:build integration
// +build integration

package oauth_test

// func TestAuthorization(t *testing.T) {
// 	configs := []string{"./config-rs-success.yaml"}

// 	server := common.NewAPI()
// 	err := app.LoadConfigFiles(server, configs)
// 	testutil.DiffError(t, nil, nil, err)

// 	ref := &kernel.Reference{
// 		APIVersion: "app/v1",
// 		Kind:       "OAuthAuthenticationHandler",
// 		Name:       "",
// 		Namespace:  "",
// 	}
// 	h, err := api.ReferTypedObject[apps.AuthenticationHandler](server, ref)
// 	testutil.DiffError(t, nil, nil, err)

// 	// authorization allowed
// 	ctx := context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 20})
// 	r1 := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil).WithContext(ctx)
// 	r1.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3Qta2V5In0.eyJhdWQiOlsiVEVTVF9DTElFTlRfSUQiXX0.cckc6-kDsBd9vrc8-PhUMW1S_XEIZrH5qPqWcLf4Dg7DPsZsWPYcv_uWOFVw_wSeM_EniqVwvajSV3PPcP_TTonKYPjh0ccodkxYFL6tW2PBfqt4JcNZqn6ipzpQA8LZcyyH-9olRDk_NN5EPCHWCJePNVs7RjobQVeVRG25s5R4b5bbyuXZwcleoG4Q6MJX2LfU33udjx-PMfjNDDmmmBfYYgomMQPbKyu89xk3be5HCyUkh41sRWGd9omNC70vzNz4QF3EWsdkJl9E8ILcIAZrUQMgvGZAFV1kQzFDbWaaa1RGDT_ja_H4ufNsvgBbh1Ph-MulTkZmFAVhr6pq9w")
// 	w1 := httptest.NewRecorder()
// 	_, result, shouldReturn, err := h.ServeAuthn(w1, r1)
// 	testutil.Diff(t, apps.AuthSucceeded, result)
// 	testutil.Diff(t, false, shouldReturn)
// 	testutil.Diff(t, nil, err)
// 	b1, _ := io.ReadAll(w1.Result().Body)
// 	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
// 	testutil.Diff(t, `{"message":"allowed"}`, string(b1))

// 	// authorization denied
// 	// ctx = context.WithValue(context.Background(), "AuthClaims", map[string]any{"age": 18})
// 	// r2 := httptest.NewRequest(http.MethodGet, "http://casbin.com/denied", nil).WithContext(ctx)
// 	// w2 := httptest.NewRecorder()
// 	// h.ServeHTTP(w2, r2)
// 	// b2, _ := io.ReadAll(w2.Result().Body)
// 	// testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
// 	// testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

// }

// func TestAuthorizationSkip(t *testing.T) {
// 	configs := []string{"./config-authorization-skip.yaml"}

// 	server := common.NewAPI()
// 	err := app.LoadConfigFiles(server, configs)
// 	testutil.DiffError(t, nil, nil, err)

// 	ref := &kernel.Reference{
// 		APIVersion: "app/v1",
// 		Kind:       "OAuthAuthenticationHandler",
// 		Name:       "test",
// 		Namespace:  "testNamespace",
// 	}
// 	m, err := api.ReferTypedObject[core.Middleware](server, ref)
// 	testutil.DiffError(t, nil, nil, err)

// 	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte(`{"message":"skipped"}`))
// 	}))

// 	h := m.Middleware(handler)

// 	// authorization skipped
// 	r1 := httptest.NewRequest(http.MethodGet, "http://casbin.com/skipped", nil)
// 	w1 := httptest.NewRecorder()
// 	h.ServeHTTP(w1, r1)
// 	b1, _ := io.ReadAll(w1.Result().Body)
// 	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
// 	testutil.Diff(t, `{"message":"skipped"}`, string(b1))

// 	// authorization Forbidden
// 	r2 := httptest.NewRequest(http.MethodGet, "http://casbin.com/forbidden", nil)
// 	w2 := httptest.NewRecorder()
// 	h.ServeHTTP(w2, r2)
// 	b2, _ := io.ReadAll(w2.Result().Body)
// 	testutil.Diff(t, http.StatusForbidden, w2.Result().StatusCode)
// 	testutil.Diff(t, `{"status":403,"statusText":"Forbidden"}`, string(b2))

// }
