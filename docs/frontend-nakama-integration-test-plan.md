# Svelte ↔ Nakama Integration — Test Plan (Day-by-Day, Handmade Hero)

**Goal:** Prove that the SvelteKit frontend (`client/` at `http://localhost:5173`) talks to the authoritative Go/Nakama backend (`docker compose` at `7350`/`7351`) for real — not just mocks — and that every player-visible state comes from `OpServerState 100`/`101` and every action goes through `1..9` → `MatchLoop` → `CheckTileConservation 106` → `OpServerError 102` when wrong.

**Rewrite note:** This file replaces `docs/nakama-gameplay-test-plan.md`’s backend-centric view with a **frontend-connected-to-backend** view. Backend plan remains valid; this file is how we _observe it from the browser_. For pure rules determinism see `docs/testing.md`; for opcodes see `docs/protocol.md:1` (`Version 1`, `1..9`/`100..103`).

**Scope:** `client/src/lib/nakama/*` (`Client`, `Session`, `Socket`, `match`, `auth`, `reconnect`), `client/src/lib/game/*` (`snapshot`, `store`, `errorStore`, `actions`), `client/src/components/*` (`Rack`, `TableBoard`, `TopBar`, `Tile`, `Toast`, `WinnerOverlay`), `client/src/routes/*` (`+page.svelte` lobby, `game/+page.svelte`, `demo/*`), `client/tests/*` (`playwright` `*.e2e.ts` + `vitest` `*.test.ts`), `compose.yml` (`postgres:15` `5433`, `nakama:3.26.0` `7350`/`7351`).

**Definition of “connected and working”:** After `docker compose up --build -d` + `make smoke` `SMOKE PASSED`, a user can at `5173`: be `● conectat` (device `rummy_device_id` → `rummy_token`/`rummy_userId`), see `Camere disponibile` live, click `＋ Creează cameră nouă` → `rummy_matchId` persisted → `/game` shows `Talon 77` + own `Rack 14/15` vs other `RackCount`, `Current: seat-N Playing/MustDraw ←`, `DESCHIS`/`NEDESCHIS`, and can `START` → `ARUNCĂ` → `TRAGE` → `ETALEAZĂ 50+` → `EXTINDE` → `IA ULTIMA`/`RIDICĂ` → `ÎNLOCUIEȘTE JOLY` → `WIN` with `Toast bg #dc2626` on `OpServerError`, reconnect keeps `OwnRack`, and `json.Marshal(PublicView)` never contains foreign `OwnRack` IDs (Network WS frames inspected).

**Mock vs Real rule:** In `vitest` (`process.env.VITEST==='true'` → `src/lib/nakama/socket.ts:82` `isTestEnv`) `createSocket()` returns `mock-match` and `match.ts:65` RPC short-circuit `mock-match`; in `playwright`/`vite preview` (`4173`) real `Client("127.0.0.1","7350","defaultkey")` is used and `createMatch()` does `rpc(session,'create_match')` → `uuid.rummy_backend` + `joinMatch` via WS. E2E must assert `mock-match` _only_ in `vitest`; real WS tests assert `*.rummy_backend` or non-empty.

---

## Phase 0 — Smoke & env (Days 1–2)

### Day 1 — Audit frontend↔backend contract

**Goal:** Map every frontend touch point to its backend owner.

- Read `docs/protocol.md` (`Envelope{v,op,requestId,payload}`, `OpClient 1..9` `Start/Discard/DrawStock/DrawPrevious/Pickup/MeldInitial/MeldNew/Extend/ReplaceJoker`, `OpServer 100 PrivateSnapshot 101 PublicSnapshot 102 Error 103 Event`), `docs/state-machine.md` (`GamePhase Waiting/OpeningDiscard/Playing/RoundComplete`, `TurnPhase MustDraw/MeldOrDiscard`, `AllowedOps`), `client/src/lib/nakama/protocol.ts:1` `Version 1` + `NewEnvelope`, `client/src/lib/game/snapshot.ts:1` `PublicSnapshot{players{RackCount,HasOpened},stockCount,discardRow,tableMelds,winner}` / `PrivateSnapshot extends +OwnRack/OwnSeat`, `internal/match/visibility.go:36` `Version 1` `PublicView`/`PrivateView`.
- Table `frontend file | backend file | opcode | store` (e.g., `Rack.svelte:canDraw` ↔ `rummy_match.go:MatchLoop op 3` ↔ `store:privateStore` `MustDraw`+`isMyTurn`).
- Record missing gaps (e.g., `jokerReps` not yet in UI, `requestId` echo not surfaced).
- **Run:** `cat docs/protocol.md | head -n 80`, `grep -R OpClient client/src | sort`.
- **Accept:** 1-page table checked into this file Day 1 section, no code.

