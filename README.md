# url-shortener-go

*Simple URL shortener written in Go*

## Overview

[*Читайте на русском здесь.*](README.ru.md)

URL Shortner provides two kinds of URL storages. The first one wraps `map` data structure native in GO and makes it thread-sage. Unfortunately, the current implementation guarantee reliable storing. So in case of persistent storage is main requirement the second one was introduced which is based on [Bolt DB](https://github.com/boltdb/bolt). Waranties and properties of shortener implementation are pointed the following out.

1. Thread-safe for new URL shornening.
2. Time to live and base URL are parameters.
3. Configuration could be specified in [toml-file](etc/url-shortener.toml).
4. Two kinds of URL storages.
5. Guarantee persistence of short URLs with storage based on Bolt DB.
6. Work correctly behind proxy or balancer.

In order to start `url-shortener-go` by oneself with persistency storing of URLS one could just run
```bash
    ./url-shortener --url-storage bolt
```
Or one could register it as a system service defined with [systemd unit file](etc/url-shortener-go.service) and could run it.

[Try it here.](https://daskol.xyz/shorten/)

### API Methods

One could create new short URL with simple request to `/shorten/` URI parametrized with `url` which is target link.
```bash
    curl -v -X POST http://localhost:8080/shorten/?url=https://google.com
```
Server replies with 201(Created) response that contains short URL in `Location` header and duplicates it in response body.

Performing request to short URL created before server replies with 302(Found/Moved Temporary) response and sets `Location` header referred to original URL.
```bash
    curl -v http://localhost:8080/ri0xJwQ6
```
