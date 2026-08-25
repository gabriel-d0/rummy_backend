// Package match — Additional melds handler (AGENTS.md Day 14, roadmap Day 68-70).
// Implements MELD_NEW: opened players may create one or more new independent
// sets/runs from rack, no score minimum, atomic, preserves TableMeld IDs per
// docs/rules-decisions.md:5.
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

// meldNewPayload is identical wire format to initial but validated differently.
type meldNewPayload struct {
	Melds []meldPayload `json:"melds"`
}

// handleMeldNew validates and applies MELD_NEW. Atomic, requires HasOpened.
func handleMeldNew(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	if st.GamePhase != PhasePlaying || st.TurnPhase != TurnMeldOrDiscard {
		sendError(dispatcher, sender, protocol.ErrCodeWrongPhase, fmt.Sprintf("meld_new only in Playing MeldOrDiscard, got %v/%v", st.GamePhase, st.TurnPhase), requestId, op, logger)
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
	if !st.Players[pIdx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeNotOpened, "must open with MELD_INITIAL first (50+ with run)", requestId, op, logger)
		return fmt.Errorf("not opened")
	}

	var req meldNewPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld_new payload bad JSON: %v", err), requestId, op, logger)
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
	// Existing meld IDs for collision check
	existingIDs := map[meld.MeldID]bool{}
	for _, m := range st.TableMelds {
		existingIDs[meld.MeldID(m.ID)] = true
	}

	scoringMelds := make([]meld.Meld, 0, len(req.Melds))
	seenTile := map[tile.TileInstanceId]bool{}
	seenMeldID := map[meld.MeldID]bool{}

	for idx, mp := range req.Melds {
		if mp.Kind != string(meld.KindRun) && mp.Kind != string(meld.KindSet) {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %d invalid kind %q must be run or set", idx, mp.Kind), requestId, op, logger)
			return fmt.Errorf("invalid kind %q", mp.Kind)
		}
		if len(mp.TileIDs) < 3 {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %d must have >=3 tiles, got %d", idx, len(mp.TileIDs)), requestId, op, logger)
			return fmt.Errorf("meld %d too few tiles", idx)
		}
		meldID := meld.MeldID(mp.ID)
		if meldID == "" {
			meldID = meld.MeldID(fmt.Sprintf("meld-%d-%d-%d", int(seat), len(st.TableMelds), idx))
		}
		if seenMeldID[meldID] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate meld id %q in batch", meldID), requestId, op, logger)
			return fmt.Errorf("duplicate meld id %q", meldID)
		}
		if existingIDs[meldID] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld id %q already exists on table", meldID), requestId, op, logger)
			return fmt.Errorf("meld id exists %q", meldID)
		}
		seenMeldID[meldID] = true

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

		reps := make(map[tile.TileInstanceId]tile.TileInstance, len(mp.JokerReps))
		for jidStr, rep := range mp.JokerReps {
			jid := tile.TileInstanceId(jidStr)
			foundInTiles := false
			for i := range tiles {
				if tiles[i].ID == jid {
					foundInTiles = true
					if !tiles[i].IsJoker {
						sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q rep for non-joker %v", meldID, jid), requestId, op, logger)
						return fmt.Errorf("rep for non-joker %v", jid)
					}
					break
				}
			}
			if !foundInTiles {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q rep for tile %v not in tileIds", meldID, jid), requestId, op, logger)
				return fmt.Errorf("rep for non-existent tile %v", jid)
			}
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
			repTile := tile.MustTile(tile.TileInstanceId(fmt.Sprintf("rep-%s", jid)), colour, rank)
			reps[jid] = repTile
		}
		for _, t := range tiles {
			if t.IsJoker {
				if _, ok := reps[t.ID]; !ok {
					sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q joker %v missing rep", meldID, t.ID), requestId, op, logger)
					return fmt.Errorf("missing rep for joker %v", t.ID)
				}
			}
		}
		kind := meld.Kind(mp.Kind)
		m, err := meld.New(meldID, kind, tiles, reps)
		if err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld %q invalid: %v", meldID, err), requestId, op, logger)
			return err
		}
		scoringMelds = append(scoringMelds, m)
	}

	// Validate ownership and each meld's run/set rules (no score/hasRun)
	batch := scoring.Batch{PlayerID: string(senderId), Melds: scoringMelds}
	if err := scoring.ValidateBatchOwnership(batch, rack); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld_new invalid: %v", err), requestId, op, logger)
		return err
	}

	// ATOMIC MUTATION
	newRack := make([]tile.TileInstance, 0, len(rack)-len(seenTile))
	for _, t := range rack {
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
			st.Racks[seat] = rack
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("table meld invalid: %v", err), requestId, op, logger)
			return err
		}
		st.TableMelds = append(st.TableMelds, tm)
	}
	// HasOpened already true, stays true; phase stays MeldOrDiscard
	logger.Info("MeldNew by %s seat %v melds %d newRack %d table %d", senderId, seat, len(scoringMelds), len(newRack), len(st.TableMelds))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"meldNew","seat":%d,"melds":%d}`, int(seat), len(scoringMelds))), nil, nil, true)
	}
	return nil
}
