# Prometheus Unbound exporter

This repository provides code for a simple Prometheus metrics exporter
for [the Unbound DNS resolver](https://unbound.net/). This exporter
connects to Unbounds TLS control socket and sends the `stats_noreset`
command, causing Unbound to return metrics as key-value pairs. The
metrics exporter converts Unbound metric names to Prometheus metric
names and labels by using a set of regular expressions.

- - - -

# Installation

To install this code and in your go environment. You can then add the binary to your `PATH`.

    go get github.com/kumina/unbound_exporter
    go install github.com/kumina/unbound_exporter

- - - -

# Usage

To show all CLI flags available

    unbound_exporter -h

For extended statistics, you may want to add the following to your unbound.conf

	# statistics
	extended-statistics: yes
