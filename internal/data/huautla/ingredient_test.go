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

type IngredienterMock struct {
	selectAllResult []types.Ingredient
	selectAllErr    error

	selectResult types.Ingredient
	selectErr    error

	insertResult types.Ingredient
	insertErr    error

	updateErr error

	deleteErr error
}

func Test_SelectAllIngredients(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.Ingredient
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Ingredient{},
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
				Ingredienter: &IngredienterMock{
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
			ha.GetAllIngredients(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Ingredient{}, &v.result)
			}
		})
	}
}

func Test_GetIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Ingredient
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Ingredient{},
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
				Ingredienter: &IngredienterMock{
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

			ha.GetIngredient(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Ingredient{}, &v.result)
			}
		})
	}
}

func Test_PostIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		ingredient *types.Ingredient
		result     types.Ingredient
		err        error
		sc         int
	}{
		"happy_path": {
			ingredient: &types.Ingredient{},
			result:     types.Ingredient{},
			sc:         http.StatusCreated,
		},
		"read_fails": {
			ingredient: &types.Ingredient{},
			sc:         http.StatusBadRequest,
		},
		"missing_Ingredient": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			ingredient: &types.Ingredient{},
			err:        fmt.Errorf("db error"),
			sc:         http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Ingredienter: &IngredienterMock{
					insertResult: tc.result,
					insertErr:    tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeIngredient(tc.ingredient)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				"url",
				bodyreader)

			ha.PostIngredient(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Ingredient{}, &tc.result)
			}
		})
	}
}

func Test_PatchIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id         types.UUID
		ingredient *types.Ingredient
		err        error
		sc         int
	}{
		"happy_path": {
			id:         "1",
			ingredient: &types.Ingredient{},
			sc:         http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"read_fails": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"missing_ingredient": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:         "1",
			ingredient: &types.Ingredient{},
			err:        fmt.Errorf("db error"),
			sc:         http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Ingredienter: &IngredienterMock{
					updateErr: tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeIngredient(tc.ingredient)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(tc.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bodyreader)

			ha.PatchIngredient(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_DeleteIngredient(t *testing.T) {
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
				Ingredienter: &IngredienterMock{
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

			ha.DeleteIngredient(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func serializeIngredient(s *types.Ingredient) []byte {
	if s == nil {
		return []byte{}
	}
	result, _ := json.Marshal(s)
	return result
}

func (vm *IngredienterMock) SelectAllIngredients(context.Context, types.CID) ([]types.Ingredient, error) {
	return vm.selectAllResult, vm.selectAllErr
}

func (vm *IngredienterMock) SelectIngredient(context.Context, types.UUID, types.CID) (types.Ingredient, error) {
	return vm.selectResult, vm.selectErr
}

func (vm *IngredienterMock) InsertIngredient(context.Context, types.Ingredient, types.CID) (types.Ingredient, error) {
	return vm.insertResult, vm.insertErr
}

func (vm *IngredienterMock) UpdateIngredient(context.Context, types.UUID, types.Ingredient, types.CID) error {
	return vm.updateErr
}

func (vm *IngredienterMock) DeleteIngredient(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
