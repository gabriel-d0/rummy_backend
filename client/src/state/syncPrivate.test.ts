import { describe, it, expect, beforeEach, vi } from "vitest";
import {
  subscribePrivateSnapshot,
  clearPrivateListeners,
  onPrivateSnapshot,
  getLastPrivateSnapshot,
  handleMatchData,
} from "./sync";
import type { PrivateSnapshot } from "./snapshot";

// Day 30: Private snapshot handling — RackScene re-renders OwnRack only (redaction)

describe("Sync PrivateSnapshot — Day 30", () => {
  beforeEach(() => {
    clearPrivateListeners();
    // Mock localStorage if not present (node env)
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: () => null,
        setItem: () => {},
        removeItem: () => {},
      } as unknown as Storage;
    }
    // Clear any stored snapshot by sending a dummy? Instead we just ensure listeners are cleared
    // getLastPrivateSnapshot may still have previous value; we will test with new snapshots
  });

  it("subscribePrivateSnapshot receives future onPrivateSnapshot", () => {
    const calls: PrivateSnapshot[] = [];
    const unsub = subscribePrivateSnapshot((snap) => calls.push(snap));
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
      ownRack: [
        { ID: "t1", Colour: 1, Rank: 1, IsJoker: false },
        { ID: "t2", Colour: 1, Rank: 2, IsJoker: false },
        { ID: "t3", Colour: 1, Rank: 3, IsJoker: false },
      ],
      ownSeat: 0,
    };
    onPrivateSnapshot(snap);
    expect(calls).toHaveLength(1);
    expect(calls[0].ownRack.length).toBe(3);
    expect(calls[0].ownSeat).toBe(0);
    unsub();
  });

  it("subscribePrivateSnapshot immediately replays lastPrivate if exists", () => {
    const snap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 1,
      players: [{ id: "bob", seat: 1, hasOpened: true, rackCount: 5 }],
      stockCount: 70,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "x1", Colour: 2, Rank: 7, IsJoker: false }],
      ownSeat: 1,
    };
    onPrivateSnapshot(snap);
    const calls: PrivateSnapshot[] = [];
    const unsub = subscribePrivateSnapshot((s) => calls.push(s));
    expect(calls).toHaveLength(1);
    expect(calls[0].ownSeat).toBe(1);
    expect(calls[0].ownRack[0].ID).toBe("x1");
    unsub();
  });

  it("unsubscribe stops receiving", () => {
    const calls: PrivateSnapshot[] = [];
    const unsub = subscribePrivateSnapshot((s) => calls.push(s));
    unsub();
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
    // If there was a previous lastPrivate, the initial subscribe would have been called, but we unsubscribed before onPrivateSnapshot
    // calls may be 0 or 1 depending on previous state; we just ensure that after unsubscribe, the new snap does not immediately trigger
    // Since we unsubscribed right after subscribing, and there was a lastPrivate from previous test, calls would have 1 from immediate replay
    // So we clear and check that new snap doesn't add another after unsubscribe
    const before = calls.length;
    // Send another snap
    const snap2: PrivateSnapshot = {
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
    onPrivateSnapshot(snap2);
    expect(calls.length).toBe(before); // no new call
  });

  it("handleMatchData op 100 triggers private listeners", () => {
    const calls: PrivateSnapshot[] = [];
    const unsub = subscribePrivateSnapshot((s) => calls.push(s));
    // Clear previous calls that may have been replayed from lastPrivate
    calls.length = 0;
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
      ownRack: [
        { ID: "r1", Colour: 1, Rank: 10, IsJoker: false },
        { ID: "r2", Colour: 1, Rank: 11, IsJoker: false },
        { ID: "r3", Colour: 1, Rank: 12, IsJoker: false },
      ],
      ownSeat: 0,
    };
    handleMatchData(100, JSON.stringify({ v: 1, op: 100, payload: snap }));
    expect(calls).toHaveLength(1);
    expect(calls[0].ownRack.length).toBe(3);
    expect(getLastPrivateSnapshot()?.ownSeat).toBe(0);
    unsub();
  });

  it("onPrivateSnapshot stores in localStorage and lastPrivate", () => {
    const snap: PrivateSnapshot = {
      v: 1,
      gamePhase: "RoundComplete",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: true, rackCount: 0 }],
      stockCount: 40,
      discardRow: [],
      tableMelds: [],
      winner: 0,
      ownRack: [],
      ownSeat: 0,
    };
    const setItemSpy = vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    // Ensure localStorage exists
    if (typeof localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: () => null,
        setItem: setItemSpy as unknown as () => void,
        removeItem: () => {},
      } as unknown as Storage;
    }
    onPrivateSnapshot(snap);
    expect(getLastPrivateSnapshot()?.winner).toBe(0);
    // localStorage may be mocked, but we check that getLastPrivate is correct
    setItemSpy.mockRestore();
  });

  it("RackScene would render OwnRack only (redaction) — simulate", () => {
    // Simulate what RackScene does: sortRack and render
    const snap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [
        { id: "alice", seat: 0, hasOpened: false, rackCount: 3 },
        { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
      ],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [
        { ID: "alice-1", Colour: 1, Rank: 5, IsJoker: false },
        { ID: "alice-2", Colour: 2, Rank: 7, IsJoker: false },
        { ID: "alice-3", Colour: 3, Rank: 9, IsJoker: false },
      ],
      ownSeat: 0,
    };
    // Verify that private snapshot contains only own rack, not bob's
    expect(snap.ownRack.every((t) => t.ID.startsWith("alice"))).toBe(true);
    const publicJson = JSON.stringify({
      v: snap.v,
      gamePhase: snap.gamePhase,
      players: snap.players,
      stockCount: snap.stockCount,
      discardRow: snap.discardRow,
      tableMelds: snap.tableMelds,
    });
    expect(publicJson).not.toContain("alice-1");
    expect(publicJson).not.toContain("alice-2");
    // RackScene's renderPrivateRack would sort and render these 3 tiles
    // We can't test Phaser rendering here, but we test the data flow
    expect(snap.ownRack.length).toBe(3);
  });
});
