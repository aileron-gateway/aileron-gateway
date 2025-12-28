// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package bodylimit

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/core"
	"github.com/aileron-gateway/aileron-gateway/internal/testutil"
	"github.com/aileron-gateway/aileron-gateway/kernel/api"
	utilhttp "github.com/aileron-gateway/aileron-gateway/util/http"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCreate(t *testing.T) {
	type condition struct {
		manifest protoreflect.ProtoMessage
	}

	type action struct {
		bl         any
		err        any
		errPattern *regexp.Regexp
	}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"create with default",
			&condition{
				manifest: Resource.Default(),
			},
			&action{
				bl: &bodyLimit{
					eh:       utilhttp.GlobalErrorHandler(utilhttp.DefaultErrorHandlerName),
					maxSize:  1 << 22,
					memLimit: 1 << 22,
					tempPath: filepath.Clean(os.TempDir()) + "/",
				},
			},
		),
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			server := api.NewContainerAPI()

			bl, err := Resource.Create(server, tt.C.manifest)
			testutil.DiffError(t, tt.A.err, tt.A.errPattern, err)

			opts := []cmp.Option{
				cmp.AllowUnexported(bodyLimit{}),
				cmp.Comparer(testutil.ComparePointer[core.ErrorHandler]),
			}
			testutil.Diff(t, tt.A.bl, bl, opts...)
		})
	}
}
