#!/bin/bash
set -e

# Start unbound with appropriate config using include directive
if [ "$UNBOUND_CONFIG" = "/etc/unbound/unbound-tcp.conf" ]; then
  # Create config that includes the TCP remote control config
  cp /etc/unbound/unbound-example.conf /tmp/unbound.conf
  echo "include: /etc/unbound/remote-control-tcp.conf" >> /tmp/unbound.conf
  exec /usr/local/sbin/unbound -v -d -c /tmp/unbound.conf
elif [ "$UNBOUND_CONFIG" = "/etc/unbound/unbound-tls.conf" ]; then
  # Create config that includes the TLS remote control config
  cp /etc/unbound/unbound-example.conf /tmp/unbound.conf
  echo "include: /etc/unbound/remote-control-tls.conf" >> /tmp/unbound.conf
  exec /usr/local/sbin/unbound -v -d -c /tmp/unbound.conf
else
  # Default: use existing config (unix socket)
  exec /usr/local/sbin/unbound -v -d -c /etc/unbound/unbound-example.conf
fi

