package replication

import (
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
			iter := r.wal.Iterate(r.state.LastCommitted() + 1)
			if iter == nil { // nothing to iterate
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// TODO: replicate logs to followers

			lastCommitted, err := r.state.Commit(iter)
			if err != nil {
				log.Errorf("failed to commit: %v", err)
			}
			log.Debugf("committed to seq %v", lastCommitted)
		}
	}
}

func (r *Replicator) Stop() {
	r.stopCh <- struct{}{}
	<-r.stopCh
}
