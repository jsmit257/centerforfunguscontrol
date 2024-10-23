package huautla

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

type huautlaMock struct {
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
	types.Vendorer
}

func checkResult(t *testing.T, b *bytes.Buffer, rx any, expected any) {
	body, err := io.ReadAll(b)
	require.Nil(t, err)
	require.Nil(t, json.Unmarshal(body, rx))
	require.Equal(t, expected, rx)
}
