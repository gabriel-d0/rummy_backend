# SvelteKit Frontend — Vertical-Slice Implementation Plan (Handmade Hero, Technically Correct)

**Target:** A small, maintainable, playable Romanian Tile Rummy web client in **SvelteKit 2 + Svelte 5 (runes) + Vite 6 + Tailwind 4 + TypeScript 5 + `@heroiclabs/nakama-js 2.8`** that is _thin and correct_: it renders only `PrivateSnapshot`/`PublicSnapshot` from the authoritative Go/Nakama backend (`docs/protocol.md` `Version 1`, `internal/protocol/opcodes.go:8` `1..9`/`100..103`, `internal/match/visibility.go:36`) and sends `1..9` via `sendMatchState`. No UI guesses — server decides; UI only enables what `AllowedOps` + `HasOpened` + `isMyTurn` allow.

**Backend prerequisite (frozen):** `rummy-mvp-rc1` at `main` (`go 1.23.5` + `nakama:3.26.0` + `postgres:15` `docker compose` `7350`/`7351` `defaultkey`), `make smoke` `SMOKE PASSED`. Tiles 106 (`TileInstance{ID,Colour 1..4,Rank 1..13,IsJoker}` unique ID), `CheckTileConservation 106` invariant.

**Why this order is “logical, natural, technically correct”:**

1. Design system before layout — `tokens` single source, no `px` hardcode, avoids churn.
2. Layout before state — `TopBar/TableBoard/Rack` exist as pure props components, so `stores` can mount without re-plumbing.
3. State (`snapshot` → `store` → `socket`) before actions — `onmatchdata` must land in `privateStore/publicStore` before any `send*` is wired, otherwise UI flakes.
4. One opcode per day, in turn-order (`Start → OpeningDiscard → DrawStock → DrawPrevious → Pickup → Discard → MeldInitial → MeldNew → Extend → ReplaceJoker → Winner`) — each day’s `AllowedOps` builds on yesterday’s `TurnPhase`, so regression is local.
5. Feedback (Turn/HasOpened/Toast/validate/scoring) _after_ actions — gives non-authoritative preview without duplicating server math; server remains source of truth (`OpServerError 102`).

**Daily loop (Handmade Hero, one commit):** `UI slice` → `Nakama wire (protocol + store + actions)` → `Playwright slice` → `npm run check (svelte-check + eslint + prettier --check + vitest --run + vite build) + npm run test:e2e green` → `git add <specific> && git commit -m "feat(client): ... (Day N)" && git push`.

**Stack anchors:** `Svelte 5 runes` (`$props`, `$state`, `$derived`, `$effect`), `svelte/store` (`writable`, `derived`, `get`), `svelte/transition|animate` (`fly`, `fade`, `flip`), `Tailwind 4` via `@tailwindcss/vite` + `@import "tailwindcss"` in `app.css`, `adapter-auto`, `svelte-check --tsconfig ./tsconfig.json`, `vitest` `projects.server` `environment: node` for `src/**/*.{test,spec}.{js,ts}`, `playwright` `webServer: npm run build && npm run preview` `port 4173` `testMatch **/*.e2e.{ts,js}`, `prettier-plugin-svelte` + `prettier-plugin-tailwindcss` + `eslint-plugin-svelte`.

---

## Phase 0 — Foundation (Days 1–8) — no game logic yet

### Day 1 — Repo audit

**Do:** Read `rummy_backend` root, `AGENTS.md` 24-day plan, `docs/protocol.md`/`state-machine.md`/`rules-decisions.md`, `compose.yml` (`postgres:15-alpine 5433→5432`, `nakama:3.26.0` `7350/7351`), `internal/match/visibility.go` (`PublicView`/`PrivateView` `Version 1`), `client/` prior `rm -rf` state. Record in `client/docs/roadmap.md` + `client/README.md` stub + `.env.example` `VITE_NAKAMA_HOST/PORT/KEY/USE_SSL`.

**UI:** none.
**Nakama:** none.
**Test:** `vitest` zero test placeholder, `playwright` `page.goto("/")` `200` smoke.
**Done:** Audit file committed, everyone knows what `Remi Etalat` means vs brand.

### Day 2 — SvelteKit skeleton (prove toolchain)

