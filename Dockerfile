FROM golang:1.17-alpine3.14 AS builder

RUN apk add --no-cache \
  build-base \
  make \
  curl

WORKDIR /workspace
RUN mkdir _out

COPY go.mod go.sum Makefile unbound_exporter.go ./
RUN make clean-compile

FROM scratch
COPY --from=builder /workspace/_out/unbound_explorer /usr/local/bin/unbound_exporter

ENTRYPOINT ["unbound_exporter"]
