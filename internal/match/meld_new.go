// Package match — Additional melds handler (AGENTS.md Day 14, refactor Day 23).
// Implements MELD_NEW: opened players may create one or more new independent
// sets/runs from rack, no score minimum, atomic, preserves TableMeld IDs per
// docs/rules-decisions.md:5.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/scoring"
	"github.com/heroiclabs/nakama-common/runtime"
)

// handleMeldNew validates and applies MELD_NEW. Atomic, requires HasOpened.
func handleMeldNew(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	seat, pIdx, err := requireMeldOrDiscardTurn(st, senderId, requestId, op, dispatcher, sender, logger)
	if err != nil {
		return err
	}
	if !st.Players[pIdx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeNotOpened, "must open with MELD_INITIAL first (50+ with run)", requestId, op, logger)
		return fmt.Errorf("not opened")
	}

	rack := st.Racks[seat]
	existingIDs := map[meld.MeldID]bool{}
	for _, m := range st.TableMelds {
		existingIDs[meld.MeldID(m.ID)] = true
	}

	scoringMelds, seenTile, err := buildMeldsFromPayload(payload, rack, existingIDs, seat, len(st.TableMelds), requestId, op, dispatcher, sender, logger)
	if err != nil {
		return err
	}

	batch := scoring.Batch{PlayerID: string(senderId), Melds: scoringMelds}
	if err := scoring.ValidateBatchOwnership(batch, rack); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("meld_new invalid: %v", err), requestId, op, logger)
		return err
	}

	if err := applyMeldBatch(st, seat, rack, scoringMelds, seenTile, requestId, op, dispatcher, sender, logger); err != nil {
		return err
	}
	logger.Info("MeldNew by %s seat %v melds %d newRack %d table %d", senderId, seat, len(scoringMelds), len(st.Racks[seat]), len(st.TableMelds))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"meldNew","seat":%d,"melds":%d}`, int(seat), len(scoringMelds))), nil, nil, true)
	}
	return nil
}
