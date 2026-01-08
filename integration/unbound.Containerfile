FROM debian:forky

ARG UNBOUND_TAG=release-1.24.1

RUN groupadd --system unbound \
  && useradd --system --no-create-home --gid unbound --shell /usr/sbin/nologin unbound

RUN apt-get -y update && apt-get install -y \
    bison flex gcc make openssl \
    libevent-dev libexpat1-dev libnghttp2-dev libngtcp2-crypto-ossl-dev libngtcp2-dev libssl-dev

ADD https://github.com/NLnetLabs/unbound.git#${UNBOUND_TAG} /src/unbound

WORKDIR /src/unbound

RUN ./configure --enable-subnet --with-libevent --with-libnghttp2 --with-libngtcp2 \
    && make && make install

COPY configs /etc/unbound/

CMD ["/usr/local/sbin/unbound", "-v", "-d", "-c", "/etc/unbound/unbound-example.conf"]
