[![Stories in Ready](https://badge.waffle.io/Akagi201/redface.png?label=ready&title=Ready)](https://waffle.io/Akagi201/redface)
# RedFace

[![Build Status](https://travis-ci.org/Akagi201/redface.svg)](https://travis-ci.org/Akagi201/redface) [![Coverage Status](https://coveralls.io/repos/github/Akagi201/redface/badge.svg?branch=master)](https://coveralls.io/github/Akagi201/redface?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/Akagi201/redface)](https://goreportcard.com/report/github.com/Akagi201/redface) [![GoDoc](https://godoc.org/github.com/Akagi201/redface?status.svg)](https://godoc.org/github.com/Akagi201/redface)

RedFace means redis interface.

It can be used as a redis server-side api in golang.

## Features

- [x] Suport tcp protocol.
- [x] Support unix socket protocol.
- [x] Support net/http like interface.
- [x] Add benchmarks.
- [ ] Support pipelining.
- [ ] Support telnet commands.
- [ ] Support redis lua script.
- [ ] Support TLS.
- [ ] Support net/context.

## Install

* `go get github.com/Akagi201/redface`

## Import

* `import "github.com/Akagi201/redface/resp"`
* `import "github.com/Akagi201/redface/server"`

## Benchmarks

### redis-benchmark

Redis: Single-threaded, no disk persistence.

```
❯ redis-server --port 6379  --appendonly no
```

```
❯ redis-benchmark -p 6379 -t set,get -n 1000000 -q -P 512 -c 512
SET: 767459.75 requests per second
GET: 941619.56 requests per second
```

RedFace: Single-threaded, no disk persistence.

```
GOMAXPROCS=1 go run example/clone/main.go
```

```
❯ redis-benchmark -p 6389 -t set,get -n 1000000 -q -P 512 -c 512
SET: 68861.04 requests per second
GET: 65261.37 requests per second
```

RedFace: Multi-threaded, no disk persistence.

```
GOMAXPROCS=0 go run example/clone/main.go
```

```
❯ redis-benchmark -p 6389 -t set,get -n 1000000 -q -P 512 -c 512
SET: 30049.88 requests per second
GET: 30422.88 requests per second
```

Hardward info

```
❯ system_profiler SPHardwareDataType
Hardware:

    Hardware Overview:

      Model Name: MacBook Pro
      Model Identifier: MacBookPro11,3
      Processor Name: Intel Core i7
      Processor Speed: 2.3 GHz
      Number of Processors: 1
      Total Number of Cores: 4
      L2 Cache (per Core): 256 KB
      L3 Cache: 6 MB
      Memory: 16 GB
      Boot ROM Version: MBP112.0138.B17
      SMC Version (system): 2.19f12
      Serial Number (system): C02MG6L8FD57
      Hardware UUID: EB84A5CF-F1BA-5604-B1A6-534E30EA95C1
```
