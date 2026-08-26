# Rummy Client — SvelteKit + Nakama Incremental Roadmap (Handmade Hero)

**Target:** Modern, maintainable, *playable* Romanian Tile Rummy web client in SvelteKit (`svelte@5` + `vite@5` + `tailwindcss@4` + `@heroiclabs/nakama-js@2.8` + `typescript@5`) that proves the authoritative Go backend (`docs/protocol.md` `Version 1`, `internal/protocol/opcodes.go:8` `1..9`/`100..199`) without copying Remi Online branding. Every day is a small, tested, committed vertical slice with `vitest` + `playwright` verification. No horizontal scroll on table — tiles auto-scale to fit.

**Backend prerequisite:** `rummy-mvp-rc1` at `main` (Go `1.23.5` + `nakama:3.26.0` + `postgres:15` at `127.0.0.1:7350` `defaultkey`) is stable. `docker compose up --build -d` + `make smoke` `SMOKE PASSED`. Client must never leak private racks: `PublicView` only `RackCount`/`StockCount`/`DiscardRow`/`TableMelds`/`Winner`, `PrivateView` adds `OwnRack` per `Seat` only (`internal/match/visibility.go:36` `Version 1`).

**Style:** Handmade Hero — one small vertical slice per day, `npm run check` (`svelte-check` + `eslint` + `prettier --check` + `vitest --run` + `vite build`) + `npm run test:e2e` green before next day, one focused commit per day, `git add <specific files>` + `git commit -m "feat(client): ... (Day N)"` + `git push`.

**How to read:** Day numbers are *client-only* (Day 1 is first client commit after reset, independent of backend days 1–78). Each day ends in a runnable, tested state with `npm run dev -- --open` at `http://localhost:5173` showing the slice and `npm run test:e2e` covering that slice.

---

## Phase 1 — Foundation and Local Development (Days 1–8)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 1 | Repo & tooling audit | Audit `rummy_backend` repo, `client/` empty after `rm -rf client`, `docs/protocol.md`/`state-machine.md`, `AGENTS.md:812` Day 24 acceptance, record in `client/docs/roadmap.md` + `client/README.md` + `.env.example` `VITE_NAKAMA_*` | `vitest` no test, `playwright` smoke `page.goto("/")` `200` |
| 2 | SvelteKit skeleton | `npx sv create client --template minimal --types ts --add prettier eslint vitest playwright tailwindcss --install npm`, `src/routes/+page.svelte` `Hello SvelteKit`, `npm run dev -- --open` shows `SvelteKit` at `localhost:5173` | `test:e2e` `page.getByText("SvelteKit")` |
| 3 | Tailwind & design tokens | `tailwindcss@4` `@tailwindcss/vite` `app.css` `@import "tailwindcss"`, `theme` `felt #0a2e1a`, `wood #5d4037`, `Inter` + `JetBrains Mono` via `googleapis`, `+page.svelte` uses `bg-felt` `max-w-[1600px]` | `vitest` `app.css` contains `@import "tailwindcss"` |
| 4 | Nakama JS client | `+ @heroiclabs/nakama-js@2.8.0`, `src/lib/nakama/client.ts` `Client("127.0.0.1","7350","defaultkey")` `authenticateDevice` `uuid` `token→localStorage`, `createSocket` lazy, `vite build` green | `test:unit` `getClient()` exists, `e2e` `localStorage rummy_device_id` |
| 5 | Project structure & lint | `src/lib/{nakama,game,ui}` + `src/components/{Tile,Table,Rack}` empty, `eslint` `svelte` + `prettier` + `svelte-check` `strict`, `npm run lint`/`check` | `npm run lint` pass (CI) |
| 6 | Env & docs | `src/lib/config.ts` `VITE_NAKAMA_*` via `$env/static/public`, `.env.example`, `client/README.md` `Quick Start` (`npm i`, `npm run dev`, `VITE_NAKAMA_*`), root `README` `Next Steps` | `e2e` `README` contains `VITE_NAKAMA_HOST` |
| 7 | Smoke test | `tests/smoke.spec.ts` verifies `npm run dev` serves `200` at `5173`, `+page.svelte` `200`, `nakama/client.ts` `Client` can `authenticateDevice` against `docker compose up` backend | `npm run test:e2e` `page.goto("/")` `200` + `SvelteKit` |
| 8 | CI baseline | `.github/workflows/client.yml` `eslint` `prettier --check` `svelte-check` `vitest --run` `vite build` `playwright` | CI green |