**UI:** `npx sv create client --template minimal --types ts --add prettier eslint vitest playwright tailwindcss --install npm` → `src/routes/+page.svelte` `Hello SvelteKit`, `src/app.html`, `src/routes/layout.css` empty, `vite.config.ts` (`tailwindcss()`, `sveltekit({compilerOptions:{runes:…}})`, `adapter: adapter()`), `svelte.config.js`, `tsconfig.json` `strict`.
**Nakama:** none.
**Test:** `npm run dev -- --open` shows `SvelteKit` at `5173`; `tests/scaffold.e2e.ts` `getByText("SvelteKit")`; `npm run check` green.
**Files:** `package.json:scripts dev/check/lint/test:unit/test:e2e/build`.
**Pitfall:** `adapter-auto` needs `vite preview` `4173` for `playwright` `webServer` — not `dev`.

### Day 3 — Tailwind & design tokens

**UI:** `tailwind.config` is _v4_ — no `tailwind.config.js`; add `app.css` `@import "tailwindcss"` + `theme` `colors.felt #0a2e1a`, `wood #5d4037`, `fonts Inter/JetBrains Mono` via `googleapis` import, `layout.css` imports `app.css`; `src/lib/ui/tokens.ts` single source `colors:{felt,wood,tile{red #e53935,yellow #f9a825,blue #1e88e5,black #212121},feltLight}`, `radius 12/16/18`, `spacing 4/8/12/16`; `+page.svelte` `bg-[#0a2e1a]` `max-w-[1600px] mx-auto` (Handmade Hero `Scale.FIT`).
**Nakama:** none.
**Test:** `vitest src/lib/ui/tokens.test.ts` `felt #0a2e1a`.
**Done:** No per-component `px` yet; tokens drive all.

### Day 4 — Nakama JS client stub (no UI, just `Client`)

**UI:** none.
**Nakama:** `npm i @heroiclabs/nakama-js@2.8.0`, `src/lib/nakama/client.ts` `Client(KEY,HOST,PORT,USE_SSL)` via `$env/static/public` (`VITE_NAKAMA_*`), `getClient()` singleton, `authenticateDevice(deviceId,true,username)` + `getOrCreateDeviceId()` `rummy_device_id` → `localStorage`, lazy `createSocket`; `src/lib/config.ts` re-export `VITE_NAKAMA_*` with `127.0.0.1:7350 defaultkey` fallback (never `defaulthttpkey`).
**Test:** `src/lib/nakama/client.test.ts` `getClient()` exists, `vitest` mock `localStorage`; `vite build` green (no `Socket` connect yet).
**Accept:** `Client` can be instantiated in `node` test without `window`.

### Day 5 — Project structure & lint strict

**UI:** Create `src/lib/nakama/`, `src/lib/game/`, `src/lib/ui/`, `src/components/{Tile,TableBoard,Rack,TopBar}/` empty shells, `src/routes/demo/` empty; `eslint.config.js` (`@eslint/js` + `typescript-eslint` + `eslint-plugin-svelte` + `eslint-config-prettier`), `.prettierignore`, `svelte-check` `strict` (no `any` without comment).
**Nakama:** none.
**Test:** `npm run lint` (`prettier --check . && eslint .`) pass.
**Done:** Every future file has a home; `lint:check` is CI.

### Day 6 — Env wiring & docs

**UI:** none.
**Nakama:** `src/lib/config.ts` typed `VITE_NAKAMA_*`, `.env.example` `VITE_NAKAMA_HOST=127.0.0.1 … KEY=defaultkey`, `client/README.md` `Quick Start` (`docker compose up --build -d && make smoke` → `cd client && npm i && cp .env.example .env && npm run dev -- --open` at `5173`), root `README.md` `Next Steps` points to `client/`.
**Test:** `tests/docs.e2e.ts` `README` contains `VITE_NAKAMA_HOST`, `.env.example` has 4 keys.
**Done:** New clone → `npm run dev` without oral knowledge.

### Day 7 — Smoke test (prove `5173` and Nakama reachable)

**UI:** `+page.svelte` still `Hello` but `src/lib/nakama/client.ts` `authenticateDevice` reachable.
**Nakama:** `tests/smoke.e2e.ts` `page.goto("/")` `200`, `playwright.config.ts` `webServer: npm run build && npm run preview` `port 4173` ensures `build` before `e2e`; manual `curl -s -X POST http://127.0.0.1:7350/v2/account/authenticate/device?create=true --user defaultkey: -d '{"id":"..."}' | jq .token`.
**Test:** `npm run test:e2e` `smoke` green, no console errors.
**Done:** `npm run dev` + `docker` are not two worlds.

### Day 8 — CI baseline

