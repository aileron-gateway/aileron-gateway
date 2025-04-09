//go:build integration

package headercert_test

import (
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/cmd/aileron/app"
	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/aileron-gateway/aileron-gateway/test/integration/common"
)

const (
	certPath = "./client.crt"
	fpPath   = "./fingerprint.txt"
)

// func testHeaderCertMiddleware(t *testing.T, m core.Middleware) {
// 	t.Helper()

// 	cert, _ := os.ReadFile(certPath)
// 	fp, _ := os.ReadFile(fpPath)

// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("ok"))
// 	})
// 	h := m.Middleware(handler)

// 	r1 := httptest.NewRequest(http.MethodGet, "http://headercert-test.com/test", nil)
// 	r1.Header.Set("X-SSL-Client-Cert", base64.URLEncoding.EncodeToString(cert))
// 	r1.Header.Set("X-SSL-Client-Fingerprint", string(fp))
// 	w1 := httptest.NewRecorder()
// 	h.ServeHTTP(w1, r1)
// 	b1, _ := io.ReadAll(w1.Result().Body)
// 	testutil.Diff(t, http.StatusOK, w1.Result().StatusCode)
// 	testutil.Diff(t, "ok", string(b1))
// 	return
// }

func TestMinimalWithoutMetadata(t *testing.T) {
	configs := []string{"./config-minimal-without-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "",
		Namespace:  "",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile(core.ErrPrefix+`failed to create HeaderCertMiddleware`), err)

}
func TestMinimalWithMetadata(t *testing.T) {
	configs := []string{"./config-minimal-with-metadata.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "testName",
		Namespace:  "testNamespace",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile(core.ErrPrefix+`failed to create HeaderCertMiddleware`), err)
}
func TestEmptyName(t *testing.T) {
	configs := []string{"./config-empty-name.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "",
		Namespace:  "testNamespace",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile(core.ErrPrefix+`failed to create HeaderCertMiddleware`), err)
}
func TestEmptyNamespace(t *testing.T) {
	configs := []string{"./config-empty-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "testName",
		Namespace:  "",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile(core.ErrPrefix+`failed to create HeaderCertMiddleware`), err)
}
func TestEmptyNameNamespace(t *testing.T) {
	configs := []string{"./config-empty-name-namespace.yaml"}

	server := common.NewAPI()
	err := app.LoadConfigFiles(server, configs)
	testutil.DiffError(t, nil, nil, err)

	ref := &kernel.Reference{
		APIVersion: "app/v1",
		Kind:       "HeaderCertMiddleware",
		Name:       "",
		Namespace:  "",
	}
	_, err = api.ReferTypedObject[core.Middleware](server, ref)
	testutil.DiffError(t, core.ErrCoreGenCreateObject, regexp.MustCompile(core.ErrPrefix+`failed to create HeaderCertMiddleware`), err)
}
