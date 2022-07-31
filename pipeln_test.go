// SPDX-License-Identifier: CC0-1.0

package pipeln

import (
	"context"
	"net/http"
	"testing"

	"golang.org/x/sys/unix"

	"go.awhk.org/core"
)

func Test(s *testing.T) {
	t := core.T{T: s}

	ln := New("test:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/endpoint", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := http.Server{Handler: mux}
	go srv.Serve(ln)

	client := http.Client{Transport: &http.Transport{Dial: ln.Dial}}

	t.Run("OK", func(t *core.T) {
		resp, err := client.Get("http://test/endpoint")
		t.AssertErrorIs(nil, err)
		t.AssertEqual(http.StatusOK, resp.StatusCode)
	})

	t.Run("Address Mismatch", func(t *core.T) {
		_, err := client.Get("http://other-test/endpoint")
		t.AssertErrorIs(unix.EINVAL, err)
	})

	srv.Shutdown(context.Background())

	t.Run("Remote Connection Closed", func(t *core.T) {
		_, err := client.Get("http://test/endpoint")
		t.AssertErrorIs(unix.ECONNREFUSED, err)
	})

	t.Run("Already-closed Listener", func(t *core.T) {
		srv = http.Server{Handler: mux}
		t.AssertErrorIs(unix.EINVAL, srv.Serve(ln))
	})
}
