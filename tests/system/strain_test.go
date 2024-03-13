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

var strains = []types.Strain{}

func Test_HappyStrain(t *testing.T) {
	url := fmt.Sprintf(`http://%s:%d/strain`, cfg.HTTPHost, cfg.HTTPPort)

	for _, s := range []types.Strain{
		{Name: "Morel", Species: "M.anatolica", Vendor: vendors[0]},
		{Name: "Hens o' the Wood", Species: "G.frondosa", Vendor: vendors[2]},
		{Name: "Reishi", Species: "G.lingzhi", Vendor: vendors[1]},
		{Name: "Turkey Tail", Species: "T.versicolor", Vendor: vendors[3]},
		{Name: "Scooby Snacks", Species: "Vanilla Wafers", Vendor: vendors[1]},
		{Name: "Shitake", Species: "L.edodes", Vendor: vendors[0]},
		{Name: "Chestnut", Species: "P.adiposa", Vendor: vendors[2]},
		{Name: "Hericium", Species: "H.abietis", Vendor: vendors[0]},
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

		strains = append(strains, s)
	}
}
