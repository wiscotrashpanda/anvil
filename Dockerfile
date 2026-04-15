# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/anvil ./cmd/anvil

FROM gcr.io/distroless/static-debian12

COPY --from=builder /out/anvil /usr/local/bin/anvil

ENTRYPOINT ["/usr/local/bin/anvil"]
CMD ["--help"]
