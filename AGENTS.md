You are the lead gameplay/server engineer working in an existing Git repository. Your task is to build a small, maintainable, multiplayer Romanian Tile Rummy game using Nakama, inspired by the social table-game experience of Remi Online, but without copying its branding, artwork, UI, proprietary assets, or implementation.

The game rules must follow Romanian Tile Rummy as documented by Pagat, with all ambiguous or incomplete source excerpts recorded as explicit product decisions before implementation. The initial target is a server-authoritative multiplayer game for 2–4 players.

IMPORTANT WORKING STYLE: “Handmade Hero” incremental development
- Work in very small, working vertical slices.
- Make only one clearly scoped change per day.
- Build on the previous day’s code; do not rewrite broad areas without a demonstrated need.
- Do not dump large amounts of code at once.
- Prefer simple, explicit code over clever abstractions.
- Avoid premature generalization, framework-heavy architecture, and hidden magic.
- Every day must end in a runnable, tested state where possible.
- Preserve a clean commit history: one focused commit per day unless a small follow-up fix is essential.
- Do not proceed to the next day until the current day’s acceptance criteria are met.
- At the end of each day: run relevant checks, `git add`, commit with the prescribed style, push the branch, and report the commit hash plus what changed.

## Product goal

Create a multiplayer Romanian Tile Rummy game with:
- 2, 3, or 4 players.
- Nakama authoritative match handling.
- A simple client-facing state protocol suitable for a future web client.
- Private player racks.
- Public table melds, stock/wall state, discard row, turn state, scoring, and end-of-round state.
- Rules enforcement on the server.
- A development-friendly Docker setup.
- Strong deterministic tests for rules and match flow.

The initial implementation is a functional gameplay core, not a full commercial/social platform. Do NOT implement payments, ads, gambling mechanics, tournaments, friends, avatars, global chat, ranking, or elaborate UI until the core game is complete.

## Source requirements and gameplay scope

Use these rules as the baseline:

### Tile set
- 106 tiles total:
  - 4 colours: red, yellow, blue, black.
  - Numbers 1–13.
  - Two copies of every colour/number combination: 104 numbered tiles.
  - 2 jokers (“Joly”).
- Each tile instance must have a unique immutable ID, even when its colour and number are identical.

### Player count and direction
- Support 2–4 players.
- Play proceeds anticlockwise.
- In a networked game, define a deterministic seat order and explicitly document how “anticlockwise” maps to the next seat.

### Valid melds
- Run (`suita`):
  - At least 3 consecutive values of the same colour.
  - `1-2-3` is valid.
  - `12-13-1` is valid.
  - `13-1-2` is invalid.
  - Aces cannot appear in the middle of a run.
- Set (`terta`):
  - Three or four tiles of the same number, each in a different colour.
- Joker:
  - Can substitute for a needed tile.
  - A meld must have at least twice as many real tiles as jokers.
  - A joker’s represented tile/value must be explicit and immutable once melded, unless a legal replacement operation is performed.
  - A joker at a run end must preserve its declared represented position; do not silently reinterpret it later.

### Initial meld
- A player’s first meld must:
  - Be made only from that player’s rack.
  - Include at least one run.
  - Have a total value of at least 50 points.
- Values:
  - 2–9: 5 points each.
  - 10–13: 10 points each.
  - Ace in `1-2-3...`: 5 points.
  - Ace in `...12-13-1`: 10 points.
  - Aces in a set of three: 25 points each.
  - Joker value equals the represented tile’s value.
- Keep the point calculation isolated in a dedicated rules module with thorough tests.

### Turn flow
- The opening player begins with 15 tiles and starts by discarding one tile.
- Other players start with 14 tiles.
- Normal turns are draw/meld/discard.
- Before initial meld, a player must draw from the wall at the start of their turn.
- After initial meld, permitted actions include:
  - Create additional valid melds.
  - Extend own or other players’ melds.
  - Take the immediately previous discard at the start of a turn instead of drawing from the stock.
  - Take an earlier discard only when the player can immediately meld it with two rack tiles into a valid set/run, then pick up all later discarded tiles into their rack.
  - Replace a joker only under the precise legal conditions and immediately use that joker in a new meld with two tiles from that player’s rack.
- Each turn must end with exactly one discard unless the game-ending rule explicitly prevents that.
- The initial discarded tile is permanently unavailable for pickup and must be represented distinctly in state.

### End condition
- A player wins when they have no tiles left after legal play, subject to the final-discard/closing behavior chosen and documented from the source rules.
- Do not guess at missing rules. If a rule is unclear due to the source excerpt ending early, create a `docs/rules-decisions.md` entry containing:
  1. source rule,
  2. ambiguity,
  3. proposed deterministic behavior,
  4. test cases affected,
  5. a TODO marker for product confirmation.

### Intentionally deferred rules/features
Do not implement these until after the baseline round is playable and tested:
- Exact physical wall construction/dealing ritual involving 15 stacks of seven.
- “Doubla” announcements/exchanges.
- Exposed-tile bonus.
- Detailed multi-round scoring and dealer rotation.
- Spectators.
- Chat, friends, player profiles, tournaments.
- Bots/AI opponents.
- Animation and polished UI.

For the initial version, use a deterministic, server-side shuffle and a simple deal:
- Randomly shuffled tile bag.
- Dealer/opening-player selection deterministic from match setup.
- Opening player receives 15 tiles; all others receive 14.
- Remaining tiles form stock.
- Opening player’s first action is discard.
Document this as a deliberate MVP simplification.

## Technology constraints

Use:
- Nakama authoritative multiplayer matches.
- Docker Compose for local development.
- Nakama database dependency as required by the selected Nakama version.
- TypeScript for Nakama runtime code unless the repository already has a clearly established supported language; if another language is already configured, preserve that choice and explain it.
- A test runner appropriate for the language/runtime.
- JSON messages for client/server match protocol unless the existing repository uses a better established compatible protocol.

Do not add a frontend framework in the first phase unless the repository already has one. The primary deliverable is the authoritative game backend plus developer tools and tests.

## Architecture principles

Use a small, explicit architecture:

```text
/src
  /match
    matchState
    matchHandler
    matchCommands
    matchEvents
    visibility
  /rules
    tile
    meldValidation
    jokerRules
    scoring
    legalMoves
  /setup
    deck
    shuffle
    deal
  /protocol
    opcodes
    messages
    schemas
/tests
/docs
/docker
```

Adapt naming to the existing repository if needed, but preserve the separation of concerns.

Rules:
- The match handler orchestrates state transitions; it must not contain complex meld math.
- Rules modules should be mostly pure and deterministic.
- A player may only see:
  - their own rack,
  - public table melds,
  - public discard row,
  - public turn and phase information,
  - public counts for other players’ racks,
  - public scores/status when available.
- Never send another player’s private rack to a client.
- Server state is authoritative. Never trust client claims that a meld is valid.
- Commands must be validated for:
  - authenticated participant,
  - match phase,
  - active turn,
  - ownership of submitted tile IDs,
  - duplicate/missing IDs,
  - legal rule constraints,
  - atomicity: invalid commands must not partially mutate state.
- Use typed command payloads and schema validation at the protocol boundary.
- All important state changes must emit explicit client events/snapshots.
- State serialization must be stable and intentional. Do not expose internal implementation details.
- Make randomness injectable or seedable for tests.

## Daily execution protocol

For every day:
1. Read the current repository state and prior documentation.
2. State the day’s goal in one sentence.
3. Describe the smallest implementation plan before writing code.
4. Implement only that day’s scope.
5. Add or update tests for new behavior.
6. Run formatting, type checks, tests, and Docker/integration checks relevant to the change.
7. Update documentation/README when the developer workflow or protocol changes.
8. Show a concise summary:
   - files changed,
   - behavior added,
   - tests run and results,
   - known limitations,
   - next day’s dependency.
9. Execute:
   - `git status`
   - `git add <specific files>`
   - `git commit -m "<type>: <focused message>"`
   - `git push`
10. Stop. Do not begin the next day unless explicitly asked, or unless instructed to execute the whole plan one day at a time.

Commit style:
- `chore: ...`
- `docs: ...`
- `feat: ...`
- `test: ...`
- `fix: ...`
- `refactor: ...`

Never use a vague commit like `update`, `changes`, or `wip`.

## Required implementation plan

Follow this plan in order. If an existing repository changes the exact setup, preserve the intent and document deviations.

### Day 1 — Repository, Nakama, Docker, and developer baseline
Goal: Create a reproducible local Nakama development environment.

Implement:
- Inspect the repository and identify current language/tooling.
- Add Docker Compose for Nakama and its required database.
- Add Nakama configuration appropriate for local development.
- Create the runtime module skeleton.
- Add `.env.example` with non-secret local defaults.
- Add `README.md` instructions:
  - prerequisites,
  - start/stop commands,
  - Nakama console URL,
  - how to inspect logs,
  - how to run tests.
- Add basic formatting/lint/type-check/test scripts.
- Add a health/smoke verification that confirms Nakama starts.

Acceptance criteria:
- `docker compose up --build` starts Nakama and database successfully.
- A developer can reach the Nakama console.
- The project has a documented local workflow.
- No game behavior is implemented yet.

