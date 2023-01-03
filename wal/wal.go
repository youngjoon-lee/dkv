package wal

import (
	pb "github.com/youngjoon-lee/dkv/pb/dkv/v0"
)

type WAL interface {
	Append(msg *pb.PutRequest, doneCh chan<- error) (Sequence, error)
	AppendWith(seq Sequence, msg *pb.PutRequest) error
	Iterate(fromSeq Sequence) Iterator
	Latest() *Element
}

type Sequence uint64

const (
	NilSequence     = Sequence(0)
	InitialSequence = Sequence(1)
)

type Iterator interface {
	Next() *Element
}

type Element struct {
	Seq    Sequence
	Msg    *pb.PutRequest
	DoneCh chan<- error
}
