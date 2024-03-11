package huautla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type strainattributerMock struct {
	knownNames []string
	namesErr   error

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
			log:   log.WithFields(log.Fields{"test": "Test_GetStrainAttributeNames", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
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
		n, v      string
		strainErr error
		attrErr   error
		sc        int
	}{
		"happy_path": {
			s:  types.Strain{UUID: "happy_path"},
			n:  "happy_path",
			v:  "squirrel",
			sc: http.StatusCreated,
		},
		"missing_strain_id": {
			sc: http.StatusBadRequest,
		},
		"missing_attribute_name": {
			s:  types.Strain{UUID: "missing_attribute_name"},
			sc: http.StatusBadRequest,
		},
		"missing_attribute_value": {
			s:  types.Strain{UUID: "happy_path"},
			n:  "missing_attribute_value",
			sc: http.StatusBadRequest,
		},
		"strain_error": {
			s:         types.Strain{UUID: "strain_error"},
			n:         "strain_error",
			v:         "squirrel",
			strainErr: fmt.Errorf("strain_error"),
			sc:        http.StatusInternalServerError,
		},
		"attribute_error": {
			s:       types.Strain{UUID: "attribute_error"},
			n:       "attribute_error",
			v:       "squirrel",
			attrErr: fmt.Errorf("attribute_error"),
			sc:      http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectResult: v.s,
					selectErr:    v.strainErr,
				},
				StrainAttributer: &strainattributerMock{addErr: v.attrErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostStrainAttribute", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id", "at_name", "at_value"}, Values: []string{
				string(v.s.UUID),
				v.n,
				v.v,
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeAttribute(nil)))

			ha.PostStrainAttribute(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}

}

func Test_PatchStrainAttribute(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s         types.Strain
		n, v      string
		strainErr error
		attrErr   error
		sc        int
	}{
		"happy_path": {
			s:  types.Strain{UUID: "happy_path"},
			n:  "happy_path",
			v:  "squirrel",
			sc: http.StatusOK,
		},
		"missing_strain_id": {
			sc: http.StatusBadRequest,
		},
		"missing_attribute_name": {
			s:  types.Strain{UUID: "missing_attribute_name"},
			sc: http.StatusBadRequest,
		},
		"missing_attribute_value": {
			s:  types.Strain{UUID: "happy_path"},
			n:  "missing_attribute_value",
			sc: http.StatusBadRequest,
		},
		"strain_error": {
			s:         types.Strain{UUID: "strain_error"},
			n:         "strain_error",
			v:         "squirrel",
			strainErr: fmt.Errorf("strain_error"),
			sc:        http.StatusInternalServerError,
		},
		"attribute_error": {
			s:       types.Strain{UUID: "attribute_error"},
			n:       "attribute_error",
			v:       "squirrel",
			attrErr: fmt.Errorf("attribute_error"),
			sc:      http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectResult: v.s,
					selectErr:    v.strainErr,
				},
				StrainAttributer: &strainattributerMock{changeErr: v.attrErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostStrainAttribute", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"st_id", "at_name", "at_value"}, Values: []string{
				string(v.s.UUID),
				v.n,
				v.v,
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeAttribute(nil)))

			ha.PatchStrainAttribute(w, r)

			require.Equal(t, v.sc, w.Code)
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

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectResult: v.s,
					selectErr:    v.strainErr,
				},
				StrainAttributer: &strainattributerMock{rmErr: v.rmErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_DeleteStrainAttribute", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"st_id", "at_id"}, Values: []string{
				string(v.s.UUID),
				v.id,
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader([]byte("")))

			ha.DeleteStrainAttribute(w, r)

			require.Equal(t, v.sc, w.Code)
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
func (sa *strainattributerMock) AddAttribute(ctx context.Context, s *types.Strain, n, v string, cid types.CID) error {
	return sa.addErr
}
func (sa *strainattributerMock) ChangeAttribute(ctx context.Context, s *types.Strain, n, v string, cid types.CID) error {
	return sa.changeErr
}
func (sa *strainattributerMock) RemoveAttribute(ctx context.Context, s *types.Strain, id types.UUID, cid types.CID) error {
	return sa.rmErr
}