**UI:** none.
**Nakama:** `.github/workflows/client.yml` `eslint` `prettier --check` `svelte-check` `vitest --run` `vite build` `playwright` (install `chromium` + `webServer`).
**Test:** CI green on push to `main`.
**Milestone:** Any commit that breaks `check` is red.

---

## Phase 1 — Design system & layout (Days 9–16) — pure props, no `store`

### Day 9 — Design tokens

**UI:** `src/lib/ui/tokens.ts` finalized `export const colors`, `fonts`, `spacing`, `radius`, `shadows`; `vitest` asserts `colors.felt === "#0a2e1a"`; usage `bg-[colors.felt]` via Tailwind arbitrary.
**Why first:** Tokens avoid per-component `hex` drift later (`tile red` must be `#e53935` everywhere, `WCAG AA` on `felt`).

### Day 10 — Layout — GameSpace `1600×900`

**UI:** `src/lib/layout.ts` `GAME_SPACE 1600×900`, `getLayout()` `TopBar 56`, `TableArea`, `PlayerRackArea 132`, `ActionBar 56`; `+page.svelte` `max-w-[1600px] mx-auto` `flex-col lg:flex-row` (no JS scale, `flex-wrap` handles overflow).
**Test:** `vitest` `getLayout().PlayerRackArea.height 132`.

### Day 11 — Tile

**UI:** `src/components/Tile.svelte` `$props<{colour 1..4, rank 1..13, isJoker?, size: 'rack'|'table', selected?, draggable?}>` → `sizeClass` `rack 64×90 vs table 48×64`, `colourMap` `1:red #e53935 border-red-500`, `rankLabel` `1→A 11→J 12→Q 13→K` else `String(r)`, `isJoker` `amber-50 border-amber-500` `JOLY`, `selected` `ring-2 sky-500 scale-1.06 ✓`, `draggable` `cursor-grab`; `tileContent` snippet.
**Test:** `vitest src/lib/ui/tile.test.ts` `red #e53935`, `playwright tests/tile.e2e.ts` `A/K/Joly`.
**Done:** `Table` uses `table` `52×72`, `Rack` uses `rack` `64×90` — auto-scale via `flex-wrap`, no `overflow-x`.

### Day 12 — Tile modern (polish)

**UI:** White `rounded-lg border` `shadow-sm hover:-translate-y-0.5 hover:shadow-md`, `Joly` centered `J` `16px` + `JOKER` `6px`; selected adds `shadow-lg`.
**Test:** `vitest selected ring-2`.

### Day 13 — TableBoard

**UI:** `src/components/TableBoard.svelte` `$props<{melds?: Meld[]}>` default 3 demo melds (`66 pct` `53 pct` `55 pct`), `bg-[#f5f1e8] rounded-[18px] border-[#e8e0c8] shadow-inner`, `ETALĂRI PE MASĂ • n SETURI` `flex-wrap gap-2.5` **no scroll** (`auto-scale` if overflow via `flex-wrap` + `max-w`), `melds.slice(0,2)` + `slice(2)` rows, `points` pill `slate-900`.
**Test:** `playwright tests/table.e2e.ts` `ETALĂRI PE MASĂ` `no scroll` (`scrollWidth<=clientWidth`).
**Why before `Rack`:** `Table` is public — no `store` leak risk.

### Day 14 — Rack

**UI:** `src/components/Rack.svelte` `$props<{tiles?: RackTile[11]}>` 11 demo tiles, `Mâna ta • n cărți` + `TREBUIE SĂ TRAGI` pill `amber-400`, `SORTEAZĂ CULOARE/NUMĂR` pills `white/10`, `flex-wrap gap-1.5 sm:gap-2` `min-h-[110px]` no scroll, `TRAGE DIN TALON` `▶` button placeholder.
**Test:** `playwright tests/rack.e2e.ts` `Mâna ta 11 cărți` `no scroll`.
**Done:** `Table`+`Rack` look like `remi-online.ro` but without its `overflow-x`.

### Day 15 — TopBar

**UI:** `src/components/TopBar.svelte` `$props<{players 4, masa 1, seconds 0}>` `h-12 bg-black/90 border-white/10`, `R` `amber-400` `REMI ETALAT ETALAT` `PREMIUM • ONLINE`, `MASA n • x JUCĂTORI` `pulse emerald`, `🕒 0s`, `REGULI` `JOC NOU`.
**Test:** `playwright tests/topbar.e2e.ts` `REMI ETALAT` `MASA 1`.

