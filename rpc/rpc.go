package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/db"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	GracefulStop()
}

func Serve(rpcPort, restPort int, db db.DB) (Server, error) {
	grpcSvr, err := serveGRPC(rpcPort, db)
	if err != nil {
		return nil, fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	if err := serveREST(restPort, rpcPort); err != nil {
		grpcSvr.GracefulStop()
		return nil, fmt.Errorf("failed to serve REST server: %w", err)
	}

	return grpcSvr, nil
}

func serveGRPC(port int, db db.DB) (*grpc.Server, error) {
	svr := grpc.NewServer()
	pb.RegisterKVStoreServer(svr, &kvStoreServer{db: db})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen port for RPC: %w", err)
	}

	go func() {
		log.Infof("gRPC server listening at %d...", port)
		if err := svr.Serve(lis); err != nil {
			log.Panicf("gRPC server shutted down: %v", err)
		}
	}()

	return svr, nil
}

func serveREST(port, rpcPort int) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterKVStoreHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%d", rpcPort), opts)
	if err != nil {
		return fmt.Errorf("failed to register grpc-gateway: %w", err)
	}

	log.Infof("REST server listening at %d...", port)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
			log.Panicf("REST server shutted down: %v", err)
		}
	}()

	return nil
}

type kvStoreServer struct {
	pb.UnimplementedKVStoreServer

	db db.DB
}

func (s *kvStoreServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutReply, error) {
	log.Infof("key:%v(%s), value:%v(%s)", req.Key, string(req.Key), req.Value, string(req.Value))
	if err := s.db.Put(req.Key, req.Value); err != nil {
		return nil, fmt.Errorf("failed to put: %w", err)
	}
	return &pb.PutReply{Message: "success"}, nil
}
