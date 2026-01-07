# Integration Tests

This directory contains integration tests for unbound_exporter that verify it can connect to and scrape metrics from Unbound using different connection types.

## Test Coverage

The integration tests verify three connection modes:

1. **Unix Socket** (`TestIntegration`) - Tests connection via Unix domain socket
2. **TCP** (`TestIntegrationTCP`) - Tests plain TCP connection without TLS
3. **TLS** (`TestIntegrationTLS`) - Tests TLS-encrypted TCP connection with mutual authentication

## Running Tests

To run the integration tests:

```bash
# Start all services (Unix socket, TCP, and TLS instances)
docker compose up --build -d

# Run integration tests
go test --tags=integration -v

# Clean up
docker compose down
```

## Architecture

The docker-compose.yml defines six services:

- `unbound` - Unbound instance with Unix socket control interface
- `unbound_exporter` - Exporter connecting to unbound via Unix socket (port 9167)
- `unbound_tcp` - Unbound instance with TCP control interface (no TLS)
- `unbound_exporter_tcp` - Exporter connecting via TCP (port 9168)
- `unbound_tls` - Unbound instance with TLS control interface
- `unbound_exporter_tls` - Exporter connecting via TLS (port 9169)

## TLS Certificates

TLS certificates are automatically generated during the Docker build process using the `configs/generate-certs.sh` script. This script creates:

- CA certificate and key (`unbound_ca_ec.pem`, `unbound_ca_ec.key`)
- Server certificate and key (`unbound_server_ec.pem`, `unbound_server_ec.key`)
- Client certificate and key (`unbound_control_ec.pem`, `unbound_control_ec.key`)

These certificates are only used for testing and are regenerated on each build.

## Configuration Files

- `configs/unbound-example.conf` - Unix socket configuration (original)
- `configs/unbound-tcp.conf` - TCP configuration without TLS
- `configs/unbound-tls.conf` - TLS configuration with mutual authentication
- `configs/generate-certs.sh` - Script to generate TLS certificates
