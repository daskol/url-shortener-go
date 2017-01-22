package main

import (
	"core"
    "errors"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var configPath = flag.String("config", "", "Path to *.toml config.")

type Config struct {
	Host string
	Port int

	ExpiringTime time.Duration `toml:"expiring_time"`
	HostName     string        `toml:"host_name"`
	UriLength    int           `toml:"uri_length"`
    UrlStorage   string        `toml:"url_storage"`

    Bolt BoltStorageConfig     `toml:"bolt-storage"`
}

type BoltStorageConfig struct {
    Database string `toml:"database"`
}

var config Config

var urlStorage core.UrlStorage

func HandleShortRequest(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No url to shorten."))
		return
	}

	uri := urlStorage.Put(core.Url(url), config.ExpiringTime)

	location := config.HostName + string(uri) + "\n"

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(location))
}

func HandleExpandRequest(w http.ResponseWriter, r *http.Request) {
	if url, ok := urlStorage.Get(core.Uri(r.URL.Path)); ok {
		w.Header().Set("Location", string(url))
		w.WriteHeader(http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func ReadConfig(path string) Config {
    log.Printf("read config from %s\n", path)

	config := Config{
		Host:         "localhost",
		Port:         8080,
        HostName:     "http://localhost:8080",
		ExpiringTime: 3600 * time.Second,
		UriLength:    8,
        UrlStorage:   "map",
        Bolt:        BoltStorageConfig{
            Database: "url-storage-bolt.db",
        },
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

    if config.ExpiringTime <= 0 {
        log.Println("ttl of url is set to store urls forever")
    }

	return config
}

func NewUrlStorage(config *Config) (core.UrlStorage, error) {
    switch config.UrlStorage {
    case "map":
        return core.NewMapStorage(config.ExpiringTime, config.UriLength)
    case "bolt":
        return core.NewBoltStorage(
            config.ExpiringTime,
            config.UriLength,
            config.Bolt.Database,
        )
    default:
        return nil, errors.New("unknown storage " + config.UrlStorage)
    }
}

func main() {
	flag.Parse()

	log.Println("    Simple URL shortner in Go")
	log.Println("    Daniel Bershatsky <daniel.bershatsky@skolkovotech.ru>, 2017")

    config = ReadConfig(*configPath)
    storage, err := NewUrlStorage(&config)

	if err != nil {
		log.Fatal(err)
	} else {
        log.Println("use storage: " + config.UrlStorage)
		urlStorage = storage
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten/", HandleShortRequest)
	mux.HandleFunc("/", HandleExpandRequest)

	server := &http.Server{
		Addr:    config.Host + ":" + strconv.Itoa(config.Port),
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
