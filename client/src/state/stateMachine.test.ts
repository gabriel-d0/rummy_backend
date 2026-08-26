import { describe, it, expect } from "vitest";
import { allowedOps, isOpAllowedForSnapshot, canPlayerAct } from "./sync";
import {
  OpClientStart,
  OpClientDiscard,
  OpClientDrawStock,
  OpClientDrawPreviousDiscard,
  OpClientPickupDiscardForMeld,
  OpClientMeldInitial,
  OpClientMeldNew,
  OpClientExtendMeld,
  OpClientReplaceJoker,
} from "../net/protocol";
import type { PrivateSnapshot } from "./snapshot";

// Day 36: State machine client — mirrors Go phases.go:15 AllowedOps

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
    ownRack: [{ ID: "t1", Colour: 1, Rank: 5, IsJoker: false }],
    ownSeat: 0,
    ...overrides,
  };
}

describe("State machine — Day 36", () => {
  it("Waiting allows only Start", () => {
    const ops = allowedOps("Waiting", "MustDraw");
    expect(ops.has(OpClientStart)).toBe(true);
    expect(ops.has(OpClientDiscard)).toBe(false);
    expect(ops.has(OpClientDrawStock)).toBe(false);
    expect(ops.size).toBe(1);
  });

  it("OpeningDiscard allows only Discard", () => {
    const ops = allowedOps("OpeningDiscard", "MustDraw");
    expect(ops.has(OpClientDiscard)).toBe(true);
    expect(ops.has(OpClientStart)).toBe(false);
    expect(ops.has(OpClientDrawStock)).toBe(false);
    expect(ops.size).toBe(1);
  });

  it("Playing MustDraw allows DrawStock, DrawPrevious, Pickup", () => {
    const ops = allowedOps("Playing", "MustDraw");
    expect(ops.has(OpClientDrawStock)).toBe(true);
    expect(ops.has(OpClientDrawPreviousDiscard)).toBe(true);
    expect(ops.has(OpClientPickupDiscardForMeld)).toBe(true);
    expect(ops.has(OpClientDiscard)).toBe(false);
    expect(ops.has(OpClientMeldInitial)).toBe(false);
    expect(ops.size).toBe(3);
  });

  it("Playing MeldOrDiscard allows Discard, MeldInitial, MeldNew, Extend, Replace", () => {
    const ops = allowedOps("Playing", "MeldOrDiscard");
    expect(ops.has(OpClientDiscard)).toBe(true);
    expect(ops.has(OpClientMeldInitial)).toBe(true);
    expect(ops.has(OpClientMeldNew)).toBe(true);
    expect(ops.has(OpClientExtendMeld)).toBe(true);
    expect(ops.has(OpClientReplaceJoker)).toBe(true);
    expect(ops.has(OpClientDrawStock)).toBe(false);
    expect(ops.size).toBe(5);
  });

  it("RoundComplete allows none", () => {
    const ops = allowedOps("RoundComplete", "MustDraw");
    expect(ops.size).toBe(0);
    const ops2 = allowedOps("RoundComplete", "MeldOrDiscard");
    expect(ops2.size).toBe(0);
  });

  it("isOpAllowedForSnapshot mirrors allowedOps", () => {
    const snap = makePrivate({ gamePhase: "Playing", turnPhase: "MustDraw" });
    expect(isOpAllowedForSnapshot(snap, OpClientDrawStock)).toBe(true);
    expect(isOpAllowedForSnapshot(snap, OpClientDiscard)).toBe(false);
    const snap2 = makePrivate({ gamePhase: "Playing", turnPhase: "MeldOrDiscard" });
    expect(isOpAllowedForSnapshot(snap2, OpClientDiscard)).toBe(true);
    expect(isOpAllowedForSnapshot(snap2, OpClientDrawStock)).toBe(false);
  });

  it("canPlayerAct checks Waiting host only", () => {
    const snap = makePrivate({
      gamePhase: "Waiting",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
    });
    expect(canPlayerAct(snap, OpClientStart)).toBe(true);
    const snapNotHost = makePrivate({
      gamePhase: "Waiting",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 1,
    });
    expect(canPlayerAct(snapNotHost, OpClientStart)).toBe(false);
  });

  it("canPlayerAct checks currentSeat for Playing", () => {
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

  it("canPlayerAct checks HasOpened for DrawPrevious/Pickup/Extend/Replace/MeldNew", () => {
    const snapNotOpened = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
    });
    expect(canPlayerAct(snapNotOpened, OpClientDrawPreviousDiscard)).toBe(false);
    expect(canPlayerAct(snapNotOpened, OpClientPickupDiscardForMeld)).toBe(false);

    const snapOpenedMustDraw = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      ownSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: true, rackCount: 14 }],
    });
    expect(canPlayerAct(snapOpenedMustDraw, OpClientDrawPreviousDiscard)).toBe(true);

    const snapNotOpenedMeld = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      ownSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
    });
    expect(canPlayerAct(snapNotOpenedMeld, OpClientMeldNew)).toBe(false);
    expect(canPlayerAct(snapNotOpenedMeld, OpClientExtendMeld)).toBe(false);
    expect(canPlayerAct(snapNotOpenedMeld, OpClientReplaceJoker)).toBe(false);
    expect(canPlayerAct(snapNotOpenedMeld, OpClientMeldInitial)).toBe(true);

    const snapOpenedMeld = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      ownSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: true, rackCount: 14 }],
    });
    expect(canPlayerAct(snapOpenedMeld, OpClientMeldNew)).toBe(true);
    expect(canPlayerAct(snapOpenedMeld, OpClientMeldInitial)).toBe(false);
  });

  it("canPlayerAct respects turn + HasOpened for MeldOrDiscard Discard", () => {
    const snap = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      ownSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
    });
    expect(canPlayerAct(snap, OpClientDiscard)).toBe(true);
    const snapNotTurn = makePrivate({
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 1,
      ownSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
    });
    expect(canPlayerAct(snapNotTurn, OpClientDiscard)).toBe(false);
  });
});
