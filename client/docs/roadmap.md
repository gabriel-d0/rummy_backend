# Rummy Client — Phaser 3 Incremental Roadmap with Playwright (Handmade Hero)

**Target:** Modern, maintainable, *playable* Romanian Tile Rummy web client in Phaser 3 (`phaser@3.80` + `vite@5` + `typescript@5` + `@heroiclabs/nakama-js@2.9`) that proves the authoritative Go backend (`docs/protocol.md` `Version 1`, `internal/protocol/opcodes.go:8` `1..9`/`100..199`) without copying Remi Online branding. Every day is a small, tested, committed vertical slice with Playwright verification.

**Backend prerequisite:** `rummy-mvp-rc1` at `main` (Go `1.23.5` + `nakama:3.26.0` + `postgres:15` at `127.0.0.1:7350` `defaultkey`) is stable. `docker compose up --build -d` + `make smoke` `SMOKE PASSED`. Client must never leak private racks: `PublicView` only `RackCount`/`StockCount`/`DiscardRow`/`TableMelds`/`Winner`, `PrivateView` adds `OwnRack` per `Seat` only (`internal/match/visibility.go:36` `Version 1`).

**Style:** Handmade Hero — one small vertical slice per day, `make client-check` (`lint` + `typecheck` + `vitest` + `build`) + `npx playwright test` green before next day, one focused commit per day, `git add <specific files>` + `git commit -m "feat(client): ... (Day N)"` + `git push`.

**How to read:** Day numbers are *client-only* (Day 1 is first client commit, independent of backend days 1–78). Each day ends in a runnable, tested state with `npm run dev` at `http://127.0.0.1:5173` showing the slice and `npx playwright test` covering that slice.

---

## Phase 1 — Foundation and Local Development (Days 1–7)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 1 | Repo & tooling audit | Audit `rummy_backend` repo, `client/` empty, `docs/protocol.md`/`state-machine.md`, `AGENTS.md:812` Day 24 acceptance, record in `client/docs/roadmap.md` + `client/README.md` | `e2e/smoke.spec.ts` checks repo exists |
| 2 | Vite + TS + Phaser scaffold | `npm init` + `vite@5` + `ts@5` + `phaser@3.80`, `vite.config.ts:5173` `host 127.0.0.1`, `index.html #game`, `src/main.ts` `new Phaser.Game({type: Phaser.AUTO, parent:"game", width:1024, height:768, scene:[Preload]})`, `npm run dev` shows `Phaser 3` `Init` | `e2e/scaffold.spec.ts` canvas exists, no-console-error |
| 3 | Nakama JS client | `+ @heroiclabs/nakama-js@2.9`, `src/net/nakama.ts` `Client("127.0.0.1","7350","defaultkey")` `authenticateDevice` `uuid` `token→localStorage`, `createSocket` lazy, `npm run build` green | `e2e/nakama.spec.ts` `getClient()` exists, `authenticate` not called yet |
| 4 | Project structure & lint/format | `src/scenes|net|state|ui`, `eslint` + `prettier` + `tsconfig strict`, `npm run lint`/`fmt`/`typecheck` + `make client-*` | `e2e/lint.spec.ts` `npm run lint` pass (via CI) |
| 5 | Dev scripts & docs | `client/README.md` `Quick Start` (`npm i`, `npm run dev`, `VITE_NAKAMA_*`), `client/.env.example`, root `README` `Next Steps` | `e2e/docs.spec.ts` `README` contains `VITE_NAKAMA_HOST` |
| 6 | Smoke test | `scripts/client-smoke.sh` or `e2e/smoke.spec.ts` verifies `npm run dev` serves `200` at `5173`, `Preload` `complete` log, `nakama.ts` `Client` can `authenticateDevice` against `docker compose up` backend | `e2e/smoke.spec.ts` `page.goto("/")` `200` + `canvas` + `Phaser.VERSION` |
| 7 | CI baseline | `.github/workflows/client.yml` `eslint` `prettier --check` `tsc --noEmit` `vitest` `vite build` | CI green |

**Milestone:** `npm run dev` shows blank Phaser canvas at `http://127.0.0.1:5173` and `nakama.ts` can `authenticateDevice` against `docker compose up` backend.

---

