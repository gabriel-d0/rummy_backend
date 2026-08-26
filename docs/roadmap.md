# Roadmap — Days to be Implemented

**Scope:** This file lists the *future* days that are **planned but not yet implemented** after the current MVP. For what is already done, see `docs/IMPLEMENTED.md` (Phases 1–16 complete, `rummy-mvp-rc1` at `cfce62e` / `41e794b`, 82 commits `36c2c59..41e794b`). This roadmap follows the Handmade Hero incremental style: one small vertical slice per day, `make check` + `docker compose build` green before next day.

**Current HEAD:** `main@41e794b` — all 24-day `AGENTS.md` days 1–24 done (foundation through minimal CLI `5ade046` + deterministic simulation `01f6a3c` + hardening `3e1b92a` + win `2d278a5`). Tagged `rummy-mvp-rc1`.

**How to read:** Day numbers are per the extended roadmap (≈Day 25–820) that follows the 24-day MVP. Each day is one focused commit; phases end in a runnable, tested state. This file now shows **per-day** breakdown (not ranges).

---

## Phase 17 — Future product features (still deferred per `AGENTS.md:100`)

These were intentionally deferred until the baseline round is playable and tested. They are the *first* days to be implemented post-MVP if product approves.

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 147 | Multi-round scoring model | Design scoring model for multi-round: per-round raw points vs cumulative, data structures `MatchScore`/`RoundScore`. | ⏳ Planned |
| 148 | Dealer rotation | Implement deterministic dealer rotation per round (seat order, opening player shifts anticlockwise). | ⏳ Planned |
| 149 | Accumulated scores | Implement cumulative `totalScore` per `PlayerId` across rounds, with reset per match. | ⏳ Planned |
| 150 | Bonus calculation | Implement bonuses for quick win / low deadwood / joker conservation (per product decision). | ⏳ Planned |
| 151 | Match winner determination | Implement `matchWinner` logic: first to reach target or highest after N rounds, tie-break. | ⏳ Planned |
| 152 | Round history storage | Persist per-round history: `roundNumber`, `winner`, `TableMelds`, `Racks` counts, `Stock` count. | ⏳ Planned |
| 153 | History API | Add `GetMatchHistory` RPC / `OpServerState` extension to return round history to clients. | ⏳ Planned |
| 154 | Multi-round tests | Test 2/3/4-player multi-round flows with `CheckTileConservation` after each round. | ⏳ Planned |
| 155 | Multi-round polish | Integration, docs update `docs/rules-decisions.md` § multi-round, `README` next steps. | ⏳ Planned |
| 156 | Social: table chat design | Define chat scope: table-only, rate limits, ephemeral vs persistent. | ⏳ Planned |
| 157 | Table chat implementation | Implement `OpClientChat` + `OpServerChat` with ordering, sender `PlayerId`, length limits. | ⏳ Planned |
| 158 | Chat reconnection | Ensure chat history re-sent on rejoin via `PrivateView` or separate log. | ⏳ Planned |
| 159 | Player profiles design | Define `PlayerProfile` fields: display name, avatar, language, region. | ⏳ Planned |
| 160 | Player profiles implementation | Implement `GetProfile`/`UpdateProfile` with validation and reserved names. | ⏳ Planned |
| 161 | Profile moderation | Add report flow and admin moderation for names/avatars. | ⏳ Planned |
| 162 | Friends design | Define friend request, accepted, blocked states. | ⏳ Planned |
| 163 | Friends implementation | Implement `AddFriend`/`AcceptFriend`/`Block` with privacy controls. | ⏳ Planned |
| 164 | Private invitations | Implement private match invitations via `PlayerId`. | ⏳ Planned |
| 165 | Social tests | Test chat rate limits, friend blocking, profile visibility. | ⏳ Planned |
| 166 | Matchmaking design | Define public tables, private rooms, skill filters, reconnect UX. | ⏳ Planned |
| 167 | Public tables | Implement lobby listing for public tables (`MatchLabelUpdate` + filter). | ⏳ Planned |
| 168 | Private rooms | Implement private `create` with `password`/`invite` and `join` via `matchId`. | ⏳ Planned |
| 169 | Skill filters | Add optional `minScore`/`maxScore` filter for matchmaking. | ⏳ Planned |
| 170 | Reconnect UX | Design UX for disconnected player: keep `Seat`/`Racks`, show `reconnecting` state. | ⏳ Planned |
| 171 | Reconnect handling | Implement grace period before ending/pausing match on disconnect. | ⏳ Planned |
| 172 | Rematch | Add `rematch` request after `RoundComplete`. | ⏳ Planned |
| 173 | Matchmaking tests | Test public/private, skill filters, reconnect, rematch. | ⏳ Planned |
| 174 | Matchmaking docs | Document lobby, private rooms, skill filters, reconnect in `docs/state-machine.md`. | ⏳ Planned |
| 175 | Matchmaking polish | Polish `README` `How to Debug a Match` with lobby examples. | ⏳ Planned |
| 176 | Tournaments: design | Define scheduled tournaments, formats: single elimination, round-robin, Swiss. | ⏳ Planned |
| 177 | Tournament creation | Implement `CreateTournament` with `name`, `format`, `startTime`, `maxPlayers`. | ⏳ Planned |
| 178 | Tournament registration | Implement `JoinTournament` and `LeaveTournament` before start. | ⏳ Planned |
| 179 | Tournament brackets | Generate brackets/standings from participants, seed by rating. | ⏳ Planned |
| 180 | Tournament matches | Auto-create rummy matches per bracket round, assign seats. | ⏳ Planned |
| 181 | Tournament scoring | Update tournament standings after each match `RoundComplete`. | ⏳ Planned |
| 182 | Tournament advancement | Implement advancement: winner moves to next bracket, loser eliminated or Swiss points. | ⏳ Planned |
| 183 | Tournament standings | Add `GetTournamentStandings` RPC. | ⏳ Planned |
| 184 | Tournament prizes | Define prize logic (only if legally appropriate; else cosmetic). | ⏳ Planned |
| 185 | Tournament prizes implementation | Implement cosmetic `Inventory` grants for winners (no real-money unless licensed). | ⏳ Planned |
| 186 | Tournament moderation | Add admin `CancelTournament` and `Disqualify` with audit. | ⏳ Planned |
| 187 | Tournament notifications | Notify participants of start, next round, completion via `OpServerEvent`. | ⏳ Planned |
| 188 | Tournament tests | Test bracket generation, advancement, standings, prizes. | ⏳ Planned |
| 189 | Tournament load test | Test with 16/32 players, parallel matches, no `TileId` leak. | ⏳ Planned |
| 190 | Tournament docs | Document tournament formats, registration, standings in `docs/rules-decisions.md`. | ⏳ Planned |
| 191 | Bots: deterministic test bot | Create test bot that plays random legal moves for load testing. | ⏳ Planned |
| 192 | Test bot integration | Integrate test bot into `MatchLoop` via same `OpClient*` path. | ⏳ Planned |
| 193 | Test bot vs test bot simulation | Run `deterministic_simulation_test.go` with test bots instead of fixed moves. | ⏳ Planned |
| 194 | AI opponent: design | Define AI opponent requirements: same protocol, no hidden info. | ⏳ Planned |
| 195 | AI opponent: one-turn lookahead | Implement simple AI: best immediate meld/discard via scoring. | ⏳ Planned |
| 196 | AI opponent: two-turn planning | Add shallow search for future stock/discards. | ⏳ Planned |
| 197 | AI opponent: difficulty levels | `beginner`/`normal`/`advanced` with different search depth and randomness. | ⏳ Planned |
| 198 | AI opponent: timing | Add `think` delay, deterministic test mode, cancellation on match end. | ⏳ Planned |
| 199 | AI opponent: tests | Test AI never leaks hidden info, always legal, fallback on failure. | ⏳ Planned |
| 200 | AI opponent: observability | Log `botId`, `seat`, `action`, `duration`, `fallback`. | ⏳ Planned |
| 201 | AI opponent: self-play runner | Script to run thousands of bot-vs-bot rounds with seeds. | ⏳ Planned |
| 202 | AI opponent: telemetry | Collect `winRate`, `invalidActionCount`, `timeoutRate`. | ⏳ Planned |
| 203 | AI opponent: benchmarks | Compare difficulties on fixed board scenarios. | ⏳ Planned |
| 204 | AI opponent: calibration | Adjust difficulty via self-play and human practice results. | ⏳ Planned |
| 205 | AI opponent: docs | Document AI opponent design, difficulty, and `TODO(product)` for ranking separation. | ⏳ Planned |
| 206 | Production ops: logging standard | Define structured fields: `requestId`, `matchId`, `playerId` hash, `op`, `result`, `duration`. | ⏳ Planned |
| 207 | Production ops: metrics foundation | Add `requests`, `matchCount`, `activeUsers`, `commandVolume`, `errors`, `latency`, `memory`, `CPU`, `DB` metrics. | ⏳ Planned |
| 208 | Production ops: match metrics | Track `creation`, `joins`, `starts`, `completion`, `abandonment`, `duration`, `playerCount`. | ⏳ Planned |
| 209 | Production ops: gameplay metrics | Track `draw` type, `pickup` use, `opening` timing, `joker` use, `win` conditions, `validation` failures. | ⏳ Planned |
| 210 | Production ops: bot metrics | Track `botDifficulty`, `decisionLatency`, `searchDepth`, `fallbackRate`, `winRate`. | ⏳ Planned |
| 211 | Production ops: dashboards | Build dashboards for infra health, live matches, gameplay funnel, bot behavior. | ⏳ Planned |
| 212 | Production ops: alerting | Alerts for DB failures, Nakama restart loops, elevated errors, stuck matches, slow bots. | ⏳ Planned |
| 213 | Production ops: tracing | Add request/match tracing for external services and AI workers. | ⏳ Planned |
| 214 | Production ops: analytics validation | Verify event correctness with controlled test matches, ensure no hidden rack details. | ⏳ Planned |
| 215 | Production ops: load testing | Load test with 100+ concurrent matches, measure `TickRate` 5, `CheckTileConservation`. | ⏳ Planned |
| 216 | Production ops: backup drills | Automate PostgreSQL backups and test restoration. | ⏳ Planned |
| 217 | Production ops: staging | Deploy production-like staging and run smoke tests. | ⏳ Planned |
| 218 | Production ops: release process | Define approvals, release notes, rollback, incident contacts. | ⏳ Planned |
| 219 | Production ops: cost monitoring | Monitor Nakama + DB + AI worker costs, set limits. | ⏳ Planned |
| 220 | Production ops: docs | Document ops, dashboards, alerts, runbooks in `docs/architecture.md`. | ⏳ Planned |
| 221 | Polish: typography | Establish typography scale for table, rack, meld displays. | ⏳ Planned |
| 222 | Polish: colours | Define colour palette with `red`/`yellow`/`blue`/`black` tile colours accessible. | ⏳ Planned |
| 223 | Polish: spacing | Define spacing scale for `TableMelds`/`DiscardRow`/`Rack`. | ⏳ Planned |
| 224 | Polish: buttons | Design `Draw`/`Discard`/`Meld`/`Extend`/`Pickup`/`Replace` buttons. | ⏳ Planned |
| 225 | Polish: modals | Design `Meld` and `Replace` modals with tile selection. | ⏳ Planned |
| 226 | Polish: table layout | Design table layout for `PublicView` (`StockCount`, `DiscardRow`, `TableMelds`). | ⏳ Planned |
| 227 | Polish: rack sorting | Implement rack sorting by `Colour` then `Rank`. | ⏳ Planned |
| 228 | Polish: drag-drop | Add drag/drop for tiles in future web client. | ⏳ Planned |
| 229 | Polish: tile skins | Add tile skins as cosmetic-only inventory. | ⏳ Planned |
| 230 | Polish: table themes | Add table themes as cosmetic-only. | ⏳ Planned |
| 231 | Polish: animations | Add subtle animations for `Draw`/`Discard`/`Meld` (respect `prefers-reduced-motion`). | ⏳ Planned |
| 232 | Polish: sound | Add optional sound for `meld` and `win` (muted by default). | ⏳ Planned |
| 233 | Polish: accessibility baseline | Add keyboard control, screen-reader labels, focus handling. | ⏳ Planned |
| 234 | Polish: colour-blind safe | Ensure tile colours have text labels + patterns. | ⏳ Planned |
| 235 | Polish: scalable text | Add scalable text for `RackCount`/`StockCount`. | ⏳ Planned |
| 236 | Polish: reduced motion | Respect `prefers-reduced-motion` for animations. | ⏳ Planned |
| 237 | Polish: responsive | Support desktop/tablet/mobile layouts without reducing rule clarity. | ⏳ Planned |
| 238 | Polish: localization foundation | Externalize UI text, support Romanian and English first. | ⏳ Planned |
| 239 | Polish: pluralization | Handle pluralization for `RackCount`/`StockCount`/`Meld` counts. | ⏳ Planned |
| 240 | Polish: docs | Document polish decisions in `docs/architecture.md` and `README` `Next Steps`. | ⏳ Planned |