**Milestone:** `npm run dev` shows SvelteKit `Hello` at `http://localhost:5173` with Tailwind `bg-felt`, and `nakama/client.ts` can `authenticateDevice` against `docker compose up` backend. No game UI yet, but `npm run check` + `npm run test:e2e` green.

---

## Phase 2 — Design System & Layout (Days 9–16)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 9 | Design tokens | `src/lib/ui/tokens.ts` `colors: felt, wood, tile{red,yellow,blue,black}, feltLight`, `fonts: Inter, JetBrains Mono`, `spacing: 4/8/12/16`, `radius: 12/16/18` | `vitest` tokens `felt #0a2e1a` |
| 10 | Layout — GameSpace | `src/lib/layout.ts` `GAME_SPACE 1600×900` + `getLayout()` `TopBar 56px`, `TableArea`, `PlayerRackArea 132px`, `ActionBar 56px` — single source, no `px` hardcode, `Scale.FIT` equivalent via `max-w-[1600px] mx-auto` | `vitest` `getLayout().PlayerRackArea.height 132` |
| 11 | Tile component | `src/components/Tile.svelte` `colour 1..4` `rank 1..13` `isJoker` `size: rack|table` `selected` `onClick` `onDrag` — `52×72` table vs `64×90` rack, `colourToHex` + `rankToLabel` | `vitest` `Tile` `A` `K` `Joly` `red #e53935` |
| 12 | Tile modern | `Tile.svelte` `48×64` white `rounded-lg` `border` + `shadow`, `Joly` `amber-50`, selected `ring-2 sky-500 scale 1.06 + ✓`, `draggable` | `vitest` `selected` `ring-2` |
| 13 | Table board | `src/components/TableBoard.svelte` `melds: Meld[]` `TableBoard` `bg-[#f5f1e8] rounded-[18px]` `ETALĂRI PE MASĂ • n SETURI` `flex-wrap` + `gap-2.5` + `no scroll` `auto-scale` `scale()` if overflow | `playwright` `ETALĂRI PE MASĂ` `no scroll` |
| 14 | Rack | `src/components/Rack.svelte` `tiles: Tile[]` `Mâna ta • n cărți` `TREBUIE SĂ TRAGI` `flex-wrap gap-1.5` `no scroll` `min-h-[110px]` + `SORTEAZĂ` buttons + `TRAGE DIN TALON` etc. | `playwright` `Mâna ta` `15 cărți` `no scroll` |
| 15 | Top bar | `src/components/TopBar.svelte` `REMI ETALAT` `MASA 1 • 4 JUCĂTORI` `0s` `REGULI` `JOC NOU` `h-12 bg-black/90` | `playwright` `REMI ETALAT` `MASA 1` |
| 16 | Visual layout | `playwright` screenshot `TopBar` + `TableBoard` + `Rack` + `Jurnal` at `1280×800` + `375×667` | `e2e/layout.spec.ts` `TopBar` `TableBoard` `Rack` visible, no overlap |

**Milestone:** App shows `TopBar` + `TableBoard` (2 melds `66 pct` + `53 pct` + `55 pct` with `Tile` `52×72`) + `Rack` (11 tiles `TRAGE DIN TALON`) + `Jurnal` placeholder, all via `Tailwind` `flex-wrap` no scroll on desktop or mobile, `vitest` + `playwright` green.

---

