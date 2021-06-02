// SPDX-License-Identifier: CC0-1.0

package pipeln

import (
	"context"
	"net"

	"golang.org/x/sys/unix"
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

// PipeListenerDialer can be used to simulate client-server interaction
// within the same process.
type PipeListenerDialer struct {
	addr  string
	conns chan net.Conn
	done  chan struct{}
	ok    bool
}

var _ net.Listener = &PipeListenerDialer{}

// See net.Listener.Accept for more details.
func (ln *PipeListenerDialer) Accept() (net.Conn, error) {
	select {
	case conn := <-ln.conns:
		return conn, nil
	case <-ln.done:
		return nil, unix.EINVAL
	}
}

// See net.Listener.Addr for more details.
func (ln *PipeListenerDialer) Addr() net.Addr {
	return addr{ln}
}

// See net.Listener.Close for more details.
func (ln *PipeListenerDialer) Close() error {
	if !ln.ok {
		return unix.EINVAL
	}
	close(ln.done)
	ln.ok = false
	return nil
}

// See net.Dialer.Dial for more details.
func (ln *PipeListenerDialer) Dial(_, addr string) (net.Conn, error) {
	if addr != ln.addr {
		return nil, unix.EINVAL
	}
	s, c := net.Pipe()
	select {
	case ln.conns <- s:
		return c, nil
	case <-ln.done:
		return nil, unix.ECONNREFUSED
	}
}

// DialContext is a dummy wrapper around Dial.
func (ln *PipeListenerDialer) DialContext(_ context.Context, network, addr string) (net.Conn, error) {
	return ln.Dial(network, addr)
}

// DialContextAddr is a dummy wrapper around Dial.
//
// This function can be passed to grpc.WithContextDialer.
func (ln *PipeListenerDialer) DialContextAddr(_ context.Context, addr string) (net.Conn, error) {
	return ln.Dial("", addr)
}

// New returns a PipeListenerDialer that will only accept connections
// made to the given addr.
func New(addr string) *PipeListenerDialer {
	return &PipeListenerDialer{addr, make(chan net.Conn), make(chan struct{}), true}
}