---

## Phase 18 — AI Bots: Foundation and Safe Integration

Bots must use the same authoritative command protocol as humans; never give bots hidden state.

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 241 | Bot architecture decision | Document bot lifecycle, identities, difficulty levels, scheduling, how bot actions enter Nakama matches. | ⏳ Planned |
| 242 | Bot lifecycle design | Design `BotId`/`Seat` assignment, think-delay scheduling, cancellation on `RoundComplete`. | ⏳ Planned |
| 243 | Bot scheduling model | Implement scheduler that queues bot actions per `CurrentSeat` and `TurnPhase`. | ⏳ Planned |
| 244 | Bot player model | Add `BotUser`/`BotProfile` metadata, display names, avatar placeholders. | ⏳ Planned |
| 245 | Bot flags | Add `IsBot` flag to `PlayerState` and `PublicPlayer{IsBot}`. | ⏳ Planned |
| 246 | Bot seat assignment | Implement safe `Seat` assignment for bots, never overwrite human `Seat`. | ⏳ Planned |
| 247 | Bot match participation: host adds | Allow host to add 1–3 bots to private/local match via `OpClientAddBot`. | ⏳ Planned |
| 248 | Bot match participation: validation | Validate bot count `2–4` total, not in `RoundComplete`, only host. | ⏳ Planned |
| 249 | Bot match participation: start | Ensure bot seats are filled before `OpClientStart` and `NewRoundState` includes bot racks. | ⏳ Planned |
| 250 | Bot match participation: tests | Test 2/3/4-player with bots, seat assignment, `CheckTileConservation`. | ⏳ Planned |
| 251 | Bot action adapter: design | Design adapter `BotDecision→OpClient*` with `requestId` and payload validation. | ⏳ Planned |
| 252 | Bot action adapter: draw | Implement `DRAW_STOCK`/`DRAW_PREVIOUS`/`PICKUP` via adapter. | ⏳ Planned |
| 253 | Bot action adapter: discard/meld | Implement `DISCARD`/`MELD_INITIAL`/`MELD_NEW`/`EXTEND`/`REPLACE` via adapter. | ⏳ Planned |
| 254 | Bot action adapter: tests | Test adapter always produces `ValidatePayload` `ok` and `ValidateRun/Set` `ok`. | ⏳ Planned |
| 255 | Bot timing: think delay | Add configurable `thinkMs` per difficulty (e.g., `beginner 500ms`, `advanced 1500ms`). | ⏳ Planned |
| 256 | Bot timing: deterministic test mode | Add `testModeNoDelay` where bots act immediately for `go test`. | ⏳ Planned |
| 257 | Bot timing: cancellation | Cancel bot action if match ends or bot leaves before delay. | ⏳ Planned |
| 258 | Bot timing: tests | Test think delay, cancellation, deterministic mode. | ⏳ Planned |
| 259 | Bot visibility boundaries | Ensure bot only receives `PrivateView` for its `Seat` + `PublicView`; add test `TestBotNoOpponentRack`. | ⏳ Planned |
| 260 | Bot visibility: joker | Ensure bot never sees `JokerReps` beyond public `TableMelds`. | ⏳ Planned |
| 261 | Bot visibility: stock | Ensure bot never sees `Stock` order, only `StockCount`. | ⏳ Planned |
| 262 | Bot visibility: tests | Tests proving no `OwnRack` leak to bot or other. | ⏳ Planned |
| 263 | Bot error recovery | If bot generates invalid action, log `invalid bot action` and fallback to `DISCARD` lowest deadwood. | ⏳ Planned |
| 264 | Bot error recovery: log | Log `matchId`, `botId`, `attemptedOp`, `errorCode`, `fallbackOp`. | ⏳ Planned |
| 265 | Bot error recovery: fallback | Select fallback legal move without corrupting match (`CheckTileConservation` still `106`). | ⏳ Planned |
| 266 | Bot error recovery: tests | Test invalid bot action → fallback and `TableMelds` unchanged for invalid. | ⏳ Planned |
| 267 | Bot observability | Add structured logs: `matchId`, `botId`, `turn`, `action`, `duration`, `fallback`, `rejection reason`. | ⏳ Planned |
| 268 | Bot observability: metrics | Add `botDecisionCount`, `botFallbackCount`, `botInvalidCount` metrics. | ⏳ Planned |
| 269 | Bot observability: tests | Test logs contain `botId` and `action` and no private `TileId`. | ⏳ Planned |
| 270 | Bot observability: docs | Document bot observability in `docs/architecture.md`. | ⏳ Planned |

---

## Phase 19 — AI Bots: Legal Move Generation and Beginner Difficulty

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 271 | Legal move generator | Create pure module `bot/legalMoves` enumerating legal actions from `PrivateView` + public. | ⏳ Planned |
| 272 | Legal move: draw | Determine available draws: stock, previous discard, earlier pickup. | ⏳ Planned |
| 273 | Legal move: discard | Generate every legal `DISCARD` candidate for bot’s `OwnRack` after draw/meld. | ⏳ Planned |
| 274 | Legal move: meld validation | Ensure each generated move passes `ValidateRun/Set` and `ValidateBatchOwnership`. | ⏳ Planned |
| 275 | Draw decision generation | Determine available draw choices with `CanPickupPreviousDiscard`/`CanPickupDiscardForMeld` checks. | ⏳ Planned |
| 276 | Draw decision: stock vs previous | Score stock vs previous discard (e.g., prefer previous if it completes a meld). | ⏳ Planned |
| 277 | Draw decision: earlier pickup | Score earlier pickup: need `discardIndex` + `tileIds[2]` to form valid meld + `laterTiles` sweep value. | ⏳ Planned |
| 278 | Draw decision: tests | Test draw choices in fixed board states. | ⏳ Planned |
| 279 | Discard candidate generation | Generate every legal discard candidate for bot’s rack after draw/meld. | ⏳ Planned |
| 280 | Discard: heuristic | Prefer discarding low-value deadwood (`2–9:5` vs `10–13:10` vs `Ace 5/10` vs `Joker rep`). | ⏳ Planned |
| 281 | Discard: avoid needed | Avoid discarding tiles needed for existing sequences/sets or near-complete melds. | ⏳ Planned |
| 282 | Discard: tests | Test discard choices never include `IsJoker` unless forced and never reveal hidden info. | ⏳ Planned |
| 283 | Initial meld candidate generation | Generate possible opening batches satisfying 50-point + at-least-one-run. | ⏳ Planned |
| 284 | Initial meld: scoring | Score batches via `TotalScore` and `ValidateBatchScore`/`ValidateBatchHasRun`. | ⏳ Planned |
| 285 | Initial meld: tests | Test opening batches: 49 rejected, 50 accepted, no-run rejected. | ⏳ Planned |
| 286 | New meld candidate generation | Generate valid sets/runs from opened rack. | ⏳ Planned |
| 287 | New meld: batch | Permit multiple new melds per batch, atomic `ValidateBatchOwnership`. | ⏳ Planned |
| 288 | New meld: tests | Test `MELD_NEW` candidates with `CheckTileConservation`. | ⏳ Planned |
| 289 | Meld extension candidate generation | Generate legal extensions to public melds, including others’ `OwnerSeat`. | ⏳ Planned |
| 290 | Meld extension: validation | Validate entire resulting meld via `meld.New` + `ValidateRun/Set`. | ⏳ Planned |
| 291 | Meld extension: tests | Test `EXTEND_MELD` candidates with `TableMeld.Kind` stable. | ⏳ Planned |
| 292 | Joker replacement candidate generation | Generate legal joker replacements (exact tile for run, exact missing colour for set). | ⏳ Planned |
| 293 | Joker replacement: validation | Validate `updatedMeld` and `newMeld` (`joker+2` tiles) via `meld.New`. | ⏳ Planned |
| 294 | Joker replacement: tests | Test `REPLACE_JOKER` candidates with `JokerReps` immutability. | ⏳ Planned |
| 295 | Beginner strategy: design | Design heuristic: open when possible, obvious melds, extend, discard low-value. | ⏳ Planned |
| 296 | Beginner strategy: implement | Implement `BeginnerStrategy` that picks first legal opening or best immediate meld. | ⏳ Planned |
| 297 | Beginner strategy: tests | Test beginner always produces legal `OpClient*` and never stalls (`TurnPhase` advances). | ⏳ Planned |
| 298 | Beginner strategy: polish | Ensure beginner never leaks hidden info and logs `fallback` if needed. | ⏳ Planned |
| 299 | Beginner bot tests | Test that bot always produces legal commands across fixed scenarios. | ⏳ Planned |
| 300 | Beginner bot: exhaustive | Test bot never reveals hidden info via `json.Marshal(PublicView)` string search. | ⏳ Planned |
| 301 | Bot turn simulation: single bot | Run deterministic game with 1 beginner bot + 1 human (alice). | ⏳ Planned |
| 302 | Bot turn simulation: multiple bots | Run deterministic games with 2–3 beginner bots. | ⏳ Planned |
| 303 | Bot turn simulation: verification | Verify rounds eventually progress or terminate safely, `CheckTileConservation` holds. | ⏳ Planned |
| 304 | Bot turn simulation: docs | Document `TestBotTurnSimulation` in `docs/testing.md`. | ⏳ Planned |
| 305 | Bot turn simulation: performance | Ensure bot turn < `100ms` in test mode. | ⏳ Planned |
| 306 | Beginner bot tests: exhaustive | Test bot always legal across `n=2,3,4 × seeds`. | ⏳ Planned |
| 307 | Beginner bot: docs | Document beginner difficulty in `docs/rules-decisions.md` and `README`. | ⏳ Planned |
| 308 | Beginner milestone polish | Polish `README` `Next Steps` to reflect beginner bot done. | ⏳ Planned |
| 309 | Beginner milestone: tag | Tag `bot-beginner-rc1` when all 271–308 pass. | ⏳ Planned |
| 310 | Bot turn simulation: final | Run full deterministic simulation with bots as in `deterministic_simulation_test.go` but with bots. | ⏳ Planned |

