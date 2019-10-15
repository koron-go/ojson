# koron-go/ojson

[![GoDoc](https://godoc.org/github.com/koron-go/ojson?status.svg)](https://godoc.org/github.com/koron-go/ojson)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/ojson/master.svg)](https://circleci.com/gh/koron-go/ojson/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/ojson)](https://goreportcard.com/report/github.com/koron-go/ojson)

**ojson** is a package to describe **o**rdered **JSON** object.

Ordered JSON

```go
// this gets `{"foo":123,"bar":"xyz"}`
json.Marshal(ojson.Object{
    {"foo", 123},
    {"bar", "xyz"},
})

// this gets `{"bar":"xyz","foo":123}`
json.Marshal(ojson.Object{
    {"bar", "xyz"},
    {"foo", 123},
})
```
