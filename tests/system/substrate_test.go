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

var substrates []types.Substrate

func Test_HappySubstrate(t *testing.T) {
	url := fmt.Sprintf(`http://%s:%d/substrate`, cfg.HTTPHost, cfg.HTTPPort)

	for _, s := range []types.Substrate{
		{Name: "5-grain", Type: types.GrainType, Vendor: vendors[0]},
		{Name: "Rye", Type: types.GrainType, Vendor: vendors[1]},
		{Name: "Millet", Type: types.GrainType, Vendor: vendors[2]},
		{Name: "Popcorn", Type: types.GrainType, Vendor: vendors[1]},
		{Name: "Hemp", Type: types.GrainType, Vendor: vendors[3]},
		{Name: "Birdseed", Type: types.GrainType, Vendor: vendors[0]},
		{Name: "Agar", Type: types.PlatingType, Vendor: vendors[2]},
		{Name: "Liquid culture", Type: types.LiquidType, Vendor: vendors[1]},
		{Name: "Liquid culture", Type: types.LiquidType, Vendor: vendors[3]},
		{Name: "Horse cookies", Type: types.BulkType, Vendor: vendors[0]},
	} {
		b, err := json.Marshal(s)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		require.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		b, err = io.ReadAll(res.Body)
		require.Nil(t, err)

		err = json.Unmarshal(b, &s)
		require.Nil(t, err)

		substrates = append(substrates, s)
	}
}
