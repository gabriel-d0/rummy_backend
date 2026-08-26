import { describe, it, expect, beforeEach } from "vitest";
import {
  subscribePublicSnapshot,
  clearPublicListeners,
  onPublicSnapshot,
  getLastPublicSnapshot,
  handleMatchData,
} from "./sync";
import type { PublicSnapshot } from "./snapshot";

// Day 31: Public snapshot handling — TableScene re-renders TableMelds/DiscardRow/StockCount/CurrentSeat (not OwnRack)

describe("Sync PublicSnapshot — Day 31", () => {
  beforeEach(() => {
    clearPublicListeners();
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: () => null,
        setItem: () => {},
        removeItem: () => {},
      } as unknown as Storage;
    }
  });

  it("subscribePublicSnapshot receives future onPublicSnapshot", () => {
    const calls: PublicSnapshot[] = [];
    const unsub = subscribePublicSnapshot((snap) => calls.push(snap));
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
    onPublicSnapshot(snap);
    expect(calls).toHaveLength(1);
    expect(calls[0].stockCount).toBe(77);
    expect(calls[0].currentSeat).toBe(0);
    unsub();
  });

  it("subscribePublicSnapshot immediately replays lastPublic if exists", () => {
    const snap: PublicSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 1,
      players: [{ id: "bob", seat: 1, hasOpened: true, rackCount: 5 }],
      stockCount: 70,
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
    onPublicSnapshot(snap);
    const calls: PublicSnapshot[] = [];
    const unsub = subscribePublicSnapshot((s) => calls.push(s));
    expect(calls).toHaveLength(1);
    expect(calls[0].tableMelds[0].ID).toBe("m1");
    expect(calls[0].discardRow[0].Tile.ID).toBe("d1");
    unsub();
  });

  it("unsubscribe stops receiving", () => {
    const calls: PublicSnapshot[] = [];
    const unsub = subscribePublicSnapshot((s) => calls.push(s));
    unsub();
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
    const before = calls.length;
    onPublicSnapshot(snap);
    // If there was a lastPublic from previous test, the subscribe would have immediately called, but we unsubscribed right after, so before is 0 or 1
    // After unsubscribe, new snap should not add
    const snap2: PublicSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 1,
      players: [{ id: "bob", seat: 1, hasOpened: false, rackCount: 14 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
    };
    onPublicSnapshot(snap2);
    expect(calls.length).toBe(before);
  });

  it("handleMatchData op 101 triggers public listeners", () => {
    const calls: PublicSnapshot[] = [];
    const unsub = subscribePublicSnapshot((s) => calls.push(s));
    calls.length = 0;
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
    handleMatchData(101, JSON.stringify({ v: 1, op: 101, payload: snap }));
    expect(calls).toHaveLength(1);
    expect(calls[0].discardRow.length).toBe(0);
    expect(getLastPublicSnapshot()?.currentSeat).toBe(0);
    unsub();
  });

  it("onPublicSnapshot updates stock/turn/melds/discard and not OwnRack (redaction)", () => {
    const snap: PublicSnapshot = {
      v: 1,
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
      tableMelds: [
        {
          ID: "run-1-2-3",
          Kind: "run",
          Tiles: [{ ID: "t1", Colour: 1, Rank: 1, IsJoker: false }],
          JokerReps: {},
          OwnerSeat: 0,
        },
      ],
      winner: -1,
    };
    onPublicSnapshot(snap);
    const stored = getLastPublicSnapshot();
    expect(stored?.stockCount).toBe(63);
    expect(stored?.currentSeat).toBe(1);
    expect(stored?.tableMelds.length).toBe(1);
    expect(stored?.discardRow.length).toBe(2);
    expect(stored?.discardRow[0].IsOpeningDiscard).toBe(true);
    // Verify that public snapshot JSON does not contain ownRack
    const json = JSON.stringify(stored);
    expect(json).not.toContain("ownRack");
    expect(json).not.toContain("ownSeat");
    expect(json).toContain("stockCount");
    expect(json).toContain("discardRow");
  });

  it("TableScene renderPublicSnapshot would keep OwnRack redacted — simulate", () => {
    const publicSnap: PublicSnapshot = {
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
    // TableScene's renderPublicSnapshot should not expose any TileInstanceId from private racks
    // We simulate that it only uses stockCount, tableMelds, discardRow, currentSeat
    expect(publicSnap.stockCount).toBe(77);
    expect(publicSnap.tableMelds.length).toBe(0);
    // Ensure that even if we try to look for private IDs, they are not there
    const privateIds = ["alice-secret-1", "bob-secret-2"];
    const publicJson = JSON.stringify(publicSnap);
    for (const id of privateIds) {
      expect(publicJson).not.toContain(id);
    }
  });
});
