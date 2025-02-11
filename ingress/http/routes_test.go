package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
)

type mockHandler struct{}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func Test_authn(t *testing.T) {

	tcs := map[string]struct {
		host     string
		port     uint16
		cookie   http.Cookie
		response *http.Response
		err      error
		pad      string // ought to be used in responses.Header["Authn-Pad"]
		sc       int
	}{
		"happy_path": {
			host: "Test_authn",
			port: 1313,
			cookie: http.Cookie{
				Name:  "us-authn",
				Value: "hmmm",
				// ...
			},
			response: &http.Response{
				Request: &http.Request{RequestURI: "/otp/12345"},
				Header: http.Header{
					"Authn-Pad": []string{"123"},
				},
			},
			// pad: "123",
			sc: http.StatusFound,
		},
		"missing_uri": {
			host:   "Test_authn",
			port:   1313,
			cookie: http.Cookie{Name: "us-authn"},
			response: &http.Response{
				Request: &http.Request{RequestURI: "/foobar"},
			},
			sc: http.StatusFound,
		},
		// this just covers some temporary logging i want to remove
		"response_nil": {
			host: "Test_authn",
			port: 1313,
			cookie: http.Cookie{
				Name:  "us-authn",
				Value: "hmmm",
			},
			sc: http.StatusFound,
		},
		"not_valid": {
			host:   "Test_authn",
			port:   1313,
			cookie: http.Cookie{Name: "us-authn"},
			sc:     http.StatusForbidden,
		},
		"no_cookie": {
			host: "Test_authn",
			port: 1313,
			sc:   http.StatusFound,
		},
	}

	for name, tc := range tcs {
		// name, tc := name, tc // do NOT parallelize

		t.Run(name, func(t *testing.T) {
			mh := &mockHandler{}
			next := http.Handler(mh)
			wrapper := authn(tc.host, tc.port)
			handler := wrapper(next)

			w := httptest.NewRecorder()
			r := httptest.NewRequestWithContext(context.WithValue(
				context.TODO(),
				metrics.Log,
				logrus.WithField("test", name)),
				http.MethodGet,
				"/valid",
				nil,
			)
			r.Response = tc.response
			if name != "no_cookie" {
				r.AddCookie(&tc.cookie)
			}

			httpmock.RegisterResponder(http.MethodGet,
				fmt.Sprintf("http://%s:%d/valid", tc.host, tc.port),
				func(r *http.Request) (*http.Response, error) {
					resp := httpmock.NewBytesResponse(tc.sc, nil)
					resp.Header.Set("Set-Cookie", tc.cookie.String())
					return resp, tc.err
				})
			httpmock.Activate()
			defer httpmock.Deactivate()

			handler.ServeHTTP(w, r)

			require.Equal(t, tc.pad, w.Header().Get("Authn-Pad"))
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
