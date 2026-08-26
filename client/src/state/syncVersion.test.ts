import { describe, it, expect, beforeEach, vi } from "vitest";
import {
  onPrivateSnapshot,
  onPublicSnapshot,
  getLastPrivateSnapshot,
  getLastPublicSnapshot,
  handleMatchData,
  clearPrivateListeners,
  clearPublicListeners,
  clearAllPrivate,
} from "./sync";
import type { PrivateSnapshot, PublicSnapshot } from "./snapshot";

// Day 35: Versioning — Snapshot Version 1 check: if snap.v !== 1 log bad_version and ignore (as in Go parser.go:22)

describe("Sync versioning — Day 35", () => {
  beforeEach(() => {
    clearPrivateListeners();
    clearPublicListeners();
    clearAllPrivate();
    // Clear lastPublic by sending a valid one then clearing? Instead we just ensure next test starts clean by not relying on previous
    // Reset by sending invalid then checking that last snapshots are as expected
    // For simplicity, we don't clear lastPublic, but we test that bad version doesn't overwrite it
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: () => null,
        setItem: () => {},
        removeItem: () => {},
        key: () => null,
        length: 0,
        clear: () => {},
      } as unknown as Storage;
    }
  });

  it("onPrivateSnapshot ignores bad version", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const badSnap = {
      v: 2,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "t1", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    } as unknown as PrivateSnapshot;
    const before = getLastPrivateSnapshot();
    onPrivateSnapshot(badSnap);
    const after = getLastPrivateSnapshot();
    // Should not have updated (still before, which is null initially)
    expect(after).toBe(before);
    const hasBadVersionLog = consoleSpy.mock.calls.some((args) =>
      String(args[0]).includes("bad_version")
    );
    expect(hasBadVersionLog).toBe(true);
    consoleSpy.mockRestore();
  });

  it("onPublicSnapshot ignores bad version", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const goodSnap: PublicSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
    };
    onPublicSnapshot(goodSnap);
    const before = getLastPublicSnapshot();
    const badSnap = { ...goodSnap, v: 999 } as unknown as PublicSnapshot;
    onPublicSnapshot(badSnap);
    const after = getLastPublicSnapshot();
    expect(after).toEqual(before);
    const hasBadVersionLog = consoleSpy.mock.calls.some((args) =>
      String(args[0]).includes("bad_version")
    );
    expect(hasBadVersionLog).toBe(true);
    consoleSpy.mockRestore();
  });

  it("handleMatchData ignores envelope with bad version", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const goodSnap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "t1", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    // First, store a good one
    handleMatchData(100, JSON.stringify({ v: 1, op: 100, payload: goodSnap }));
    const before = getLastPrivateSnapshot();
    expect(before?.ownRack[0].ID).toBe("t1");
    // Now send bad version envelope
    const badEnvelope = JSON.stringify({ v: 999, op: 100, payload: { ...goodSnap, v: 999 } });
    handleMatchData(100, badEnvelope);
    const after = getLastPrivateSnapshot();
    expect(after).toEqual(before);
    const hasBadVersionLog = consoleSpy.mock.calls.some((args) =>
      String(args[0]).includes("bad_version")
    );
    expect(hasBadVersionLog).toBe(true);
    consoleSpy.mockRestore();
  });

  it("handleMatchData ignores payload with bad version even if envelope is ok", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const goodSnap: PublicSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
    };
    handleMatchData(101, JSON.stringify({ v: 1, op: 101, payload: goodSnap }));
    const before = getLastPublicSnapshot();
    const badPayload = { ...goodSnap, v: 2 };
    handleMatchData(101, JSON.stringify({ v: 1, op: 101, payload: badPayload }));
    const after = getLastPublicSnapshot();
    expect(after).toEqual(before);
    const hasBadVersionLog = consoleSpy.mock.calls.some((args) =>
      String(args[0]).includes("bad_version")
    );
    expect(hasBadVersionLog).toBe(true);
    consoleSpy.mockRestore();
  });

  it("valid version 1 is accepted", () => {
    const goodPriv: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 3 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "t1", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    onPrivateSnapshot(goodPriv);
    expect(getLastPrivateSnapshot()?.v).toBe(1);
    const goodPub: PublicSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 3 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
    };
    onPublicSnapshot(goodPub);
    expect(getLastPublicSnapshot()?.v).toBe(1);
  });
});