## Phase 2 — Phaser Scene & Asset Pipeline (Days 8–14)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 8 | Preload scene | `src/scenes/Preload.ts` `preload()` loads `tile.png` 32×48 `1×1` placeholder + `create()` `start TableScene`, console `Preload complete` | `e2e/preload.spec.ts` `console` contains `Preload complete` |
| 9 | Tile & joker sprites | `public/assets/tile.png` `red-1` + `joker.png` `Joly`, `Preload` loads both, `TableScene` `add.image(100,100,"tile")` | `e2e/assets.spec.ts` `tile.png` `200` + `joker.png` `200` |
| 10 | Table background | `public/assets/table.png` green felt `1024×768`, `TableScene` `add.image(512,384,"table")` behind tiles | `e2e/table.spec.ts` `table.png` loaded, canvas has green pixel |
| 11 | Rack background | `public/assets/rack.png` brown `800×120`, `RackScene` `add.image(512,680,"rack")` parallel scene | `e2e/rack.spec.ts` `rack.png` loaded |
| 12 | Asset manifest | `src/scenes/assets.ts` `{tile:"assets/tile.png", joker:"...", table:"...", rack:"..."}` + `Preload.test.ts` loads 4 without `404` | `e2e/manifest.spec.ts` `ASSET_MANIFEST` 4 entries `200` |
| 13 | Vite publicDir | `vite.config.ts` `publicDir:"public"` + `assets` under `public/assets` | `e2e/public.spec.ts` `fetch /assets/tile.png` `ok` |
| 14 | Asset build | `npm run build` includes `tile.png` `joker.png` `table.png` `rack.png` in `dist/assets` | `e2e/build.spec.ts` `dist` contains 4 assets |

**Milestone:** Phaser shows static table with one tile and a rack, no networking, all assets `200`.

---

## Phase 3 — Rack & Table Rendering — Static (Days 15–22)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 15 | Rack rendering (static) | `src/ui/Rack.ts` `renderRack(tiles:TileInstance[], seat:Seat)` `this.add.image` at `x=100+i*50`, `RackScene` mock `PrivateView.OwnRack` 2 tiles `red-5` `blue-7` | `e2e/rack.spec.ts` 2 tiles visible via `renderRack` count |
| 16 | Rack sorting | `src/ui/Rack.ts` `sortRack` by `Colour` then `Rank` `1..13`, test `red-13, red-1, blue-5 → red-1, red-13, blue-5` | `e2e/sort.spec.ts` `sortRack` order via `page.evaluate` |
| 17 | Table melds (static) | `src/scenes/TableScene.ts` `renderTableMelds(melds:TableMeld[])` at `y=100+i*80`, mock `PublicView.TableMelds` `run 1-2-3 red` `set 7 R/Y/B` | `e2e/melds.spec.ts` 2 melds `run` `set` via `getData("isTableMeld")` count |
| 18 | Discard row (static) | `src/ui/DiscardRow.ts` `renderDiscardRow(discardRow:DiscardEntry[])` at `x=100+i*40 y=300`, mock `IsOpeningDiscard` red border | `e2e/discard.spec.ts` 3 discards, `IsOpeningDiscard` has red border |
| 19 | Stock & turn (static) | `src/ui/StockCount.ts` `renderStockCount(count)` `Stock:77` + `TurnIndicator` `Current: seat-0` `GamePhase/TurnPhase` mock `PublicView` | `e2e/stock.spec.ts` `Stock: 77` `seat-0` visible |
| 20 | Layout helper | `src/ui/Layout.ts` `getLayout(w,h)` + `getSubspaces()` `GAME_SPACE 1000×1000` `Scale.FIT` | `e2e/layout.spec.ts` `getLayout(1000,1000).rackSlots.length 14` |
| 21 | Responsive — desktop | `main.ts` `scale:{mode:Fit, autoCenter:CenterBoth, width:1000, height:1000}` centered | `e2e/responsive.spec.ts` `1000×1000` `FIT` |
| 22 | Visual check | `playwright` screenshot `rack` + `table` + `discard` + `stock` at `1000×1000` | `e2e/visual.spec.ts` screenshot `rack-table` |

**Milestone:** Client renders static hand, table, discard, stock, turn from mock `PublicView`/`PrivateView` without networking, centered via `GAME_SPACE`.

---

