package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"github.com/youngjoon-lee/dkv/server"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Info("starting Distributed Key-Value Store...")

	db, err := bbolt.Open("my.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Info("closing BoltDB...")
		db.Close()
	}()

	grpcPort := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal(err)
	}

	grpcSvr := grpc.NewServer()
	pb.RegisterGreeterServer(grpcSvr, &server.GRPCServer{})
	log.Infof("gRPC server listening at %d...", grpcPort)
	go func() {
		if err := grpcSvr.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterGreeterHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		log.Fatal(err)
	}

	httpPort := grpcPort + 1
	log.Infof("HTTP server listening at %d...", httpPort)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux); err != nil {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-sigCh:
		log.Info("signal detected. starting graceful shutdown...")
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//TODO: do something with ctx
}
