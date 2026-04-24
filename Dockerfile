# ── Stage 1: build ────────────────────────────────────────────────────────────
FROM golang:alpine AS builder

WORKDIR /build

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o mcp-gateway .

# ── Stage 2: runtime ──────────────────────────────────────────────────────────
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/mcp-gateway .

# Create default directories — will be overridden by volumes in compose
RUN mkdir -p \
    artifacts \
    workflows \
    prompts \
    context/global \
    runs

EXPOSE 8099

CMD ["./mcp-gateway"]
