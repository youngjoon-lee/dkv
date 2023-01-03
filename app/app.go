package app

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	"github.com/youngjoon-lee/dkv/config"
	"github.com/youngjoon-lee/dkv/replication"
	"github.com/youngjoon-lee/dkv/rpc"
	"github.com/youngjoon-lee/dkv/state"
	"github.com/youngjoon-lee/dkv/wal"
)

type App struct {
	conf       config.Config
	cluster    *cluster.Cluster
	wal        wal.WAL
	state      *state.State
	replicator *replication.Replicator
	rpcSvr     rpc.Server
}

func New(conf config.Config) (*App, error) {
	cluster, err := cluster.New(conf.NodeID, conf.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to init cluster info: %w", err)
	}
	log.Debugf("cluster info initialized: %v", cluster.Status())

	wal := wal.NewMemWAL()

	state, err := state.New(conf.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init state: %w", err)
	}

	var replicator *replication.Replicator
	if cluster.IsLeader() {
		log.Info("starting replicator as a leader...")
		replicator = replication.NewReplicator(wal, state)
	}

	rpcSvr, err := rpc.Serve(conf.RPCPort, conf.RESTPort, wal, state, cluster)
	if err != nil {
		state.Close()
		return nil, fmt.Errorf("failed to serve RPC: %w", err)
	}

	return &App{
		conf:       conf,
		cluster:    cluster,
		wal:        wal,
		state:      state,
		replicator: replicator,
		rpcSvr:     rpcSvr,
	}, nil
}

func (s *App) Close() {
	log.Info("closing RPC...")
	s.rpcSvr.GracefulStop()

	if s.replicator != nil {
		log.Info("stopping replicator...")
		s.replicator.Stop()
	}

	log.Info("closing state...")
	s.state.Close()
}
