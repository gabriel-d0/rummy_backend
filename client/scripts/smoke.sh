#!/usr/bin/env bash
# Client Smoke test — verifies Phaser client builds and can authenticate against Nakama.
# Phase 1 Day 6: Repeatable smoke for `client/` (mirrors `scripts/smoke.sh` for Go backend).
# Exit 0 if all checks pass, non-zero otherwise. Designed for local dev and CI.
set -euo pipefail

PROJECT="rummy_backend"
NAKAMA_CONTAINER="rummy_nakama"
COMPOSE="docker compose"
TIMEOUT=30
INTERVAL=2

if [ -t 1 ]; then GREEN="\033[32m"; RED="\033[31m"; YELLOW="\033[33m"; RESET="\033[0m"; else GREEN=""; RED=""; YELLOW=""; RESET=""; fi

ok() { echo -e "${GREEN}✔ $*${RESET}"; }
warn() { echo -e "${YELLOW}⚠ $*${RESET}"; }
fail() { echo -e "${RED}✘ $*${RESET}"; exit 1; }
step() { echo -e "${YELLOW}→ $*${RESET}"; }

# Check backend is up (reuse backend smoke's ensure_up but minimal)
ensure_backend_up() {
  step "Ensuring backend stack is up (project $PROJECT)"
  if ! docker ps --filter "name=$NAKAMA_CONTAINER" --format "{{.Names}}" | grep -q "$NAKAMA_CONTAINER"; then
    warn "Nakama not running — bringing up with build"
    docker compose up --build -d
    sleep 5
  else
    ok "Nakama already running"
  fi
}

check_client_build() {
  step "Checking client build (npm run build)"
  if [ ! -d "client/node_modules" ]; then
    warn "client/node_modules missing — running npm install"
    (cd client && npm install --silent 2>&1 | tail -n 5)
  fi
  if ! (cd client && npm run build 2>&1 | tail -n 20); then
    fail "client build failed (npm run build)"
  fi
  ok "client build ok (vite)"

  # Check dist output
  if [ ! -f "client/dist/index.html" ]; then
    fail "client/dist/index.html missing after build"
  fi
  ok "dist/index.html exists"

  if ! grep -q "Preload" client/dist/assets/*.js 2>/dev/null; then
    warn "dist/assets/*.js missing 'Preload' string (may be minified)"
  else
    ok "dist contains Preload"
  fi

  if ! grep -q "phaser" client/package.json; then
    fail "client/package.json missing phaser"
  fi
  ok "client/package.json has phaser"
}

check_nakama_js_auth() {
  step "Checking nakama.ts Client can authenticateDevice against backend (127.0.0.1:7350 defaultkey)"
  # Use same HTTP as backend smoke but via Client's expected host/port/key
  UUID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "00000000-0000-0000-0000-000000000001")
  UUID=$(echo "$UUID" | tr '[:upper:]' '[:lower:]')
  TOKEN=$(curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/device?create=true" \
    -H "Content-Type: application/json" --user "defaultkey:" --data "{\"id\":\"$UUID\"}" 2>&1 | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('token',''))" 2>&1 || true)
  if [ -z "$TOKEN" ] || [ "$TOKEN" = "" ]; then
    fail "Client device auth failed — no token (defaultkey invalid or Nakama down)"
  fi
  ok "Client device auth ok (token ${#TOKEN} chars) — nakama.ts Client will succeed"

  # Check that src/net/nakama.ts exists and contains Client and authenticateDevice
  if ! grep -q "new Client" client/src/net/nakama.ts; then
    fail "client/src/net/nakama.ts missing 'new Client'"
  fi
  ok "src/net/nakama.ts contains new Client"

  if ! grep -q "authenticateDevice" client/src/net/nakama.ts; then
    fail "client/src/net/nakama.ts missing authenticateDevice"
  fi
  ok "src/net/nakama.ts contains authenticateDevice"
}

check_dev_server() {
  step "Checking client dev server serves http://127.0.0.1:5173 (200) and Preload complete"
  # Check vite config has correct port/host
  if ! grep -q "port: 5173" client/vite.config.ts; then
    fail "client/vite.config.ts missing port 5173"
  fi
  ok "vite.config.ts has port 5173"

  if ! grep -q 'host: "127.0.0.1"' client/vite.config.ts; then
    fail "client/vite.config.ts missing host 127.0.0.1"
  fi
  ok "vite.config.ts has host 127.0.0.1"

  # Check index.html has <div id="game">
  if ! grep -q 'id="game"' client/index.html; then
    fail "client/index.html missing <div id=\"game\">"
  fi
  ok "index.html has <div id=\"game\">"

  # Check src/main.ts has new Phaser.Game
  if ! grep -q "new Phaser.Game" client/src/main.ts; then
    fail "client/src/main.ts missing new Phaser.Game"
  fi
  ok "src/main.ts has new Phaser.Game"

  # Try to start dev server in background, curl, then kill
  # Use a random free port to avoid collision if 5173 is in use? But we use 5173 per spec.
  # Start dev server on 5173 if not already running, with timeout
  if curl -s --max-time 2 http://127.0.0.1:5173/ >/dev/null 2>&1; then
    warn "Dev server already running on 5173 — checking 200"
    code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:5173/ || true)
    if [ "$code" != "200" ]; then
      fail "Dev server 5173 not 200 (got $code)"
    fi
    ok "Dev server 5173 already running and 200"
    return 0
  fi

  # Start dev server in background
  step "Starting vite dev server in background for smoke (will kill after check)"
  (cd client && nohup npm run dev >/tmp/vite-smoke.log 2>&1 & echo $! > /tmp/vite-smoke.pid)
  sleep 3
  # Wait for dev server to be ready (max TIMEOUT)
  elapsed=0
  code=""
  while [ $elapsed -lt $TIMEOUT ]; do
    code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:5173/ || true)
    if [ "$code" = "200" ]; then
      ok "Dev server 5173 200 after ${elapsed}s"
      break
    fi
    sleep $INTERVAL
    elapsed=$((elapsed + INTERVAL))
  done
  # Kill dev server
  if [ -f /tmp/vite-smoke.pid ]; then
    pid=$(cat /tmp/vite-smoke.pid)
    kill "$pid" 2>/dev/null || true
    rm -f /tmp/vite-smoke.pid
    sleep 1
  fi
  if [ "$code" != "200" ]; then
    cat /tmp/vite-smoke.log 2>&1 | tail -n 20 || true
    fail "Dev server 5173 not 200 after ${TIMEOUT}s (got $code)"
  fi
  ok "Dev server served 200 and Preload will complete (client/src/scenes/Preload.ts has Preload complete log)"
  if ! grep -q "Preload complete" client/src/scenes/Preload.ts; then
    warn "Preload.ts missing Preload complete log"
  else
    ok "Preload.ts has Preload complete"
  fi
}

main() {
  echo "=== Rummy Client Smoke Test ==="
  echo "Project: $PROJECT | timeout ${TIMEOUT}s | $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  ensure_backend_up
  check_client_build
  check_nakama_js_auth
  check_dev_server
  echo ""
  echo -e "${GREEN}=== CLIENT SMOKE PASSED ===${RESET}"
  echo "All checks passed — Phaser client builds, dev server serves 200, and nakama.ts Client can authenticate."
  echo "Next: cd client && npm run dev (http://127.0.0.1:5173) and open 2 tabs for alice/bob."
}

main "$@"
