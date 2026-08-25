package protocol

import (
	"encoding/json"
	"testing"
)

func mustPayload(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func TestValidatePayloadStart(t *testing.T) {
	if err := ValidatePayload(OpClientStart, nil); err != nil {
		t.Fatalf("start nil: %v", err)
	}
	if err := ValidatePayload(OpClientStart, json.RawMessage(`{}`)); err != nil {
		t.Fatalf("start {}: %v", err)
	}
	if err := ValidatePayload(OpClientStart, mustPayload(map[string]interface{}{"extra": 1})); err != nil {
		t.Fatalf("start extra fields should still pass (object): %v", err)
	}
}

func TestValidatePayloadDiscard(t *testing.T) {
	ok := mustPayload(map[string]string{"tileId": "t-123"})
	if err := ValidatePayload(OpClientDiscard, ok); err != nil {
		t.Fatalf("discard ok: %v", err)
	}
	for _, bad := range []json.RawMessage{
		nil,
		json.RawMessage(`{}`),
		mustPayload(map[string]string{"tileId": ""}),
		mustPayload(map[string]string{"wrong": "x"}),
		json.RawMessage(`notjson`),
	} {
		if err := ValidatePayload(OpClientDiscard, bad); err == nil {
			t.Fatalf("discard bad %s should fail", string(bad))
		}
	}
}

func TestValidatePayloadDrawStock(t *testing.T) {
	if err := ValidatePayload(OpClientDrawStock, nil); err != nil {
		t.Fatalf("drawStock nil: %v", err)
	}
	if err := ValidatePayload(OpClientDrawStock, json.RawMessage(`{}`)); err != nil {
		t.Fatalf("drawStock {}: %v", err)
	}
	if err := ValidatePayload(OpClientDrawStock, mustPayload(map[string]int{"extra": 1})); err == nil {
		t.Fatalf("drawStock with extra should fail")
	}
}

func TestValidatePayloadPickup(t *testing.T) {
	ok := mustPayload(map[string]interface{}{"discardIndex": 2, "tileIds": []string{"t1", "t2"}})
	if err := ValidatePayload(OpClientPickupDiscardForMeld, ok); err != nil {
		t.Fatalf("pickup ok: %v", err)
	}
	badCases := []json.RawMessage{
		mustPayload(map[string]interface{}{"tileIds": []string{"t1", "t2"}}),                          // missing index
		mustPayload(map[string]interface{}{"discardIndex": -1, "tileIds": []string{"t1", "t2"}}),      // negative
		mustPayload(map[string]interface{}{"discardIndex": 0, "tileIds": []string{"t1"}}),             // only 1
		mustPayload(map[string]interface{}{"discardIndex": 0, "tileIds": []string{"t1", ""}}),         // empty id
		mustPayload(map[string]interface{}{"discardIndex": 0, "tileIds": []string{"t1", "t2", "t3"}}), // 3
	}
	for _, bad := range badCases {
		if err := ValidatePayload(OpClientPickupDiscardForMeld, bad); err == nil {
			t.Fatalf("pickup bad %s should fail", string(bad))
		}
	}
}

func TestValidatePayloadMeld(t *testing.T) {
	ok := mustPayload(map[string]interface{}{"melds": []interface{}{map[string]interface{}{"tiles": []string{"t1", "t2", "t3"}}}})
	if err := ValidatePayload(OpClientMeldInitial, ok); err != nil {
		t.Fatalf("meld ok: %v", err)
	}
	if err := ValidatePayload(OpClientMeldInitial, mustPayload(map[string]interface{}{"melds": []interface{}{}})); err == nil {
		t.Fatalf("empty melds should fail")
	}
	if err := ValidatePayload(OpClientMeldInitial, mustPayload(map[string]interface{}{})); err == nil {
		t.Fatalf("missing melds should fail")
	}
}

func TestValidatePayloadExtendAndReplace(t *testing.T) {
	if err := ValidatePayload(OpClientExtendMeld, mustPayload(map[string]interface{}{"meldId": "m1", "tileIds": []string{"t1"}})); err != nil {
		t.Fatalf("extend ok: %v", err)
	}
	if err := ValidatePayload(OpClientExtendMeld, mustPayload(map[string]interface{}{"meldId": "", "tileIds": []string{"t1"}})); err == nil {
		t.Fatalf("extend empty meldId should fail")
	}
	if err := ValidatePayload(OpClientExtendMeld, mustPayload(map[string]interface{}{"meldId": "m1", "tileIds": []string{}})); err == nil {
		t.Fatalf("extend empty tileIds should fail")
	}
	if err := ValidatePayload(OpClientReplaceJoker, mustPayload(map[string]interface{}{"targetMeldId": "m1", "tileId": "t1", "newMeldTiles": []string{"a", "b"}})); err != nil {
		t.Fatalf("replace ok: %v", err)
	}
	if err := ValidatePayload(OpClientReplaceJoker, mustPayload(map[string]interface{}{"targetMeldId": "m1", "tileId": "t1", "newMeldTiles": []string{"a"}})); err == nil {
		t.Fatalf("replace 1 tile should fail")
	}
}

func TestValidateEnvelopeIntegration(t *testing.T) {
	envBytes, _ := NewEnvelope(OpClientDiscard, map[string]string{"tileId": "t1"})
	if _, err := ValidateEnvelope(envBytes); err != nil {
		t.Fatalf("ValidateEnvelope discard ok: %v", err)
	}
	// bad payload via envelope
	envBytes, _ = NewEnvelope(OpClientDiscard, map[string]string{"wrong": "x"})
	if _, err := ValidateEnvelope(envBytes); err == nil {
		t.Fatalf("ValidateEnvelope should fail bad payload")
	}
	// bad opcode via envelope
	badEnv := Envelope{Version: Version, OpCode: 99, Payload: json.RawMessage(`{}`)}
	b, _ := json.Marshal(badEnv)
	if _, err := ValidateEnvelope(b); err == nil {
		t.Fatalf("unknown opcode should fail")
	}
}
