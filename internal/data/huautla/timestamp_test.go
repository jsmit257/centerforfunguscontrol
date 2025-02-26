package huautla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

type (
	mockTS struct {
		updErr, undelErr error
	}
)

func Test_PatchTS(t *testing.T) {
	t.Parallel()

	ref := time.Now().UTC()

	tcs := map[string]struct {
		table string
		id    types.UUID
		org   *time.Time
		flds  []string
		err   error
		sc    int
	}{
		"happy_path": {
			table: "testtable",
			id:    "1",
			sc:    http.StatusNoContent,
			flds:  []string{"mtime"},
			org:   &ref,
		},
		"update_fails": {
			table: "testtable",
			id:    "1",
			err:   fmt.Errorf("some error"),
			flds:  []string{"mtime"},
			org:   &ref,
			sc:    http.StatusInternalServerError,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"missing_table": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"mangled_id": {
			id:    "%zzz",
			table: "testtable",
			sc:    http.StatusBadRequest,
		},
		"mangled_table": {
			id:    "1",
			table: "%zzz",
			sc:    http.StatusBadRequest,
		},
		"read_fails": {
			id:    "1",
			table: "testtable",
			sc:    http.StatusBadRequest,
		},
		"unmarshal_fails": {
			id:    "1",
			table: "testtable",
			sc:    http.StatusBadRequest,
		},
		"validate_fails": {
			id:    "1",
			table: "testtable",
			sc:    http.StatusBadRequest,
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ha := &HuautlaAdaptor{db: &huautlaMock{
				Timestamper: &mockTS{updErr: tc.err},
			}}

			body := serializeTimestamp(types.Timestamp{
				Fields: tc.flds,
				Origin: tc.org,
			})
			if name == "unmarshal_fails" {
				body = body[1:]
			}
			bodyreader := io.Reader(bytes.NewReader([]byte(body)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"table", "id"}, Values: []string{tc.table, string(tc.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPatch,
				"url",
				bodyreader)

			ha.PatchTS(w, r)
			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_Undel(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		id    types.UUID
		table string
		err   error
		sc    int
	}{
		"happy_path": {
			id:    "1",
			table: "testtable",
			sc:    http.StatusNoContent,
		},
		"query_fails": {
			id:    "1",
			table: "testtable",
			err:   fmt.Errorf("some error"),
			sc:    http.StatusInternalServerError,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"mangled_id": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"missing_table": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"mangled_table": {
			id:    "1",
			table: "%zzz",
			sc:    http.StatusBadRequest,
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {

			ha := &HuautlaAdaptor{db: &huautlaMock{
				Timestamper: &mockTS{undelErr: tc.err},
			}}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"table", "id"}, Values: []string{tc.table, string(tc.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPatch,
				"url",
				nil)

			ha.Undel(w, r)
			require.Equal(t, tc.sc, w.Code)
		})
	}
}
func serializeTimestamp(ts types.Timestamp) []byte {
	result, _ := json.Marshal(ts)
	return result
}

func (ts *mockTS) UpdateTimestamps(context.Context, string, types.UUID, types.Timestamp) error {
	return ts.updErr
}

func (ts *mockTS) Undelete(context.Context, string, types.UUID) error {
	return ts.undelErr
}
