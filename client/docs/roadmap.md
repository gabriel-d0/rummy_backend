# Rummy Client — Phaser 3 Incremental Roadmap

**Target:** Minimal, maintainable, *playable* Romanian Tile Rummy web client in Phaser 3 that proves the authoritative Go backend protocol (`docs/protocol.md` `Version 1`, `internal/protocol/opcodes.go:8` `1..9`/`100..199`) without copying Remi Online branding. Follows the same *Handmade Hero* incremental style as the backend: one small vertical slice per day, `make check` + `make dev` green before next day, one focused commit per day.

**Backend prerequisite:** `rummy-mvp-rc1` at `main@cfce62e` (82 commits `36c2c59..41e794b` + `rummy-mvp-rc1` tag) is stable. Backend runs at `127.0.0.1:7350` (`defaultkey`) via `docker compose up --build -d` + `make smoke` `SMOKE PASSED`. Client must never leak private racks: `PublicView` only `RackCount`/`StockCount`/`DiscardRow`/`TableMelds`/`Winner`, `PrivateView` adds `OwnRack` per `Seat` only (`internal/match/visibility.go:36` `Version 1`).

**How to read:** Day numbers are *client-only* (Day 1 is first client commit, independent of backend days 1–24). Each day ends in a runnable, tested state with `npm run dev` showing a working vertical slice.

---

## Phase 1 — Foundation and local development

| Day | Focus | Deliverable |
|---:|---|---|
| 1 | Repository and tooling audit | Review `rummy_backend` repo, `client/` empty, backend `docs/protocol.md`/`state-machine.md`, `AGENTS.md:812` Day 24 acceptance, record in `client/docs/roadmap.md` and `client/README.md`. |
| 2 | Vite + TypeScript + Phaser scaffold | `npm init` + `phaser@3.80` + `vite@5` + `typescript@5`, `vite.config.ts` with `server.port 5173` `host 127.0.0.1`, `index.html` with `<div id="game">`, `src/main.ts` with `new Phaser.Game({type: Phaser.AUTO, parent: "game", width: 1024, height: 768, scene: [Preload]})`, `npm run dev` shows `Phaser 3` `Init` log. |
| 3 | Nakama JS client dependency | Add `@heroiclabs/nakama-js@2.9`, create `src/net/nakama.ts` with `Client("127.0.0.1","7350","defaultkey")` and `authenticateDevice` via `deviceId` `uuid`, store `token` in `localStorage`, `socket` lazy init. `npm run build` green. |
| 4 | Project structure and lint/format | Add `src/scenes` `src/net` `src/state` `src/ui` dirs, `eslint` flat + `prettier`, `tsconfig.json` `strict`, `npm run lint` + `npm run fmt` + `npm run typecheck` scripts, `make` alias `make client-lint` if needed, `gofmt -l` still clean for backend. |
| 5 | Developer scripts and docs | Add `client/README.md` `Quick Start` (`npm install`, `npm run dev`, `make cli` backend CLI vs this Phaser client), `client/.env.example` (`VITE_NAKAMA_HOST=127.0.0.1` `VITE_NAKAMA_PORT=7350` `VITE_NAKAMA_KEY=defaultkey`), update root `README` `Next Steps` to point to `client/docs/roadmap.md`. |

**Milestone:** `npm run dev` shows a blank Phaser canvas at `http://127.0.0.1:5173` and `nakama.ts` can `authenticateDevice` against `docker compose up` backend.

---

## Phase 2 — Phaser scene and asset pipeline

| Day | Focus | Deliverable |
|---:|---|---|
| 6 | Preload scene | `src/scenes/Preload.ts` with `preload()` that loads a single 1x1 tile sprite (placeholder `assets/tile.png` 32x48) and `create()` that starts `TableScene`, `npm run dev` shows `Preload` `complete` in console. |
| 7 | Tile sprite | Add `assets/tile.png` for `red-1` and `assets/joker.png` for `Joly`, `Preload` loads both, `TableScene` renders one `this.add.image(100,100,"tile")`, `npm run dev` shows a tile. |
| 8 | Table background | Add `assets/table.png` (green felt 1024x768), `TableScene` renders `this.add.image(512,384,"table")` behind tiles. |
| 9 | Rack background | Add `assets/rack.png` (brown 800x120), `RackScene` renders `this.add.image(512,680,"rack")` as UI layer, `TableScene` launches `RackScene` as parallel scene. |
| 10 | Asset manifest | Create `src/scenes/Preload.ts` manifest object `{tile: "assets/tile.png", joker: "assets/joker.png", table: "assets/table.png", rack: "assets/rack.png"}` and test `preload` loads all 4 without 404. |

