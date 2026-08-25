package match

import (
	"context"
	"database/sql"
	"testing"

	"github.com/heroiclabs/nakama-common/runtime"
)

// testLogger is a no-op logger satisfying runtime.Logger.
type testLogger struct{}

func (l *testLogger) Debug(format string, v ...interface{})                   {}
func (l *testLogger) Info(format string, v ...interface{})                    {}
func (l *testLogger) Warn(format string, v ...interface{})                    {}
func (l *testLogger) Error(format string, v ...interface{})                   {}
func (l *testLogger) WithField(key string, v interface{}) runtime.Logger      { return l }
func (l *testLogger) WithFields(fields map[string]interface{}) runtime.Logger { return l }
func (l *testLogger) Fields() map[string]interface{}                          { return nil }

// mockPresence implements runtime.Presence for tests.
type mockPresence struct {
	userId    string
	sessionId string
	username  string
	node      string
}

func (m *mockPresence) GetUserId() string                 { return m.userId }
func (m *mockPresence) GetSessionId() string              { return m.sessionId }
func (m *mockPresence) GetUsername() string               { return m.username }
func (m *mockPresence) GetNodeId() string                 { return m.node }
func (m *mockPresence) GetHidden() bool                   { return false }
func (m *mockPresence) GetPersistence() bool              { return false }
func (m *mockPresence) GetStatus() string                 { return "" }
func (m *mockPresence) GetReason() runtime.PresenceReason { return runtime.PresenceReasonUnknown }

func newPresence(userId string) *mockPresence {
	return &mockPresence{userId: userId, sessionId: "sess-" + userId, username: "user-" + userId, node: "node1"}
}

// mockDispatcher records label updates.
type mockDispatcher struct {
	lastLabel string
}

func (m *mockDispatcher) BroadcastMessage(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	return nil
}
func (m *mockDispatcher) BroadcastMessageDeferred(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	return nil
}
func (m *mockDispatcher) MatchKick(presences []runtime.Presence) error { return nil }
func (m *mockDispatcher) MatchLabelUpdate(label string) error {
	m.lastLabel = label
	return nil
}

func TestMatchInit(t *testing.T) {
	m := &RummyMatch{}
	state, tickRate, label := m.MatchInit(context.Background(), &testLogger{}, nil, nil, nil)
	if tickRate != 5 {
		t.Fatalf("tickRate %d want 5", tickRate)
	}
	if label != "rummy" {
		t.Fatalf("label %q want rummy", label)
	}
	st, ok := state.(*RoundState)
	if !ok {
		t.Fatalf("state %T want *RoundState", state)
	}
	if st.GamePhase != PhaseWaiting {
		t.Fatalf("phase %v want Waiting", st.GamePhase)
	}
	if len(st.Players) != 0 {
		t.Fatalf("players %d want 0", len(st.Players))
	}
}

func TestMatchJoinAttempt(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _, _ := m.MatchInit(context.Background(), logger, nil, nil, nil)

	// first join should be allowed
	_, ok, _ := m.MatchJoinAttempt(context.Background(), logger, nil, nil, nil, 0, state, newPresence("alice"), nil)
	if !ok {
		t.Fatalf("first join should be allowed")
	}

	// simulate 4 players already in state
	st := state.(*RoundState)
	st.Players = []PlayerState{{ID: "p1", Seat: 0}, {ID: "p2", Seat: 1}, {ID: "p3", Seat: 2}, {ID: "p4", Seat: 3}}
	_, ok, reason := m.MatchJoinAttempt(context.Background(), logger, nil, nil, nil, 0, st, newPresence("p5"), nil)
	if ok || reason != "match full" {
		t.Fatalf("full should reject, got ok %v reason %q", ok, reason)
	}

	// duplicate
	st.Players = []PlayerState{{ID: "alice", Seat: 0}}
	_, ok, reason = m.MatchJoinAttempt(context.Background(), logger, nil, nil, nil, 0, st, newPresence("alice"), nil)
	if ok || reason != "already joined" {
		t.Fatalf("duplicate should reject, got ok %v reason %q", ok, reason)
	}

	// game already started
	st.GamePhase = PhasePlaying
	_, ok, reason = m.MatchJoinAttempt(context.Background(), logger, nil, nil, nil, 0, st, newPresence("bob"), nil)
	if ok || reason != "game already started" {
		t.Fatalf("started game should reject, got ok %v reason %q", ok, reason)
	}
}

func TestMatchJoinAndLeave(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _, _ := m.MatchInit(context.Background(), logger, nil, nil, nil)
	dispatcher := &mockDispatcher{}
	// Join 3 players via MatchJoin (batch)
	pres := []runtime.Presence{newPresence("alice"), newPresence("bob"), newPresence("carol")}
	state = m.MatchJoin(context.Background(), logger, nil, nil, dispatcher, 0, state, pres)
	st := state.(*RoundState)
	if len(st.Players) != 3 {
		t.Fatalf("players %d want 3", len(st.Players))
	}
	if st.Players[0].ID != "alice" || st.Players[0].Seat != 0 {
		t.Fatalf("alice seat %v", st.Players[0].Seat)
	}
	if st.Players[2].ID != "carol" || st.Players[2].Seat != 2 {
		t.Fatalf("carol seat")
	}
	// Ensure Racks entries created
	if _, ok := st.Racks[0]; !ok {
		t.Fatalf("rack 0 missing")
	}
	// Duplicate join should be idempotent (no duplicate)
	state = m.MatchJoin(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.Presence{newPresence("alice")})
	if len(st.Players) != 3 {
		t.Fatalf("duplicate join created duplicate, len %d", len(st.Players))
	}
	// Leave bob
	state = m.MatchLeave(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.Presence{newPresence("bob")})
	st = state.(*RoundState)
	if len(st.Players) != 2 {
		t.Fatalf("after leave bob len %d", len(st.Players))
	}
	for _, p := range st.Players {
		if p.ID == "bob" {
			t.Fatalf("bob should be removed")
		}
	}
	if _, ok := st.Racks[1]; ok {
		t.Fatalf("rack 1 should be deleted after bob leaves")
	}
}

func TestMatchLoopAndTerminate(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _, _ := m.MatchInit(context.Background(), logger, nil, nil, nil)
	// Loop with no messages should return same state
	next := m.MatchLoop(context.Background(), logger, nil, nil, nil, 0, state, nil)
	if next != state {
		t.Fatalf("loop should return same state when no messages")
	}
	// Terminate
	next = m.MatchTerminate(context.Background(), logger, nil, nil, nil, 0, state, 0)
	if next != state {
		t.Fatalf("terminate should return same state")
	}
	// Signal
	next, out := m.MatchSignal(context.Background(), logger, nil, nil, nil, 0, state, "test")
	if next != state || out != "" {
		t.Fatalf("signal")
	}
	// Ensure dispatcher and logger nil safe (previous tests used nil dispatcher)
	_, _, _ = m.MatchJoinAttempt(context.Background(), logger, &sql.DB{}, nil, nil, 0, state, newPresence("x"), nil)
}
