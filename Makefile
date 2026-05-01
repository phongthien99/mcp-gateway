.PHONY: all up down mcp-build mcp-run hugo-book-run hugo-book-docker

# ── Monorepo ──────────────────────────────────────────────────────────────────
all: up

up:
	docker compose up

down:
	docker compose down

# ── MCP Server + Dashboard ────────────────────────────────────────────────────
mcp-build:
	cd apps/mcp-server && go build -o ../../bin/mcp-workbench .

mcp-run:
	cd apps/mcp-server && DASHBOARD_SITE_DIR=../../apps/hugo-book-site go run .

# ── Standalone Hugo Book app ─────────────────────────────────────────────────
hugo-book-run:
	cd apps/hugo-book-site && hugo server --bind=0.0.0.0 --port=$${HUGO_BOOK_PORT:-1314} --disableFastRender --noBuildLock --poll=700ms

hugo-book-docker:
	docker compose up hugo-book-site
