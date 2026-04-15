# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/anvil ./cmd/anvil

FROM gcr.io/distroless/static-debian12

COPY --from=builder /out/anvil /usr/local/bin/anvil

ENTRYPOINT ["/usr/local/bin/anvil"]
CMD ["--help"]
