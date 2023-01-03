package replication

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/state"
	"github.com/youngjoon-lee/dkv/wal"
)

type Replicator struct {
	stopCh chan struct{}

	wal   wal.WAL
	state *state.State
}

func NewReplicator(wal wal.WAL, state *state.State) *Replicator {
	r := &Replicator{
		stopCh: make(chan struct{}),
		wal:    wal,
		state:  state,
	}

	go r.start()

	return r
}

func (r *Replicator) start() {
	for {
		select {
		case <-r.stopCh:
			r.stopCh <- struct{}{}
		default:
			if err := r.runIteration(); err != nil {
				log.Error(err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (r *Replicator) Stop() {
	r.stopCh <- struct{}{}
	<-r.stopCh
}

func (r *Replicator) runIteration() error {
	// TODO: replicate logs to followers

	from := r.state.LastCommitted() + 1
	iter := r.wal.Iterate(from, wal.NilSequence)

	lastCommitted, err := r.state.Commit(iter)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	} else if from <= lastCommitted {
		log.Debugf("committed to seq %v~%v", from, lastCommitted)
	}

	return nil
}