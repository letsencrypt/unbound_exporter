# Prometheus Unbound exporter

This repository provides code for a simple Prometheus metrics exporter
for [the Unbound DNS resolver](https://unbound.net/). This exporter
connects to Unbounds TLS control socket and sends the `stats_noreset`
command, causing Unbound to return metrics as key-value pairs. The
metrics exporter converts Unbound metric names to Prometheus metric
names and labels by using a set of regular expressions.

- - - -

# Installation

Install [Bazel](https://bazel.build/) on your system and run:

    $ bazel build //:unbound_exporter
    ...
    Target //:unbound_exporter up-to-date:
      bazel-bin/linux_amd64_pure_stripped/unbound_exporter
    ...
    $ install -m 555 bazel-bin/linux_amd64_pure_stripped/unbound_exporter /usr/bin/unbound_exporter

- - - -

# Usage

To show all CLI flags available

    unbound_exporter -h
