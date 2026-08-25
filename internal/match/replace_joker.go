// Package match — Joker replacement handler (AGENTS.md Day 18).
// Implements REPLACE_JOKER: opened player in MeldOrDiscard may replace a joker
// in a run (exact tile) or set (missing colour) and must immediately form a new
// valid meld with the recovered joker plus two rack tiles. Atomic.
package match

import (
	"encoding/json"
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

type replaceJokerPayload struct {
	TargetMeldId string                     `json:"targetMeldId"`
	TileId       string                     `json:"tileId"`       // real tile replacing joker
	NewMeldTiles []string                   `json:"newMeldTiles"` // exactly 2 rack tiles for new meld
	JokerReps    map[string]jokerRepPayload `json:"jokerReps,omitempty"`
	NewMeldKind  *string                    `json:"newMeldKind,omitempty"` // optional "run" or "set"
}

func handleReplaceJoker(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	seat, pIdx, err := requireMeldOrDiscardTurn(st, senderId, requestId, op, dispatcher, sender, logger)
	if err != nil {
		return err
	}
	if !st.Players[pIdx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeNotOpened, "must open with MELD_INITIAL first (50+ with run)", requestId, op, logger)
		return fmt.Errorf("not opened")
	}

	var req replaceJokerPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("replace payload bad JSON: %v", err), requestId, op, logger)
		return err
	}
	if req.TargetMeldId == "" || req.TileId == "" || len(req.NewMeldTiles) != 2 {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "targetMeldId, tileId and newMeldTiles[2] required", requestId, op, logger)
		return fmt.Errorf("bad payload")
	}

	// Find target meld
	targetIdx := -1
	for i, m := range st.TableMelds {
		if m.ID == req.TargetMeldId {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld %q not found", req.TargetMeldId), requestId, op, logger)
		return fmt.Errorf("meld not found")
	}
	target := st.TableMelds[targetIdx]

	// Must contain at least one joker
	hasJoker := false
	for _, tl := range target.Tiles {
		if tl.IsJoker {
			hasJoker = true
			break
		}
	}
	if !hasJoker {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld %q has no joker to replace", req.TargetMeldId), requestId, op, logger)
		return fmt.Errorf("no joker")
	}

	rack := st.Racks[seat]
	rackByID := make(map[tile.TileInstanceId]tile.TileInstance, len(rack))
	for _, t := range rack {
		rackByID[t.ID] = t
	}

	// Validate tileId owned and not duplicate with newMeldTiles
	seen := map[tile.TileInstanceId]bool{}
	tileId := tile.TileInstanceId(req.TileId)
	if tileId == "" {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "tileId empty", requestId, op, logger)
		return fmt.Errorf("empty tileId")
	}
	replTile, ok := rackByID[tileId]
	if !ok {
		sendError(dispatcher, sender, protocol.ErrCodeInvalidTile, fmt.Sprintf("tile %v not in rack", tileId), requestId, op, logger)
		return fmt.Errorf("tile not in rack")
	}
	if replTile.IsJoker {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "tileId must be real tile, not joker", requestId, op, logger)
		return fmt.Errorf("tileId is joker")
	}
	seen[tileId] = true

	newTiles := make([]tile.TileInstance, 0, 2)
	for _, tidStr := range req.NewMeldTiles {
		tid := tile.TileInstanceId(tidStr)
		if tid == "" {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "newMeldTiles empty", requestId, op, logger)
			return fmt.Errorf("empty newMeldTile")
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
		newTiles = append(newTiles, t)
	}
	// Ensure tileId not duplicate with newMeldTiles (already via seen)

	// Find joker in target whose rep matches replTile (exact colour/rank)
	var jokerToReplace tile.TileInstance
	var jokerRep tile.TileInstance
	jokerFound := false
	for _, tl := range target.Tiles {
		if tl.IsJoker {
			rep, ok := target.JokerReps[tl.ID]
			if !ok {
				continue
			}
			if rep.Colour == replTile.Colour && rep.Rank == replTile.Rank {
				jokerToReplace = tl
				jokerRep = rep
				jokerFound = true
				break
			}
		}
	}
	if !jokerFound {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("no joker in meld %q represents tile %v %v", req.TargetMeldId, replTile.Colour, replTile.Rank), requestId, op, logger)
		return fmt.Errorf("no matching joker")
	}
	_ = jokerRep

	// Build updated target meld: replace joker tile with replTile, remove jokerId from JokerReps
	updatedTiles := make([]tile.TileInstance, 0, len(target.Tiles))
	for _, tl := range target.Tiles {
		if tl.ID == jokerToReplace.ID {
			updatedTiles = append(updatedTiles, replTile)
		} else {
			updatedTiles = append(updatedTiles, tl)
		}
	}
	updatedReps := map[tile.TileInstanceId]tile.TileInstance{}
	for k, v := range target.JokerReps {
		if k != jokerToReplace.ID {
			updatedReps[k] = v
		}
	}
	// Determine target kind (use stored Kind or infer)
	targetKindStr := target.Kind
	if targetKindStr == "" {
		tmpRun := meld.Meld{ID: meld.MeldID(target.ID), Kind: meld.KindRun, Tiles: target.Tiles, JokerReps: target.JokerReps}
		if err := meld.ValidateRun(tmpRun); err == nil {
			targetKindStr = string(meld.KindRun)
		} else {
			tmpSet := meld.Meld{ID: meld.MeldID(target.ID), Kind: meld.KindSet, Tiles: target.Tiles, JokerReps: target.JokerReps}
			if err := meld.ValidateSet(tmpSet); err == nil {
				targetKindStr = string(meld.KindSet)
			} else {
				sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("existing meld %q invalid", req.TargetMeldId), requestId, op, logger)
				return fmt.Errorf("invalid existing meld")
			}
		}
	}
	targetKind := meld.Kind(targetKindStr)
	updatedMeld, err := meld.New(meld.MeldID(target.ID), targetKind, updatedTiles, updatedReps)
	if err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("replaced meld invalid: %v", err), requestId, op, logger)
		return err
	}
	if targetKind == meld.KindRun {
		if err := meld.ValidateRun(updatedMeld); err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("replaced run invalid: %v", err), requestId, op, logger)
			return err
		}
	} else {
		if err := meld.ValidateSet(updatedMeld); err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("replaced set invalid: %v", err), requestId, op, logger)
			return err
		}
	}

	// Build new meld with recovered joker + 2 new tiles
	// Recovered joker tile instance is jokerToReplace (original joker)
	recoveredJoker := jokerToReplace
	// Need rep for recovered joker in new meld context
	newMeldTiles := append([]tile.TileInstance{recoveredJoker}, newTiles...)
	newMeldReps := map[tile.TileInstanceId]tile.TileInstance{}
	// Check payload jokerReps for recovered joker
	if req.JokerReps != nil {
		for jidStr, rep := range req.JokerReps {
			jid := tile.TileInstanceId(jidStr)
			if jid != recoveredJoker.ID {
				// Only recovered joker should be in reps for this new meld (since newTiles are real tiles per test, but could also be jokers? For MVP new meld has exactly 1 joker (recovered) +2 real, so no other jokers)
				// If payload contains rep for other joker not in new meld, ignore or error
				// For strictness, require jid == recovered joker
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("new meld jokerRep %v must be recovered joker %v", jid, recoveredJoker.ID), requestId, op, logger)
				return fmt.Errorf("joker rep mismatch")
			}
			colour, ok := tile.ParseColour(rep.Colour)
			if !ok {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("new meld joker %v invalid colour %q", jid, rep.Colour), requestId, op, logger)
				return fmt.Errorf("invalid colour")
			}
			rank := tile.Rank(rep.Rank)
			if !rank.IsValid() {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("new meld joker %v invalid rank %d", jid, rep.Rank), requestId, op, logger)
				return fmt.Errorf("invalid rank")
			}
			newMeldReps[jid] = tile.MustTile(tile.TileInstanceId(fmt.Sprintf("rep-%s", jid)), colour, rank)
		}
	}
	// Ensure recovered joker has rep
	if _, ok := newMeldReps[recoveredJoker.ID]; !ok {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("new meld missing rep for recovered joker %v", recoveredJoker.ID), requestId, op, logger)
		return fmt.Errorf("missing rep for recovered joker")
	}
	// Also ensure newTiles that are jokers (should not happen, since newTiles are from rack and rack could contain jokers? But new meld requires 2 rack tiles, they could be jokers, but ratio would still hold? For MVP new meld has 1 joker +2 real, ratio ok. If newTiles contains joker, they'd need reps too, but our payload only allows rep for recovered joker. For simplicity, require newTiles are not jokers for now.
	for _, nt := range newTiles {
		if nt.IsJoker {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "newMeldTiles must be real tiles, not jokers (MVP)", requestId, op, logger)
			return fmt.Errorf("new tiles must be real")
		}
	}

	// Determine new meld kind (try run then set, or use payload kind if provided)
	var newKindToTry []meld.Kind
	if req.NewMeldKind != nil && *req.NewMeldKind != "" {
		if *req.NewMeldKind != string(meld.KindRun) && *req.NewMeldKind != string(meld.KindSet) {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("newMeldKind must be run or set, got %q", *req.NewMeldKind), requestId, op, logger)
			return fmt.Errorf("invalid kind")
		}
		newKindToTry = []meld.Kind{meld.Kind(*req.NewMeldKind)}
	} else {
		newKindToTry = []meld.Kind{meld.KindRun, meld.KindSet}
	}
	var newMeld meld.Meld
	foundValid := false
	var lastErr error
	for _, k := range newKindToTry {
		mID := meld.MeldID(fmt.Sprintf("meld-replace-new-%d", len(st.TableMelds)))
		// Ensure unique
		collIdx := 0
		origID := mID
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
		m, err := meld.New(mID, k, newMeldTiles, newMeldReps)
		if err != nil {
			lastErr = err
			continue
		}
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
		newMeld = m
		foundValid = true
		break
	}
	if !foundValid {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("new joker meld invalid: %v", lastErr), requestId, op, logger)
		return fmt.Errorf("new meld invalid: %v", lastErr)
	}

	// Atomic mutation
	origRack := append([]tile.TileInstance(nil), rack...)
	origTable := append([]TableMeld(nil), st.TableMelds...)
	// New rack = old rack minus 3 tiles (tileId + 2 new)
	newRack := make([]tile.TileInstance, 0, len(rack)-3)
	for _, t := range rack {
		if t.ID == tileId || t.ID == tile.TileInstanceId(req.NewMeldTiles[0]) || t.ID == tile.TileInstanceId(req.NewMeldTiles[1]) {
			continue
		}
		newRack = append(newRack, t)
	}
	st.Racks[seat] = newRack
	// Update target meld
	updatedTM := TableMeld{
		ID:        string(updatedMeld.ID),
		Kind:      string(updatedMeld.Kind),
		Tiles:     updatedMeld.Tiles,
		JokerReps: updatedMeld.JokerReps,
		OwnerSeat: target.OwnerSeat,
	}
	if err := updatedTM.Validate(); err != nil {
		st.Racks[seat] = origRack
		st.TableMelds = origTable
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("updated meld invalid: %v", err), requestId, op, logger)
		return err
	}
	st.TableMelds[targetIdx] = updatedTM
	// Add new meld
	newTM := TableMeld{
		ID:        string(newMeld.ID),
		Kind:      string(newMeld.Kind),
		Tiles:     newMeld.Tiles,
		JokerReps: newMeld.JokerReps,
		OwnerSeat: seat,
	}
	if err := newTM.Validate(); err != nil {
		st.Racks[seat] = origRack
		st.TableMelds = origTable
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("new meld invalid: %v", err), requestId, op, logger)
		return err
	}
	st.TableMelds = append(st.TableMelds, newTM)

	logger.Info("ReplaceJoker by %s seat %v target %q repl %v newMeld %q kind %v", senderId, seat, req.TargetMeldId, tileId, newMeld.ID, newMeld.Kind)
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"replaceJoker","seat":%d,"target":"%s","newMeld":"%s"}`, int(seat), req.TargetMeldId, newMeld.ID)), nil, nil, true)
	}
	return nil
}
