package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type GRPCServer struct {
	pb.UnimplementedKVStoreServer
}

func (s *GRPCServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetReply, error) {
	log.Infof("key:%v(%s), value:%v(%s)", req.Key, string(req.Key), req.Value, string(req.Value))
	return &pb.SetReply{Message: "success"}, nil
}
