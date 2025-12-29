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

var generations []types.Generation

func Test_HappyGeneration(t *testing.T) {
	url := fmt.Sprintf(`http://%s:%d/generation`, cfg.HTTPHost, cfg.HTTPPort)

	for _, gen := range []types.Generation{
		{PlatingSubstrate: substrates[6], LiquidSubstrate: substrates[7]}, // spores from 2 lifecycles
		{PlatingSubstrate: substrates[6], LiquidSubstrate: substrates[8]}, // one spore from a lifecycle
		{PlatingSubstrate: substrates[6], LiquidSubstrate: substrates[7]}, // one clone from a lifecycle
		{PlatingSubstrate: substrates[6], LiquidSubstrate: substrates[8]}, // one clone from a strain
		{PlatingSubstrate: substrates[6], LiquidSubstrate: substrates[7]}, // one spore from a strain
		{PlatingSubstrate: substrates[6], LiquidSubstrate: substrates[8]}, // spores from 2 strains
	} {
		b, err := json.Marshal(gen)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		require.Nil(t, err)
		req.AddCookie(cookie)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		b, err = io.ReadAll(res.Body)
		require.Nil(t, err)

		err = json.Unmarshal(b, &gen)
		require.Nil(t, err)

		generations = append(generations, gen)
	}
}

func Test_HappyGeneratedStrain(t *testing.T) {

	for i, gen := range generations[0:3] {
		url := fmt.Sprintf(
			`http://%s:%d/strain/%s/generation/%s`,
			cfg.HTTPHost,
			cfg.HTTPPort,
			strains[2-i].UUID,
			gen.UUID)

		req, err := http.NewRequest(http.MethodPatch, url, nil)
		require.Nil(t, err)
		req.AddCookie(cookie)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	}
}
