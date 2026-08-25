package match

import "testing"

func TestAllowedOps(t *testing.T) {
	// Waiting: only start
	if !AllowedOps(PhaseWaiting, TurnMustDraw)[1] {
		t.Fatalf("Waiting should allow start")
	}
	if AllowedOps(PhaseWaiting, TurnMustDraw)[2] {
		t.Fatalf("Waiting should not allow discard")
	}
	// OpeningDiscard: only discard
	if !AllowedOps(PhaseOpeningDiscard, TurnMustDraw)[2] {
		t.Fatalf("OpeningDiscard should allow discard")
	}
	if AllowedOps(PhaseOpeningDiscard, TurnMustDraw)[3] {
		t.Fatalf("OpeningDiscard should not allow draw")
	}
	// Playing MustDraw: draw ops
	mustDraw := AllowedOps(PhasePlaying, TurnMustDraw)
	if !mustDraw[3] || !mustDraw[4] || !mustDraw[5] {
		t.Fatalf("MustDraw should allow draws, got %v", mustDraw)
	}
	if mustDraw[2] {
		t.Fatalf("MustDraw should not allow discard")
	}
	// Playing MeldOrDiscard: meld + discard
	meldOrDiscard := AllowedOps(PhasePlaying, TurnMeldOrDiscard)
	if !meldOrDiscard[2] || !meldOrDiscard[6] || !meldOrDiscard[7] {
		t.Fatalf("MeldOrDiscard should allow discard/meld, got %v", meldOrDiscard)
	}
	if meldOrDiscard[3] {
		t.Fatalf("MeldOrDiscard should not allow draw")
	}
	// RoundComplete: none
	if len(AllowedOps(PhaseRoundComplete, TurnMustDraw)) != 0 {
		t.Fatalf("RoundComplete should allow none")
	}
}

func TestValidatePhase(t *testing.T) {
	if err := ValidatePhase(PhaseWaiting, TurnMustDraw); err != nil {
		t.Fatalf("Waiting valid: %v", err)
	}
	if err := ValidatePhase(PhasePlaying, TurnMustDraw); err != nil {
		t.Fatalf("Playing MustDraw valid: %v", err)
	}
	if err := ValidatePhase(PhasePlaying, TurnPhase(99)); err == nil {
		t.Fatalf("invalid TurnPhase should fail")
	}
	if err := ValidatePhase(GamePhase(99), TurnMustDraw); err == nil {
		t.Fatalf("invalid GamePhase should fail")
	}
}
