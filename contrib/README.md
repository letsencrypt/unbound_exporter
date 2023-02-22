# Contrib
This collection of scripts and files helps us further configure our unbounds and unbound_exporters.

## unbound-control-setup.sh

From [Golang 1.15 docs:](https://golang.google.cn/doc/go1.15#commonname)
> X.509 CommonName deprecation
> The deprecated, legacy behavior of treating the CommonName field on X.509 certificates as a host name when no Subject Alternative Names are present is now disabled by default. It can be temporarily re-enabled by adding the value x509ignoreCN=0 to the GODEBUG environment variable.
> Note that if the CommonName is an invalid host name, it's always ignored, regardless of GODEBUG settings. Invalid names include those with any characters other than letters, digits, hyphens and underscores, and those with empty labels or trailing dots.

Unbound still ships with an `unbound-control-setup` that generates a problematic keypair. This script will generate a keypair that satisfies newer versions of Golang.

Generate the new keypair
```
$ bash unbound-control-setup.sh
```

You'll then want to configure `/etc/unbound/unbound.conf` with the following stanza

```
$ cat /etc/unbound/unbound.conf
...
remote-control:
    control-enable: yes
    control-use-cert: yes
    server-key-file: "/etc/unbound/unbound_server_ec.key"
    server-cert-file: "/etc/unbound/unbound_server_ec.pem"
    control-key-file: "/etc/unbound/unbound_control_ec.key"
    control-cert-file: "/etc/unbound/unbound_control_ec.pem"
```

Test that you can still communicate with unbound via `unbound_control`. You should be able to see metrics.
```
$ unbound-control stats_noreset
thread0.num.queries=35
thread0.num.queries_ip_ratelimited=0
thread0.num.cachehits=25
thread0.num.cachemiss=10
thread0.num.prefetch=0
thread0.num.expired=0
...

```

To reconfigure `unbound_exporter` as a systemd service, see [this file](unbound_exporter.service).