## Phase 3 — State & Networking — Svelte Stores (Days 17–28)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 17 | Snapshot types | `src/lib/game/snapshot.ts` `PublicSnapshot{v,gamePhase,turnPhase,currentSeat,players{RackCount,HasOpened},stockCount,discardRow,tableMelds,winner}` `PrivateSnapshot extends PublicSnapshot {ownRack, ownSeat}` `SnapshotVersion=1` | `vitest` `SnapshotVersion 1` `isValidPublicSnapshot` |
| 18 | Nakama envelope | `src/lib/nakama/protocol.ts` `Envelope{v,op,requestId,payload}` `Version 1`, `OpClientStart 1`…`ReplaceJoker 9`, `OpServer* 100..103`, `NewEnvelope(op,payload)` | `vitest` `NewEnvelope(1,{})` `v:1` |
| 19 | Auth store | `src/lib/nakama/auth.ts` `authStore` `writable<Session|null>` `authenticate()` `deviceId→localStorage rummy_device_id` `token→rummy_token` | `vitest` `authStore` `localStorage rummy_device_id` |
| 20 | Socket store | `src/lib/nakama/socket.ts` `socketStore` `writable<Socket|null>` `createSocket()` `onmatchdata` → `gameStore` | `vitest` `createSocket` resolves |
| 21 | Match create/join | `src/lib/nakama/match.ts` `createMatch()` `socket.createMatch()` `joinMatch(id)` `matchId→rummy_matchId` | `playwright` `2 browsers` `alice createMatch` `bob joinMatch` same `matchId` |
| 22 | Game store — private | `src/lib/game/store.ts` `privateStore: writable<PrivateSnapshot\|null>` `onPrivateSnapshot` `lastPrivate` `derived isMyTurn` | `vitest` `privateStore` `ownRack 3` only local |
| 23 | Game store — public | `publicStore: writable<PublicSnapshot\|null>` `onPublicSnapshot` `TableBoard` subscribes, not `OwnRack` | `vitest` `publicStore` `TableMelds 2` |
| 24 | Redaction | `checkNoLeak(publicJson, privateIds)` `JSON.stringify(PublicSnapshot)` search for `OwnRack` IDs | `vitest` `publicJson` not `alice-secret` |
| 25 | Reconnection — store | `privateBySeat: Map<seat,PrivateSnapshot>` `localStorage rummy_lastPrivate:${seat}` `socket.onDisconnect` keep `matchId/userId` | `vitest` `rummy_lastPrivate:0` |
| 26 | Reconnection — rejoin | `reconnect()` `socket.connect`+`joinMatch(matchId)` expects `OpServerState 100` `PrivateSnapshot` for that `Seat` only, `Rack` rehydrates from `OwnRack` not old | `playwright` `reconnect()` `ownRack new-1` not `old-1` |
| 27 | Versioning | `SnapshotVersion=1` check: if `snap.v !==1` `bad_version` ignore | `vitest` `v:2` ignored |
| 28 | Visual sync | `playwright` `PrivateSnapshot` `Rack` only local, `PublicSnapshot` `Table` for all | `e2e/sync.spec.ts` `alice` vs `bob` different `Rack` |

**Milestone:** Two `page` can `joinMatch` and receive `PrivateSnapshot` (only own `Rack`) + `PublicSnapshot` (table), `Rack` and `TableBoard` re-render via `stores`, no leak, reconnection keeps `ownRack`.

---

