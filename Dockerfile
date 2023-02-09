# syntax = docker/dockerfile:1.3-labs

ARG VERSION=0.0.0
ARG COMMIT=''
ARG DATE=''

FROM golang:1-alpine as builder
WORKDIR /go/src/mysqlrouter_exporter
COPY . .
RUN apk --no-cache add git openssh build-base
RUN go build -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o app .

FROM alpine as production
LABEL maintainer="rluisr" \
  org.opencontainers.image.url="https://github.com/rluisr/mysqlrouter_exporter" \
  org.opencontainers.image.source="https://github.com/rluisr/mysqlrouter_exporter" \
  org.opencontainers.image.vendor="rluisr" \
  org.opencontainers.image.title="mysqlrouter_exporter" \
  org.opencontainers.image.description="Prometheus exporter for MySQL Router." \
  org.opencontainers.image.licenses="AGPL"
RUN <<EOF
    apk add --no-cache ca-certificates libc6-compat \
    rm -rf /var/cache/apk/*
EOF
COPY --from=builder /go/src/mysqlrouter_exporter/app /app
ENTRYPOINT ["/app"]

