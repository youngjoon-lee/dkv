package wal

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

var _ WAL = (*memWAL)(nil)

type memWAL struct {
	mu      sync.RWMutex
	logs    []*Element
	nextSeq Sequence
}

func NewMemWAL() WAL {
	return &memWAL{
		logs:    make([]*Element, 0),
		nextSeq: InitialSequence,
	}
}

func (w *memWAL) Append(msg *pb.PutRequest, doneCh chan<- error) (Sequence, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	elem := &Element{
		Seq:    w.nextSeq,
		Msg:    msg,
		DoneCh: doneCh,
	}
	w.logs = append(w.logs, elem)
	w.nextSeq++

	log.Debugf("wal: appended as seq %v", elem.Seq)
	return elem.Seq, nil
}

func (w *memWAL) AppendWith(seq Sequence, message *pb.PutRequest) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.nextSeq < seq {
		return fmt.Errorf("next seq must be %v", w.nextSeq)
	} else if w.nextSeq > seq { // already appended
		log.Debugf("wal: seq %v is already appended", seq)
		return nil
	}

	elem := &Element{
		Seq: seq,
		Msg: message,
	}
	w.logs = append(w.logs, elem)
	w.nextSeq++

	log.Debugf("wal: appended with seq %v", elem.Seq)
	return nil
}

func (w *memWAL) Iterate(from, to Sequence) Iterator {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i, elem := range w.logs {
		if elem.Seq >= from {
			return newMemWALIterator(w, i, to)
		}
	}
	return nil
}

func (w *memWAL) Latest() *Element {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.logs[len(w.logs)-1]
}

var _ Iterator = (*memWALIterator)(nil)

type memWALIterator struct {
	wal     *memWAL
	nextIdx int
	toSeq   Sequence
}

func newMemWALIterator(wal *memWAL, fromIdx int, toSeq Sequence) *memWALIterator {
	return &memWALIterator{
		wal:     wal,
		nextIdx: fromIdx,
		toSeq:   toSeq,
	}
}

func (m *memWALIterator) Next() *Element {
	m.wal.mu.RLock()
	defer m.wal.mu.RUnlock()

	//TODO: handle the case when logs are truncated
	if m.nextIdx < len(m.wal.logs) {
		elem := m.wal.logs[m.nextIdx]
		if m.toSeq != NilSequence && elem.Seq > m.toSeq {
			return nil
		}
		m.nextIdx++
		return elem
	}

	return nil
}
