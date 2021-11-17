FROM golang:1.17.3 as builder

WORKDIR /workspace/

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o unbound_exporter ./

FROM scratch
COPY --from=builder /workspace/unbound_exporter unbound_exporter

ENTRYPOINT [ "/unbound_exporter" ]
