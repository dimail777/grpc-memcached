package app

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"internal/gocache"
	"internal/memcache"
	"net"
	"time"
)

var kacp = keepalive.ServerParameters{
	Time:              10 * time.Second,
	Timeout:           5 * time.Second,
	MaxConnectionIdle: 5 * time.Minute,
}

func RunMemcached(connConfig *memcache.ConnConfig, poolConfig *memcache.PoolConfig) {
	handler := &memcachedGrpcHandler{client: memcache.InitClient(*connConfig, *poolConfig)}
	run(handler)
}

func RunGoCached(concurrency int) {
	handler := &goCachedGrpcHandler{container: gocache.InitContainer(concurrency)}
	run(handler)
}

func run(handler ICacheServer) {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.KeepaliveParams(kacp))
	RegisterICacheServer(server, handler)
	fmt.Println("Server started at :50051")
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}