---

## Phase 20 — AI Bots: Intermediate Strategy and Difficulty Tiers

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 311 | Evaluation model: design | Define hand-evaluation: meld potential, deadwood, joker value, discard risk. | ⏳ Planned |
| 312 | Evaluation model: meld potential | Score `run` vs `set` potential per `Colour`/`Rank`. | ⏳ Planned |
| 313 | Evaluation model: deadwood | Score `deadwood` value via `ScoreTile` (5/10/25) minus meld potential. | ⏳ Planned |
| 314 | Evaluation model: joker value | Score `Joker` value via `represented` tile plus flexibility. | ⏳ Planned |
| 315 | Candidate scoring: design | Design `ScoreCandidate` that scores each legal `OpClient*` via evaluation model. | ⏳ Planned |
| 316 | Candidate scoring: implementation | Implement scoring for `MELD_INITIAL`/`MELD_NEW`/`EXTEND`/`REPLACE`/`DISCARD` etc. | ⏳ Planned |
| 317 | Candidate scoring: tests | Test scoring picks best immediate meld vs discard. | ⏳ Planned |
| 318 | Candidate scoring: docs | Document scoring model in `docs/architecture.md`. | ⏳ Planned |
| 319 | Better discard heuristic | Avoid discarding tiles likely needed for existing sequences/sets or near-complete melds. | ⏳ Planned |
| 320 | Discard heuristic: tests | Test discard heuristic never discards needed tile if alternative. | ⏳ Planned |
| 321 | Discard heuristic: docs | Document heuristic in `docs/rules-decisions.md`. | ⏳ Planned |
| 322 | Opening strategy | Decide when to preserve tiles vs open immediately with 50-point batch. | ⏳ Planned |
| 323 | Opening strategy: tests | Test opening strategy: 49 vs 50, no-run, duplicate. | ⏳ Planned |
| 324 | Opening strategy: docs | Document opening threshold per `ScoreTile` context. | ⏳ Planned |
| 325 | Pickup strategy | Score stock vs latest vs earlier pickup (sweep value vs meld value). | ⏳ Planned |
| 326 | Pickup strategy: tests | Test pickup strategy picks earlier when it completes a meld and sweep is low deadwood. | ⏳ Planned |
| 327 | Pickup strategy: docs | Document pickup scoring in `docs/state-machine.md`. | ⏳ Planned |
| 328 | Table exploitation | Prefer safe extensions that reduce `RackCount` or improve future options. | ⏳ Planned |
| 329 | Table exploitation: tests | Test extension candidate scoring. | ⏳ Planned |
| 330 | Table exploitation: docs | Document exploitation in `docs/architecture.md`. | ⏳ Planned |
| 331 | Joker strategy | Evaluate retain/meld/replace/expose jokers. | ⏳ Planned |
| 332 | Joker strategy: tests | Test joker strategy: never discard joker unless forced. | ⏳ Planned |
| 333 | Joker strategy: docs | Document joker strategy per `docs/rules-decisions.md:1.3`. | ⏳ Planned |
| 334 | Difficulty configuration | Add `beginner`/`normal`/`advanced` profiles with search depth/randomness/delay. | ⏳ Planned |
| 335 | Difficulty: implementation | Implement profile config via `BotConfig{ThinkMs, SearchDepth, Randomness}`. | ⏳ Planned |
| 336 | Difficulty: tests | Test `beginner` vs `advanced` pick different moves on same board. | ⏳ Planned |
| 337 | Controlled randomness | Add seeded `Rand` weighted choice, reproducible. | ⏳ Planned |
| 338 | Controlled randomness: tests | Test seeded choice deterministic. | ⏳ Planned |
| 339 | Controlled randomness: docs | Document seeded choice in `docs/testing.md`. | ⏳ Planned |
| 340 | Strategy benchmark suite | Create fixed board-state scenarios, compare `beginner`/`normal`/`advanced` choices. | ⏳ Planned |
| 341 | Benchmark: implementation | Implement `bot/benchmark_test.go` with `n=10` fixed scenarios. | ⏳ Planned |
| 342 | Benchmark: docs | Document benchmark suite in `docs/testing.md`. | ⏳ Planned |
| 343 | Benchmark: polish | Ensure `advanced` measurably better without cheating (no hidden info). | ⏳ Planned |
| 344 | Controlled randomness: polish | Ensure not mechanically identical but reproducible. | ⏳ Planned |
| 345 | Benchmark: final | Run benchmark and record `winRate` per difficulty. | ⏳ Planned |
| 346 | Intermediate milestone polish | Polish `README` `Next Steps` to reflect intermediate bot done. | ⏳ Planned |
| 347 | Intermediate milestone: docs | Update `docs/rules-decisions.md` with intermediate strategy decisions. | ⏳ Planned |
| 348 | Intermediate milestone: tag | Tag `bot-intermediate-rc1` when all 311–348 pass. | ⏳ Planned |
| 349 | Intermediate final tests | Run full `go test ./...` with intermediate bots. | ⏳ Planned |
| 350 | Documentation polish | Polish `docs/architecture.md` with evaluation model. | ⏳ Planned |

---

## Phase 21 — AI Bots: Search, Simulation, and Training Infrastructure

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 351 | Search budget design | Define max CPU time, candidate count, simulation count, cancellation. | ⏳ Planned |
| 352 | Search budget: implementation | Implement `BotSearchBudget{MaxMs, MaxCandidates, MaxSimulations}`. | ⏳ Planned |
| 353 | Search budget: tests | Test budget enforcement (timeout, max candidates). | ⏳ Planned |
| 354 | Search budget: docs | Document budget in `docs/architecture.md`. | ⏳ Planned |
| 355 | One-turn lookahead | Implement bounded evaluation of immediate actions. | ⏳ Planned |
| 356 | One-turn: tests | Test lookahead picks better immediate meld vs discard. | ⏳ Planned |
| 357 | One-turn: docs | Document lookahead in `docs/testing.md`. | ⏳ Planned |
| 358 | Two-turn approximate planning | Add limited future-turn estimation without hidden info. | ⏳ Planned |
| 359 | Two-turn: tests | Test two-turn planning improves over one-turn in fixed scenario. | ⏳ Planned |
| 360 | Two-turn: docs | Document two-turn planning. | ⏳ Planned |
| 361 | Monte Carlo rollout prototype | Simulate plausible unseen `Stock`/`Racks` distributions from public info only. | ⏳ Planned |
| 362 | Monte Carlo: implementation | Implement `MonteCarloRollout` that shuffles `unknown` tiles from public counts. | ⏳ Planned |
| 363 | Monte Carlo: tests | Test rollout never uses hidden `TileId` not in public `StockCount`/`RackCount`. | ⏳ Planned |
| 364 | Monte Carlo: docs | Document Monte Carlo and `TileId` privacy. | ⏳ Planned |
| 365 | Simulation state cloning | Build efficient deterministic copies for simulations (`RoundState.Clone`). | ⏳ Planned |
| 366 | Simulation cloning: tests | Test `Clone` deep copies `Racks`/`Stock`/`DiscardRow`/`TableMelds` and `CheckTileConservation`. | ⏳ Planned |
| 367 | Simulation cloning: docs | Document `Clone` in `docs/architecture.md`. | ⏳ Planned |
| 368 | Search cancellation | Cancel search on `RoundComplete`/`CurrentSeat` change/`MatchLeave`/`MatchTerminate`/timeout. | ⏳ Planned |
| 369 | Search cancellation: tests | Test cancellation on match end and timeout. | ⏳ Planned |
| 370 | Search cancellation: docs | Document cancellation in `docs/state-machine.md`. | ⏳ Planned |
| 371 | Fallback policy | Guarantee legal fallback if simulation fails/times out/`foundValid==false`. | ⏳ Planned |
| 372 | Fallback: tests | Test fallback always legal and `CheckTileConservation`. | ⏳ Planned |
| 373 | Fallback: docs | Document fallback in `docs/rules-decisions.md`. | ⏳ Planned |
| 374 | Offline self-play runner | Script `scripts/bot-self-play.sh` to run thousands of bot-vs-bot rounds with seeds. | ⏳ Planned |
| 375 | Self-play: metrics | Collect `winRate`, `turnLength`, `invalidActionCount`, `timeoutRate`, `duration`. | ⏳ Planned |
| 376 | Self-play: tests | Test runner with `seed 42` deterministic. | ⏳ Planned |
| 377 | Self-play: docs | Document runner in `docs/testing.md`. | ⏳ Planned |
| 378 | Bot telemetry | Collect `botDecisionCount`, `botFallbackCount`, `botInvalidCount` metrics. | ⏳ Planned |
| 379 | Telemetry: tests | Test metrics increment. | ⏳ Planned |
| 380 | Telemetry: docs | Document telemetry in `docs/architecture.md`. | ⏳ Planned |
| 381 | Performance tests | Verify bot turns meet latency/CPU budgets under concurrent load. | ⏳ Planned |
| 382 | Performance: test | Test with 4 concurrent matches, 3 bots each, `go test -run TestBotPerformance`. | ⏳ Planned |
| 383 | Performance: docs | Document budgets in `docs/testing.md`. | ⏳ Planned |
| 384 | Performance: polish | Ensure bot never blocks `MatchLoop` tick 5. | ⏳ Planned |
| 385 | Docs polish | Polish `docs/architecture.md` with search infrastructure. | ⏳ Planned |
| 386 | Final docs | Update `README` `Next Steps` to reflect search done. | ⏳ Planned |
| 387 | Tag | Tag `bot-search-rc1` when all 351–386 pass. | ⏳ Planned |
| 388 | Final tests | Run `go test ./...` with search bots. | ⏳ Planned |
| 389 | Polish | Polish `docs/roadmap.md` for search. | ⏳ Planned |
| 390 | Release | Prepare release notes for search. | ⏳ Planned |

