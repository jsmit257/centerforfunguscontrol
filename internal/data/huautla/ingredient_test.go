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
			log:   log.WithFields(log.Fields{"test": "Test_GetAllIngredients", "case": k}),
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
			log:   log.WithFields(log.Fields{"test": "Test_GetIngredient", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
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
		Ingredient *types.Ingredient
		result     types.Ingredient
		err        error
		sc         int
	}{
		"happy_path": {
			Ingredient: &types.Ingredient{},
			result:     types.Ingredient{},
			sc:         http.StatusCreated,
		},
		"missing_Ingredient": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			Ingredient: &types.Ingredient{},
			err:        fmt.Errorf("db error"),
			sc:         http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Ingredienter: &IngredienterMock{
					insertResult: v.result,
					insertErr:    v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostIngredient", "case": k}),
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
				bytes.NewReader(serializeIngredient(v.Ingredient)))

			ha.PostIngredient(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Ingredient{}, &v.result)
			}
		})
	}
}

func Test_PatchIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id         types.UUID
		Ingredient *types.Ingredient
		err        error
		sc         int
	}{
		"happy_path": {
			id:         "1",
			Ingredient: &types.Ingredient{},
			sc:         http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"missing_Ingredient": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:         "1",
			Ingredient: &types.Ingredient{},
			err:        fmt.Errorf("db error"),
			sc:         http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Ingredienter: &IngredienterMock{
					updateErr: v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PatchIngredient", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader(serializeIngredient(v.Ingredient)))

			ha.PatchIngredient(w, r)

			require.Equal(t, v.sc, w.Code)
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
			log:   log.WithFields(log.Fields{"test": "Test_DeleteIngredient", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
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