Commit:
`chore: bootstrap local nakama development environment`

### Day 2 — Rules specification and domain model
Goal: Establish a shared, testable vocabulary for tiles and game decisions.

Implement:
- `docs/rules-decisions.md`.
- Capture the source-derived rules, MVP simplifications, deferred rules, and unresolved ambiguities.
- Define domain types:
  - `TileColor`,
  - `TileRank`,
  - `Tile`,
  - `JokerTile`,
  - `TileInstanceId`,
  - `Meld`,
  - `RunMeld`,
  - `SetMeld`,
  - `PlayerId`,
  - `Seat`,
  - `GamePhase`,
  - `TurnPhase`.
- Define explicit state invariants in documentation.
- Add unit tests for type-level/domain construction where meaningful.

Acceptance criteria:
- The team can tell exactly what is in and out of the MVP.
- There is no implicit “joker means anything” representation.
- Ambiguous rules have documented decisions/TODOs.

Commit:
`docs: define rummy rules and domain model`

### Day 3 — Deck creation and deterministic shuffle
Goal: Create the complete tile set correctly and make it testable.

Implement:
- Deck factory for 106 unique tile instances.
- Two copies of every numbered colour/rank tile.
- Two unique jokers.
- Seeded/injectable shuffle function.
- Tests for:
  - total count 106,
  - 104 numbered tiles,
  - 2 jokers,
  - exactly two of each colour/rank combination,
  - all instance IDs unique,
  - deterministic output for a fixed seed,
  - no lost/duplicated tiles after shuffle.

Acceptance criteria:
- Deck correctness is proven by tests.
- Randomness is controllable in tests.

Commit:
`feat: add deterministic romanian tile deck`

### Day 4 — MVP dealing and initial round state
Goal: Produce a legal initial game state for 2–4 players.

Implement:
- Deterministic seat assignment.
- MVP dealer/opening-player selection.
- Deal logic:
  - opening player: 15 tiles,
  - others: 14 tiles,
  - remainder: stock.
- Initial game state:
  - player racks,
  - stock,
  - empty public melds,
  - discard row with a dedicated “opening discard unavailable” marker prepared for use,
  - current player,
  - opening-discard phase.
- Invariant checker ensuring every tile is in exactly one location.

Tests:
- 2, 3, and 4 player deals.
- Correct rack sizes.
- Correct stock sizes.
- No tile duplication/loss.
- Correct opening player and phase.

Acceptance criteria:
- A round can be initialized deterministically.
- The tile conservation invariant passes after setup.

Commit:
`feat: create initial rummy round state`

### Day 5 — Nakama authoritative match skeleton
Goal: Run and join a real authoritative Nakama match.

Implement:
- Nakama runtime initialization.
- Match create/join/leave/terminate lifecycle.
- Lobby/waiting state for 2–4 participants.
- Start condition for local development: allow the match creator to issue a start command once at least 2 players join.
- Stable protocol opcodes and typed message envelope.
- Public match snapshot that reveals no private racks.
- Player-specific snapshot sent to each participant containing only their own rack.

Tests:
- Unit-test state transitions where possible.
- Add a simple integration/smoke script or documented manual test using Nakama console/client tooling.

Acceptance criteria:
- Two players can join a match.
- A match can start and each player receives an appropriate state view.
- Private rack data is not broadcast to other players.

Commit:
`feat: add authoritative nakama match lifecycle`

### Day 6 — Protocol validation and safe command dispatch
Goal: Establish a robust foundation for all game actions.

Implement:
- Command envelope with version, opcode, request ID, and payload.
- Runtime schema validation.
- Standard error response format:
  - error code,
  - human-readable message,
  - optional field/details,
  - request ID.
- Dispatcher that checks membership and routes commands.
- Reject malformed JSON, unknown opcode, invalid payload, non-member actions, and out-of-phase actions.
- Idempotency strategy for request IDs where needed; document its scope.

Acceptance criteria:
- Invalid client input cannot crash or corrupt a match.
- Errors are predictable and test-covered.
- No gameplay mutation occurs yet except safe command handling.

Commit:
`feat: validate and dispatch match commands`

### Day 7 — Opening discard
Goal: Complete the first mandatory move of the round.

Implement:
- `DISCARD` command.
- Opening player must discard exactly one tile from their 15-tile rack.
- Mark this discard as the opening discard and permanently unavailable for pickup.
- Advance turn anticlockwise.
- Change phase to normal turn start.
- Broadcast public discard update and private rack update.

Tests:
- Only current opening player can discard.
- Tile must belong to their rack.
- Opening player’s rack becomes 14.
- The opening discard cannot later be selected by pickup actions.
- Tile conservation holds.

Acceptance criteria:
- A newly started match progresses from setup to the first normal turn.

Commit:
`feat: implement opening discard turn`

### Day 8 — Draw from stock
Goal: Implement the default start-of-turn action.

Implement:
- Normal turn phases: `MUST_DRAW`, `MELD_OR_DISCARD`, `TURN_COMPLETE`.
- `DRAW_STOCK` command.
- Only current player in `MUST_DRAW` can draw.
- Add one tile to player rack, remove it from stock, transition to `MELD_OR_DISCARD`.
- Handle empty-stock behavior as a documented MVP decision; do not improvise silently.

Tests:
- Turn/phase validation.
- Stock count and rack count changes.
- Tile conservation.
- Empty-stock result behavior.

Acceptance criteria:
- A normal player can start a turn by drawing one stock tile.

Commit:
`feat: implement stock draw action`

### Day 9 — Basic discard and turn rotation
Goal: Make ordinary non-melding turns playable.

Implement:
- Extend `DISCARD` for normal turns.
- A player in `MELD_OR_DISCARD` can discard exactly one owned tile.
- Append discard to public row.
- Advance active player anticlockwise and reset to `MUST_DRAW`.
- Keep discard history ordered and distinguish the opening discard.

Tests:
- Reject discard before draw.
- Reject non-current player.
- Reject foreign/nonexistent tile.
- Verify turn order for 2, 3, 4 players.
- Verify row ordering and conservation.

Acceptance criteria:
- Players can loop through draw-from-stock and discard turns.

Commit:
`feat: add normal discard and turn rotation`

### Day 10 — Pure validation for sets
Goal: Implement reliable validation of `terta` melds.

Implement:
- Pure set validation:
  - exactly 3 or 4 tiles after joker substitution,
  - same rank,
  - distinct colours among real/replaced tiles,
  - legal joker assignments,
  - real-to-joker ratio rule.
- Return structured validation results, not only booleans.
- Define canonical representation of joker substitution in a meld.

Tests:
- Valid 3- and 4-colour sets.
- Duplicate colour rejection.
- Rank mismatch rejection.
- Legal/illegal joker use.
- Two jokers in one set rejection.
- Ratio-rule rejection.

Acceptance criteria:
- Set validation is independently tested and reusable by match code.

Commit:
`feat: validate rummy set melds`

### Day 11 — Pure validation for runs
Goal: Implement reliable validation of `suita` melds, including ace edge cases.

Implement:
- Pure run validation:
  - same colour,
  - length at least 3,
  - legal consecutive sequence,
  - legal `1-2-3`,
  - legal `12-13-1`,
  - reject `13-1-2`,
  - explicit joker represented rank/colour,
  - real-to-joker ratio rule,
  - immutable declared joker interpretation.
- Prefer a clear enumerative/canonical algorithm over a compact but opaque one.

Tests:
- Standard runs.
- Low-ace runs.
- high-ace runs.
- invalid ace-middle runs.
- joker-gap runs.
- invalid reinterpretation scenarios.
- multiple-joker ratio cases.

Acceptance criteria:
- All important run and ace cases are proven by tests.

Commit:
`feat: validate rummy run melds`

### Day 12 — Meld scoring and initial-meld eligibility
Goal: Calculate 50-point initial meld eligibility correctly.

Implement:
- Pure scoring functions for tile values in context.
- Score sets and runs, including joker represented values.
- Evaluate a proposed first meld batch:
  - all tiles from player rack,
  - at least one run,
  - total score >= 50,
  - every submitted meld independently valid,
  - no duplicate tile use across melds.

Tests:
- Numeric tile scoring.
- Low/high ace scoring.
- Ace sets.
- Joker scoring.
- exactly 50 accepted, 49 rejected.
- no-run rejection.
- duplicated tile across proposed melds rejection.

Acceptance criteria:
- Initial meld eligibility is a pure, tested rules operation.

Commit:
`feat: add initial meld scoring rules`

### Day 13 — Submit initial melds to the table
Goal: Let a player legally open their table play.

Implement:
- `MELD_INITIAL` command that submits a batch of melds.
- Require active player and `MELD_OR_DISCARD`.
- Validate using Day 12 rules.
- Atomically remove rack tiles and add public melds.
- Mark player as “opened”.
- Keep player in `MELD_OR_DISCARD`; they must still discard to end the turn.
- Emit state updates with table meld IDs and canonical joker assignments.

