# Romanian Tile Rummy — Rules Decisions

**Source:** Romanian Tile Rummy as documented by Pagat (`https://www.pagat.com/rummy/romtile.html`) — authoritative baseline per `AGENTS.md:3`.
**Scope:** Server-authoritative 2–4 players, Nakama Go runtime (`go 1.23.5`), Docker local dev. Remi Online is inspiration only for social table experience — no branding/assets copied.
**Status:** Day 8 — Phase 2 Rules source documentation. Decisions here are binding for all `internal/*` pure rules modules and `main.go` match handler. Ambiguities are marked `TODO(product):` and must not be guessed.

---

## 1. Source Rules (confirmed, verbatim intent)

### 1.1 Tile set — `AGENTS.md:35`
- 106 tiles: 4 colours (red, yellow, blue, black) × ranks 1–13 × 2 copies = 104 numbered + 2 jokers (“Joly”).
- Each tile instance has a **unique immutable ID** (`TileInstanceId`) even when colour+rank identical — required by `AGENTS.md:41` and tile conservation invariant.

### 1.2 Player count & direction — `AGENTS.md:43`
- 2–4 players, anticlockwise. In networked game, seat order is deterministic join order `0..n-1`; anticlockwise = `nextSeat = (current+1) % n` (low index → next clockwise on screen, but server logs it as anticlockwise per Pagat). Documented here to avoid UI confusion.

### 1.3 Valid melds — `AGENTS.md:48`

**Run (`suita`):**
- Same colour, at least 3 consecutive values.
- `1-2-3` valid, `12-13-1` valid, `13-1-2` invalid, Ace never in middle.
- No silent reinterpretation of a joker at a run end once placed.

**Set (`terta`):**
- 3 or 4 tiles, same rank, each a different colour.

**Joker:**
- Substitutes for a needed tile; **real:joker ≥ 2:1** after substitution (e.g. 3-tile meld may have at most 1 joker, 5-tile meld at most 1 joker? Actually 3 tiles → 1 joker needs 2 real, 4 tiles → 1 joker needs 2 real, 6 tiles → 2 jokers needs 4 real). Enforced in pure validator.
- Joker’s `representedTile` (colour+rank) is explicit and immutable once melded unless a legal replacement occurs.
- Joker at run end preserves its declared represented rank — never silently shift `12-13-J` from `12-13-1` to `12-13-14`.

### 1.4 Initial meld — `AGENTS.md:63`
- Must be from rack only, include **at least one run**, total **≥50 points**.
- Point values:

| Tile | Points |
|------|--------|
| 2–9 | 5 each |
| 10–13 | 10 each |
| Ace `1-2-3…` | 5 |
| Ace `…12-13-1` | 10 |
| Ace set (three Aces, three colours) | 25 each |
| Joker (any context) | value of represented tile |

- Scoring is pure and isolated (`internal/rules/scoring`); tests must cover exactly 50 accepted / 49 rejected, low/high Ace, ace-set, joker.

### 1.5 Turn flow — `AGENTS.md:77`
- Opening player starts with **15** tiles, discards one tile to start. Others 14.
- Normal turns: **draw → meld (optional) → discard**; must end with exactly one discard unless win rule prevents it.
- Before opening, must draw from stock at turn start.
- After opening, allowed:
  - Create additional valid melds.
  - Extend own or others’ public melds.
  - Take **immediately previous discard** at start of turn instead of stock (only if opened, not opening discard).
  - Take **earlier discard** only if can immediately meld it with **exactly two rack tiles** into a valid set/run, then pick up **all later discards** into rack.
  - Replace a joker only under precise legal conditions and **immediately** use that joker in a new meld with **two rack tiles**.
- The **initial discarded tile is permanently unavailable** for pickup and must be represented distinctly in state (flag `IsOpeningDiscard` on first discard).

### 1.6 End condition — `AGENTS.md:91`
- A player wins when they have **no tiles left after legal play**, subject to final-discard behavior. Source excerpt ends early — see ambiguity §6.

---

## 2. MVP Simplifications (deliberate, not source)

Recorded per `AGENTS.md:111`:

1. **Shuffle:** Uniform Fisher–Yates over 106 unique IDs with injectable seeded `Rand` (deterministic for tests), not physical 15 stacks of seven.
2. **Dealer / opening player:** Deterministic: lowest seat index (host / first joiner) is opening player. Dealer rotation deferred.
3. **Deal:** Opening player 15 tiles, others 14, remainder = stock (ordered array, top = last shuffled element). No “wall” geometry.
4. **Opening discard:** First action is discard from 15 to 14; that discard is flagged unavailable.
5. **Stock exhaustion:** Deferred to §6 ambiguity — MVP will be defined deterministically (see §6.2) rather than improvising silently.

