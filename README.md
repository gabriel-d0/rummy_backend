# Rummy Backend — Romanian Tile Rummy (Nakama authoritative, Go)

Server-authoritative multiplayer Romanian Tile Rummy for 2–4 players, implemented as a Nakama Go runtime plugin. This is the **Go** version — migrated from TypeScript on 2026-08-25 (`55c7f3b`). See `docs/project-baseline.md` §13 for migration rationale.

The game follows Romanian Tile Rummy (106 tiles, 2 jokers, 50-point opening meld with at least one run, anticlockwise turns) per `AGENTS.md` and will be built incrementally (“Handmade Hero” vertical slices). Current phase is **Phase 1 — Foundation & local dev**; no gameplay rules beyond skeleton RPCs are implemented yet.

## Prerequisites

- Docker Desktop 4.x+ (`docker --version` 29.x, `docker compose version` v5.x) — tested on `darwin 25.6.0 arm64`.
- Go 1.23.x (`go version` — `go.mod` pins `1.23.5` to match `heroiclabs/nakama:3.26.0` runtime `go1.23.5`; local `1.27` works but `go vet` should pass on `1.23`).
- Git.
- No local `psql` required — Postgres runs in Docker at `5433`. Host `5432` is left for the unrelated `tinybot` stack.

## Quick Start

```bash
# 1. Clone (already done)
git clone git@github.com:gabriel-d0/rummy_backend.git
cd rummy_backend

# 2. Copy env defaults (non-secret, local only)
cp .env.example .env   # optional — compose.yml has defaults via ${VAR:-default}

# 3. Build Go plugin and start DB + Nakama
docker compose up --build -d

# 4. Watch startup (should end with "Startup done" and Go module log)
docker compose logs -f nakama --tail=50
# expected: Rummy backend InitModule — Romanian Tile Rummy v0.1.0 Go Day 3 skeleton starting
# expected: Found runtime modules count 1 [rummy_backend.so]
# expected: Registered Go runtime RPC function invocation id health/version

# 5. Verify console and API
open http://127.0.0.1:7351   # Nakama console, admin/password (local)
curl -s http://127.0.0.1:7350/ | head   # API gateway (empty 200/404 is ok, not connection refused)
```

**First-time DB init** takes ~10s for Postgres healthcheck + Nakama `migrate up`. Re-running `docker compose up -d` without `--build` reuses the existing `rummy_backend_pgdata` volume and skips migrations (`Successfully applied migration count 0`).

## Services & Ports

| Service | Container | Image / Build | Host port → Container | Purpose |
|---|---|---|---|---|
| `postgres` | `rummy_postgres` | `postgres:15-alpine` | `5433 → 5432` | Nakama DB — avoids `tinybot`’s `5432` (see `docs/project-baseline.md:6`). Uses volume `rummy_backend_pgdata`. |
| `nakama` | `rummy_nakama` | Built `rummy_backend:local` via `Dockerfile` (`nakama-pluginbuilder:3.26.0` → `nakama:3.26.0`) | `7350 → 7350` (HTTP API), `7351 → 7351` (Console), `7349 → 7349` (gRPC) | Authority, runs `main.go:12` `InitModule` → `health`/`version` RPCs. Plugin baked as `/nakama/data/modules/rummy_backend.so`; only `local.yml` is host-mounted. |

Compose project name is `rummy_backend` (`compose.yml:1` `name:`) to isolate volumes/networks from `tinybot`.

## Development Workflow

### Rebuild after Go changes

Go plugin is **baked into the image**, not volume-mounted like the previous TS `build/index.js` flow (TS mounted `./nakama/data` and hid the image; Go failed with `plugin was built with ... pragma` until we baked). After editing `main.go` or any Go file:

```bash
docker compose up --build -d
docker compose logs nakama --tail=40 | grep -E "Rummy|modules|Registered"
# expect: Found runtime modules count 1 [rummy_backend.so]
```

**Do NOT** run `go build --buildmode=plugin` on macOS host — the `.so` would be `darwin` and fail inside `linux` container. Always build via `Dockerfile`.

### Start / Stop / Logs / Clean

```bash
docker compose up -d                # start (with existing image)
docker compose up --build -d        # start + rebuild Go plugin
docker compose ps                   # status, health (both should be healthy)
docker compose logs nakama --tail=100   # Nakama logs
docker compose logs postgres --tail=50
docker compose restart nakama        # reload after local.yml edit (no rebuild)
docker compose down                  # stop, preserve PG volume + host plugin baked in image
docker compose down -v               # stop AND wipe DB (pgdata volume) — next up will initdb
docker compose exec postgres psql -U postgres -d nakama -c "SELECT * FROM migrate;"  # DB inspect
```

### Go commands (on host, for lint/vet/test before building image)

