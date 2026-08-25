// Package match — Initial meld handler (Day 65, AGENTS.md Day 13).
// Implements MELD_INITIAL: batch of melds from player's rack, >=50 points, at least one run,
// atoms, marks HasOpened, stays in MeldOrDiscard per docs/rules-decisions.md:1.4 and 5.
// Pure validation delegated to scoring.ValidateInitialBatch; here we handle wire payload
// parsing, tile lookup, TableMeld construction, and state mutation.
package match

import (
	"encoding/json"
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/scoring"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

// meldInitialPayload is the client wire format for MELD_INITIAL.
// It mirrors the protocol's {melds: [...]} payload; each meld carries tileIds
// (server resolves colour/rank via rack lookup, never trusts client face values)
// and explicit jokerReps (jokerId -> represented colour/rank) per docs/rules-decisions.md:3.
type meldInitialPayload struct {
	Melds []meldPayload `json:"melds"`
}

type meldPayload struct {
	ID        string                     `json:"id"`                  // stable meld ID; if empty server generates one
	Kind      string                     `json:"kind"`                // "run" or "set"
	TileIDs   []string                   `json:"tileIds"`             // 3+ tile instance IDs, owned in rack
	JokerReps map[string]jokerRepPayload `json:"jokerReps,omitempty"` // jokerId -> represented tile
}

type jokerRepPayload struct {
	Colour string `json:"colour"` // "red"/"yellow"/"blue"/"black"
	Rank   int    `json:"rank"`   // 1..13
}

// handleMeldInitial validates and applies MELD_INITIAL. It is atomic: on any
// validation error state is unchanged and an OpServerError is sent.
func handleMeldInitial(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	// Phase/turn pre-check (MatchLoop already did, but double-check for direct calls)
	if st.GamePhase != PhasePlaying || st.TurnPhase != TurnMeldOrDiscard {
		sendError(dispatcher, sender, protocol.ErrCodeWrongPhase, fmt.Sprintf("meld_initial only in Playing MeldOrDiscard, got %v/%v", st.GamePhase, st.TurnPhase), requestId, op, logger)
		return fmt.Errorf("wrong phase")
	}
	seat := SeatOfPlayer(st.Players, senderId)
	if seat == SeatInvalid {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "not in match", requestId, op, logger)
		return fmt.Errorf("not member")
	}
	if seat != st.CurrentSeat {
		sendError(dispatcher, sender, protocol.ErrCodeNotYourTurn, fmt.Sprintf("not your turn: current %v sender %v", st.CurrentSeat, seat), requestId, op, logger)
		return fmt.Errorf("not your turn")
	}
	// Find player index and check HasOpened
	pIdx := -1
	for i, p := range st.Players {
		if p.ID == senderId {
			pIdx = i
			break
		}
	}
	if pIdx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "player not found", requestId, op, logger)
		return fmt.Errorf("player not found")
	}
	if st.Players[pIdx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeAlreadyOpened, "already opened: use MELD_NEW", requestId, op, logger)
		return fmt.Errorf("already opened")
	}

	// Parse payload
	var req meldInitialPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld_initial payload bad JSON: %v", err), requestId, op, logger)
		return err
	}
	if len(req.Melds) == 0 {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "melds must have at least 1 meld", requestId, op, logger)
		return fmt.Errorf("no melds")
	}

	rack := st.Racks[seat]
	rackByID := make(map[tile.TileInstanceId]tile.TileInstance, len(rack))
	for _, t := range rack {
		rackByID[t.ID] = t
	}

	// Build meld.Meld objects for scoring validation
	scoringMelds := make([]meld.Meld, 0, len(req.Melds))
	seenTile := map[tile.TileInstanceId]bool{}
	seenMeldID := map[meld.MeldID]bool{}

	for idx, mp := range req.Melds {
		// Validate kind
		if mp.Kind != string(meld.KindRun) && mp.Kind != string(meld.KindSet) {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %d invalid kind %q must be run or set", idx, mp.Kind), requestId, op, logger)
			return fmt.Errorf("invalid kind %q", mp.Kind)
		}
		if len(mp.TileIDs) < 3 {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %d must have >=3 tiles, got %d", idx, len(mp.TileIDs)), requestId, op, logger)
			return fmt.Errorf("meld %d too few tiles", idx)
		}
		// Meld ID: require non-empty for determinism in tests; generate fallback if empty
		meldID := meld.MeldID(mp.ID)
		if meldID == "" {
			// Generate deterministic ID from seat + index + existing table size
			meldID = meld.MeldID(fmt.Sprintf("meld-%d-%d-%d", int(seat), len(st.TableMelds), idx))
		}
		if seenMeldID[meldID] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate meld id %q", meldID), requestId, op, logger)
			return fmt.Errorf("duplicate meld id %q", meldID)
		}
		seenMeldID[meldID] = true

		// Resolve tiles from rack
		tiles := make([]tile.TileInstance, 0, len(mp.TileIDs))
		for _, tidStr := range mp.TileIDs {
			tid := tile.TileInstanceId(tidStr)
			if tid == "" {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q has empty tileId", meldID), requestId, op, logger)
				return fmt.Errorf("empty tileId")
			}
			if seenTile[tid] {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate tile %v across melds", tid), requestId, op, logger)
				return fmt.Errorf("duplicate tile %v", tid)
			}
			t, ok := rackByID[tid]
			if !ok {
				sendError(dispatcher, sender, protocol.ErrCodeInvalidTile, fmt.Sprintf("tile %v not in rack", tid), requestId, op, logger)
				return fmt.Errorf("tile %v not in rack", tid)
			}
			seenTile[tid] = true
			tiles = append(tiles, t)
		}

		// Build JokerReps map for this meld
		reps := make(map[tile.TileInstanceId]tile.TileInstance, len(mp.JokerReps))
		for jidStr, rep := range mp.JokerReps {
			jid := tile.TileInstanceId(jidStr)
			// Check joker is in tileIds and is actually a joker
			foundInTiles := false
			var jokerTile *tile.TileInstance
			for i := range tiles {
				if tiles[i].ID == jid {
					foundInTiles = true
					if !tiles[i].IsJoker {
						sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q rep for non-joker %v", meldID, jid), requestId, op, logger)
						return fmt.Errorf("rep for non-joker %v", jid)
					}
					jokerTile = &tiles[i]
					break
				}
			}
			if !foundInTiles {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q rep for tile %v not in tileIds", meldID, jid), requestId, op, logger)
				return fmt.Errorf("rep for non-existent tile %v", jid)
			}
			_ = jokerTile
			// Parse colour
			colour, ok := tile.ParseColour(rep.Colour)
			if !ok {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v invalid colour %q", meldID, jid, rep.Colour), requestId, op, logger)
				return fmt.Errorf("invalid colour %q", rep.Colour)
			}
			rank := tile.Rank(rep.Rank)
			if !rank.IsValid() {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v invalid rank %d", meldID, jid, rep.Rank), requestId, op, logger)
				return fmt.Errorf("invalid rank %d", rep.Rank)
			}
			// Rep ID is synthetic; content matters only for colour/rank and IsJoker false
			repTile := tile.MustTile(tile.TileInstanceId(fmt.Sprintf("rep-%s", jid)), colour, rank)
			// But MustTile validates that ID non-empty; we use rep- prefix
			// However we should not reuse MustTile with rep ID that collides; use a unique placeholder
			// Actually tile.MustTile requires ID, but rep's ID is not checked for uniqueness; we use repTile as value
			reps[jid] = repTile
		}
		// Also need to check that every joker in tiles has a rep, and no extra reps
		// meld.New will validate, but we can give a better error here
		// Count jokers in tiles
		for _, t := range tiles {
			if t.IsJoker {
				if _, ok := reps[t.ID]; !ok {
					sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v missing rep", meldID, t.ID), requestId, op, logger)
					return fmt.Errorf("missing rep for joker %v", t.ID)
				}
			}
		}
		// Create meld.Meld via meld.New (validates structure, duplicate tiles, missing rep, etc.)
		kind := meld.Kind(mp.Kind)
		m, err := meld.New(meldID, kind, tiles, reps)
		if err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld %q invalid: %v", meldID, err), requestId, op, logger)
			return err
		}
		scoringMelds = append(scoringMelds, m)
	}

	// Now validate the whole batch via scoring rules: ownership, score >=50, at least one run, duplicate, each meld valid
	batch := scoring.Batch{PlayerID: string(senderId), Melds: scoringMelds}
	if err := scoring.ValidateInitialBatch(batch, rack); err != nil {
		// Map to appropriate error code: score/ hasRun -> bad_request, ownership/invalid -> specific
		code := protocol.ErrCodeBadRequest
		// Use string inspection to pick more specific codes if needed, but keep bad_request for most
		// For duplicate tile across melds or not in rack, we already handled earlier, so remaining is score/hasRun or meld validation
		// Keep code as bad_request; client can display message
		sendError(dispatcher, sender, code, fmt.Sprintf("initial batch invalid: %v", err), requestId, op, logger)
		return err
	}

	// ATOMIC MUTATION: all validation passed, now apply
	// Remove tiles from rack
	newRack := make([]tile.TileInstance, 0, len(rack)-len(seenTile))
	for _, t := range rack {
		if !seenTile[t.ID] {
			newRack = append(newRack, t)
		}
	}
	st.Racks[seat] = newRack
	// Append TableMelds
	for _, m := range scoringMelds {
		tm := TableMeld{
			ID:        string(m.ID),
			Tiles:     m.Tiles,
			JokerReps: m.JokerReps,
			OwnerSeat: seat,
		}
		// Validate TableMeld structure (should pass)
		if err := tm.Validate(); err != nil {
			// This should not happen due to earlier validation, but if it does, rollback?
			// For atomicity we would need to revert rack; but since we already mutated rack, we log and revert
			// Revert rack to original for safety
			st.Racks[seat] = rack
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("table meld invalid: %v", err), requestId, op, logger)
			return err
		}
		st.TableMelds = append(st.TableMelds, tm)
	}
	// Mark opened
	st.Players[pIdx].HasOpened = true
	// Keep in MeldOrDiscard (do not advance turn; player must discard)
	// st.TurnPhase stays MeldOrDiscard, st.CurrentSeat unchanged

	logger.Info("MeldInitial by %s seat %v melds %d newRack %d table %d", senderId, seat, len(scoringMelds), len(newRack), len(st.TableMelds))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"meldInitial","seat":%d,"melds":%d}`, int(seat), len(scoringMelds))), nil, nil, true)
	}
	return nil
}
