// Package match — Shared meld parsing (Day 23 refactor).
// Extracted from meld_initial.go / meld_new.go to remove duplication.
// Keeps match orchestration thin and rules validation pure per AGENTS.md.
package match

import (
	"encoding/json"
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

// clientMeldPayload is the wire format for both MELD_INITIAL and MELD_NEW.
// tileIds are validated against the sender's rack; jokerReps maps jokerId -> represented tile.
type clientMeldPayload struct {
	ID        string                     `json:"id"`                  // stable meld ID, server generates if empty
	Kind      string                     `json:"kind"`                // "run" or "set"
	TileIDs   []string                   `json:"tileIds"`             // 3+ owned tile instance IDs
	JokerReps map[string]jokerRepPayload `json:"jokerReps,omitempty"` // jokerId -> represented colour/rank
}

type jokerRepPayload struct {
	Colour string `json:"colour"` // "red"/"yellow"/"blue"/"black"
	Rank   int    `json:"rank"`   // 1..13
}

type meldBatchPayload struct {
	Melds []clientMeldPayload `json:"melds"`
}

// findPlayerIndex returns the index of senderId in st.Players, or -1 if not found.
func findPlayerIndex(st *RoundState, senderId PlayerId) int {
	for i, p := range st.Players {
		if p.ID == senderId {
			return i
		}
	}
	return -1
}

// requireMeldOrDiscardTurn validates that the state is Playing/MeldOrDiscard and sender is the current player.
// It sends OpServerError on failure and returns the Seat or SeatInvalid.
func requireMeldOrDiscardTurn(st *RoundState, senderId PlayerId, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) (Seat, int, error) {
	if st.GamePhase != PhasePlaying || st.TurnPhase != TurnMeldOrDiscard {
		sendError(dispatcher, sender, protocol.ErrCodeWrongPhase, fmt.Sprintf("meld only in Playing MeldOrDiscard, got %v/%v", st.GamePhase, st.TurnPhase), requestId, op, logger)
		return SeatInvalid, -1, fmt.Errorf("wrong phase")
	}
	seat := SeatOfPlayer(st.Players, senderId)
	if seat == SeatInvalid {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "not in match", requestId, op, logger)
		return SeatInvalid, -1, fmt.Errorf("not member")
	}
	if seat != st.CurrentSeat {
		sendError(dispatcher, sender, protocol.ErrCodeNotYourTurn, fmt.Sprintf("not your turn: current %v sender %v", st.CurrentSeat, seat), requestId, op, logger)
		return SeatInvalid, -1, fmt.Errorf("not your turn")
	}
	idx := findPlayerIndex(st, senderId)
	if idx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "player not found", requestId, op, logger)
		return SeatInvalid, -1, fmt.Errorf("player not found")
	}
	return seat, idx, nil
}

