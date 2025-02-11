package metrics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type (
	mockContextHandler struct {
		sc int
	}

	mockRoutes struct {
		routes  []chi.Route
		middles chi.Middlewares
		matched bool
		found   string
	}
)

func Test_WrapContext(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		mr *mockRoutes
		sc int
	}{
		"happy_path": {
			sc: http.StatusTeapot,
			mr: &mockRoutes{},
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wrap := WrapContext(logrus.WithField("test", name))
			handle := wrap(&mockContextHandler{sc: tc.sc})
			w := httptest.NewRecorder()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.TODO(),
					chi.RouteCtxKey,
					&chi.Context{
						URLParams: chi.RouteParams{Keys: []string{"user_id"}, Values: []string{name}},
						Routes:    tc.mr,
					}),
				http.MethodGet,
				"tc.url",
				nil,
			)
			r.URL.RawPath = "not-empty"

			handle.ServeHTTP(w, r)
			require.Equal(t, tc.sc, w.Result().StatusCode)
		})
	}
}

func (f *mockContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(f.sc)
}

func (mr *mockRoutes) Routes() []chi.Route {
	return mr.routes
}
func (mr *mockRoutes) Middlewares() chi.Middlewares {
	return mr.middles
}
func (mr *mockRoutes) Match(rctx *chi.Context, method, path string) bool {
	return mr.matched
}
func (mr *mockRoutes) Find(rctx *chi.Context, method, path string) string {
	return mr.found
}
