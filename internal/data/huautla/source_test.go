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

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

type sourcerMock struct {
	add    types.Source
	addErr error

	changeErr,
	removeErr error
}

func Test_PostSource(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		gid    types.UUID
		origin string
		s      *types.Source
		result types.Source
		addErr error
		sc     int
	}{
		"happy_event_path": {
			gid:    "happy_event_path",
			origin: "event",
			s:      &types.Source{UUID: "happy_event_path"},
			result: types.Source{UUID: "happy_event_path"},
			sc:     http.StatusCreated,
		},
		"happy_strain_path": {
			gid:    "happy_strain_path",
			origin: "strain",
			s:      &types.Source{UUID: "happy_strain_path"},
			result: types.Source{UUID: "happy_strain_path"},
			sc:     http.StatusCreated,
		},
		"gid_urldecode_error": {
			gid: "%zzz",
			sc:  http.StatusBadRequest,
		},
		"origin_missing": {
			gid: "happy_path",
			sc:  http.StatusBadRequest,
		},
		"origin_urldecode_error": {
			gid:    "happy_path",
			origin: "%zzz",
			sc:     http.StatusBadRequest,
		},
		"origin_not_allowed": {
			gid:    "happy_path",
			origin: "not allowed",
			sc:     http.StatusBadRequest,
		},
		"missing_body": {
			gid:    "missing_body",
			origin: "event",
			sc:     http.StatusBadRequest,
		},
		"read_fails": {
			gid:    "read_fails",
			origin: "event",
			sc:     http.StatusBadRequest,
		},
		"unmarshal_fails": {
			gid:    "unmarshal_fails",
			origin: "event",
			s:      &types.Source{},
			sc:     http.StatusBadRequest,
		},
		"add_fails": {
			gid:    "add_fails",
			origin: "event",
			s:      &types.Source{UUID: "happy_path"},
			addErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{db: &huautlaMock{
			Sourcer: &sourcerMock{
				add:    tc.result,
				addErr: tc.addErr,
			},
		}}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id", "origin"}, Values: []string{string(tc.gid), tc.origin}}

			body := serializeSource(tc.s)
			if name == "unmarshal_fails" {
				body = body[1:]
			}

			bodyreader := io.Reader(bytes.NewReader([]byte(body)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bodyreader)

			ha.PostSource(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusCreated {
				checkResult(t, w.Body, &types.Source{}, &tc.result)
			}
		})
	}
}

func Test_PatchSource(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		gid       types.UUID
		origin    string
		sid       types.UUID
		s         *types.Source
		changeErr error
		sc        int
	}{
		"happy_event_path": {
			gid:    "happy_event_path",
			origin: "event",
			sid:    "happy_event_path",
			s:      &types.Source{UUID: "happy_event_path"},
			sc:     http.StatusNoContent,
		},
		"happy_strain_path": {
			gid:    "happy_strain_path",
			origin: "strain",
			sid:    "happy_strain_path",
			s:      &types.Source{UUID: "happy_strain_path"},
			sc:     http.StatusNoContent,
		},
		"gid_urldecode_error": {
			gid: "%zzz",
			sc:  http.StatusBadRequest,
		},
		"origin_missing": {
			gid: "happy_path",
			sc:  http.StatusBadRequest,
		},
		"origin_urldecode_error": {
			gid:    "happy_path",
			origin: "%zzz",
			sc:     http.StatusBadRequest,
		},
		"origin_not_allowed": {
			gid:    "happy_path",
			origin: "not allowed",
			sc:     http.StatusBadRequest,
		},
		"sid_urldecode_error": {
			gid:    "happy_strain_path",
			origin: "strain",
			sid:    "%zzz",
			sc:     http.StatusBadRequest,
		},
		"missing_body": {
			gid:    "missing_body",
			origin: "event",
			sid:    "mossing_body",
			sc:     http.StatusBadRequest,
		},
		"read_fails": {
			gid:    "read_fails",
			origin: "event",
			sid:    "read_fails",
			sc:     http.StatusBadRequest,
		},
		"unmarshal_fails": {
			gid:    "unmarshal_body",
			origin: "event",
			sid:    "unmarshal_body",
			s:      &types.Source{},
			sc:     http.StatusBadRequest,
		},
		"update_fails": {
			gid:       "happy_event_path",
			origin:    "event",
			sid:       "happy_event_path",
			s:         &types.Source{UUID: "happy_event_path"},
			changeErr: fmt.Errorf("some error"),
			sc:        http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{db: &huautlaMock{
			Generationer: &generationerMock{
				sel: types.Generation{},
			},
			Sourcer: &sourcerMock{changeErr: tc.changeErr},
		}}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{
				Keys:   []string{"g_id", "origin", "s_id"},
				Values: []string{string(tc.gid), tc.origin, string(tc.sid)},
			}

			body := serializeSource(tc.s)
			if name == "unmarshal_fails" {
				body = body[1:]
			}

			bodyreader := io.Reader(bytes.NewReader([]byte(body)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bodyreader)

			ha.PatchSource(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_DeleteSource(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		gid    types.UUID
		sid    types.UUID
		genErr error
		srcErr error
		sc     int
	}{
		"happy_path": {
			gid: "happy_path",
			sid: "happy_path",
			sc:  http.StatusOK,
		},
		"source_error": {
			gid:    "source_error",
			sid:    "source_error",
			srcErr: fmt.Errorf("source_error"),
			sc:     http.StatusInternalServerError,
		},
		"generation_error": {
			gid:    "generation_error",
			sid:    "lifecycle_error",
			genErr: fmt.Errorf("lifecycle_error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_generation": {
			sc: http.StatusBadRequest,
		},
		"gid_decode_error": {
			gid: "%zzz",
			sc:  http.StatusBadRequest,
		},
		"missing_source": {
			gid: "happy_path",
			sc:  http.StatusBadRequest,
		},
		"sid_decode_error": {
			gid: "happy_path",
			sid: "%zzz",
			sc:  http.StatusBadRequest,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					sel:    types.Generation{UUID: tc.gid},
					selErr: tc.genErr,
				},
				Sourcer: &sourcerMock{removeErr: tc.srcErr},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"g_id", "s_id"}, Values: []string{
				string(tc.gid),
				string(tc.sid),
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				nil)

			ha.DeleteSource(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func serializeSource(e *types.Source) []byte {
	if e == nil {
		return []byte{}
	}
	result, _ := json.Marshal(e)
	return result
}

func (sm *sourcerMock) InsertSource(context.Context, types.UUID, string, types.Source, types.CID) (types.Source, error) {
	return sm.add, sm.addErr
}
func (sm *sourcerMock) UpdateSource(context.Context, string, types.Source, types.CID) error {
	return sm.changeErr
}
func (sm *sourcerMock) RemoveSource(context.Context, *types.Generation, types.UUID, types.CID) error {
	return sm.removeErr
}
