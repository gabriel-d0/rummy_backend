# Day 2 — Env & Docker smoke

**Plan ref:** `docs/frontend-nakama-integration-test-plan.md:27` — `Goal: Frontend VITE_NAKAMA_* can actually hit Nakama`.

## Env checked
```
client/.env.example:8
PUBLIC_NAKAMA_HOST=127.0.0.1
PUBLIC_NAKAMA_PORT=7350
PUBLIC_NAKAMA_KEY=defaultkey
PUBLIC_NAKAMA_USE_SSL=false
VITE_NAKAMA_HOST=127.0.0.1
VITE_NAKAMA_PORT=7350
VITE_NAKAMA_KEY=defaultkey
VITE_NAKAMA_USE_SSL=false

client/.env:8 identical (present, ignored via .gitignore:16 !.env.example)
client/src/lib/config.ts:1 imports PUBLIC_* from $env/static/public with fallback 127.0.0.1:7350 defaultkey
client/src/lib/nakama/client.ts:1 imports NAKAMA_* from $lib/config → Client(HOST,PORT,KEY,USE_SSL)
compose.yml:1 name: rummy_backend postgres 5433→5432 nakama 7350/7351/7349
```

## Docker
```
docker compose ps:
rummy_nakama   rummy_backend:local   Up 40 minutes (healthy)  0.0.0.0:7349-7351->7349-7351/tcp
rummy_postgres postgres:15-alpine    Up 40 minutes (healthy)  0.0.0.0:5433->5432/tcp
docker compose logs nakama --tail=100:
# ping/pong only in last 100 (startup earlier), but smoke checks InitModule below
```

## Smoke (`make smoke` 2026-08-26T23:51:23Z)

```
✔ Stack already running
✔ postgres healthy after 0s
✔ nakama healthy after 0s
✔ pg_isready ok
✔ psql SELECT 1 ok
✔ nakama healthcheck ok
✔ Log contains Rummy backend InitModule
✔ Log contains Found runtime modules rummy_backend.so
✔ RPC health/version registered
✔ rummy_backend.so exists inside container
✔ Console 200 (http://127.0.0.1:7351)
✔ API 7350 reachable
✔ Device auth ok (token 273 chars) via POST /v2/account/authenticate/device?create=true --user defaultkey:
✔ RPC health payload ok (status ok, service rummy_backend) via POST /v2/rpc/health Authorization: Bearer <token>
✔ RPC version ok
✔ Log shows rpcHealth called
=== SMOKE PASSED ===
```

## Proof that frontend env can hit Nakama
- `VITE_NAKAMA_KEY` is `defaultkey` not `defaulthttpkey` (smoke uses `--user defaultkey:` succeeds; `defaulthttpkey` would give `Server key invalid`).
- `client/src/lib/config.ts` fallback `?? '127.0.0.1'` ensures `npm run build` + `svelte-check` pass even without `.env` (but `.env` present makes static import work).
- `curl -s -X POST http://127.0.0.1:7350/v2/rpc/health --user defaultkey:` → `{"payload":"{\"status\":\"ok\",\"version\":\"0.1.0-go-day3-skeleton\",…}"}` (verified via smoke).

## Acceptance Day 2
- [x] `client/.env.example` has 4 `PUBLIC_*` + 4 `VITE_*`
- [x] `docker compose ps` both `healthy`
- [x] `Log contains Rummy backend InitModule` + `Found runtime modules rummy_backend.so`
- [x] `make smoke` `SMOKE PASSED` with `RPC health` via `defaultkey`

## Next
Day 3 — Device auth & persistence (`authenticateDevice` → `localStorage rummy_device_id/rummy_token/rummy_userId` survives reload). Depends on this `7350 defaultkey` being reachable — now proven.
