# Build context is the monorepo root.
# ── Stage 1: build ────────────────────────────────────────────────────────────
FROM golang:alpine AS builder

WORKDIR /build

COPY apps/mcp-server/go.mod apps/mcp-server/go.sum ./
RUN go mod download

COPY apps/mcp-server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o mcp-workbench .

# ── Stage 2: runtime ──────────────────────────────────────────────────────────
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/mcp-workbench .

RUN mkdir -p \
    artifacts \
    workflows \
    prompts \
    context/global \
    runs

EXPOSE 8099 8110

CMD ["./mcp-workbench"]
