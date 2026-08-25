// Package match — Initial meld handler (AGENTS.md Day 13, refactor Day 23).
// Implements MELD_INITIAL: batch of melds from player's rack, >=50 points with a run,
// atomic, marks HasOpened, stays in MeldOrDiscard per docs/rules-decisions.md:1.4 and 5.
// Pure validation delegated to scoring.ValidateInitialBatch; parsing delegated to meld_common.go.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/scoring"
	"github.com/heroiclabs/nakama-common/runtime"
)

// handleMeldInitial validates and applies MELD_INITIAL. It is atomic: on any
// validation error state is unchanged and an OpServerError is sent.
func handleMeldInitial(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	seat, pIdx, err := requireMeldOrDiscardTurn(st, senderId, requestId, op, dispatcher, sender, logger)
	if err != nil {
		return err
	}
	if st.Players[pIdx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeAlreadyOpened, "already opened: use MELD_NEW", requestId, op, logger)
		return fmt.Errorf("already opened")
	}

	rack := st.Racks[seat]
	scoringMelds, seenTile, err := buildMeldsFromPayload(payload, rack, nil, seat, len(st.TableMelds), requestId, op, dispatcher, sender, logger)
	if err != nil {
		return err
	}

	batch := scoring.Batch{PlayerID: string(senderId), Melds: scoringMelds}
	if err := scoring.ValidateInitialBatch(batch, rack); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("initial batch invalid: %v", err), requestId, op, logger)
		return err
	}

	if err := applyMeldBatch(st, seat, rack, scoringMelds, seenTile, requestId, op, dispatcher, sender, logger); err != nil {
		return err
	}
	st.Players[pIdx].HasOpened = true

	logger.Info("MeldInitial by %s seat %v melds %d newRack %d table %d", senderId, seat, len(scoringMelds), len(st.Racks[seat]), len(st.TableMelds))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"meldInitial","seat":%d,"melds":%d}`, int(seat), len(scoringMelds))), nil, nil, true)
	}
	return nil
}
