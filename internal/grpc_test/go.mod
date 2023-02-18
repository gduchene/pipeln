module go.awhk.org/pipeln/internal/grpc_test

go 1.18

require go.awhk.org/pipeln v0.0.0

replace go.awhk.org/pipeln => ../../

require (
	go.awhk.org/core v0.5.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230216225411-c8e22ba71e44 // indirect
)