### Day 2 — Env & Docker smoke (proves “backend is up” before any UI)

**Goal:** Frontend `VITE_NAKAMA_*` can actually hit Nakama.

- **Tasks:** `client/.env.example` `VITE_NAKAMA_HOST/PORT/KEY/USE_SSL` → `127.0.0.1:7350 defaultkey` (not `defaulthttpkey`), `compose.yml:1` `name: rummy_backend` `postgres 5433→5432` `nakama 7350/7351/7349`; `docker compose up --build -d` → `docker compose ps` both `healthy` → `docker compose logs nakama --tail=50 | grep -E "Rummy|Found runtime modules"` expects `rummy_backend.so` + `health`; `make smoke` `pg_isready` + `nakama healthcheck` + `console 200` + `RPC health {"status":"ok"}` via `curl --user defaultkey: POST /v2/rpc/health`.
- **Files:** `client/.env*`, `compose.yml`, `Dockerfile` (`nakama-pluginbuilder:3.26.0` → `nakama:3.26.0`, plugin baked `/nakama/data/modules/rummy_backend.so`), `scripts/smoke.sh`, `Makefile`.
- **Run:** `docker compose up --build -d && make smoke`.
- **Accept:** `SMOKE PASSED` with `rummy_backend:local` image, `7351` console `admin/password` reachable, `7350` `health` works with `defaultkey:`.

---

## Phase 1 — Auth & socket (Days 3–4)

### Day 3 — Device auth & persistence (the only way to get a `Session`)

**Goal:** `authenticateDevice` → `localStorage` survives reload.

- **Tasks:** `client/src/lib/nakama/client.ts:1` `new Client(KEY,HOST,PORT,USE_SSL)`, `client.ts:getOrCreateDeviceId` `rummy_device_id` ↔ `client/src/lib/nakama/auth.ts:7` `authenticate()` `rummy_device_id→rummy_token/rummy_userId` + `authStore: writable<Session|null>`; `client/src/lib/nakama/auth.test.ts` `rummy_device_id` persists; `client/tests/sync.e2e.ts` not yet — manual: open `5173` → `● conectat` + `Application→Local Storage` shows `rummy_device_id`/`rummy_token`/`rummy_userId`; reload still `● conectat` without re-login.
- **Run:** `npm run test:unit -- --run src/lib/nakama/auth.test.ts && npm run test:e2e -- tests/sync.e2e.ts` (expect `mock-match` in unit, real `defaultkey:` in e2e).
- **Accept:** `localStorage` has `rummy_device_id` 36-char UUID, `token` non-empty, `Network → POST /v2/account/authenticate/device?create=true` `200` with `Authorization: Bearer …` for later `rpc`/`listMatches`.

### Day 4 — Socket & `onmatchdata → store` wiring

**Goal:** WS frames actually update Svelte stores.

- **Tasks:** `client/src/lib/nakama/socket.ts:13` `socketStore`, `setMatchDataHandler`, `createSocket(session)` `client.createSocket(false,false)` `sock.connect(session,true)` + `sock.onmatchdata → _matchDataHandler` + `ondisconnect`; `client/src/lib/game/store.ts:14` `privateStore`/`publicStore` `isMyTurn` derived, `handleMatchData(opCode, rawData)` handles `OpServerState 100 → onPrivateSnapshot`, `101 → onPublicSnapshot`, `102 → onServerError` (Day 45), `initGameStore()` auto-wire in browser; `socket.test.ts` mock `onmatchdata` forwarding.
- **Verify connected:** Open DevTools `WS` → `wss://127.0.0.1:7350/ws...` frames: after `joinMatch` expect `op 100` JSON `{"v":1,"gamePhase":"Waiting",...,"ownRack":...}` swallowed into `privateStore`, `op 101` for opponents; `console` no `onmatchdata is null`.
- **Run:** `npm run test:unit -- --run src/lib/nakama/socket.test.ts && npm run test:unit -- --run src/lib/game/store.test.ts`.
- **Accept:** `privateStore` set from `100`, `publicStore` from `101`, `isMyTurn = currentSeat===ownSeat` true only on own turn, `lastPrivate` mirrored, no throw on `Uint8Array` vs `string`.

