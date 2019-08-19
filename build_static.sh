#!/bin/sh

docker run -i -v `pwd`:/unbound_exporter alpine:edge /bin/sh << 'EOF'
set -ex

# Install prerequisites for the build process.
apk update
apk add ca-certificates git go libc-dev
update-ca-certificates

# Build the unbound_exporter.
cd /unbound_exporter
go build --ldflags '-extldflags "-static"'
strip unbound_exporter
EOF