---

## Phase 22 — Matchmaking, Ratings, and Bot-Assisted Player Experience

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 391 | Match mode design | Define `casual`, `private`, `practice`, `ranked`, `bot-only` modes. | ⏳ Planned |
| 392 | Match mode: documentation | Document modes in `docs/rules-decisions.md` and `docs/state-machine.md`. | ⏳ Planned |
| 393 | Match mode: tests | Test mode selection via `MatchSignal` `mode:casual` etc. | ⏳ Planned |
| 394 | Match mode: polish | Polish `README` `Next Steps`. | ⏳ Planned |
| 395 | Practice mode | Allow player to start private practice match with 1–3 bots via `OpClientAddBot`. | ⏳ Planned |
| 396 | Practice mode: validation | Validate `HasOpened` not required for bot practice, only host. | ⏳ Planned |
| 397 | Practice mode: tests | Test `2+1` bot practice with `CheckTileConservation`. | ⏳ Planned |
| 398 | Practice mode: docs | Document practice mode in `README` `Minimal Test Client`. | ⏳ Planned |
| 399 | Bot transparency UX/protocol | Mark bot players clearly: `PublicPlayer{IsBot}` + `PlayerState{IsBot}` + `OpServerState` `isBot`. | ⏳ Planned |
| 400 | Bot transparency: tests | Test `PublicView` contains `IsBot` and `PrivateView` per bot not leak. | ⏳ Planned |
| 401 | Bot transparency: docs | Document `IsBot` in `docs/protocol.md` and `docs/terminology.md`. | ⏳ Planned |
| 402 | Bot transparency: polish | Ensure UI shows bot avatar distinct. | ⏳ Planned |
| 403 | Matchmaking queues | Add queue metadata: `playerCount`, `region`, `mode`, `difficulty`. | ⏳ Planned |
| 404 | Matchmaking queues: implementation | Implement `MatchmakingQueue` struct and `Enqueue`/`Dequeue`. | ⏳ Planned |
| 405 | Matchmaking queues: tests | Test queue with `region` and `difficulty`. | ⏳ Planned |
| 406 | Matchmaking queues: docs | Document queues in `docs/architecture.md`. | ⏳ Planned |
| 407 | Queue timeout policy | Define when matchmaking offers bots after waiting threshold; require opt-in where appropriate. | ⏳ Planned |
| 408 | Queue timeout: implementation | Implement `QueueTimeout` `5s` → offer `bot` via `OpServerEvent` `offerBot`. | ⏳ Planned |
| 409 | Queue timeout: tests | Test timeout policy with `offerBot` and `decline`. | ⏳ Planned |
| 410 | Queue timeout: docs | Document timeout in `docs/state-machine.md`. | ⏳ Planned |
| 411 | Backfill bots | Support filling empty seats with bots only before `PhasePlaying` or per explicit match rules. | ⏳ Planned |
| 412 | Backfill: tests | Test backfill before `Start` vs after `Start` rejected. | ⏳ Planned |
| 413 | Backfill: docs | Document backfill in `docs/rules-decisions.md`. | ⏳ Planned |
| 414 | Backfill: polish | Polish `README` `Next Steps`. | ⏳ Planned |
| 415 | Player rating model | Add `Rating` `Elo`/`Glicko` for human-vs-human ranked games; keep bot games separate. | ⏳ Planned |
| 416 | Rating model: implementation | Implement `Rating{PlayerId, Value, Deviation}` and `UpdateRating` after `RoundComplete`. | ⏳ Planned |
| 417 | Rating model: tests | Test `Rating` updates with `CheckTileConservation` still `106`. | ⏳ Planned |
| 418 | Rating model: docs | Document rating in `docs/architecture.md`. | ⏳ Planned |
| 419 | Rating protection | Ensure bot matches do not inflate/deflate ranked ratings unless explicitly designed. | ⏳ Planned |
| 420 | Rating protection: tests | Test `bot` match `Rating` unchanged. | ⏳ Planned |
| 421 | Rating protection: docs | Document `bot` vs `ranked` separation in `docs/rules-decisions.md`. | ⏳ Planned |
| 422 | Bot calibration | Adjust bot difficulty using self-play and human practice results. | ⏳ Planned |
| 423 | Bot calibration: implementation | Implement `CalibrateBot` that adjusts `ThinkMs`/`SearchDepth` via `winRate`. | ⏳ Planned |
| 424 | Bot calibration: tests | Test calibration improves `beginner` `winRate` vs `advanced`. | ⏳ Planned |
| 425 | Bot calibration: docs | Document calibration in `docs/testing.md`. | ⏳ Planned |
| 426 | Matchmaking tests | Test queue formation, bot insertion, rating separation, transparency. | ⏳ Planned |
| 427 | Matchmaking: exhaustive | Test `n=2,3,4 × botCount` with `CheckTileConservation`. | ⏳ Planned |
| 428 | Matchmaking: docs | Update `README` `Next Steps` to reflect matchmaking done. | ⏳ Planned |
| 429 | Matchmaking: polish | Polish `docs/architecture.md` with matchmaking. | ⏳ Planned |
| 430 | Matchmaking: tag | Tag `matchmaking-rc1` when all 391–429 pass. | ⏳ Planned |

---

## Phase 23 — Production AI Operations, Safety, and Continuous Improvement

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 431 | AI feature flags | Add server flags for bot availability, strategy version, delays, search depth. | ⏳ Planned |
| 432 | AI feature flags: implementation | Implement `BotFeatureFlags{Enabled, StrategyVersion, ThinkMs, SearchDepth}`. | ⏳ Planned |
| 433 | AI feature flags: tests | Test flag `Enabled=false` → no bot actions. | ⏳ Planned |
| 434 | AI feature flags: docs | Document flags in `docs/architecture.md`. | ⏳ Planned |
| 435 | Strategy versioning | Version bot policies; live match records `strategyVersion` that generated decisions. | ⏳ Planned |
| 436 | Strategy versioning: implementation | Implement `StrategyVersion` in `BotDecision` logs and `TableMeld` metadata. | ⏳ Planned |
| 437 | Strategy versioning: tests | Test `strategyVersion` logged correctly. | ⏳ Planned |
| 438 | Strategy versioning: docs | Document versioning in `docs/rules-decisions.md`. | ⏳ Planned |
| 439 | Safe rollout | Staged: local → dev → staging → limited prod cohort → full release. | ⏳ Planned |
| 440 | Safe rollout: docs | Document rollout in `docs/architecture.md`. | ⏳ Planned |
| 441 | Safe rollout: tests | Test rollout flag `cohort=10%`. | ⏳ Planned |
| 442 | Safe rollout: polish | Polish `README` `Next Steps`. | ⏳ Planned |
| 443 | Bot health dashboards | Monitor bot action latency, failures, invalid decisions, fallback rate, queue impact. | ⏳ Planned |
| 444 | Bot health: implementation | Implement `BotHealthDashboard` struct and `UpdateBotHealth` after each `BotDecision`. | ⏳ Planned |
| 445 | Bot health: tests | Test dashboard metrics increment. | ⏳ Planned |
| 446 | Bot health: docs | Document dashboards in `docs/architecture.md`. | ⏳ Planned |
| 447 | Alerts and safeguards | Alerts for timeouts, elevated invalid command rate, stuck matches, win-rate shifts. | ⏳ Planned |
| 448 | Alerts: implementation | Implement `BotAlerts` with thresholds `invalidRate>5%` etc. | ⏳ Planned |
| 449 | Alerts: tests | Test alerts triggered on `invalidRate`. | ⏳ Planned |
| 450 | Alerts: docs | Document alerts in `docs/architecture.md`. | ⏳ Planned |
| 451 | Fairness audits | Verify bots use only `PrivateView` + public, cannot inspect hidden `TileId`/`Stock` order. | ⏳ Planned |
| 452 | Fairness audits: implementation | Implement `FairnessAudit` that checks `bot` never accesses `Racks[other]` or `Stock` order. | ⏳ Planned |
| 453 | Fairness audits: tests | Test `TestBotNoOpponentRack` with `json.Marshal(PublicView)` string search. | ⏳ Planned |
| 454 | Fairness audits: docs | Document audits in `docs/testing.md`. | ⏳ Planned |
| 455 | Replay and audit logs | Store privacy-safe action traces for debugging disputed games. | ⏳ Planned |
| 456 | Replay: implementation | Implement `ReplayLog{MatchId, Actions[]}` with `TileId` hashed. | ⏳ Planned |
| 457 | Replay: tests | Test replay log never contains `OwnRack` of other. | ⏳ Planned |
| 458 | Replay: docs | Document replay in `docs/architecture.md`. | ⏳ Planned |
| 459 | Human feedback loop | Post-match bot feedback: difficulty, suspicious behavior, fun rating. | ⏳ Planned |
| 460 | Feedback: implementation | Implement `BotFeedback{MatchId, BotId, Rating, Comment}`. | ⏳ Planned |
| 461 | Feedback: tests | Test feedback stored. | ⏳ Planned |
| 462 | Feedback: docs | Document feedback in `docs/state-machine.md`. | ⏳ Planned |
| 463 | Automated regression corpus | Convert production bugs and strange bot behavior into deterministic regression scenarios. | ⏳ Planned |
| 464 | Regression: implementation | Implement `RegressionCorpus` that stores `seed`, `initialRacks`, `actions`, `expectedWinner`. | ⏳ Planned |
| 465 | Regression: tests | Test corpus replay with `CheckTileConservation`. | ⏳ Planned |
| 466 | Regression: docs | Document corpus in `docs/testing.md`. | ⏳ Planned |
| 467 | Cost and capacity management | Concurrency limits, execution budgets, worker scaling, cost monitoring. | ⏳ Planned |
| 468 | Cost: implementation | Implement `BotCostManager` with `MaxConcurrentBots`, `MaxMsPerTurn`, `WorkerScale`. | ⏳ Planned |
| 469 | Cost: tests | Test cost manager enforces `MaxConcurrentBots`. | ⏳ Planned |
| 470 | Cost: docs | Document cost in `docs/architecture.md`. | ⏳ Planned |
| 471 | Disaster recovery | Behavior if bot service fails: pause/remove/substitute fallback/safely end match. | ⏳ Planned |
| 472 | Disaster: implementation | Implement `BotDisasterRecovery` with `fallback: remove bot, use test bot`. | ⏳ Planned |
| 473 | Disaster: tests | Test disaster recovery with `BotServiceDown`. | ⏳ Planned |
| 474 | Disaster: docs | Document disaster recovery in `docs/state-machine.md`. | ⏳ Planned |
| 475 | AI release candidate | Load tests, fairness tests, replay audits, security review, staged-release checklist. | ⏳ Planned |
| 476 | AI release: load tests | Run load tests with 100 bot matches. | ⏳ Planned |
| 477 | AI release: fairness tests | Run fairness tests with `TestBotNoOpponentRack`. | ⏳ Planned |
| 478 | AI release: security review | Review `Bot` cannot access hidden `TileId` or `Stock` order. | ⏳ Planned |
| 479 | AI release: checklist | Create `docs/ai-release-checklist.md`. | ⏳ Planned |
| 480 | AI release candidate: tag | Tag `ai-rc1` when all 431–479 pass. | ⏳ Planned |

