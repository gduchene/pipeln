// SPDX-License-Identifier: CC0-1.0

package pipeln

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	resp, err := client.Get("http://test/endpoint")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	srv.Shutdown(context.Background())
}
