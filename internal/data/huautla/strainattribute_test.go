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

type strainattributerMock struct {
	knownNames []string
	namesErr   error

	addResult types.StrainAttribute

	addErr,
	changeErr,
	rmErr error
}

func Test_GetStrainAttributeNames(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		knownNames []string
		namesErr   error
		sc         int
	}{
		"happy_path": {
			sc: http.StatusOK,
		},
		"sad_path": {
			namesErr: fmt.Errorf("some error"),
			sc:       http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				StrainAttributer: &strainattributerMock{
					knownNames: v.knownNames,
					namesErr:   v.namesErr,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				"url",
				bytes.NewReader([]byte("")))

			ha.GetStrainAttributeNames(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PostStrainAttribute(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s         types.Strain
		a         types.StrainAttribute
		strainErr error
		attrErr   error
		sc        int
	}{
		"happy_path": {
			s:  types.Strain{UUID: "happy_path"},
			a:  types.StrainAttribute{Name: "happy_path", Value: "squirrel"},
			sc: http.StatusCreated,
		},
		"missing_strain_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			s:  types.Strain{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"missing_attribute_name": {
			s:  types.Strain{UUID: "missing_attribute_name"},
			sc: http.StatusBadRequest,
		},
		"missing_attribute_value": {
			s:  types.Strain{UUID: "happy_path"},
			a:  types.StrainAttribute{Name: "missing_attribute_value"},
			sc: http.StatusBadRequest,
		},
		"strain_error": {
			s:         types.Strain{UUID: "strain_error"},
			a:         types.StrainAttribute{Name: "strain_error", Value: "squirrel"},
			strainErr: fmt.Errorf("strain_error"),
			sc:        http.StatusInternalServerError,
		},
		"attribute_error": {
			s:       types.Strain{UUID: "attribute_error"},
			a:       types.StrainAttribute{Name: "attribute_error", Value: "squirrel"},
			attrErr: fmt.Errorf("attribute_error"),
			sc:      http.StatusInternalServerError,
		},
		"read_fails": {
			s:  types.Strain{UUID: "read_fails"},
			a:  types.StrainAttribute{Name: "read_fails", Value: "squirrel"},
			sc: http.StatusBadRequest,
		},
		"unmarshal_fails": {
			s:  types.Strain{UUID: "unmarshal_fails"},
			a:  types.StrainAttribute{Name: "unmarshal_fails", Value: "squirrel"},
			sc: http.StatusBadRequest,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectResult: tc.s,
					selectErr:    tc.strainErr,
				},
				StrainAttributer: &strainattributerMock{addErr: tc.attrErr},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(tc.s.UUID)}}

			body := serializeAttribute(&tc.a)
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

			ha.PostStrainAttribute(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}

}

