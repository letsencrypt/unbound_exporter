FROM library/golang:1.17 as build-env

RUN apt-get install -yq --no-install-recommends git

# Copy source + vendor
COPY . /go/src/github.com/letsencrypt/unbound_exporter
WORKDIR /go/src/github.com/letsencrypt/unbound_exporter

# Build
ENV GOPATH=/go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -a -ldflags "-s -w" -o /go/bin/unbound_exporter .

FROM scratch
COPY --from=build-env /go/bin/unbound_exporter /usr/bin/unbound_exporter
ENTRYPOINT ["unbound_exporter"]
CMD ["-unbound.host", "tcp://localhost:8953"]