### Day 16 — Visual layout (responsive snapshot)

**UI:** `+page.svelte` wires `TopBar` + `TableBoard` + `Rack` + `Jurnal` placeholder `hidden lg:flex w-[300px]` `JURNAL DE JOC`, `max-w-[1600px] mx-auto p-3 lg:flex-row`.
**Test:** `playwright tests/layout.e2e.ts` `1280×800` + `375×667` all visible, no overlap, `TopBar` `56px`, `Rack` `132px` via `layout.ts`.
**Milestone:** Static game table looks right on desktop+mobile without `store`.

---

## Phase 2 — State & networking — Svelte stores (Days 17–28) — mount `snapshot` before `actions`

### Day 17 — Snapshot types

**UI:** none (type slice).
**Nakama:** `src/lib/game/snapshot.ts` `SnapshotVersion 1`, `TileInstance{ID,Colour,Rank,IsJoker}`, `DiscardEntry{Tile,IsOpeningDiscard,Index}`, `TableMeld{ID,Kind,Tiles,JokerReps,OwnerSeat}`, `PublicPlayer{id,seat,hasOpened,rackCount}` `PublicSnapshot{v,gamePhase,turnPhase,currentSeat,players,stockCount,discardRow,tableMelds,winner}` `PrivateSnapshot extends +ownRack,ownSeat`, `isValidPublic/Private` (checks `v`, `gamePhase`, `currentSeat`, `players[]`, `discardRow[]`, `tableMelds[]`).
**Test:** `vitest src/lib/game/snapshot.test.ts` `SnapshotVersion 1` `isValidPrivate false on bad`.

### Day 18 — Nakama envelope

**Nakama:** `src/lib/nakama/protocol.ts` `Version 1`, `OpClientStart 1 … ReplaceJoker 9`, `OpServerState 100 StatePublic 101 Error 102 Event 103`, `type Envelope{v,op,requestId,payload}`, `NewEnvelope(op,payload,requestId?)` `JSON.stringify({v,op,requestId,payload})`, `parseEnvelope` throws `bad envelope`.
**Test:** `vitest src/lib/nakama/protocol.test.ts` `NewEnvelope(1,{}) v:1`.

### Day 19 — Auth store

**Nakama:** `src/lib/nakama/auth.ts` `authStore writable<Session|null>`, `getOrCreateDeviceId()` `rummy_device_id` → `localStorage` `crypto.randomUUID` fallback, `authenticate(username?)` `Client.authenticateDevice(deviceId,true,username)` → `localStorage rummy_token/rummy_userId` + `authStore.set`.
**Test:** `vitest src/lib/nakama/auth.test.ts` `localStorage rummy_device_id` persists across `authenticate` calls; unit uses `node` `localStorage` mock.

### Day 20 — Socket store

**Nakama:** `src/lib/nakama/socket.ts` `socketStore writable<Socket|null>`, `_socket` singleton, `_matchDataHandler`, `setMatchDataHandler(handler)` + wiring to `(_socket as any).onmatchdata`, `createSocket(sessionOverride?)` `await authenticate()` → `getClient().createSocket(false,false)` → `sock.connect(session,true)` → `onmatchdata/onDisconnect`; `isTestEnv` guard `VITEST || NODE_ENV=test` returns `createMockSocket()` (`onmatchdata null, connect noop, createMatch→{match:{matchId:'mock-match'}}, joinMatch, leaveMatch, sendMatchState noop`) to keep `vitest` deterministic.
**Test:** `vitest src/lib/nakama/socket.test.ts` `createSocket resolves, mock onmatchdata forwarded`.

### Day 21 — Match create/join (lobby primitive)

**Nakama:** `src/lib/nakama/match.ts` `matchStore writable<string|null>` `loadStoredMatchId()` `rummy_matchId`, `persistMatchId`, `extractMatchId(result,fallback)` handles `{match:{matchId|match_id|id}}` | `{matchId}` | `{id}`; `createMatch()` **branch** `isTestEnv → sock.createMatch() mock` vs real `rpc(session,'create_match','')` → `payload.matchId` → `sock.joinMatch(matchId)` → `persist`; `joinMatch(matchId)` `sock.joinMatch(matchId)` → `persist`; `listAvailableMatches()` `listMatches(session,10,auth,label,1,4,'')` 6 combos deduped `AvailableMatch{matchId,label,size}`.
**Why branch:** `vitest` never hits real `7350`; `playwright` always hits real `rpc` then WS.
**Test:** `vitest src/lib/nakama/match.test.ts` `createMatch → mock-match`, `alice create + bob join same`; `playwright` manual `lobby` `100` vs `101` later.

