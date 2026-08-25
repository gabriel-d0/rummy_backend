// Package scoring — Run scoring (Day 57).
package scoring

import (
	"sort"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ScoreRun returns the total points for a run meld.
// It validates the meld is a run (via meld.ValidateRun) and sums ScoreTile
// with correct Ace low (1-2-3 =>5) vs high (12-13-1 =>10) context.
// Joker value is its represented tile's value in the run's Ace context.
func ScoreRun(m meld.Meld) (int, error) {
	if err := meld.ValidateRun(m); err != nil {
		return 0, err
	}
	// Collect effective ranks to determine low vs high Ace
	effectiveRanks := make([]int, 0, len(m.Tiles))
	for _, t := range m.Tiles {
		var r tile.Rank
		if t.IsJoker {
			rep := m.JokerReps[t.ID]
			r = rep.Rank
		} else {
			r = t.Rank
		}
		effectiveRanks = append(effectiveRanks, int(r))
	}
	sort.Ints(effectiveRanks)
	isLow := isLowAceRun(effectiveRanks)
	isHigh := isHighAceRun(effectiveRanks)
	// Score each tile
	total := 0
	for _, t := range m.Tiles {
		var rep *tile.TileInstance
		if t.IsJoker {
			r := m.JokerReps[t.ID]
			rep = &r
		}
		total += ScoreTile(t, isLow, isHigh, false, rep)
	}
	return total, nil
}

func isLowAceRun(ranks []int) bool {
	// Sorted ranks, low Ace means consecutive starting at 1 (1,2,3...)
	if len(ranks) == 0 || ranks[0] != 1 {
		return false
	}
	for i := 1; i < len(ranks); i++ {
		if ranks[i] != ranks[i-1]+1 {
			return false
		}
	}
	return true
}

func isHighAceRun(ranks []int) bool {
	// Contains Ace (1) and when Ace treated as 14, ranks are consecutive ending at 14
	if !containsRank(ranks, 1) {
		return false
	}
	high := make([]int, len(ranks))
	for i, r := range ranks {
		if r == 1 {
			high[i] = 14
		} else {
			high[i] = r
		}
	}
	sort.Ints(high)
	for i := 1; i < len(high); i++ {
		if high[i] != high[i-1]+1 {
			return false
		}
	}
	// Must contain 13 and be consecutive ending at 14, and not also be low (low would be 1,2,3)
	// If it was low, isLowAceRun already true, but high should not also be true for low
	// For 1,2,3 high would be 2,3,14 not consecutive, so false — correct
	return true
}

func containsRank(ranks []int, v int) bool {
	for _, r := range ranks {
		if r == v {
			return true
		}
	}
	return false
}
