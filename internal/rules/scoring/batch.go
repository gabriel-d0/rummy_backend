// Package scoring — Opening meld batch model (Day 60).
package scoring

import "github.com/gabriel-d0/rummy_backend/internal/rules/meld"

// Batch is the proposed initial melds from one player's rack.
// It is the payload for OpClientMeldInitial per docs/rules-decisions.md:4.
type Batch struct {
	PlayerID string      `json:"playerId"`
	Melds    []meld.Meld `json:"melds"`
}

// TotalScore sums ScoreRun/ScoreSet for each meld in the batch.
// It is pure and validates each meld is a run or set before scoring;
// it does NOT yet check 50-point, run requirement, or tile ownership — those are Day 61-64.
func TotalScore(batch Batch) (int, error) {
	total := 0
	for _, m := range batch.Melds {
		var s int
		var err error
		switch m.Kind {
		case meld.KindRun:
			s, err = ScoreRun(m)
		case meld.KindSet:
			s, err = ScoreSet(m)
		default:
			return 0, err
		}
		if err != nil {
			return 0, err
		}
		total += s
	}
	return total, nil
}