These are not Pagat-authentic but are testable, deterministic, and documented as MVP.

---

## 3. Valid Melds — Canonical Representations (for validators)

**TileInstance:** `ID` (unique), `Colour` (`red|yellow|blue|black`), `Rank` (`1..13`), `IsJoker` bool. Joker has no colour/rank until declared.

**Meld:** `ID` stable UUID, `Kind` (`run|set`), `Tiles` ordered, `JokerRepresentations` map `TileID → {Colour,Rank}` explicit.

Validators are pure `func ValidateRun(tiles []TileInstance, reps map) error` / `ValidateSet` returning structured `ValidationError` with code + field, never bool-only. Ratio check: `real*2 >= joker*?` Actually requirement `real >= 2*joker` — e.g. `len(real) >= 2*len(joker)`. So 3-tile meld max 1 joker, 5-tile meld max 1 joker? 6-tile meld max 2 jokers. Tests will enumerate.

Ace handling: Run is **not circular**; it is either low-starting `1-2-3…` or high-ending `12-13-1` (or longer `10-11-12-13-1` etc.) but never wraps `13-1-2`. Implementation: generate canonical sequences `1..13` and also `12,13,1` as endpoint, then check consecutive without gaps except jokers filling exactly the declared represented rank.

---

## 4. Scoring — Isolated Rules Module

- `ScoreTile(tile, context)` → points, where context is `isAceLowRun`/`isAceHighRun`/`isAceSet` and `represented` for joker.
- `ScoreRun` / `ScoreSet` sum with Joker delegation.
- Initial batch: `ValidateInitialMeld(racks, proposedMelds)` checks all tiles owned, each meld valid, at least one run, total ≥50, no duplicate TileID across melds, atomic (no partial mutation).

Tests: low Ace `A-2-3` (5), high Ace `Q-K-A` (10), ace set `A(red)+A(blue)+A(black)` = 75, joker as `1` in `1-2-3` = 5, exactly 50 passes.

---

## 5. Turn Flow — Phases and Legal Commands

**Seats:** `Seat 0..n-1` join order, `Current == (Current+1)%n` is next anticlockwise.

**Phases:**
- `Waiting` — lobby, 2–4 players.
- `OpeningDiscard` — opening player must `DISCARD` one tile (15→14).
- `MustDraw` — active player’s turn start; legal: `DRAW_STOCK`, `DRAW_PREVIOUS_DISCARD` (if opened, not opening), `PICKUP_DISCARD_FOR_MELD` (if opened).
- `MeldOrDiscard` — after draw; legal: `MELD_INITIAL` (once, if not opened), `MELD_NEW`, `EXTEND_MELD`, `REPLACE_JOKER`, `DISCARD` (must end turn).
- `RoundComplete` — winner empty rack, no further gameplay.

Invariants per move: tile conservation `racks+stock+discardRow+tableMelds == 106`, no duplicate TileID, atomic command handling.

---

## 6. Ambiguities — Explicit Decisions with TODOs

Per `AGENTS.md:92` format: source → ambiguity → deterministic MVP behavior → tests affected → TODO.

### 6.1 Final discard / closing behavior

- **Source rule:** “A player wins when they have no tiles left after legal play” (`AGENTS.md:92`).
- **Ambiguity:** Does the winning move require a final discard (leaving 0 after discard) or can a player win by melding/emptying rack without discarding? Pagat excerpt truncated; some Rummy variants require discard, others allow melding out.
- **Proposed deterministic MVP:** **Require final discard to reach 0** unless the winning meld itself consumes the last tile and the player has already opened and the discard would be the winning player’s last tile — then **allow winning without final discard** but still emit `ROUND_COMPLETE` after the melding action. In other words: `rack == 0` after either `DISCARD` or after a successful `MELD_*`/`EXTEND`/`REPLACE_JOKER` that empties rack is a win. The server will check `rackEmpty` after **every** state-mutating command in `MeldOrDiscard` as well as after `DISCARD`.
- **Tests affected:** `TestWinAfterMeldWithoutDiscard`, `TestWinAfterDiscardToZero`, `TestNoGameplayAfterRoundComplete`, tile conservation at win.
- **TODO(product):** Confirm closing rule with product/design — does this match intended “Remi Etalat” social table behaviour? Reference Remi Online distinction `Remi Etalat` vs `Remi Pe Tabla` (targeting tile-rummy/meld-on-table). Update this doc and `internal/match/win.go` when confirmed.

### 6.2 Stock exhaustion / dead round

