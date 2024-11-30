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

var sources []types.Source

func Test_HappyStrainSource(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/generation/%%s/sources/strain`, cfg.HTTPHost, cfg.HTTPPort)

	for k, v := range map[int][]types.Source{
		3: {
			{Type: "Clone", Strain: strains[0]},
		},
		4: {
			{Type: "Spore", Strain: strains[0]},
		},
		5: {
			{Type: "Spore", Strain: strains[0]},
			{Type: "Spore", Strain: strains[1]},
		},
	} {
		for _, s := range v {
			url := fmt.Sprintf(urlfmt, generations[k].UUID)

			b, err := json.Marshal(s)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &s)
			require.Nil(t, err)

			sources = append(sources, s)
		}
	}
}

func Test_HappyEventSource(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/generation/%%s/sources/event`, cfg.HTTPHost, cfg.HTTPPort)

	var s types.Source

	for k, v := range map[int][]types.Event{
		2: {
			findEvent("Clone", "Generation", lifecycles[2].Events),
		},
		1: {
			findEvent("Spore print", "Generation", lifecycles[2].Events),
		},
		0: {
			findEvent("Spore print", "Generation", lifecycles[2].Events),
			findEvent("Spore print", "Generation", lifecycles[0].Events),
		},
	} {
		for _, e := range v {
			url := fmt.Sprintf(urlfmt, generations[k].UUID)

			b, err := json.Marshal(e)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode, "%#v", e)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &s)
			require.Nil(t, err)

			sources = append(sources, s)
		}
	}
}

func findEvent(name, sev string, events []types.Event) types.Event {
	for _, e := range events {
		if e.EventType.Severity == sev && e.EventType.Name == name {
			return e
		}
	}
	return types.Event{UUID: "not-found"}
}
