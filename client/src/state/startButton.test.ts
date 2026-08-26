import { describe, it, expect } from "vitest";
import { shouldShowStartButton } from "./sync";
import type { PrivateSnapshot } from "./snapshot";
import { OpClientStart } from "../net/protocol";

// Day 37: Start match — Start button visible only if Waiting + host Seat 0 + >=2 players, sends OpClientStart 1 {}

function makePrivate(overrides: Partial<PrivateSnapshot>): PrivateSnapshot {
  return {
    v: 1,
    gamePhase: "Waiting",
    turnPhase: "MustDraw",
    currentSeat: 0,
    players: [
      { id: "alice", seat: 0, hasOpened: false, rackCount: 0 },
      { id: "bob", seat: 1, hasOpened: false, rackCount: 0 },
    ],
    stockCount: 0,
    discardRow: [],
    tableMelds: [],
    winner: -1,
    ownRack: [],
    ownSeat: 0,
    ...overrides,
  };
}

describe("Start button — Day 37", () => {
  it("shouldShowStartButton true only for Waiting host with >=2 players", () => {
    const snap = makePrivate({
      gamePhase: "Waiting",
      ownSeat: 0,
      players: [
        { id: "a", seat: 0, hasOpened: false, rackCount: 0 },
        { id: "b", seat: 1, hasOpened: false, rackCount: 0 },
      ],
    });
    expect(shouldShowStartButton(snap)).toBe(true);
  });

  it("shouldShowStartButton false when not Waiting", () => {
    const snap = makePrivate({
      gamePhase: "Playing",
      ownSeat: 0,
      players: [
        { id: "a", seat: 0, hasOpened: false, rackCount: 0 },
        { id: "b", seat: 1, hasOpened: false, rackCount: 0 },
      ],
    });
    expect(shouldShowStartButton(snap)).toBe(false);
    const snap2 = makePrivate({
      gamePhase: "OpeningDiscard",
      ownSeat: 0,
      players: [
        { id: "a", seat: 0, hasOpened: false, rackCount: 0 },
        { id: "b", seat: 1, hasOpened: false, rackCount: 0 },
      ],
    });
    expect(shouldShowStartButton(snap2)).toBe(false);
    const snap3 = makePrivate({
      gamePhase: "RoundComplete",
      ownSeat: 0,
      players: [
        { id: "a", seat: 0, hasOpened: false, rackCount: 0 },
        { id: "b", seat: 1, hasOpened: false, rackCount: 0 },
      ],
    });
    expect(shouldShowStartButton(snap3)).toBe(false);
  });

  it("shouldShowStartButton false when not host", () => {
    const snap = makePrivate({
      gamePhase: "Waiting",
      ownSeat: 1,
      players: [
        { id: "a", seat: 0, hasOpened: false, rackCount: 0 },
        { id: "b", seat: 1, hasOpened: false, rackCount: 0 },
      ],
    });
    expect(shouldShowStartButton(snap)).toBe(false);
    const snap2 = makePrivate({
      gamePhase: "Waiting",
      ownSeat: 2,
      players: [
        { id: "a", seat: 0, hasOpened: false, rackCount: 0 },
        { id: "b", seat: 1, hasOpened: false, rackCount: 0 },
        { id: "c", seat: 2, hasOpened: false, rackCount: 0 },
      ],
    });
    expect(shouldShowStartButton(snap2)).toBe(false);
  });

  it("shouldShowStartButton false when <2 players", () => {
    const snap = makePrivate({
      gamePhase: "Waiting",
      ownSeat: 0,
      players: [{ id: "a", seat: 0, hasOpened: false, rackCount: 0 }],
    });
    expect(shouldShowStartButton(snap)).toBe(false);
    const snap0 = makePrivate({ gamePhase: "Waiting", ownSeat: 0, players: [] });
    expect(shouldShowStartButton(snap0)).toBe(false);
  });

  it("OpClientStart is 1 and is the only allowed op in Waiting", async () => {
    expect(OpClientStart).toBe(1);
    // Verify via allowedOps that Waiting only allows Start
    const { allowedOps } = await import("./sync");
    const ops = allowedOps("Waiting", "MustDraw");
    expect(ops.has(OpClientStart)).toBe(true);
    expect(ops.size).toBe(1);
  });

  it("Start button would send OpClientStart via sendMatchState (code inspection)", async () => {
    const mod = await import("../net/protocol");
    expect(typeof mod.sendMatchState).toBe("function");
    // Use fs to read TableScene source without importing Phaser (which needs window)
    const fs = await import("fs");
    const path = await import("path");
    const filePath = path.resolve(__dirname, "../scenes/TableScene.ts");
    const src = fs.readFileSync(filePath, "utf-8");
    expect(src).toContain("OpClientStart");
    expect(src).toContain("sendMatchState");
    expect(src).toContain("shouldShowStartButton");
    expect(src).toContain("[Start]");
  });
});