---

## Phase 24 — Production Deployment and Infrastructure

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 481 | Production architecture | Document target hosting, regions, domains, TLS, database topology, backups, scaling, failure boundaries. | ⏳ Planned |
| 482 | Production architecture: implementation | Write `docs/production-architecture.md` with hosting diagram. | ⏳ Planned |
| 483 | Production architecture: tests | Test `docs/production-architecture.md` exists. | ⏳ Planned |
| 484 | Production architecture: polish | Polish `README` `Next Steps`. | ⏳ Planned |
| 485 | Infrastructure as Code | Add Terraform/Pulumi for networking, DB, container runtime, secrets, observability. | ⏳ Planned |
| 486 | IaC: networking | Implement `terraform/network.tf` with `vpc` and `subnets`. | ⏳ Planned |
| 487 | IaC: DB | Implement `terraform/db.tf` with `postgres` `15` `pgdata` volume. | ⏳ Planned |
| 488 | IaC: container runtime | Implement `terraform/nakama.tf` with `nakama` `3.26.0` `backend.so`. | ⏳ Planned |
| 489 | Container hardening | Production Docker images pinned, non-root, minimal layers, health checks, image scanning. | ⏳ Planned |
| 490 | Container hardening: Dockerfile | Update `Dockerfile` with `non-root` `USER` and `HEALTHCHECK`. | ⏳ Planned |
| 491 | Container hardening: scanning | Add `trivy` scan in CI for `backend.so`. | ⏳ Planned |
| 492 | Container hardening: tests | Test `docker compose build` with `non-root` and `HEALTHCHECK`. | ⏳ Planned |
| 493 | Environment management | Separate local/dev/staging/production configs; prevent dev settings reaching prod. | ⏳ Planned |
| 494 | Environment: implementation | Create `nakama/data/production.yml` and `staging.yml` with `DEBUG=false`. | ⏳ Planned |
| 495 | Environment: tests | Test `production.yml` never contains `admin/password`. | ⏳ Planned |
| 496 | Environment: docs | Document environments in `docs/architecture.md`. | ⏳ Planned |
| 497 | Secrets management | Managed secret store for DB credentials, Nakama keys, signing secrets; never commit. | ⏳ Planned |
| 498 | Secrets: implementation | Integrate `aws secretsmanager` or `gcp secretmanager` for `POSTGRES_PASSWORD` and `NAKAMA_KEY`. | ⏳ Planned |
| 499 | Secrets: tests | Test `secrets` not in `git log --all`. | ⏳ Planned |
| 500 | Secrets: docs | Document secrets in `docs/architecture.md`. | ⏳ Planned |
| 501 | TLS and edge security | HTTPS, certificates, secure headers, CORS, rate limits, WAF/CDN. | ⏳ Planned |
| 502 | TLS: implementation | Add `nginx` or `traefik` with `letsencrypt` for `7350`/`7351`. | ⏳ Planned |
| 503 | TLS: tests | Test `curl -k https://` with valid cert. | ⏳ Planned |
| 504 | TLS: docs | Document TLS in `docs/architecture.md`. | ⏳ Planned |
| 505 | Database migrations | Versioned migrations, rollback strategy, migration checks, backup-before-migrate. | ⏳ Planned |
| 506 | Migrations: implementation | Add `migrations/001_init.sql` with `migrate up/down`. | ⏳ Planned |
| 507 | Migrations: tests | Test `migrate up` and `migrate down` with `CheckTileConservation` still `106`. | ⏳ Planned |
| 508 | Migrations: docs | Document migrations in `docs/testing.md`. | ⏳ Planned |
| 509 | Backup and restore drills | Automate PostgreSQL backups and test restoration into isolated environment. | ⏳ Planned |
| 510 | Backup: implementation | Add `scripts/backup.sh` with `pg_dump` and `scripts/restore.sh`. | ⏳ Planned |
| 511 | Backup: tests | Test `backup` and `restore` with `CheckTileConservation`. | ⏳ Planned |
| 512 | Backup: docs | Document backup drills in `docs/architecture.md`. | ⏳ Planned |
| 513 | Staging environment | Deploy production-like staging and run automated smoke tests against it. | ⏳ Planned |
| 514 | Staging: implementation | Deploy `staging` with `docker compose -f compose.staging.yml up --build -d`. | ⏳ Planned |
| 515 | Staging: tests | Test `staging` `make smoke` `SMOKE PASSED`. | ⏳ Planned |
| 516 | Staging: docs | Document staging in `README` `How to Inspect`. | ⏳ Planned |
| 517 | Production release process | Approvals, release notes, rollback process, incident contacts, deployment checklist. | ⏳ Planned |
| 518 | Release process: implementation | Write `docs/release-process.md` with `approvals` and `rollback`. | ⏳ Planned |
| 519 | Release process: tests | Test `release process` exists. | ⏳ Planned |
| 520 | Release process: docs | Document release process in `README` `Next Steps`. | ⏳ Planned |

---

## Phase 25 — Observability, Analytics, and Game Telemetry

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 521 | Logging standard | Structured fields: `requestId`, `matchId`, `playerId` hash, `botId`, `opcode`, `transition result`, `duration`, `error code`. | ⏳ Planned |
| 522 | Logging: implementation | Implement `logger.WithField("matchId", matchId)` in `MatchLoop`. | ⏳ Planned |
| 523 | Logging: tests | Test logs contain `matchId` and `playerId` hash and no private `TileId`. | ⏳ Planned |
| 524 | Logging: docs | Document logging in `docs/architecture.md`. | ⏳ Planned |
| 525 | Metrics foundation | Service metrics: `requests`, `matchCount`, `activeUsers`, `commandVolume`, `errors`, `latency`, `memory`, `CPU`, `DB` health. | ⏳ Planned |
| 526 | Metrics: implementation | Implement `metrics.RequestCount` with `prometheus` or `opentelemetry`. | ⏳ Planned |
| 527 | Metrics: tests | Test `metrics` increment. | ⏳ Planned |
| 528 | Metrics: docs | Document metrics in `docs/architecture.md`. | ⏳ Planned |
| 529 | Match metrics | Track `creation`, `joins`, `starts`, `completion`, `abandonment`, `duration`, `playerCount`. | ⏳ Planned |
| 530 | Match metrics: implementation | Implement `matchMetrics.CreationCount` etc. | ⏳ Planned |
| 531 | Match metrics: tests | Test match metrics. | ⏳ Planned |
| 532 | Match metrics: docs | Document match metrics in `docs/architecture.md`. | ⏳ Planned |
| 533 | Gameplay metrics | Track `draw` type, `pickup` use, `opening` timing, `joker` use, `win` conditions, `validation` failures. | ⏳ Planned |
| 534 | Gameplay metrics: implementation | Implement `gameplayMetrics.DrawType` etc. | ⏳ Planned |
| 535 | Gameplay metrics: tests | Test gameplay metrics. | ⏳ Planned |
| 536 | Gameplay metrics: docs | Document gameplay metrics. | ⏳ Planned |
| 537 | Bot metrics | Track `botDifficulty`, `decisionLatency`, `searchDepth`, `fallbackRate`, `winRate`, `invalidActionCount`. | ⏳ Planned |
| 538 | Bot metrics: implementation | Implement `botMetrics` with `prometheus`. | ⏳ Planned |
| 539 | Bot metrics: tests | Test bot metrics. | ⏳ Planned |
| 540 | Bot metrics: docs | Document bot metrics. | ⏳ Planned |
| 541 | Product analytics privacy design | Event retention, player consent, anonymization, aggregation, PII boundaries. | ⏳ Planned |
| 542 | Product analytics: implementation | Write `docs/analytics-privacy.md` with `retention` and `PII` rules. | ⏳ Planned |
| 543 | Product analytics: tests | Test analytics never contains `OwnRack` `TileId`. | ⏳ Planned |
| 544 | Product analytics: docs | Document analytics privacy. | ⏳ Planned |
| 545 | Dashboards | Infrastructure health, live matches, gameplay funnel, bot behavior, error rates. | ⏳ Planned |
| 546 | Dashboards: implementation | Create `dashboards/` with `grafana` JSON for `infra` and `gameplay`. | ⏳ Planned |
| 547 | Dashboards: tests | Test dashboards exist. | ⏳ Planned |
| 548 | Dashboards: docs | Document dashboards in `docs/architecture.md`. | ⏳ Planned |
| 549 | Alerting | Alerts for `DB` failures, `Nakama` restart loops, elevated command errors, stuck matches, slow bot actions. | ⏳ Planned |
| 550 | Alerting: implementation | Implement `alerting` with `prometheus` `AlertManager`. | ⏳ Planned |
| 551 | Alerting: tests | Test alerts triggered on `DB` failure. | ⏳ Planned |
| 552 | Alerting: docs | Document alerting. | ⏳ Planned |
| 553 | Distributed tracing | Request/match tracing where feasible, especially for external services and AI workers. | ⏳ Planned |
| 554 | Tracing: implementation | Implement `tracing` with `opentelemetry` `Trace` for `MatchLoop`. | ⏳ Planned |
| 555 | Tracing: tests | Test tracing contains `matchId`. | ⏳ Planned |
| 556 | Tracing: docs | Document tracing. | ⏳ Planned |
| 557 | Analytics validation | Verify event correctness using controlled test matches and ensure no hidden rack details. | ⏳ Planned |
| 558 | Analytics validation: implementation | Run controlled matches with `TestDeterministicSimulation` and check analytics events. | ⏳ Planned |
| 559 | Analytics validation: tests | Test analytics events correct. | ⏳ Planned |
| 560 | Analytics validation: docs | Document validation in `docs/testing.md`. | ⏳ Planned |