### Day 22 — Game store — private

**Nakama:** `src/lib/game/store.ts` `privateStore/privateBySeat/lastPrivate`, `onPrivateSnapshot(snap)` validates `isValidPrivateSnapshot` → `privateStore.set`, `privateBySeat.set(ownSeat)`, **also** `publicStore.set({v,gamePhase,turnPhase,currentSeat,players,stockCount,discardRow,tableMelds,winner})` (private contains public), `localStorage rummy_lastPrivate:${seat}`; `getLastPrivateForSeat`, `isMyTurn = derived(privateStore, currentSeat===ownSeat)`, `myRack`, `mySeat`.
**Test:** `vitest src/lib/game/store.test.ts` `privateStore ownRack 3 only local`.

### Day 23 — Game store — public (TableBoard subscribes)

**UI → Nakama:** `TableBoard.svelte` now `displayMelds = derived($publicStore)` `pub.tableMelds.map(tm→{id:tm.ID, kind: tm.Kind==='set'?set:run, tiles: tm.Tiles.map(t→{id:t.ID,colour:t.Colour,rank:t.Rank,isJoker:t.IsJoker})})` fallback to demo 3 melds only when `pub==null`.
**Test:** `vitest src/lib/game/store.test.ts` `publicStore TableMelds 2`, `playwright tests/table.e2e.ts` still `ETALĂRI` but now via `publicStore`.

### Day 24 — Redaction (prove no leak)

**UI:** `src/lib/game/snapshot.ts` `checkNoLeak(publicJson, privateIds)` `publicJson.includes(id)` loop.
**Test:** `vitest src/lib/game/redaction.test.ts` `6 tests` `publicJson not alice-secret`, exhaustive `2,3,4×seed` `PublicView` never contains `rack`/`stock` IDs.

### Day 25 — Reconnection — store (persist)

**Nakama:** `store.ts` `privateBySeat: Map<number,PrivateSnapshot>` + `localStorage rummy_lastPrivate:${seat}` on every `onPrivateSnapshot`; `getStoredMatchIdSafe()/getStoredUserIdSafe()` helpers; `socket.ts` `ondisconnect` clears `_socket` but **not** `localStorage rummy_matchId` (so `hasStoredMatch` stays true).
**Test:** `vitest` `rummy_lastPrivate:0` after `onPrivateSnapshot`, `reconnect.test.ts` `keeps matchId` after `disconnect`.

### Day 26 — Reconnection — rejoin (WS re-connect)

**Nakama:** `src/lib/nakama/reconnect.ts` `reconnect()` `getStoredMatchId() ?? getMatchId()` → `createSocket()` → `joinMatch(storedId)` expects `OpServerState 100` new `OwnRack` not old, `simulatePrivateAfterReconnect(newIds,seat)` helper; `src/routes/+page.svelte` `onMount: await authenticate(); await reconnect()` + `initGameStore()`.
**Test:** `vitest reconnect.test.ts` `reconnect ownRack new-1 not old-1`, `playwright` `F5` → `Se conectează la masa …` spinner until `100`.

### Day 27 — Versioning

**Nakama:** `snapshot.ts` `SnapshotVersion 1` guard `if snap.v!==1 bad_version ignore`; `store.ts` `onPrivate/onPublic` return `false` on bad, caller `handleMatchData` ignores.
**Test:** `vitest src/lib/game/versioning.test.ts` `10 tests` `v:2 ignored`, `v:1 ok`.

### Day 28 — Visual sync (2 clients, isolated demos)

**UI:** `src/routes/demo/sync/+page.svelte` two buttons `alice PrivateSnapshot ownRack [a1]` vs `bob [b1]` driving `onPrivate/onPublic`; `src/lib/game/store.ts` `initGameStore()` auto-wire `setMatchDataHandler` when `window !== undefined && process === undefined` (browser not test).
**Test:** `playwright tests/sync.e2e.ts` `alice vs bob different Rack` `PrivateSnapshot` `Rack only local`, `PublicSnapshot` `Table` shared; `demo/sync` screenshot.
**Milestone:** Stores are correct and isolated demos exist for every later opcode.

---

## Phase 3 — Game actions — Svelte → Nakama `1..9` (Days 29–42) — one opcode = one day