- **Source rule:** Normal draw is from wall/stock; no rule for empty stock given.
- **Ambiguity:** What happens when stock becomes empty? Does the round draw? Is it a dead round? Can players still use discard pickups?
- **Proposed MVP:** **Stock empty → player in `MustDraw` may only perform discard-pickup moves (`DRAW_PREVIOUS_DISCARD` or `PICKUP_DISCARD_FOR_MELD`) if eligible; if no eligible discard pickup, the player must pass? For MVP simplicity, if stock is empty and no discard pickup is legal, the server transitions to `ROUND_COMPLETE` with no winner (“dead round”) and emits `Reason: stockExhausted`. Tile conservation still holds (discard row + racks + melds = 106). This avoids infinite loop.
- **Tests affected:** `TestDrawStockEmptyTransitions`, `TestStockExhaustedNoWinner`, `TestDiscardPickupStillAllowedWhenStockEmpty`.
- **TODO(product):** Decide if stock exhaustion should be a reshuffle or true dead round; affects multi-round scoring.

### 6.3 Doubla / 15 stacks / exposed-tile

- **Source note:** `AGENTS.md:101` lists “15 stacks of seven”, “Doubla”, “exposed-tile bonus” as deferred.
- **Ambiguity:** Not enough source excerpt to define deterministic behavior; do not implement in MVP.
- **Proposed MVP:** Do not implement; these mechanics are not available via protocol (`OPCODE` not defined) and validators will reject unknown melds.
- **Tests affected:** N/A — no tests for deferred features.
- **TODO(product):** Revisit post-MVP (Phase 17).

### 6.4 Joker replacement in sets — exact missing colours

- **Source rule:** “Replace a joker in a three-tile set by supplying both missing colours of that rank, then immediately use recovered joker with two rack tiles in a new valid meld” (`AGENTS.md:87`).
- **Ambiguity:** Does the replacement require the **exact two missing colours** (i.e., if set is `5 red + 5 blue + J(black)` missing `yellow`, player must supply `5 yellow`? Actually three-tile set has 1 joker → 2 real colours present, 1 colour missing, but rule says “both missing colours” suggests a 3-tile set with 1 joker has 2 missing colours? Need clarification — the excerpt may refer to 2-colour set?).
- **Proposed MVP:** Require **exactly the colour(s) of the represented joker(s)**. For a 3-tile set `5 red + 5 yellow + J(black→5 black)` missing `blue`, replacing requires `5 blue` (the represented colour) — not “both missing”. For a 4-tile set `5 red +5 yellow+ J(blue)+J(black)`? Actually 4-tile set cannot have joker due to ratio (needs 2 real per joker → 4 tiles max 1 joker). So only 3-tile set with 1 joker is legal. Replace requires the single missing colour. Documented as `Run joker: exact represented tile; Set joker: exact missing colour`.
- **Tests affected:** `TestReplaceSetJokerValidMissingColour`, `TestReplaceSetJokerWrongColourRejected`, `TestReplaceSetJokerNewMeldRequiresTwoRackTiles`.
- **TODO(product):** Confirm pagat “both missing colours” wording — does it imply a 2-tile set with joker? Keep this doc as MVP until confirmed.

---

## 7. Intentionally Deferred (excluded from MVP)

Per `AGENTS.md:100`:

- Exact physical wall construction/dealing ritual (15 stacks of seven)
- “Doubla” announcements/exchanges
- Exposed-tile bonus
- Detailed multi-round scoring and dealer rotation
- Spectators, chat, friends, player profiles, tournaments
- Bots/AI opponents
- Animation and polished UI

These opcodes will not be accepted; added only after baseline round is playable and tested (`docs/project-baseline.md:5`).

---

## 8. State Invariants (must hold after every command)

- **Tile conservation:** Every of the 106 `TileInstanceId` appears exactly once across `stock + each rack + discardRow + tableMelds` (and temporary holding for atomic checks).
- **No duplicate TileID** in any meld/rack.
- **Joker representation immutable** once melded (unless via `REPLACE_JOKER` flow).
- **Opening discard unavailable:** First discard entry has `IsOpeningDiscard: true` and is never selectable by `DRAW_PREVIOUS_DISCARD` or `PICKUP_DISCARD_FOR_MELD`.
- **Phase legality:** Only current `Seat` in correct `TurnPhase` may act.
- **Atomicity:** Invalid commands do not partially mutate state.

---

## 9. Next Artefacts

- `Phase 2 Day 9`: Domain terminology doc (`tile, rack, stock, discard row, meld…`).
- `Phase 2 Day 10`: Go types `TileColour`, `TileRank`, `TileInstanceId`, `Joker`, `Meld`, `RunMeld`, `SetMeld`, `Seat`, `GamePhase`, `TurnPhase` in `internal/rules/tile`.
- Scoring isolated in `internal/rules/scoring` with table above.

---

*Last updated: 2026-08-25 (Day 8 Phase 2). TODOs require product confirmation — this doc is the blocking record before any rule code is merged.*
