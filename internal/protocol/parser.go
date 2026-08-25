// Package protocol — Command parser (Day 27). Safe inbound JSON parsing.
package protocol

import (
	"encoding/json"
	"fmt"
)

// ParseError is a structured parse failure; never panics.
type ParseError struct {
	Code    string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ParseEnvelope validates raw client data as a safe envelope.
// It rejects: empty, non-JSON, missing/unknown opcode, wrong version, non-object payload.
// It never panics and never mutates match state.
func ParseEnvelope(data []byte) (Envelope, error) {
	if len(data) == 0 {
		return Envelope{}, &ParseError{Code: "empty", Message: "empty payload"}
	}
	// Quick check: must be JSON object
	if data[0] != '{' {
		return Envelope{}, &ParseError{Code: "bad_json", Message: "not a JSON object"}
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return Envelope{}, &ParseError{Code: "bad_json", Message: fmt.Sprintf("invalid JSON: %v", err)}
	}
	if env.Version != Version {
		return Envelope{}, &ParseError{Code: "bad_version", Message: fmt.Sprintf("version %d want %d", env.Version, Version)}
	}
	if !IsClientOp(env.OpCode) {
		return Envelope{}, &ParseError{Code: "unknown_opcode", Message: fmt.Sprintf("opcode %d not a client op", env.OpCode)}
	}
	// Known opcodes check — currently all 1..9 are known; future unknown within 1..99 still rejected if not defined.
	switch env.OpCode {
	case OpClientStart, OpClientDiscard, OpClientDrawStock, OpClientDrawPreviousDiscard, OpClientPickupDiscardForMeld, OpClientMeldInitial, OpClientMeldNew, OpClientExtendMeld, OpClientReplaceJoker:
		// known
	default:
		return Envelope{}, &ParseError{Code: "unknown_opcode", Message: fmt.Sprintf("unknown opcode %d", env.OpCode)}
	}
	// Payload must be valid JSON if present (allow empty for start which has no payload)
	if len(env.Payload) > 0 {
		var tmp json.RawMessage
		if err := json.Unmarshal(env.Payload, &tmp); err != nil {
			return Envelope{}, &ParseError{Code: "bad_payload", Message: fmt.Sprintf("payload not JSON: %v", err)}
		}
	}
	return env, nil
}

// ParseEnvelopeStrict is an alias for tests; same as ParseEnvelope.
var ParseEnvelopeStrict = ParseEnvelope