## Phase 4 — Input Handling — Local Only (Days 23–29)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 23 | Tile selection | `src/ui/Rack.ts` `selected:Set<TileInstanceId>` `onTileClicked` `tint 0xffff00`, `RackScene` `pointerdown` logs `selected` | `e2e/selection.spec.ts` click tile → `selected` size 1, tint |
| 24 | Drag start | `this.input.setDraggable(tileImage)` per tile, `RackScene` `dragstart` logs `tileId` | `e2e/drag.spec.ts` `dragstart` event |
| 25 | Drop zone (local) | `TableScene` `drop` zone at `y=200` `setInteractive({dropZone:true})`, logs `tileId` dropped | `e2e/drop.spec.ts` `drop` zone `600×44` `isDropZone` |
| 26 | Discard selection | `RackScene` `discardSelected()` validates `selected.size===1` logs `DISCARD {tileId}` | `e2e/discard-select.spec.ts` `selected 1` → `DISCARD` else `null` |
| 27 | Meld selection | `RackScene` `meldSelected(kind)` validates `selected.size>=3` logs `MELD_INITIAL {kind, tileIds}` | `e2e/meld-select.spec.ts` `selected 3` `run/set` else `null` |
| 28 | Keyboard | `RackScene` `1..9` selects `rack[i]`, `Enter` discard, `Tab` switch seat (local 2p) | `e2e/keyboard.spec.ts` `press 1` selects |
| 29 | Visual input | `playwright` drag `tile` `100,680 → 500,400` + `drop` log | `e2e/drag-visual.spec.ts` drag screenshot |

**Milestone:** Player can click/drag tiles and see logs for `discard`/`meld` intents, no server call, selection tint works.

---

## Phase 5 — Networking: Nakama JS & Envelope (Days 30–38)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 30 | Device auth | `src/net/nakama.ts` `authenticate()` `client.authenticateDevice(deviceId, true, username)` `defaultkey` `127.0.0.1:7350` `token→localStorage` | `e2e/auth.spec.ts` `Client` exists, `localStorage rummy_device_id` |
| 31 | Socket connect | `createSocket()` `client.createSocket(false,false)` + `socket.connect(session, true)` `appearOnline=false` | `e2e/socket.spec.ts` `createSocket` resolves, `socket.connected` |
| 32 | Create match | `createMatch(socket)` `socket.createMatch()` `matchId→localStorage rummy_matchId` | `e2e/create-match.spec.ts` `matchId` `rummy_` prefix |
| 33 | Join match | `joinMatch(socket, matchId)` `socket.joinMatch(matchId)` `match_presence_event` | `e2e/join.spec.ts` 2nd `joinMatch` same `matchId` |
| 34 | Envelope | `src/net/protocol.ts` `Envelope{v,op,requestId,payload}` `Version 1`, `OpClientStart 1`…`ReplaceJoker 9`, `OpServer* 100..103`, `NewEnvelope(op,payload)` | `e2e/envelope.spec.ts` `NewEnvelope(1,{})` `json` `v:1` |
| 35 | Send match state | `sendMatchState(socket, matchId, op, payload)` `socket.sendMatchState(matchId, op, JSON.stringify(Envelope))` `requestId` `uuid` | `e2e/send.spec.ts` `sendMatchState` calls `socket.sendMatchState` |
| 36 | Receive match state | `src/net/nakama.ts` `socket.onmatchdata` parses `Envelope` routes `op 100/101/102/103` to `src/state/sync.ts` handlers | `e2e/receive.spec.ts` `onmatchdata` `100` → `onPrivateSnapshot` |
| 37 | Error display (net) | `src/ui/ErrorToast.ts` `OpServerError 102` `code/message/requestId/op` red toast `3s` | `e2e/error-toast.spec.ts` `LEAKED` toast `data-error-code` |
| 38 | E2E net | `playwright` `2 browsers` `alice createMatch` `bob joinMatch` `send Envelope` `receive` | `e2e/net-e2e.spec.ts` 2 `page` `alice/bob` `matchId` same |

**Milestone:** Client can `authenticateDevice` → `createSocket` → `createMatch`/`joinMatch` → `send Envelope` → `receive OpServerError/Event` via WebSocket.

---

