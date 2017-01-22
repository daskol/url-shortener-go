package core

import (
    "bytes"
    "encoding/gob"
    "errors"
    "github.com/boltdb/bolt"
    "log"
    "time"
)

type BoltStorage struct {
    urls         *bolt.DB
	expiringTime time.Duration
    uriLength    int
}

func NewBoltStorage(ttl time.Duration, length int, database string) (*BoltStorage, error) {
    db, err := bolt.Open(
        database,
        0600,
        &bolt.Options{Timeout: 1 * time.Second})

    if err != nil {
        return nil, err
    }

    storage := BoltStorage{
        expiringTime: ttl,
        uriLength:    length,
        urls:         db,
    }

    if err := storage.urls.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("urls"))
        return err
    }); err != nil {
        return nil, err
    }

    return &storage, nil
}

func (b *BoltStorage) Put(url Url, exp time.Duration) Uri {
    uri := Uri("")
    desc := UrlDesc{
        Url: url,
        ExpiresAt: time.Now().Add(exp),
    }

    if err := b.urls.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("urls"))

        if bucket == nil {
            return errors.New("no bucket `urls`")
        }

        for {
            uri = NewUri(b.uriLength)

            if bucket.Get([]byte(uri)) != nil {
                continue
            } else {
                var buffer bytes.Buffer
                encoder := gob.NewEncoder(&buffer)
                encoder.Encode(desc)
                return bucket.Put([]byte(uri), buffer.Bytes())
            }
        }
    }); err != nil {
        return uri
    } else {
        return uri
    }
}

func (b *BoltStorage) Contains(uri Uri) bool {
    return false
}

func (b *BoltStorage) Get(uri Uri) (Url, bool) {
    var desc UrlDesc

    if err := b.urls.View(func(tx *bolt.Tx) error {
        if bucket := tx.Bucket([]byte("urls")); bucket == nil {
            return errors.New("no bucket `urls`")
        } else if value := bucket.Get([]byte(uri)); value != nil {
            buffer := bytes.NewBuffer(value)
            decoder := gob.NewDecoder(buffer)
            decoder.Decode(&desc)
            return nil
        } else {
            return errors.New("key not in bucket: " + string(uri))
        }
    }); err != nil {
        return Url(""), false
    }

    if b.expiringTime <= 0 || desc.ExpiresAt.After(time.Now()) {
        return Url(desc.Url), true
    }

    if err := b.urls.Update(func(tx *bolt.Tx) error {
        if bucket := tx.Bucket([]byte("urls")); bucket == nil {
            return errors.New("no bucket `urls`")
        } else {
            return bucket.Delete([]byte(uri))
        }
    }); err != nil {
        log.Println(err)
    }

    return Url(desc.Url), false
}