---

## Phase 26 — Security, Privacy, and Anti-Cheat

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 561 | Threat model | Document threats: account abuse, token theft, malformed input, collusion, automation abuse, data leakage, DDoS, admin compromise. | ⏳ Planned |
| 562 | Threat model: implementation | Write `docs/threat-model.md` with `threat` matrix. | ⏳ Planned |
| 563 | Threat model: tests | Test `threat-model.md` exists. | ⏳ Planned |
| 564 | Threat model: docs | Document threat model in `README`. | ⏳ Planned |
| 565 | Authentication review | Harden account/session handling, token expiry, refresh flow, password policies, recovery. | ⏳ Planned |
| 566 | Authentication: implementation | Implement `auth` with `tokenExpiry` `15m` and `refresh` flow. | ⏳ Planned |
| 567 | Authentication: tests | Test `token` expiry and `refresh`. | ⏳ Planned |
| 568 | Authentication: docs | Document authentication. | ⏳ Planned |
| 569 | Authorization review | Verify all RPCs, match commands, admin endpoints, storage access enforce least privilege. | ⏳ Planned |
| 570 | Authorization: implementation | Implement `authorization` checks for `MatchJoin`/`MatchLoop`/`MatchSignal`. | ⏳ Planned |
| 571 | Authorization: tests | Test `authorization` with `not_member` and `not_your_turn`. | ⏳ Planned |
| 572 | Authorization: docs | Document authorization. | ⏳ Planned |
| 573 | Protocol fuzzing | Fuzz malformed, oversized, duplicated, reordered, adversarial match commands. | ⏳ Planned |
| 574 | Protocol fuzzing: implementation | Add `go test -fuzz` for `ParseEnvelope` and `ValidatePayload`. | ⏳ Planned |
| 575 | Protocol fuzzing: tests | Test fuzz never panics. | ⏳ Planned |
| 576 | Protocol fuzzing: docs | Document fuzzing in `docs/testing.md`. | ⏳ Planned |
| 577 | Rate limiting | Limits for login, registration, match creation, joining, command spam, chat, bot-related endpoints. | ⏳ Planned |
| 578 | Rate limiting: implementation | Implement `rateLimiter` with `5/s` for `login`, `10/s` for `match` `join`. | ⏳ Planned |
| 579 | Rate limiting: tests | Test rate limiting with `TestRateLimit`. | ⏳ Planned |
| 580 | Rate limiting: docs | Document rate limiting. | ⏳ Planned |
| 581 | Anti-cheat audit | Verify server authority across all actions (racks, stock, meld validation, discard pickup, joker replacement, win detection). | ⏳ Planned |
| 582 | Anti-cheat: implementation | Audit `racks`/`stock`/`meld` validation, ensure no `TileId` leak. | ⏳ Planned |
| 583 | Anti-cheat: tests | Test `TestAntiCheat` with `CheckTileConservation` and `PublicView` redaction. | ⏳ Planned |
| 584 | Anti-cheat: docs | Document anti-cheat in `docs/architecture.md`. | ⏳ Planned |
| 585 | Replay integrity | Sign or hash match event streams to detect tampering in stored replays. | ⏳ Planned |
| 586 | Replay integrity: implementation | Implement `ReplaySigner` with `hmac` `sha256` for `OpServerEvent` stream. | ⏳ Planned |
| 587 | Replay integrity: tests | Test replay signature valid. | ⏳ Planned |
| 588 | Replay integrity: docs | Document replay integrity. | ⏳ Planned |
| 589 | Privacy controls | Account-data export/delete workflows and retention policies per applicable privacy laws. | ⏳ Planned |
| 590 | Privacy controls: implementation | Implement `ExportAccount`/`DeleteAccount` with `GDPR` compliance. | ⏳ Planned |
| 591 | Privacy controls: tests | Test export/delete. | ⏳ Planned |
| 592 | Privacy controls: docs | Document privacy controls. | ⏳ Planned |
| 593 | Dependency security | Dependency scanning, container scanning, security update policy, critical-CVE response. | ⏳ Planned |
| 594 | Dependency: implementation | Add `govulncheck` and `trivy` scan in CI for `backend.so` and `postgres`. | ⏳ Planned |
| 595 | Dependency: tests | Test `govulncheck` no `critical`. | ⏳ Planned |
| 596 | Dependency: docs | Document dependency security. | ⏳ Planned |
| 597 | Security test and review | Penetration testing, tabletop incident exercises, security regression tests. | ⏳ Planned |
| 598 | Security test: implementation | Create `docs/security-review.md` with `pen-test` checklist. | ⏳ Planned |
| 599 | Security test: tests | Test `TestSecurity` with `CheckTileConservation` and `PublicView` redaction. | ⏳ Planned |
| 600 | Security test: docs | Document security test. | ⏳ Planned |

---

## Phase 27 — Player Accounts, Profiles, and Social Foundation

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 601 | Account model | Define guest, registered, linked-provider, verified, suspended, deleted states. | ⏳ Planned |
| 602 | Account model: implementation | Implement `Account{Id, State, Provider, Verified}`. | ⏳ Planned |
| 603 | Account model: tests | Test `Account` states. | ⏳ Planned |
| 604 | Account model: docs | Document account model in `docs/architecture.md`. | ⏳ Planned |
| 605 | Registration and login UX/API | Secure email/password and/or platform-provider login flows. | ⏳ Planned |
| 606 | Registration: implementation | Implement `Register`/`Login` with `bcrypt` and `JWT` `15m`. | ⏳ Planned |
| 607 | Registration: tests | Test `Register` and `Login`. | ⏳ Planned |
| 608 | Registration: docs | Document registration. | ⏳ Planned |
| 609 | Guest accounts | Low-friction guest play with upgrade/link path. | ⏳ Planned |
| 610 | Guest: implementation | Implement `GuestAuth` with `device` `id` and `upgrade` to `registered`. | ⏳ Planned |
| 611 | Guest: tests | Test `Guest` upgrade. | ⏳ Planned |
| 612 | Guest: docs | Document guest accounts. | ⏳ Planned |
| 613 | Player profile | Display name, avatar selection/upload policy, language, region, safe public fields. | ⏳ Planned |
| 614 | Player profile: implementation | Implement `PlayerProfile{DisplayName, Avatar, Language, Region}`. | ⏳ Planned |
| 615 | Player profile: tests | Test `PlayerProfile` validation. | ⏳ Planned |
| 616 | Player profile: docs | Document profile in `docs/terminology.md`. | ⏳ Planned |
| 617 | Profile moderation controls | Validation, reserved names, report flow, avatar limits, admin moderation tools. | ⏳ Planned |
| 618 | Profile moderation: implementation | Implement `ModerateProfile` with `reserved` list and `report`. | ⏳ Planned |
| 619 | Profile moderation: tests | Test moderation. | ⏳ Planned |
| 620 | Profile moderation: docs | Document moderation. | ⏳ Planned |
| 621 | Friends model | Friend requests, accepted friends, blocked users, privacy controls. | ⏳ Planned |
| 622 | Friends: implementation | Implement `Friend{UserId, State, Blocked}`. | ⏳ Planned |
| 623 | Friends: tests | Test `Friend` `Add`/`Accept`/`Block`. | ⏳ Planned |
| 624 | Friends: docs | Document friends. | ⏳ Planned |
| 625 | Presence and invitations | Opt-in online presence and private match invitations. | ⏳ Planned |
| 626 | Presence: implementation | Implement `Presence{UserId, Online, MatchId}` and `Invite`. | ⏳ Planned |
| 627 | Presence: tests | Test `Presence` and `Invite`. | ⏳ Planned |
| 628 | Presence: docs | Document presence. | ⏳ Planned |
| 629 | Recent players | Store recent opponents with add friend/block/report/invite again. | ⏳ Planned |
| 630 | Recent players: implementation | Implement `RecentPlayers` with `MatchId` and `PlayerId` list. | ⏳ Planned |
| 631 | Recent players: tests | Test `RecentPlayers`. | ⏳ Planned |
| 632 | Recent players: docs | Document recent players. | ⏳ Planned |
| 633 | Social privacy tests | Verify blocks, privacy settings, invitations, profile visibility. | ⏳ Planned |
| 634 | Social privacy: implementation | Implement tests for `Block` and `Invite` visibility. | ⏳ Planned |
| 635 | Social privacy: tests | Test `TestSocialPrivacy` with `CheckTileConservation` still `106`. | ⏳ Planned |
| 636 | Social privacy: docs | Document social privacy. | ⏳ Planned |
| 637 | Account lifecycle documentation | Recovery, deletion, moderation, support workflows, audit trails. | ⏳ Planned |
| 638 | Account lifecycle: implementation | Write `docs/account-lifecycle.md` with `recovery` and `deletion`. | ⏳ Planned |
| 639 | Account lifecycle: tests | Test `docs/account-lifecycle.md` exists. | ⏳ Planned |
| 640 | Account lifecycle: docs | Document lifecycle in `README`. | ⏳ Planned |

---

