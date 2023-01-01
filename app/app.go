package app

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	"github.com/youngjoon-lee/dkv/config"
	"github.com/youngjoon-lee/dkv/db"
	"github.com/youngjoon-lee/dkv/rpc"
)

type App struct {
	conf    config.Config
	db      db.DB
	cluster *cluster.Cluster
	rpcSvr  rpc.Server
}

func New(conf config.Config) (*App, error) {
	db, err := db.NewBoltDB(conf.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init DB: %w", err)
	}

	cluster, err := cluster.New(conf.NodeID, conf.Cluster)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to init cluster info: %w", err)
	}
	log.Debugf("cluster info initialized: %v", cluster.Status())

	rpcSvr, err := rpc.Serve(conf.RPCPort, conf.RESTPort, db, cluster)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to serve RPC: %w", err)
	}

	return &App{
		conf:    conf,
		db:      db,
		cluster: cluster,
		rpcSvr:  rpcSvr,
	}, nil
}

func (s *App) Close() {
	log.Info("closing RPC...")
	s.rpcSvr.GracefulStop()

	log.Info("closing DB...")
	s.db.Close()
}
