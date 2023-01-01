package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	"github.com/youngjoon-lee/dkv/config"
	"github.com/youngjoon-lee/dkv/db"
	"github.com/youngjoon-lee/dkv/rpc"
)

type Service struct {
	conf       config.Config
	db         db.DB
	clusterMap *cluster.Cluster
	rpcSvr     rpc.Server
}

func New(conf config.Config) (*Service, error) {
	db, err := db.NewBoltDB(conf.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init DB: %w", err)
	}

	clusterMap, err := cluster.NewClusterMap(conf.NodeID, conf.Cluster)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to init cluster info: %w", err)
	}
	log.Debugf("cluster info initialized: %v", clusterMap.ClusterInfo())

	rpcSvr, err := rpc.Serve(conf.RPCPort, conf.RESTPort, db, clusterMap)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to serve RPC: %w", err)
	}

	return &Service{
		conf:       conf,
		db:         db,
		clusterMap: clusterMap,
		rpcSvr:     rpcSvr,
	}, nil
}

func (s Service) Close() {
	log.Info("closing RPC...")
	s.rpcSvr.GracefulStop()

	log.Info("closing DB...")
	s.db.Close()
}
