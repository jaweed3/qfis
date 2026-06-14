#!/bin/sh
set -e

echo "==> Building frontend..."
cd frontend && npm ci && npm run build && cd ..

echo "==> Copying frontend dist to backend..."
rm -rf backend/frontend/dist
cp -r frontend/dist backend/frontend/dist

echo "==> Cross-compiling Go backend for linux/amd64..."
cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o qfis-linux . && cd ..

echo ""
echo "==> Done! Binary: backend/qfis-linux"
echo ""
echo "── Deploy ──────────────────────────────────────"
echo "  scp backend/qfis-linux user@server:~/qfis/"
echo "  ssh user@server"
echo "  mkdir -p ~/qfis/data"
  echo "  PORT=21465 GIN_MODE=release ./qfis-linux"
echo ""
echo "── Or via Docker ───────────────────────────────"
echo "  docker compose up -d"
echo "────────────────────────────────────────────────"
echo ""
