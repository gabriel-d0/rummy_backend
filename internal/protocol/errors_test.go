package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	err := NewError(ErrCodeBadPayload, "tileId required", "req-123", map[string]string{"field": "tileId"})
	if err.Code != ErrCodeBadPayload || err.Message != "tileId required" || err.RequestId != "req-123" {
		t.Fatalf("NewError fields %+v", err)
	}
	if err.Details["field"] != "tileId" {
		t.Fatalf("details")
	}
	// Encode
	b := EncodeError(err)
	var decoded ErrorResponse
	if err2 := json.Unmarshal(b, &decoded); err2 != nil {
		t.Fatalf("unmarshal %v", err2)
	}
	if decoded.Code != ErrCodeBadPayload || decoded.RequestId != "req-123" {
		t.Fatalf("decoded %+v", decoded)
	}
}

func TestNewErrorForEnvelope(t *testing.T) {
	env := Envelope{Version: Version, OpCode: OpClientDiscard, RequestId: "req-456"}
	err := NewErrorForEnvelope(env, ErrCodeUnknownOpcode, "unknown", nil)
	if err.RequestId != "req-456" || err.OpCode != OpClientDiscard {
		t.Fatalf("for envelope %+v", err)
	}
	if err.Code != ErrCodeUnknownOpcode {
		t.Fatalf("code")
	}
}

func TestErrorResponseJSONOmitEmpty(t *testing.T) {
	err := NewError(ErrCodeBadRequest, "bad", "", nil)
	b := EncodeError(err)
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if _, ok := m["requestId"]; ok {
		t.Fatalf("requestId should be omitted when empty")
	}
	if _, ok := m["details"]; ok {
		t.Fatalf("details should be omitted when nil")
	}
	if m["code"] != ErrCodeBadRequest {
		t.Fatalf("code")
	}
}

func TestEnvelopeRequestIdRoundTrip(t *testing.T) {
	b, _ := NewEnvelopeWithRequestId(OpClientStart, "req-789", map[string]string{"a": "b"})
	env, err := ParseEnvelope(b)
	if err != nil {
		t.Fatalf("ParseEnvelope with requestId: %v", err)
	}
	if env.RequestId != "req-789" {
		t.Fatalf("RequestId %q want req-789", env.RequestId)
	}
	// Error correlation
	perr := NewErrorForEnvelope(env, ErrCodeBadPayload, "payload bad", nil)
	if perr.RequestId != "req-789" {
		t.Fatalf("error should echo requestId")
	}
}

func TestErrorCodesStable(t *testing.T) {
	// Ensure codes are not empty and stable
	codes := []string{ErrCodeBadRequest, ErrCodeBadJSON, ErrCodeBadVersion, ErrCodeUnknownOpcode, ErrCodeBadPayload, ErrCodeNotMember, ErrCodeNotYourTurn, ErrCodeWrongPhase}
	seen := map[string]bool{}
	for _, c := range codes {
		if c == "" {
			t.Fatalf("empty code")
		}
		if seen[c] {
			t.Fatalf("duplicate code %q", c)
		}
		seen[c] = true
	}
}
