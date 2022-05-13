FROM golang:1.18.1-bullseye AS builder

WORKDIR /build

COPY . .

FROM golang:1.18.1-bullseye

WORKDIR /app

COPY --from=builder /build/api /app/api

COPY docker-entrypoint.sh .

ENTRYPOINT ["/app/docker-entrypoint.sh"]