Pattern per day: **(a) UI `canX` derived from `privateStore`** → **(b) `src/lib/game/actions.ts` `sendX()` `NewEnvelope(op,payload,requestId)` `getSocket() ?? createSocket()` `sendMatchState(matchId,op,envelope)` → `lastSent` store** → **(c) `src/routes/demo/x/+page.svelte` isolated `onPrivateSnapshot` builder + buttons → `playwright tests/x.e2e.ts` asserts `canX` + `lastSent` contains `"op":N` + `mock-match`**.

### Day 29 — Start

`TopBar.svelte` `canStart = derived($privateStore=> gamePhase==='Waiting' && ownSeat===0 && players.length>=2)` → `sendStart()` `NewEnvelope(1,{},requestId)` `sendMatchState(matchId,1,env)` `OpClientStart`; `demo/start` `Waiting Host 2p` vs `Waiting Guest 2p` vs `1p` → `playwright start.e2e.ts` `START` visible only host `2p`.

### Day 30 — Opening discard

`Rack.svelte` `isOpeningDiscard && isMyTurn && rackCount===15 && selected.size===1` → `sendDiscard(tileId)` `NewEnvelope(2,{tileId})`; `demo/opening` `OpeningDiscard 15` `discard-tile` click; `tests/opening.e2e.ts` `OpeningDiscard MyTurn 15→14` `canDiscard`.

### Day 31 — Draw stock

`Rack.svelte: canDraw = Playing && MustDraw && isMyTurn` → `▶ TRAGE DIN TALON` `OpClientDrawStock 3 {}` `demo/draw` `MustDraw MyTurn` → `playwright draw.e2e.ts` `draw-btn` enabled only then, `op 3` `mock-match`.

### Day 32 — Draw previous (HasOpened gate)

`hasOpened` `players.find(p=>p.seat==ownSeat).hasOpened`, `canDrawPrevious = Playing MustDraw isMyTurn hasOpened last !IsOpeningDiscard` → `↩ IA ULTIMA` `4 {}`; `demo/draw-previous` `HasOpened+Discard` vs `NotOpened` vs `OpeningDiscard` vs `Empty`; `tests/draw-prev.e2e.ts`.

### Day 33 — Pickup for meld (selected 2 + TableBoard discard click)

`pickupDiscardIndex writable`, `TableBoard` `discard-tile-${idx}` `onClick pickupDiscardIndex.set(entry.Index)` `border-sky-500` when selected, `IsOpeningDiscard border-red-300`; `Rack` `canPickup = MustDraw HasOpened selected 2 + idx in range !IsOpeningDiscard` → `⬆ RIDICĂ PENTRU ETALARE` `OpClientPickupDiscardForMeld 5 {discardIndex,tileIds:[2]}`; `demo/pickup` `selected 2 + discardIndex`; `tests/pickup.e2e.ts`.

### Day 34 — Normal discard + turn rotation

`canDiscard extended` `Playing MeldOrDiscard isMyTurn selected 1` → `ARUNCĂ CARTEA` `OpClientDiscard 2 {tileId}` `CurrentSeat→(cur+1)%n` seen via `publicStore`; `demo/discard` `MeldOrDiscard 0→1`; `tests/discard.e2e.ts`.

### Day 35 — MeldInitial (`!HasOpened` `50+` `≥1 run`)

`canMeldInitial = !hasOpened && Playing MeldOrDiscard isMyTurn selected>=3` → `ETALEAZĂ SELECTATE` `sendMeldInitial([{kind:'run',tileIds}])` `NewEnvelope(6,{melds:[{id,kind,tileIds,jokerReps:{}}]})`; `demo/meld` `NotOpened selected 3`; `tests/meld-initial.e2e.ts` `op 6`.

### Day 36 — MeldNew (`HasOpened`)

`canMeldNew = hasOpened && … selected>=3` → `sendMeldNew` `7`; `demo/meld-new` `HasOpened`; `tests/meld-new.e2e.ts` `op 7`.

### Day 37 — Extend (TableBoard `selectedMeldId`)

`selectedMeldId` `writable null`, `TableBoard` `meld-m1 onClick selectedMeldId.set(meld.ID)` `border-sky-500`; `Rack` `canExtend = hasOpened MeldOrDiscard isMyTurn selected>=1 && meld exists` → `EXTINDE ETALAREA` `OpClientExtendMeld 8 {meldId,tileIds}`; `demo/extend` `HasOpened+meld`; `tests/extend.e2e.ts`.

### Day 38 — ReplaceJoker (joker tile click)