## Phase 6 — State Synchronization — Authoritative (Days 39–48)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 39 | Public snapshot types | `src/state/snapshot.ts` `PublicSnapshot{v,gamePhase,turnPhase,currentSeat,players{RackCount,HasOpened},stockCount,discardRow,tableMelds,winner}` `SnapshotVersion=1` | `e2e/snapshot.spec.ts` `SnapshotVersion 1` |
| 40 | Private snapshot handling | `src/state/sync.ts` `onPrivateSnapshot(snap)` `lastPrivate:PrivateSnapshot` `RackScene` re-renders `OwnRack` only | `e2e/private.spec.ts` `ownRack 3` only local `Stock:77` still public |
| 41 | Public snapshot handling | `onPublicSnapshot(snap)` `lastPublic:PublicSnapshot` `TableScene` re-renders `TableMelds/DiscardRow/StockCount/CurrentSeat` not `OwnRack` | `e2e/public.spec.ts` `TableMelds 2` `DiscardRow 3` `Stock 77` |
| 42 | Redaction check | `src/state/sync.ts` `checkNoLeak(publicJson)` `JSON.stringify(PublicSnapshot)` search for `OwnRack` IDs as in `visibility_test.go` `no leak`/`LEAKED` | `e2e/redaction.spec.ts` `publicJson` not `alice-secret` |
| 43 | Reconnection: store last snapshot | `src/state/sync.ts` `lastPrivate:Map<seat,PrivateSnapshot>` `localStorage rummy_lastPrivate:${seat}` + `rummy_lastPrivate:map`, `socket.onDisconnect` keep `matchId/userId` | `e2e/store.spec.ts` `localStorage rummy_lastPrivate:0` |
| 44 | Reconnection: rejoin | `src/net/nakama.ts` `reconnect()` `socket.connect`+`joinMatch(matchId)` expects `OpServerState 100` `PrivateSnapshot` for that `Seat` only (`rummy_match.go:79`), `sync` rehydrates `RackScene` from `OwnRack` not old | `e2e/reconnect.spec.ts` `reconnect()` `ownRack new-1` not `old-1` |
| 45 | Versioning | `snapshot.ts` `SnapshotVersion=1` check: if `snap.v !==1` `bad_version` ignore as in `parser.go:22` | `e2e/version.spec.ts` `v:2` ignored, `v:1` stored |
| 46 | State machine client | `src/state/sync.ts` `allowedOps(gamePhase,turnPhase)` mirror `phases.go:15` (`Waiting→Start` host 0, `OpeningDiscard→Discard` seat 0, `Playing MustDraw→DrawStock/DrawPrevious/Pickup`, `MeldOrDiscard→Discard/MeldInitial/MeldNew/Extend/Replace`, `RoundComplete→none`), client disables buttons per `CurrentSeat/TurnPhase/HasOpened` | `e2e/state-machine.spec.ts` `Waiting` `Start` host only, `Playing MustDraw` `Draw` enabled |
| 47 | Visual sync | `playwright` `PrivateSnapshot` `Rack` only local, `PublicSnapshot` `Table` for all, no leak | `e2e/sync-visual.spec.ts` screenshot `Rack` `alice` vs `bob` different |
| 48 | E2E sync | `2 browsers` `joinMatch` `PrivateSnapshot` per seat `OwnRack` different, `PublicSnapshot` same | `e2e/sync-e2e.spec.ts` `alice ownRack 0` not `bob` |

**Milestone:** Client receives `PrivateSnapshot` on join/reconnect (only own rack) and `PublicSnapshot` for table, re-renders correctly, never shows foreign `OwnRack`, handles `RoundComplete` `Winner`.

---

