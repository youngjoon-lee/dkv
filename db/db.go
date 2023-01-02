package db

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type DB interface {
	Close()
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
}

type boltDB struct {
	db         *bbolt.DB
	bucketName []byte
}

func NewBoltDB(path string) (DB, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	bucketName := []byte("default")
	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bucket: %w", err)
	}

	return &boltDB{
		db:         db,
		bucketName: bucketName,
	}, nil
}

func (b boltDB) Close() {
	b.db.Close()
}

func (b boltDB) Put(key, value []byte) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucketName)
		return bucket.Put(key, value)
	})
	if err != nil {
		return fmt.Errorf("failed to put kv to boltDB: %w", err)
	}

	return nil
}

func (b boltDB) Get(key []byte) ([]byte, error) {
	var out []byte

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucketName)
		value := bucket.Get(key)
		if value == nil {
			out = nil
			return nil
		}

		out = make([]byte, len(value))
		copy(out, value)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put value from boltDB: %w", err)
	}

	return out, nil
}
