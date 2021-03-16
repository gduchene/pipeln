// SPDX-License-Identifier: CC0-1.0

package pipeln

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
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
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Address Mismatch", func(t *testing.T) {
		_, err := client.Get("http://other-test/endpoint")
		assert.ErrorIs(t, err, unix.EINVAL)
	})

	srv.Shutdown(context.Background())

	t.Run("Remote Connection Closed", func(t *testing.T) {
		_, err := client.Get("http://test/endpoint")
		assert.ErrorIs(t, err, unix.ECONNREFUSED)
	})

	t.Run("Already-closed Listener", func(t *testing.T) {
		srv = http.Server{Handler: mux}
		assert.ErrorIs(t, srv.Serve(ln), unix.EINVAL)
	})
}