## Phase 7 — Game Actions — Client → Server (Days 49–62)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 49 | Start match | `TurnIndicator` `Start` visible only if `Waiting` `OwnSeat==0` `Players>=2`, click `OpClientStart 1 {}` via `sendMatchState`, disable after `started` | `e2e/start.spec.ts` `Waiting` `Start` visible host `2p` |
| 50 | Opening discard | `RackScene` `discardSelected()` when `OpeningDiscard` `OwnSeat==CurrentSeat` `OwnRack 15` `OpClientDiscard 2 {tileId}` clear `selected` | `e2e/opening.spec.ts` `OpeningDiscard` `15→14` `IsOpeningDiscard` |
| 51 | Draw stock | `Draw` `Visible` only if `Playing MustDraw` `OwnSeat==CurrentSeat`, `OpClientDrawStock 3 {}` disable until `MeldOrDiscard` | `e2e/draw-stock.spec.ts` `MustDraw` `Draw` enabled, `MeldOrDiscard` disabled |
| 52 | Draw previous discard | `DrawPrevious` `HasOpened` `DiscardRow not empty` `!IsOpeningDiscard`, `OpClientDrawPreviousDiscard 4 {}` | `e2e/draw-prev.spec.ts` `HasOpened true` `disc-open` false `DrawPrev` enabled |
| 53 | Pickup discard for meld | `pickupSelected()` `selected 2` + `discardIndex` via `TableScene` click, `OpClientPickupDiscardForMeld 5 {discardIndex, tileIds:[2]}` `+jokerReps` | `e2e/pickup.spec.ts` `selected 2` `discardIndex 1` `Pickup` `TableMelds+1` |
| 54 | Normal discard | `discardSelected()` when `MeldOrDiscard` `OwnSeat==CurrentSeat` `OpClientDiscard 2 {tileId}` `CurrentSeat→(current+1)%n` | `e2e/normal-discard.spec.ts` `MeldOrDiscard` `Discard` `CurrentSeat 0→1` |
| 55 | Meld initial | `meldSelected("run")` when `!HasOpened` `selected>=3` `OpClientMeldInitial 6 {melds:[{id, kind:run, tileIds}]}` | `e2e/meld-initial.spec.ts` `!HasOpened` `selected 3` `run` `total>=50` |
| 56 | Meld new | when `HasOpened` `OpClientMeldNew 7 {melds:[...]}` `HasOpened` stays true | `e2e/meld-new.spec.ts` `HasOpened` `MeldNew` `HasOpened true` |
| 57 | Extend meld | `TableScene onMeldClicked(meldId)` highlight + `extendSelected(meldId)` `selected>=1` `HasOpened` `OpClientExtendMeld 8 {meldId, tileIds}` | `e2e/extend.spec.ts` `HasOpened` `selected 1` `Extend` `TableMeld` revalidated |
| 58 | Replace joker | `onJokerClicked(meldId,jokerId)` highlight + `replaceSelected(targetMeldId, tileId, new1, new2)` `OpClientReplaceJoker 9 {targetMeldId, tileId, newMeldTiles[2], jokerReps, newMeldKind}` | `e2e/replace.spec.ts` `Replace` `target` `tileId` `newMeldTiles 2` |
| 59 | Winner | `on PrivateSnapshot` `GamePhase=="RoundComplete"` `Winner` overlay `RackCount` `TableMelds` `DiscPlayer` disabled, `TableScene` `Winner: alice` | `e2e/winner.spec.ts` `RoundComplete` `Winner 0` `overlay` |
| 60 | Visual actions | `playwright` `start→opening discard→draw→meld 50+→extend→previous/pickup→replace→discard→win` via `MeldOrDiscard` loop | `e2e/actions-visual.spec.ts` screenshot `action` sequence |
| 61 | E2E actions | `2 browsers` `alice/bob` `draw`/`discard`/`meld`/`win` via `MatchLoop` `requestId` | `e2e/actions-e2e.spec.ts` `alice` `bob` `win` `RoundComplete` |
| 62 | Error on invalid | `OpServerError 102` `code=bad_payload` field highlight `already_opened` → `meld` button | `e2e/invalid.spec.ts` `DISCARD` `tileId not in rack` `OpServerError` |

**Milestone:** Player can `start` → `opening discard` → `draw` → `meld 50+ run` → `extend` → `previous/pickup` → `replace` → `discard` → `win` (`RoundComplete`) via `MeldOrDiscard` loop, all through `MatchLoop` with `requestId`.

---

## Phase 8 — UI Feedback & Validation — Client-Side (Days 63–72)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 63 | Turn indicator | `TurnIndicator` `CurrentSeat` `GamePhase/TurnPhase` `← current`, `Draw` disabled if not `MustDraw`/`CurrentSeat`, `Discard` disabled if not `MeldOrDiscard`/`CurrentSeat`/`selected!=1` | `e2e/turn.spec.ts` `turn` `← current` |
| 64 | HasOpened indicator | `TurnIndicator` `HasOpened` per `PublicPlayer`, disables `DrawPrevious`/`Pickup`/`Extend`/`Replace` if `!HasOpened` | `e2e/opened.spec.ts` `HasOpened false` buttons disabled |
| 65 | Error toast | `ErrorToast` `OpServerError` `code/message/details/requestId/op` `3s` `bg #dc2626` `data-error-code` | `e2e/error-toast.spec.ts` `LEAKED` `data-error-code` |
| 66 | Meld validation | `validateRun`/`validateSet` mirror `rules/meld/run.go:16`/`set.go:16` (`1-2-3` `12-13-1` `13-1-2` `joker` `real>=2*joker`) immediate feedback | `e2e/meld-validation.spec.ts` `1-2-3` valid `13-1-2` invalid |
| 67 | Scoring preview | `Rack` `scorePreview(selected)` `TotalScore` `2-9:5 10-13:10 Ace 5/10/25 Joker rep` `total>=50` `≥1 run` `MELD_INITIAL` enable | `e2e/scoring.spec.ts` `selected 3` `total 50` enabled `49` disabled |
| 68 | Duplicate tile | `RackScene` `selected Set` prevents duplicate `TileId`, toast `duplicate tile` | `e2e/duplicate.spec.ts` `selected` duplicate `LEAKED` |
| 69 | Discard highlight | `TableScene` highlights `IsOpeningDiscard` red `2px` + `Index` + disables `DrawPrevious`/`Pickup` if `discardRow[0].IsOpeningDiscard` `length==1` | `e2e/discard-highlight.spec.ts` `disc-open` red `2px` |
| 70 | Meld highlight | `TableScene` highlights `selectedMeldId` `run/set` `JokerReps ->colour-rank` | `e2e/meld-highlight.spec.ts` `selectedMeldId` `Kind` |
| 71 | Visual feedback | `playwright` `selected 3` valid `score 50` `turn` `HasOpened` `toast` | `e2e/feedback-visual.spec.ts` screenshot `feedback` |
| 72 | E2E feedback | `2 browsers` `selected` `HasOpened` `toast` `MELD_INITIAL 49` `OpServerError` | `e2e/feedback-e2e.spec.ts` `MELD_INITIAL 49` `bad_request` |

