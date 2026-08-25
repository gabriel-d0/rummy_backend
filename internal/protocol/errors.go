// Package protocol — Standard error protocol (Day 29).
// Error responses are server→client OpServerError 102 with code, message,
// details, and requestId correlation per AGENTS.md:368.
package protocol

import (
	"encoding/json"
	"fmt"
)

// ErrorResponse is the payload for OpServerError (102).
type ErrorResponse struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestId string            `json:"requestId,omitempty"`
	OpCode    int64             `json:"op,omitempty"`
}

// Error implements error interface.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a server error response.
func NewError(code, message string, requestId string, details map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestId: requestId,
	}
}

// NewErrorForEnvelope creates an error response correlated to a client envelope.
func NewErrorForEnvelope(env Envelope, code, message string, details map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestId: env.RequestId,
		OpCode:    env.OpCode,
	}
}

// EncodeError marshals an ErrorResponse as JSON for BroadcastMessage payload.
func EncodeError(err *ErrorResponse) []byte {
	b, _ := json.Marshal(err)
	return b
}

// Common error codes (stable, per AGENTS.md:370).
const (
	ErrCodeBadRequest    = "bad_request"
	ErrCodeBadJSON       = "bad_json"
	ErrCodeBadVersion    = "bad_version"
	ErrCodeUnknownOpcode = "unknown_opcode"
	ErrCodeBadPayload    = "bad_payload"
	ErrCodeNotMember     = "not_member"
	ErrCodeNotYourTurn   = "not_your_turn"
	ErrCodeWrongPhase    = "wrong_phase"
	ErrCodeInvalidTile   = "invalid_tile"
	ErrCodeNotOpened     = "not_opened"
	ErrCodeAlreadyOpened = "already_opened"
)