// buildMeldsFromPayload parses the batch payload, resolves tileIds against rack, builds JokerReps,
// validates each meld via meld.New, and checks duplicate tileIds and meldId uniqueness.
// existingIDs is the set of TableMeld IDs already on the table (for collision check); may be nil for initial.
// It sends OpServerError on any validation failure and is atomic (no state mutated).
func buildMeldsFromPayload(payload []byte, rack []tile.TileInstance, existingIDs map[meld.MeldID]bool, seat Seat, tableSize int, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) ([]meld.Meld, map[tile.TileInstanceId]bool, error) {
	var req meldBatchPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld payload bad JSON: %v", err), requestId, op, logger)
		return nil, nil, err
	}
	if len(req.Melds) == 0 {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "melds must have at least 1 meld", requestId, op, logger)
		return nil, nil, fmt.Errorf("no melds")
	}

	rackByID := make(map[tile.TileInstanceId]tile.TileInstance, len(rack))
	for _, t := range rack {
		rackByID[t.ID] = t
	}

	scoringMelds := make([]meld.Meld, 0, len(req.Melds))
	seenTile := map[tile.TileInstanceId]bool{}
	seenMeldID := map[meld.MeldID]bool{}

	for idx, mp := range req.Melds {
		if mp.Kind != string(meld.KindRun) && mp.Kind != string(meld.KindSet) {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %d invalid kind %q must be run or set", idx, mp.Kind), requestId, op, logger)
			return nil, nil, fmt.Errorf("invalid kind %q", mp.Kind)
		}
		if len(mp.TileIDs) < 3 {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %d must have >=3 tiles, got %d", idx, len(mp.TileIDs)), requestId, op, logger)
			return nil, nil, fmt.Errorf("meld %d too few tiles", idx)
		}
		meldID := meld.MeldID(mp.ID)
		if meldID == "" {
			meldID = meld.MeldID(fmt.Sprintf("meld-%d-%d-%d", int(seat), tableSize, idx))
		}
		if seenMeldID[meldID] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate meld id %q in batch", meldID), requestId, op, logger)
			return nil, nil, fmt.Errorf("duplicate meld id %q", meldID)
		}
		if existingIDs != nil && existingIDs[meldID] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld id %q already exists on table", meldID), requestId, op, logger)
			return nil, nil, fmt.Errorf("meld id exists %q", meldID)
		}
		seenMeldID[meldID] = true

		tiles := make([]tile.TileInstance, 0, len(mp.TileIDs))
		for _, tidStr := range mp.TileIDs {
			tid := tile.TileInstanceId(tidStr)
			if tid == "" {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q has empty tileId", meldID), requestId, op, logger)
				return nil, nil, fmt.Errorf("empty tileId")
			}
			if seenTile[tid] {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate tile %v across melds", tid), requestId, op, logger)
				return nil, nil, fmt.Errorf("duplicate tile %v", tid)
			}
			t, ok := rackByID[tid]
			if !ok {
				sendError(dispatcher, sender, protocol.ErrCodeInvalidTile, fmt.Sprintf("tile %v not in rack", tid), requestId, op, logger)
				return nil, nil, fmt.Errorf("tile %v not in rack", tid)
			}
			seenTile[tid] = true
			tiles = append(tiles, t)
		}

		reps := make(map[tile.TileInstanceId]tile.TileInstance, len(mp.JokerReps))
		for jidStr, rep := range mp.JokerReps {
			jid := tile.TileInstanceId(jidStr)
			found := false
			for _, tl := range tiles {
				if tl.ID == jid {
					found = true
					if !tl.IsJoker {
						sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q rep for non-joker %v", meldID, jid), requestId, op, logger)
						return nil, nil, fmt.Errorf("rep for non-joker %v", jid)
					}
					break
				}
			}
			if !found {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q rep for tile %v not in tileIds", meldID, jid), requestId, op, logger)
				return nil, nil, fmt.Errorf("rep for non-existent tile %v", jid)
			}
			colour, ok := tile.ParseColour(rep.Colour)
			if !ok {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v invalid colour %q", meldID, jid, rep.Colour), requestId, op, logger)
				return nil, nil, fmt.Errorf("invalid colour %q", rep.Colour)
			}
			rank := tile.Rank(rep.Rank)
			if !rank.IsValid() {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v invalid rank %d", meldID, jid, rep.Rank), requestId, op, logger)
				return nil, nil, fmt.Errorf("invalid rank %d", rep.Rank)
			}
			repTile := tile.MustTile(tile.TileInstanceId(fmt.Sprintf("rep-%s", jid)), colour, rank)
			reps[jid] = repTile
		}
		for _, t := range tiles {
			if t.IsJoker {
				if _, ok := reps[t.ID]; !ok {
					sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v missing rep", meldID, t.ID), requestId, op, logger)
					return nil, nil, fmt.Errorf("missing rep for joker %v", t.ID)
				}
			}
		}
		kind := meld.Kind(mp.Kind)
		m, err := meld.New(meldID, kind, tiles, reps)
		if err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld %q invalid: %v", meldID, err), requestId, op, logger)
			return nil, nil, err
		}
		scoringMelds = append(scoringMelds, m)
	}
	return scoringMelds, seenTile, nil
}

// applyMeldBatch atomically removes seen tiles from rack and appends TableMelds.
// It rolls back rack mutation if any TableMeld.Validate fails (should not happen after meld.New).
func applyMeldBatch(st *RoundState, seat Seat, originalRack []tile.TileInstance, scoringMelds []meld.Meld, seenTile map[tile.TileInstanceId]bool, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	newRack := make([]tile.TileInstance, 0, len(originalRack)-len(seenTile))
	for _, t := range originalRack {
		if !seenTile[t.ID] {
			newRack = append(newRack, t)
		}
	}
	st.Racks[seat] = newRack
	for _, m := range scoringMelds {
		tm := TableMeld{
			ID:        string(m.ID),
			Tiles:     m.Tiles,
			JokerReps: m.JokerReps,
			OwnerSeat: seat,
		}
		if err := tm.Validate(); err != nil {
			st.Racks[seat] = originalRack
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("table meld invalid: %v", err), requestId, op, logger)
			return err
		}
		st.TableMelds = append(st.TableMelds, tm)
	}
	return nil
}
