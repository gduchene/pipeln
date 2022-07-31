// SPDX-License-Identifier: CC0-1.0

package grpc_test

//go:generate protoc --go_out=. --go-grpc_out=. echo.proto

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.awhk.org/core"
	"go.awhk.org/pipeln"
)

type impl struct{ UnimplementedEchoServer }

var _ EchoServer = impl{}

func (impl) Echo(_ context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Message: req.Message}, nil
}

func Test(s *testing.T) {
	t := core.T{T: s}

	ln := pipeln.New("test-backend-name")
	ret := make(chan error)
	srv := grpc.NewServer()
	RegisterEchoServer(srv, &impl{})
	go func() { ret <- srv.Serve(ln) }()

	opts := []grpc.DialOption{grpc.WithContextDialer(ln.DialContextAddr), grpc.WithInsecure()}
	req := &EchoRequest{Message: "Hello World!"}

	t.Run("OK", func(t *core.T) {
		conn, err := grpc.Dial("test-backend-name", opts...)
		t.Must(t.AssertErrorIs(nil, err))
		defer conn.Close()

		client := NewEchoClient(conn)
		resp, err := client.Echo(context.Background(), req)
		t.AssertErrorIs(nil, err)
		t.AssertEqual(req.Message, resp.Message)
	})

	t.Run("Address Mismatch", func(t *core.T) {
		conn, err := grpc.Dial("bad-backend-name", opts...)
		t.Must(t.AssertErrorIs(nil, err))
		defer conn.Close()

		client := NewEchoClient(conn)
		_, err = client.Echo(context.Background(), req)
		st := status.Convert(err)
		t.Must(t.AssertNotEqual(nil, st))
		t.AssertEqual(codes.Unavailable, st.Code())
		t.Assert(strings.Contains(st.Message(), "invalid argument"))
	})

	srv.GracefulStop()
	t.Must(t.AssertErrorIs(nil, <-ret))

	t.Run("Remote Connection Closed", func(t *core.T) {
		conn, err := grpc.Dial("test-backend-name", opts...)
		t.Must(t.AssertErrorIs(nil, err))
		defer conn.Close()

		client := NewEchoClient(conn)
		_, err = client.Echo(context.Background(), req)
		st := status.Convert(err)
		t.Must(t.AssertNotEqual(nil, st))
		t.AssertEqual(codes.Unavailable, st.Code())
		t.Assert(strings.Contains(st.Message(), "connection refused"))
	})
}
