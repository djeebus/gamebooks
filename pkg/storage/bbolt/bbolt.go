package bbolt

import (
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
	"go.etcd.io/bbolt"
	"strings"
)

var bucketName = []byte("game")

type Storage struct {
	db *bbolt.DB
}

func ensureBucket(db *bbolt.DB) error {
	if err := db.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket != nil {
			return nil
		}
		_, err := tx.CreateBucket(bucketName)
		if err != nil {
			return errors.Wrap(err, "failed to create bucket")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to create bucket")
	}

	return nil
}

func (b *Storage) Get(key string) (interface{}, error) {
	var data []byte

	if err := b.db.Batch(func(tx *bbolt.Tx) error {
		data = tx.Bucket(bucketName).Get([]byte(key))

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to get key")
	}

	if len(data) == 0 {
		return nil, nil
	}

	var result interface{}
	if err := msgpack.Unmarshal(data, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal data")
	}

	return result, nil
}

func (b *Storage) Set(key string, value interface{}) error {
	data, err := msgpack.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	if err = b.db.Batch(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), data)
	}); err != nil {
		return errors.Wrap(err, "failed to set data")
	}

	return nil
}

func (b *Storage) Remove(key string) error {
	return b.db.Batch(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketName).Delete([]byte(key))
	})
}

func (b *Storage) Clear(keyPrefix string) error {
	return b.db.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if err := bucket.ForEach(func(k, _ []byte) error {
			if key := string(k); strings.HasPrefix(key, keyPrefix) {
				if err := bucket.Delete(k); err != nil {
					return errors.Wrap(err, "failed to delete key")
				}
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "failed to delete keys")
		}

		return nil
	})
}

var _ Storage = new(Storage)

func New(path string) (*Storage, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	if err = ensureBucket(db); err != nil {
		return nil, errors.Wrap(err, "failed to create bucket")
	}

	return &Storage{db}, nil
}
