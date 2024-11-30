package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/userservice/shared/v1"
)

var (
	cookie *http.Cookie

	cfg = config.NewConfig()
)

func init() {
	if cfg.AuthnHost == "" || cfg.AuthnPort == 0 {
		return
	}

	auth := shared.BasicAuth{Name: "testuser"}
	b, _ := json.Marshal(auth)
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

	b, _ = json.Marshal(auth)
	url = fmt.Sprintf("http://%s:%d/auth", cfg.AuthnHost, cfg.AuthnPort)
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("code isn't OK: %s, %s, %v", url, resp.Status, auth))
	} else if header := resp.Header.Get("Set-Cookie"); len(header) == 0 {
		panic("no Set-Cookie header found")
	} else if cookie, err = http.ParseSetCookie(header); err != nil {
		panic(fmt.Errorf("cookie was unparseable %s", header))
	}
}