```bash
go vet ./...          # static checks (required before commit)
go test ./...         # no tests yet — Day 3 skeleton; future days add go test with deterministic seeds
go fmt ./...          # format (use go fmt, not prettier)
go mod tidy           # sync deps after editing go.mod (keeps go 1.23.5, nakama-common v1.36.0, protobuf v1.36.4)
```

Protobuf is pinned to `v1.36.4` — `v1.36.6` caused `plugin was built with a different version of pragma` crash on `heroiclabs/nakama:3.26.0` (`go1.23.5`). Keep `go 1.23` line in `go.mod` aligned with builder.

### RPC smoke test (proves runtime loaded)

Nakama HTTP key is `defaultkey` (not `defaulthttpkey` — the latter gives `Server key invalid`).

```bash
UUID=$(uuidgen | tr '[:upper:]' '[:lower:]')
TOKEN=$(curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/device?create=true" \
  -H "Content-Type: application/json" --user "defaultkey:" --data "{\"id\":\"$UUID\"}" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")

curl -s -X POST "http://127.0.0.1:7350/v2/rpc/health" -H "Authorization: Bearer $TOKEN" --data '""' | python3 -m json.tool
# {"payload": "{\"status\":\"ok\",\"version\":\"0.1.0-go-day3-skeleton\",...}"}

curl -s -X POST "http://127.0.0.1:7350/v2/rpc/version" -H "Authorization: Bearer $TOKEN" --data '""' | python3 -m json.tool
# {"payload": "{\"runtime\":\"go\",\"phase\":\"1-day3-go\",...}"}
```

Or use Nakama Console → API Explorer → RPC → `health` with the same `$TOKEN`.

## Project Structure

```
.
├── main.go                # Nakama InitModule (Go) — registers health/version RPCs (main.go:22)
├── go.mod / go.sum        # module github.com/gabriel-d0/rummy_backend, nakama-common v1.36.0
├── Dockerfile             # multi-stage pluginbuilder:3.26.0 → nakama:3.26.0, outputs backend.so
├── compose.yml            # name rummy_backend, postgres 5433 + nakama 7350/7351, builds Dockerfile
├── nakama/data/local.yml  # minimal DEBUG config, mounted ro
├── .env.example           # non-secret local defaults (5433, 7350/7351, admin/password)
├── .gitignore / .dockerignore # Go ignores (*.so, vendor) + annotated legacy Node
├── docs/project-baseline.md # Day 1 audit + §13 Go migration amendment
├── docs/                  # future: rules-decisions.md, architecture.md, protocol.md
└── AGENTS.md              # product spec — source of truth, includes 24-day plan
```

Future `internal/` layout per `AGENTS.md:133` will be Go packages: `internal/match`, `internal/rules`, `internal/setup`, `internal/protocol` (mirrors previous `src/` plan).

## How to Inspect

- **Nakama console:** `http://127.0.0.1:7351` user `admin` password `password` (local only).
- **Logs:** `docker compose logs -f nakama` — look for `InitModule` and `Found runtime modules`.
- **Healthcheck:** `docker inspect rummy_nakama --format '{{json .State.Health}}' | python3 -m json.tool` and `docker exec rummy_nakama /nakama/nakama healthcheck`.
- **DB:** `docker compose exec postgres pg_isready -U postgres -d nakama` or `psql` inside container.

## Troubleshooting

- **`plugin was built with a different version of ... pragma`** → protobuf mismatch; ensure `google.golang.org/protobuf v1.36.4` and `nakama-common v1.36.0` in `go.mod` and rebuild `docker compose build --no-cache`.
- **`Server key invalid` on device auth** → use `--user defaultkey:` not `defaulthttpkey:`.
- **`mkdir /nakama/data/modules: read-only file system`** (legacy TS bug) → fixed by mounting only `local.yml` not whole `data` dir (Go) — see `compose.yml:51`.
- **Port collision with tinybot** → rummy uses `5433` not `5432`; check `docker ps` and `.env.example`.
- **Nakama not loading new Go code** → you edited `main.go` but forgot `docker compose up --build -d`; host `go run` does not affect container.

## Docs & Decisions

- `docs/project-baseline.md` — Day 1 audit + language amendment §13.
- `AGENTS.md` — full product spec, tile set, meld rules, increment plan.
- Next docs (planned): `docs/rules-decisions.md` (50-point opening, jokers), `docs/architecture.md`, `docs/protocol.md`.

## Next Steps (Roadmap Phase 1)

- **Day 4** — Developer scripts: `Makefile` helpers for `go vet/test/fmt` + compose shortcuts.
- **Day 5** — Local setup docs: expand this `README.md` with architecture overview & protocol stub.
- **Day 6** — CI: GitHub Actions `go vet` + `docker compose build` smoke.

---

*Stack: Nakama 3.26.0 (`go1.23.5`) on `postgres:15-alpine`, Go 1.23.5 plugin (`rummy_backend.so`), Docker Compose `rummy_backend:local`.*