## Phase 4 — Game Actions — Svelte → Nakama (Days 29–42)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 29 | Start | `TopBar` `Start` visible only if `Waiting` `ownSeat==0` `players>=2`, `click` `OpClientStart 1 {}` `sendMatchState` | `playwright` `Waiting` `Start` host `2p` |
| 30 | Opening discard | `Rack` `discardSelected()` when `OpeningDiscard` `ownSeat==currentSeat` `ownRack 15` `OpClientDiscard 2 {tileId}` | `playwright` `OpeningDiscard` `15→14` |
| 31 | Draw stock | `Rack` `Draw` visible only if `Playing MustDraw` `ownSeat==currentSeat`, `OpClientDrawStock 3 {}` disable until `MeldOrDiscard` | `playwright` `MustDraw` `Draw` enabled |
| 32 | Draw previous | `DrawPrevious` `HasOpened` `discardRow not empty` `!IsOpeningDiscard`, `OpClientDrawPreviousDiscard 4 {}` | `playwright` `HasOpened` `DrawPrev` enabled |
| 33 | Pickup for meld | `Rack` `selected 2` + `discardIndex` via `TableBoard` click, `OpClientPickupDiscardForMeld 5 {discardIndex, tileIds:[2]}` | `playwright` `selected 2` `Pickup` |
| 34 | Normal discard | `Rack` `discardSelected()` when `MeldOrDiscard` `ownSeat==currentSeat` `OpClientDiscard 2 {tileId}` `CurrentSeat→(current+1)%n` | `playwright` `MeldOrDiscard` `Discard` `0→1` |
| 35 | Meld initial | `Rack` `selected>=3` `!HasOpened` `OpClientMeldInitial 6 {melds:[{kind:run, tileIds}]}` `total>=50` `≥1 run` | `playwright` `!HasOpened` `MeldInitial` `50+` |
| 36 | Meld new | `HasOpened` `OpClientMeldNew 7 {melds}` | `playwright` `HasOpened` `MeldNew` |
| 37 | Extend meld | `TableBoard` `onMeldClicked(meldId)` + `Rack` `selected>=1` `HasOpened` `OpClientExtendMeld 8 {meldId, tileIds}` | `playwright` `HasOpened` `Extend` |
| 38 | Replace joker | `TableBoard` `onJokerClicked` + `Rack` `replaceSelected(targetMeldId, tileId, new1, new2)` `OpClientReplaceJoker 9 {targetMeldId, tileId, newMeldTiles[2]}` | `playwright` `Replace` |
| 39 | Winner | `PrivateSnapshot` `GamePhase=="RoundComplete"` `Winner` overlay `RESTART MASA` | `playwright` `RoundComplete` `Winner 0` |
| 40 | Visual actions | `playwright` `start→opening→draw→meld→extend→prev/pickup→replace→discard→win` | `playwright` screenshot `actions` |
| 41 | E2E actions | `2 browsers` `alice/bob` `draw`/`discard`/`meld`/`win` via `MatchLoop` `requestId` | `playwright` `alice` `bob` `win` |
| 42 | Error on invalid | `OpServerError 102` `code=bad_payload` `toast` `data-error-code` | `playwright` `DISCARD tileId not in rack` `OpServerError` |

**Milestone:** Player can `start` → `opening discard` → `draw` → `meld 50+` → `extend` → `prev/pickup` → `replace` → `discard` → `win` via `stores` + `sendMatchState`.

---

## Phase 5 — UI Feedback & Validation — Svelte (Days 43–52)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 43 | Turn indicator | `TopBar` `Current: seat-0` `Playing/MustDraw` `← current`, `Draw` disabled if not `MustDraw`/`CurrentSeat` | `playwright` `turn` `← current` |
| 44 | HasOpened | `Rack` `HasOpened` per `PublicPlayer` disables `Prev`/`Pickup`/`Extend`/`Replace` if `!HasOpened` | `playwright` `HasOpened false` disabled |
| 45 | Error toast | `src/components/Toast.svelte` `OpServerError` `code/message/details/requestId/op` `3s` `bg #dc2626` | `playwright` `LEAKED` `data-error-code` |
| 46 | Meld validation | `src/lib/game/validate.ts` `validateRun`/`validateSet` `1-2-3` `12-13-1` `13-1-2` `joker` `real>=2*joker` | `vitest` `1-2-3` valid |
| 47 | Scoring preview | `Rack` `scorePreview(selected)` `TotalScore` `2-9:5 10-13:10 A 5/10/25` `total>=50` `≥1 run` | `playwright` `selected 3` `50` enabled `49` disabled |
| 48 | Duplicate | `Rack` `selected Set` prevents `TileId` duplicate, `toast duplicate` | `playwright` `duplicate` |
| 49 | Discard highlight | `TableBoard` highlights `IsOpeningDiscard` `ring-red-500` + `Index` | `playwright` `disc-open` `ring` |
| 50 | Meld highlight | `TableBoard` highlights `selectedMeldId` `run/set` `JokerReps` | `playwright` `selectedMeldId` |
| 51 | Visual feedback | `playwright` `selected 3` valid `score 50` `turn` `HasOpened` `toast` | `playwright` screenshot |
| 52 | E2E feedback | `2 browsers` `selected` `HasOpened` `toast` `MELD_INITIAL 49` `OpServerError` | `playwright` `MELD_INITIAL 49` |