**Milestone:** Client gives immediate, non-authoritative feedback for `selected` validity and `Draw`/`Meld` enablement, but always shows server `OpServerError` as canonical.

---

## Phase 9 — Polish, Responsiveness & Accessibility (Days 73–95)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 73 | Typography | `Inter`/`JetBrains Mono` scale for `RackCount`/`StockCount`/`TurnIndicator` | `e2e/typography.spec.ts` `fontFamily Inter` |
| 74 | Colours | `red #e53935`/`yellow #f9a825`/`blue #1e88e5`/`black #212121` `green felt #0a4d2e` contrast WCAG AA | `e2e/colours.spec.ts` `tile red #e53935` |
| 75 | Spacing | `TableMelds y=100+i*80` `DiscardRow x=100+i*40` `Rack x=100+i*50` via `getSubspaces()` `GAME_SPACE 1000×1000` | `e2e/spacing.spec.ts` `MeldArea` `DiscardRowArea` |
| 76 | Buttons | `Phaser.Rectangle` + `Text` `setInteractive` `hover #228a3e` `disabled #2a2a2a` `alpha 0.35` | `e2e/buttons.spec.ts` `Draw` hover `#228a3e` |
| 77 | Modals | `MeldModal` `jokerReps` `colour`/`rank` dropdown for `jokerId` | `e2e/modal.spec.ts` `jokerReps` modal `colour` |
| 78 | Rack polish | `rack.png` 9-slice `5d4037` rounded `12` inner `6d4c41` slots `46×68` `6` gap | `e2e/rack-polish.spec.ts` `rack.png` `rounded` |
| 79 | Table polish | `table.png` green felt `stock pile` `44×22` `#4444aa` `222266` | `e2e/table-polish.spec.ts` `stock pile` `44×22` |
| 80 | Tiles polish | Per `colour/rank` `Container` `48×64` `white` `stroke 2` `rank 20px` `●` `Joly` | `e2e/tiles.spec.ts` `A` `K` `5` `J` |
| 81 | Animations | `tween` `Draw` `stock→rack` `200ms` `Back.easeOut` `prefers-reduced-motion` | `e2e/animations.spec.ts` `tween` `200ms` |
| 82 | Sound | `this.sound.add` `meld`/`discard`/`win` muted default `toggle` | `e2e/sound.spec.ts` `sound` muted |
| 83 | Responsive desktop | `1024×768` `Scale.FIT` `CENTER_BOTH` `1000×1000` centered | `e2e/responsive-desktop.spec.ts` `1000×1000` `FIT` |
| 84 | Responsive tablet | `768×1024` portrait `Rack y=650` `Meld` scroll | `e2e/responsive-tablet.spec.ts` `768×1024` `Rack y=650` |
| 85 | Responsive mobile | `375×667` `Rack` horizontal scroll `TableMelds` stacked `StockCount` hide | `e2e/responsive-mobile.spec.ts` `375×667` `Rack` scroll |
| 86 | A11y keyboard | `Phaser.Input.Keyboard` `1..9` select `Enter` discard/meld `Tab` switch seat (local 2p) | `e2e/keyboard.spec.ts` `press 1` select |
| 87 | A11y screen-reader | `aria-label` `canvas` wrapper + `Text` `name:Tile red-5 id red-05-1` `axe` | `e2e/a11y.spec.ts` `axe` `0 violations` |
| 88 | A11y colour-blind | `R`/`Y`/`B`/`K` text labels + `black` outline not colour alone | `e2e/colour-blind.spec.ts` `R` `Y` `B` `K` |
| 89 | A11y reduced motion | `matchMedia("(prefers-reduced-motion: reduce)")` disable `tween` | `e2e/reduced-motion.spec.ts` `prefers-reduced-motion` |
| 90 | Localization | `locales/en.json` `locales/ro.json` `t("draw")` `ro` `suita/terta/Joly` | `e2e/localization.spec.ts` `ro` `suita` |
| 91 | Pluralization | `Intl.PluralRules` `1 tile` vs `5 tiles` `RackCount` | `e2e/plural.spec.ts` `1 tile` `5 tiles` |
| 92 | Connection UX | `reconnecting` overlay `socket.onDisconnect` `resync` spinner until `OpServerState 100` `PrivateSnapshot` | `e2e/connection.spec.ts` `reconnecting` overlay |
| 93 | Visual polish | `playwright` `Inter` `colours` `spacing` `buttons` `rack` `table` `tiles` `animations` | `e2e/polish-visual.spec.ts` screenshot `polish` |
| 94 | A11y e2e | `2 browsers` `keyboard` `screen-reader` `colour-blind` | `e2e/a11y-e2e.spec.ts` `keyboard` `axe` |
| 95 | E2E polish | `2 browsers` `responsive` `localization` `connection` | `e2e/polish-e2e.spec.ts` `responsive` |

