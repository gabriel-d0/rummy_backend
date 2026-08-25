// Package scoring — Set scoring (Day 58).
package scoring

import (
	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ScoreSet returns total points for a set meld.
// It validates via meld.ValidateSet and sums ScoreTile with isAceSet context
// for the special 25-point Ace set (three Aces, len 3, all rank 1).
// Joker value is its represented tile's value in set context.
func ScoreSet(m meld.Meld) (int, error) {
	if err := meld.ValidateSet(m); err != nil {
		return 0, err
	}
	// Determine if this is an Ace set: all effective ranks are 1 and len 3 or 4 (per spec "set of three" but 4 Aces also 25)
	isAceSet := false
	if len(m.Tiles) == 3 || len(m.Tiles) == 4 {
		allAce := true
		for _, t := range m.Tiles {
			var r tile.Rank
			if t.IsJoker {
				rep := m.JokerReps[t.ID]
				r = rep.Rank
			} else {
				r = t.Rank
			}
			if r != tile.RankAce {
				allAce = false
				break
			}
		}
		if allAce {
			isAceSet = true
		}
	}
	total := 0
	for _, t := range m.Tiles {
		var rep *tile.TileInstance
		if t.IsJoker {
			r := m.JokerReps[t.ID]
			rep = &r
		}
		total += ScoreTile(t, false, false, isAceSet, rep)
	}
	return total, nil
}
