package cluster

import (
	"fmt"
	"sync"

	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	nodeID string

	mutex   sync.RWMutex
	cluster *pb.Cluster
}

func New(nodeID string, nodes []*pb.Node) (*Cluster, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no node specified")
	}

	leader := nodes[0]
	followers := make(map[string]*pb.Node, len(nodes)-1)
	for i := 1; i < len(nodes); i++ {
		followers[nodes[i].Id] = nodes[i]
	}

	return &Cluster{
		nodeID: nodeID,
		cluster: &pb.Cluster{
			Leader:    leader,
			Followers: followers,
		},
	}, nil
}

func (c *Cluster) IsLeader() bool {
	return c.nodeID == c.cluster.Leader.Id
}

func (c *Cluster) NodeID() string {
	return c.nodeID
}

func (c *Cluster) Status() *pb.Cluster {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return proto.Clone(c.cluster).(*pb.Cluster)
}
