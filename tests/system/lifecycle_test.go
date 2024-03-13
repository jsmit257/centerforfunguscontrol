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

var lifecycles []types.Lifecycle

func Test_HappyLifecycle(t *testing.T) {
	url := fmt.Sprintf(`http://%s:%d/lifecycle`, cfg.HTTPHost, cfg.HTTPPort)

	for _, l := range []types.Lifecycle{
		{Location: "1st chair, 2nd violin", Strain: strains[0], GrainSubstrate: substrates[0], BulkSubstrate: substrates[9]},
		{Location: "cat box", Strain: strains[3], GrainSubstrate: substrates[2], BulkSubstrate: substrates[8]},
		{Location: "6 underground", Strain: strains[5], GrainSubstrate: substrates[2], BulkSubstrate: substrates[9]},
	} {
		b, err := json.Marshal(l)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		require.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		b, err = io.ReadAll(res.Body)
		require.Nil(t, err)

		err = json.Unmarshal(b, &l)
		require.Nil(t, err)

		lifecycles = append(lifecycles, l)
	}
}
