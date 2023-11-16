FROM docker.io/library/golang:1.21.4-bookworm AS build

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY *.go .

ENV CGO_ENABLED=0

RUN go build -v -o /go/bin/unbound_exporter ./...

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/unbound_exporter /

ENTRYPOINT ["/unbound_exporter"]