**Milestone:** Phaser shows a static table with one tile and a rack, no networking yet.

---

## Phase 3 — Rack and table rendering (static)

| Day | Focus | Deliverable |
|---:|---|---|
| 11 | Rack rendering (static) | `src/ui/Rack.ts` with `renderRack(tiles: TileInstance[], seat: Seat)` that draws 14 `this.add.image` per tile at `x = 100 + i*50`, `RackScene` calls it with mock `PrivateView.OwnRack` (2 tiles `red-5`, `blue-7`), `npm run dev` shows rack. |
| 12 | Rack sorting | Add `sortRack(tiles)` by `Colour` then `Rank` (`tile.Rank 1..13`), `RackScene` sorts before render, test with `red-13`, `red-1`, `blue-5` → `red-1`, `red-13`, `blue-5`. |
| 13 | Table melds rendering (static) | `src/scenes/TableScene.ts` `renderTableMelds(melds: TableMeld[])` that draws each `TableMeld` at `y = 100 + i*80` with `Tiles` as `this.add.image` per tile, mock `PublicView.TableMelds` with 1 run `1-2-3 red` and 1 set `7 red/yellow/blue`, `npm run dev` shows melds. |
| 14 | Discard row rendering (static) | `TableScene` `renderDiscardRow(discardRow: DiscardEntry[])` that draws `DiscardRow` at `x = 100 + i*40` `y = 300`, mock `DiscardRow` with `IsOpeningDiscard` flagged `disc-open` in red border. |
| 15 | Stock and turn indicator (static) | `TableScene` `renderStockCount(count: number)` and `TurnIndicator` at `x=800 y=50` showing `CurrentSeat` and `GamePhase`/`TurnPhase` from mock `PublicView`, `npm run dev` shows `Stock: 77` and `Current: seat-0`. |

**Milestone:** Client renders a static hand, table, discard, stock, and turn from mock `PublicView`/`PrivateView` without networking.

---

## Phase 4 — Input handling (local only)

| Day | Focus | Deliverable |
|---:|---|---|
| 16 | Tile selection | `src/ui/Rack.ts` `onTileClicked(tileId)` toggles `selected: Set<TileInstanceId>` and tints selected tile `0xffff00`, `RackScene` logs `selected` on click. |
| 17 | Drag start | Add `this.input.setDraggable(tileImage)` per tile, `RackScene` `dragstart` logs `tileId`, no drop yet. |
| 18 | Drag drop to table (local) | `TableScene` `drop` zone at `y=200` logs `tileId` dropped, no server call yet. |
| 19 | Discard selection | `RackScene` `discardSelected()` validates exactly 1 `selected` and logs `DISCARD {tileId}`, no server call yet. |
| 20 | Meld selection | `RackScene` `meldSelected(kind)` validates `selected.size >=3` and logs `MELD_INITIAL {kind, tileIds}`, no server call yet. |

**Milestone:** Player can click/drag tiles and see logs for `discard`/`meld` intents, no validation yet.

---

## Phase 5 — Networking: Nakama JS and envelope