---

## Phase 2 — Lobby & room (Days 5–6)

### Day 5 — Lobby: `create`/`list`/`join` rooms

**Goal:** 2 browsers can form a room without typing IDs.

- **Tasks:** `client/src/lib/nakama/match.ts:64` `createMatch()` real path `rpc(session,'create_match')` → `matchId` `uuid.rummy_backend` → `sock.joinMatch(matchId)` + `persistMatchId` `rummy_matchId`; fallback `sock.createMatch()` mock; `listAvailableMatches()` `listMatches(session,10,…)` deduped `AvailableMatch{matchId,label,size}`; `client/src/routes/+page.svelte:29` `lobby` `＋ Creează cameră nouă` → `goto('/game')`, `Camere disponibile` `flip`/`fade` with `reîmprospătează` (no `interval` hard reload), `availableMatches` filtered `size 1..4`; `+page.svelte:116` auto `goto('/game')` when `priv||pub`; `game/+page.svelte` shows `Se conectează la masa …` until `100` then `Lobby — Așteptare`.
- **E2E:** `tests/start.e2e.ts` `Waiting Host 2p` → `START` visible, guest not, `1p` not; manual 2 browsers `alice: Creează` → `bob: Intră` same `matchId.slice(0,8)`; `Application→localStorage rummy_matchId` same on both.
- **Run:** `npm run test:e2e -- tests/start.e2e.ts && npm run test:unit -- --run src/lib/nakama/match.test.ts` (unit expects `mock-match`, e2e expects real).
- **Accept:** `lobby` shows `CAMERE DISPONIBILE` auto without `location.reload`, `create` → `rummy_matchId` persisted, `bob` sees `1/4` then `2/4` after join, both see `Waiting`.

### Day 6 — Start & opening discard (first server-authoritative move)

**Goal:** `START` flips both clients to `OpeningDiscard 15→14`.

- **Tasks:** `client/src/components/TopBar.svelte:11` `canStart = Waiting && ownSeat==0 && players>=2` → `sendStart()` `NewEnvelope(1,{},requestId)` `sendMatchState(matchId,1,env)`; `client/src/components/Rack.svelte:46` `isOpeningDiscard && isMyTurn && rackCount==15` → `canDiscard` + `ARUNCĂ CARTEA`; `client/src/lib/game/actions.ts:29` `sendDiscard(tileId)` `NewEnvelope(2,{tileId})`; `internal/match/opening_discard_test.go` asserts `IsOpeningDiscard` flag.
- **E2E:** `tests/opening.e2e.ts` `OpeningDiscard MyTurn 15→14` `discard-btn` enabled only for opener; `game/+page.svelte` after `START` shows `Faza: OpeningDiscard` `Aruncă 1 din 15` + `Talon 77` (106-15-14-0? Actually `stock 77` after `15+14` deal) + `TableBoard` empty; WS `op 101` shows `DiscardRow[0].IsOpeningDiscard true`; next `Current: seat-1` `MustDraw`.
- **Run:** `npm run test:e2e -- tests/opening.e2e.ts`.
- **Accept:** Non-opener `discard` → `OpServerError not_your_turn` `Toast`, opener `15→14` `stock 77` unchanged, `IsOpeningDiscard` never pickable (`draw-prev` stays disabled).

---

## Phase 3 — Turn loop (Days 7–8)

### Day 7 — Draw stock & normal discard + turn rotation

**Goal:** `MustDraw → MeldOrDiscard → Discard → next seat` is server-driven.

