package db

type DB interface {
	Close()
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
}
