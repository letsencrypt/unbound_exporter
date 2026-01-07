FROM debian:forky

RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*

COPY setup-certs.sh /setup-certs.sh

CMD ["/setup-certs.sh"]
