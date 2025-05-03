package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
)

type mockHandler struct {
	sc   int
	body []byte
}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(mh.sc)
	_, _ = w.Write(mh.body)
}

func Test_authn(t *testing.T) {
	tcs := map[string]struct {
		host     string
		port     uint16
		cookie   *http.Cookie
		location string
		mh       *mockHandler
		sc       int
		code     int
	}{
		"happy_path": {
			host: "Test_authn",
			port: 1313,
			cookie: &http.Cookie{
				Name:  "us-authn",
				Value: "happy_path",
			},
			mh:   &mockHandler{sc: http.StatusOK},
			sc:   http.StatusNoContent,
			code: http.StatusOK,
		},
		"no_auth": {
			host: "Test_authn",
			port: 1314,
			cookie: &http.Cookie{
				Name:  "us-authn",
				Value: "no_auth",
			},
			location: "/location",
			sc:       http.StatusTemporaryRedirect,
			code:     http.StatusForbidden,
		},
		"service_error": {
			host: "\t",
			code: http.StatusInternalServerError,
		},
		"unexpected_error": {
			host: "Test_authn",
			port: 1313,
			cookie: &http.Cookie{
				Name:  "us-authn",
				Value: "weird_status",
			},
			mh:   &mockHandler{sc: http.StatusOK},
			sc:   http.StatusBadGateway,
			code: http.StatusBadGateway,
		},
	}

	for name, tc := range tcs {
		// name, tc := name, tc // do NOT parallelize

		t.Run(name, func(t *testing.T) {
			handler := authn(tc.host, tc.port, "logon")(tc.mh)

			log := logrus.WithField("test", name)
			w := httptest.NewRecorder()
			r := httptest.NewRequestWithContext(
				context.WithValue(context.TODO(), metrics.Log, log),
				http.MethodGet,
				"/valid",
				nil)

			httpmock.RegisterResponder(http.MethodGet,
				fmt.Sprintf("http://%s:%d/valid", tc.host, tc.port),
				func(r *http.Request) (*http.Response, error) {
					resp := httpmock.NewBytesResponse(tc.sc, nil)
					if tc.cookie != nil {
						resp.Header.Set("Set-Cookie", tc.cookie.String())
					}
					resp.Header.Set("Location", tc.location)
					return resp, nil
				})
			httpmock.Activate()
			defer httpmock.Deactivate()

			handler.ServeHTTP(w, r)

			if tc.cookie != nil {
				require.Contains(t, tc.cookie.String(), w.Header().Get("Set-Cookie"))
			}
			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.location, w.Header().Get("Location"))
		})
	}
}

func Test_newHuautla(t *testing.T) {
	// TODO: give it a whirl
	newHuautla(&config.Config{
		AuthnHost: "Test_newHuautla",
		AuthnPort: 12000,
	}, nil, logrus.WithField("test", "Test_newHuautla"))
}

func Test_Settings(t *testing.T) {
	t.Parallel()

	handler := settings(&config.Config{
		AuthnPath:   "foobar",
		AuthnPort:   1234,
		HuautlaHost: "quux",
	})

	tcs := map[string]struct {
		name  string
		value map[string]interface{}
		sc    int
	}{
		"get_all": {
			name: "*",
			value: map[string]interface{}{
				"authn_path":   "foobar",
				"authn_port":   float64(1234),
				"huautla_host": "quux",
			},
			sc: http.StatusOK,
		},
		"get_string": {
			name:  "authn_path",
			value: map[string]interface{}{"value": "foobar"},
			sc:    http.StatusOK,
		},
		"get_number": {
			name:  "authn_port",
			value: map[string]interface{}{"value": float64(1234)},
			sc:    http.StatusOK,
		},
		"not_found": {
			name:  "bad_key",
			value: map[string]interface{}{"error": "no value for key: bad_key"},
			sc:    http.StatusNotFound,
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"name"}, Values: []string{tc.name}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				nil)

			handler(w, r)

			body, err := io.ReadAll(w.Body)
			require.Nil(t, err)
			t.Log(string(body))
			var result any
			err = json.Unmarshal(body, &result)
			require.Nil(t, err, string(body))

			require.Equal(t, tc.sc, w.Code)
			require.Equal(t, tc.value, result)
		})
	}
}