| Day | Focus | Deliverable |
|---:|---|---|
| 21 | Device auth | `src/net/nakama.ts` `authenticate()` does `client.authenticateDevice(deviceId, create=true, username)` via `defaultkey` at `127.0.0.1:7350`, stores `token` and `userId`, `npm run dev` console shows `authenticated alice` with `userId`. |
| 22 | Socket connect | `src/net/nakama.ts` `connectSocket(token)` does `client.createSocket(false,false)` + `socket.connect(session, true)` with `appearOnline=false`, logs `socket connected`. |
| 23 | Create match | `src/net/nakama.ts` `createMatch(socket)` does `socket.createMatch()` and logs `matchId`, stores `matchId` in `src/state/sync.ts`. |
| 24 | Join match | `src/net/nakama.ts` `joinMatch(socket, matchId)` does `socket.joinMatch(matchId)` and logs `joined`, handles `match_presence_event`. |
| 25 | Envelope | `src/net/protocol.ts` `Envelope{v,op,requestId,payload}` with `Version 1`, `OpClientStart 1`/`Discard 2`/`DrawStock 3`/`DrawPrevious 4`/`Pickup 5`/`MeldInitial 6`/`MeldNew 7`/`Extend 8`/`ReplaceJoker 9`, `OpServerState 100`/`StatePublic 101`/`Error 102`/`Event 103`, `NewEnvelope(op,payload)` helper. |
| 26 | Send match state | `src/net/protocol.ts` `sendMatchState(socket, matchId, op, payload)` does `socket.sendMatchState(matchId, op, JSON.stringify(Envelope{v,op,requestId,payload}))`, logs `sent op 2`. |
| 27 | Receive match state | `src/net/nakama.ts` `socket.onmatchdata` parses `Envelope` and routes `op 100`/`101`/`102`/`103` to `src/state/sync.ts` handlers, logs `received op 102`. |
| 28 | Error display (network) | `src/ui/ErrorToast.ts` shows `OpServerError 102` `code`/`message`/`requestId`/`op` as red toast for 3s, logs `not_your_turn` etc. |

**Milestone:** Client can `authenticateDevice`, `createSocket`, `createMatch`/`joinMatch`, send `Envelope` and receive `OpServerError`/`OpServerEvent` via WebSocket.

---

## Phase 6 — State synchronization (authoritative)

| Day | Focus | Deliverable |
|---:|---|---|
| 29 | Public snapshot types | `src/state/snapshot.ts` `PublicSnapshot{Version,GamePhase,TurnPhase,CurrentSeat,Players{RackCount,HasOpened},StockCount,DiscardRow,TableMelds,Winner}` and `PrivateSnapshot{PublicSnapshot,OwnRack,OwnSeat}` mirroring `internal/match/visibility.go:36`. |
| 30 | Private snapshot handling | `src/state/sync.ts` `onPrivateSnapshot(snap)` stores `PrivateSnapshot` per `OwnSeat`, `RackScene` re-renders `OwnRack` only, `TableScene` re-renders `PublicSnapshot` parts. |
| 31 | Public snapshot handling | `src/state/sync.ts` `onPublicSnapshot(snap)` stores `PublicSnapshot`, `TableScene` re-renders `TableMelds`/`DiscardRow`/`StockCount`/`CurrentSeat`, not `OwnRack`. |
| 32 | Redaction check (client) | `src/state/sync.ts` `checkNoLeak(publicJson)` does `json.stringify(PublicSnapshot)` string search for `OwnRack` IDs as in `visibility_test.go`, logs `no leak` or `LEAKED` in console. |
| 33 | Reconnection: store last snapshot | `src/state/sync.ts` stores `lastPrivateSnapshot` per `Seat` in `localStorage` or memory, on `socket.onDisconnect` keep `matchId` and `userId`. |
| 34 | Reconnection: rejoin | `src/net/nakama.ts` `reconnect()` does `socket.connect` + `joinMatch(matchId)` and expects `OpServerState 100` `PrivateSnapshot` for that `Seat` only (per `rummy_match.go:79`), `sync.ts` rehydrates `RackScene` from `OwnRack`, not from old. |
| 35 | Versioning | `snapshot.ts` `Version 1` check: if `snap.Version !== 1` log `bad_version` and ignore, as in `parser.go:22`. |
| 36 | State machine client | `src/state/sync.ts` `AllowedOps` mirror `phases.go:15` (e.g., `Waiting→Start` only host Seat 0, `OpeningDiscard→Discard` only Seat 0, `Playing MustDraw→DrawStock/DrawPrevious/Pickup`, `MeldOrDiscard→Discard/MeldInitial/MeldNew/Extend/ReplaceJoker`, `RoundComplete→none`), client disables buttons per `CurrentSeat`/`TurnPhase`/`HasOpened`. |

