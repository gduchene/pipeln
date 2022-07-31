module go.awhk.org/pipeln/internal/grpc_test

go 1.18

require go.awhk.org/pipeln v0.0.0

replace go.awhk.org/pipeln => ../../

require (
	go.awhk.org/core v0.2.0
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	golang.org/x/net v0.0.0-20220728211354-c7608f3a8462 // indirect
	golang.org/x/sys v0.0.0-20220730100132-1609e554cd39 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220728213248-dd149ef739b9 // indirect
)
