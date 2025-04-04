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

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
)

type vendorerMock struct {
	selectAllResult []types.Vendor
	selectAllErr    error

	selectResult types.Vendor
	selectErr    error

	insertResult types.Vendor
	insertErr    error

	updateErr error

	deleteErr error

	vr    types.Entity
	vrErr error
}

func Test_GetAllVendors(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Vendor
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Vendor{},
			sc:     http.StatusOK,
		},
		"db_error": {
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Vendorer: &vendorerMock{
					selectAllResult: v.result,
					selectAllErr:    v.err,
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

			ha.GetAllVendors(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Vendor{}, &v.result)
			}
		})
	}
}

func Test_GetVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Vendor
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Vendor{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Vendorer: &vendorerMock{
					selectResult: v.result,
					selectErr:    v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				bytes.NewReader([]byte("")))

			ha.GetVendor(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Vendor{}, &v.result)
			}
		})
	}
}

func Test_PostVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		vendor *types.Vendor
		result types.Vendor
		err    error
		sc     int
	}{
		"happy_path": {
			vendor: &types.Vendor{},
			result: types.Vendor{},
			sc:     http.StatusCreated,
		},
		"missing_vendor": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			vendor: &types.Vendor{},
			err:    fmt.Errorf("db error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Vendorer: &vendorerMock{
					insertResult: v.result,
					insertErr:    v.err,
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
				bytes.NewReader(serializeVendor(v.vendor)))

			ha.PostVendor(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Vendor{}, &v.result)
			}
		})
	}
}

func Test_PatchVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		vendor *types.Vendor
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			vendor: &types.Vendor{},
			sc:     http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"missing_vendor": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:     "1",
			vendor: &types.Vendor{},
			err:    fmt.Errorf("db error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Vendorer: &vendorerMock{
					updateErr: v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader(serializeVendor(v.vendor)))

			ha.PatchVendor(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  string
		err error
		sc  int
	}{
		"happy_path": {
			id: "1",
			sc: http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Vendorer: &vendorerMock{
					deleteErr: v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader([]byte("")))

			ha.DeleteVendor(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_GetVendorReport(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Entity
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Entity{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Vendorer: &vendorerMock{
					vr:    v.result,
					vrErr: v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				nil)

			ha.GetVendorReport(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &v.result)
			}
		})
	}
}

func serializeVendor(v *types.Vendor) []byte {
	if v == nil {
		return []byte{}
	}
	result, _ := json.Marshal(v)
	return result
}

func (vm *vendorerMock) SelectAllVendors(context.Context, types.CID) ([]types.Vendor, error) {
	return vm.selectAllResult, vm.selectAllErr
}
func (vm *vendorerMock) SelectVendor(context.Context, types.UUID, types.CID) (types.Vendor, error) {
	return vm.selectResult, vm.selectErr
}
func (vm *vendorerMock) InsertVendor(context.Context, types.Vendor, types.CID) (types.Vendor, error) {
	return vm.insertResult, vm.insertErr
}
func (vm *vendorerMock) UpdateVendor(context.Context, types.UUID, types.Vendor, types.CID) error {
	return vm.updateErr
}
func (vm *vendorerMock) DeleteVendor(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
func (vm *vendorerMock) VendorReport(context.Context, types.UUID, types.CID) (types.Entity, error) {
	return vm.vr, vm.vrErr
}