- **Tasks:** `Rack.svelte:74 canDraw = Playing && MustDraw && isMyTurn` → `▶ TRAGE DIN TALON` `OpClientDrawStock 3 {}`; `canDiscard = (OpeningDiscard 15→14) || (Playing MeldOrDiscard && isMyTurn)` → `ARUNCĂ CARTEA` `OpClientDiscard 2 {tileId}`; `TopBar Day43` `Current: seat-N Playing/MustDraw ←` (emerald `← rândul tău` when `isMyTurn` else muted `← seat-N`); `game/+page.svelte` `StockCount`/`DiscardRow`/`CurrentSeat`.
- **E2E:** `tests/draw.e2e.ts` `MustDraw MyTurn` → `draw-btn` enabled else disabled (`MustDraw NotMyTurn`, `MeldOrDiscard`); `tests/discard.e2e.ts` `MeldOrDiscard Discard 0→1`; `tests/turn.e2e.ts` `Current: seat-0 Playing/MustDraw ← rândul tău` vs `← seat-1` + `draw-btn` disabled when not `MustDraw`/`CurrentSeat`.
- **Run:** `npm run test:e2e -- tests/draw.e2e.ts tests/discard.e2e.ts tests/turn.e2e.ts`.
- **Accept:** `draw` pops `Stock 70→69` `Rack 14→15` `MustDraw→MeldOrDiscard`; `discard` appends `DiscardRow` ordered `Index` + advances `CurrentSeat=(cur+1)%n` `MeldOrDiscard→MustDraw` for next seat.

### Day 8 — HasOpened gating & DrawPrevious gated by `HasOpened`

**Goal:** Frontend respects `PublicPlayer.hasOpened`.

- **Tasks:** `Rack.svelte:81 hasOpened = players.find(p=>p.seat==ownSeat).hasOpened` → badge `DESCHIS ✓` emerald vs `NEDESCHIS • 50+ RUN` amber + helper; `canDrawPrevious = Playing MustDraw isMyTurn hasOpened last !IsOpeningDiscard`, `canPickup = … hasOpened selected 2 + discardIndex !IsOpeningDiscard`, `canExtend/canReplace` require `hasOpened && MeldOrDiscard && isMyTurn`; `tests/hasopened.e2e.ts` `NEDESCHIS` → `IA ULTIMA/RIDICĂ/EXTINDE/ÎNLOCUIEȘTE` disabled.
- **Run:** `npm run test:e2e -- tests/hasopened.e2e.ts && npm run test:e2e -- tests/draw-prev.e2e.ts`.
- **Accept:** `!HasOpened` never shows `IA ULTIMA` enabled even with discard; `HasOpened` false helper text visible.

---

## Phase 4 — Melds (Days 9–10)

### Day 9 — MeldInitial (`!HasOpened 50+ run`) vs MeldNew (`HasOpened`)

**Goal:** Frontend sends the right opcode `6` vs `7` and server decides.

- **Tasks:** `Rack.svelte:178 canMeldInitial = !hasOpened && Playing MeldOrDiscard isMyTurn selected>=3`, `canMeldNew = hasOpened && …`; `meldInitial()` `hasOpened? sendMeldNew([{kind:'run',tileIds}]) : sendMeldInitial([...])`; `actions.ts:188 sendMeldInitial/SendMeldNew` payload `{melds:[{id,kind,tileIds,jokerReps:{}}]}` `NewEnvelope(6|7, payload, requestId)`; `internal/match/meld_*_test.go` scoring `49` rejected.
- **E2E:** `tests/meld-initial.e2e.ts` `!HasOpened selected 3 → ETALEAZĂ SELECTATE` `OpClientMeldInitial 6` `mock-match`; `tests/meld-new.e2e.ts` `HasOpened` → `7`; `demo/meld` `TotalScore` preview (Day 47 future).
- **Run:** `npm run test:e2e -- tests/meld-initial.e2e.ts tests/meld-new.e2e.ts`.
- **Accept:** `!HasOpened` never sends `7`, `HasOpened` never sends `6`; `49` → `Toast bad_payload` `code` + `requestId` echo, `Racks` unchanged (atomic).

### Day 10 — Extend & ReplaceJoker (via `TableBoard` click)

**Goal:** Extending others’ melds and joker recovery use the same path.

- **Tasks:** `TableBoard.svelte:83 pickupDiscardIndex/selectedMeldId/replaceTargetMeldId` `onMeldClicked→selectedMeldId`, `onJokerClicked→replaceTargetMeldId`; `Rack.svelte:220 canExtend = hasOpened MeldOrDiscard isMyTurn selected>=1 && meld exists`, `canReplace = hasOpened MeldOrDiscard selected==3 && meld.Tiles.some(joker)`; `actions:sendExtendMeld(8,{meldId,tileIds})`, `sendReplaceJoker(9,{targetMeldId,tileId,newMeldTiles:[2]})`.
- **E2E:** `tests/extend.e2e.ts` `HasOpened selected>=1 + meldId → EXTINDE 8`; `tests/joker.e2e.ts` `joker click + Rack 3 tiles → ÎNLOCUIEȘTE 9`; `TableBoard` highlight `selectedMeldId` `border-sky-500`.
- **Run:** `npm run test:e2e -- tests/extend.e2e.ts tests/joker.e2e.ts`.
- **Accept:** Clicking `meld-m1` sets `selectedMeldId`, joker tile `joker-m1-*` sets `replaceTargetMeldId`, both require `HasOpened` else disabled.

