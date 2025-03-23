package huautla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

type (
	errReader string

	huautlaMock struct {
		types.EventTyper
		types.Generationer
		types.GenerationEventer
		types.Ingredienter
		types.LifecycleEventer
		types.Lifecycler
		types.Noter
		types.Observer
		types.Photoer
		types.Sourcer
		types.Stager
		types.StrainAttributer
		types.Strainer
		types.SubstrateIngredienter
		types.Substrater
		types.Timestamper
		types.Vendorer
	}
)

func (er errReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("%s", er)
}

func checkResult(t *testing.T, b *bytes.Buffer, rx any, expected any) {
	body, err := io.ReadAll(b)
	require.Nil(t, err)
	require.Nil(t, json.Unmarshal(body, rx), string(body))
	require.Equal(t, expected, rx)
}
