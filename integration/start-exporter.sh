#!/bin/sh
set -e

if [ -n "$EXPORTER_COMMAND" ]; then
  exec /unbound_exporter $EXPORTER_COMMAND
else
  exec /unbound_exporter -unbound.host=unix:///var/run/socket/unbound.ctl
fi
