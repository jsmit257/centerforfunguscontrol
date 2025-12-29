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

	for gen, srcs := range map[int][]types.Source{
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
		for _, src := range srcs {
			url := fmt.Sprintf(urlfmt, generations[gen].UUID)

			b, err := json.Marshal(src)
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &src)
			require.Nil(t, err)

			sources = append(sources, src)
		}
	}
}

func Test_HappyEventSource(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/generation/%%s/sources/event`, cfg.HTTPHost, cfg.HTTPPort)

	var s types.Source

	for gen, evts := range map[int][]types.Event{
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
		for _, evt := range evts {
			url := fmt.Sprintf(urlfmt, generations[gen].UUID)

			b, err := json.Marshal(types.Source{
				Lifecycle: &types.Lifecycle{
					Events: []types.Event{evt},
				},
				Type: func(s string) string {
					return s[0:5]
				}(evt.EventType.Name),
			})
			require.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode, "%d, %s", gen, b)

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
