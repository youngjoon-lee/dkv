package server

import (
	"context"

	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type GRPCServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GRPCServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return nil, nil
}
