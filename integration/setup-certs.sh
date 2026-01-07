#!/bin/bash
set -e

cd /certs

if [ -f unbound_ca_ec.pem ]; then
  echo 'Certificates already exist'
  exit 0
fi

# Generate CA key and certificate
openssl ecparam -genkey -name secp384r1 -out unbound_ca_ec.key
openssl req -new -x509 -days 397 -key unbound_ca_ec.key -out unbound_ca_ec.pem \
  -subj "/CN=unbound-ca" \
  -addext "subjectKeyIdentifier=hash" \
  -addext "authorityKeyIdentifier=keyid:always,issuer:always" \
  -addext "basicConstraints=critical,CA:true"

# Generate server key and certificate
openssl ecparam -genkey -name secp384r1 -out unbound_server_ec.key
openssl req -new -key unbound_server_ec.key -out server.csr -subj "/CN=unbound"
openssl x509 -req -in server.csr -CA unbound_ca_ec.pem -CAkey unbound_ca_ec.key \
  -CAcreateserial -out unbound_server_ec.pem -days 397 -sha256 \
  -addext "subjectAltName=DNS:unbound,DNS:localhost" \
  -addext "subjectKeyIdentifier=hash" \
  -addext "authorityKeyIdentifier=keyid:always,issuer:always" \
  -addext "extendedKeyUsage=serverAuth"

# Append CA cert to server cert for chain
cat unbound_ca_ec.pem >> unbound_server_ec.pem

# Generate client key and certificate
openssl ecparam -genkey -name secp384r1 -out unbound_control_ec.key
openssl req -new -key unbound_control_ec.key -out client.csr -subj "/CN=unbound-control"
openssl x509 -req -in client.csr -CA unbound_ca_ec.pem -CAkey unbound_ca_ec.key \
  -CAcreateserial -out unbound_control_ec.pem -days 397 -sha256 \
  -addext "extendedKeyUsage=clientAuth"

# Set permissions
chmod 0640 *.key
chmod 0644 *.pem

# Cleanup
rm -f *.csr

echo 'Certificates generated successfully'

