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

func Test_HappyIngredient(t *testing.T) {
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