`replaceTargetMeldId`, `TableBoard` `joker-m1-* onClick e.stopPropagation() replaceTargetMeldId.set(meldId), selectedMeldId.set(meldId)` `amber-400 ring`; `Rack` `canReplace = hasOpened MeldOrDiscard selected==3 && meld.Tiles.some(joker)` → `ÎNLOCUIEȘTE JOLY` `OpClientReplaceJoker 9 {targetMeldId,tileId,newMeldTiles:[2]}`; `demo/joker` `HasOpened+joker`; `tests/joker.e2e.ts`.

### Day 39 — Winner `RoundComplete`

`PrivateSnapshot gamePhase=="RoundComplete" winner>=0` → `WinnerOverlay.svelte` fixed overlay `CÂȘTIGĂTOR Seat N` `RESTART MASA` `localStorage.removeItem('rummy_matchId')` + `resetGame()` + `clearMatchId()`; `demo/winner` `Winner 0`; `tests/winner.e2e.ts`.

### Day 40 — Visual actions (screenshot)

No new opcode; `tests/visual-actions.e2e.ts` `start→opening→draw→meld→extend→prev/pickup→replace→discard→win` sequence via `demo/*` `onPrivateSnapshot` builders, screenshot `TopBar+TableBoard+Rack`.

### Day 41 — E2E actions (2 browsers, `requestId`)

`tests/e2e-actions.e2e.ts` `2 BrowserContext` `alice createMatch + bob joinMatch same matchId` `alice Start → Opening` etc. via isolated `onPrivateSnapshot` + `lastSent.requestId` echo check; proves `rpc create_match` → `joinMatch` → `sendMatchState` chain works real.

### Day 42 — Error on invalid (first `OpServerError` toast)

`src/lib/game/errorStore.ts` + `src/components/Toast.svelte` `data-error-code` (`bg #dc2626`) introduced via `handleMatchData 102→onServerError`; `tests` `DISCARD tileId not in rack → OpServerError bad_payload` `data-error-code`.
**Milestone:** Full `1..9` loop playable over `MatchLoop`, `Playable` `↔` not just visual.

---

## Phase 4 — UI feedback & validation — Svelte derived (Days 43–52) — _after_ server works

### Day 43 — Turn indicator (today `TopBar` Day 43)

**UI:** `TopBar` `derived($privateStore+$publicStore => gamePhase/turnPhase/currentSeat/isMyTurn)` `showTurn=gamePhase!==''&&currentSeat>=0` → pill `Current: seat-N Playing/MustDraw ← rândul tău|← seat-N` `emerald` when `isMyTurn` else `white/10` + mobile `S N ← tu`.
**Test:** `playwright demo/turn` `MustDraw MyTurn → ← rândul tău` vs `NotMyTurn → ← seat-1` + `draw-btn` disabled when not `MustDraw/CurrentSeat`.

### Day 44 — HasOpened

**UI:** `Rack` header `DESCHIS ✓` emerald vs `NEDESCHIS • 50+ RUN` amber `data-testid="hasopened-badge"` + helper `deschiderea cere ≥50 pct cu cel puțin o suită`; `canDrawPrevious/canPickup/canExtend/canReplace` already gated (Days 32/33/37/38) — now visible.
**Test:** `playwright demo/hasopened` `NEDESCHIS → Prev/Pickup/Extend/Replace disabled`.

### Day 45 — Error toast (hardened)

**UI:** `Toast.svelte` `onServerError{raw: string|Uint8Array|object}` parses `{code,message,details,requestId,op}` → `errorStore` 3s auto-clear `clearTimeout` + `fade/fly` `bg-[#dc2626]` `role=alert` `aria-live assertive` `data-error-code/requestId/op` `details k: v •` + ✕ `clearError`; `+layout.svelte` global `<Toast/>`.
**Test:** `playwright demo/error` `Trigger bad_payload → toast visible bg rgb(220,38,38) data-error-code → LEAKED → bad_request → Clear → auto-dismiss 3.5s`.

### Day 46 — Meld validation (client preview, not trust)

`src/lib/game/validate.ts` `validateRun(1-2-3 ✓, 12-13-1 ✓, 13-1-2 ✗)/validateSet(3 vs 4 distinct)` `joker real>=2*joker` + `vitest` `validate.test.ts` `1-2-3` valid.

### Day 47 — Scoring preview (non-authoritative)

