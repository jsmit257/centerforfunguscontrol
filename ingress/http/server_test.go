package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
)

func Test_startServer(t *testing.T) {
	router := chi.NewMux()
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wg := &sync.WaitGroup{}

	host, port := "localhost", rand.IntN(math.MaxUint16)
	go startServer(&config.Config{
		HTTPHost: host,
		HTTPPort: port,
	},
		router,
		wg,
		logrus.WithField("test", "Test_startServer"))

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s:%d/test", host, port),
		nil)
	require.Nil(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.Nilf(t, err, "%q", err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	err = syscall.Kill(os.Getpid(), syscall.SIGINT)
	require.Nil(t, err)

	select {
	case <-c:
	case <-time.After(6 * time.Second):
		require.Fail(t, "shouldn't have timed out")
	}
}