**Milestone:** Client receives `PrivateSnapshot` on join/reconnect (only own rack) and `PublicSnapshot` for table, re-renders correctly, never shows foreign `OwnRack`, and handles `RoundComplete` `Winner`.

---

## Phase 7 — Game actions (client → server)

| Day | Focus | Deliverable |
|---:|---|---|
| 37 | Start match | `src/ui/TurnIndicator.ts` `Start` button visible only if `GamePhase==Waiting` and `OwnSeat==0` and `Players.length>=2`, click sends `OpClientStart 1` `{}` via `sendMatchState`, disables after `started` `OpServerEvent`. |
| 38 | Opening discard | `RackScene` `discardSelected()` when `GamePhase==OpeningDiscard` and `OwnSeat==CurrentSeat` and `OwnRack.length==15` sends `OpClientDiscard 2` `{"tileId": selectedId}`, clears `selected`, waits for `OpServerEvent` `phase: Playing`. |
| 39 | Draw stock | `src/ui/Rack.ts` `Draw` button visible only if `GamePhase==Playing` and `TurnPhase==MustDraw` and `OwnSeat==CurrentSeat`, click sends `OpClientDrawStock 3` `{}`, disables until `MeldOrDiscard`. |
| 40 | Draw previous discard | `DrawPrevious` button visible only if `HasOpened` and `DiscardRow` not empty and not `IsOpeningDiscard`, click sends `OpClientDrawPreviousDiscard 4` `{}`, `RackScene` adds one tile. |
| 41 | Pickup discard for meld | `RackScene` `pickupSelected()` validates `selected.size==2` and `DiscardRow` index `discardIndex` selected via `TableScene` discard click, sends `OpClientPickupDiscardForMeld 5` `{"discardIndex":2,"tileIds":["id1","id2"]}` (plus `jokerReps` if selected includes `IsJoker`), `sync.ts` expects `TableMelds+1` and `DiscardRow` truncated. |
| 42 | Normal discard | `RackScene` `discardSelected()` when `Playing/MeldOrDiscard` and `OwnSeat==CurrentSeat` and `OwnRack.length>=1` sends `OpClientDiscard 2` `{"tileId": selectedId}`, `CurrentSeat` should become `(current+1)%n` via `OpServerEvent` `nextSeat`. |
| 43 | Meld initial | `RackScene` `meldSelected("run")` when `!HasOpened` and `selected.size>=3` sends `OpClientMeldInitial 6` `{"melds":[{"id":"cli-run-0","kind":"run","tileIds": selectedIds}]}`, validates `Kind` `run`/`set` and `real>=2*joker` locally but server is authoritative, shows `OpServerError` `bad_request` if `total>=50` fails. |
| 44 | Meld new | `RackScene` `meldSelected` when `HasOpened` sends `OpClientMeldNew 7` `{"melds":[{...}]}` with `id` `cli-<kind>-<n>` `kind` `tileIds`, server validates each `ValidateRun/Set` and `meldId` not colliding, `HasOpened` stays true. |
| 45 | Extend meld | `TableScene` `onMeldClicked(meldId)` highlights, `RackScene` `extendSelected(meldId)` validates `selected.size>=1` and `HasOpened`, sends `OpClientExtendMeld 8` `{"meldId":"...","tileIds": selectedIds}` (plus `jokerReps` if needed), server revalidates entire resulting meld. |
| 46 | Replace joker | `TableScene` `onJokerClicked(meldId, jokerId)` highlights, `RackScene` `replaceSelected(targetMeldId, tileId, new1, new2)` sends `OpClientReplaceJoker 9` `{"targetMeldId":"...","tileId":"...","newMeldTiles":["...","..."],"jokerReps":{"jokerId":{"colour":"red","rank":7}},"newMeldKind":"run"}` per `validator.go:142` and `docs/rules-decisions.md:6.4`. |
| 47 | Winner and round complete | `src/state/sync.ts` on `OpServerState 100` with `GamePhase=="RoundComplete"` shows `Winner` overlay and `Final DiscardRow`/`TableMelds`/`RackCount`, disables all `MeldOrDiscard` buttons, `TableScene` shows `Winner: alice`. |

