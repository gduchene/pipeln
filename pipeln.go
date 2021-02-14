// SPDX-License-Identifier: CC0-1.0

package pipeln

import (
	"context"
	"errors"
	"net"
)

var (
	ErrBadAddress = errors.New("bad address")
	ErrClosed     = errors.New("closed listener")
)

type addr struct {
	ln *PipeListenerDialer
}

var _ net.Addr = addr{}

func (addr) Network() string {
	return "pipe"
}

func (a addr) String() string {
	return "pipe:" + a.ln.addr
}

// PipeListener can be used to simulate client-server interaction within
// the same process. Useful for testing. Somehow the Go standard library
// provides net.Pipe but no net.PipeListener.
type PipeListenerDialer struct {
	addr  string
	conns chan net.Conn
	done  chan struct{}
}

var _ net.Listener = &PipeListenerDialer{}

func (ln *PipeListenerDialer) Accept() (net.Conn, error) {
	select {
	case conn := <-ln.conns:
		return conn, nil
	case <-ln.done:
		return nil, ErrClosed
	}
}

func (ln *PipeListenerDialer) Addr() net.Addr {
	return addr{ln}
}

func (ln *PipeListenerDialer) Close() error {
	close(ln.done)
	return nil
}

func (ln *PipeListenerDialer) Dial(_, addr string) (net.Conn, error) {
	if addr != ln.addr {
		return nil, ErrBadAddress
	}
	s, c := net.Pipe()
	select {
	case ln.conns <- s:
		return c, nil
	case <-ln.done:
		return nil, ErrClosed
	}
}

func (ln *PipeListenerDialer) DialContext(_ context.Context, network, addr string) (net.Conn, error) {
	return ln.Dial(network, addr)
}

func (ln *PipeListenerDialer) DialContextAddr(_ context.Context, addr string) (net.Conn, error) {
	return ln.Dial("", addr)
}

func New(addr string) *PipeListenerDialer {
	return &PipeListenerDialer{addr, make(chan net.Conn), make(chan struct{})}
}
