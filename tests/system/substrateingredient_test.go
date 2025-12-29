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

func Test_HappySubstrateIngredient(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/substrate/%%s/ingredients`, cfg.HTTPHost, cfg.HTTPPort)

	for sub, ings := range map[int][]types.Ingredient{
		0: {
			ingredients[2],
			ingredients[4],
			ingredients[9],
			ingredients[11],
			ingredients[13],
		},
		2: {
			ingredients[2],
			ingredients[4],
			ingredients[6],
		},
		4: {
			ingredients[3],
			ingredients[9],
			ingredients[10],
		},
		5: {
			ingredients[5],
			ingredients[8],
			ingredients[12],
		},
		7: {
			ingredients[1],
			ingredients[5],
			ingredients[6],
		},
		8: {
			ingredients[7],
			ingredients[8],
			ingredients[9],
		},
	} {
		for _, ing := range ings {
			b, err := json.Marshal(ing)
			require.Nil(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf(urlfmt, substrates[sub].UUID),
				bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &substrates[sub])
			require.Nil(t, err)
		}

		require.Equal(t, len(ings), len(substrates[sub].Ingredients))
	}
}
