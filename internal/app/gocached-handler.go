package app

import (
	"context"
	"internal/gocache"
)

type goCachedGrpcHandler struct {
	container gocache.Container
}

func (s *goCachedGrpcHandler) mustEmbedUnimplementedICacheServer() {}

func (s *goCachedGrpcHandler) Get(_ context.Context, req *GetCommand) (*GetResult, error) {
	value := s.container.Get(req.GetKey())
	return &GetResult{Key: req.GetKey(), Value: value}, nil
}

func (s *goCachedGrpcHandler) Set(_ context.Context, req *SetCommand) (*SetResult, error) {
	s.container.Set(req.GetKey(), req.GetValue())
	return &SetResult{Key: req.GetKey(), Done: true}, nil
}

func (s *goCachedGrpcHandler) Del(_ context.Context, req *DelCommand) (*DelResult, error) {
	s.container.Del(req.GetKey())
	return &DelResult{Key: req.GetKey(), Done: true}, nil
}
