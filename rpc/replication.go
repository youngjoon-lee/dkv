package rpc

import (
	"context"
	"io"

	"github.com/youngjoon-lee/dkv/cluster"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"github.com/youngjoon-lee/dkv/state"
	"github.com/youngjoon-lee/dkv/wal"
)

type replicationServer struct {
	pb.UnimplementedReplicationServer

	wal     wal.WAL
	state   *state.State
	cluster *cluster.Cluster
}

func (s *replicationServer) AppendLogs(stream pb.Replication_AppendLogsServer) error {
	lastAppended := wal.NilSequence
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AppendLogsReply{LastAppendedSeq: uint64(lastAppended)})
		}
		if err != nil {
			return err
		}

		if err := s.wal.AppendWith(wal.Sequence(req.Sequence), req.Msg); err != nil {
			return err
		}
		lastAppended = wal.Sequence(req.Sequence)
	}
}
func (s *replicationServer) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.CommitReply, error) {
	iter := s.wal.Iterate(s.state.LastCommitted()+1, wal.Sequence(req.ToSequence))
	lastCommitted, err := s.state.Commit(iter)
	if err != nil {
		return nil, err
	}

	return &pb.CommitReply{LastCommittedSeq: uint64(lastCommitted)}, nil
}
