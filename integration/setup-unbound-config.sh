#!/bin/bash
set -e

# Generate unbound config based on environment variable
if [ "$UNBOUND_CONFIG" = "/etc/unbound/unbound-tcp.conf" ]; then
  cp /etc/unbound/unbound-example.conf /tmp/unbound.conf
  sed -i 's|control-interface: /var/run/socket/unbound.ctl|control-interface: 0.0.0.0\n    control-port: 8953\n    control-use-cert: no|' /tmp/unbound.conf
  exec /usr/local/sbin/unbound -v -d -c /tmp/unbound.conf
elif [ "$UNBOUND_CONFIG" = "/etc/unbound/unbound-tls.conf" ]; then
  cp /etc/unbound/unbound-example.conf /tmp/unbound.conf
  sed -i 's|control-interface: /var/run/socket/unbound.ctl|control-interface: 0.0.0.0\n    control-port: 8954\n    control-use-cert: yes\n    server-key-file: "/certs/unbound_server_ec.key"\n    server-cert-file: "/certs/unbound_server_ec.pem"\n    control-key-file: "/certs/unbound_control_ec.key"\n    control-cert-file: "/certs/unbound_control_ec.pem"|' /tmp/unbound.conf
  exec /usr/local/sbin/unbound -v -d -c /tmp/unbound.conf
else
  exec /usr/local/sbin/unbound -v -d -c /etc/unbound/unbound-example.conf
fi
