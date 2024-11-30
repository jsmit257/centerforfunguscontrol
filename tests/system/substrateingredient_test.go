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

	for s, v := range map[int][]types.Ingredient{
		0: {
			ingredients[2],
			ingredients[4],
			ingredients[9],
			ingredients[11],
			ingredients[13],
		},
		4: {
			ingredients[9],
			ingredients[10],
		},
		7: {
			ingredients[1],
		},
		8: {
			ingredients[7],
		},
	} {
		for _, i := range v {
			b, err := json.Marshal(i)
			require.Nil(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf(urlfmt, substrates[s].UUID),
				bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &substrates[s])
			require.Nil(t, err)
		}
	}
}
