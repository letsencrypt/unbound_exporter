#!/usr/bin/env bash

# Script to generate certificates for TLS testing
# Based on contrib/unbound-cert-setup.sh

set -e

# directory for files
DESTDIR="/etc/unbound/certs"

# validity period for certificates
DAYS=397

# hash algorithm
HASH=sha256

# base name for unbound CA keys
CA_BASE=unbound_ca_ec

# base name for unbound server keys
SVR_BASE=unbound_server_ec

# base name for unbound-control keys
CTL_BASE=unbound_control_ec

# we want -rw-r----- access (say you run this as root: grp=yes (server), all=no).
umask 0027

# functions:
error ( ) {
        echo "$0 fatal error: ${1}"
        exit 1
}

# go!:
echo "setup in directory ${DESTDIR}"
mkdir -p "${DESTDIR}"
cd "${DESTDIR}" || error "could not cd to ${DESTDIR}"

# create certificate keys; do not recreate if they already exist.
if test -f "${CA_BASE}.key"; then
        echo "${CA_BASE}.key exists"
else
        echo "generating ${CA_BASE}.key"
        openssl ecparam -genkey -name secp384r1 > ${CA_BASE}.key || error "could not gen ecdsa"
fi
if test -f "${SVR_BASE}.key"; then
        echo "${SVR_BASE}.key exists"
else
        echo "generating ${SVR_BASE}.key"
        openssl ecparam -genkey -name secp384r1 > ${SVR_BASE}.key || error "could not gen ecdsa"
fi
if test -f "${CTL_BASE}.key"; then
        echo "${CTL_BASE}.key exists"
else
        echo "generating ${CTL_BASE}.key"
        openssl ecparam -genkey -name secp384r1 > ${CTL_BASE}.key || error "could not gen ecdsa"
fi

# create self-signed cert CSR for server
cat > ca_request.cfg <<EOCAConfig
[req]
prompt                 = no
distinguished_name     = req_distinguished_name
x509_extensions        = req_v3_extensions

[req_distinguished_name]
commonName             = unbound-ca

[req_v3_extensions]
subjectKeyIdentifier   = hash
authorityKeyIdentifier = keyid:always,issuer:always
basicConstraints       = CA:true
EOCAConfig
test -f ca_request.cfg || error "could not create ca_request.cfg"

echo "creating ${CA_BASE}.pem (self signed certificate)"
openssl req -key "${CA_BASE}.key" -config ca_request.cfg  -new -x509 -days "${DAYS}" -out "${CA_BASE}.pem" || error "could not create ${CA_BASE}.pem"

# --------------

# create server cert CSR and sign it, piped
cat > server_request.cfg <<EOServerConfig
[req]
prompt                 = no
distinguished_name     = req_distinguished_name

[req_distinguished_name]
commonName             = unbound
EOServerConfig
test -f server_request.cfg || error "could not create server_request.cfg"

cat > server_exts.cfg <<EOServerConfig
subjectAltName         = @alt_names
subjectKeyIdentifier   = hash
authorityKeyIdentifier = keyid:always,issuer:always
extendedKeyUsage       = serverAuth

[alt_names]
DNS.1                  = unbound_tls
DNS.2                  = unbound
DNS.3                  = localhost
EOServerConfig
test -f server_exts.cfg || error "could not create server_exts.cfg"

echo "create ${SVR_BASE}.pem (signed server certificate)"
openssl req -key "${SVR_BASE}.key" -config server_request.cfg -new | openssl x509 -req -days "${DAYS}" -CA "${CA_BASE}.pem" -CAkey "${CA_BASE}.key" -CAcreateserial -"${HASH}" -extfile server_exts.cfg -out "${SVR_BASE}.pem"
test -f ${SVR_BASE}.pem || error "could not create ${SVR_BASE}.pem"

# --------------

# create client cert CSR and sign it, piped
cat > client_request.cfg <<EOCertConfig
[req]
prompt                 = no
distinguished_name     = req_distinguished_name

[req_distinguished_name]
commonName             = unbound-control
EOCertConfig
test -f client_request.cfg || error "could not create client_request.cfg"

cat > client_exts.cfg <<EOCertConfig
extendedKeyUsage       = clientAuth
EOCertConfig
test -f client_exts.cfg || error "could not create client_exts.cfg"

echo "create ${CTL_BASE}.pem (signed client certificate)"
openssl req -key "${CTL_BASE}.key" -config client_request.cfg -new | openssl x509 -req -days "${DAYS}" -CA "${CA_BASE}.pem" -CAkey "${CA_BASE}.key" -CAcreateserial -"${HASH}" -extfile client_exts.cfg -out "${CTL_BASE}.pem"
test -f "${CTL_BASE}.pem" || error "could not create ${CTL_BASE}.pem"

# --------------

# set desired permissions
chmod 0644 "${CA_BASE}.key" "${SVR_BASE}.key" "${CTL_BASE}.key"
chmod 0644 "${CA_BASE}.pem" "${SVR_BASE}.pem" "${CTL_BASE}.pem"

# cleanup
rm -f ca_request.cfg client_request.cfg client_exts.cfg server_request.cfg server_exts.cfg  *.srl

echo "Satisfy unbound daemon/remote.c SSL_CTX_use_certificate_chain_file by appending ${CA_BASE}.pem to ${SVR_BASE}.pem"
cat "${CA_BASE}.pem" >> "${SVR_BASE}.pem"

echo "Setup success. Certificates created in ${DESTDIR}."
