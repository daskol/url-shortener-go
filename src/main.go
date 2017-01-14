package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
    "sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

var chars = []rune(letters + digits)

var configPath = flag.String("config", "", "Path to *.toml config.")

type Config struct {
	Host string
	Port int

	ExpiringTime time.Duration `toml:"expiring_time"`
	HostName     string        `toml:"host_name"`
	UriLength    int           `toml:"uri_length"`
}

type Uri string

type Url string

type UrlDesc struct {
	url       Url
	expiresAt time.Time
}

type UrlStorage struct {
	urls         map[Uri]UrlDesc
	expiringTime time.Duration
    mutex        sync.Mutex
}

var config Config

var urlStorage UrlStorage

func NewUri(length int) Uri {
	b := make([]rune, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return Uri("/" + string(b))
}

func NewUrlStorage(expiringTime time.Duration) UrlStorage {
    return UrlStorage{
        urls: make(map[Uri]UrlDesc),
        expiringTime: expiringTime,
	}
}

func (u *UrlStorage) Put(url Url, exp time.Duration) Uri {
    desc := UrlDesc{
        url: url,
        expiresAt: time.Now().Add(exp),
    }

    u.mutex.Lock()
    defer u.mutex.Unlock()

    for {
        uri := NewUri(config.UriLength)

        if val, ok := u.urls[uri]; !ok {
            u.urls[uri] = desc
            return uri
        } else if ok && val.url == url {
            u.urls[uri] = desc
            return uri
        }
    }
}

func (u *UrlStorage) Contains(uri Uri) bool {
	_, ok := u.urls[uri]
	return ok
}

func (u *UrlStorage) Get(uri Uri) *UrlDesc {
    desc, ok := u.urls[uri]

    if !ok {
        return nil
    } else if ok && desc.expiresAt.Before(time.Now()) {
        delete(u.urls, uri)
        return nil
    } else {
        return &desc
    }
}

func HandleShortRequest(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No url to shorten."))
		return
	}

    uri := urlStorage.Put(Url(url), config.ExpiringTime)

	location := config.HostName + string(uri) + "\n"

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(location))
}

func HandleExpandRequest(w http.ResponseWriter, r *http.Request) {
	if desc := urlStorage.Get(Uri(r.URL.Path)); desc != nil {
        w.Header().Set("Location", string(desc.url))
        w.WriteHeader(http.StatusFound)
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

func ReadConfig(path string) Config {
    config := Config{
        Host:         "localhost",
        Port:         8080,
        HostName:     "http://localhost",
        ExpiringTime: 3600 * time.Second,
        UriLength:    8,
    }

	if len(path) == 0 {
		return config
	}

	if _, err := os.Stat(path); err != nil {
		log.Fatal("Config file is missing: ", path)
	}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
	}

	config.ExpiringTime *= time.Second

	return config
}

func main() {
	flag.Parse()

	log.Println("    Simple URL shortner in Go")
	log.Println("    Daniel Bershatsky <daniel.bershatsky@skolkovotech.ru>, 2017")

	config = ReadConfig(*configPath)
	urlStorage = NewUrlStorage(config.ExpiringTime)

	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten/", HandleShortRequest)
	mux.HandleFunc("/", HandleExpandRequest)

	server := &http.Server{
		Addr:    config.Host + ":" + strconv.Itoa(config.Port),
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
