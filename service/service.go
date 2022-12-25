package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/db"
	"github.com/youngjoon-lee/dkv/rpc"
)

type Service struct {
	db     db.DB
	rpcSvr rpc.Server
}

//TODO: config
func New() (*Service, error) {
	db, err := db.NewBoltDB("my.db")
	if err != nil {
		return nil, fmt.Errorf("failed to init DB: %w", err)
	}

	rpcSvr, err := rpc.Serve(8080, 8081, db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to serve RPC: %w", err)
	}

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


