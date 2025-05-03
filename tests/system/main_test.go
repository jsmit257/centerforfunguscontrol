package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/userservice/shared/v1"
)

var (
	cookie *http.Cookie = &http.Cookie{}

	cfg *config.Config
)

func init() {
	var err error
	if cfg == nil { // was it better when NewConfig panicked
		cfg, err = config.NewConfig()
		if err != nil {
			panic(err)
		}
	}

	if cfg.AuthnHost == "" || cfg.AuthnPort == 0 {
		return
	}

	addr := shared.Email("foobar@example.com")
	auth := shared.BasicAuth{Name: "testuser"}
	b, _ := json.Marshal(shared.User{
		Name:  auth.Name,
		Email: &addr,
	})
	// create a user
	url := fmt.Sprintf("http://%s:%d/user", cfg.AuthnHost, cfg.AuthnPort)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusCreated {
		panic(fmt.Errorf("code isn't created: %d %s", resp.StatusCode, resp.Status))
	} else if id, err := io.ReadAll(resp.Body); err != nil {
		panic(err)
	} else if auth.UUID = shared.UUID(id); auth.UUID == "" {
		panic(fmt.Errorf("no userid returned"))
	}

	<-time.After(time.Second / 2)

	b, _ = json.Marshal(auth)
	// log new user in to get cookie for all subsequent requests
	url = fmt.Sprintf("http://%s:%d/auth", cfg.AuthnHost, cfg.AuthnPort)
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	} else if resp, err := (&http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}).Do(req); err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusMovedPermanently {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Errorf("%s, code isn't OK: sc: %d, body: %s, cookie: '%s', resp: %s",
			url,
			resp.StatusCode,
			string(b),
			resp.Header.Get("Set-Cookie"),
			body))
	} else if header := resp.Header.Get("Set-Cookie"); len(header) == 0 {
		panic("no Set-Cookie header found")
	} else if cookie, err = http.ParseSetCookie(header); err != nil {
		panic(fmt.Errorf("cookie was unparseable %s", header))
	}
}
