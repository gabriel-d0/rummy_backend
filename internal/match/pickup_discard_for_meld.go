// Package match — Discard row pickup for meld handler (AGENTS.md Day 17).
// Implements PICKUP_DISCARD_FOR_MELD: opened active player in MustDraw may
// take a non-opening discard plus all later discards, but must immediately
// meld the selected discard with exactly two rack tiles into a valid set/run.
package match

import (
	"encoding/json"
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

type pickupDiscardPayload struct {
	DiscardIndex int                        `json:"discardIndex"`
	TileIDs      []string                   `json:"tileIds"` // exactly 2
	JokerReps    map[string]jokerRepPayload `json:"jokerReps,omitempty"`
	Kind         *string                    `json:"kind,omitempty"` // optional "run" or "set"
}

// handlePickupDiscardForMeld validates and applies PICKUP_DISCARD_FOR_MELD atomically.
func handlePickupDiscardForMeld(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	if st.GamePhase != PhasePlaying || st.TurnPhase != TurnMustDraw {
		sendError(dispatcher, sender, protocol.ErrCodeWrongPhase, fmt.Sprintf("pickup only in Playing MustDraw, got %v/%v", st.GamePhase, st.TurnPhase), requestId, op, logger)
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
	idx := findPlayerIndex(st, senderId)
	if idx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "player not found", requestId, op, logger)
		return fmt.Errorf("player not found")
	}
	if !st.Players[idx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeNotOpened, "must open with MELD_INITIAL first (50+ with run)", requestId, op, logger)
		return fmt.Errorf("not opened")
	}

	var req pickupDiscardPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("pickup payload bad JSON: %v", err), requestId, op, logger)
		return err
	}
	if len(req.TileIDs) != 2 {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "tileIds must have exactly 2", requestId, op, logger)
		return fmt.Errorf("tileIds len !=2")
	}
	if perr := CanPickupDiscardForMeld(st, req.DiscardIndex); perr != nil {
		sendError(dispatcher, sender, perr.Code, perr.Message, requestId, op, logger)
		return fmt.Errorf("%s: %s", perr.Code, perr.Message)
	}

	// Validate tileIds owned and not duplicate
	rack := st.Racks[seat]
	rackByID := make(map[tile.TileInstanceId]tile.TileInstance, len(rack))
	for _, t := range rack {
		rackByID[t.ID] = t
	}
	seen := map[tile.TileInstanceId]bool{}
	rackTiles := make([]tile.TileInstance, 0, 2)
	for _, tidStr := range req.TileIDs {
		tid := tile.TileInstanceId(tidStr)
		if tid == "" {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "tileId empty", requestId, op, logger)
			return fmt.Errorf("empty tileId")
		}
		if seen[tid] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate tile %v", tid), requestId, op, logger)
			return fmt.Errorf("duplicate tile %v", tid)
		}
		t, ok := rackByID[tid]
		if !ok {
			sendError(dispatcher, sender, protocol.ErrCodeInvalidTile, fmt.Sprintf("tile %v not in rack", tid), requestId, op, logger)
			return fmt.Errorf("tile %v not in rack", tid)
		}
		seen[tid] = true
		rackTiles = append(rackTiles, t)
	}
	// Also ensure rack tiles not duplicate with discard tile ID? That would be duplicate tile across game, but check anyway
	discardTile := st.DiscardRow[req.DiscardIndex].Tile
	if seen[discardTile.ID] {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("tile %v duplicate with discard", discardTile.ID), requestId, op, logger)
		return fmt.Errorf("duplicate with discard")
	}

	// Build combined tiles for meld validation: discard + 2 rack tiles
	combinedTiles := make([]tile.TileInstance, 0, 3)
	combinedTiles = append(combinedTiles, discardTile)
	combinedTiles = append(combinedTiles, rackTiles...)

	// Build JokerReps map for any joker among the three
	combinedReps := map[tile.TileInstanceId]tile.TileInstance{}
	// req.JokerReps may contain reps for any of the three jokers
	for jidStr, rep := range req.JokerReps {
		jid := tile.TileInstanceId(jidStr)
		// Must be one of the three tiles and be a joker
		found := false
		for _, tl := range combinedTiles {
			if tl.ID == jid {
				found = true
				if !tl.IsJoker {
					sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("pickup rep for non-joker %v", jid), requestId, op, logger)
					return fmt.Errorf("rep for non-joker %v", jid)
				}
				break
			}
		}
		if !found {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("pickup rep for tile %v not in meld (discard + tileIds)", jid), requestId, op, logger)
			return fmt.Errorf("rep for non-existent tile %v", jid)
		}
		colour, ok := tile.ParseColour(rep.Colour)
		if !ok {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("pickup joker %v invalid colour %q", jid, rep.Colour), requestId, op, logger)
			return fmt.Errorf("invalid colour %q", rep.Colour)
		}
		rank := tile.Rank(rep.Rank)
		if !rank.IsValid() {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("pickup joker %v invalid rank %d", jid, rep.Rank), requestId, op, logger)
			return fmt.Errorf("invalid rank %d", rep.Rank)
		}
		repTile := tile.MustTile(tile.TileInstanceId(fmt.Sprintf("rep-%s", jid)), colour, rank)
		combinedReps[jid] = repTile
	}
	for _, tl := range combinedTiles {
		if tl.IsJoker {
			if _, ok := combinedReps[tl.ID]; !ok {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("pickup joker %v missing rep", tl.ID), requestId, op, logger)
				return fmt.Errorf("missing rep for joker %v", tl.ID)
			}
		}
	}

	// Determine kind and validate meld
	var kindToTry []meld.Kind
	if req.Kind != nil && *req.Kind != "" {
		if *req.Kind != string(meld.KindRun) && *req.Kind != string(meld.KindSet) {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("kind must be run or set, got %q", *req.Kind), requestId, op, logger)
			return fmt.Errorf("invalid kind")
		}
		kindToTry = []meld.Kind{meld.Kind(*req.Kind)}
	} else {
		kindToTry = []meld.Kind{meld.KindRun, meld.KindSet}
	}
	var validatedMeld meld.Meld
	var validatedKind meld.Kind
	foundValid := false
	var lastErr error
	for _, k := range kindToTry {
		mID := meld.MeldID(fmt.Sprintf("meld-pickup-%d-%d", req.DiscardIndex, int(seat)))
		// Ensure unique ID by appending if collision (unlikely)
		origID := mID
		collIdx := 0
		for {
			collision := false
			for _, em := range st.TableMelds {
				if meld.MeldID(em.ID) == mID {
					collision = true
					break
				}
			}
			if !collision {
				break
			}
			collIdx++
			mID = meld.MeldID(fmt.Sprintf("%s-%d", origID, collIdx))
		}
		m, err := meld.New(mID, k, combinedTiles, combinedReps)
		if err != nil {
			lastErr = err
			continue
		}
		// Also call specific validator for structured errors
		if k == meld.KindRun {
			if err := meld.ValidateRun(m); err != nil {
				lastErr = err
				continue
			}
		} else {
			if err := meld.ValidateSet(m); err != nil {
				lastErr = err
				continue
			}
		}
		validatedMeld = m
		validatedKind = k
		foundValid = true
		break
	}
	if !foundValid {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("pickup meld invalid: %v", lastErr), requestId, op, logger)
		return fmt.Errorf("meld invalid: %v", lastErr)
	}

	// Atomic mutation: remove 2 rack tiles, remove discard slice, add later discards to rack, add new meld, transition
	// Save originals for rollback
	origRack := append([]tile.TileInstance(nil), rack...)
	origDiscard := append([]DiscardEntry(nil), st.DiscardRow...)

	// 1. New rack = old rack minus 2 tiles + later discards
	laterTiles := []tile.TileInstance{}
	for i := req.DiscardIndex + 1; i < len(st.DiscardRow); i++ {
		laterTiles = append(laterTiles, st.DiscardRow[i].Tile)
	}
	newRack := make([]tile.TileInstance, 0, len(rack)-2+len(laterTiles))
	for _, t := range rack {
		if !seen[t.ID] {
			newRack = append(newRack, t)
		}
	}
	// Append later discards in order
	newRack = append(newRack, laterTiles...)
	st.Racks[seat] = newRack

	// 2. New discard row = up to discardIndex (exclusive)
	newDiscard := make([]DiscardEntry, req.DiscardIndex)
	copy(newDiscard, st.DiscardRow[:req.DiscardIndex])
	// Reindex
	for i := range newDiscard {
		newDiscard[i].Index = i
	}
	st.DiscardRow = newDiscard

	// 3. Add new meld
	tm := TableMeld{
		ID:        string(validatedMeld.ID),
		Kind:      string(validatedKind),
		Tiles:     validatedMeld.Tiles,
		JokerReps: validatedMeld.JokerReps,
		OwnerSeat: seat,
	}
	if err := tm.Validate(); err != nil {
		// Rollback
		st.Racks[seat] = origRack
		st.DiscardRow = origDiscard
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("table meld invalid: %v", err), requestId, op, logger)
		return err
	}
	st.TableMelds = append(st.TableMelds, tm)
	st.TurnPhase = TurnMeldOrDiscard

	logger.Info("PickupDiscardForMeld by %s seat %v discardIdx %d tile %v meld %q kind %v later %d newRack %d", senderId, seat, req.DiscardIndex, discardTile.ID, validatedMeld.ID, validatedKind, len(laterTiles), len(newRack))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"pickupDiscardForMeld","seat":%d,"discardIndex":%d,"meldId":"%s","later":%d}`, int(seat), req.DiscardIndex, validatedMeld.ID, len(laterTiles))), nil, nil, true)
	}
	return nil
}
