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

	for strain, attrs := range map[int][]types.StrainAttribute{
		0: {
			{Name: "Contamination resistance", Value: "high"},
			{Name: "Harvest when", Value: "cap wrinkles"},
		},
		3: {
			{Name: "Shape", Value: "tall/thin"},
			{Name: "Spore color", Value: "purple"},
		},
		4: {
			{Name: "Daphne", Value: "Hot"},
		},
		6: {
			{Name: "Headroom", Value: "18cm"},
			{Name: "Color", Value: "chestnut"},
			{Name: "Yield", Value: "high"},
		},
		7: {
			{Name: "Growth rate", Value: "slow"},
			{Name: "Yield", Value: "medium"},
		},
	} {
		url := fmt.Sprintf(urlfmt, strains[strain].UUID)
		for _, attr := range attrs {
			b, err := json.Marshal(attr)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &attr)
			require.Nil(t, err)

			require.NotEmpty(t, attr.UUID)
		}
	}
}
