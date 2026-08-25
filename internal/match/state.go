// Package match — Game state model (Day 12).
// Defines server-side round state structures per docs/terminology.md and docs/rules-decisions.md:5.
// Pure data, no I/O; invariant checker lands Day 13.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// GamePhase is the top-level round state.
type GamePhase int

const (
	PhaseWaiting        GamePhase = iota // lobby, 2–4 players not yet started
	PhaseOpeningDiscard                  // opening player must discard from 15
	PhasePlaying                         // normal MustDraw ↔ MeldOrDiscard loop
	PhaseRoundComplete                   // winner or dead round
)

func (p GamePhase) String() string {
	switch p {
	case PhaseWaiting:
		return "Waiting"
	case PhaseOpeningDiscard:
		return "OpeningDiscard"
	case PhasePlaying:
		return "Playing"
	case PhaseRoundComplete:
		return "RoundComplete"
	default:
		return fmt.Sprintf("GamePhase(%d)", int(p))
	}
}

func (p GamePhase) IsValid() bool {
	return p >= PhaseWaiting && p <= PhaseRoundComplete
}

// TurnPhase is the phase within a Playing turn.
type TurnPhase int

const (
	TurnMustDraw      TurnPhase = iota // at turn start — must draw
	TurnMeldOrDiscard                  // after draw — may meld/extend/replace, must discard
)

func (t TurnPhase) String() string {
	switch t {
	case TurnMustDraw:
		return "MustDraw"
	case TurnMeldOrDiscard:
		return "MeldOrDiscard"
	default:
		return fmt.Sprintf("TurnPhase(%d)", int(t))
	}
}

func (t TurnPhase) IsValid() bool {
	return t == TurnMustDraw || t == TurnMeldOrDiscard
}

// DiscardEntry is one public discard in order. IsOpeningDiscard marks the
// very first discard (opening player's 15→14) as permanently unavailable per
// docs/rules-decisions.md:1.5.
type DiscardEntry struct {
	Tile             tile.TileInstance
	IsOpeningDiscard bool
	Index            int // 0-based position in DiscardRow, for tests
}

func (d DiscardEntry) Validate() error {
	if err := d.Tile.Validate(); err != nil {
		return fmt.Errorf("discard %d tile invalid: %w", d.Index, err)
	}
	return nil
}

// TableMeld is a public meld on the table. Joker representations are explicit
// and immutable once placed — see docs/rules-decisions.md:1.3.
type TableMeld struct {
	ID        string                                    // stable meld UUID
	Tiles     []tile.TileInstance                       // ordered tiles as placed
	JokerReps map[tile.TileInstanceId]tile.TileInstance // joker ID → represented tile (colour+rank)
	OwnerSeat Seat                                      // seat who created it (for info, extensions allowed from any opened player)
}

func (m TableMeld) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("meld ID empty")
	}
	if len(m.Tiles) < 3 {
		return fmt.Errorf("meld %q has %d tiles, need >=3", m.ID, len(m.Tiles))
	}
	// Tile validation + duplicate check
	seen := map[tile.TileInstanceId]bool{}
	for _, t := range m.Tiles {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("meld %q tile %v invalid: %w", m.ID, t.ID, err)
		}
		if seen[t.ID] {
			return fmt.Errorf("meld %q duplicate tile %v", m.ID, t.ID)
		}
		seen[t.ID] = true
		if t.IsJoker {
			if _, ok := m.JokerReps[t.ID]; !ok {
				return fmt.Errorf("meld %q joker %v missing representation", m.ID, t.ID)
			}
			rep := m.JokerReps[t.ID]
			if rep.IsJoker {
				return fmt.Errorf("meld %q joker %v rep cannot be joker", m.ID, t.ID)
			}
			if !rep.Colour.IsValid() || !rep.Rank.IsValid() {
				return fmt.Errorf("meld %q joker %v rep invalid", m.ID, t.ID)
			}
		}
	}
	for jid := range m.JokerReps {
		if !seen[jid] {
			return fmt.Errorf("meld %q rep for non-joker %v", m.ID, jid)
		}
	}
	return nil
}

// RoundState is the authoritative server-side round state per docs/terminology.md.
type RoundState struct {
	Players     []PlayerState                // 2–4, deterministic seat order
	Racks       map[Seat][]tile.TileInstance // private, indexed by Seat
	Stock       []tile.TileInstance          // remaining tiles, top = last element
	DiscardRow  []DiscardEntry               // ordered, Index tracks position
	TableMelds  []TableMeld                  // public melds
	CurrentSeat Seat                         // active seat, validated vs Players
	GamePhase   GamePhase
	TurnPhase   TurnPhase // only meaningful when GamePhase==PhasePlaying; ignored otherwise
	Winner      Seat      // SeatInvalid if none yet; set when PhaseRoundComplete
}

// Rack returns the rack for a seat (nil if not found, not an error — caller checks seat validity first).
func (s *RoundState) Rack(seat Seat) []tile.TileInstance {
	return s.Racks[seat]
}

// Validate performs structural checks (counts, phases, seat validity, duplicate tiles across locations
// is deferred to Day 13 invariant, but we check local duplicates here).
func (s *RoundState) Validate() error {
	if err := ValidatePlayers(s.Players); err != nil {
		return fmt.Errorf("players: %w", err)
	}
	n := len(s.Players)
	if !s.CurrentSeat.IsValid(n) && s.GamePhase != PhaseWaiting && s.GamePhase != PhaseRoundComplete {
		return fmt.Errorf("current seat %v invalid for n=%d phase %v", s.CurrentSeat, n, s.GamePhase)
	}
	if !s.GamePhase.IsValid() {
		return fmt.Errorf("invalid GamePhase %v", s.GamePhase)
	}
	if s.GamePhase == PhasePlaying && !s.TurnPhase.IsValid() {
		return fmt.Errorf("invalid TurnPhase %v for Playing", s.TurnPhase)
	}
	if s.Winner != SeatInvalid && !s.Winner.IsValid(n) {
		return fmt.Errorf("winner %v invalid for n=%d", s.Winner, n)
	}
	// Racks: each seat must exist, 13–15 tiles allowed at construction (opening 15, others 14), but Validate does not enforce exact deal sizes — Day 13 will.
	for _, p := range s.Players {
		rack := s.Racks[p.Seat]
		if rack == nil {
			return fmt.Errorf("missing rack for seat %v", p.Seat)
		}
		seen := map[tile.TileInstanceId]bool{}
		for _, t := range rack {
			if err := t.Validate(); err != nil {
				return fmt.Errorf("rack seat %v tile %v: %w", p.Seat, t.ID, err)
			}
			if seen[t.ID] {
				return fmt.Errorf("rack seat %v duplicate tile %v", p.Seat, t.ID)
			}
			seen[t.ID] = true
		}
	}
	// Stock
	for _, t := range s.Stock {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("stock tile %v: %w", t.ID, err)
		}
	}
	// DiscardRow ordered and opening flag only on first entry
	for i, d := range s.DiscardRow {
		if err := d.Validate(); err != nil {
			return err
		}
		if d.Index != i {
			return fmt.Errorf("discard index %d mismatch expected %d", d.Index, i)
		}
		if i > 0 && d.IsOpeningDiscard {
			return fmt.Errorf("only discard 0 may be opening discard, found at %d", i)
		}
	}
	// TableMelds
	for _, m := range s.TableMelds {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}
