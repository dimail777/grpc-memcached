package app

import (
	"context"
	"internal/memcache"
)

type memcachedGrpcHandler struct {
	client memcache.Client
}

func (s *memcachedGrpcHandler) mustEmbedUnimplementedICacheServer() {}

func (s *memcachedGrpcHandler) Get(ctx context.Context, req *GetCommand) (*GetResult, error) {
	deadline, _ := ctx.Deadline()
	value, err := s.client.Get(req.GetKey(), deadline)
	if err != nil {
		return nil, err
	}
	return &GetResult{Key: req.GetKey(), Value: value}, nil
}

func (s *memcachedGrpcHandler) Set(ctx context.Context, req *SetCommand) (*SetResult, error) {
	deadline, _ := ctx.Deadline()
	err := s.client.Set(req.GetKey(), req.GetValue(), deadline)
	if err != nil {
		return nil, err
	}
	return &SetResult{Key: req.GetKey(), Done: true}, nil
}

func (s *memcachedGrpcHandler) Del(ctx context.Context, req *DelCommand) (*DelResult, error) {
	deadline, _ := ctx.Deadline()
	err := s.client.Del(req.GetKey(), deadline)
	if err != nil {
		return nil, err
	}
	return &DelResult{Key: req.GetKey(), Done: true}, nil
}
