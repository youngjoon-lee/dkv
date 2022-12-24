package server

import (
	"context"
	"fmt"

	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type GRPCServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GRPCServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("hello, %s", req.Name)}, nil
}
