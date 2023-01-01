package rpc

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	"github.com/youngjoon-lee/dkv/db"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type kvStoreServer struct {
	pb.UnimplementedKVStoreServer

	db      db.DB
	cluster *cluster.Cluster
}

func (s *kvStoreServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutReply, error) {
	if !s.cluster.IsLeader() {
		return nil, fmt.Errorf("only leader(%v) can handle write operations", s.cluster.Status().Leader)
	}

	log.Debugf("key:%v(%s), value:%v(%s)", req.Key, string(req.Key), req.Value, string(req.Value))
	if err := s.db.Put(req.Key, req.Value); err != nil {
		return nil, fmt.Errorf("failed to put: %w", err)
	}
	return &pb.PutReply{Message: "success"}, nil
}
