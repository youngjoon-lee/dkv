package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/cluster"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
	"github.com/youngjoon-lee/dkv/state"
	"github.com/youngjoon-lee/dkv/wal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	GracefulStop()
}

func Serve(rpcPort, restPort int, wal wal.WAL, state *state.State, cluster *cluster.Cluster) (Server, error) {
	grpcSvr, err := serveGRPC(rpcPort, wal, state, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	if err := serveREST(restPort, rpcPort); err != nil {
		grpcSvr.GracefulStop()
		return nil, fmt.Errorf("failed to serve REST server: %w", err)
	}

	return grpcSvr, nil
}

func serveGRPC(port int, wal wal.WAL, state *state.State, cluster *cluster.Cluster) (*grpc.Server, error) {
	svr := grpc.NewServer()
	pb.RegisterKVStoreServer(svr, &kvStoreServer{wal: wal, state: state, cluster: cluster})
	pb.RegisterReplicationServer(svr, &replicationServer{wal: wal, state: state, cluster: cluster})
	pb.RegisterClusterServiceServer(svr, &clusterServiceServer{cluster: cluster})

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
	rpcEndpoint := fmt.Sprintf("localhost:%d", rpcPort)

	if err := pb.RegisterKVStoreHandlerFromEndpoint(context.Background(), mux, rpcEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register grpc-gateway for kvstore: %w", err)
	}

	if err := pb.RegisterClusterServiceHandlerFromEndpoint(context.Background(), mux, rpcEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register grpc-gateway for cluster: %w", err)
	}

	log.Infof("REST server listening at %d...", port)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
			log.Panicf("REST server shutted down: %v", err)
		}
	}()

	return nil
}
