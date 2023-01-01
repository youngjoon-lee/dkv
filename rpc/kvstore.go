package rpc

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/db"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type kvStoreServer struct {
	pb.UnimplementedKVStoreServer

	db db.DB
}

func (s *kvStoreServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutReply, error) {
	log.Infof("key:%v(%s), value:%v(%s)", req.Key, string(req.Key), req.Value, string(req.Value))
	if err := s.db.Put(req.Key, req.Value); err != nil {
		return nil, fmt.Errorf("failed to put: %w", err)
	}
	return &pb.PutReply{Message: "success"}, nil
}

