package rpc

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"github.com/youngjoon-lee/dkv/state"
	"github.com/youngjoon-lee/dkv/wal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type kvStoreServer struct {
	pb.UnimplementedKVStoreServer

	wal     wal.WAL
	state   *state.State
	cluster *cluster.Cluster
}

func (s *kvStoreServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutReply, error) {
	log.Debugf("received a put request: %v", req)

	if !s.cluster.IsLeader() {
		return nil, fmt.Errorf("only leader(%v) can handle write operations", s.cluster.Status().Leader)
	}

	doneCh := make(chan error)
	if _, err := s.wal.Append(req, doneCh); err != nil {
		return nil, fmt.Errorf("failed to append to WAL: %w", err)
	}

	if err := <-doneCh; err != nil {
		return nil, fmt.Errorf("failed to put KV: %w", err)
	}

	log.Debugf("put succeeded: %v", req)
	return &pb.PutReply{Message: "success"}, nil
}

func (s *kvStoreServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	log.Debugf("received a get request: %v", req)

	value, err := s.state.Get(req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get value from db: %w", err)
	}
	if value == nil {
		return nil, status.Errorf(codes.NotFound, "key is not found in the DB")
	}

	return &pb.GetReply{Key: req.Key, Value: value}, nil
}
