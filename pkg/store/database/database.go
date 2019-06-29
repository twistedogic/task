package database

import (
	"bytes"
	"fmt"
	"time"

	bolt "github.com/etcd-io/bbolt"
	"github.com/twistedogic/task/pkg/store"
)

type DB struct {
	*bolt.DB
}

func New(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Set(bucket, key, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("failed to create bucket '%s': %v", bucket, err)
		}
		if err := bk.Put(key, value); err != nil {
			return fmt.Errorf("failed to insert '%s': %v", key, err)
		}
		return nil
	})
}

func (db *DB) Get(bucket, key []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}
		b := bk.Get(key)
		if b == nil {
			return fmt.Errorf("key '%s' not found", key)
		}
		if _, err := buf.Write(b); err != nil {
			return err
		}
		return nil
	})
	return buf.Bytes(), err
}

func (db *DB) Delete(bucket, key []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}
		return bk.Delete(key)
	})
}

func (db *DB) List(bucket []byte) ([]store.Item, error) {
	output := []store.Item{}
	err := db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return nil
		}
		c := bk.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			output = append(output, store.Item{Key: k, Value: v})
		}
		return nil
	})
	return output, err
}

func (db *DB) Move(src, dst, key []byte) error {
	b, err := db.Get(src, key)
	if err != nil {
		return err
	}
	if err := db.Set(dst, key, b); err != nil {
		return err
	}
	return db.Delete(src, key)
}
