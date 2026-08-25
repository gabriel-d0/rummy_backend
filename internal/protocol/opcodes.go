// Package protocol defines stable client/server opcodes and envelope versioning.
// Day 26 — Protocol opcodes. All opcodes are stable and must not be reused.
// Client opcodes are 1..99, server opcodes 100..199 per AGENTS.md:152-157.
// JSON messages use OpCode field; Version is envelope version.
package protocol

// Version is the protocol envelope version. Bump only on breaking changes.
const Version = 1

// Client → Server opcodes (authoritative match loop messages).
const (
	// OpClientStart requests to start the waiting room (host Seat 0, 2..4 players).
	// Payload: {} . Temporary alias for Day 24 start via MatchLoop op 1.
	OpClientStart int64 = 1

	// Reserved for turn actions (implemented Day  7..18, opcodes stable now)
	OpClientDiscard              int64 = 2 // DISCARD (opening and normal)
	OpClientDrawStock            int64 = 3 // DRAW_STOCK
	OpClientDrawPreviousDiscard  int64 = 4 // DRAW_PREVIOUS_DISCARD
	OpClientPickupDiscardForMeld int64 = 5 // PICKUP_DISCARD_FOR_MELD
	OpClientMeldInitial          int64 = 6 // MELD_INITIAL
	OpClientMeldNew              int64 = 7 // MELD_NEW
	OpClientExtendMeld           int64 = 8 // EXTEND_MELD
	OpClientReplaceJoker         int64 = 9 // REPLACE_JOKER
)

// Server → Client opcodes (dispatcher.BroadcastMessage).
const (
	OpServerState       int64 = 100 // Full public + private snapshot (Day 30)
	OpServerStatePublic int64 = 101 // Public snapshot broadcast
	OpServerError       int64 = 102 // Standard error (Day 29)
	OpServerEvent       int64 = 103 // Generic event (e.g. started, discarded)
)

// IsClientOp reports whether op is a client→server opcode.
func IsClientOp(op int64) bool {
	return op >= 1 && op <= 99
}

// IsServerOp reports whether op is a server→client opcode.
func IsServerOp(op int64) bool {
	return op >= 100 && op <= 199
}
