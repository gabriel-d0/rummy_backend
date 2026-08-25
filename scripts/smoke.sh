#!/usr/bin/env bash
# Smoke test — verifies Nakama starts, runtime loads, DB connects. Repeatable.
# Phase 1 Day 7: Replaces manual curl checks with a single `make smoke` / `./scripts/smoke.sh`.
# Exit 0 if all checks pass, non-zero otherwise. Designed for local dev and CI (no external deps beyond curl/docker).
set -euo pipefail

PROJECT="rummy_backend"
NAKAMA_CONTAINER="rummy_nakama"
POSTGRES_CONTAINER="rummy_postgres"
COMPOSE="docker compose"
CURL="curl -s --max-time 5"
TIMEOUT=90
INTERVAL=2

# Colors if stdout is tty
if [ -t 1 ]; then GREEN="\033[32m"; RED="\033[31m"; YELLOW="\033[33m"; RESET="\033[0m"; else GREEN=""; RED=""; YELLOW=""; RESET=""; fi

ok() { echo -e "${GREEN}✔ $*${RESET}"; }
warn() { echo -e "${YELLOW}⚠ $*${RESET}"; }
fail() { echo -e "${RED}✘ $*${RESET}"; exit 1; }
step() { echo -e "${YELLOW}→ $*${RESET}"; }

# Ensure compose is up (with build if needed). If already healthy, skip rebuild.
ensure_up() {
  step "Ensuring compose stack is up (project $PROJECT)"
  if ! docker ps --filter "name=$POSTGRES_CONTAINER" --format "{{.Names}}" | grep -q "$POSTGRES_CONTAINER" || \
     ! docker ps --filter "name=$NAKAMA_CONTAINER" --format "{{.Names}}" | grep -q "$NAKAMA_CONTAINER"; then
    warn "Stack not running — bringing up with build"
    docker compose up --build -d
  else
    ok "Stack already running"
  fi
}

wait_for_health() {
  local container=$1
  local label=$2
  step "Waiting for $label ($container) to be healthy (timeout ${TIMEOUT}s)"
  local elapsed=0
  while [ $elapsed -lt $TIMEOUT ]; do
    health=$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}no-health{{end}}' "$container" 2>/dev/null || echo "missing")
    if [ "$health" = "healthy" ]; then
      ok "$label healthy after ${elapsed}s"
      return 0
    fi
    if [ "$health" = "missing" ] || [ "$health" = "no-health" ]; then
      # Fallback: check running
      if docker ps --filter "name=$container" --format "{{.Status}}" | grep -q "Up"; then
        # container up but no health — wait a bit then treat as ok if nakama not yet healthcheck
        sleep $INTERVAL
        elapsed=$((elapsed + INTERVAL))
        continue
      else
        fail "$label container not found/running"
      fi
    fi
    sleep $INTERVAL
    elapsed=$((elapsed + INTERVAL))
    if [ $((elapsed % 10)) -eq 0 ]; then echo "  ... $elapsed s ($health)"; fi
  done
  docker compose logs --tail=30 "$label" 2>&1 | tail -n 40 || true
  fail "Timeout waiting for $label healthy (last: $health)"
}

check_db_connectivity() {
  step "Checking DB connectivity (postgres:5433 PgIsReady + psql SELECT 1)"
  # pg_isready inside container
  if ! docker exec "$POSTGRES_CONTAINER" pg_isready -U postgres -d nakama -q 2>/dev/null; then
    docker compose logs --tail=20 postgres 2>&1 | tail -n 20
    fail "pg_isready failed inside $POSTGRES_CONTAINER"
  fi
  ok "pg_isready ok"

  # psql SELECT 1 via exec
  if ! docker compose exec -T postgres psql -U postgres -d nakama -tAc "SELECT 1" 2>/dev/null | grep -q "1"; then
    fail "psql SELECT 1 failed"
  fi
  ok "psql SELECT 1 ok"

  # compose exec host-side also via 5433
  if command -v psql >/dev/null 2>&1; then
    if PGPASSWORD=localdb psql -h 127.0.0.1 -p 5433 -U postgres -d nakama -tAc "SELECT 1" 2>/dev/null | grep -q "1"; then
      ok "host psql 127.0.0.1:5433 ok (optional)"
    else
      warn "host psql not reachable (optional, not failing)"
    fi
  fi
}

check_nakama_healthcheck() {
  step "Checking Nakama internal healthcheck (docker exec /nakama/nakama healthcheck)"
  if ! docker exec "$NAKAMA_CONTAINER" /nakama/nakama healthcheck >/dev/null 2>&1; then
    docker exec "$NAKAMA_CONTAINER" /nakama/nakama healthcheck 2>&1 || true
    fail "nakama healthcheck binary failed"
  fi
  ok "nakama healthcheck ok"
}

