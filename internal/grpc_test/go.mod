module go.awhk.org/pipeln/internal/grpc_test

go 1.18

require go.awhk.org/pipeln v1.0.5

replace go.awhk.org/pipeln => ../../

require (
	go.awhk.org/core v0.7.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
)
