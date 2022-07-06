# Prometheus Unbound exporter

This repository provides code for a simple Prometheus metrics exporter
for [the Unbound DNS resolver](https://unbound.net/). This exporter
connects to Unbounds TLS control socket and sends the `stats_noreset`
command, causing Unbound to return metrics as key-value pairs. The
metrics exporter converts Unbound metric names to Prometheus metric
names and labels by using a set of regular expressions.

- - - -

# Prerequisites

Go 1.16 or above is required.

# Installation

    go install github.com/letsencrypt/unbound_exporter@latest

This will install the binary in `$GOBIN`, or `$HOME/go/bin` if
`$GOBIN` is unset.

# Updating dependencies

```
go get -u
go mod tidy
```

- - - -

# Usage

To show all CLI flags available

    $ unbound_exporter -h

Scrape metrics from the exporter

    $ curl 127.0.0.1:9167/metrics | grep '^unbound_up'
    unbound_up 1

From the Unbound [statistics doc](https://www.nlnetlabs.nl/documentation/unbound/howto-statistics/): Unbound has an option to enable extended statistics collection. If enabled, more statistics are collected, for example what types of queries are sent to the resolver. Otherwise, only the total number of queries is collected. Add the following to your `unbound.conf`.

    server:
	    extended-statistics: yes
