package core

import (
    "sync"
    "time"
)

type MapStorage struct {
	urls         map[Uri]UrlDesc
	expiringTime time.Duration
    mutex        sync.Mutex
    uriLength    int
}

func NewMapStorage(ttl time.Duration, length int) (*MapStorage, error) {
    return &MapStorage{
        urls: make(map[Uri]UrlDesc),
        uriLength: length,
        expiringTime: ttl,
	}, nil
}

func (u *MapStorage) Put(url Url, exp time.Duration) Uri {
    desc := UrlDesc{
        Url: url,
        ExpiresAt: time.Now().Add(exp),
    }

    u.mutex.Lock()
    defer u.mutex.Unlock()

    for {
        uri := NewUri(u.uriLength)

        if val, ok := u.urls[uri]; !ok {
            u.urls[uri] = desc
            return uri
        } else if ok && val.Url == url {
            u.urls[uri] = desc
            return uri
        }
    }
}

func (u *MapStorage) Contains(uri Uri) bool {
	_, ok := u.urls[uri]
	return ok
}

func (u *MapStorage) Get(uri Uri) (Url, bool) {
    u.mutex.Lock()
    defer u.mutex.Unlock()

    desc, ok := u.urls[uri]

    if !ok {
        return Url(""), false
    } else if ok && u.expiringTime > 0 && desc.ExpiresAt.Before(time.Now()) {
        delete(u.urls, uri)
        return desc.Url, false
    } else {
        return desc.Url, true
    }
}
