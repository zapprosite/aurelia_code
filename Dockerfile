# Stage 1: Build
# modernc.org/sqlite is pure Go — CGO_ENABLED=0 works and produces a static binary.
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Cache deps before copying source
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o aurelia ./cmd/aurelia

# Stage 2: Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 1000 aurelia

WORKDIR /app
COPY --from=builder /build/aurelia /app/aurelia

# Dashboard (3334) + Health (8484)
EXPOSE 3334 8484

USER aurelia

# AURELIA_HOME must be a writable volume in production (SQLite + config)
ENV AURELIA_HOME=/home/aurelia/.aurelia

ENTRYPOINT ["/app/aurelia"]
