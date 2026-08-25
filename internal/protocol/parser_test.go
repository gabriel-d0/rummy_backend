package protocol

import (
	"encoding/json"
	"testing"
)

func TestParseEnvelopeValid(t *testing.T) {
	for _, op := range []int64{OpClientStart, OpClientDiscard, OpClientDrawStock, OpClientMeldInitial} {
		env := Envelope{Version: Version, OpCode: op, Payload: json.RawMessage(`{}`)}
		b, _ := json.Marshal(env)
		got, err := ParseEnvelope(b)
		if err != nil {
			t.Fatalf("op %d should parse, err %v", op, err)
		}
		if got.OpCode != op || got.Version != Version {
			t.Fatalf("got %+v want op %d", got, op)
		}
	}
	// No payload also valid for start
	env := Envelope{Version: Version, OpCode: OpClientStart}
	b, _ := json.Marshal(env)
	if _, err := ParseEnvelope(b); err != nil {
		t.Fatalf("start no payload should parse: %v", err)
	}
}

func TestParseEnvelopeRejectsMalformed(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		code string
	}{
		{"empty", []byte{}, "empty"},
		{"not json", []byte("notjson"), "bad_json"},
		{"not object", []byte(`"string"`), "bad_json"},
		{"bad json", []byte(`{op:1`), "bad_json"},
		{"wrong version", MustEnvelopeWithVersion(99, OpClientStart, nil), "bad_version"},
		{"unknown opcode server", MustEnvelopeWithVersion(Version, OpServerState, nil), "unknown_opcode"},
		{"unknown client 99 not defined", MustEnvelopeWithVersion(Version, 99, nil), "unknown_opcode"},
		{"bad payload outer invalid", []byte(`{"v":1,"op":1,"payload":notjson}`), "bad_json"},
	}
	for _, tc := range cases {
		_, err := ParseEnvelope(tc.data)
		if err == nil {
			t.Fatalf("%s should fail", tc.name)
		}
		pe, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("%s err type %T", tc.name, err)
		}
		if pe.Code != tc.code {
			t.Fatalf("%s code %q want %q err %v", tc.name, pe.Code, tc.code, err)
		}
	}
}

func MustEnvelopeWithVersion(v int64, op int64, payload interface{}) []byte {
	env := Envelope{Version: v, OpCode: op}
	if payload != nil {
		b, _ := json.Marshal(payload)
		env.Payload = b
	} else {
		env.Payload = json.RawMessage(`{}`)
		if op == OpClientStart {
			env.Payload = nil
		}
	}
	b, _ := json.Marshal(env)
	return b
}

func TestParseEnvelopeDoesNotPanicOnRandomBytes(t *testing.T) {
	// Fuzz-like: random bytes should not panic
	for _, data := range [][]byte{
		[]byte{0, 1, 2, 3},
		[]byte(`{"v":1,"op":1,"payload":`),
		[]byte(`{"v":1,"op":9999999999999999999999}`),
	} {
		_, err := ParseEnvelope(data)
		if err == nil {
			// some may be valid? but should not panic
		}
	}
}
