package core

import (
	"math/rand"
    "time"
)

type Uri string

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

var chars = []rune(letters + digits)

func NewUri(length int) Uri {
	b := make([]rune, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return Uri("/" + string(b))
}

type Url string

type UrlDesc struct {
	Url       Url
	ExpiresAt time.Time
}

type UrlStorage interface {
    Put(url Url, exp time.Duration) Uri;
    Get(uri Uri) (Url, bool);
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
