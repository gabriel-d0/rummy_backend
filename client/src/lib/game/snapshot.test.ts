import { describe, it, expect } from "vitest";
import { SnapshotVersion, isValidPublicSnapshot, isValidPrivateSnapshot, checkNoLeak } from "./snapshot";
import type { PublicSnapshot, PrivateSnapshot } from "./snapshot";

describe("Snapshot — Day 17", () => {
  it("SnapshotVersion is 1", () => {
    expect(SnapshotVersion).toBe(1);
  });

  it("isValidPublicSnapshot accepts correct Version 1", () => {
    const snap: PublicSnapshot = {
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
    expect(isValidPublicSnapshot(snap)).toBe(true);
  });

  it("isValidPublicSnapshot rejects bad version", () => {
    const snap = { v: 2, gamePhase: "Playing", turnPhase: "MustDraw", currentSeat: 0, players: [], stockCount: 77, discardRow: [], tableMelds: [], winner: -1 };
    expect(isValidPublicSnapshot(snap)).toBe(false);
  });

  it("isValidPrivateSnapshot requires ownRack and ownSeat", () => {
    const priv: PrivateSnapshot = {
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
    expect(isValidPrivateSnapshot(priv)).toBe(true);
    expect(isValidPrivateSnapshot({ ...priv, ownRack: undefined } as unknown as PrivateSnapshot)).toBe(false);
  });

  it("checkNoLeak detects private IDs in public JSON", () => {
    const publicJson = JSON.stringify({ v: 1, players: [{ rackCount: 14 }], stockCount: 77 });
    expect(checkNoLeak(publicJson, ["alice-secret"])).toBe(true);
    const leakJson = JSON.stringify({ v: 1, ownRack: [{ ID: "alice-secret" }] });
    expect(checkNoLeak(leakJson, ["alice-secret"])).toBe(false);
  });
});