**Milestone:** Client is usable with keyboard, screen-reader, colour-blind, mobile, subtle animations and no private leak, with `Inter`/`JetBrains Mono` polished.

---

## Phase 10 — Testing, Docs & Release (Days 96–115)

| Day | Focus | Deliverable | Playwright |
|---:|---|---|---|
| 96 | Unit: tile | `src/state/snapshot.test.ts` `TileInstance` `Colour`/`Rank` `IsJoker` `Validate` mirror `tile_test.go` | `vitest` `TileInstance` |
| 97 | Unit: meld | `src/rules/meld.test.ts` `ValidateSet`/`ValidateRun` `1-2-3`/`12-13-1`/`13-1-2`/`joker`/`real>=2*joker` mirror `matrix_test.go` | `vitest` `1-2-3` |
| 98 | Unit: scoring | `src/rules/scoring.test.ts` `ScoreTile` `5/10/25` `joker` `TotalScore 30/90` `ValidateInitialBatch 49/50` | `vitest` `5/10/25` |
| 99 | Integration: sync | `src/state/sync.test.ts` `PublicView`/`PrivateView` `Version 1` `reconnection` `PrivateSnapshot` per seat `redaction` | `vitest` `redaction` |
| 100 | Integration: flow | `src/state/sync.test.ts` `TestDeterministicSimulation` `MatchInit→Start→OpeningDiscard→Draw→InitialMeld→Extend→Previous→Pickup→Win` `CheckTileConservation 106` | `vitest` `deterministic` |
| 101 | E2E: auth | `e2e/auth.spec.ts` `authenticateDevice` `defaultkey` `127.0.0.1:7350` `token` `userId` | `e2e/auth.spec.ts` `token` |
| 102 | E2E: create/join | `e2e/match.spec.ts` `createMatch`/`joinMatch` `alice`/`bob` `OpServerState 100` `PrivateSnapshot` per seat | `e2e/match.spec.ts` `alice` `bob` `matchId` |
| 103 | E2E: draw/discard | `e2e/draw-discard.spec.ts` `DRAW_STOCK` `DISCARD` `CurrentSeat 0→1` `anticlockwise` | `e2e/draw-discard.spec.ts` `CurrentSeat 0→1` |
| 104 | E2E: melding | `e2e/meld.spec.ts` `MELD_INITIAL 50+ run` `MELD_NEW` `EXTEND` `HasOpened` `Kind` | `e2e/meld.spec.ts` `MELD_INITIAL` `50+` |
| 105 | E2E: reconnection | `e2e/reconnect.spec.ts` `socket.disconnect`+`connect`+`joinMatch` `PrivateSnapshot` `OwnRack` `Winner` | `e2e/reconnect.spec.ts` `reconnect` |
| 106 | E2E: error states | `e2e/error.spec.ts` `DISCARD` `tileId not in rack` `MELD_INITIAL 49` `EXTEND` wrong colour `OpServerError 102` | `e2e/error.spec.ts` `OpServerError` |
| 107 | Docs: README | `client/README.md` `Next Steps` to `Day 96–106` `How to test` `npm run test` `npm run e2e` `make client-*` | `e2e/docs.spec.ts` `README` `npm run test` |
| 108 | Docs: roadmap | `client/docs/roadmap.md` Phase 10 `Day 96–106` `✅ Done` `Current State` | `e2e/roadmap.spec.ts` `roadmap` `✅ Done` |
| 109 | Visual regression | `playwright` `visual` `rack` `table` `melds` `discard` `stock` `turn` `error toast` `winner` | `e2e/visual.spec.ts` `visual` |
| 110 | E2E full | `playwright` `2 browsers` `full` `auth` `create/join` `draw/discard` `meld` `reconnect` `error` | `e2e/full.spec.ts` `full` |
| 111 | Release candidate | `client-mvp-rc1` tag when `npm run typecheck`+`test`+`e2e`+`build` green + `go test` `no leak` | `git tag client-mvp-rc1` |
| 112 | Release notes | `client/docs/release-notes.md` `Stack: Phaser 3.80+ Vite 5 nakama-js 2.9 TS 5 Go rummy-mvp-rc1` `make dev` `known limitations` | `e2e/release.spec.ts` `release-notes` |
| 113 | Docs: testing | `docs/testing.md` `unit` `integration` `simulation` `Docker` `e2e` | `e2e/testing.spec.ts` `testing.md` |
| 114 | Visual docs | `playwright` `docs` screenshots | `e2e/docs-visual.spec.ts` `docs` |
| 115 | E2E docs | `playwright` `docs` `README` `roadmap` | `e2e/docs-e2e.spec.ts` `docs` |

