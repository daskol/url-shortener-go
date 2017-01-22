package core

import (
    "github.com/boltdb/bolt"
    "time"
)

type BoltDbStorage struct {
    urls         *bolt.DB
	expiringTime time.Duration
    uriLength    int
}

func NewBoltDbStorage(ttl time.Duration, length int) (*BoltDbStorage, error) {
    db, err := bolt.Open(
        "url-storage-bolt.db",
        0600,
        &bolt.Options{Timeout: 1 * time.Second})

    if err != nil {
        return nil, err
    }

    storage := BoltDbStorage{
        expiringTime: ttl,
        uriLength:    length,
        urls:         db,
    }

    return &storage, nil
}

func (b *BoltDbStorage) Put(url Url, exp time.Duration) Uri {
    return NewUri(b.uriLength)
}

func (b *BoltDbStorage) Contains(uri Uri) bool {
    return false
}

func (b *BoltDbStorage) Get(uri Uri) (Url, bool) {
    return Url(""), false
}
