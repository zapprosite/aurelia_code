# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Instala dependências de build se necessário
RUN apk add --no-cache git

# Copia arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário estático (SOTA 2026 Hardening)
# -ldflags="-s -w" remove símbolos de debug conforme padrão industrial
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o aurelia-api main.go

# Stage 2: Runtime (Hardened)
FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/aurelia-api .

# Expõe a porta padrão de governança
EXPOSE 8080

# Go Runtime Polish (Padrão SOTA 2026.2)
ENV GOMEMLIMIT=512MiB
ENV GOGC=100

# Executa o binário
CMD ["./aurelia-api"]
