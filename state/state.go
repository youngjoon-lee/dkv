package state

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/state/db"
	"github.com/youngjoon-lee/dkv/wal"
)

type State struct {
	mu            sync.RWMutex
	db            db.DB
	lastCommitted wal.Sequence
}

func New(dbPath string) (*State, error) {
	db, err := db.NewBoltDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	return &State{
		db:            db,
		lastCommitted: wal.NilSequence,
	}, nil
}

func (s *State) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Info("closing DB...")
	s.db.Close()
}

func (s *State) Commit(iter wal.Iterator) (wal.Sequence, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		elem := iter.Next()
		if elem == nil {
			return s.lastCommitted, nil
		}

		if err := s.db.Put(elem.Msg.Key, elem.Msg.Value); err != nil {
			err = fmt.Errorf("failed to put KV to DB: %w", err)
			if elem.DoneCh != nil {
				elem.DoneCh <- err
			}
			return s.lastCommitted, err
		}

		s.lastCommitted = elem.Seq
		if elem.DoneCh != nil {
			elem.DoneCh <- nil
		}
	}
}

func (s *State) LastCommitted() wal.Sequence {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastCommitted
}

func (s *State) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.Get(key)
}
