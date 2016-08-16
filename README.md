# RedFace

[![Build Status](https://travis-ci.org/Akagi201/redface.svg)](https://travis-ci.org/Akagi201/redface) [![Coverage Status](https://coveralls.io/repos/github/Akagi201/redface/badge.svg?branch=master)](https://coveralls.io/github/Akagi201/redface?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/Akagi201/redface)](https://goreportcard.com/report/github.com/Akagi201/redface) [![GoDoc](https://godoc.org/github.com/Akagi201/redface?status.svg)](https://godoc.org/github.com/Akagi201/redface)

RedFace means redis interface.

It can be used as a redis server-side api in golang.

## Features

- [x] Suport tcp protocol.
- [x] Support unix socket protocol.
- [x] Support net/http like interface.
- [ ] Support TLS.
- [ ] Support net/context.

## Install

* `go get github.com/Akagi201/redface`

## Import

* `import "github.com/Akagi201/redface/resp"`
* `import "github.com/Akagi201/redface/server"`
