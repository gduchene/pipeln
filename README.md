`pipeln` implements a trivial type, `PipeListenerDialer`, that can be
used both as a `net.Listener` and as a dialer. It uses `net.Pipe` to
connect server and clients so that testing client-server communication
becomes easier.

Several dialer methods are available, and happen (as there is no
`net.Dialer` interface) to be compatible with both `net.Transport` and
`grpc.WithContextDialer`.

For instance:

```go
func TestHTTP(t *testing.T) {
	ln := pipeln.New("test:80")

	srv := http.Server{}
	go srv.Serve(ln)

	client := http.Client{Transport: &http.Transport{DialContext: ln.DialContext}}

	// ...
}

func TestGRPC(t *testing.T) {
	ln := pipeln.New("test")

	srv := grpc.NewServer()
	go srv.Serve(ln)

	client, _ := grpc.Dial("test", grpc.WithContextDialer(ln.DialContextAddr), grpc.WithInsecure())

	// ...
}
```
