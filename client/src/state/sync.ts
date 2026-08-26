import type { PrivateSnapshot, PublicSnapshot } from "./snapshot";
import { showErrorToast, type ServerError } from "../ui/ErrorToast";

// Day 27: Receive match state — parses Envelope and routes op 100/101/102/103 to handlers
// Mirrors Go internal/protocol/opcodes.go:8 Version 1 and internal/match/visibility.go:36
// Day 28: OpServerError 102 is routed to ErrorToast (red toast 3s)

let lastPrivateSnapshot: PrivateSnapshot | null = null;
let lastPublicSnapshot: PublicSnapshot | null = null;

export function onPrivateSnapshot(snap: PrivateSnapshot): void {
  lastPrivateSnapshot = snap;
  console.log(
    `received op 100 PrivateSnapshot v=${snap.v} gamePhase=${snap.gamePhase} currentSeat=${snap.currentSeat} ownSeat=${snap.ownSeat} rack=${snap.ownRack.length} — Day 27`
  );
  // In a real scene, this would trigger RackScene re-render and TableScene re-render
  // For Day 27, we just log and store for reconnection
  localStorage.setItem("rummy_lastPrivate", JSON.stringify(snap));
}

export function onPublicSnapshot(snap: PublicSnapshot): void {
  lastPublicSnapshot = snap;
  console.log(`received op 101 PublicSnapshot v=${snap.v} gamePhase=${snap.gamePhase} — Day 27`);
}

export function onServerError(error: ServerError): void {
  // Store and toast via ErrorToast (Day 28)
  try {
    showErrorToast(error);
  } catch {
    console.log(
      `received op 102 OpServerError code=${error.code} message=${error.message} requestId=${error.requestId} op=${error.op} — Day 28 (toast failed)`
    );
  }
}

export function onServerEvent(event: unknown): void {
  console.log(`received op 103 OpServerEvent`, event, "— Day 27");
}

export function getLastPrivateSnapshot(): PrivateSnapshot | null {
  return lastPrivateSnapshot;
}

export function getLastPublicSnapshot(): PublicSnapshot | null {
  return lastPublicSnapshot;
}

export function handleMatchData(opCode: number, data: Uint8Array | string): void {
  // Data is Envelope JSON string or Uint8Array
  let jsonStr: string;
  if (data instanceof Uint8Array) {
    jsonStr = new TextDecoder().decode(data);
  } else {
    jsonStr = data as string;
  }
  try {
    const envelope = JSON.parse(jsonStr);
    const op = envelope.op ?? opCode;
    const payload = envelope.payload ?? envelope;
    // Merge envelope requestId/op into payload for OpServerError correlation (Day 28)
    const mergeError = (p: unknown): ServerError => {
      const obj = (p ?? {}) as Record<string, unknown>;
      return {
        code: (obj.code as string) ?? "unknown",
        message: (obj.message as string) ?? "",
        requestId: (obj.requestId as string) ?? (envelope.requestId as string | undefined),
        op: (obj.op as number) ?? (envelope.op as number | undefined) ?? op,
        details: obj.details,
      };
    };
    // Also handle case where data is already the payload (when opCode is used directly)
    // For Day 27, we route based on op
    if (op === 100) {
      onPrivateSnapshot(payload as PrivateSnapshot);
    } else if (op === 101) {
      onPublicSnapshot(payload as PublicSnapshot);
    } else if (op === 102) {
      onServerError(mergeError(payload));
    } else if (op === 103) {
      onServerEvent(payload);
    } else {
      console.log(`received unknown op ${op}`, payload);
    }
  } catch (e) {
    console.log(`failed to parse match data op ${opCode}`, e);
  }
}
