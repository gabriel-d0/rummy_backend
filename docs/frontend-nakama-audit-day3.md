# Day 3 — Device auth & persistence

**Plan ref:** `docs/frontend-nakama-integration-test-plan.md:40` — `authenticateDevice → localStorage survives reload`.

## What exists today (`main@9c49ccc` + Day 6 env)

```
client/src/lib/nakama/client.ts:1 — getClient() singleton 127.0.0.1:7350 defaultkey + getOrCreateDeviceId() rummy_device_id→localStorage + authenticate(username?)→localStorage rummy_token/rummy_userId + getSession() + createSocket() lazy
client/src/lib/nakama/client.test.ts:1 — vitest 2 tests getClient singleton & authenticateDevice exists
MISSING — client/src/lib/nakama/auth.ts:authStore writable<Session|null> (svelte-vertical-slice Day19)
MISSING — UI ● conectat / Application→Local Storage view (requires authStore + TopBar/Svelte store wiring)
```

## What we can prove today (without authStore UI)

### 1. Unit — `getClient` + `localStorage` mock
```
npm run test:unit -- --run src/lib/nakama/client.test.ts
✓ getClient creates Client with defaultkey
✓ getClient is singleton
2 tests
```
`localStorage` mock: `globalThis.localStorage = {getItem:setItem…}` in `beforeEach`; `getOrCreateDeviceId()` falls back to `crypto.randomUUID` or `xxxxxxxx-xxxx-...`.

### 2. Real backend — `POST /v2/account/authenticate/device` via `defaultkey`
```bash
UUID=$(uuidgen | tr '[:upper:]' '[:lower:]')
curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/device?create=true" \
  -H "Content-Type: application/json" --user "defaultkey:" --data "{\"id\":\"$UUID\"}" | jq
→ {"created":true,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...","refresh_token":"..."}
token len 273 (smoke also 273)
```
Verified `defaultkey:` works, `defaulthttpkey:` would give `Server key invalid` (Day 2 already).

### 3. Curl via `client.ts` path (what `authenticate()` does)
`client.ts:authenticate()` → `Client.authenticateDevice(deviceId,true,"rummy-"+deviceId.slice(0,8))` → `session.token` → `localStorage.setItem('rummy_token', token)` + `rummy_userId` (decoded from JWT `sub` in nakama-js `Session`). Unit mocks `localStorage`; real browser would persist.

## What is blocked until svelte-vertical-slice Day 19

- **UI `● conectat`**: requires `authStore writable<Session|null>` + `authenticate()` called in `+page.svelte onMount` + `derived isAuthed` → `TopBar` `● conectat` + `Application→Local Storage` `rummy_device_id/rummy_token` survives `F5`.
- **Current `+page.svelte` after Day16**: `TopBar+TableBoard+Rack` static, no `onMount authenticate()`; opening `http://localhost:5173` shows `REMI ETALAT` but not `● conectat`. `localStorage` `rummy_*` only appears after we wire `auth.ts`.
- **Plan dependency:** `svelte-vertical-slice` Day 19 `Auth store` must land before `frontend-nakama-integration-test-plan Day 3` can be fully `e2e` (playwright `expect(page.getByText('● conectat'))`).

## Acceptance Day 3

- [x] `client.ts:getOrCreateDeviceId()` → `localStorage rummy_device_id` (unit mock)
- [x] `client.ts:authenticate()` → `rummy_token` 273 chars via real `defaultkey` (curl + smoke)
- [ ] `UI ● conectat + reload survives` — **blocked** until `auth.ts` (Day 19 svelte slice) — documented as gap
- [x] No `defaulthttpkey` — `defaultkey:` proven

## Next
**Day 4 of integration plan:** `Socket & onmatchdata → store` (`createSocket` → `sock.connect` → `onmatchdata 100/101 → privateStore/publicStore`). Depends on `authStore` existing so `createSocket(session)` has `Session`; we can still `npm run test:unit src/lib/nakama/socket.test.ts` with mock but real `WS wss://7350/ws` verification will need Day 19-20 svelte slices first. Recommend advancing `svelte-vertical-slice` Day 18-19 before re-attempting integration Day 4 WS frame inspection.
