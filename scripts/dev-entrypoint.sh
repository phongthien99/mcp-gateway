#!/bin/sh
set -e

# Hugo hot reload
cd /hugo-src
hugo server \
  --bind=0.0.0.0 \
  --port=1313 \
  --disableFastRender \
  --noBuildLock \
  --poll=700ms &

# MCP server binary
/app/server &

# Exit if either process dies
wait -n 2>/dev/null || wait