---

## Phase 5 — Discard & win (Days 11–13)

### Day 11 — Discard row pickup (previous + sweep)

**Goal:** `Pickup` moves `DiscardRow[discardIndex+1:]` to rack.

- **Tasks:** `TableBoard.svelte:191 discardRow` buttons `discard-tile-${idx}` `border-sky-500` when `pickupDiscardIndex==Index`, `IsOpeningDiscard` `border-red-300`; `Rack.svelte:146 canPickup = MustDraw HasOpened selected 2 + index in range !IsOpeningDiscard`; `actions:sendPickupDiscardForMeld(5,{discardIndex,tileIds:[2]})` sweep later tiles server-side.
- **E2E:** `tests/pickup.e2e.ts` `selected 2 + discardIndex → RIDICĂ 5` `mock-match`; `game` after `PICKUP` `RackCount` increases by `laterCount+1` (discard + swept) and `DiscardRow` truncated `[:discardIndex]`.
- **Run:** `npm run test:e2e -- tests/pickup.e2e.ts`.
- **Accept:** Clicking opening discard never enables `RIDICĂ`, latest only via `IA ULTIMA` path vs earlier via `RIDICĂ`.

### Day 12 — Win / lose `RoundComplete`

**Goal:** Empty rack (`Rack 0`) → `Winner` overlay everywhere, no further moves.

- **Tasks:** `internal/match/win.go` `checkWinAndComplete` `rack==0` after `DISCARD` or melding without discard (`docs/rules-decisions.md:6.1`); `client/src/components/WinnerOverlay.svelte` `GamePhase=="RoundComplete"` `Winner 0` `RESTART MASA` → `localStorage.removeItem('rummy_matchId')` + `resetGame()`; `game/+page.svelte` `StockCount` still `PublicView` winner visible; `tests/winner.e2e.ts` `RoundComplete Winner 0`.
- **Run:** `npm run test:e2e -- tests/winner.e2e.ts && go test ./internal/match -run TestWin -v`.
- **Accept:** After `win` any `DISCARD`/`MELD` → `OpServerError wrong_phase` `Toast`, `WinnerOverlay` shows on both browsers same `Winner`.

### Day 13 — Errors → Toast (`OpServerError 102`)

**Goal:** Every `102` is user-visible, 3s, `bg #dc2626`, `data-error-code`.

- **Tasks:** `client/src/lib/game/errorStore.ts:15` `onServerError(raw)` parses `{code,message,details,requestId,op}` `errorStore` 3s auto-clear; `src/components/Toast.svelte:5` `errorStore` `role=alert` `data-error-code` `data-request-id` `data-op` `fly`/`fade` `bg-[#dc2626]`; `store.ts:131` `handleMatchData(102)` → `onServerError`; `+layout.svelte:3` global `<Toast/>`; `demo/error` 5 triggers (`bad_payload not_your_turn bad_request LEAKED with details`) + `Clear`.
- **E2E:** `tests/error.e2e.ts` `Trigger bad_payload → toast visible data-error-code=bad_payload bg rgb(220,38,38) → LEAKED → bad_request → Clear → auto-dismiss 3.5s` ✓ (30 tests).
- **Run:** `npm run test:e2e -- tests/error.e2e.ts`.
- **Accept:** `requestId`/`op` echoed, `LEAKED` not suppressed, no private `TileId` in payload.

---

## Phase 6 — Resilience & parity (Days 14–17)

### Day 14 — Visibility / no-leak (the only security gate that matters)

**Goal:** `PublicSnapshot` never carries `OwnRack`.

