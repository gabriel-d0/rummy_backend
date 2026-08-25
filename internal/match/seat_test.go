package match

import "testing"

func TestSeatsForCount(t *testing.T) {
	for _, n := range []int{2, 3, 4} {
		seats, err := SeatsForCount(n)
		if err != nil {
			t.Fatalf("SeatsForCount(%d) err %v", n, err)
		}
		if len(seats) != n {
			t.Fatalf("SeatsForCount(%d) len %d want %d", n, len(seats), n)
		}
		for i, s := range seats {
			if int(s) != i {
				t.Fatalf("seat %d = %v want %d", i, s, i)
			}
			if !s.IsValid(n) {
				t.Fatalf("seat %v should be valid for n=%d", s, n)
			}
		}
	}
	if _, err := SeatsForCount(1); err == nil {
		t.Fatalf("expected error for n=1")
	}
	if _, err := SeatsForCount(5); err == nil {
		t.Fatalf("expected error for n=5")
	}
}

func TestAssignSeats(t *testing.T) {
	players, err := AssignSeats([]PlayerId{"alice", "bob", "carol"})
	if err != nil {
		t.Fatalf("AssignSeats err %v", err)
	}
	if len(players) != 3 {
		t.Fatalf("len %d", len(players))
	}
	if players[0].Seat != 0 || players[1].Seat != 1 || players[2].Seat != 2 {
		t.Fatalf("seats not deterministic: %+v", players)
	}
	if players[0].ID != "alice" || players[1].ID != "bob" {
		t.Fatalf("ID mismatch")
	}
	// duplicate ID
	if _, err := AssignSeats([]PlayerId{"alice", "alice"}); err == nil {
		t.Fatalf("expected duplicate error")
	}
	// empty ID
	if _, err := AssignSeats([]PlayerId{"", "bob"}); err == nil {
		t.Fatalf("expected empty ID error")
	}
}

func TestNextSeatAnticlockwise(t *testing.T) {
	// 2 players: 0→1→0
	if s, _ := NextSeat(0, 2); s != 1 {
		t.Fatalf("NextSeat 0,2 = %v want 1", s)
	}
	if s, _ := NextSeat(1, 2); s != 0 {
		t.Fatalf("NextSeat 1,2 = %v want 0", s)
	}
	// 3 players
	if s, _ := NextSeat(2, 3); s != 0 {
		t.Fatalf("NextSeat 2,3 = %v want 0", s)
	}
	if s, _ := NextSeat(1, 4); s != 2 {
		t.Fatalf("NextSeat 1,4 = %v want 2", s)
	}
	// loop covers all
	for n := 2; n <= 4; n++ {
		for seat := 0; seat < n; seat++ {
			next, err := NextSeat(Seat(seat), n)
			if err != nil {
				t.Fatalf("NextSeat %d %d err %v", seat, n, err)
			}
			prev, _ := PrevSeat(next, n)
			if prev != Seat(seat) {
				t.Fatalf("PrevSeat(NextSeat) failed n=%d seat %d", n, seat)
			}
		}
	}
	// invalid
	if _, err := NextSeat(SeatInvalid, 2); err == nil {
		t.Fatalf("expected error for invalid seat")
	}
	if _, err := NextSeat(5, 4); err == nil {
		t.Fatalf("expected error for out of range")
	}
}

func TestSeatOfPlayer(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"p1", "p2", "p3", "p4"})
	if SeatOfPlayer(players, "p1") != 0 {
		t.Fatalf("SeatOfPlayer p1")
	}
	if SeatOfPlayer(players, "p4") != 3 {
		t.Fatalf("SeatOfPlayer p4")
	}
	if SeatOfPlayer(players, "unknown") != SeatInvalid {
		t.Fatalf("unknown should be invalid")
	}
}

func TestValidatePlayers(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	if err := ValidatePlayers(players); err != nil {
		t.Fatalf("valid players err %v", err)
	}
	// duplicate ID
	dup := []PlayerState{{ID: "a", Seat: 0}, {ID: "a", Seat: 1}}
	if err := ValidatePlayers(dup); err == nil {
		t.Fatalf("expected duplicate ID error")
	}
	// duplicate seat
	dupSeat := []PlayerState{{ID: "a", Seat: 0}, {ID: "b", Seat: 0}}
	if err := ValidatePlayers(dupSeat); err == nil {
		t.Fatalf("expected duplicate seat error")
	}
	// invalid count
	if err := ValidatePlayers([]PlayerState{{ID: "a", Seat: 0}}); err == nil {
		t.Fatalf("expected count error")
	}
}
