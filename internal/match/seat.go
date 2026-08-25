// Package match defines player identity, seating, and anticlockwise turn order.
// Day 11 — pure deterministic seating per docs/terminology.md and docs/rules-decisions.md:1.2.
package match

import "fmt"

// PlayerId is a Nakama userId string.
type PlayerId string

func (id PlayerId) IsValid() bool {
	return id != ""
}

func (id PlayerId) String() string {
	return string(id)
}

// Seat is deterministic position 0..n-1 by join order.
// SeatInvalid is used for unseated or error cases.
type Seat int

const (
	SeatInvalid Seat = -1
	SeatMin     Seat = 0
)

func (s Seat) IsValid(n int) bool {
	return n >= 2 && n <= 4 && s >= 0 && int(s) < n
}

func (s Seat) String() string {
	if s == SeatInvalid {
		return "seat-invalid"
	}
	return fmt.Sprintf("seat-%d", int(s))
}

// PlayerState holds per-player seating and open flag.
// Rack and other state will be added Day 12; here we keep minimal for Day 11.
type PlayerState struct {
	ID        PlayerId
	Seat      Seat
	HasOpened bool
}

// Validate checks ID non-empty and seat in range for n players.
func (p PlayerState) Validate(n int) error {
	if !p.ID.IsValid() {
		return fmt.Errorf("player seat %v has empty ID", p.Seat)
	}
	if !p.Seat.IsValid(n) {
		return fmt.Errorf("player %v seat %v invalid for n=%d", p.ID, p.Seat, n)
	}
	return nil
}

// SeatsForCount returns deterministic seats 0..n-1 for n=2..4.
func SeatsForCount(n int) ([]Seat, error) {
	if n < 2 || n > 4 {
		return nil, fmt.Errorf("player count %d must be 2..4", n)
	}
	seats := make([]Seat, n)
	for i := 0; i < n; i++ {
		seats[i] = Seat(i)
	}
	return seats, nil
}

// AssignSeats assigns seats to playerIds in join order.
// Deterministic: first joiner = Seat 0 (opening player per docs/rules-decisions.md:2).
func AssignSeats(playerIds []PlayerId) ([]PlayerState, error) {
	n := len(playerIds)
	if n < 2 || n > 4 {
		return nil, fmt.Errorf("player count %d must be 2..4", n)
	}
	seen := map[PlayerId]bool{}
	players := make([]PlayerState, n)
	for i, id := range playerIds {
		if !id.IsValid() {
			return nil, fmt.Errorf("player %d has empty ID", i)
		}
		if seen[id] {
			return nil, fmt.Errorf("duplicate player ID %v", id)
		}
		seen[id] = true
		players[i] = PlayerState{ID: id, Seat: Seat(i), HasOpened: false}
	}
	return players, nil
}

// NextSeat returns the next seat anticlockwise per docs/rules-decisions.md:1.2.
// Anticlockwise is defined as (current+1)%n for deterministic server order.
func NextSeat(current Seat, n int) (Seat, error) {
	if !current.IsValid(n) {
		return SeatInvalid, fmt.Errorf("current seat %v invalid for n=%d", current, n)
	}
	return Seat((int(current) + 1) % n), nil
}

// PrevSeat returns the previous seat (clockwise) — useful for validation/tests.
func PrevSeat(current Seat, n int) (Seat, error) {
	if !current.IsValid(n) {
		return SeatInvalid, fmt.Errorf("current seat %v invalid for n=%d", current, n)
	}
	return Seat((int(current) + n - 1) % n), nil
}

// SeatOfPlayer returns the seat for a given playerId, or SeatInvalid if not found.
func SeatOfPlayer(players []PlayerState, id PlayerId) Seat {
	for _, p := range players {
		if p.ID == id {
			return p.Seat
		}
	}
	return SeatInvalid
}

// ValidatePlayers checks count 2..4, no duplicate IDs, seats 0..n-1, and unique seats.
func ValidatePlayers(players []PlayerState) error {
	n := len(players)
	if n < 2 || n > 4 {
		return fmt.Errorf("player count %d must be 2..4", n)
	}
	seenID := map[PlayerId]bool{}
	seenSeat := map[Seat]bool{}
	for _, p := range players {
		if err := p.Validate(n); err != nil {
			return err
		}
		if seenID[p.ID] {
			return fmt.Errorf("duplicate player ID %v", p.ID)
		}
		seenID[p.ID] = true
		if seenSeat[p.Seat] {
			return fmt.Errorf("duplicate seat %v", p.Seat)
		}
		seenSeat[p.Seat] = true
	}
	return nil
}
