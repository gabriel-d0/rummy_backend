// Package protocol — Command schema validation (Day 28).
// Validates per-opcode payloads after ParseEnvelope per AGENTS.md:176.
package protocol

import (
	"encoding/json"
	"fmt"
)

// ValidatePayload checks that payload matches the expected schema for op.
// It is pure and never mutates state. On failure returns *ParseError with Code bad_payload.
func ValidatePayload(op int64, payload json.RawMessage) error {
	// Empty payload handling: some ops allow empty
	isEmpty := len(payload) == 0 || string(payload) == "null" || string(payload) == "{}" || string(payload) == `""`
	switch op {
	case OpClientStart:
		// Start allows empty or {} — no required fields
		if isEmpty {
			return nil
		}
		var tmp map[string]interface{}
		if err := json.Unmarshal(payload, &tmp); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("start payload must be object: %v", err)}
		}
		return nil

	case OpClientDiscard:
		if len(payload) == 0 || string(payload) == "null" {
			return &ParseError{Code: "bad_payload", Message: "discard payload missing: need {tileId}"}
		}
		var p struct {
			TileId *string `json:"tileId"`
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("discard payload bad JSON: %v", err)}
		}
		if p.TileId == nil || *p.TileId == "" {
			return &ParseError{Code: "bad_payload", Message: "discard.tileId required non-empty string"}
		}
		return nil

	case OpClientDrawStock, OpClientDrawPreviousDiscard:
		// No payload required; allow empty or {}
		if isEmpty {
			return nil
		}
		var tmp map[string]interface{}
		if err := json.Unmarshal(payload, &tmp); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("payload must be empty object: %v", err)}
		}
		if len(tmp) != 0 {
			return &ParseError{Code: "bad_payload", Message: "payload must be empty for this op"}
		}
		return nil

	case OpClientPickupDiscardForMeld:
		var p struct {
			DiscardIndex *int     `json:"discardIndex"`
			TileIds      []string `json:"tileIds"`
			JokerReps    map[string]struct {
				Colour string `json:"colour"`
				Rank   int    `json:"rank"`
			} `json:"jokerReps,omitempty"`
			Kind *string `json:"kind,omitempty"` // optional "run" or "set", if omitted server tries both
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("pickup payload bad JSON: %v", err)}
		}
		if p.DiscardIndex == nil {
			return &ParseError{Code: "bad_payload", Message: "pickup.discardIndex required"}
		}
		if *p.DiscardIndex < 0 {
			return &ParseError{Code: "bad_payload", Message: "pickup.discardIndex must be >=0"}
		}
		if len(p.TileIds) != 2 {
			return &ParseError{Code: "bad_payload", Message: "pickup.tileIds must have exactly 2 entries"}
		}
		for i, id := range p.TileIds {
			if id == "" {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("pickup.tileIds[%d] empty", i)}
			}
		}
		if p.Kind != nil && *p.Kind != "" && *p.Kind != "run" && *p.Kind != "set" {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("pickup.kind must be run or set, got %q", *p.Kind)}
		}
		for jid, rep := range p.JokerReps {
			if jid == "" {
				return &ParseError{Code: "bad_payload", Message: "pickup.jokerReps key empty"}
			}
			if rep.Colour == "" {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("pickup.jokerReps[%q] colour required", jid)}
			}
			if rep.Rank < 1 || rep.Rank > 13 {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("pickup.jokerReps[%q] rank must be 1..13", jid)}
			}
		}
		return nil

	case OpClientMeldInitial, OpClientMeldNew:
		var p struct {
			Melds []json.RawMessage `json:"melds"`
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("meld payload bad JSON: %v", err)}
		}
		if len(p.Melds) == 0 {
			return &ParseError{Code: "bad_payload", Message: "melds must have at least 1 meld"}
		}
		return nil

	case OpClientExtendMeld:
		var p struct {
			MeldId    *string  `json:"meldId"`
			TileIds   []string `json:"tileIds"`
			JokerReps map[string]struct {
				Colour string `json:"colour"`
				Rank   int    `json:"rank"`
			} `json:"jokerReps,omitempty"`
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("extend payload bad JSON: %v", err)}
		}
		if p.MeldId == nil || *p.MeldId == "" {
			return &ParseError{Code: "bad_payload", Message: "extend.meldId required"}
		}
		if len(p.TileIds) == 0 {
			return &ParseError{Code: "bad_payload", Message: "extend.tileIds must have at least 1"}
		}
		for jid, rep := range p.JokerReps {
			if jid == "" {
				return &ParseError{Code: "bad_payload", Message: "extend.jokerReps key empty"}
			}
			if rep.Colour == "" {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("extend.jokerReps[%q] colour required", jid)}
			}
			if rep.Rank < 1 || rep.Rank > 13 {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("extend.jokerReps[%q] rank must be 1..13", jid)}
			}
		}
		return nil

	case OpClientReplaceJoker:
		var p struct {
			TargetMeldId *string  `json:"targetMeldId"`
			TileId       *string  `json:"tileId"`
			NewMeldTiles []string `json:"newMeldTiles"`
			JokerReps    map[string]struct {
				Colour string `json:"colour"`
				Rank   int    `json:"rank"`
			} `json:"jokerReps,omitempty"`
			NewMeldKind *string `json:"newMeldKind,omitempty"`
		}
		if err := json.Unmarshal(payload, &p); err != nil {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("replace payload bad JSON: %v", err)}
		}
		if p.TargetMeldId == nil || *p.TargetMeldId == "" {
			return &ParseError{Code: "bad_payload", Message: "replace.targetMeldId required"}
		}
		if p.TileId == nil || *p.TileId == "" {
			return &ParseError{Code: "bad_payload", Message: "replace.tileId required"}
		}
		if len(p.NewMeldTiles) != 2 {
			return &ParseError{Code: "bad_payload", Message: "replace.newMeldTiles must have exactly 2"}
		}
		for i, id := range p.NewMeldTiles {
			if id == "" {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("replace.newMeldTiles[%d] empty", i)}
			}
		}
		if p.NewMeldKind != nil && *p.NewMeldKind != "" && *p.NewMeldKind != "run" && *p.NewMeldKind != "set" {
			return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("replace.newMeldKind must be run or set, got %q", *p.NewMeldKind)}
		}
		for jid, rep := range p.JokerReps {
			if jid == "" {
				return &ParseError{Code: "bad_payload", Message: "replace.jokerReps key empty"}
			}
			if rep.Colour == "" {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("replace.jokerReps[%q] colour required", jid)}
			}
			if rep.Rank < 1 || rep.Rank > 13 {
				return &ParseError{Code: "bad_payload", Message: fmt.Sprintf("replace.jokerReps[%q] rank must be 1..13", jid)}
			}
		}
		return nil

	default:
		return &ParseError{Code: "unknown_opcode", Message: fmt.Sprintf("unknown opcode %d", op)}
	}
}

// ValidateEnvelope is a convenience that parses and validates payload in one call.
func ValidateEnvelope(data []byte) (Envelope, error) {
	env, err := ParseEnvelope(data)
	if err != nil {
		return Envelope{}, err
	}
	if err := ValidatePayload(env.OpCode, env.Payload); err != nil {
		return Envelope{}, err
	}
	return env, nil
}
