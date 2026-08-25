// Package match — Extend meld handler (AGENTS.md Day 15).
// Implements EXTEND_MELD: opened players may add rack tiles to any public meld
// (own or other player's). The entire resulting meld is revalidated per
// docs/rules-decisions.md:1.3 and 5. Atomic, preserves joker rep immutability.
package match

import (
	"encoding/json"
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

type extendMeldPayload struct {
	MeldId    string                     `json:"meldId"`
	TileIDs   []string                   `json:"tileIds"`
	JokerReps map[string]jokerRepPayload `json:"jokerReps,omitempty"`
}

// handleExtendMeld validates and applies EXTEND_MELD. Atomic.
func handleExtendMeld(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	seat, pIdx, err := requireMeldOrDiscardTurn(st, senderId, requestId, op, dispatcher, sender, logger)
	if err != nil {
		return err
	}
	if !st.Players[pIdx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeNotOpened, "must open with MELD_INITIAL first (50+ with run)", requestId, op, logger)
		return fmt.Errorf("not opened")
	}

	var req extendMeldPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("extend payload bad JSON: %v", err), requestId, op, logger)
		return err
	}
	if req.MeldId == "" {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "meldId required", requestId, op, logger)
		return fmt.Errorf("meldId empty")
	}
	if len(req.TileIDs) == 0 {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "tileIds must have at least 1", requestId, op, logger)
		return fmt.Errorf("no tileIds")
	}

	// Find target meld
	targetIdx := -1
	for i, m := range st.TableMelds {
		if m.ID == req.MeldId {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld %q not found", req.MeldId), requestId, op, logger)
		return fmt.Errorf("meld not found")
	}
	target := st.TableMelds[targetIdx]

	rack := st.Racks[seat]
	rackByID := make(map[tile.TileInstanceId]tile.TileInstance, len(rack))
	for _, t := range rack {
		rackByID[t.ID] = t
	}

	// Check all tileIds owned and not duplicate
	seenNew := map[tile.TileInstanceId]bool{}
	newTiles := make([]tile.TileInstance, 0, len(req.TileIDs))
	for _, tidStr := range req.TileIDs {
		tid := tile.TileInstanceId(tidStr)
		if tid == "" {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q has empty tileId", req.MeldId), requestId, op, logger)
			return fmt.Errorf("empty tileId")
		}
		if seenNew[tid] {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("duplicate tile %v in extend payload", tid), requestId, op, logger)
			return fmt.Errorf("duplicate tile %v", tid)
		}
		// Ensure tile not already in target meld (no duplicate TileId across table)
		for _, et := range target.Tiles {
			if et.ID == tid {
				sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("tile %v already in meld %q", tid, req.MeldId), requestId, op, logger)
				return fmt.Errorf("tile already in meld")
			}
		}
		// Also ensure tile not in other melds (global conservation ensures not, but check rack ownership)
		t, ok := rackByID[tid]
		if !ok {
			sendError(dispatcher, sender, protocol.ErrCodeInvalidTile, fmt.Sprintf("tile %v not in rack", tid), requestId, op, logger)
			return fmt.Errorf("tile %v not in rack", tid)
		}
		seenNew[tid] = true
		newTiles = append(newTiles, t)
	}

	// Build combined JokerReps: copy existing, then add new reps for new jokers
	combinedReps := make(map[tile.TileInstanceId]tile.TileInstance, len(target.JokerReps)+len(req.JokerReps))
	for k, v := range target.JokerReps {
		combinedReps[k] = v
	}
	// Validate new joker reps
	for jidStr, rep := range req.JokerReps {
		jid := tile.TileInstanceId(jidStr)
		// Must be one of the new tiles and be a joker
		found := false
		for _, nt := range newTiles {
			if nt.ID == jid {
				found = true
				if !nt.IsJoker {
					sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("extend meld %q rep for non-joker %v", req.MeldId, jid), requestId, op, logger)
					return fmt.Errorf("rep for non-joker %v", jid)
				}
				break
			}
		}
		if !found {
			// Also check if jid is existing joker and client tries to override rep — that would be joker immutability violation
			if _, isExistingJoker := target.JokerReps[jid]; isExistingJoker {
				// Client is trying to change existing joker's rep — reject unless value matches existing exactly
				existingRep := target.JokerReps[jid]
				colour, ok := tile.ParseColour(rep.Colour)
				if !ok || colour != existingRep.Colour || tile.Rank(rep.Rank) != existingRep.Rank {
					sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("extend meld %q joker %v rep immutable (existing %v %v)", req.MeldId, jid, existingRep.Colour, existingRep.Rank), requestId, op, logger)
					return fmt.Errorf("joker rep immutable")
				}
				// If matches, no need to add (already in combinedReps)
				continue
			}
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("extend meld %q rep for tile %v not in tileIds and not existing joker", req.MeldId, jid), requestId, op, logger)
			return fmt.Errorf("rep for non-existent tile %v", jid)
		}
		colour, ok := tile.ParseColour(rep.Colour)
		if !ok {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("extend meld %q joker %v invalid colour %q", req.MeldId, jid, rep.Colour), requestId, op, logger)
			return fmt.Errorf("invalid colour %q", rep.Colour)
		}
		rank := tile.Rank(rep.Rank)
		if !rank.IsValid() {
			sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("extend meld %q joker %v invalid rank %d", req.MeldId, jid, rep.Rank), requestId, op, logger)
			return fmt.Errorf("invalid rank %d", rep.Rank)
		}
		repTile := tile.MustTile(tile.TileInstanceId(fmt.Sprintf("rep-%s", jid)), colour, rank)
		combinedReps[jid] = repTile
	}
	// Ensure every new joker has a rep
	for _, nt := range newTiles {
		if nt.IsJoker {
			if _, ok := combinedReps[nt.ID]; !ok {
				sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("extend meld %q new joker %v missing rep", req.MeldId, nt.ID), requestId, op, logger)
				return fmt.Errorf("missing rep for joker %v", nt.ID)
			}
		}
	}

	// Build combined tiles
	combinedTiles := make([]tile.TileInstance, 0, len(target.Tiles)+len(newTiles))
	combinedTiles = append(combinedTiles, target.Tiles...)
	combinedTiles = append(combinedTiles, newTiles...)

	// Determine kind: use target.Kind if set, else infer via validation
	kindStr := target.Kind
	if kindStr == "" {
		// Infer: try run, then set
		tmpRun := meld.Meld{ID: meld.MeldID(target.ID), Kind: meld.KindRun, Tiles: target.Tiles, JokerReps: target.JokerReps}
		if err := meld.ValidateRun(tmpRun); err == nil {
			kindStr = string(meld.KindRun)
		} else {
			tmpSet := meld.Meld{ID: meld.MeldID(target.ID), Kind: meld.KindSet, Tiles: target.Tiles, JokerReps: target.JokerReps}
			if err := meld.ValidateSet(tmpSet); err == nil {
				kindStr = string(meld.KindSet)
			} else {
				sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("existing meld %q has invalid kind (neither run nor set)", req.MeldId), requestId, op, logger)
				return fmt.Errorf("invalid existing meld kind")
			}
		}
	}
	kind := meld.Kind(kindStr)
	if kind != meld.KindRun && kind != meld.KindSet {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, fmt.Sprintf("meld %q has invalid kind %q", req.MeldId, kindStr), requestId, op, logger)
		return fmt.Errorf("invalid kind")
	}

	// Validate entire resulting meld
	newMeld, err := meld.New(meld.MeldID(target.ID), kind, combinedTiles, combinedReps)
	if err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("extend meld %q resulting meld invalid: %v", req.MeldId, err), requestId, op, logger)
		return err
	}
	// Additional validation via specific validator to ensure structured errors
	if kind == meld.KindRun {
		if err := meld.ValidateRun(newMeld); err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("extend meld %q run invalid: %v", req.MeldId, err), requestId, op, logger)
			return err
		}
	} else {
		if err := meld.ValidateSet(newMeld); err != nil {
			sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("extend meld %q set invalid: %v", req.MeldId, err), requestId, op, logger)
			return err
		}
	}

	// Atomic mutation: remove from rack, update table meld
	newRack := make([]tile.TileInstance, 0, len(rack)-len(newTiles))
	for _, t := range rack {
		if !seenNew[t.ID] {
			newRack = append(newRack, t)
		}
	}
	st.Racks[seat] = newRack
	updated := TableMeld{
		ID:        target.ID,
		Kind:      string(kind),
		Tiles:     newMeld.Tiles,
		JokerReps: newMeld.JokerReps,
		OwnerSeat: target.OwnerSeat, // keep original owner
	}
	if err := updated.Validate(); err != nil {
		st.Racks[seat] = rack
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("extended table meld invalid: %v", err), requestId, op, logger)
		return err
	}
	st.TableMelds[targetIdx] = updated

	logger.Info("ExtendMeld by %s seat %v meld %q tiles %d newRack %d", senderId, seat, req.MeldId, len(newTiles), len(newRack))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"extendMeld","seat":%d,"meldId":"%s","tiles":%d}`, int(seat), req.MeldId, len(newTiles))), nil, nil, true)
	}
	return nil
}
