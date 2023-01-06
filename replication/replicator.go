package replication

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"github.com/youngjoon-lee/dkv/state"
	"github.com/youngjoon-lee/dkv/wal"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Replicator struct {
	stopCh chan struct{}

	wal     wal.WAL
	state   *state.State
	cluster *cluster.Cluster
}

func NewReplicator(wal wal.WAL, state *state.State, cluster *cluster.Cluster) *Replicator {
	r := &Replicator{
		stopCh:  make(chan struct{}),
		wal:     wal,
		state:   state,
		cluster: cluster,
	}

	go r.start()

	return r
}

func (r *Replicator) start() {
	for {
		select {
		case <-r.stopCh:
			r.stopCh <- struct{}{}
		default:
			if err := r.runIteration(); err != nil {
				log.Error(err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (r *Replicator) Stop() {
	r.stopCh <- struct{}{}
	<-r.stopCh
}

func (r *Replicator) runIteration() error {
	fromSeq := r.state.LastCommitted() + 1
	followers := r.cluster.Status().Followers

	// broadcast AppendLogs to follower
	// wait for responses from followers
	lastAppendedSeq, err := r.sendAppendLogsToAll(fromSeq, followers)
	if err != nil {
		return fmt.Errorf("failed to replicate: %w", err)
	}
	log.Debugf("replicated to seq %v", lastAppendedSeq)

	toSeq := lastAppendedSeq
	iter := r.wal.Iterate(fromSeq, toSeq)
	lastCommitted, err := r.state.Commit(iter)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	} else if fromSeq <= lastCommitted {
		log.Debugf("committed to seq %v~%v", fromSeq, lastCommitted)
	}

	//  broadcast commit to followers
	if err = r.sendCommitToAll(toSeq, followers); err != nil {
		return fmt.Errorf("failed to commit for followers: %w", err)
	}

	return nil
}

func (r *Replicator) sendAppendLogsToAll(fromSeq wal.Sequence, followers map[string]*pb.Node) (wal.Sequence, error) {
	results := make(map[string]wal.Sequence, len(followers))

	g, _ := errgroup.WithContext(context.Background())

	for id, follower := range followers {
		g.Go(func() error {
			lastAppendedSeq, err := r.sendAppendLogs(fromSeq, follower)
			if err != nil {
				log.Error(err)
				return err
			}

			results[id] = lastAppendedSeq
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return wal.NilSequence, fmt.Errorf("failed from one of followers: %w", err)
	}

	minLastAppendedSeq := wal.NilSequence
	for _, result := range results {
		if minLastAppendedSeq == wal.NilSequence || minLastAppendedSeq > result {
			minLastAppendedSeq = result
		}
	}
	return minLastAppendedSeq, nil
}

func (r *Replicator) sendAppendLogs(fromSeq wal.Sequence, follower *pb.Node) (wal.Sequence, error) {
	iter := r.wal.Iterate(fromSeq, wal.NilSequence)
	if iter == nil {
		log.Debugf("nothing to iterate")
		return wal.NilSequence, nil
	}

	conn, err := grpc.Dial(follower.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return wal.NilSequence, fmt.Errorf("failed to dial to %v: %w", follower.Addr, err)
	}
	defer conn.Close()

	stream, err := pb.NewReplicationClient(conn).AppendLogs(context.Background())
	if err != nil {
		return wal.NilSequence, fmt.Errorf("failed to open a stream to %v: %w", follower.Addr, err)
	}

	for {
		elem := iter.Next()
		if elem == nil {
			break
		}

		log.Debugf("sending append logs from seq %v", elem.Seq)
		req := &pb.AppendLogRequest{Sequence: uint64(elem.Seq), Msg: elem.Msg}
		if err := stream.Send(req); err != nil {
			return wal.NilSequence, err
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		return wal.NilSequence, err
	}
	log.Debugf("replication reply: %v", reply)

	return wal.Sequence(reply.LastAppendedSeq), nil
}

func (r *Replicator) sendCommitToAll(toSeq wal.Sequence, followers map[string]*pb.Node) error {
	var wg sync.WaitGroup

	for _, follower := range followers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := r.sendCommit(toSeq, follower); err != nil {
				log.Errorf("failed to commit to the follower %v to the seq %v: %w", follower, toSeq, err)
				return
			}
			log.Debugf("committed successfully: follower:%v, toSeq:%v", follower, toSeq)
		}()
	}

	wg.Wait()
	return nil
}

func (r *Replicator) sendCommit(toSeq wal.Sequence, follower *pb.Node) error {
	conn, err := grpc.Dial(follower.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to dial to %v: %w", follower.Addr, err)
	}
	defer conn.Close()

	req := &pb.CommitRequest{ToSequence: uint64(toSeq)}
	reply, err := pb.NewReplicationClient(conn).Commit(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to send commit to %v: %w", follower.Addr, err)
	}
	log.Debugf("commit reply: %v", reply)

	return nil
}
