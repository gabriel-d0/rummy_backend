// Package protocol — Envelope (Day 27). Stable JSON envelope for client→server.
package protocol

import "encoding/json"

// Envelope is the wire format for client→server match messages.
// All messages are JSON objects with version, opcode, and optional payload.
// Server validates envelope before any state mutation per AGENTS.md:176.
type Envelope struct {
	Version int64           `json:"v"`
	OpCode  int64           `json:"op"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// NewEnvelope creates an envelope for tests/clients.
func NewEnvelope(op int64, payload interface{}) ([]byte, error) {
	pBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	env := Envelope{Version: Version, OpCode: op, Payload: pBytes}
	return json.Marshal(env)
}

// MustEnvelope panics on error — for tests.
func MustEnvelope(op int64, payload interface{}) []byte {
	b, err := NewEnvelope(op, payload)
	if err != nil {
		panic(err)
	}
	return b
}
