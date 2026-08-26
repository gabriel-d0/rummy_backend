// Day 25: Envelope and opcodes mirror Go internal/protocol/opcodes.go:8 Version 1, 1..9/100..199
export const Version = 1;
export const OpClientStart = 1;
export const OpClientDiscard = 2;
export const OpClientDrawStock = 3;
export const OpClientDrawPreviousDiscard = 4;
export const OpClientPickupDiscardForMeld = 5;
export const OpClientMeldInitial = 6;
export const OpClientMeldNew = 7;
export const OpClientExtendMeld = 8;
export const OpClientReplaceJoker = 9;
export const OpServerState = 100;
export const OpServerStatePublic = 101;
export const OpServerError = 102;
export const OpServerEvent = 103;

export interface Envelope {
  v: number;
  op: number;
  requestId?: string;
  payload?: unknown;
}

export function NewEnvelope(op: number, payload?: unknown, requestId?: string): string {
  return JSON.stringify({ v: Version, op, requestId, payload });
}

export function NewEnvelopeWithRequestId(op: number, requestId: string, payload?: unknown): string {
  return JSON.stringify({ v: Version, op, requestId, payload });
}
