package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/config"
	"github.com/youngjoon-lee/dkv/db"
	"github.com/youngjoon-lee/dkv/rpc"
)

type Service struct {
	db     db.DB
	rpcSvr rpc.Server
}

func New(conf config.Config) (*Service, error) {
	db, err := db.NewBoltDB(conf.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init DB: %w", err)
	}

	rpcSvr, err := rpc.Serve(conf.RPCPort, conf.RESTPort, db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to serve RPC: %w", err)
	}

	//TODO: handle leader_addr
	log.Infof("leader_addr: %v", conf.LeaderAddr)

	return &Service{
		db:     db,
		rpcSvr: rpcSvr,
	}, nil
}

func (s Service) Close() {
	log.Info("closing RPC...")
	s.rpcSvr.GracefulStop()

	log.Info("closing DB...")
	s.db.Close()
}