check_runtime_loaded() {
  step "Checking Go runtime loaded (InitModule + rummy_backend.so)"
  # Give Nakama a moment to finish runtime init after healthy
  sleep 2
  logs=$(docker compose logs nakama 2>&1 | tail -n 200)
  if ! echo "$logs" | grep -q "Rummy backend InitModule"; then
    echo "$logs" | tail -n 40
    fail "Log missing 'Rummy backend InitModule' — Go plugin not loaded"
  fi
  ok "Log contains Rummy backend InitModule"

  if ! echo "$logs" | grep -q 'Found runtime modules.*rummy_backend.so'; then
    echo "$logs" | grep -E "Found runtime modules|Rummy" | tail -n 20
    fail "Log missing 'Found runtime modules count 1 [rummy_backend.so]'"
  fi
  ok "Log contains Found runtime modules rummy_backend.so"

  if ! echo "$logs" | grep -q 'Registered Go runtime RPC.*health'; then
    fail "Log missing Registered Go RPC health"
  fi
  ok "RPC health/version registered"

  # Also verify .so exists inside container (baked)
  if ! docker exec "$NAKAMA_CONTAINER" test -f /nakama/data/modules/rummy_backend.so; then
    fail "/nakama/data/modules/rummy_backend.so not found inside container"
  fi
  ok "rummy_backend.so exists inside container"
}

check_console_and_api() {
  step "Checking console http://127.0.0.1:7351 (200) and API http://127.0.0.1:7350"
  code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:7351/ || true)
  if [ "$code" != "200" ]; then
    curl -v http://127.0.0.1:7351/ 2>&1 | head -n 20 || true
    fail "Console not 200 (got $code)"
  fi
  ok "Console 200"

  # API gateway root may be 404 but should not be connection refused; check that 7350 is listening
  if ! curl -s --max-time 3 http://127.0.0.1:7350/ >/dev/null 2>&1; then
    # Try health endpoint via nakama healthcheck already passed, so warn but not fail if 404
    warn "API root curl failed (may be 404, not critical if healthcheck passed)"
  else
    ok "API 7350 reachable"
  fi
}

check_rpc_health() {
  step "Checking RPC health (device auth via defaultkey + POST /v2/rpc/health)"
  UUID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "00000000-0000-0000-0000-000000000001")
  UUID=$(echo "$UUID" | tr '[:upper:]' '[:lower:]')
  TOKEN=$(curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/device?create=true" \
    -H "Content-Type: application/json" --user "defaultkey:" --data "{\"id\":\"$UUID\"}" 2>&1 | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('token',''))" 2>&1 || true)
  if [ -z "$TOKEN" ] || [ "$TOKEN" = "" ]; then
    fail "Device auth failed — no token (defaultkey invalid or Nakama down)"
  fi
  ok "Device auth ok (token ${#TOKEN} chars)"

  payload=$(curl -s -X POST "http://127.0.0.1:7350/v2/rpc/health" -H "Authorization: Bearer $TOKEN" --data '""' 2>&1 || true)
  # Validate outer+inner JSON — outer payload is stringified JSON containing status ok
  if ! echo "$payload" | python3 -c "import sys,json; d=json.load(sys.stdin); p=json.loads(d['payload']); assert p['status']=='ok' and p['service']=='rummy_backend'" 2>&1; then
    echo "$payload" | python3 -m json.tool 2>&1 | head -n 30 || echo "$payload"
    fail "RPC health payload validation failed (expected status ok, service rummy_backend)"
  fi
  ok "RPC health payload ok (status ok, service rummy_backend)"

  # Also version RPC as secondary check
  ver_payload=$(curl -s -X POST "http://127.0.0.1:7350/v2/rpc/version" -H "Authorization: Bearer $TOKEN" --data '""' 2>&1 || true)
  if ! echo "$ver_payload" | grep -q "rummy_backend"; then
    echo "$ver_payload"
    fail "RPC version missing rummy_backend"
  fi
  ok "RPC version ok"

  # Verify logs contain rpcHealth call (async)
  sleep 1
  if ! docker compose logs nakama 2>&1 | tail -n 30 | grep -q "rpcHealth called"; then
    warn "Log missing rpcHealth called (may be timing)"
  else
    ok "Log shows rpcHealth called"
  fi
}

main() {
  echo "=== Rummy Backend Smoke Test ==="
  echo "Project: $PROJECT | timeout ${TIMEOUT}s | $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  ensure_up
  wait_for_health "$POSTGRES_CONTAINER" "postgres"
  wait_for_health "$NAKAMA_CONTAINER" "nakama"
  check_db_connectivity
  check_nakama_healthcheck
  check_runtime_loaded
  check_console_and_api
  check_rpc_health
  echo ""
  echo -e "${GREEN}=== SMOKE PASSED ===${RESET}"
  echo "All checks passed — Nakama Go runtime, DB, and RPCs healthy."
  echo "Next: try 'open http://127.0.0.1:7351' (admin/password) or 'make health'."
}

main "$@"
