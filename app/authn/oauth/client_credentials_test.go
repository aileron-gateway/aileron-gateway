// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package oauth

// func TestYyyyy(t *testing.T) {

// 	e, _ := url.Parse("http://localhost:18080/realms/aileron/.well-known/openid-configuration")
// 	p := &provider{
// 		DiscoveryEP: e,
// 	}
// 	p.discover()

// 	skipper, _ := matcher.NewSkipper(nil, nil)

// 	f := &clientCredentialsHandler{
// 		redeemTokenClient: &redeemTokenClient{
// 			provider: p,
// 			client: &client{
// 				ID:     "oauth_client_credentials",
// 				Secret: "j9wEw6Zj3dhGIhVzGCD1jpsDLDkMl8wD",
// 				Scope:  "openid",
// 			},
// 			rt:        http.DefaultTransport,
// 			grantType: "client_credentials",
// 		},
// 		skipper: skipper,
// 	}

// 	r, _ := http.NewRequestWithContext(context.Background(),http.MethodGet, "http://localhost:8888/auth", nil)
// 	r, _, _, err := f.ServeAuthn(nil, r)
// 	t.Error(err)
// }
