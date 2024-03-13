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

func Test_HappyStrainAttribute(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/strain/%%s/attribute`, cfg.HTTPHost, cfg.HTTPPort)

	for s, v := range map[int][]types.StrainAttribute{
		4: {
			{Name: "Daphne", Value: "Hot"},
		},
		6: {
			{Name: "Headroom", Value: "18cm"},
			{Name: "Color", Value: "chestnut"},
			{Name: "Yield", Value: "high"},
		},
	} {
		url := fmt.Sprintf(urlfmt, strains[s].UUID)
		for _, a := range v {
			b, err := json.Marshal(a)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			require.Nil(t, err)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &strains[s])
			require.Nil(t, err)
		}
	}
}
