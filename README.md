# url-shortner-go

*Simple URL shortner written in Go*

## Overview

[*Читайте на русском здесь.*](README.md.ru)

### API Methods

One could create new short URL with simple request to `/shorten/` URI parametrized with `url` which is target link.
```bash
    curl -v http://localhost:8080/shorten/?url=https://google.com
```
Server replies with 201(Created) response that contains short URL in `Location` header and duplicates it in response body.

Performing request to short URL created before server replies with 302(Found/Moved Temporary) response and sets `Location` header referred to original URL.
```bash
    curl -v http://localhost:8080/ri0xJwQ6
```
