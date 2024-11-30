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

var eventtypes = map[string]map[string]types.EventType{
	"Any":          {},
	"Gestation":    {},
	"Colonization": {},
	"Majority":     {},
	"Vacation":     {},
}

func init() {
	url := fmt.Sprintf(`http://%s:%d/eventtypes`, cfg.HTTPHost, cfg.HTTPPort)

	var et []types.EventType

	if req, err := http.NewRequest(
		http.MethodGet,
		url,
		bytes.NewReader(nil)); err != nil {
		panic(err)
	} else if req.AddCookie(cookie); false {
	} else if res, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("status was not OK: %d", res.StatusCode))
	} else if b, err := io.ReadAll(res.Body); err != nil {
		panic(err)
	} else if err = json.Unmarshal(b, &et); err != nil {
		panic(err)
	}

	for _, v := range et {
		eventtypes[v.Stage.Name][v.Name] = v
	}
}

func Test_HappyLifecycleEvent(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/lifecycle/%%s/events`, cfg.HTTPHost, cfg.HTTPPort)

	for lc, v := range map[int][]types.Event{
		0: {
			{Humidity: 12, Temperature: 77, EventType: eventtypes["Colonization"]["Innoculation"]},
			{Humidity: 22, Temperature: 76, EventType: eventtypes["Any"]["50% colonization"]},
			{Humidity: 32, Temperature: 75, EventType: eventtypes["Colonization"]["Redistribute substrate"]},
			{Humidity: 42, Temperature: 74, EventType: eventtypes["Majority"]["Binning"]},
			{Humidity: 52, Temperature: 73, EventType: eventtypes["Majority"]["Harvesting"]},
			{Humidity: 52, Temperature: 73, EventType: eventtypes["Majority"]["Spore print"]},
			{Humidity: 22, Temperature: 76, EventType: eventtypes["Any"]["Clone"]},
		},
		2: {
			{Humidity: 12, Temperature: 77, EventType: eventtypes["Colonization"]["Innoculation"]},
			{Humidity: 22, Temperature: 76, EventType: eventtypes["Vacation"]["Chill"]},
			{Humidity: 52, Temperature: 73, EventType: eventtypes["Majority"]["Spore print"]},
			{Humidity: 52, Temperature: 73, EventType: eventtypes["Majority"]["Spore print"]},
			{Humidity: 22, Temperature: 76, EventType: eventtypes["Any"]["Clone"]},
		},
	} {
		url := fmt.Sprintf(urlfmt, lifecycles[lc].UUID)
		for _, e := range v {
			b, err := json.Marshal(e)
			require.Nil(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				url,
				bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &lifecycles[lc])
			require.Nil(t, err)
		}
	}
}

func Test_HappyGenerationEvent(t *testing.T) {
	urlfmt := fmt.Sprintf(`http://%s:%d/generation/%%s/events`, cfg.HTTPHost, cfg.HTTPPort)

	for g, v := range map[int][]types.Event{
		0: {
			{Humidity: 12, Temperature: 77, EventType: eventtypes["Gestation"]["Agar sampling"]},
			{Humidity: 32, Temperature: 75, EventType: eventtypes["Any"]["50% colonization"]},
			{Humidity: 22, Temperature: 76, EventType: eventtypes["Gestation"]["Liquid innoculation"]},
		},
		2: {
			{Humidity: 12, Temperature: 77, EventType: eventtypes["Gestation"]["Agar sampling"]},
			{Humidity: 22, Temperature: 76, EventType: eventtypes["Any"]["100% colonization"]},
		},
	} {
		url := fmt.Sprintf(urlfmt, generations[g].UUID)

		for _, e := range v {
			b, err := json.Marshal(e)
			require.Nil(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				url,
				bytes.NewReader(b))
			require.Nil(t, err)
			req.AddCookie(cookie)

			res, err := http.DefaultClient.Do(req)
			require.Nil(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)

			b, err = io.ReadAll(res.Body)
			require.Nil(t, err)

			err = json.Unmarshal(b, &generations[g])
			require.Nil(t, err)
		}
	}
}
