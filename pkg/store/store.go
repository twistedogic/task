package store

type Item struct {
	Key   []byte
	Value []byte
}

type Store interface {
	Move(src, dst, key []byte) error
	Get(table, key []byte) ([]byte, error)
	Set(table, key, value []byte) error
	Delete(table, key []byte) error
	List([]byte) ([]Item, error)
}