**Milestone:** Player can `start` → `opening discard` → `draw` → `meld` (50+ run) → `extend` → `previous`/`pickup` → `replace` → `discard` → `win` (`RoundComplete`) via `MeldOrDiscard` loop, all through `MatchLoop` with `requestId` correlation.

---

## Phase 8 — UI feedback and validation (client-side)

| Day | Focus | Deliverable |
|---:|---|---|
| 48 | Turn indicator | `src/ui/TurnIndicator.ts` shows `CurrentSeat` `GamePhase`/`TurnPhase` with `← current` marker, disables `Draw` if not `MustDraw` or not `CurrentSeat`, disables `Discard` if not `MeldOrDiscard` or not `CurrentSeat` or `selected.size!=1`. |
| 49 | HasOpened indicator | `TurnIndicator` shows `HasOpened` per `PublicPlayer` and disables `DrawPrevious`/`Pickup`/`Extend`/`Replace` if `!HasOpened`, and disables `MeldInitial` if `HasOpened` else `MeldNew`. |
| 50 | Error toast | `src/ui/ErrorToast.ts` renders `OpServerError` `code`/`message`/`details`/`requestId`/`op` as in `protocol/errors.go:12`, with `bad_payload` field highlight and `already_opened`/`not_opened` mapping to `meld` button. |
| 51 | Meld validation (client) | `src/state/snapshot.ts` `validateRun`/`validateSet` mirror `rules/meld/run.go:16`/`set.go:16` for immediate feedback (e.g., `1-2-3` valid low, `12-13-1` valid high, `13-1-2` invalid), but never trust; server error still canonical. |
| 52 | Scoring preview (client) | `src/ui/Rack.ts` `scorePreview(selected)` shows `TotalScore` via `ScoreRun`/`ScoreSet` logic (2–9:5, 10–13:10, Ace 5/10/25, Joker rep) per `scoring/scoring.go:16`, with `total>=50` and `≥1 run` check for `MELD_INITIAL` button enable. |
| 53 | Duplicate tile check (client) | `RackScene` `selected` `Set` prevents duplicate `TileId` across meld batch, shows `duplicate tile` toast. |
| 54 | Discard row highlight | `TableScene` highlights `IsOpeningDiscard` with red border and `Index` label, disables `DrawPrevious`/`Pickup` if `discardRow[0].IsOpeningDiscard` and `discardRow.length==1`. |
| 55 | Table meld highlight | `TableScene` highlights `selectedMeldId` on `TableMeld` click for `Extend` target, shows `Kind` `run`/`set` and `JokerReps` `->colour-rank`. |

**Milestone:** Client gives immediate, non-authoritative feedback for `selected` validity and `Disc`/`Meld` enablement, but always shows server `OpServerError` as canonical.

---

## Phase 9 — Polish, responsiveness, and accessibility

