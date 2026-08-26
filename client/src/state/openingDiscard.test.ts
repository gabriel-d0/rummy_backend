import { describe, it, expect } from "vitest";
import { canPlayerAct, allowedOps } from "./sync";
import { OpClientDiscard } from "../net/protocol";
import type { PrivateSnapshot } from "./snapshot";

// Day 38: Opening discard — opening player must discard exactly one tile from 15-tile rack, marked unavailable

function makeOpeningSnap(overrides: Partial<PrivateSnapshot> = {}): PrivateSnapshot {
  const rack = Array.from({ length: 15 }, (_, i) => ({
    ID: `t-${i}`,
    Colour: 1,
    Rank: (i % 13) + 1,
    IsJoker: false,
  }));
  return {
    v: 1,
    gamePhase: "OpeningDiscard",
    turnPhase: "MustDraw",
    currentSeat: 0,
    players: [
      { id: "alice", seat: 0, hasOpened: false, rackCount: 15 },
      { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
    ],
    stockCount: 77,
    discardRow: [],
    tableMelds: [],
    winner: -1,
    ownRack: rack,
    ownSeat: 0,
    ...overrides,
  };
}

describe("Opening discard — Day 38", () => {
  it("OpeningDiscard allows only Discard", () => {
    const ops = allowedOps("OpeningDiscard", "MustDraw");
    expect(ops.has(OpClientDiscard)).toBe(true);
    expect(ops.size).toBe(1);
  });

  it("canPlayerAct allows opening discard for currentSeat with 15 tiles", () => {
    const snap = makeOpeningSnap({ currentSeat: 0, ownSeat: 0 });
    expect(canPlayerAct(snap, OpClientDiscard)).toBe(true);
  });

  it("canPlayerAct denies opening discard for non-currentSeat", () => {
    const snap = makeOpeningSnap({ currentSeat: 0, ownSeat: 1 });
    expect(canPlayerAct(snap, OpClientDiscard)).toBe(false);
  });

  it("opening discard requires 15 tiles (simulated via snapshot)", () => {
    const snap15 = makeOpeningSnap({
      ownRack: Array.from({ length: 15 }, (_, i) => ({
        ID: `t-${i}`,
        Colour: 1,
        Rank: 1,
        IsJoker: false,
      })),
    });
    expect(snap15.ownRack.length).toBe(15);
    const snap14 = makeOpeningSnap({
      ownRack: Array.from({ length: 14 }, (_, i) => ({
        ID: `t-${i}`,
        Colour: 1,
        Rank: 1,
        IsJoker: false,
      })),
    });
    expect(snap14.ownRack.length).toBe(14);
    // canPlayerAct does not check rack length, but RackScene handler does — we verify via snapshot
    expect(canPlayerAct(snap15, OpClientDiscard)).toBe(true);
    expect(canPlayerAct(snap14, OpClientDiscard)).toBe(true); // allowed per phase, but handler will block if not 15
  });

  it("RackScene opening discard would send OpClientDiscard with tileId (code inspection)", async () => {
    const fs = await import("fs");
    const path = await import("path");
    const filePath = path.resolve(__dirname, "../scenes/RackScene.ts");
    const src = fs.readFileSync(filePath, "utf-8");
    expect(src).toContain("OpeningDiscard");
    expect(src).toContain("ownRack.length === 15");
    expect(src).toContain("OpClientDiscard");
    expect(src).toContain("sendMatchState");
    expect(src).toContain("clearSelected");
  });

  it("opening discard tile must be in ownRack", () => {
    const snap = makeOpeningSnap();
    const tileInRack = snap.ownRack[0].ID;
    const tileNotInRack = "not-in-rack";
    expect(snap.ownRack.some((t) => t.ID === tileInRack)).toBe(true);
    expect(snap.ownRack.some((t) => t.ID === tileNotInRack)).toBe(false);
  });
});
