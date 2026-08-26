# Rummy Backend — Developer Scripts (Go + Docker)
# Phase 1 Day 4+7: start, stop, logs, clean/reset, build, test, lint, type-check, format + smoke
# Usage: `make help` for list; all targets are explicit and non-magic per Handmade Hero.

.PHONY: help build up up-build down restart logs logs-nakama logs-postgres ps clean clean-all reset db-shell vet fmt fmt-check test tidy health smoke check

# Default: show help
help:
	@echo "Rummy Backend — Go + Nakama dev commands"
	@echo ""
	@echo "Docker (Nakama + Postgres 15 on 5433 | 7350/7351):"
	@echo "  make up          - build Go plugin image and start db+nakama (docker compose up --build -d)"
	@echo "  make up-build    - alias for up"
	@echo "  make down        - stop without wiping DB volume"
	@echo "  make restart     - restart nakama (e.g. after local.yml edit, no rebuild)"
	@echo "  make ps          - show compose ps + health"
	@echo "  make logs        - tail both nakama+postgres logs (100 lines)"
	@echo "  make logs-nakama - tail nakama logs (alias for inspecting Rummy InitModule)"
	@echo "  make logs-postgres - tail postgres"
	@echo "  make clean       - down (preserve PG volume)"
	@echo "  make clean-all   - down -v (WIPE DB volume rummy_backend_pgdata) + remove Go build artifacts"
	@echo "  make reset       - clean-all + up --build -d (fresh DB + fresh plugin)"
	@echo "  make db-shell    - psql into nakama DB (postgres:5433)"
	@echo ""
	@echo "Go (local, no Docker):"
	@echo "  make vet         - go vet ./... (type-check + static checks, required before commit)"
	@echo "  make fmt         - go fmt ./... (format, writes files)"
	@echo "  make fmt-check   - check if formatted (fails if go fmt would change)"
	@echo "  make test        - go test ./... (no tests yet Day 4; placeholder)"
	@echo "  make tidy        - go mod tidy"
	@echo "  make check       - vet + fmt-check + test (CI baseline)"
	@echo "  make health      - curl health: console 200, device auth + RPC health"
	@echo "  make smoke       - full smoke: DB+nakama healthy, runtime loaded, DB conn, console+RPC (scripts/smoke.sh)"
	@echo "  make cli         - run minimal Rummy CLI (local simulation, 2 players, shows Private vs Public)"
	@echo "  make cli-help    - show CLI help"
	@echo ""

# --- Docker ---
build:
	docker compose build

up: build
	docker compose up -d

up-build: up

down:
	docker compose down

restart:
	docker compose restart nakama

ps:
	docker compose ps
	docker compose logs --tail=20 nakama | tail -n 20

logs:
	docker compose logs --tail=100

logs-nakama:
	docker compose logs --tail=100 nakama

logs-postgres:
	docker compose logs --tail=100 postgres

clean:
	docker compose down

clean-all:
	docker compose down -v
	rm -f backend.so

reset: clean-all
	docker compose up --build -d
	@echo "Reset done — Postgres re-init + Go plugin rebuilt. Wait 8s then 'make logs-nakama'."

db-shell:
	docker compose exec postgres psql -U postgres -d nakama

# --- Go ---
vet:
	go vet ./...

fmt:
	go fmt ./...

fmt-check:
	@test -z "$$(gofmt -l . | grep -v vendor | head -n 20)" || (echo "Not formatted — run 'make fmt'"; gofmt -l . | head -n 20; exit 1)

test:
	go test ./...

tidy:
	go mod tidy

check: vet fmt-check test

# --- Health / smoke (uses host curl, not inside Docker) ---
health:
	@echo "Checking console http://127.0.0.1:7351 ..."
	@curl -s -o /dev/null -w "console %{http_code}\n" http://127.0.0.1:7351/ || true
	@echo "Device auth via defaultkey + RPC health ..."
	@UUID=$$(uuidgen | tr '[:upper:]' '[:lower:]') && \
	TOKEN=$$(curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/device?create=true" \
	  -H "Content-Type: application/json" --user "defaultkey:" --data "{\"id\":\"$$UUID\"}" | python3 -c "import sys,json; print(json.load(sys.stdin).get('token',''))") && \
	if [ -z "$$TOKEN" ]; then echo "auth failed"; exit 1; fi && \
	echo "token $${#TOKEN} chars" && \
	curl -s -X POST "http://127.0.0.1:7350/v2/rpc/health" -H "Authorization: Bearer $$TOKEN" --data '""' | python3 -m json.tool | head -n 20

smoke:
	./scripts/smoke.sh

# --- Minimal test client (Day 24) ---
cli:
	@echo "Starting minimal Rummy CLI (local simulation) — see cmd/rummy-cli/main.go"
	@go run ./cmd/rummy-cli

cli-help:
	@go run ./cmd/rummy-cli --help
