package setup

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/match"
)

// TestSnapshotRedactionExhaustive verifies that public and per-player private
// snapshots never reveal another player's rack tile IDs, for 2/3/4 players
// and multiple seeds using the real 106-tile deck.
func TestSnapshotRedactionExhaustive(t *testing.T) {
	seeds := []int64{42, 7, 123}
	for _, n := range []int{2, 3, 4} {
		for _, seed := range seeds {
			names := []match.PlayerId{"alice", "bob", "carol", "dana"}
			ids := names[:n]
			state, _, err := NewRoundState(ids, seed)
			if err != nil {
				t.Fatalf("n=%d seed=%d NewRoundState: %v", n, seed, err)
			}
			// Collect all rack tile IDs by seat
			rackIDsBySeat := map[match.Seat]map[string]bool{}
			for _, p := range state.Players {
				m := map[string]bool{}
				for _, tl := range state.Racks[p.Seat] {
					m[string(tl.ID)] = true
				}
				rackIDsBySeat[p.Seat] = m
			}

			// Public view must not contain any rack tile ID
			pub := match.PublicView(state)
			pubJSON, _ := json.Marshal(pub)
			pubStr := string(pubJSON)
			for seat, m := range rackIDsBySeat {
				for id := range m {
					if strings.Contains(pubStr, id) {
						t.Fatalf("n=%d seed=%d public leaked seat %v tile %v", n, seed, seat, id)
					}
				}
			}
			// Each private view must contain own and not others
			for _, p := range state.Players {
				priv := match.PrivateView(state, p.Seat)
				privJSON, _ := json.Marshal(priv)
				privStr := string(privJSON)
				ownIDs := rackIDsBySeat[p.Seat]
				for id := range ownIDs {
					if !strings.Contains(privStr, id) {
						t.Fatalf("n=%d seed=%d seat %v private missing own %v", n, seed, p.Seat, id)
					}
				}
				for otherSeat, m := range rackIDsBySeat {
					if otherSeat == p.Seat {
						continue
					}
					for id := range m {
						if strings.Contains(privStr, id) {
							t.Fatalf("n=%d seed=%d seat %v private leaked other seat %v tile %v", n, seed, p.Seat, otherSeat, id)
						}
					}
				}
				// Stock IDs must also not leak (public only has count)
				for _, tl := range state.Stock {
					if strings.Contains(privStr, string(tl.ID)) {
						t.Fatalf("n=%d seed=%d seat %v private leaked stock %v", n, seed, p.Seat, tl.ID)
					}
				}
				if !strings.Contains(pubStr, `"stockCount"`) {
					t.Fatalf("public missing stockCount")
				}
			}
		}
	}
}