- **Tasks:** `client/src/lib/game/snapshot.test.ts` `isValidPublic/Private`, `client/src/lib/game/redaction.test.ts` `checkNoLeak(publicJson, privateIds)` `JSON.stringify(PublicSnapshot)` search; `src/lib/game/store.test.ts` `PrivateSnapshot ownRack 3 only local`; `internal/setup/redaction_test.go` exhaustive `2,3,4×seed`; `client/tests/sync.e2e.ts` `alice vs bob different Rack`.
- **Run:** `npm run test:unit -- --run && go test ./internal/setup -run TestRedaction -v`.
- **Accept:** `Public JSON` contains `stockCount`/`RackCount`/`DiscardRow`/`TableMelds`/`Winner` only, `Private 100` contains `OwnRack` for that `Seat` only; WS `Network` frames inspected manually show no foreign `ID`.

### Day 15 — Reconnection (keep `Seat`+`Rack`)

**Goal:** Refresh or briefly disconnect → still your rack.

- **Tasks:** `client/src/lib/nakama/reconnect.ts:9` `reconnect()` `getStoredMatchId() ?? getMatchId()` → `createSocket()` → `joinMatch(storedId)` expects `OpServerState 100` new `OwnRack` not old; `client/src/lib/game/store.ts:20 privateBySeat Map + localStorage rummy_lastPrivate:${seat}`; `+page.svelte:52` `onMount authenticate() → reconnect()` + `Se conectează` / `Înapoi la Lobby`; `client/tests` `reconnect.test.ts` `new-1 not old-1`.
- **E2E:** Manual: create 2p, `F5` on `bob` → `Se conectează la masa …` → `Rack` same count, `HasOpened` preserved; WS `ondisconnect` clears `socketStore` but not `rummy_matchId`.
- **Run:** `npm run test:unit -- --run src/lib/nakama/reconnect.test.ts`.
- **Accept:** `MatchLeave` keeps `Players`/`Racks` (`6f2af5b`), `MatchJoinAttempt` allows existing in `Playing`, `localStorage rummy_lastPrivate:0` updated.

### Day 16 — UI feedback that _must_ match server (Turn + HasOpened)

**Goal:** What the user sees (disabled buttons) matches what server will accept.

- **Tasks:** `TopBar Day43` `Current: seat-N Playing/MustDraw ← rândul tău|← seat-N` `aria-label`; `Rack Day44` `DESCHIS/NEDESCHIS` + helper `50+ RUN`; `tests/turn.e2e.ts` + `tests/hasopened.e2e.ts` + future `validate/scoring preview` (Day 46+).
- **Run:** `npm run test:e2e -- tests/turn.e2e.ts tests/hasopened.e2e.ts`.
- **Accept:** `draw-btn` disabled when `!MustDraw`/`NotMyTurn` matches `wrong_phase`/`not_your_turn` server; `Prev/Pickup/Extend/Replace` disabled when `!HasOpened` matches `not_opened`.

### Day 17 — 2-browser real WS flows (the “does it actually work” proof)

**Goal:** `alice` and `bob` play to win through real Nakama via Svelte.

- **Tasks:** `client/playwright.config.ts:3` `webServer: npm run build && npm run preview port 4173`; `client/tests/e2e-actions.e2e.ts` 2 `BrowserContext` `alice createMatch + bob joinMatch same matchId` `alice Start → OpeningDiscard → DrawStock 3 → MeldInitial 6 → MeldNew 7 → Extend 8 → DrawPrevious 4/Pickup 5 → Replace 9 → Discard 2 → win`; `client/tests/visual-actions.e2e.ts` screenshot `TopBar+TableBoard+Rack`; `internal/match/deterministic_simulation_test.go` mirrored in Go with `CheckTileConservation`.
- **Run:** `npm run test:e2e -- tests/e2e-actions.e2e.ts tests/visual-actions.e2e.ts` (both `2 browsers`, `requestId`).
- **Accept:** Both see same `TableMelds`/`DiscardRow`/`CurrentSeat`, `Rack` distinct, `Winner` same, `lastSent requestId` present in `OpServerError`.

---

## Phase 7 — Manual & CI (Days 18–20)

### Day 18 — Manual QA with 2 laptops + mobile

**Goal:** No `interval` hard reload, `flip`/`fade` only reîmprospătează.

