// Package protocol — Envelope (Day 27). Stable JSON envelope for client→server.
package protocol

import "encoding/json"

// Envelope is the wire format for client→server match messages.
// All messages are JSON objects with version, opcode, optional payload, and
// optional requestId for correlation per AGENTS.md:368. Version is stable 1.
type Envelope struct {
	Version   int64           `json:"v"`
	OpCode    int64           `json:"op"`
	RequestId string          `json:"requestId,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
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

// NewEnvelopeWithRequestId creates an envelope with requestId for correlation.
func NewEnvelopeWithRequestId(op int64, requestId string, payload interface{}) ([]byte, error) {
	pBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	env := Envelope{Version: Version, OpCode: op, RequestId: requestId, Payload: pBytes}
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