| Day | Focus | Deliverable |
|---:|---|---|
| 56 | Typography | Establish typography scale for `RackCount`/`StockCount`/`TurnIndicator` (Phaser `Text` with `fontFamily` `Inter`/`monospace`). |
| 57 | Colours | Define palette: `red`/`yellow`/`blue`/`black` tile colours plus `green` table, ensure contrast per `docs/rules-decisions.md:1.1`. |
| 58 | Spacing | Define spacing scale for `TableMelds` `y = 100 + i*80`, `DiscardRow` `x = 100 + i*40`, `Rack` `x = 100 + i*50`. |
| 59 | Buttons | Style `Draw`/`Discard`/`Meld`/`Extend`/`Pickup`/`Replace` buttons with `Phaser.GameObjects.Rectangle` + `Text` and `setInteractive`. |
| 60 | Modals | Add `MeldModal` for `jokerReps` selection (when `selected` contains `IsJoker`, show `colour`/`rank` dropdown for that `jokerId`). |
| 61 | Rack background polish | Replace placeholder `rack.png` with 9-slice `Rack` with shadow and `selected` tint `0xffff00`. |
| 62 | Table background polish | Replace `table.png` with green felt texture and `Stock` pile visualization (`StockCount` as stacked tiles). |
| 63 | Tile sprites polish | Replace `tile.png`/`joker.png` with per-colour/rank sprites (or `Text` fallback for MVP: `red-5` with `red` tint). |
| 64 | Animations | Add `tween` for `Draw` (stock→rack) and `Discard` (rack→discardRow) with `duration 200ms`, respect `prefers-reduced-motion`. |
| 65 | Sound | Add optional `meld`/`discard`/`win` sound via `this.sound.add`, muted by default, toggle button. |
| 66 | Responsive: desktop | Ensure `1024x768` is centered with `scale: {mode: Phaser.Scale.FIT, autoCenter: Phaser.Scale.CENTER_BOTH}`. |
| 67 | Responsive: tablet | Test `768x1024` portrait with `Rack` at `y=650` and `TableMelds` scroll. |
| 68 | Responsive: mobile | Test `375x667` with `Rack` horizontal scroll and `TableMelds` stacked, hide `StockCount` text if needed. |
| 69 | Accessibility: keyboard | Add `Phaser.Input.Keyboard` for `1`–`9` to select rack tiles, `Enter` to `discard`/`meld`, `Tab` to `switch` seat for local 2-player test. |
| 70 | Accessibility: screen-reader | Add `aria-label` to `canvas` wrapper `div` and `Text` objects with `name` `Tile red-5 id red-05-1` etc., test with `axe`. |
| 71 | Accessibility: colour-blind | Ensure tile colours have text labels `R`/`Y`/`B`/`K` plus `black` outline, not colour alone. |
| 72 | Accessibility: reduced motion | Respect `window.matchMedia("(prefers-reduced-motion: reduce)")` to disable `tween`. |
| 73 | Localization foundation | Externalize UI text to `locales/en.json` and `locales/ro.json` with `t("draw")` helper, support `ro` for `Remi` terms (`suita`/`terta`/`Joly`). |
| 74 | Pluralization | Handle `RackCount` plural (`1 tile` vs `5 tiles`) per `Intl.PluralRules`. |
| 75 | Connection UX | Add `reconnecting` overlay when `socket.onDisconnect`, show `resync` spinner until `OpServerState 100` `PrivateSnapshot` rehydrates `RackScene`. |

**Milestone:** Client is usable with keyboard, screen-reader, colour-blind, and mobile, with subtle animations and no private leak.

---

## Phase 10 — Testing, docs, and release

