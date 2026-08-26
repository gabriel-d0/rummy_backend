import { describe, it, expect } from "vitest";
import { canDrawPrevious, canPlayerAct, allowedOps } from "./sync";
import { OpClientDrawPreviousDiscard } from "../net/protocol";
import type { PrivateSnapshot } from "./snapshot";

// Day 40: Draw previous discard — DrawPrevious button visible only if HasOpened and DiscardRow not empty and not IsOpeningDiscard

function makePrivate(overrides: Partial<PrivateSnapshot>): PrivateSnapshot {
  return {
    v: 1,
    gamePhase: "Playing",
    turnPhase: "MustDraw",
    currentSeat: 0,
    players: [
      { id: "alice", seat: 0, hasOpened: true, rackCount: 14 },
      { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
    ],
    stockCount: 77,
    discardRow: [
      { Tile: { ID: "d1", Colour: 1, Rank: 7, IsJoker: false }, IsOpeningDiscard: false, Index: 0 },
    ],
    tableMelds: [],
    winner: -1,
    ownRack: Array.from({ length: 14 }, (_, i) => ({
      ID: `t-${i}`,
      Colour: 1,
      Rank: 1,
      IsJoker: false,
    })),
    ownSeat: 0,
    ...overrides,
  };
}

describe("Draw previous discard — Day 40", () => {
  it("allowedOps MustDraw includes DrawPrevious", () => {
    const ops = allowedOps("Playing", "MustDraw");
    expect(ops.has(OpClientDrawPreviousDiscard)).toBe(true);
  });

  it("canDrawPrevious true only when HasOpened and discardRow not empty and not opening", () => {
    const snap = makePrivate({});
    expect(canDrawPrevious(snap)).toBe(true);
  });

  it("canDrawPrevious false when not HasOpened", () => {
    const snap = makePrivate({
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
    });
    expect(canDrawPrevious(snap)).toBe(false);
    expect(canPlayerAct(snap, OpClientDrawPreviousDiscard)).toBe(false);
  });

  it("canDrawPrevious false when discardRow empty", () => {
    const snap = makePrivate({ discardRow: [] });
    expect(canDrawPrevious(snap)).toBe(false);
  });

  it("canDrawPrevious false when last discard is opening", () => {
    const snap = makePrivate({
      discardRow: [
        {
          Tile: { ID: "disc-open", Colour: 1, Rank: 7, IsJoker: false },
          IsOpeningDiscard: true,
          Index: 0,
        },
      ],
    });
    expect(canDrawPrevious(snap)).toBe(false);
  });

  it("canDrawPrevious false when not MustDraw or not currentSeat", () => {
    const snapNotMustDraw = makePrivate({ turnPhase: "MeldOrDiscard" });
    expect(canDrawPrevious(snapNotMustDraw)).toBe(false);
    const snapNotTurn = makePrivate({ currentSeat: 1, ownSeat: 0 });
    expect(canDrawPrevious(snapNotTurn)).toBe(false);
  });

  it("canDrawPrevious handles multiple discards, only last matters", () => {
    const snap = makePrivate({
      discardRow: [
        {
          Tile: { ID: "disc-open", Colour: 1, Rank: 7, IsJoker: false },
          IsOpeningDiscard: true,
          Index: 0,
        },
        {
          Tile: { ID: "disc-1", Colour: 2, Rank: 3, IsJoker: false },
          IsOpeningDiscard: false,
          Index: 1,
        },
      ],
    });
    expect(canDrawPrevious(snap)).toBe(true);
    const snapOnlyOpening = makePrivate({
      discardRow: [
        {
          Tile: { ID: "disc-open", Colour: 1, Rank: 7, IsJoker: false },
          IsOpeningDiscard: true,
          Index: 0,
        },
      ],
    });
    expect(canDrawPrevious(snapOnlyOpening)).toBe(false);
  });

  it("RackScene DrawPrev button would send OpClientDrawPreviousDiscard (code inspection)", async () => {
    const fs = await import("fs");
    const path = await import("path");
    const filePath = path.resolve(__dirname, "../scenes/RackScene.ts");
    const src = fs.readFileSync(filePath, "utf-8");
    expect(src).toContain("OpClientDrawPreviousDiscard");
    expect(src).toContain("sendMatchState");
    expect(src).toContain("[Prev]");
    expect(src).toContain("canDrawPrevious");
  });
});
