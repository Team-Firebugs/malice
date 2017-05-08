Malice API client
=================

This package provides the `api` package which attempts to
provide programmatic access to the full Malice API.

Currently, all of the Malice APIs included in version 0.3.0 are supported.

Documentation
=============

The full documentation is available on [Godoc](https://godoc.org/github.com/maliceio/malice/api)

Usage
=====

Below is an example of using the Malice client:

```go
// Get a new client
client, err := api.NewClient(api.DefaultConfig())
if err != nil {
    panic(err)
}

// Get a handle to the KV API
kv := client.KV()

```