// SPDX-License-Identifier: CC0-1.0

package grpc_test

//go:generate protoc --go_out=. --go-grpc_out=. echo.proto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.awhk.org/pipeln"
)

type impl struct{ UnimplementedEchoServer }

var _ EchoServer = impl{}

func (impl) Echo(_ context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Message: req.Message}, nil
}

func Test(t *testing.T) {
	ln := pipeln.New("test-backend-name")

	ret := make(chan error)
	srv := grpc.NewServer()
	RegisterEchoServer(srv, &impl{})
	go func() { ret <- srv.Serve(ln) }()

	opts := []grpc.DialOption{grpc.WithContextDialer(ln.DialContextAddr), grpc.WithInsecure()}
	req := &EchoRequest{Message: "Hello World!"}

	t.Run("OK", func(t *testing.T) {
		conn, err := grpc.Dial("test-backend-name", opts...)
		require.NoError(t, err)
		defer conn.Close()

		client := NewEchoClient(conn)
		resp, err := client.Echo(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, req.Message, resp.Message)
	})

	t.Run("Address Mismatch", func(t *testing.T) {
		conn, err := grpc.Dial("bad-backend-name", opts...)
		require.NoError(t, err)
		defer conn.Close()

		client := NewEchoClient(conn)
		_, err = client.Echo(context.Background(), req)
		require.Error(t, err)
		st := status.Convert(err)
		require.NotNil(t, st)
		assert.Equal(t, codes.Unavailable, st.Code())
		assert.Contains(t, st.Message(), "invalid argument")
	})

	srv.GracefulStop()
	assert.NoError(t, <-ret)

	t.Run("Remote Connection Closed", func(t *testing.T) {
		conn, err := grpc.Dial("test-backend-name", opts...)
		require.NoError(t, err)
		defer conn.Close()

		client := NewEchoClient(conn)
		_, err = client.Echo(context.Background(), req)
		require.Error(t, err)
		st := status.Convert(err)
		require.NotNil(t, st)
		assert.Equal(t, codes.Unavailable, st.Code())
		assert.Contains(t, st.Message(), "connection refused")
	})
}