`Rack` `scorePreview(selected)` `TotalScore` `2–9:5 10–13:10 A low 5 high 10 set 25 joker rep` `total>=50 && ≥1 run` enables `ETALEAZĂ`; `playwright` `selected 3 50 enabled 49 disabled`.

### Days 48–52 — Discard/meld highlights, `duplicate` guard, visual + 2-browser feedback

`TableBoard` `IsOpeningDiscard ring-red-500`, `selectedMeldId ring-sky-500`, `Rack` `SvelteSet` duplicate `toast duplicate`, `2 browsers` `MELD_INITIAL 49 → OpServerError`.

---

## Phase 5 — Polish … (Days 53–75, 76–95 as in `client/docs/roadmap.md:110`)

Typography `Inter/JetBrains Mono`, `WCAG AA` `red #e53935`, `Table gap-2.5` no scroll `52px table vs 64px rack`, buttons `rounded-xl hover disabled`, modals `jokerReps`, `svelte/transition` `prefers-reduced-motion`, responsive `max-w-[1600px] 1280/768/375`, `a11y` `aria-label Tile red-5`, `i18n ro suita/terta/Joly`, `reconnecting` overlay.

Testing & release: `unit tile/meld/scoring`, `integration sync/flow CheckTileConservation 106`, `e2e auth/create/join/draw/meld/reconnect/error`, `visual regression`, `client-mvp-rc1` tag `vite build` green.

---

## Why each file lives where (architecture)

```
client/src/lib/nakama/  — only Nakama SDK wrappers (Client/Session/Socket) + localStorage persist
  client.ts   → getClient() singleton (VITE_NAKAMA_*)
  auth.ts     → rummy_device_id → authenticateDevice → rummy_token
  socket.ts   → createSocket(session) + setMatchDataHandler
  match.ts    → createMatch rpc→join + joinMatch + listAvailableMatches (isTestEnv branch)
  protocol.ts → Version 1 + NewEnvelope(v,op,requestId,payload)
  reconnect.ts→ getStoredMatchId → joinMatch
client/src/lib/game/     — pure game state (no SDK)
  snapshot.ts → Public/Private Snapshot + isValid + SnapshotVersion 1
  store.ts    → privateStore/publicStore + onPrivate/onPublic + handleMatchData(op 100/101/102) + isMyTurn
  errorStore.ts→ onServerError → errorStore 3s
  actions.ts  → sendStart/Discard/DrawStock/DrawPrevious/Pickup/MeldInitial/MeldNew/Extend/ReplaceJoker
  validate.ts → validateRun/Set (client preview)
  scoring.ts  → ScoreTile/TotalScore (preview)
client/src/components/   — props + derived stores, no SDK calls except via actions
  Tile.svelte, TableBoard.svelte (publicStore), Rack.svelte (privateStore + selected SvelteSet + canX), TopBar.svelte (private+public currentSeat), Toast.svelte (errorStore), WinnerOverlay.svelte
client/src/routes/demo/* — isolated onPrivateSnapshot builders + button → e2e-driven, no backend
client/tests/*.e2e.ts    — one file per Day, data-testid + getByRole exact:true to avoid substring flake
```

**Rule:** `components` never import `nakama-js` directly — only `actions.ts` and `store.ts`. `store.ts:handleMatchData` never trusts client `tile face values` — server `TableMeld.ID`/`JokerReps` is truth.

---

## Daily “what to commit” template

```bash
npm run check   # svelte-check + eslint + prettier --check + vitest --run + vite build must be green
npm run test:e2e -- tests/x.e2e.ts --project=chromium  # one file, headed for debug: --headed
git add client/src/... client/tests/... client/src/routes/demo/...
git commit -m "feat(client): <UI slice> + <Nakama wire> + <playwright> (Day N)"
git push
# next day only after this Day is green; never land 2 opcodes in one Day
```

**Milestone per phase:** Phase 0 `5173` + `SMOKE PASSED` + `SvelteKit Hello`; Phase 1 `TopBar/Table/Rack` no scroll; Phase 2 `alice vs bob different Rack` no leak; Phase 3 `start→win` loop over real `MatchLoop`; Phase 4 `turn/HasOpened/toast` matches server `AllowedOps`.

---

_How this file differs from `client/docs/roadmap.md`: that file is the table; this file is the **why and how** — which file, which derived, which `op`, which `data-testid`, which `vitest` vs `playwright`, and which `isTestEnv` branch. Last updated: 2026-08-27 after Day 45 (`5959151` Toast) for `client/` at `5173` vs `rummy:7350`._