## Phase 28 — Communication, Moderation, and Community Safety

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 641 | Communication scope | Define channels: table chat, private messages, friends chat, lobby chat, system announcements. | ⏳ Planned |
| 642 | Communication: implementation | Write `docs/communication-scope.md` with `channels` matrix. | ⏳ Planned |
| 643 | Communication: tests | Test `docs/communication-scope.md` exists. | ⏳ Planned |
| 644 | Communication: docs | Document communication scope. | ⏳ Planned |
| 645 | Table chat MVP | Implement match/table chat with limits, ordering, reconnection behavior. | ⏳ Planned |
| 646 | Table chat: implementation | Implement `OpClientChat` + `OpServerChat` with `limit 200` chars and `ordering`. | ⏳ Planned |
| 647 | Table chat: tests | Test chat `limit` and `ordering`. | ⏳ Planned |
| 648 | Table chat: docs | Document chat MVP. | ⏳ Planned |
| 649 | Chat persistence decision | Document whether messages are ephemeral, match-scoped, retained, and retention. | ⏳ Planned |
| 650 | Chat persistence: implementation | Write `docs/chat-persistence.md` with `ephemeral` vs `retained`. | ⏳ Planned |
| 651 | Chat persistence: tests | Test `docs/chat-persistence.md` exists. | ⏳ Planned |
| 652 | Chat persistence: docs | Document persistence. | ⏳ Planned |
| 653 | Spam protection | Rate limits, duplicate-message controls, cooldowns, flood detection. | ⏳ Planned |
| 654 | Spam protection: implementation | Implement `SpamProtection` with `5/s` and `duplicate` check. | ⏳ Planned |
| 655 | Spam protection: tests | Test `SpamProtection` with `TestSpam`. | ⏳ Planned |
| 656 | Spam protection: docs | Document spam protection. | ⏳ Planned |
| 657 | Player blocking | Ensure blocked users cannot message/invite/interact. | ⏳ Planned |
| 658 | Player blocking: implementation | Implement `Block` check in `Chat` and `Invite`. | ⏳ Planned |
| 659 | Player blocking: tests | Test `Block` prevents `Chat`. | ⏳ Planned |
| 660 | Player blocking: docs | Document blocking. | ⏳ Planned |
| 661 | Reporting workflow | Report categories, evidence capture, moderation queue, player feedback states. | ⏳ Planned |
| 662 | Reporting: implementation | Implement `Report{Category, Evidence, Queue}`. | ⏳ Planned |
| 663 | Reporting: tests | Test `Report` flow. | ⏳ Planned |
| 664 | Reporting: docs | Document reporting. | ⏳ Planned |
| 665 | Automated content filtering | Configurable profanity/spam filters with language-aware policies. | ⏳ Planned |
| 666 | Content filtering: implementation | Implement `ContentFilter` with `profanity` list per `language`. | ⏳ Planned |
| 667 | Content filtering: tests | Test `ContentFilter` with `TestContentFilter`. | ⏳ Planned |
| 668 | Content filtering: docs | Document content filtering. | ⏳ Planned |
| 669 | Moderator tools | Review queues, mutes, warnings, suspensions, bans, audit logs. | ⏳ Planned |
| 670 | Moderator tools: implementation | Implement `ModeratorTools` with `queue` and `actions`. | ⏳ Planned |
| 671 | Moderator tools: tests | Test `ModeratorTools` with `TestModeratorTools`. | ⏳ Planned |
| 672 | Moderator tools: docs | Document moderator tools. | ⏳ Planned |
| 673 | Appeals and enforcement policy | Document moderation standards, appeal process, staff permissions, evidence retention. | ⏳ Planned |
| 674 | Appeals: implementation | Write `docs/appeals-enforcement.md` with `appeal` flow. | ⏳ Planned |
| 675 | Appeals: tests | Test `docs/appeals-enforcement.md` exists. | ⏳ Planned |
| 676 | Appeals: docs | Document appeals. | ⏳ Planned |
| 677 | Community safety tests | Test rate limits, blocking, reporting, moderation actions, auditability. | ⏳ Planned |
| 678 | Community safety: implementation | Implement `TestCommunitySafety` with `CheckTileConservation` still `106`. | ⏳ Planned |
| 679 | Community safety: tests | Test `TestCommunitySafety` with `TestRateLimit` etc. | ⏳ Planned |
| 680 | Community safety: docs | Document community safety. | ⏳ Planned |

---

## Phase 29 — Economy, Progression, and Cosmetics

*Strictly non-gambling. No cash wagering, real-money stakes, loot boxes, or gambling-like mechanics.*

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 681 | Economy principles | Document non-gambling rules, fair-play constraints, virtual-currency policy, regional compliance, parental controls. | ⏳ Planned |
| 682 | Economy principles: implementation | Write `docs/economy-principles.md` with `non-gambling` rules. | ⏳ Planned |
| 683 | Economy principles: tests | Test `docs/economy-principles.md` exists. | ⏳ Planned |
| 684 | Economy principles: docs | Document economy principles in `README`. | ⏳ Planned |
| 685 | Player progression | Experience, levels, badges, non-competitive achievement tracking. | ⏳ Planned |
| 686 | Player progression: implementation | Implement `PlayerProgression{Exp, Level, Badges}`. | ⏳ Planned |
| 687 | Player progression: tests | Test `PlayerProgression` with `CheckTileConservation` still `106`. | ⏳ Planned |
| 688 | Player progression: docs | Document progression. | ⏳ Planned |
| 689 | Match rewards | Non-monetary progression rewards for completed matches, good sportsmanship, tutorials. | ⏳ Planned |
| 690 | Match rewards: implementation | Implement `MatchRewards` with `exp` and `badge` on `RoundComplete`. | ⏳ Planned |
| 691 | Match rewards: tests | Test `MatchRewards` with `TestMatchRewards`. | ⏳ Planned |
| 692 | Match rewards: docs | Document match rewards. | ⏳ Planned |
| 693 | Cosmetics model | Cosmetic-only items: table themes, tile skins, avatars, emotes, profile frames, victory effects. | ⏳ Planned |
| 694 | Cosmetics: implementation | Write `docs/cosmetics-model.md` with `tableThemes` and `tileSkins`. | ⏳ Planned |
| 695 | Cosmetics: tests | Test `docs/cosmetics-model.md` exists. | ⏳ Planned |
| 696 | Cosmetics: docs | Document cosmetics. | ⏳ Planned |
| 697 | Inventory service | Ownership, unlocks, equipped items, grants, revocations, audit history. | ⏳ Planned |
| 698 | Inventory: implementation | Implement `Inventory{PlayerId, Items[], Equipped}`. | ⏳ Planned |
| 699 | Inventory: tests | Test `Inventory` with `TestInventory`. | ⏳ Planned |
| 700 | Inventory: docs | Document inventory. | ⏳ Planned |
| 701 | Daily/weekly challenges | Optional goals rewarding engagement without pressuring risky spending. | ⏳ Planned |
| 702 | Challenges: implementation | Implement `Challenges{Daily, Weekly}` with `exp` rewards. | ⏳ Planned |
| 703 | Challenges: tests | Test `Challenges` with `TestChallenges`. | ⏳ Planned |
| 704 | Challenges: docs | Document challenges. | ⏳ Planned |
| 705 | Achievement system | Rule-based achievements, progress tracking, notifications, anti-exploit checks. | ⏳ Planned |
| 706 | Achievements: implementation | Implement `Achievements{PlayerId, Achievements[]}` with `Rule` checks. | ⏳ Planned |
| 707 | Achievements: tests | Test `Achievements` with `TestAchievements`. | ⏳ Planned |
| 708 | Achievements: docs | Document achievements. | ⏳ Planned |
| 709 | Store design, if applicable | If monetization approved, cosmetic purchases with transparent pricing, receipts, refunds, parental controls, no pay-to-win. | ⏳ Planned |
| 710 | Store: implementation | Write `docs/store-design.md` with `pricing` and `receipts`. | ⏳ Planned |
| 711 | Store: tests | Test `docs/store-design.md` exists. | ⏳ Planned |
| 712 | Store: docs | Document store design. | ⏳ Planned |
| 713 | Economy abuse prevention | Detect duplicate grants, replayed purchase events, reward farming, inventory inconsistencies. | ⏳ Planned |
| 714 | Economy abuse: implementation | Implement `EconomyAbusePrevention` with `duplicate` and `farming` checks. | ⏳ Planned |
| 715 | Economy abuse: tests | Test `EconomyAbusePrevention` with `TestEconomyAbuse`. | ⏳ Planned |
| 716 | Economy abuse: docs | Document abuse prevention. | ⏳ Planned |
| 717 | Economy tests and compliance review | Test transactions, entitlement handling; legal/compliance review before release. | ⏳ Planned |
| 718 | Economy tests: implementation | Create `docs/economy-tests.md` with `transaction` checklist. | ⏳ Planned |
| 719 | Economy tests: docs | Document economy tests. | ⏳ Planned |
| 720 | Economy tests: polish | Polish `README` `Next Steps`. | ⏳ Planned |

---