func Test_PatchStrainAttribute(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s         types.Strain
		a         types.StrainAttribute
		strainErr error
		attrErr   error
		sc        int
	}{
		"happy_path": {
			s:  types.Strain{UUID: "happy_path"},
			a:  types.StrainAttribute{UUID: "happy_path", Name: "happy_path", Value: "squirrel"},
			sc: http.StatusOK,
		},
		"missing_strain_id": {
			sc: http.StatusBadRequest,
		},
		"missing_attr_id": {
			s:  types.Strain{UUID: "missing_attr_id"},
			sc: http.StatusBadRequest,
		},
		"strain_urldecode_error": {
			s:  types.Strain{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"attr_urldecode_error": {
			s:  types.Strain{UUID: "attr_urldecode_error"},
			a:  types.StrainAttribute{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"missing_attribute_name": {
			s:  types.Strain{UUID: "missing_attribute_name"},
			a:  types.StrainAttribute{UUID: "missing_attr_name", Value: "squirrel"},
			sc: http.StatusBadRequest,
		},
		"missing_attribute_value": {
			s:  types.Strain{UUID: "missing_attr_value"},
			a:  types.StrainAttribute{UUID: "missing_attr_value", Name: "missing_attribute_value"},
			sc: http.StatusBadRequest,
		},
		"strain_error": {
			s:         types.Strain{UUID: "strain_error"},
			a:         types.StrainAttribute{UUID: "strain_error", Name: "strain_error", Value: "squirrel"},
			strainErr: fmt.Errorf("strain_error"),
			sc:        http.StatusInternalServerError,
		},
		"attribute_error": {
			s:       types.Strain{UUID: "attribute_error"},
			a:       types.StrainAttribute{UUID: "attribute_error", Name: "attribute_error", Value: "squirrel"},
			attrErr: fmt.Errorf("attribute_error"),
			sc:      http.StatusInternalServerError,
		},
		"read_fails": {
			s:  types.Strain{UUID: "read_fails"},
			a:  types.StrainAttribute{UUID: "read_fails", Value: "squirrel"},
			sc: http.StatusBadRequest,
		},
		"unmarshal_fails": {
			s:  types.Strain{UUID: "unmarshal_fails"},
			a:  types.StrainAttribute{UUID: "unmarshal_fails", Value: "squirrel"},
			sc: http.StatusBadRequest,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectResult: tc.s,
					selectErr:    tc.strainErr,
				},
				StrainAttributer: &strainattributerMock{changeErr: tc.attrErr},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"st_id", "at_id"}, Values: []string{
				string(tc.s.UUID),
				string(tc.a.UUID),
			}}

			body := serializeAttribute(&tc.a)
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

			ha.PatchStrainAttribute(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_DeleteStrainAttribute(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s         types.Strain
		id        string
		strainErr error
		rmErr     error
		sc        int
	}{
		"happy_path": {
			s:  types.Strain{UUID: "happy_path"},
			id: "max vertical",
			sc: http.StatusOK,
		},
		"missing_strainid": {
			sc: http.StatusBadRequest,
		},
		"urldecode_strain_id": {
			s:  types.Strain{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"urldecode_attribute_id": {
			s:  types.Strain{UUID: "happy_path"},
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"missing_attributeid": {
			s:  types.Strain{UUID: "missing_attributeid"},
			sc: http.StatusBadRequest,
		},
		"strain_fails": {
			s:         types.Strain{UUID: "strain_fails"},
			id:        "max vertical",
			strainErr: fmt.Errorf("some error"),
			sc:        http.StatusInternalServerError,
		},
		"remove_fails": {
			s:     types.Strain{UUID: "remove_fails"},
			id:    "max vertical",
			rmErr: fmt.Errorf("some error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectResult: tc.s,
					selectErr:    tc.strainErr,
				},
				StrainAttributer: &strainattributerMock{rmErr: tc.rmErr},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"st_id", "at_id"}, Values: []string{
				string(tc.s.UUID),
				tc.id,
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader([]byte("")))

			ha.DeleteStrainAttribute(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func serializeAttribute(sa *types.StrainAttribute) []byte {
	if sa == nil {
		return []byte{}
	}
	result, _ := json.Marshal(sa)
	return result
}

func (sa *strainattributerMock) KnownAttributeNames(ctx context.Context, cid types.CID) ([]string, error) {
	return sa.knownNames, sa.namesErr
}
func (sa *strainattributerMock) GetAllAttributes(ctx context.Context, s *types.Strain, cid types.CID) error {
	return nil
}
func (sa *strainattributerMock) AddAttribute(ctx context.Context, s *types.Strain, a types.StrainAttribute, cid types.CID) (types.StrainAttribute, error) {
	return sa.addResult, sa.addErr
}
func (sa *strainattributerMock) ChangeAttribute(ctx context.Context, s *types.Strain, a types.StrainAttribute, cid types.CID) error {
	return sa.changeErr
}
func (sa *strainattributerMock) RemoveAttribute(ctx context.Context, s *types.Strain, id types.UUID, cid types.CID) error {
	return sa.rmErr
}
