module grpc.com/memcached

go 1.20

replace internal/app => ./internal/app

replace internal/gocache => ./internal/pkg/gocache

replace internal/memcache => ./internal/pkg/memcache

require (
	internal/app v0.0.0-00010101000000-000000000000
	internal/memcache v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	internal/gocache v0.0.0-00010101000000-000000000000 // indirect
)