**Milestone:** Client gives immediate, non-authoritative feedback for `selected` validity and `Draw`/`Meld` enablement, server `OpServerError` canonical.

---

## Phase 6 — Polish, Responsiveness & Accessibility (Days 53–75)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 53 | Typography | `Inter`/`JetBrains Mono` scale `RackCount`/`StockCount` | `playwright` `fontFamily Inter` |
| 54 | Colours | `red #e53935`/`yellow #f9a825`/`blue #1e88e5`/`black #212121` `felt #0a2e1a` WCAG AA | `playwright` `tile red #e53935` |
| 55 | Spacing | `TableBoard` `gap-2.5` `Rack` `gap-1.5 sm:gap-2` no scroll on table (52px vs 64px) | `playwright` `no scroll` |
| 56 | Buttons | `Rack` `Draw` `Prev` `Meld Set/Run` `Discard` `rounded-xl` `hover` `disabled` | `playwright` `Draw` hover |
| 57 | Modals | `Modal` `jokerReps` `colour`/`rank` dropdown | `playwright` `jokerReps` |
| 58 | Rack polish | `Rack` `rounded-2xl bg-[#1a1a1a] border` `shadow-xl` | `playwright` `rounded-2xl` |
| 59 | Table polish | `TableBoard` `rounded-[18px] bg-[#f5f1e8] border-[#e8e0c8]` `shadow-inner` | `playwright` `rounded-[18px]` |
| 60 | Tiles polish | `Tile` `w-[64px] h-[90px]` `rounded-lg border` `rank 22px` `●` | `playwright` `Tile` `A` `K` `5` `J` |
| 61 | Animations | `svelte/transition` `fly/scale` `Draw stock→rack 200ms` `prefers-reduced-motion` | `playwright` `transition` |
| 62 | Sound | `+melt`/`discard`/`win` muted `toggle` | `playwright` `sound` muted |
| 63 | Responsive desktop | `max-w-[1600px] mx-auto` `flex-row` `3 cols` `TopBar 56px` `Rack 132px` | `playwright` `1280×800` `3 cols` |
| 64 | Responsive tablet | `768×1024` portrait `Rack` `min-h-[320px]` `Table` scroll | `playwright` `768×1024` |
| 65 | Responsive mobile | `375×667` `Rack` `overflow-x-auto` `Table` stacked `Jurnal` `hidden lg:flex` | `playwright` `375×667` `Rack` scroll |
| 66 | A11y keyboard | `1..9` select `Enter` discard/meld `Tab` switch seat | `playwright` `press 1` select |
| 67 | A11y screen-reader | `aria-label` `Tile red-5 id` `axe` | `playwright` `axe` `0 violations` |
| 68 | A11y colour-blind | `R`/`Y`/`B`/`K` text labels + `black` outline not colour alone | `playwright` `R` `Y` `B` `K` |
| 69 | A11y reduced motion | `matchMedia("(prefers-reduced-motion: reduce)")` disable `transition` | `playwright` `prefers-reduced-motion` |
| 70 | Localization | `src/lib/i18n/{en,ro}.json` `t("draw")` `ro` `suita/terta/Joly` | `playwright` `ro` `suita` |
| 71 | Pluralization | `Intl.PluralRules` `1 carte` vs `5 cărți` `RackCount` | `playwright` `1 carte` `5 cărți` |
| 72 | Connection UX | `reconnecting` overlay `socket.onDisconnect` `resync` spinner until `OpServerState 100` | `playwright` `reconnecting` |
| 73 | Visual polish | `playwright` `Inter` `colours` `spacing` `buttons` `rack` `table` `tiles` | `playwright` screenshot |
| 74 | A11y e2e | `2 browsers` `keyboard` `screen-reader` `colour-blind` | `playwright` `keyboard` |
| 75 | E2E polish | `2 browsers` `responsive` `localization` `connection` | `playwright` `responsive` |

