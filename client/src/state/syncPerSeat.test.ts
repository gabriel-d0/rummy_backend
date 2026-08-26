import { describe, it, expect, beforeEach, vi } from "vitest";
import {
  onPrivateSnapshot,
  getLastPrivateBySeat,
  getAllPrivateBySeat,
  clearAllPrivate,
  clearPrivateListeners,
} from "./sync";
import type { PrivateSnapshot } from "./snapshot";

// Day 33: Reconnection — store last snapshot per Seat in localStorage/memory, keep matchId/userId on disconnect

describe("Sync per-Seat — Day 33", () => {
  beforeEach(() => {
    clearPrivateListeners();
    clearAllPrivate();
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: (_k: string) => null as unknown as string | null,
        setItem: vi.fn(),
        removeItem: vi.fn(),
        key: () => null,
        length: 0,
        clear: () => {},
      } as unknown as Storage;
    }
  });

  it("stores PrivateSnapshot per Seat in memory", () => {
    const snap0: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "a1", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    const snap1: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 1,
      players: [{ id: "bob", seat: 1, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "b1", Colour: 2, Rank: 7, IsJoker: false }],
      ownSeat: 1,
    };
    onPrivateSnapshot(snap0);
    onPrivateSnapshot(snap1);
    expect(getLastPrivateBySeat(0)?.ownRack[0].ID).toBe("a1");
    expect(getLastPrivateBySeat(1)?.ownRack[0].ID).toBe("b1");
    expect(getAllPrivateBySeat().size).toBe(2);
  });

  it("stores per-Seat in localStorage with seat suffix", () => {
    const calls: string[] = [];
    const ls = globalThis.localStorage as unknown as {
      setItem: (k: string, v: string) => void;
      getItem: (k: string) => string | null;
      removeItem: (k: string) => void;
      key: (i: number) => string | null;
      length: number;
    };
    const origSetItem = ls.setItem;
    const spy = vi.fn((k: string, v: string) => {
      calls.push(k);
      try {
        return (origSetItem as unknown as (k: string, v: string) => void).call(ls, k, v);
      } catch {
        // ignore if mock throws
      }
    });
    (globalThis as unknown as Record<string, unknown>).localStorage = {
      ...ls,
      setItem: spy as unknown as typeof ls.setItem,
      getItem: ls.getItem,
      removeItem: ls.removeItem,
      key: ls.key,
      length: ls.length,
      clear: () => {},
    } as unknown as Storage;
    const snap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 3 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "t1", Colour: 1, Rank: 1, IsJoker: false }],
      ownSeat: 2,
    };
    onPrivateSnapshot(snap);
    expect(calls.some((k) => k === "rummy_lastPrivate:2")).toBe(true);
    expect(calls.some((k) => k === "rummy_lastPrivate")).toBe(true);
    expect(calls.some((k) => k === "rummy_lastPrivate:map")).toBe(true);
    // restore
    (globalThis as unknown as Record<string, unknown>).localStorage = ls as unknown as Storage;
  });

  it("clearAllPrivate removes all per-Seat entries", () => {
    const snap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "a1", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    onPrivateSnapshot(snap);
    expect(getAllPrivateBySeat().size).toBe(1);
    clearAllPrivate();
    expect(getAllPrivateBySeat().size).toBe(0);
    expect(getLastPrivateBySeat(0)).toBeNull();
  });

  it("overwrites same seat with latest OwnRack", () => {
    const snap1: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "old", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    const snap2: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: true, rackCount: 12 }],
      stockCount: 70,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "new", Colour: 2, Rank: 7, IsJoker: false }],
      ownSeat: 0,
    };
    onPrivateSnapshot(snap1);
    expect(getLastPrivateBySeat(0)?.ownRack[0].ID).toBe("old");
    onPrivateSnapshot(snap2);
    expect(getLastPrivateBySeat(0)?.ownRack[0].ID).toBe("new");
    expect(getLastPrivateBySeat(0)?.stockCount).toBe(70);
  });
});
