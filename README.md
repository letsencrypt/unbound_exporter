# Prometheus Unbound exporter

This repository provides code for a simple Prometheus metrics exporter
for [the Unbound DNS resolver](https://unbound.net/). This exporter
connects to Unbounds TLS control socket and sends the `stats_noreset`
command, causing Unbound to return metrics as key-value pairs. The
metrics exporter converts Unbound metric names to Prometheus metric
names and labels by using a set of regular expressions.

Please refer to this utility's `main()` function for a list of supported
command line flags.

- - - -

# Building

To build this code, you'll need a working `go` environment. Once you have a suitable environment, simply run the following commands

    go get
    go build unbound_exporter.go

If your build is working successfully, you should see the following output.

    ./unbound_exporter -h
    Usage of ./unbound_exporter:
      -unbound.ca string
        Unbound server certificate. (default "/etc/unbound/unbound_server.pem")
      -unbound.cert string
        Unbound client certificate. (default "/etc/unbound/unbound_control.pem")
      -unbound.host string
        Unbound control socket hostname and port number. (default "localhost:8953")
      -unbound.key string
        Unbound client key. (default "/etc/unbound/unbound_control.key")
      -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9167")
      -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