- Checklist: `5173` lobby `● conectat` → `＋ Creează` → `Camere disponibile` appears without `F5` on 2nd laptop → `Intră` → `/game` `Lobby — Așteptare 2/4` → host `START` → `Talon 77` → opener `ARUNCĂ` `15→14` → other `TRAGE` → `ETALEAZĂ 50+ run` → `EXTINDE` others’ meld → `IA ULTIMA` → `RIDICĂ` sweep → `ÎNLOCUIEȘTE` → `ARUNCĂ` → `WIN` `RESTART` → `F5` reconnect → `Toast` on `DISCARD not in rack`.
- **Accept:** No full `location.reload`, `Network→WS` shows `op 1..9` with `requestId`, `Preserves log` shows `OpServerError` toast 3s.

### Day 19 — CI gate (no oral knowledge)

**Goal:** One command proves everything.

- **Tasks:** `.github/workflows/client.yml` `eslint` `prettier --check` `svelte-check` `vitest --run` `vite build` `playwright` (install `chromium`); root `.github/workflows/ci.yml` `go vet` `gofmt -l` `go test` `go mod tidy` + `docker/build`; `client/README.md:6` `npm run check && npm run test:unit -- --run && npm run build && npm run test:e2e`.
- **Run:** `make check && npm run check && npm run test:unit -- --run && npm run build && npm run test:e2e` (same as CI).
- **Accept:** All 3 workflows green on `main`, `vite build` no `svelte-check` errors.

### Day 20 — Tag & “what’s still deferred”

**Goal:** Mark verified, document deferred.

- **Tasks:** `git tag svelte-nakama-verified` when Days 1–19 green + `go test ./... -race` green + `docker compose build` green; update `client/README.md:6` `Next Steps` to `Day 46 — Meld validation` etc. per `client/docs/roadmap.md:97`; list deferred: `jokerReps` dropdown, scoring preview `5/10/25`, duplicate guard, `Rack` sorting, `prefers-reduced-motion` (Phase 6).
- **Accept:** `git log --oneline --grep="Day 45"` shows `5959151` etc., this file checked, `docs/rules-decisions.md` §Deferred up to date.

---

## Appendix — How to prove “frontend is connected” (diagnostics)

```bash
# 1. Backend up?
docker compose up --build -d && make smoke # expects SMOKE PASSED + rummy_backend.so

# 2. Env correct?
cat client/.env.example | grep VITE_NAKAMA # HOST 127.0.0.1 PORT 7350 KEY defaultkey (not defaulthttpkey)

# 3. Auth?
open http://localhost:5173 → ● conectat + Application→Local Storage rummy_* + Network 200 POST /v2/account/authenticate/device?create=true --user defaultkey:

# 4. Create/join?
# alice click ＋ Creează → Network WS wss://127.0.0.1:7350/ws... op 100 {"v":1,"gamePhase":"Waiting","ownRack":[]} + localStorage rummy_matchId
# bob click Intră same 8-char → both see MASA 1 • 2 JUCĂTORI → host START → both see Talon 77

# 5. Store?
# DevTools→WS frames: op 1 Start, op 2 Discard {tileId}, op 3 DrawStock {}, op 6 MeldInitial {melds:[{kind:run,tileIds}]}} with requestId, responses 100/101/102

# 6. Unit vs Real?
npm run test:unit -- --run # expects mock-match (isTestEnv)
npm run test:e2e -- tests/e2e-actions.e2e.ts # expects *.rummy_backend uuid via real rpc

# 7. Full gate?
make check && npm run check && npm run test:unit -- --run && npm run build && npm run test:e2e # 30 tests (Day45)
```

**Invariant every check:** `docker compose logs nakama --tail=100 | grep -E "MatchLoop op=|Sent error|Rummy"` shows `MatchLoop op=6` etc., no `OwnRack` in `PublicView` JSON (see `visibility.go:36`).

---

_Refs: `client/src/lib/nakama/{auth,client,socket,match,protocol}.ts`, `client/src/lib/game/{snapshot,store,errorStore}.ts`, `client/src/components/{Rack,TableBoard,TopBar,Toast}.svelte`, `client/src/routes/demo/*` (`turn/hasopened/error`), `client/tests/*.e2e.ts` (`turn/hasopened/error/visual-actions/e2e-actions/winner`), `internal/match/*_test.go` `CheckTileConservation 106`, `docs/protocol.md` `Version 1` `1..9`/`100..103`, `docs/state-machine.md` `AllowedOps`, `AGENTS.md` 24-day plan, `compose.yml:1` `rummy_backend` `5433`._