## Phase 30 — Competitive Play, Rankings, and Tournaments

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 721 | Competitive ruleset | Define ranked match eligibility, player count, disconnect policy, bot policy, scoring format, seasonal structure. | ⏳ Planned |
| 722 | Competitive ruleset: implementation | Write `docs/competitive-ruleset.md` with `eligibility` and `season` rules. | ⏳ Planned |
| 723 | Competitive ruleset: tests | Test `docs/competitive-ruleset.md` exists. | ⏳ Planned |
| 724 | Competitive ruleset: docs | Document competitive ruleset. | ⏳ Planned |
| 725 | Rating algorithm | Select and implement Elo/Glicko/TrueSkill-like logic for multiplayer. | ⏳ Planned |
| 726 | Rating algorithm: implementation | Implement `Rating{PlayerId, Value, Deviation}` and `UpdateRating` after `RoundComplete`. | ⏳ Planned |
| 727 | Rating algorithm: tests | Test `Rating` updates with `CheckTileConservation` still `106`. | ⏳ Planned |
| 728 | Rating algorithm: docs | Document rating algorithm in `docs/architecture.md`. | ⏳ Planned |
| 729 | Rating event pipeline | Record rating-relevant results atomically and make recalculation/recovery possible. | ⏳ Planned |
| 730 | Rating event: implementation | Implement `RatingEvent{MatchId, PlayerIds, Winner, Timestamp}`. | ⏳ Planned |
| 731 | Rating event: tests | Test `RatingEvent` with `TestRatingEvent`. | ⏳ Planned |
| 732 | Rating event: docs | Document rating event pipeline. | ⏳ Planned |
| 733 | Leaderboards | Global, regional, friends, seasonal, game-mode leaderboards. | ⏳ Planned |
| 734 | Leaderboards: implementation | Implement `Leaderboard{Type, PlayerIds, Scores}` and `UpdateLeaderboard` after `RoundComplete`. | ⏳ Planned |
| 735 | Leaderboards: tests | Test `Leaderboard` with `TestLeaderboard`. | ⏳ Planned |
| 736 | Leaderboards: docs | Document leaderboards. | ⏳ Planned |
| 737 | Ranked matchmaking | Match via rating, region, latency, queue duration, anti-smurf controls. | ⏳ Planned |
| 738 | Ranked matchmaking: implementation | Implement `RankedMatchmaking` with `rating` and `region` filters. | ⏳ Planned |
| 739 | Ranked matchmaking: tests | Test `RankedMatchmaking` with `TestRankedMatchmaking`. | ⏳ Planned |
| 740 | Ranked matchmaking: docs | Document ranked matchmaking. | ⏳ Planned |
| 741 | Season system | Season start/end, soft reset policy, rewards, historical ranks, announcement flow. | ⏳ Planned |
| 742 | Season: implementation | Implement `Season{Start, End, SoftReset, Rewards}`. | ⏳ Planned |
| 743 | Season: tests | Test `Season` with `TestSeason`. | ⏳ Planned |
| 744 | Season: docs | Document season system. | ⏳ Planned |
| 745 | Tournament design | Formats: scheduled tables, Swiss, bracket, elimination, points-based. | ⏳ Planned |
| 746 | Tournament design: implementation | Write `docs/tournament-design.md` with `formats` and `rules`. | ⏳ Planned |
| 747 | Tournament design: tests | Test `docs/tournament-design.md` exists. | ⏳ Planned |
| 748 | Tournament design: docs | Document tournament design. | ⏳ Planned |
| 749 | Tournament engine | Registration, seating, round creation, advancement, standings, tie-breaking. | ⏳ Planned |
| 750 | Tournament engine: implementation | Implement `TournamentEngine` with `Register`, `Seat`, `CreateRound`, `Advance`, `Standings`. | ⏳ Planned |
| 751 | Tournament engine: tests | Test `TournamentEngine` with `TestTournamentEngine`. | ⏳ Planned |
| 752 | Tournament engine: docs | Document tournament engine. | ⏳ Planned |
| 753 | Tournament integrity | Check-in, disconnect rules, bot restrictions, anti-collusion signals, staff intervention, replay review. | ⏳ Planned |
| 754 | Tournament integrity: implementation | Implement `TournamentIntegrity` with `CheckIn` and `Disconnect` rules. | ⏳ Planned |
| 755 | Tournament integrity: tests | Test `TournamentIntegrity` with `TestTournamentIntegrity`. | ⏳ Planned |
| 756 | Tournament integrity: docs | Document tournament integrity. | ⏳ Planned |
| 757 | Tournament rewards | Cosmetic/badge rewards and transparent eligibility. Avoid real-money prizes unless legal/licensed. | ⏳ Planned |
| 758 | Tournament rewards: implementation | Implement `TournamentRewards` with `cosmetic` and `badge` grants. | ⏳ Planned |
| 759 | Tournament rewards: tests | Test `TournamentRewards` with `TestTournamentRewards`. | ⏳ Planned |
| 760 | Tournament rewards: docs | Document tournament rewards. | ⏳ Planned |
| 761 | Competitive testing | Simulate rating updates, tournament progression, dropouts, ties, rollback. | ⏳ Planned |
| 762 | Competitive testing: implementation | Create `TestCompetitive` with `CheckTileConservation` still `106`. | ⏳ Planned |
| 763 | Competitive testing: tests | Test `TestCompetitive` with `TestRatingUpdates` etc. | ⏳ Planned |
| 764 | Competitive testing: docs | Document competitive testing. | ⏳ Planned |
| 765 | Competitive launch readiness | Load test queues/leaderboards, audit integrity, publish rules, run limited beta season. | ⏳ Planned |
| 766 | Competitive launch: implementation | Write `docs/competitive-launch-readiness.md` with `load test` and `audit` checklist. | ⏳ Planned |
| 767 | Competitive launch: tests | Test `competitive launch readiness` exists. | ⏳ Planned |
| 768 | Competitive launch: docs | Document competitive launch. | ⏳ Planned |
| 769 | Competitive launch: polish | Polish `README` `Next Steps`. | ⏳ Planned |
| 770 | Competitive launch: tag | Tag `competitive-rc1` when all 721–769 pass. | ⏳ Planned |

---

## Phase 31 — Client Platforms, UX, Accessibility, and Localization

| Day | Focus | Deliverables | Status |
|---:|---|---|---|
| 771 | Client platform strategy | Decide web/desktop/Android/iOS priorities; define shared protocol/client SDK strategy. | ⏳ Planned |
| 772 | Client platform: implementation | Write `docs/client-platform-strategy.md` with `web` first, `mobile` second. | ⏳ Planned |
| 773 | Client platform: tests | Test `docs/client-platform-strategy.md` exists. | ⏳ Planned |
| 774 | Client platform: docs | Document client platform strategy. | ⏳ Planned |
| 775 | Design system | Reusable typography, colours, spacing, buttons, modals, notifications, tile components, table layouts. | ⏳ Planned |
| 776 | Design system: implementation | Create `design-system/` with `typography` and `colours`. | ⏳ Planned |
| 777 | Design system: tests | Test `design-system` exists. | ⏳ Planned |
| 778 | Design system: docs | Document design system in `docs/architecture.md`. | ⏳ Planned |
| 779 | Rummy table UX | Core table interaction: rack sorting, drag/drop, tile selection, undo-before-submit, turn feedback. | ⏳ Planned |
| 780 | Rummy table UX: implementation | Implement `RummyTable` with `rackSorting` and `dragDrop`. | ⏳ Planned |
| 781 | Rummy table UX: tests | Test `RummyTable` with `TestRummyTable`. | ⏳ Planned |
| 782 | Rummy table UX: docs | Document table UX. | ⏳ Planned |
| 783 | Meld-building UX | Accessible creation of runs/sets, explicit joker assignment, invalid-state feedback, server-confirmed actions. | ⏳ Planned |
| 784 | Meld-building UX: implementation | Implement `MeldBuilder` with `joker` assignment and `validation`. | ⏳ Planned |
| 785 | Meld-building UX: tests | Test `MeldBuilder` with `TestMeldBuilder`. | ⏳ Planned |
| 786 | Meld-building UX: docs | Document meld-building UX. | ⏳ Planned |
| 787 | Responsive layouts | Desktop, tablet, mobile layouts without reducing rule clarity. | ⏳ Planned |
| 788 | Responsive: implementation | Implement `responsive` with `desktop`/`tablet`/`mobile` `media` queries. | ⏳ Planned |
| 789 | Responsive: tests | Test `responsive` with `TestResponsive`. | ⏳ Planned |
| 790 | Responsive: docs | Document responsive layouts. | ⏳ Planned |
| 791 | Accessibility baseline | Keyboard control, screen-reader labels, focus handling, colour-blind-safe tile identifiers, scalable text, reduced motion. | ⏳ Planned |
| 792 | Accessibility: implementation | Implement `accessibility` with `keyboard` and `screenReader` support. | ⏳ Planned |
| 793 | Accessibility: tests | Test `accessibility` with `TestAccessibility`. | ⏳ Planned |
| 794 | Accessibility: docs | Document accessibility. | ⏳ Planned |
| 795 | Localization foundation | Externalize UI text, support Romanian and English first, pluralization/date/number rules. | ⏳ Planned |
| 796 | Localization: implementation | Create `locales/` with `en.json` and `ro.json` and `pluralization` logic. | ⏳ Planned |
| 797 | Localization: tests | Test `localization` with `TestLocalization`. | ⏳ Planned |
| 798 | Localization: docs | Document localization. | ⏳ Planned |
| 799 | Connection UX | Reconnecting states, resync UI, action-pending feedback, network error handling, match recovery. | ⏳ Planned |
| 800 | Connection UX: implementation | Implement `ConnectionUX` with `reconnecting` and `resync` states. | ⏳ Planned |
| 801 | Connection UX: tests | Test `ConnectionUX` with `TestConnectionUX`. | ⏳ Planned |
| 802 | Connection UX: docs | Document connection UX. | ⏳ Planned |
| 803 | Onboarding and tutorial | Guided rules introduction, practice match with beginner bots, contextual help, rule glossary. | ⏳ Planned |
| 804 | Onboarding: implementation | Create `tutorial` with `guided` `rules` and `practice` `match`. | ⏳ Planned |
| 805 | Onboarding: tests | Test `Onboarding` with `TestOnboarding`. | ⏳ Planned |
| 806 | Onboarding: docs | Document onboarding. | ⏳ Planned |
| 807 | Notifications | Opt-in invitations, turn reminders, friend activity, tournaments, moderation actions. | ⏳ Planned |
| 808 | Notifications: implementation | Implement `Notifications` with `invite` and `turnReminder`. | ⏳ Planned |
| 809 | Notifications: tests | Test `Notifications` with `TestNotifications`. | ⏳ Planned |
| 810 | Notifications: docs | Document notifications. | ⏳ Planned |
| 811 | Client QA automation | End-to-end client tests for login, match creation, joining, draw/discard, melding, reconnection, error states. | ⏳ Planned |
| 812 | Client QA: implementation | Create `client-qa/` with `e2e` tests for `login` and `match`. | ⏳ Planned |
| 813 | Client QA: tests | Test `ClientQA` with `TestClientQA`. | ⏳ Planned |
| 814 | Client QA: docs | Document client QA. | ⏳ Planned |
| 815 | UX beta and polish | Usability testing, fix high-impact friction, improve performance, prepare release notes. | ⏳ Planned |
| 816 | UX beta: implementation | Run `ux beta` with `5` users and `fix` friction. | ⏳ Planned |
| 817 | UX beta: tests | Test `UX beta` with `TestUXBeta`. | ⏳ Planned |
| 818 | UX beta: docs | Document UX beta. | ⏳ Planned |
| 819 | Polish | Polish `README` `Next Steps` to reflect UX beta done. | ⏳ Planned |
| 820 | Polish: final | Tag `ux-rc1` when all 771–819 pass. | ⏳ Planned |

---

## Extended Roadmap Summary (Future)

| Phase | Days | Primary Outcome | Status |
|---:|---|---|---|
| 17 | 147–240 | Future product features (scoring, social, matchmaking, bots, production ops, polish) | ⏳ Not started |
| 18 | 241–480 | AI Bots: foundation through production AI operations (Phases 18–23 above) | ⏳ Not started |
| 24 | 481–520 | Production infrastructure and deployment | ⏳ Not started |
| 25 | 521–560 | Observability, analytics, telemetry | ⏳ Not started |
| 26 | 561–600 | Security, privacy, anti-cheat | ⏳ Not started |
| 27 | 601–640 | Accounts, profiles, social foundation | ⏳ Not started |
| 28 | 641–680 | Chat, moderation, community safety | ⏳ Not started |
| 29 | 681–720 | Safe progression, cosmetics, optional economy (non-gambling) | ⏳ Not started |
| 30 | 721–770 | Ranked play, leaderboards, tournaments | ⏳ Not started |
| 31 | 771–820 | Multi-platform client, UX, accessibility, localization | ⏳ Not started |

By **Day 820**, the project evolves from the current `rummy-mvp-rc1` authoritative backend (`Phases 1–16` done) into a production-ready, AI-supported, secure, social, competitive, accessible multiplayer platform. Execution remains Handmade Hero: one vertical slice per day, `make check` + `docker compose build` green, focused commit, push.

---

*This roadmap lists days that **will be** implemented. For what **is** implemented, see `docs/IMPLEMENTED.md`. Last updated: 2026-08-26 at `main@2c5cca7` (and `rummy-mvp-rc1` `cfce62e`).*

