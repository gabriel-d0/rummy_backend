import type { PrivateSnapshot, PublicSnapshot } from "./snapshot";
import {
  SnapshotVersion,
  isValidPublicSnapshot,
  isValidPrivateSnapshot,
  checkNoLeak as checkNoLeakSnapshot,
} from "./snapshot";
import { showErrorToast, type ServerError } from "../ui/ErrorToast";

// Day 27: Receive match state — parses Envelope and routes op 100/101/102/103 to handlers
// Mirrors Go internal/protocol/opcodes.go:8 Version 1 and internal/match/visibility.go:36
// Day 28: OpServerError 102 is routed to ErrorToast (red toast 3s)

let lastPrivateSnapshot: PrivateSnapshot | null = null;
let lastPublicSnapshot: PublicSnapshot | null = null;

// Day 33: per-Seat storage for reconnection (ownSeat -> PrivateSnapshot)
const lastPrivateBySeat = new Map<number, PrivateSnapshot>();

const privateListeners: ((snap: PrivateSnapshot) => void)[] = [];
const publicListeners: ((snap: PublicSnapshot) => void)[] = [];

export function subscribePrivateSnapshot(cb: (snap: PrivateSnapshot) => void): () => void {
  privateListeners.push(cb);
  if (lastPrivateSnapshot) {
    try {
      cb(lastPrivateSnapshot);
    } catch {
      // ignore listener error during immediate replay
    }
  }
  return () => {
    const idx = privateListeners.indexOf(cb);
    if (idx !== -1) privateListeners.splice(idx, 1);
  };
}

export function subscribePublicSnapshot(cb: (snap: PublicSnapshot) => void): () => void {
  publicListeners.push(cb);
  if (lastPublicSnapshot) {
    try {
      cb(lastPublicSnapshot);
    } catch {
      // ignore
    }
  }
  return () => {
    const idx = publicListeners.indexOf(cb);
    if (idx !== -1) publicListeners.splice(idx, 1);
  };
}

// For tests: clear listeners without affecting stored snapshot
export function clearPrivateListeners(): void {
  privateListeners.length = 0;
}

export function clearPublicListeners(): void {
  publicListeners.length = 0;
}

// Day 33: per-Seat accessors for reconnection
export function getLastPrivateBySeat(seat: number): PrivateSnapshot | null {
  return lastPrivateBySeat.get(seat) ?? null;
}

export function getAllPrivateBySeat(): ReadonlyMap<number, PrivateSnapshot> {
  return new Map(lastPrivateBySeat);
}

export function clearAllPrivate(): void {
  lastPrivateBySeat.clear();
  lastPrivateSnapshot = null;
  try {
    const keys: string[] = [];
    for (let i = 0; i < localStorage.length; i++) {
      const k = localStorage.key(i);
      if (k && k.startsWith("rummy_lastPrivate")) keys.push(k);
    }
    for (const k of keys) localStorage.removeItem(k);
  } catch {
    // ignore in node test env
  }
}

export function onPrivateSnapshot(snap: PrivateSnapshot): void {
  // Day 35: Version 1 check — if snap.v !== SnapshotVersion log bad_version and ignore (as in Go parser.go:22)
  if (!isValidPrivateSnapshot(snap)) {
    if ((snap as unknown as { v?: unknown })?.v !== SnapshotVersion) {
      console.log(
        `bad_version: PrivateSnapshot v=${(snap as unknown as { v?: unknown })?.v} want ${SnapshotVersion} — ignore — Day 35`
      );
    } else {
      console.log("bad_version: invalid PrivateSnapshot — ignore — Day 35");
    }
    return;
  }
  lastPrivateSnapshot = snap;
  lastPrivateBySeat.set(snap.ownSeat, snap);
  console.log(
    `received op 100 PrivateSnapshot v=${snap.v} gamePhase=${snap.gamePhase} currentSeat=${snap.currentSeat} ownSeat=${snap.ownSeat} rack=${snap.ownRack.length} — Day 33`
  );
  try {
    localStorage.setItem("rummy_lastPrivate", JSON.stringify(snap));
    localStorage.setItem(`rummy_lastPrivate:${snap.ownSeat}`, JSON.stringify(snap));
    const mapObj: Record<string, PrivateSnapshot> = {};
    for (const [seat, s] of lastPrivateBySeat) mapObj[String(seat)] = s;
    localStorage.setItem("rummy_lastPrivate:map", JSON.stringify(mapObj));
  } catch {
    // ignore in node test env where localStorage is mocked or missing
  }
  for (const cb of [...privateListeners]) {
    try {
      cb(snap);
    } catch (e) {
      console.log("privateListener error", e);
    }
  }
}

export function onPublicSnapshot(snap: PublicSnapshot): void {
  // Day 35: Version 1 check
  if (!isValidPublicSnapshot(snap)) {
    if ((snap as unknown as { v?: unknown })?.v !== SnapshotVersion) {
      console.log(
        `bad_version: PublicSnapshot v=${(snap as unknown as { v?: unknown })?.v} want ${SnapshotVersion} — ignore — Day 35`
      );
    } else {
      console.log("bad_version: invalid PublicSnapshot — ignore — Day 35");
    }
    return;
  }
  lastPublicSnapshot = snap;
  console.log(
    `received op 101 PublicSnapshot v=${snap.v} gamePhase=${snap.gamePhase} currentSeat=${snap.currentSeat} tableMelds=${snap.tableMelds.length} discardRow=${snap.discardRow.length} — Day 31`
  );
  // Day 32: redaction check — ensure PublicSnapshot JSON does not contain OwnRack IDs (as in visibility_test.go)
  if (lastPrivateSnapshot) {
    const privateIds = lastPrivateSnapshot.ownRack.map((t) => t.ID);
    const publicJson = JSON.stringify(snap);
    const ok = checkNoLeakSnapshot(publicJson, privateIds);
    if (ok) {
      console.log("checkNoLeak: no leak — PublicSnapshot does not contain OwnRack IDs");
    } else {
      console.log("LEAKED: PublicSnapshot contains OwnRack ID — redaction failure");
    }
  }
  for (const cb of [...publicListeners]) {
    try {
      cb(snap);
    } catch (e) {
      console.log("publicListener error", e);
    }
  }
}

// Day 32: client-side redaction check — mirrors Go visibility_test.go string search for OwnRack IDs
export function checkNoLeak(publicJson: string, privateIds: string[]): boolean {
  const ok = checkNoLeakSnapshot(publicJson, privateIds);
  if (ok) {
    console.log("checkNoLeak: no leak");
  } else {
    console.log("LEAKED: publicJson contains privateId");
  }
  return ok;
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
    // Day 35: Version 1 check for envelope (as in Go parser.go:22)
    if (envelope.v !== undefined && envelope.v !== SnapshotVersion) {
      console.log(
        `bad_version: envelope v=${envelope.v} want ${SnapshotVersion} — ignore — Day 35`
      );
      return;
    }
    const op = envelope.op ?? opCode;
    const payload = envelope.payload ?? envelope;
    // Day 35: also check payload version for snapshots (op 100/101)
    if (
      (op === 100 || op === 101) &&
      payload &&
      typeof payload === "object" &&
      (payload as { v?: unknown }).v !== undefined &&
      (payload as { v?: unknown }).v !== SnapshotVersion
    ) {
      console.log(
        `bad_version: payload v=${(payload as { v?: unknown }).v} want ${SnapshotVersion} — ignore — Day 35`
      );
      return;
    }
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
