// Package meld — Meld representation (Day 43).
// Canonical representation for table melds with stable IDs and explicit
// joker substitutions per docs/rules-decisions.md:1.3 and 3.
package meld

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Kind distinguishes run vs set.
type Kind string

const (
	KindRun Kind = "run"
	KindSet Kind = "set"
)

func (k Kind) IsValid() bool {
	return k == KindRun || k == KindSet
}

// MeldID is a stable UUID for a table meld. It never changes after creation.
type MeldID string

func (id MeldID) IsValid() bool {
	return id != ""
}

// Meld is the canonical table representation. JokerReps maps each joker
// TileInstanceId in Tiles to its declared represented tile (colour+rank, never joker).
// The mapping is immutable once melded unless a legal replacement occurs.
type Meld struct {
	ID        MeldID
	Kind      Kind
	Tiles     []tile.TileInstance
	JokerReps map[tile.TileInstanceId]tile.TileInstance
}

// Validate checks structural invariants (not full game-rule validation which is Day 44+).
func (m Meld) Validate() error {
	if !m.ID.IsValid() {
		return fmt.Errorf("meld ID empty")
	}
	if !m.Kind.IsValid() {
		return fmt.Errorf("meld %q invalid kind %q", m.ID, m.Kind)
	}
	if len(m.Tiles) < 3 {
		return fmt.Errorf("meld %q has %d tiles, need >=3", m.ID, len(m.Tiles))
	}
	seen := map[tile.TileInstanceId]bool{}
	tileByID := map[tile.TileInstanceId]tile.TileInstance{}
	for _, t := range m.Tiles {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("meld %q tile %v: %w", m.ID, t.ID, err)
		}
		if seen[t.ID] {
			return fmt.Errorf("meld %q duplicate tile %v", m.ID, t.ID)
		}
		seen[t.ID] = true
		tileByID[t.ID] = t
		if t.IsJoker {
			rep, ok := m.JokerReps[t.ID]
			if !ok {
				return fmt.Errorf("meld %q joker %v missing rep", m.ID, t.ID)
			}
			if rep.IsJoker {
				return fmt.Errorf("meld %q joker %v rep cannot be joker", m.ID, t.ID)
			}
			if !rep.Colour.IsValid() || !rep.Rank.IsValid() {
				return fmt.Errorf("meld %q joker %v rep invalid %v %v", m.ID, t.ID, rep.Colour, rep.Rank)
			}
		}
	}
	for jid := range m.JokerReps {
		t, ok := tileByID[jid]
		if !ok {
			return fmt.Errorf("meld %q rep for non-existent tile %v", m.ID, jid)
		}
		if !t.IsJoker {
			return fmt.Errorf("meld %q rep for non-joker %v", m.ID, jid)
		}
	}
	return nil
}

// New creates a meld with the given ID, kind, tiles, and joker reps (may be nil for no jokers).
func New(id MeldID, kind Kind, tiles []tile.TileInstance, reps map[tile.TileInstanceId]tile.TileInstance) (Meld, error) {
	if reps == nil {
		reps = map[tile.TileInstanceId]tile.TileInstance{}
	}
	// Copy tiles and reps to ensure immutability of caller slices
	tilesCopy := append([]tile.TileInstance(nil), tiles...)
	repsCopy := make(map[tile.TileInstanceId]tile.TileInstance, len(reps))
	for k, v := range reps {
		repsCopy[k] = v
	}
	m := Meld{ID: id, Kind: kind, Tiles: tilesCopy, JokerReps: repsCopy}
	if err := m.Validate(); err != nil {
		return Meld{}, err
	}
	return m, nil
}