| Day | Focus | Deliverable |
|---:|---|---|
| 76 | Unit tests: tile | Add `src/state/snapshot.test.ts` for `TileInstance` `Colour`/`Rank` `IsJoker` `Validate` (mirror `internal/rules/tile/tile_test.go`). |
| 77 | Unit tests: meld | Add `src/rules/meld.test.ts` for `ValidateSet`/`ValidateRun` `1-2-3`/`12-13-1`/`13-1-2`/`joker`/`real>=2*joker` (mirror `matrix_test.go`). |
| 78 | Unit tests: scoring | Add `src/rules/scoring.test.ts` for `ScoreTile` `5/10/25` and `joker` rep, `TotalScore` 30/90, `ValidateInitialBatch` 49 vs 50. |
| 79 | Integration: state sync | Add `src/state/sync.test.ts` for `PublicView`/`PrivateView` versioned `1`, `reconnection` `PrivateSnapshot` per seat, `redaction` `PublicView` never contains `OwnRack` IDs. |
| 80 | Integration: match flow | Add `src/state/sync.test.ts` `TestDeterministicSimulation` equivalent: `MatchInit`→`Start`→`OpeningDiscard`→`Draw`→`InitialMeld`→`Extend`→`Previous`→`Pickup`→`Win` with `CheckTileConservation` style `RackCount+StockCount+DiscardRow+TableMelds` `106`. |
| 81 | E2E: auth | Add `e2e/auth.spec.ts` with `playwright` that does `authenticateDevice` via `defaultkey` at `127.0.0.1:7350` and asserts `token` and `userId`. |
| 82 | E2E: create/join | Add `e2e/match.spec.ts` that does `createMatch`/`joinMatch` for `alice`/`bob` and asserts `OpServerState 100` `PrivateSnapshot` per seat with `OwnRack` and `PublicView` no leak. |
| 83 | E2E: draw/discard | Add `e2e/draw_discard.spec.ts` that does `DRAW_STOCK` then `DISCARD` and asserts `CurrentSeat` advances `0→1` anticlockwise `(current+1)%n`. |
| 84 | E2E: melding | Add `e2e/meld.spec.ts` that does `MELD_INITIAL` 50+ with run and `MELD_NEW` and `EXTEND` and asserts `TableMelds` `Kind` stable and `HasOpened`. |
| 85 | E2E: reconnection | Add `e2e/reconnect.spec.ts` that does `socket.disconnect` then `socket.connect`+`joinMatch` and asserts `PrivateSnapshot` rehydrated with `OwnRack` and `Winner` if `RoundComplete`. |
| 86 | E2E: error states | Add `e2e/error.spec.ts` that does invalid `DISCARD` `tileId` not in rack, `MELD_INITIAL` 49, `EXTEND` wrong colour, asserts `OpServerError 102` `code`/`message`/`requestId`/`op`. |
| 87 | Docs: client README | Update `client/README.md` `Next Steps` to reflect `Day 76–86` done, add `How to test` (`npm run test` `npm run e2e` `make client-test`). |
| 88 | Docs: roadmap update | Mark `client/docs/roadmap.md` Phase 10 `Day 76–86` as `✅ Done` with commit hashes, update `Current State` to after Day 86. |
| 89 | Release candidate tag | Create `client-mvp-rc1` tag when `npm run typecheck` + `npm run test` + `npm run e2e` + `npm run build` green and `go run ./cmd/rummy-cli` manual two-user flow still `no private leak` as in backend `rummy-mvp-rc1`. |
| 90 | Release notes | Write `client/docs/release-notes.md` with `Stack: Phaser 3.80+, Vite 5, nakama-js 2.9+, TypeScript 5, Go backend rummy-mvp-rc1 at 127.0.0.1:7350`, `make dev` health check, and `known limitations` (no `Doubla`, no prod deployment). |

**Milestone:** Two local users can run `npm run dev` at `http://127.0.0.1:5173`, create/join a match via `defaultkey`, see `PrivateView.OwnRack` vs `PublicView` counts, `draw`/`discard`/`meld`/`extend`/`pickup`/`replace`/`win` through the authoritative `MatchLoop`, see `OpServerError` toasts, and `go run ./cmd/rummy-cli` still works for headless protocol validation.

---

## Extended Roadmap Summary (Client Future)

| Phase | Days | Primary Outcome | Status |
|---:|---|---|---|
| 1 | 1–5 | Foundation and local development | ⏳ Not started (this file is the plan) |
| 2 | 6–10 | Phaser scene and asset pipeline | ⏳ Not started |
| 3 | 11–15 | Rack and table rendering (static) | ⏳ Not started |
| 4 | 16–20 | Input handling (local only) | ⏳ Planned |
| 5 | 21–28 | Networking: Nakama JS and envelope | ⏳ Planned |
| 6 | 29–36 | State synchronization (authoritative) | ⏳ Planned |
| 7 | 37–47 | Game actions (client → server) | ⏳ Planned |
| 8 | 48–55 | UI feedback and validation (client-side) | ⏳ Planned |
| 9 | 56–75 | Polish, responsiveness, accessibility | ⏳ Planned |
| 10 | 76–90 | Testing, docs, and release | ⏳ Planned |

By **Day 90** (client), the project evolves from the current backend `rummy-mvp-rc1` into a playable, accessible, localized Phaser web client that proves the protocol and can be run via `npm run dev` against `docker compose up --build -d` without oral knowledge. Execution remains Handmade Hero: one vertical slice per day, `npm run dev` green, focused commit, push.

---

*This roadmap lists days that **will be** implemented for the Phaser client. For what **is** implemented on the backend, see `docs/IMPLEMENTED.md`; for backend future days, see `docs/roadmap.md`; for client execution, see `client/docs/roadmap.md` (this file). Last updated: 2026-08-26 at `main@41e794b` (and `rummy-mvp-rc1` `cfce62e` + `client` scaffolding `5ade046`).*
