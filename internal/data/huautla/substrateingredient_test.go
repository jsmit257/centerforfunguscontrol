package huautla

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type ingredienterMock struct {
	getErr,
	addErr,
	changeErr,
	rmErr error
}

func Test_PostSubstrateIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s             types.Substrate
		i             *types.Ingredient
		substrateErr  error
		ingredientErr error
		sc            int
	}{
		"happy_path": {
			s:  types.Substrate{UUID: "happy"},
			i:  &types.Ingredient{UUID: "rye", Name: "rye"},
			sc: http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"missing_ingredient": {
			s:  types.Substrate{UUID: "happy"},
			sc: http.StatusBadRequest,
		},
		"substrate_error": {
			s:            types.Substrate{UUID: "happy"},
			i:            &types.Ingredient{UUID: "rye", Name: "rye"},
			substrateErr: fmt.Errorf("some error"),
			sc:           http.StatusInternalServerError,
		},
		"add_error": {
			s:             types.Substrate{UUID: "happy"},
			i:             &types.Ingredient{UUID: "rye", Name: "rye"},
			ingredientErr: fmt.Errorf("some error"),
			sc:            http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			ha := &HuautlaAdaptor{
				db: &huautlaMock{
					Substrater: &substraterMock{
						selectResult: v.s,
						selectErr:    v.substrateErr,
					},
					SubstrateIngredienter: &ingredienterMock{
						addErr: v.ingredientErr,
					},
				},
				log: log.WithFields(log.Fields{"test": "Test_PostSubstrateIngredient", "case": k}),
			}
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.s.UUID)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader(serializeIngredient(v.i)))

			ha.PostSubstrateIngredient(w, r)
			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PatchSubstrateIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s             types.Substrate
		i             *types.Ingredient
		id            string
		substrateErr  error
		ingredientErr error
		sc            int
	}{
		"happy_path": {
			s:  types.Substrate{UUID: "happy"},
			id: "rye",
			i:  &types.Ingredient{UUID: "millet", Name: "millet"},
			sc: http.StatusOK,
		},
		"missing_substrate": {
			sc: http.StatusBadRequest,
		},
		"missing_old_ingredient": {
			s:  types.Substrate{UUID: "missing_old_ingredient"},
			sc: http.StatusBadRequest,
		},
		"no_data": {
			s:  types.Substrate{UUID: "no_data"},
			id: "rye",
			sc: http.StatusBadRequest,
		},
		"failed_substrate": {
			s:            types.Substrate{UUID: "failed_substrate"},
			id:           "rye",
			i:            &types.Ingredient{UUID: "millet", Name: "millet"},
			substrateErr: fmt.Errorf("some error"),
			sc:           http.StatusInternalServerError,
		},
		"failed_update": {
			s:             types.Substrate{UUID: "failed_update"},
			id:            "rye",
			i:             &types.Ingredient{UUID: "millet", Name: "millet"},
			ingredientErr: fmt.Errorf("some error"),
			sc:            http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			ha := &HuautlaAdaptor{
				db: &huautlaMock{
					Ingredienter: &IngredienterMock{
						selectResult: types.Ingredient{UUID: types.UUID(v.id)},
					},
					Substrater: &substraterMock{
						selectResult: v.s,
						selectErr:    v.substrateErr,
					},
					SubstrateIngredienter: &ingredienterMock{
						changeErr: v.ingredientErr,
					},
				},
				log: log.WithFields(log.Fields{"test": "Test_PatchSubstrateIngredient", "case": k}),
			}
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"su_id", "ig_id"}, Values: []string{
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
				bytes.NewReader(serializeIngredient(v.i)))

			ha.PatchSubstrateIngredient(w, r)
			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteSubstrateIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s             types.Substrate
		id            string
		substrateErr  error
		ingredientErr error
		sc            int
	}{
		"happy_path": {
			s:  types.Substrate{UUID: "happy"},
			id: "rye",
			sc: http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"substrate_error": {
			s:            types.Substrate{UUID: "happy"},
			id:           "rye",
			substrateErr: fmt.Errorf("some error"),
			sc:           http.StatusInternalServerError,
		},
		"missing_ingredient": {
			s:  types.Substrate{UUID: "happy"},
			sc: http.StatusBadRequest,
		},
		"remove_error": {
			s:             types.Substrate{UUID: "happy"},
			id:            "rye",
			ingredientErr: fmt.Errorf("some error"),
			sc:            http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			ha := &HuautlaAdaptor{
				db: &huautlaMock{
					Substrater: &substraterMock{
						selectResult: v.s,
						selectErr:    v.substrateErr,
					},
					SubstrateIngredienter: &ingredienterMock{
						rmErr: v.ingredientErr,
					},
				},
				log: log.WithFields(log.Fields{"test": "Test_AddIngredient", "case": k}),
			}
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"su_id", "ig_id"}, Values: []string{
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

			ha.DeleteSubstrateIngredient(w, r)
			require.Equal(t, v.sc, w.Code)
		})
	}

}

func (im *ingredienterMock) GetAllIngredients(ctx context.Context, s *types.Substrate, cid types.CID) error {
	return im.getErr
}

func (im *ingredienterMock) AddIngredient(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error {
	return im.addErr
}

func (im *ingredienterMock) ChangeIngredient(ctx context.Context, s *types.Substrate, oldI, newI types.Ingredient, cid types.CID) error {
	return im.changeErr
}

func (im *ingredienterMock) RemoveIngredient(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error {
	return im.rmErr
}
