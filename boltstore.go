package boltstore

import (
	"encoding/json"
	"fmt"
	"sync"

	bolt "go.etcd.io/bbolt"
)

// NoSuchKeyError is thrown when calling Get with invalid key
type NoSuchKeyError struct {
	key string
}

func (err NoSuchKeyError) Error() string {
	return "BoltStore: no such key \"" + err.key + "\""
}

// BoltStore is the basic store object.
type BoltStore struct {
	bucket string
	db     *bolt.DB
	sync.RWMutex
}

// Open will load a BoltStore from a file.
func Open(filename string) (s *BoltStore, err error) {
	s = new(BoltStore)
	s.Lock()
	defer s.Unlock()
	s.db, err = bolt.Open(filename, 0600, nil)
	if err != nil {
		return
	}
	s.bucket = "BoltStore"
	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return
}

// Set saves a value at the given key.
func (s *BoltStore) Set(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()
	bValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		err := b.Put([]byte(key), bValue)
		return err
	})
}

// Get will return the value associated with a key.
func (s *BoltStore) Get(key string, v interface{}) (err error) {
	s.RLock()
	defer s.RUnlock()
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		val := b.Get([]byte(key))
		if val == nil {
			return NoSuchKeyError{key}
		}
		return json.Unmarshal(val, &v)
	})
	return
}

// Keys returns all the keys currently in map
func (s *BoltStore) Keys() []string {
	s.RLock()
	defer s.RUnlock()
	numKeys := 0

	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(s.bucket))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			numKeys++
		}
		return nil
	})

	keys := make([]string, numKeys)
	i := 0
	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(s.bucket))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys[i] = string(k)
			i++
		}
		return nil
	})

	return keys
}

// Delete removes a key from the store.
func (s *BoltStore) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		return b.Delete([]byte(key))
	})
}