**Milestone:** Client is usable with keyboard, screen-reader, colour-blind, mobile, subtle animations and no private leak, with `Inter`/`JetBrains Mono` polished — inspired by `remi-online.ro` but better (no scroll on table, tiles auto-scale to fit).

---

## Phase 7 — Testing, Docs & Release (Days 76–95)

| Day | Focus | Deliverable | Test |
|---:|---|---|---|
| 76 | Unit: tile | `src/lib/game/tile.test.ts` `TileInstance` `Colour`/`Rank` `IsJoker` | `vitest` `TileInstance` |
| 77 | Unit: meld | `src/lib/game/meld.test.ts` `ValidateSet`/`ValidateRun` `1-2-3`/`12-13-1`/`13-1-2`/`joker` | `vitest` `1-2-3` |
| 78 | Unit: scoring | `src/lib/game/scoring.test.ts` `ScoreTile` `5/10/25` `TotalScore 30/90` `ValidateInitialBatch 49/50` | `vitest` `5/10/25` |
| 79 | Integration: sync | `src/lib/game/store.test.ts` `PublicView`/`PrivateView` `Version 1` `reconnection` `redaction` | `vitest` `redaction` |
| 80 | Integration: flow | `src/lib/game/flow.test.ts` `TestDeterministicSimulation` `MatchInit→Win` `CheckTileConservation 106` | `vitest` `deterministic` |
| 81 | E2E: auth | `tests/auth.spec.ts` `authenticateDevice` `defaultkey` `token` `userId` | `playwright` `token` |
| 82 | E2E: create/join | `tests/match.spec.ts` `createMatch`/`joinMatch` `alice`/`bob` `PrivateSnapshot` per seat | `playwright` `alice` `bob` |
| 83 | E2E: draw/discard | `tests/draw-discard.spec.ts` `DRAW_STOCK` `DISCARD` `CurrentSeat 0→1` | `playwright` `CurrentSeat 0→1` |
| 84 | E2E: melding | `tests/meld.spec.ts` `MELD_INITIAL 50+ run` `MELD_NEW` `EXTEND` | `playwright` `MELD_INITIAL` |
| 85 | E2E: reconnection | `tests/reconnect.spec.ts` `socket.disconnect`+`connect`+`joinMatch` `OwnRack` `Winner` | `playwright` `reconnect` |
| 86 | E2E: error states | `tests/error.spec.ts` `DISCARD tileId not in rack` `MELD_INITIAL 49` `OpServerError 102` | `playwright` `OpServerError` |
| 87 | Docs: README | `client/README.md` `Next Steps` to `Day 76–86` `How to test` `npm run test:unit` `npm run test:e2e` | `playwright` `README` |
| 88 | Docs: roadmap | `client/docs/roadmap.md` `Phase 7` `Day 76–86` `✅ Done` `Current State` | `playwright` `roadmap` |
| 89 | Visual regression | `playwright` `visual` `rack` `table` `melds` `discard` `stock` `turn` | `playwright` `visual` |
| 90 | E2E full | `playwright` `2 browsers` `full` `auth` `create/join` `draw/discard` `meld` `reconnect` `error` | `playwright` `full` |
| 91 | Release candidate | `client-mvp-rc1` tag when `npm run check`+`test:unit`+`test:e2e`+`vite build` green + `go test` `no leak` | `git tag client-mvp-rc1` |
| 92 | Release notes | `client/docs/release-notes.md` `Stack: SvelteKit 5 + Tailwind 4 + nakama-js 2.8 + TS 5 + Go rummy-mvp-rc1` `known limitations` | `playwright` `release-notes` |
| 93 | Docs: testing | `docs/testing.md` `unit` `integration` `simulation` `Docker` `e2e` | `playwright` `testing.md` |
| 94 | Visual docs | `playwright` `docs` screenshots | `playwright` `docs` |
| 95 | E2E docs | `playwright` `docs` `README` `roadmap` | `playwright` `docs` |

