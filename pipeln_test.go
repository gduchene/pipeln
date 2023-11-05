// SPDX-License-Identifier: CC0-1.0

package pipeln

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"testing"
)

func Test(t *testing.T) {
	ln := New("test:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/endpoint", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := http.Server{Handler: mux}
	go srv.Serve(ln)

	client := http.Client{Transport: &http.Transport{Dial: ln.Dial}}

	t.Run("OK", func(t *testing.T) {
		resp, err := client.Get("http://test/endpoint")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Error(resp.StatusCode)
		}
	})

	t.Run("Address Mismatch", func(t *testing.T) {
		_, err := client.Get("http://other-test/endpoint")
		if !errors.Is(err, syscall.EINVAL) {
			t.Fatal(err)
		}
	})

	srv.Shutdown(context.Background())

	t.Run("Remote Connection Closed", func(t *testing.T) {
		_, err := client.Get("http://test/endpoint")
		if !errors.Is(err, syscall.ECONNREFUSED) {
			t.Fatal(err)
		}
	})

	t.Run("Already-closed Listener", func(t *testing.T) {
		srv = http.Server{Handler: mux}
		if err := srv.Serve(ln); !errors.Is(err, syscall.EINVAL) {
			t.Fatal(err)
		}
	})
}