**Milestone:** Two local users can `npm run dev` at `http://127.0.0.1:5173`, `create/join` via `defaultkey`, see `PrivateView.OwnRack` vs `PublicView` counts, `draw`/`discard`/`meld`/`extend`/`pickup`/`replace`/`win` through `MatchLoop`, `OpServerError` toasts, `go test` still `no leak`, `client-mvp-rc1` tagged.

---

## Extended Roadmap Summary (Client Future)

| Phase | Days | Primary Outcome | Playwright |
|---:|---|---|---|
| 1 | 1–7 | Foundation | `smoke` `scaffold` `nakama` `lint` |
| 2 | 8–14 | Assets | `preload` `assets` `table` `rack` |
| 3 | 15–22 | Rendering (static) | `rack` `sort` `melds` `discard` `stock` `layout` |
| 4 | 23–29 | Input (local) | `selection` `drag` `drop` `discard` `meld` `keyboard` |
| 5 | 30–38 | Networking | `auth` `socket` `create` `join` `envelope` `send` `receive` `error` |
| 6 | 39–48 | Sync (authoritative) | `snapshot` `private` `public` `redaction` `store` `reconnect` `version` `state-machine` |
| 7 | 49–62 | Actions | `start` `opening` `draw` `prev` `pickup` `discard` `meld` `extend` `replace` `winner` |
| 8 | 63–72 | Feedback | `turn` `opened` `toast` `validation` `scoring` `duplicate` `discard` `meld` |
| 9 | 73–95 | Polish | `typography` `colours` `spacing` `buttons` `modal` `rack` `table` `tiles` `animations` `sound` `responsive` `a11y` |
| 10 | 96–115 | Testing & release | `unit` `integration` `e2e` `docs` `rc1` |

By **Day 115** (client), the project evolves from empty `client/` to a playable, accessible, localized Phaser web client that proves the protocol `Version 1` and can be run via `npm run dev` against `docker compose up --build -d` without oral knowledge. Execution remains Handmade Hero: one vertical slice per day, `npm run dev` green, focused commit, push, **Playwright after each slice**.

---

## Daily Definition of Done

Every day must end with:

1. A small, working change (no dump).
2. New or updated **Playwright** `e2e/*.spec.ts` covering that slice (plus `vitest` where pure).
3. `npm run lint` `npm run typecheck` `npm run test` `npm run build` `npx playwright test` green.
4. `docs` update if behavior/protocol changed.
5. Focused `git commit -m "feat(client): ... (Day N)"` + `git push`.
6. Screenshot or `trace` for visual slices.

## Commit Pattern

```
feat(client): add Vite + TS + Phaser scaffold (Day 2)
feat(client): add tile and joker sprites (Day 9)
feat(client): add rack rendering (Day 15)
feat(client): handle private snapshot with Rack re-render (Day 40)
test(client): add modern UI e2e and fix rack size expectation (UX)
fix(client): unify to GAME_SPACE layout and fix alignment/drag-drop
```

---

*This roadmap lists days that **will be** implemented for the Phaser client. For what **is** implemented on the backend, see `docs/IMPLEMENTED.md`; for backend roadmap, see `docs/roadmap.md`. Last updated: 2026-08-26 for `client/` empty → Day 1 start with Playwright.*
