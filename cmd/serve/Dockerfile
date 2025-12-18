ARG HTTP_PORT=4001
FROM golang:1.25.5-alpine AS builder
RUN apk add --no-cache ca-certificates git
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux

WORKDIR /app

COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /app/main ./cmd/serve/main.go

FROM alpine:3.19 AS production
ARG HTTP_PORT
ENV HTTP_PORT=${HTTP_PORT}
LABEL org.opencontainers.image.source="https://github.com/salihguru/idiogo"

RUN addgroup -S app && adduser -S app -G app

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /usr/local/bin/main
COPY --from=builder /app/assets /app/assets

RUN chown -R app:app /app /usr/local/bin/main

USER app
EXPOSE ${HTTP_PORT}

CMD ["/usr/local/bin/main"]

FROM golang:1.25.5-alpine AS development
ARG HTTP_PORT
ENV HTTP_PORT=${HTTP_PORT}
RUN apk add --no-cache ca-certificates git curl
WORKDIR /app

COPY --from=builder /app /app

EXPOSE ${HTTP_PORT}

CMD ["/app/main"]