# Stage 1: Hugo binary
FROM mcp-gateway-hugo-site:latest AS hugo-src

# Stage 2: Build Go binary
FROM golang:alpine AS builder

RUN apk add --no-cache git
WORKDIR /build
COPY apps/mcp-server/go.mod apps/mcp-server/go.sum ./
RUN go mod download
COPY apps/mcp-server/ .
RUN go build -o /server .

# Stage 3: Slim runtime — binary + Hugo + baked Hugo site
FROM alpine:latest

RUN apk add --no-cache libstdc++ libgcc
COPY --from=hugo-src /usr/bin/hugo /usr/bin/hugo
COPY --from=builder /server /app/server
COPY apps/hugo-book-site/ /hugo-src/

COPY scripts/dev-entrypoint.sh /dev-entrypoint.sh
RUN chmod +x /dev-entrypoint.sh

EXPOSE 8099 8110 1313

ENTRYPOINT ["/dev-entrypoint.sh"]
