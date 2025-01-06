//go:build example
// +build example

package example_test

// func TestHTTP3Server(t *testing.T) {

// 	env := []string{}
// 	config := []string{"../../_example/http3-server/"}
// 	entrypoint := getEntrypointRunner(t, env, config)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	timer := time.AfterFunc(5*time.Second, cancel)

// 	pem, _ := os.ReadFile("../../_example/http3-server/pki/cert.crt")
// 	pool := x509.NewCertPool()
// 	pool.AppendCertsFromPEM(pem)
// 	h3t := &http3.Transport{
// 		TLSClientConfig: &tls.Config{
// 			RootCAs: pool,
// 		},
// 	}

// 	var status int
// 	var body []byte
// 	var err error
// 	go func() {
// 		req, _ := http.NewRequest(http.MethodGet, "https://localhost:8443/", nil)
// 		resp, e := h3t.RoundTrip(req)
// 		err = e
// 		status = resp.StatusCode
// 		body, _ = io.ReadAll(resp.Body)
// 		timer.Stop()
// 		cancel()
// 	}()

// 	if err := entrypoint.Run(ctx); err != nil {
// 		t.Error(err)
// 	}

// 	assert.Nil(t, err)
// 	assert.Equal(t, "", string(body))
// 	assert.Equal(t, http.StatusInternalServerError, status)

// }
