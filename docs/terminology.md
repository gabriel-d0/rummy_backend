# Domain Terminology — Shared Vocabulary

**Scope:** Canonical names for Romanian Tile Rummy implementation. All `internal/*` Go packages, `main.go` match handler, tests, and `docs/rules-decisions.md` **must** use these terms verbatim to avoid translation drift (`tila`, `joker`, `suita` are documented but code uses English). This is Day 9 — pure terminology, no types yet (types land Day 10).

## Core Tiles

| Term | Definition | Code hint (Day 10) | Notes |
|------|------------|--------------------|-------|
| **Tile** | Abstract physical tile: colour + rank or joker. Two tiles can be face-identical but are distinct instances. | `Tile` | 106 total per `docs/rules-decisions.md:1.1`. |
| **TileInstanceId** | Unique immutable ID for a single tile copy. String/UUID or monotonic int, stable across shuffle/deal. | `TileInstanceId` | Required by `AGENTS.md:41` for conservation checks. Never reuse. |
| **TileInstance** | Concrete tile at a location: `ID` + `Colour`/`Rank` or `IsJoker`. | `TileInstance` | What is shuffled, dealt, drawn, melded, discarded. |
| **Colour** | One of 4: `red`, `yellow`, `blue`, `black`. | `TileColour` | No other colours; jokers have no colour until represented. |
| **Rank** | 1–13 (Ace=1, 2…10, J=11, Q=12, K=13). In rules, Ace dual-low/high handled by meld validator, not by rank mutation. | `TileRank` | Integer 1–13. |
| **Joker (“Joly”)** | Substitute tile. Has `IsJoker` true, `ID` unique, and when melded carries `RepresentedTile` (colour+rank) explicit, immutable. | `JokerTile` / `IsJoker` | Ratio `real ≥ 2*joker` per `docs/rules-decisions.md:1.3`. |

## Player & Table State

| Term | Definition | Visibility | Code hint |
|------|------------|------------|-----------|
| **PlayerId** | Nakama `userId` string for authenticated participant. | Private ↔ server | `PlayerId` |
| **Seat** | Deterministic position `0..n-1` by join order; `n` = 2–4. | Public (player list order) | `Seat` |
| **Anticlockwise** | Next seat = `(current+1)%n` per `docs/rules-decisions.md:1.2`. Maps to visually clockwise progression but server logs anticlockwise per Pagat — do not invert. | Public | `NextSeat()` helper |
| **Rack** | Private tiles held by a player (`[]TileInstance`). Never sent to other players. | Private to owner | `Rack` |
| **Stock** (also “wall”) | Remaining shuffled tiles not dealt; draw source. Ordered stack, top = last element. | Public count only | `Stock` |
| **Discard Row** (also Discard Pile) | Ordered public history of discards, oldest→newest. Index matters for `PICKUP_DISCARD_FOR_MELD`. | Public ordered list | `DiscardRow` |
| **Opening Discard** | First discard entry from opening player’s 15 tiles. Flag `IsOpeningDiscard` true, permanently unavailable for `DRAW_PREVIOUS_DISCARD` / `PICKUP_DISCARD_FOR_MELD`. | Public (flagged) | `IsOpeningDiscard` |
| **Table Melds** | Public melds placed on table (`[]Meld`). Each has stable `ID`, never moves between owners. | Public | `TableMelds` |
| **Meld** | Generic set or run placed on table. | Public | `Meld` |
| **Run (`suita`)** | Same colour, ≥3 consecutive ranks, may include jokers with explicit reps; allowed `1-2-3` and `12-13-1`, not `13-1-2`. | Public | `RunMeld` |
| **Set (`terta`)** | Same rank, 3–4 tiles, each different colour, may include at most 1 joker (ratio). | Public | `SetMeld` |
| **Opening Meld** | Player’s first `MELD_INITIAL` batch: racks only, ≥1 run, total ≥50 points per scoring table. Marks `HasOpened`. | Public | `HasOpened` flag per player |
| **HasOpened** | Per-player bool, false until `MELD_INITIAL` succeeds. Gates `DRAW_PREVIOUS_DISCARD` etc. | Public bool | `PlayerState.HasOpened` |

## Turn & Round

| Term | Definition | Code hint |
|------|------------|-----------|
| **GamePhase** | Top-level round state: `Waiting` (lobby 2–4), `OpeningDiscard` (first discard), `Playing` (normal loop), `RoundComplete` (winner or dead round). | `GamePhase` |
| **TurnPhase** | Within a player’s turn: `MustDraw` (start — must `DRAW_STOCK` or discard pickup) → `MeldOrDiscard` (after draw — may meld/extend/replace-joker, must end with `DISCARD` unless win) → (next player `MustDraw`). | `TurnPhase` |
| **Current Turn** | `{Seat, PlayerId, Phase}` triple that the server validates. Only `Current` may act. | `CurrentTurn` |
| **Draw** | Take one tile from `Stock` (`DRAW_STOCK`) or discard row (`DRAW_PREVIOUS_DISCARD`, `PICKUP_DISCARD_FOR_MELD`). | `Draw` |
| **Discard** | Place one owned `TileInstanceId` onto `DiscardRow` (`DISCARD`). Opening discard is `DISCARD` in `OpeningDiscard`. | `Discard` |
| **Extend** | Add owned tiles to an existing public meld (`EXTEND_MELD`) — whole meld revalidated. | `Extend` |
| **Replace Joker** | Swap exact represented tile for joker in a meld, then immediately meld the freed joker with two rack tiles (`REPLACE_JOKER`). | `ReplaceJoker` |
| **Win** | `rack == 0` after a legal `DISCARD` or melding action per `docs/rules-decisions.md:6.1` → `RoundComplete` with `WinnerId`, `WinningAction`. | `WinnerId` |

## Naming Conventions (Go)

- Types use `PascalCase` (`TileInstanceId` not `tileInstanceID`) to match `AGENTS.md:259` list.
- Package `internal/rules` will own pure validators/scores; `internal/setup` owns deck/shuffle/deal; `internal/match` owns state machine; `internal/protocol` owns opcodes/messages. Terminology here maps 1:1 to those packages (Day 10+).
- Romanian terms `suita`/`terta`/`Joly` appear in comments but code uses `Run`/`Set`/`Joker` for English consistency.

## Out of Scope (explicitly not in MVP)

`Doubla`, `15 stacks of seven` wall ritual, `exposed-tile bonus`, `detailed multi-round scoring/dealer rotation`, `spectators`, `chat`, `friends`, `bots` — per `docs/rules-decisions.md:7`. If mentioned, prefix with `Deferred:`.

---

*Day 9 deliverable — terminology only. No code. Types follow Day 10: `TileColour`, `TileRank`, `TileInstanceId`, `TileInstance`, `JokerTile`, `Meld`, `RunMeld`, `SetMeld`, `PlayerId`, `Seat`, `GamePhase`, `TurnPhase`.*
