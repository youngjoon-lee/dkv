package server

import (
	"context"

	"github.com/youngjoon-lee/dkv/types"
)

type GRPCServer struct {
	types.UnimplementedGreeterServer
}

func (s *GRPCServer) SayHello(ctx context.Context, req *types.HelloRequest) (*types.HelloReply, error) {
	return nil, nil
}