Tests:
- Successful initial meld.
- Invalid initial meld leaves state unchanged.
- Player cannot initial-meld twice.
- Player still must discard after meld.
- Other players see public meld but not rack details.

Acceptance criteria:
- A player can draw, open with >=50 points including a run, then discard.

Commit:
`feat: allow initial table melds`

### Day 14 — Additional melds after opening
Goal: Allow opened players to create further melds.

Implement:
- `MELD_NEW` command for players who have already opened.
- Validate one or more independent sets/runs from rack.
- No minimum total score, but all normal meld validity rules apply.
- Apply atomically.
- Preserve immutable table meld IDs and joker assignments.

Tests:
- Opened player creates valid additional set/run.
- Unopened player is rejected.
- Invalid batch is atomic.
- Duplicate tile ID rejected.
- Meld IDs remain stable.

Acceptance criteria:
- An opened player can add legal new melds before discarding.

Commit:
`feat: allow additional melds after opening`

### Day 15 — Extend existing public melds
Goal: Support adding tiles to own or other players’ table melds.

Implement:
- `EXTEND_MELD` command.
- Allow only opened players.
- Client submits target meld ID and rack tile IDs plus explicit intended resulting meld representation if needed.
- Validate the entire resulting meld, not just the added tile.
- Ensure all added tiles come from actor’s rack.
- Preserve tile ownership history separately if desired, but treat table melds as public game objects.
- Never permit a tile to belong to two melds.

Tests:
- Extend run at either legal end.
- Extend set to four colours.
- Extend another player’s meld.
- Invalid extension does not mutate.
- Joker interpretation cannot be silently changed.
- Tile conservation and uniqueness.

Acceptance criteria:
- Public melds can be legally extended with server validation.

Commit:
`feat: support extending public melds`

### Day 16 — Pick up the immediately previous discard
Goal: Implement the simplest discard pickup after a player has opened.

Implement:
- `DRAW_PREVIOUS_DISCARD` command.
- Allowed only:
  - active player,
  - `MUST_DRAW`,
  - player has already opened,
  - previous discard exists and is not the opening discard.
- Move only the latest eligible discard into the player’s rack.
- Transition to `MELD_OR_DISCARD`.

Tests:
- Unopened player rejection.
- Opening discard rejection.
- Only latest discard can be taken by this command.
- Rack/discard counts and conservation.
- Cannot draw twice in same turn.

Acceptance criteria:
- Opened players can choose stock or latest-discard draw at turn start.

Commit:
`feat: allow pickup of previous discard`

### Day 17 — Pick up an earlier discard with required immediate meld
Goal: Implement the controlled multi-discard pickup rule.

Implement:
- `PICKUP_DISCARD_FOR_MELD` command.
- Allowed only for opened active player in `MUST_DRAW`.
- Client identifies a non-opening discard and submits:
  - exactly two rack tiles,
  - a new meld using those two tiles plus selected discard,
  - any resulting picked-up later discards are added to rack.
- Validate:
  - selected discard is eligible,
  - exactly two supplied tiles are from actor’s rack,
  - the three tiles form a valid meld,
  - selected discard is immediately placed in that meld,
  - all later discards transfer to rack in correct order,
  - mutation is atomic.
- Transition to `MELD_OR_DISCARD`.

Tests:
- Valid three-tile set pickup.
- Valid three-tile run pickup.
- Invalid combination rejection.
- Opening discard cannot be selected.
- Correct collection of all later discards.
- No partial mutation.

Acceptance criteria:
- The special discard-row pickup works deterministically and safely.

Commit:
`feat: implement discard row pickup meld`

### Day 18 — Joker replacement rules
Goal: Implement legal joker recovery without loopholes.

Implement:
- `REPLACE_JOKER` command(s), with an explicit, understandable payload.
- Cover:
  1. Replacing a joker in a run using the exact represented tile, then immediately using recovered joker with two tiles from the actor’s rack to make a new valid meld.
  2. Replacing a joker in a three-tile set by supplying both missing colours of that rank, then immediately using recovered joker with two actor-rack tiles in a new valid meld.
- Require actor has opened and is in `MELD_OR_DISCARD`.
- Validate every affected table meld and the new joker meld atomically.
- Document any source-rule ambiguity not fully available in the supplied excerpt.

Tests:
- Valid run joker replacement.
- Wrong represented tile rejection.
- Valid set joker replacement.
- Missing required colours rejection.
- New joker meld must use two rack tiles.
- Atomic rollback on failure.

Acceptance criteria:
- Jokers cannot be extracted or reassigned illegally.

Commit:
`feat: enforce joker replacement rules`

### Day 19 — Round completion and winner state
Goal: End a round deterministically.

Implement:
- Define/document MVP win condition precisely.
- Detect legal empty-rack state after a player action.
- Transition match to `ROUND_COMPLETE`.
- Record winner and round summary:
  - winner ID,
  - final public melds,
  - remaining rack counts,
  - action that ended the round.
- Prevent further gameplay commands after completion.
- Emit a final public/private result snapshot.

Tests:
- Win after valid meld/discard flow according to documented decision.
- No gameplay after round complete.
- Winner is correctly recorded.
- Tile conservation remains valid at completion.

Acceptance criteria:
- A complete playable round has a clear, authoritative winner.

Commit:
`feat: complete rummy rounds with winner state`

### Day 20 — State visibility, reconnection, and snapshot hardening
Goal: Ensure clients can safely recover state without hidden-information leaks.

Implement:
- Centralized public/private view projection.
- Reconnection handling:
  - a returning player receives current private rack plus public state,
  - other players never receive it.
- Version state snapshots/messages.
- Add a redaction test suite that serializes views for several players and proves no foreign rack tile IDs are present.
- Document client synchronization expectations.

Acceptance criteria:
- Reconnect can restore a player’s game state.
- Snapshot tests prove hidden tiles do not leak.

Commit:
`feat: harden match snapshots and reconnection`

### Day 21 — Test harness and deterministic end-to-end simulation
Goal: Test a complete game flow without manual clicking.

Implement:
- Deterministic test harness for creating match state with fixed deck/order.
- Simulate:
  - match start,
  - opening discard,
  - draw/discard turns,
  - initial meld,
  - extension,
  - previous-discard pickup,
  - special discard-row pickup,
  - joker replacement,
  - round completion.
- Keep the test readable: use named tiles/builders, not raw opaque IDs everywhere.
- Add invariant assertion after every action.

Acceptance criteria:
- A deterministic end-to-end test demonstrates major game mechanics.
- Failures identify the transition/rule that broke.

Commit:
`test: add deterministic rummy match simulation`

### Day 22 — Developer tooling and operator documentation
Goal: Make the project easy for another engineer to run and modify.

Implement:
- Improve README with:
  - architecture overview,
  - local commands,
  - test commands,
  - debugging a match,
  - protocol overview,
  - how to add a new command safely.
- Add `docs/protocol.md`.
- Add `docs/state-machine.md` with turn phases and allowed commands.
- Add `docs/testing.md`.
- Add a small local script/tool that can create a match and send sample commands, if feasible without adding excessive complexity.

Acceptance criteria:
- A new developer can start the environment, understand the state machine, and run tests from documentation alone.

Commit:
`docs: document rummy development and protocol workflow`

### Day 23 — Refactoring only where proven necessary
Goal: Improve clarity without changing behavior.

Implement:
- Review code for duplication and unclear naming.
- Refactor only after tests are in place.
- Keep rules pure and match orchestration thin.
- Add comments only for non-obvious Romanian rummy rules or deliberate decisions.
- Do not add features.

Acceptance criteria:
- All tests pass unchanged in intent.
- Diff is focused on readability/maintainability.
- No protocol-breaking changes.

Commit:
`refactor: clarify rummy match and rules modules`

### Day 24 — Optional minimal playable client adapter
Goal: Provide the smallest possible client integration reference.

Only begin this day if the backend core is stable and all prior acceptance criteria pass.

Implement one of:
- a minimal web test client, or
- a CLI client.

It must:
- authenticate/create users in local development,
- create/join a match,
- show only the local player’s rack,
- show public melds/discards/current turn,
- send draw, discard, and basic meld commands,
- display server errors.

Do not spend time on visual polish. The purpose is protocol validation and manual testing.

Acceptance criteria:
- Two local users can play a small manual test flow against the authoritative server.
- No private data leak is visible in network payloads or UI.

Commit:
`feat: add minimal rummy match test client`

## Quality bar

Before declaring any day complete:
- Code is formatted.
- Static checks pass.
- New logic has tests.
- Existing tests remain green.
- Invalid client input is handled safely.
- Tile conservation invariant is checked in state-changing tests.
- No other player’s private tiles are included in public payloads.
- Documentation is updated if behavior/protocol changed.
- The commit is focused and pushed.

## Response format when executing a day

Use exactly this structure:

# Day N — <title>

## Goal
One sentence.

## Plan
1. ...
2. ...

## Implementation
- ...
- ...

## Tests and verification
- `command` — result
- ...

## Acceptance criteria
- [x] ...
- [ ] ... if blocked, explain why

## Git
- Branch: `...`
- Commit: `...`
- Push: successful/blocked with reason

