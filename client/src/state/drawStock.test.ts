import { describe, it, expect } from "vitest";
import { canPlayerAct, allowedOps } from "./sync";
import { OpClientDrawStock } from "../net/protocol";
import type { PrivateSnapshot } from "./snapshot";

// Day 39: Draw stock — Draw button visible only if Playing MustDraw and ownSeat==currentSeat, sends OpClientDrawStock 3 {}

function makePrivate(overrides: Partial<PrivateSnapshot>): PrivateSnapshot {
  return {
    v: 1,
    gamePhase: "Playing",
    turnPhase: "MustDraw",
    currentSeat: 0,
    players: [
      { id: "alice", seat: 0, hasOpened: false, rackCount: 14 },
      { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
    ],
    stockCount: 77,
    discardRow: [],
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

describe("Draw stock — Day 39", () => {
  it("allowedOps Playing MustDraw includes DrawStock", () => {
    const ops = allowedOps("Playing", "MustDraw");
    expect(ops.has(OpClientDrawStock)).toBe(true);
    expect(ops.has(OpClientDrawStock)).toBe(true);
    expect(ops.size).toBe(3);
  });

  it("allowedOps Playing MeldOrDiscard does not include DrawStock", () => {
    const ops = allowedOps("Playing", "MeldOrDiscard");
    expect(ops.has(OpClientDrawStock)).toBe(false);
  });

  it("canPlayerAct DrawStock true only when MustDraw and currentSeat", () => {
    const snap = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(snap, OpClientDrawStock)).toBe(true);
    const snapNotTurn = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 1,
      ownSeat: 0,
    });
    expect(canPlayerAct(snapNotTurn, OpClientDrawStock)).toBe(false);
  });

  it("canPlayerAct DrawStock false in MeldOrDiscard", () => {
    const snap = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(snap, OpClientDrawStock)).toBe(false);
  });

  it("canPlayerAct DrawStock false in Waiting/OpeningDiscard/RoundComplete", () => {
    const waiting = makePrivate({
      gamePhase: "Waiting",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(waiting, OpClientDrawStock)).toBe(false);
    const opening = makePrivate({
      gamePhase: "OpeningDiscard",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(opening, OpClientDrawStock)).toBe(false);
    const complete = makePrivate({
      gamePhase: "RoundComplete",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(complete, OpClientDrawStock)).toBe(false);
  });

  it("RackScene Draw button would send OpClientDrawStock via sendMatchState (code inspection)", async () => {
    const fs = await import("fs");
    const path = await import("path");
    const filePath = path.resolve(__dirname, "../scenes/RackScene.ts");
    const src = fs.readFileSync(filePath, "utf-8");
    expect(src).toContain("OpClientDrawStock");
    expect(src).toContain("sendMatchState");
    expect(src).toContain("Draw");
    expect(src).toContain("req-draw-");
  });

  it("Draw button disables after click until MeldOrDiscard (simulated via canPlayerAct)", () => {
    const mustDrawSnap = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(mustDrawSnap, OpClientDrawStock)).toBe(true);
    const meldSnap = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(meldSnap, OpClientDrawStock)).toBe(false);
    // After draw, the server transitions to MeldOrDiscard, so Draw should be disabled
  });
});
