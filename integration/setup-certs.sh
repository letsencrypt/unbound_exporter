#!/bin/bash
set -e

# Generate CA key and certificate
openssl ecparam -genkey -name secp384r1 -out unbound_ca.key

openssl req -new -x509 -days 90 -key unbound_ca.key -out unbound_ca.pem \
  -subj "/CN=unbound-ca" \
  -addext "basicConstraints=critical,CA:true"

# Generate server key and certificate
openssl ecparam -genkey -name secp384r1 -out unbound_server.key

openssl req -new -key unbound_server.key -out server.csr -subj "/" \
  -addext "subjectAltName = DNS:unbound,DNS:localhost" \
  -addext "extendedKeyUsage = serverAuth"

openssl x509 -req -in server.csr -CA unbound_ca.pem -CAkey unbound_ca.key \
  -CAcreateserial -out unbound_server.pem -days 90 \

# Generate client key and certificate
openssl ecparam -genkey -name secp384r1 -out unbound_control.key

openssl req -new -key unbound_control.key -out client.csr -subj "/CN=unbound-control"

openssl x509 -req -in client.csr -CA unbound_ca.pem -CAkey unbound_ca.key \
  -CAcreateserial -out unbound_control.pem -days 90 \

# Set permissions
chmod 0640 *.key
chmod 0644 *.pem

# Cleanup
rm -f *.csr

echo 'Certificates generated successfully'

