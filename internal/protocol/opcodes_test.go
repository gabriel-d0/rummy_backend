package protocol

import "testing"

func TestOpcodesStableAndUnique(t *testing.T) {
	seen := map[int64]string{}
	check := func(op int64, name string) {
		if op < 1 || op > 199 {
			t.Fatalf("opcode %s %d out of range 1..199", name, op)
		}
		if prev, dup := seen[op]; dup {
			t.Fatalf("duplicate opcode %d: %s and %s", op, prev, name)
		}
		seen[op] = name
	}
	check(OpClientStart, "OpClientStart")
	check(OpClientDiscard, "OpClientDiscard")
	check(OpClientDrawStock, "OpClientDrawStock")
	check(OpClientDrawPreviousDiscard, "OpClientDrawPreviousDiscard")
	check(OpClientPickupDiscardForMeld, "OpClientPickupDiscardForMeld")
	check(OpClientMeldInitial, "OpClientMeldInitial")
	check(OpClientMeldNew, "OpClientMeldNew")
	check(OpClientExtendMeld, "OpClientExtendMeld")
	check(OpClientReplaceJoker, "OpClientReplaceJoker")
	check(OpServerState, "OpServerState")
	check(OpServerStatePublic, "OpServerStatePublic")
	check(OpServerError, "OpServerError")
	check(OpServerEvent, "OpServerEvent")

	// Client range
	for op := range seen {
		if op >= 1 && op <= 99 {
			if !IsClientOp(op) {
				t.Fatalf("IsClientOp false for %d", op)
			}
		}
		if op >= 100 && op <= 199 {
			if !IsServerOp(op) {
				t.Fatalf("IsServerOp false for %d", op)
			}
		}
	}
	if Version != 1 {
		t.Fatalf("Version %d want 1", Version)
	}
	// Ensure start still 1 (stable, Day 24 used magic 1)
	if OpClientStart != 1 {
		t.Fatalf("OpClientStart must stay 1 for Day 24 compatibility, got %d", OpClientStart)
	}
}
