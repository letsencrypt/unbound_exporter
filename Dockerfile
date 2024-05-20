FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.21.4-bookworm AS build
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETPLATFORM go build -v -o /go/bin/unbound_exporter

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/unbound_exporter /
ENTRYPOINT ["/unbound_exporter"]
