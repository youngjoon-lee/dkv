package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/server"
	"github.com/youngjoon-lee/dkv/types"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

func main() {
	log.Info("starting Distributed Key-Value Store...")

	db, err := bbolt.Open("my.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	grpcSvr := grpc.NewServer()
	types.RegisterGreeterServer(grpcSvr, &server.GRPCServer{})
	log.Infof("gRPC server listening at %d...", port)
	if err := grpcSvr.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
