package config

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type LogLevel log.Level
type Cluster []*pb.Node

type Config struct {
	LogLevel LogLevel `envconfig:"LOG_LEVEL" default:"info"`
	NodeID   string   `envconfig:"NODE_ID" required:"true"`
	RPCPort  int      `envconfig:"RPC_PORT" required:"true"`
	RESTPort int      `envconfig:"REST_PORT" required:"true"`
	DBPath   string   `envconfig:"DB_PATH" required:"true"`
	Cluster  Cluster  `envconfig:"CLUSTER" required:"true"`
}

func (l *LogLevel) Decode(value string) error {
	level, err := log.ParseLevel(value)
	if err != nil {
		return fmt.Errorf("failed to parse LogLevel: %w", err)
	}

	*l = LogLevel(level)
	return nil
}

func (c *Cluster) Decode(value string) error {
	cluster := make([]*pb.Node, 0)

	for _, nodeStr := range strings.Split(value, ",") {
		node, err := parseNode(nodeStr)
		if err != nil {
			return fmt.Errorf("invalid node str(%v): %w", nodeStr, err)
		}
		cluster = append(cluster, node)
	}

	*c = cluster
	return nil
}

func parseNode(str string) (*pb.Node, error) {
	idAndAddr := strings.Split(str, "@")
	if len(idAndAddr) != 2 {
		return nil, fmt.Errorf("failed to parse a node string")
	}

	return &pb.Node{Id: idAndAddr[0], Addr: idAndAddr[1]}, nil
}