## Notes for Day N+1
- ...

## First action

Start with Day 1 only. Inspect the repository before making assumptions. If files already exist, adapt rather than replace them. Do not implement Day 2 or later until Day 1 is complete, committed, and pushed.
```

The prompt deliberately treats Remi Online as inspiration for an online/social rummy-table experience while using Romanian Tile Rummy rules as the gameplay authority: the source describes 106 tiles, 2–4 players, anticlockwise turns, runs/sets, jokers, a 50-point opening meld with at least one run, discard pickup, and joker-replacement mechanics.[[1](https://www.pagat.com/rummy/romtile.html)] Remi Online also distinguishes “Remi Etalat” from “Remi Pe Tabla”; this plan targets the tile-rummy/meld-on-table style rather than attempting to reproduce both game variants at once.[[2](https://www.remi-online.ro)]





# now make the implementation plan roadmap from day 1 to nnn.



# Romanian Tile Rummy in Nakama — Incremental Implementation Roadmap

This roadmap follows a **Handmade Hero-style development process**: every day produces a small, testable, committed improvement. Do not skip ahead, do not introduce large abstractions early, and do not add UI/social features until the authoritative game core is reliable.

**Target:** server-authoritative Romanian Tile Rummy for 2–4 players using Nakama and Docker.

---

## Phase 1 — Foundation and local development

| Day | Focus | Deliverable |
|---:|---|---|
| 1 | Repository and environment audit | Review existing project, language, tooling, branching, CI, and deployment assumptions. Record findings in `docs/project-baseline.md`. |
| 2 | Docker Compose | Add local Docker Compose setup for Nakama and PostgreSQL. Verify containers start and persist data correctly. |
| 3 | Nakama runtime skeleton | Add Nakama runtime module structure, build process, and a minimal `InitModule` log message. |
| 4 | Developer scripts | Add scripts for start, stop, logs, clean/reset database, build, test, lint, type-check, and format. |
| 5 | Local setup documentation | Write clear README instructions for prerequisites, environment variables, Docker startup, Nakama Console access, and test execution. |
| 6 | CI baseline | Add CI pipeline for formatting, linting, type checking, unit tests, and build verification. |
| 7 | Smoke test | Add a repeatable smoke test that verifies Nakama starts, runtime code loads, and database connectivity works. |

**Milestone:** A new developer can clone the repository, run one command, start Nakama locally, and verify the backend is healthy.

---

## Phase 2 — Rules specification and core domain

| Day | Focus | Deliverable |
|---:|---|---|
| 8 | Rules source documentation | Create `docs/rules-decisions.md` containing the Romanian Tile Rummy rules, MVP decisions, assumptions, ambiguities, and deferred features. |
| 9 | Domain terminology | Define shared names and terminology: tile, rack, stock, discard row, meld, run, set, joker, opening meld, seat, turn, round. |
| 10 | Tile domain model | Implement types for colours, ranks, tile IDs, numbered tiles, joker tiles, and unique tile instances. |
| 11 | Player and seat model | Implement player identity, deterministic seating, player state, and anticlockwise turn-order helpers. |
| 12 | Game state model | Create initial server-side state structures for round state, stock, racks, discards, public melds, current player, and turn phase. |
| 13 | State invariants | Implement an invariant checker ensuring every tile exists in exactly one location: rack, stock, discard row, table meld, or reserved setup location. |
| 14 | Domain tests | Add tests for tile identity, seats, turn direction, base state construction, and invariant failures. |

**Milestone:** The game has a clear, documented vocabulary and a safe internal state model before multiplayer behavior begins.

---

## Phase 3 — Deck, randomness, and round setup

| Day | Focus | Deliverable |
|---:|---|---|
| 15 | Full deck factory | Create the 106-tile Romanian Tile Rummy deck: 104 numbered tiles plus 2 jokers. |
| 16 | Deck correctness tests | Verify exactly two copies of every colour/rank tile, two jokers, 106 total tiles, and unique instance IDs. |
| 17 | Deterministic random source | Add injectable/seedable randomness for reproducible shuffles and tests. |
| 18 | Shuffle implementation | Implement Fisher–Yates or equivalent clear shuffle logic with deterministic test support. |
| 19 | Deal logic | Implement MVP dealing: opening player gets 15 tiles, all other players get 14, remainder becomes stock. |
| 20 | Round initialization | Build `createRoundState()` with seats, dealer/opening player, racks, stock, empty table, empty discard row, and opening turn phase. |
| 21 | Setup invariants | Test dealing for 2, 3, and 4 players; verify tile conservation and expected stock counts. |

**Milestone:** A complete legal round can be initialized deterministically without Nakama.

---

## Phase 4 — Nakama match foundation

| Day | Focus | Deliverable |
|---:|---|---|
| 22 | Runtime module initialization | Register authoritative match handlers in Nakama. |
| 23 | Match lifecycle skeleton | Implement match init, join attempt, join, leave, loop, signal, and terminate callbacks. |
| 24 | Waiting room state | Add lobby state for 2–4 players with seat allocation and player-ready information. |
| 25 | Match start command | Allow host/creator to start a match when at least 2 players are present. |
| 26 | Protocol opcodes | Define stable client/server opcodes and message envelope versioning. |
| 27 | Command parser | Parse inbound JSON safely and reject malformed payloads without crashing the match. |
| 28 | Command schema validation | Add payload validation for every currently available command. |
| 29 | Standard error protocol | Add consistent error responses: code, message, request ID, and optional structured details. |
| 30 | Basic match snapshots | Send public game state and player-specific private state snapshots when a round starts. |
| 31 | Hidden-information test | Verify public snapshots and opponent snapshots never reveal another player’s rack tiles. |

**Milestone:** Two or more players can create, join, start, and receive a secure authoritative match state.

---

## Phase 5 — Turn state machine and basic actions

| Day | Focus | Deliverable |
|---:|---|---|
| 32 | Turn state machine | Define explicit phases: waiting, opening discard, must draw, may meld/discard, round complete. |
| 33 | Active-player validation | Reject actions from non-active players. |
| 34 | Phase validation | Reject actions that are not legal in the current turn phase. |
| 35 | Opening discard command | Implement the opening player’s mandatory first discard from 15 tiles. |
| 36 | Opening discard protection | Mark the first discard as permanently unavailable for normal discard pickup. |
| 37 | Turn advance | Advance active seat anticlockwise after opening discard. |
| 38 | Draw from stock | Implement `DRAW_STOCK` for active player during `MUST_DRAW`. |
| 39 | Normal discard | Implement `DISCARD` after drawing or melding. |
| 40 | Discard row ordering | Preserve public discard history in exact chronological order. |
| 41 | Turn-loop tests | Test a complete sequence: opening discard → draw stock → discard → next player. |
| 42 | Empty-stock decision | Document and implement deterministic MVP behavior for an exhausted stock. |

**Milestone:** Players can play a valid draw-and-discard loop with a server-authoritative turn order.

---

## Phase 6 — Meld rules: sets, runs, and jokers

| Day | Focus | Deliverable |
|---:|---|---|
| 43 | Meld representation | Define canonical representation for table melds, including stable meld IDs and joker substitutions. |
| 44 | Basic set validation | Validate 3- and 4-tile sets of equal rank with distinct colours. |
| 45 | Set validation errors | Return structured reasons for invalid sets: rank mismatch, duplicate colour, invalid size, duplicate tile. |
| 46 | Set joker support | Support jokers in sets with explicit represented rank/colour. |
| 47 | Set joker ratio rule | Enforce at least twice as many real tiles as jokers. |
| 48 | Basic run validation | Validate same-colour consecutive runs of length 3 or greater. |
| 49 | Low-ace runs | Support `1-2-3` and longer low-ace runs. |
| 50 | High-ace runs | Support `12-13-1` and longer high-ace runs according to the selected canonical model. |
| 51 | Invalid ace-middle runs | Explicitly reject sequences such as `13-1-2`. |
| 52 | Run joker support | Support jokers in runs with explicit represented colour and rank. |
| 53 | Run joker ratio rule | Enforce the real-tile-to-joker ratio for runs. |
| 54 | Immutable joker mapping | Prevent a tabled joker from silently changing represented value later. |
| 55 | Meld test matrix | Add a comprehensive test matrix for valid and invalid sets, runs, ace cases, and joker cases. |

**Milestone:** The game can independently and reliably determine whether any proposed meld is legal.

---

## Phase 7 — Opening meld and scoring

| Day | Focus | Deliverable |
|---:|---|---|
| 56 | Tile scoring model | Implement value rules for ranks 2–9, 10–13, aces, and jokers. |
| 57 | Run scoring | Score low-ace and high-ace runs correctly. |
| 58 | Set scoring | Score sets, including special ace-set value rules. |
| 59 | Joker scoring | Score a joker according to its declared represented tile. |
| 60 | Opening meld batch model | Define a batch of proposed initial melds from one player rack. |
| 61 | Opening meld validation | Require all opening meld tiles to come from the player’s rack. |
| 62 | Opening meld minimum score | Require total value of at least 50 points. |
| 63 | Opening meld run requirement | Require at least one valid run in the opening batch. |
| 64 | Duplicate tile prevention | Reject any tile used more than once across the batch. |
| 65 | Initial meld command | Add `MELD_INITIAL` to the Nakama command handler. |
| 66 | Atomic initial meld mutation | Remove rack tiles and add public melds only if every validation passes. |
| 67 | Opening meld tests | Cover 49 vs. 50 points, no-run rejection, joker scoring, duplicate use, and atomic rollback. |

**Milestone:** A player can legally “open” with at least 50 points and at least one run.

---

## Phase 8 — Post-opening table play

| Day | Focus | Deliverable |
|---:|---|---|
| 68 | Opened-player flag | Track whether each player has completed their initial meld. |
| 69 | Additional new melds | Add `MELD_NEW` for opened players to create new valid table melds. |
| 70 | Batch new melds | Permit multiple new melds in one command while preserving atomicity. |
| 71 | Table meld extension model | Define how a rack tile is added to an existing public meld. |
| 72 | Extend set meld | Allow legal extension of a three-tile set to a four-colour set. |
| 73 | Extend run meld | Allow legal extension at a valid run endpoint. |
| 74 | Extend any player’s meld | Permit opened players to extend public melds regardless of original owner. |
| 75 | Extension command | Add `EXTEND_MELD` protocol handling and server validation. |
| 76 | Extension rollback tests | Ensure invalid extensions never partially mutate state. |
| 77 | Table state projection | Improve public state messages for meld IDs, tiles, joker assignments, and updates. |

**Milestone:** Opened players can build on their own or others’ table melds safely.

---

## Phase 9 — Discard pickup rules

| Day | Focus | Deliverable |
|---:|---|---|
| 78 | Latest discard pickup rule | Document exact MVP behavior for taking the immediately preceding discard. |
| 79 | Previous-discard draw command | Implement `DRAW_PREVIOUS_DISCARD`. |
| 80 | Previous-discard validation | Require opened player, active turn, `MUST_DRAW`, non-opening discard, and one draw per turn. |
| 81 | Previous-discard tests | Test valid pickup, opening-discard rejection, unopened-player rejection, and turn-phase rejection. |
| 82 | Earlier discard pickup model | Define payload for selecting an earlier discard plus exactly two rack tiles. |
| 83 | Immediate pickup meld validation | Require selected discard plus exactly two rack tiles to form a legal meld immediately. |
| 84 | Later discard collection | Move all discards after the selected discard into the player’s rack, preserving order. |
| 85 | Earlier discard pickup command | Implement `PICKUP_DISCARD_FOR_MELD`. |
| 86 | Atomic discard-row pickup | Ensure selected discard, generated meld, rack additions, and discard-row removal occur atomically. |
| 87 | Discard pickup test suite | Cover valid set/run pickup, invalid pickup, opening discard exclusion, and conservation checks. |

**Milestone:** The special Romanian discard-row pickup mechanic works safely and deterministically.

---

## Phase 10 — Joker replacement mechanics

| Day | Focus | Deliverable |
|---:|---|---|
| 88 | Joker replacement rules document | Document exact supported replacement cases and unresolved source ambiguities. |
| 89 | Run joker replacement validation | Allow replacement only with the exact tile represented by a joker in a run. |
| 90 | New joker meld requirement | Require the recovered joker to immediately form a new valid meld with two tiles from the actor’s rack. |
| 91 | Set joker replacement validation | For a joker in a three-tile set, require the needed missing colours/rank as defined in the rules decision. |
| 92 | Joker replacement command | Implement `REPLACE_JOKER`. |
| 93 | Atomic replacement mutation | Apply table replacement and new joker meld as one all-or-nothing state transition. |
| 94 | Joker replacement tests | Test valid run replacement, wrong-tile rejection, valid set replacement, missing-tile rejection, and failed atomic rollback. |

**Milestone:** Jokers cannot be illegally removed, reassigned, or exploited.

---

## Phase 11 — Win conditions, round completion, and results

| Day | Focus | Deliverable |
|---:|---|---|
| 95 | Closing rule decision | Finalize/document MVP interpretation of winning with no remaining tiles and final-discard behavior. |
| 96 | Win detection | Detect when a player has legally emptied their rack. |
| 97 | Round-complete state | Add `ROUND_COMPLETE`, winner ID, ending action, timestamps, and summary information. |
| 98 | Post-game command blocking | Reject gameplay commands after a round ends. |
| 99 | Final state broadcast | Send final public state and each player’s final private view. |
| 100 | Round completion tests | Test win detection, post-game rejection, conservation, and final event payloads. |
| 101 | Optional dead-round behavior | If stock exhaustion behavior requires it, implement and test draw/round-ending state according to documented MVP decision. |

**Milestone:** A full round can end authoritatively and report a winner.

---

## Phase 12 — Security, reconnection, and reliability

| Day | Focus | Deliverable |
|---:|---|---|
| 102 | Central view projection | Consolidate public and per-player-private snapshot generation in one module. |
| 103 | Snapshot redaction tests | Serialize views for all players and prove foreign rack IDs are absent. |
| 104 | Reconnection handling | Send the current private rack and public state to a returning player. |
| 105 | Disconnect handling | Define/document behavior when a player disconnects during a live round. |
| 106 | Grace-period support | Add optional reconnect grace period before ending or pausing a match. |
| 107 | Command request IDs | Add request IDs and consistent response correlation. |
| 108 | Idempotency behavior | Add limited idempotency for repeated client requests where appropriate. |
| 109 | Abuse protection review | Add payload-size limits, rate-safe command handling, and rejection of impossible tile IDs. |
| 110 | Runtime error hardening | Ensure malformed commands and unexpected internal errors are logged without leaking private state. |

**Milestone:** The match is safe to reconnect to and does not leak hidden information.

---

## Phase 13 — Test harness and automated simulations

| Day | Focus | Deliverable |
|---:|---|---|
| 111 | Test tile builders | Add readable test helpers for named tiles, racks, stock, melds, and discard rows. |
| 112 | Fixed-deck scenario tools | Allow tests to define exact deck sequences and initial hands. |
| 113 | Action simulation helpers | Add helpers for executing commands and asserting state transitions. |
| 114 | Invariant-after-action helper | Run tile conservation and state consistency checks after every simulated action. |
| 115 | Basic full-round simulation | Simulate opening discard, draw, discard, opening meld, and round completion. |
| 116 | Extension simulation | Simulate creating and extending public melds. |
| 117 | Discard pickup simulation | Simulate latest and earlier discard pickup paths. |
| 118 | Joker replacement simulation | Simulate legal and illegal joker replacement scenarios. |
| 119 | Fuzz/property tests | Add constrained random action tests for invariants and non-crashing validation. |
| 120 | Regression test documentation | Document how to add a regression scenario when a rules bug is found. |

**Milestone:** Game rules and state transitions are covered by deterministic and repeatable end-to-end simulations.

---

## Phase 14 — Developer experience and documentation

| Day | Focus | Deliverable |
|---:|---|---|
| 121 | Architecture documentation | Create `docs/architecture.md` describing Nakama lifecycle, rules modules, state ownership, and visibility model. |
| 122 | State machine documentation | Create `docs/state-machine.md` with phases, legal commands, and transition diagrams. |
| 123 | Protocol documentation | Create `docs/protocol.md` listing opcodes, payloads, snapshots, events, and errors. |
| 124 | Rules documentation cleanup | Review `docs/rules-decisions.md`; separate confirmed rules, MVP choices, and unresolved TODOs. |
| 125 | Testing documentation | Create `docs/testing.md` for unit, integration, simulation, and Docker tests. |
| 126 | Operations documentation | Add local troubleshooting: logs, database reset, runtime rebuild, match debugging, and common errors. |
| 127 | Developer command tool | Optionally add a CLI/script to create matches and send sample commands to Nakama. |

**Milestone:** Another engineer can run, understand, test, and modify the game without oral knowledge.

---

## Phase 15 — Refactoring and stabilization

| Day | Focus | Deliverable |
|---:|---|---|
| 128 | Naming review | Rename unclear variables/functions/types without changing behavior. |
| 129 | Rules module cleanup | Remove duplication across set, run, scoring, and joker-validation modules. |
| 130 | Match handler cleanup | Keep Nakama callbacks thin and move remaining pure logic into testable modules. |
| 131 | Protocol compatibility review | Confirm versioning and client error behavior are stable. |
| 132 | Performance review | Measure state size, command cost, and worst-case validation paths. |
| 133 | Logging review | Add structured debug logs for match ID, player ID, command, rejection reason, and transition result. |
| 134 | Final backend regression pass | Run all tests, Docker smoke checks, simulation tests, and hidden-information checks. |
| 135 | Release candidate tag | Create a documented backend MVP release candidate. |

**Milestone:** The backend is stable enough for a minimal client and manual multiplayer testing.

---

# Optional Phase 16 — Minimal playable test client

Only start after Day 135 is stable.

| Day | Focus | Deliverable |
|---:|---|---|
| 136 | Client technology choice | Choose minimal CLI or web client based on repository context. |
| 137 | Local authentication | Add simple local Nakama account/session flow. |
| 138 | Match browser/create/join | Create and join a local rummy match. |
| 139 | Public game view | Show player list, turn, stock count, discard row, and public melds. |
| 140 | Private rack view | Show only the local player’s rack. |
| 141 | Draw/discard controls | Add controls for drawing stock, picking latest discard, and discarding. |
| 142 | Initial meld controls | Add basic tile selection and submit initial meld batch. |
| 143 | Table extension controls | Add basic support for extending melds. |
| 144 | Error display | Render server validation errors clearly. |
| 145 | Manual two-player test | Validate a real two-player session against local Nakama. |
| 146 | Client protocol hardening | Confirm no hidden opponent data is present in browser/CLI logs or network payloads. |

**Milestone:** Two local users can manually play a basic authoritative game flow.

---

# Optional Phase 17 — Future product features

These are explicitly outside the core MVP and should only begin after the gameplay backend is reliable.

| Day range | Feature area | Examples |
|---:|---|---|
| 147–155 | Multi-round match scoring | Dealer rotation, accumulated scores, bonuses, match winner. |
| 156–165 | Social features | Table chat, player profiles, friends, private invitations. |
| 166–175 | Matchmaking and lobbies | Public tables, private rooms, skill filters, reconnect UX. |
| 176–190 | Tournaments | Scheduled tournaments, brackets, standings, prizes only if legally/compliantly appropriate. |
| 191–205 | Bots | Deterministic test bots, then AI opponents. |
| 206–220 | Production operations | Metrics, monitoring, alerting, backups, deployment, load testing. |
| 221–240 | Polish | Animations, sound, accessibility, mobile UX, localization. |

---

# Daily Definition of Done

Every day must end with:

1. A small, working change.
2. New or updated tests.
3. Formatting, linting, type-checking, and relevant test suite passing.
4. Documentation updates where behavior, rules, setup, or protocol changed.
5. A focused Git commit.
6. A push to the active branch.
7. A short handoff note for the next day.

Recommended command sequence:

```bash
git status
npm run format
npm run lint
npm run typecheck
npm test
docker compose up --build -d
docker compose logs --tail=100 nakama
git add <specific-files>
git commit -m "feat: implement stock draw action"
git push
```

---

# Recommended Git Commit Pattern

```text
chore: bootstrap local nakama development environment
docs: define rummy rules and MVP decisions
feat: add deterministic romanian tile deck
feat: create initial rummy round state
feat: add authoritative nakama match lifecycle
feat: implement opening discard turn
feat: validate rummy set melds
feat: validate rummy run melds
feat: allow initial table melds
feat: implement discard row pickup meld
feat: enforce joker replacement rules
test: add deterministic rummy match simulation
docs: document rummy development workflow
refactor: clarify match state transitions
```

This roadmap reaches a stable, fully server-authoritative rummy MVP around **Day 135**, with a minimal manual client beginning around **Day 136**.





# add after Phase 17: 18, 19, 20,21, 22, 23. we need some AI boots features and so on



# Phase 18 — AI Bots: Foundation and Safe Integration

Build bots as normal match participants wherever possible. Bots must use the same authoritative command protocol and rules validation as human players; never give bots direct access to hidden server internals that a real player would not have.

| Day range | Focus | Deliverables |
|---:|---|---|
| 241–243 | Bot architecture decision | Document bot lifecycle, bot identities, difficulty levels, scheduling model, and how bot actions enter Nakama matches. |
| 244–246 | Bot player model | Add bot user/profile metadata, display names, avatar placeholders, bot flags, and safe seat assignment. |
| 247–250 | Bot match participation | Allow a host to add one or more bots to an eligible private/local match. Support 2–4 total seats. |
| 251–254 | Bot action adapter | Implement an adapter that turns bot decisions into normal validated game commands such as `DRAW_STOCK`, `DISCARD`, `MELD_INITIAL`, and `MELD_NEW`. |
| 255–258 | Bot timing | Add configurable think delays, deterministic test mode with no delay, and cancellation if a match ends or bot leaves. |
| 259–262 | Bot visibility boundaries | Ensure a bot only receives its rack plus public state, exactly like a human client. Add tests proving no opponent rack access. |
| 263–266 | Bot error recovery | If a bot generates an invalid action, log it, safely recover, and select a fallback legal move without corrupting the match. |
| 267–270 | Bot observability | Add structured bot-decision logs: match ID, bot ID, turn, selected action, duration, fallback usage, and rejection reason. |

**Milestone:** Bots can join normal matches and take safe, valid actions through the same command path used by humans.

---

# Phase 19 — AI Bots: Legal Move Generation and Beginner Difficulty

Start with correctness, not strength. The first bot should make legal moves predictably and never stall a match.

| Day range | Focus | Deliverables |
|---:|---|---|
| 271–274 | Legal move generator | Create a pure module that enumerates legal actions from a bot-visible game state. |
| 275–278 | Draw decision generation | Determine available draw choices: stock, previous discard, or eligible earlier-discard pickup. |
| 279–282 | Discard candidate generation | Generate every legal discard candidate for the bot’s rack after draw/meld actions. |
| 283–286 | Initial meld candidate generation | Generate possible opening meld batches that satisfy the 50-point and at-least-one-run rule. |
| 287–290 | New meld candidate generation | Generate valid sets/runs from an opened player’s rack. |
| 291–294 | Meld extension candidate generation | Generate legal extensions to public melds, including other players’ melds. |
| 295–298 | Joker replacement candidate generation | Generate legal joker replacement opportunities without allowing illegal joker reuse. |
| 299–302 | Beginner strategy | Implement a simple heuristic bot: prefer opening when possible, play obvious melds, extend melds, then discard a low-value/least-useful tile. |
| 303–306 | Beginner bot tests | Test that the bot always produces legal commands across fixed scenarios and never reveals hidden information. |
| 307–310 | Bot turn simulation | Run deterministic simulated games containing 1–3 beginner bots and verify rounds eventually progress or terminate safely. |

**Milestone:** A beginner bot can complete legal turns, open, meld, discard, and finish games without human intervention.

---

# Phase 20 — AI Bots: Intermediate Strategy and Difficulty Tiers

Add stronger decision-making while preserving understandable, testable logic.

| Day range | Focus | Deliverables |
|---:|---|---|
| 311–314 | Evaluation model | Define an explainable hand-evaluation model: meld potential, opening progress, deadwood value, joker value, run/set flexibility, and discard risk. |
| 315–318 | Candidate scoring | Score all legal candidate actions using the evaluation model. |
| 319–322 | Better discard heuristic | Avoid discarding tiles likely needed for existing sequences, sets, or near-complete melds. |
| 323–326 | Opening strategy | Decide when to preserve tiles versus opening immediately with a legal 50-point batch. |
| 327–330 | Pickup strategy | Score stock draw versus latest discard versus earlier discard-row pickup. |
| 331–334 | Table exploitation | Prefer safe extensions to public melds when they reduce rack size or improve future options. |
| 335–338 | Joker strategy | Evaluate when retaining, melding, replacing, or exposing jokers is beneficial. |
| 339–342 | Difficulty configuration | Add `beginner`, `normal`, and `advanced` profiles with different search depth, randomness, delay, and heuristic quality. |
| 343–346 | Controlled randomness | Add seeded weighted choice so bots are not mechanically identical but remain reproducible in tests. |
| 347–350 | Strategy benchmark suite | Create fixed board-state scenarios and compare candidate choices across difficulty levels. |

**Milestone:** Multiple bot difficulties exist, and higher difficulty makes measurably better decisions without using unfair information.

---

# Phase 21 — AI Bots: Search, Simulation, and Training Infrastructure

Use deeper analysis only after the heuristic bot is stable. Keep compute budgets bounded so bots do not affect match reliability.

| Day range | Focus | Deliverables |
|---:|---|---|
| 351–354 | Search budget design | Define maximum CPU time, candidate count, simulation count, and cancellation behavior per bot turn. |
| 355–358 | One-turn lookahead | Implement bounded evaluation of immediate actions and resulting rack quality. |
| 359–362 | Two-turn approximate planning | Add limited future-turn estimation without assuming knowledge of hidden opponent racks or stock order. |
| 363–366 | Monte Carlo rollout prototype | Simulate plausible unseen stock/opponent distributions from public information only. |
| 367–370 | Simulation state cloning | Build efficient, safe, deterministic copies of game state for bot simulations. |
| 371–374 | Search cancellation | Cancel search immediately when match ends, bot loses turn, player disconnects, or time budget is reached. |
| 375–378 | Fallback policy | Guarantee a legal fallback action if simulation fails, times out, or has no high-confidence result. |
| 379–382 | Offline self-play runner | Add a local script to run thousands of bot-vs-bot rounds with seeds and collect outcomes. |
| 383–386 | Bot telemetry | Collect aggregate metrics: win rates, turn length, invalid-action count, timeout rate, round duration, and action distribution. |
| 387–390 | Performance tests | Verify bot turns meet configured latency and CPU budgets under concurrent local-match load. |

**Milestone:** Advanced bots can use bounded search and offline self-play without harming live match performance or fairness.

---

# Phase 22 — Matchmaking, Ratings, and Bot-Assisted Player Experience

Use bots to improve queue times and onboarding, but never disguise a bot as a human player.

| Day range | Focus | Deliverables |
|---:|---|---|
| 391–394 | Match mode design | Define casual, private, practice, ranked, and bot-only match modes. |
| 395–398 | Practice mode | Allow a player to start a private practice match with configurable bots. |
| 399–402 | Bot transparency UX/protocol | Mark bot players clearly in player lists, match events, logs, and client payloads. |
| 403–406 | Matchmaking queues | Add queue metadata for player count, region, preferred mode, and desired bot difficulty. |
| 407–410 | Queue timeout policy | Define when matchmaking offers bots after a waiting threshold; require user opt-in where appropriate. |
| 411–414 | Backfill bots | Support filling empty seats with bots only before game start or according to explicit match rules. |
| 415–418 | Player rating model | Add an initial rating system for human-vs-human ranked games. Keep bot games separate by default. |
| 419–422 | Rating protection | Ensure bot matches do not inflate/deflate ranked ratings unless explicitly designed and communicated. |
| 423–426 | Bot calibration | Adjust bot difficulty using self-play and human opt-in practice results. |
| 427–430 | Matchmaking tests | Test queue formation, bot insertion rules, rating separation, and bot identity transparency. |

**Milestone:** Players can reliably start practice games with bots, while competitive integrity remains protected.

---

# Phase 23 — Production AI Operations, Safety, and Continuous Improvement

Turn the bot system into a maintainable production capability.

| Day range | Focus | Deliverables |
|---:|---|---|
| 431–434 | AI feature flags | Add server-side feature flags for bot availability, strategy version, delays, search depth, and experimental modes. |
| 435–438 | Strategy versioning | Version bot policies so a live match records exactly which bot strategy generated decisions. |
| 439–442 | Safe rollout process | Add staged rollout: local → development → staging → limited production cohort → full release. |
| 443–446 | Bot health dashboards | Monitor bot action latency, failures, invalid decisions, fallback rate, queue impact, and match completion rate. |
| 447–450 | Alerts and safeguards | Add alerts for abnormal bot timeouts, elevated invalid command rate, stuck matches, or unexpected win-rate shifts. |
| 451–454 | Fairness audits | Periodically verify bots use only player-visible/public information and cannot inspect hidden racks or future stock order. |
| 455–458 | Replay and audit logs | Store privacy-safe action traces for debugging disputed games and diagnosing bot decisions. |
| 459–462 | Human feedback loop | Add optional post-match bot feedback: difficulty too easy/hard, suspicious behavior, fun rating. |
| 463–466 | Automated regression corpus | Convert production bugs and player-reported strange bot behavior into deterministic regression scenarios. |
| 467–470 | Cost and capacity management | Set concurrency limits, execution budgets, worker scaling policies, and cost monitoring for bot computation. |
| 471–474 | Disaster recovery | Define behavior if bot service/scheduler fails: pause, remove bot, substitute basic fallback, or safely end match. |
| 475–480 | AI release candidate | Conduct load tests, fairness tests, replay audits, security review, and staged-release checklist. |

**Milestone:** Bot features are observable, fair, resilient, configurable, and safe to operate in production.

---

# AI Bot Design Rules

Apply these throughout Phases 18–23:

1. **Bots must follow the same rules as humans.**  
   Bot actions must pass through the normal Nakama authoritative command validation path.

2. **Bots must not cheat.**  
   A bot may use:
   - its own rack;
   - public melds;
   - public discard row;
   - public stock count;
   - public turn/player status;
   - its own prior observations.  
   
   A bot must not use:
   - other players’ racks;
   - future stock order;
   - server-only state unavailable to players;
   - hidden random seed data;
   - internal validation shortcuts.

3. **Bot identity must be transparent.**  
   A human must be able to tell that a participant is a bot.

4. **Bot computation must be bounded.**  
   Every decision has a deadline, maximum search budget, and legal fallback action.

5. **Every bot decision must be reproducible in test mode.**  
   Use seeded random choices, versioned strategies, and saved action traces.

6. **Bots should fail safely.**  
   If AI computation fails, select a basic legal action—usually draw stock if required, then a valid discard.

7. **Bot games and ranked games must remain separate by default.**  
   Do not let bot matches affect competitive rating unless rules explicitly say so.

---

This extends the roadmap from the original core MVP through **Day 480**, covering beginner bots, stronger heuristic bots, bounded search/simulation, practice-mode matchmaking, bot-assisted queues, and production AI operations.





# now phase 24,25, 26, 27, 28, 29, 30, 31



# Phase 24 — Production Deployment and Infrastructure

Move from local/staging development to a repeatable, secure production deployment.

| Day range | Focus | Deliverables |
|---:|---|---|
| 481–484 | Production architecture | Document target hosting, regions, domains, TLS, database topology, backups, scaling, and failure boundaries. |
| 485–488 | Infrastructure as Code | Add Terraform, Pulumi, or equivalent IaC for networking, database, container runtime, secrets, and observability resources. |
| 489–492 | Container hardening | Create production Docker images with pinned versions, non-root execution where supported, minimal layers, health checks, and image scanning. |
| 493–496 | Environment management | Separate local, development, staging, and production configurations. Prevent development settings from reaching production. |
| 497–500 | Secrets management | Store database credentials, Nakama keys, API keys, and signing secrets in a managed secret store; never commit secrets. |
| 501–504 | TLS and edge security | Configure HTTPS, certificates, secure headers, CORS rules, rate limits, and WAF/CDN behavior as appropriate. |
| 505–508 | Database migrations | Establish versioned migrations, rollback strategy, migration checks, and backup-before-migrate procedures. |
| 509–512 | Backup and restore drills | Automate PostgreSQL backups and test restoration into an isolated environment. |
| 513–516 | Staging environment | Deploy a production-like staging environment and run automated smoke tests against it. |
| 517–520 | Production release process | Define approvals, release notes, rollback process, incident contacts, and production deployment checklist. |

**Milestone:** The game backend can be deployed, updated, backed up, and rolled back safely in production.

---

# Phase 25 — Observability, Analytics, and Game Telemetry

Make live behavior measurable without exposing private player information.

| Day range | Focus | Deliverables |
|---:|---|---|
| 521–524 | Logging standard | Define structured logging fields: request ID, match ID, player ID hash, bot ID, command opcode, transition result, duration, and error code. |
| 525–528 | Metrics foundation | Add service metrics for requests, match count, active users, command volume, errors, latency, memory, CPU, and database health. |
| 529–532 | Match metrics | Track match creation, joins, starts, disconnects, completion, abandonment, average duration, and player-count distribution. |
| 533–536 | Gameplay metrics | Track draw type, discard pickup use, opening-meld timing, joker usage, win conditions, and rule-validation failures. |
| 537–540 | Bot metrics | Track bot difficulty, decision latency, search depth, fallback rate, win rate, and invalid-action attempts. |
| 541–544 | Product analytics privacy design | Define event retention, player consent, anonymization, aggregation, and personally identifiable information boundaries. |
| 545–548 | Dashboards | Build dashboards for infrastructure health, live matches, gameplay funnel, bot behavior, and error rates. |
| 549–552 | Alerting | Configure alerts for database failures, Nakama restart loops, elevated command errors, stuck matches, slow bot actions, and unusual abandonment. |
| 553–556 | Distributed tracing | Add request/match tracing where technically feasible, especially for external services and AI workers. |
| 557–560 | Analytics validation | Verify event correctness using controlled test matches and ensure no hidden rack details are sent to analytics. |

**Milestone:** Operators can understand system health, game health, and bot health from dashboards and alerts.

---

# Phase 26 — Security, Privacy, and Anti-Cheat

Protect player accounts, match integrity, infrastructure, and private game state.

| Day range | Focus | Deliverables |
|---:|---|---|
| 561–564 | Threat model | Document threats: account abuse, token theft, malformed protocol input, collusion, automation abuse, data leakage, DDoS, and admin compromise. |
| 565–568 | Authentication review | Harden account/session handling, token expiry, refresh flow, password policies, and account recovery. |
| 569–572 | Authorization review | Verify all RPCs, match commands, admin endpoints, and storage access enforce least privilege. |
| 573–576 | Protocol fuzzing | Fuzz malformed, oversized, duplicated, reordered, and adversarial match commands. |
| 577–580 | Rate limiting | Add limits for login, registration, match creation, joining, command spam, chat, and bot-related endpoints. |
| 581–584 | Anti-cheat audit | Verify server authority across all actions, including racks, stock, meld validation, discard pickup, joker replacement, and win detection. |
| 585–588 | Replay integrity | Sign or hash match event streams to detect tampering in stored replays. |
| 589–592 | Privacy controls | Add account-data export/delete workflows and retention policies as required by applicable privacy laws. |
| 593–596 | Dependency security | Add dependency scanning, container scanning, security update policy, and critical-CVE response procedure. |
| 597–600 | Security test and review | Run penetration testing, tabletop incident exercises, and security regression tests. |

**Milestone:** The game has a documented threat model, hardened authoritative gameplay, and practical privacy/security controls.

---

# Phase 27 — Player Accounts, Profiles, and Social Foundation

Add persistent player identity and social systems after game integrity is stable.

| Day range | Focus | Deliverables |
|---:|---|---|
| 601–604 | Account model | Define guest, registered, linked-provider, verified, suspended, and deleted account states. |
| 605–608 | Registration and login UX/API | Support secure email/password and/or platform-provider login flows. |
| 609–612 | Guest accounts | Allow low-friction guest play with a documented upgrade/link path. |
| 613–616 | Player profile | Add display name, avatar selection/upload policy, language, region, and safe public profile fields. |
| 617–620 | Profile moderation controls | Add validation, reserved names, report flow, avatar limits, and admin moderation tools. |
| 621–624 | Friends model | Add friend requests, accepted friends, blocked users, and privacy controls. |
| 625–628 | Presence and invitations | Show opt-in online presence and allow private match invitations. |
| 629–632 | Recent players | Store recent opponents with options to add friend, block, report, or invite again. |
| 633–636 | Social privacy tests | Verify blocks, privacy settings, invitations, and profile visibility work correctly. |
| 637–640 | Account lifecycle documentation | Document account recovery, deletion, moderation, support workflows, and audit trails. |

**Milestone:** Players have durable identities and can safely interact socially around games.

---

# Phase 28 — Communication, Moderation, and Community Safety

Introduce chat and community features with safety controls from the first implementation day.

| Day range | Focus | Deliverables |
|---:|---|---|
| 641–644 | Communication scope | Define supported channels: table chat, private messages, friends chat, lobby chat, and system announcements. |
| 645–648 | Table chat MVP | Implement match/table chat with message limits, ordering, and reconnection behavior. |
| 649–652 | Chat persistence decision | Document whether messages are ephemeral, match-scoped, retained, and how long retention lasts. |
| 653–656 | Spam protection | Add rate limits, duplicate-message controls, cooldowns, and flood detection. |
| 657–660 | Player blocking | Ensure blocked users cannot message, invite, or interact with each other where applicable. |
| 661–664 | Reporting workflow | Add report categories, evidence capture, moderation queue, and player feedback states. |
| 665–668 | Automated content filtering | Add configurable profanity/spam filters with language-aware policies and safe false-positive handling. |
| 669–672 | Moderator tools | Add review queues, temporary mutes, warnings, suspensions, bans, and audit logs. |
| 673–676 | Appeals and enforcement policy | Document moderation standards, appeal process, staff permissions, and evidence retention. |
| 677–680 | Community safety tests | Test rate limits, blocking, reporting, moderation actions, and auditability. |

**Milestone:** Players can communicate while the platform has practical tools to prevent spam, harassment, and abuse.

---

# Phase 29 — Economy, Progression, and Cosmetics

Keep the game non-gambling. Do not introduce cash wagering, real-money stakes, loot boxes, or mechanics that resemble gambling.

| Day range | Focus | Deliverables |
|---:|---|---|
| 681–684 | Economy principles | Document non-gambling rules, fair-play constraints, virtual-currency policy, regional compliance review, and parental-control considerations. |
| 685–688 | Player progression | Add experience, levels, badges, and non-competitive achievement tracking. |
| 689–692 | Match rewards | Award non-monetary progression rewards for completed matches, good sportsmanship, and tutorials. |
| 693–696 | Cosmetics model | Define cosmetic-only items: table themes, tile skins, avatars, emotes, profile frames, and victory effects. |
| 697–700 | Inventory service | Implement ownership, unlocks, equipped items, grants, revocations, and inventory audit history. |
| 701–704 | Daily/weekly challenges | Add optional goals that reward engagement without pressuring risky spending or exploitative play patterns. |
| 705–708 | Achievement system | Add rule-based achievements, progress tracking, notifications, and anti-exploit checks. |
| 709–712 | Store design, if applicable | If monetization is approved, implement cosmetic purchases with transparent pricing, receipts, refunds, parental controls, and no pay-to-win effect. |
| 713–716 | Economy abuse prevention | Detect duplicate grants, replayed purchase events, reward farming, and inventory inconsistencies. |
| 717–720 | Economy tests and compliance review | Test transactions and entitlement handling; conduct legal/compliance review before release. |

**Milestone:** Players can progress and personalize their experience without affecting rummy fairness or creating gambling-like mechanics.

---

# Phase 30 — Competitive Play, Rankings, and Tournaments

Build competitive systems only after match integrity, anti-cheat, moderation, and operations are mature.

| Day range | Focus | Deliverables |
|---:|---|---|
| 721–724 | Competitive ruleset | Define ranked match eligibility, player count, disconnect policy, bot policy, scoring format, and seasonal structure. |
| 725–728 | Rating algorithm | Select and implement a rating model such as Elo, Glicko, or TrueSkill-like logic appropriate for multiplayer games. |
| 729–732 | Rating event pipeline | Record rating-relevant match results atomically and make recalculation/recovery possible. |
| 733–736 | Leaderboards | Add global, regional, friends, seasonal, and game-mode leaderboards. |
| 737–740 | Ranked matchmaking | Match players using rating, region, latency, queue duration, and anti-smurf controls. |
| 741–744 | Season system | Add season start/end dates, soft reset policy, rewards, historical ranks, and announcement flow. |
| 745–748 | Tournament design | Define tournament formats: scheduled tables, Swiss, bracket, elimination, or points-based events. |
| 749–752 | Tournament engine | Implement registration, seating, round creation, advancement, standings, and tie-breaking. |
| 753–756 | Tournament integrity | Enforce check-in, disconnect rules, bot restrictions, anti-collusion signals, staff intervention, and replay review. |
| 757–760 | Tournament rewards | Add cosmetic/badge rewards and transparent eligibility rules. Avoid real-money prizes unless legal, licensed, and separately approved. |
| 761–764 | Competitive testing | Simulate rating updates, tournament progression, player dropouts, ties, and rollback/recovery cases. |
| 765–770 | Competitive launch readiness | Load test queues/leaderboards, audit integrity controls, publish rules, and run a limited beta season. |

**Milestone:** The platform can support fair ranked play, leaderboards, seasons, and controlled tournament events.

---

# Phase 31 — Client Platforms, UX, Accessibility, and Localization

Turn the backend MVP into a polished, accessible, multi-platform product.

| Day range | Focus | Deliverables |
|---:|---|---|
| 771–774 | Client platform strategy | Decide web, desktop, Android, iOS, and/or native-client priorities; define shared protocol/client SDK strategy. |
| 775–778 | Design system | Establish reusable typography, colours, spacing, buttons, modals, notifications, tile components, and table layouts. |
| 779–782 | Rummy table UX | Build the core table interaction model: rack sorting, drag/drop, tile selection, undo-before-submit, and turn feedback. |
| 783–786 | Meld-building UX | Support accessible creation of runs/sets, explicit joker assignment, invalid-state feedback, and server-confirmed actions. |
| 787–790 | Responsive layouts | Support desktop, tablet, and mobile layouts without reducing rule clarity or touch usability. |
| 791–794 | Accessibility baseline | Add keyboard control, screen-reader labels, focus handling, colour-blind-safe tile identifiers, scalable text, and reduced motion. |
| 795–798 | Localization foundation | Externalize UI text, support Romanian and English first, and define pluralization/date/number formatting rules. |
| 799–802 | Connection UX | Add reconnecting states, resync UI, action-pending feedback, network error handling, and match recovery. |
| 803–806 | Onboarding and tutorial | Create guided rules introduction, practice match with beginner bots, contextual help, and rule glossary. |
| 807–810 | Notifications | Add opt-in notifications for invitations, turn reminders where appropriate, friend activity, tournaments, and moderation actions. |
| 811–814 | Client QA automation | Add end-to-end client tests for login, match creation, joining, draw/discard, melding, reconnection, and error states. |
| 815–820 | UX beta and polish | Conduct usability testing, fix high-impact friction, improve performance, and prepare release notes. |

**Milestone:** Players can use a polished, accessible, localized client to play Romanian Tile Rummy across target platforms.

---

# Extended Roadmap Summary

| Phase | Days | Primary Outcome |
|---|---:|---|
| 24 | 481–520 | Production infrastructure and deployment |
| 25 | 521–560 | Observability, analytics, and telemetry |
| 26 | 561–600 | Security, privacy, and anti-cheat |
| 27 | 601–640 | Accounts, profiles, and social foundation |
| 28 | 641–680 | Chat, moderation, and community safety |
| 29 | 681–720 | Safe progression, cosmetics, and optional economy |
| 30 | 721–770 | Ranked play, leaderboards, and tournaments |
| 31 | 771–820 | Multi-platform client, UX, accessibility, localization |

By the end of **Day 820**, the project evolves from a Nakama-authoritative rummy backend into a production-ready, AI-supported, secure, social, competitive, and accessible multiplayer platform.




## 参考资料

1. [Romanian Tile Rummy - rules of the game](https://www.pagat.com/rummy/romtile.html) 
2. [Remi Online - rummy jocuri si socializare](https://www.remi-online.ro) 
