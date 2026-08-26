import { describe, it, expect } from "vitest";
import {
  SnapshotVersion,
  isValidPublicSnapshot,
  isValidPrivateSnapshot,
  checkNoLeak,
} from "./snapshot";
import type { PublicSnapshot, PrivateSnapshot } from "./snapshot";

// Day 29: Public snapshot types — mirrors Go internal/match/visibility.go:36 Version 1

describe("Snapshot — Day 29", () => {
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
    const snap = {
      v: 2,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
    };
    expect(isValidPublicSnapshot(snap)).toBe(false);
  });

  it("isValidPublicSnapshot rejects missing fields", () => {
    expect(isValidPublicSnapshot({ v: 1, gamePhase: "Playing" })).toBe(false);
    expect(isValidPublicSnapshot(null)).toBe(false);
    expect(isValidPublicSnapshot("bad")).toBe(false);
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
    expect(isValidPrivateSnapshot({ ...priv, ownRack: undefined })).toBe(false);
    expect(
      isValidPrivateSnapshot({
        v: 1,
        gamePhase: "Playing",
        turnPhase: "MustDraw",
        currentSeat: 0,
        players: [],
        stockCount: 77,
        discardRow: [],
        tableMelds: [],
        winner: -1,
      })
    ).toBe(false);
  });

  it("checkNoLeak detects private IDs in public JSON (redaction)", () => {
    const publicJson = JSON.stringify({
      v: 1,
      gamePhase: "Playing",
      players: [{ rackCount: 14 }],
      stockCount: 77,
      discardRow: [{ Tile: { ID: "disc-1", Colour: 1, Rank: 7, IsJoker: false } }],
      tableMelds: [],
    });
    expect(checkNoLeak(publicJson, ["disc-1"])).toBe(false); // disc-1 is public, so it leaks if we consider it private
    expect(checkNoLeak(publicJson, ["alice-secret"])).toBe(true);
    expect(checkNoLeak(publicJson, ["alice-secret", "bob-secret"])).toBe(true);
    const publicJsonWithLeak = JSON.stringify({ v: 1, ownRack: [{ ID: "alice-secret" }] });
    expect(checkNoLeak(publicJsonWithLeak, ["alice-secret"])).toBe(false);
  });

  it("PublicSnapshot mirrors Go visibility.go fields", () => {
    const snap: PublicSnapshot = {
      v: SnapshotVersion,
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 1,
      players: [
        { id: "alice", seat: 0, hasOpened: true, rackCount: 12 },
        { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
      ],
      stockCount: 63,
      discardRow: [
        {
          Tile: { ID: "d1", Colour: 1, Rank: 7, IsJoker: false },
          IsOpeningDiscard: true,
          Index: 0,
        },
      ],
      tableMelds: [{ ID: "m1", Kind: "run", Tiles: [], JokerReps: {}, OwnerSeat: 0 }],
      winner: -1,
    };
    const json = JSON.stringify(snap);
    expect(json).toContain("gamePhase");
    expect(json).toContain("rackCount");
    expect(json).not.toContain("ownRack");
    expect(isValidPublicSnapshot(JSON.parse(json))).toBe(true);
  });

  it("PrivateSnapshot extends PublicSnapshot with ownRack/ownSeat only", () => {
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
      ownRack: [
        { ID: "t1", Colour: 1, Rank: 1, IsJoker: false },
        { ID: "j1", Colour: 0, Rank: 0, IsJoker: true },
      ],
      ownSeat: 0,
    };
    expect(priv.ownRack.length).toBe(2);
    expect(priv.ownRack[1].IsJoker).toBe(true);
    expect(isValidPrivateSnapshot(priv)).toBe(true);
    const publicPart = { ...priv, ownRack: undefined, ownSeat: undefined } as unknown;
    expect(isValidPublicSnapshot(publicPart)).toBe(true);
    expect(isValidPrivateSnapshot(publicPart)).toBe(false);
  });
});