**Milestone:** Two `pc` can `npm run dev` at `http://localhost:5173`, `create/join` via `defaultkey`, see `PrivateView.OwnRack` vs `PublicView` counts, `draw`/`discard`/`meld`/`extend`/`pickup`/`replace`/`win` through `MatchLoop`, `OpServerError` toasts, `client-mvp-rc1` tagged — better than `remi-online.ro` (no table scroll, tiles `52px` vs `64px`, auto-wrap).

---

## Extended Roadmap Summary (SvelteKit Future)

| Phase | Days | Primary Outcome | Test |
|---:|---|---|---|
| 1 | 1–8 | Foundation — SvelteKit + Tailwind + Nakama | `smoke` `scaffold` |
| 2 | 9–16 | Design system & Layout — tokens, Tile, Table, Rack, TopBar | `layout` `Tile` |
| 3 | 17–28 | State & Networking — Svelte stores, envelope, auth, sync | `sync` `redaction` |
| 4 | 29–42 | Game actions — Svelte → Nakama 9 ops + winner | `actions` `invalid` |
| 5 | 43–52 | UI feedback — Turn, HasOpened, Toast, validate, scoring | `feedback` |
| 6 | 53–75 | Polish — typography, colours, spacing, a11y, responsive | `polish` |
| 7 | 76–95 | Testing & release — unit, integration, e2e, docs, rc1 | `full` |

By **Day 95** (client), the project evolves from empty `client/` to a playable, accessible, localized SvelteKit web client (better than `remi-online.ro` — no table horizontal scroll, tiles `52px` table vs `64px` rack auto-scale via `flex-wrap` gap, `max-w-[1600px]` responsive) that proves the protocol `Version 1` and can be run via `npm run dev` against `docker compose up --build -d` without oral knowledge. Execution remains Handmade Hero: one vertical slice per day, `npm run dev` green, focused commit, push, **Playwright after each slice**.

---

## Daily Definition of Done

Every day must end with:

1. A small, working change (no dump).
2. New or updated `tests/**/*.spec.ts` (Playwright) + `src/**/*.test.ts` (vitest) where pure.
3. `npm run check` (`svelte-check` + `eslint` + `prettier --check` + `vitest --run` + `vite build`) + `npm run test:e2e` green.
4. Docs update if behavior/protocol changed.
5. Focused `git commit -m "feat(client): ... (Day N)"` + `git push`.
6. Screenshot or trace for visual slices.

## Commit Pattern

```
feat(client): add SvelteKit skeleton (Day 2)
feat(client): add Tile component with auto-scale (Day 11)
feat(client): add TableBoard with no scroll (Day 13)
feat(client): handle private snapshot with Rack re-render (Day 22)
test(client): add no-scroll e2e for table vs rack (Day 13)
fix(client): use Svelte stores for sync and fix alignment
```

---

*This roadmap lists days that **will be** implemented for the SvelteKit client. For what **is** implemented on the backend, see `docs/IMPLEMENTED.md`; for backend roadmap, see `docs/roadmap.md`. Last updated: 2026-08-26 for `client/` SvelteKit Day 1 start with Playwright.*
