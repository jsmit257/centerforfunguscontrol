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

var vendors []types.Vendor

func init() {
	if req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(`http://%s:%d/vendors`, cfg.HTTPHost, cfg.HTTPPort),
		nil); err != nil {
		panic(err)
	} else if res, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else if b, err := io.ReadAll(res.Body); err != nil {
		panic(err)
	} else if err = json.Unmarshal(b, &vendors); err != nil {
		panic(err)
	}
}

func Test_HappyVendor(t *testing.T) {
	url := fmt.Sprintf(`http://%s:%d/vendor`, cfg.HTTPHost, cfg.HTTPPort)

	for _, v := range []types.Vendor{
		{Name: "Fun Guys", Website: "http://www.example.com"},
		{Name: "Nuthin but Fungus", Website: "http://www.example.com"},
		{Name: "Mycellium Emporium", Website: "http://www.example.com"},
	} {
		b, err := json.Marshal(v)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		require.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		b, err = io.ReadAll(res.Body)
		require.Nil(t, err)

		err = json.Unmarshal(b, &v)
		require.Nil(t, err)

		vendors = append(vendors, v)
	}
}
