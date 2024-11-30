package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var ingredients []types.Ingredient

func init() {
	url := fmt.Sprintf(`http://%s:%d/ingredients`, cfg.HTTPHost, cfg.HTTPPort)

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		panic(err)
	} else if req.AddCookie(cookie); false {
	} else if res, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else if http.StatusOK != res.StatusCode {
		panic(fmt.Errorf("expected sc: 200, got: %d", res.StatusCode))
	} else if b, err := io.ReadAll(res.Body); err != nil {
		panic(err)
	} else if err = json.Unmarshal(b, &ingredients); err != nil {
		panic(err)
	}
}

func Test_HappyIngredient(t *testing.T) {
	t.Skip()
	url := fmt.Sprintf(`http://%s:%d/ingredient`, cfg.HTTPHost, cfg.HTTPPort)

	for _, i := range []types.Ingredient{
		{Name: "Vermiculite"},
		{Name: "Maltodextrin"},
		{Name: "Rye"},
		{Name: "White Millet"},
		{Name: "Popcorn"},
		{Name: "Manure"},
		{Name: "Coir"},
		{Name: "Honey"},
		{Name: "Agar"},
		{Name: "Rice Flour"},
		{Name: "Hemp Seeds"},
		{Name: "White Milo"},
		{Name: "Red Milo"},
		{Name: "Red Millet"},
		{Name: "Gypsum"},
		{Name: "Calcium phosphate"},
		{Name: "Diammonium phosphate"},
	} {
		b, err := json.Marshal(i)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		require.Nil(t, err)
		req.AddCookie(cookie)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		b, err = io.ReadAll(res.Body)
		require.Nil(t, err)

		err = json.Unmarshal(b, &i)
		require.Nil(t, err)

		ingredients = append(ingredients, i)
	}
}
