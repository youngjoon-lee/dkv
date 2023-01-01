package rpc

import (
	"context"

	"github.com/youngjoon-lee/dkv/cluster"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type clusterServiceServer struct {
	pb.UnimplementedClusterServiceServer

	cluster *cluster.Cluster
}

func (s *clusterServiceServer) Status(ctx context.Context, req *pb.StatusRequest) (*pb.StatusReply, error) {
	return &pb.StatusReply{
		NodeId:  s.cluster.NodeID(),
		Cluster: s.cluster.Status(),
	}, nil
}
