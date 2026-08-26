import { describe, it, expect, beforeEach } from "vitest";
import { reconnect, getStoredMatchId, getStoredUserId, isReconnectionAvailable } from "./nakama";
import type { PrivateSnapshot } from "../state/snapshot";

// Day 34: Reconnection — rejoin via socket.connect + joinMatch(matchId), expect PrivateSnapshot for ownSeat

describe("Nakama reconnect — Day 34", () => {
  beforeEach(() => {
    // Setup in-memory localStorage that actually stores values (for isReconnectionAvailable test)
    const store = new Map<string, string>();
    (globalThis as unknown as Record<string, unknown>).localStorage = {
      getItem: (k: string) => store.get(k) ?? null,
      setItem: (k: string, v: string) => {
        store.set(k, v);
      },
      removeItem: (k: string) => {
        store.delete(k);
      },
      key: (i: number) => Array.from(store.keys())[i] ?? null,
      get length() {
        return store.size;
      },
      clear: () => store.clear(),
    } as unknown as Storage;
    try {
      localStorage.removeItem("rummy_matchId");
      localStorage.removeItem("rummy_token");
      localStorage.removeItem("rummy_userId");
    } catch {
      // ignore
    }
  });

  it("getStoredMatchId / getStoredUserId read from localStorage", () => {
    localStorage.setItem("rummy_matchId", "match-123");
    localStorage.setItem("rummy_userId", "user-abc");
    expect(getStoredMatchId()).toBe("match-123");
    expect(getStoredUserId()).toBe("user-abc");
  });

  it("isReconnectionAvailable true only when both matchId and userId present", () => {
    localStorage.removeItem("rummy_matchId");
    localStorage.removeItem("rummy_userId");
    expect(isReconnectionAvailable()).toBe(false);
    localStorage.setItem("rummy_matchId", "m1");
    expect(isReconnectionAvailable()).toBe(false);
    localStorage.setItem("rummy_userId", "u1");
    expect(isReconnectionAvailable()).toBe(true);
    localStorage.removeItem("rummy_matchId");
    expect(isReconnectionAvailable()).toBe(false);
  });

  it("reconnect returns null when no stored matchId", async () => {
    localStorage.removeItem("rummy_matchId");
    const result = await reconnect();
    expect(result).toBeNull();
  });

  it("reconnect exposes function and handles stored matchId (code inspection)", async () => {
    const mod = await import("./nakama");
    expect(typeof mod.reconnect).toBe("function");
    expect(typeof mod.getStoredMatchId).toBe("function");
    expect(typeof mod.isReconnectionAvailable).toBe("function");
    // Check that reconnect source contains joinMatch and reauth logic
    const src = mod.reconnect.toString();
    expect(src).toContain("joinMatch");
    expect(src).toContain("rummy_matchId");
    expect(src).toContain("OpServerState 100");
  });

  it("sync rehydrates RackScene from OwnRack not old (per-Seat)", async () => {
    const sync = await import("../state/sync");
    sync.clearAllPrivate();
    sync.clearPrivateListeners();
    // Simulate initial private snapshot for seat 0 with old rack
    const oldSnap: PrivateSnapshot = {
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
        { ID: "old-1", Colour: 1, Rank: 5, IsJoker: false },
        { ID: "old-2", Colour: 2, Rank: 7, IsJoker: false },
      ],
      ownSeat: 0,
    };
    sync.onPrivateSnapshot(oldSnap);
    expect(sync.getLastPrivateBySeat(0)?.ownRack[0].ID).toBe("old-1");

    // Simulate reconnect: server sends new PrivateSnapshot for same seat with new rack (not old)
    const newSnap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MeldOrDiscard",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: true, rackCount: 2 }],
      stockCount: 70,
      discardRow: [
        {
          Tile: { ID: "disc-1", Colour: 1, Rank: 7, IsJoker: false },
          IsOpeningDiscard: true,
          Index: 0,
        },
      ],
      tableMelds: [{ ID: "m1", Kind: "run", Tiles: [], JokerReps: {}, OwnerSeat: 0 }],
      winner: -1,
      ownRack: [
        { ID: "new-1", Colour: 3, Rank: 9, IsJoker: false },
        { ID: "new-2", Colour: 4, Rank: 10, IsJoker: false },
      ],
      ownSeat: 0,
    };

    let rehydrated: PrivateSnapshot | null = null;
    const unsub = sync.subscribePrivateSnapshot((s) => {
      rehydrated = s;
    });
    // The subscribe should have immediately replayed oldSnap, so rehydrated is old
    expect((rehydrated as unknown as PrivateSnapshot).ownRack[0].ID).toBe("old-1");
    // Now simulate server sending new snapshot after reconnect
    sync.onPrivateSnapshot(newSnap);
    expect((rehydrated as unknown as PrivateSnapshot).ownRack[0].ID).toBe("new-1");
    expect((rehydrated as unknown as PrivateSnapshot).ownRack.length).toBe(2);
    expect(sync.getLastPrivateBySeat(0)?.ownRack[0].ID).toBe("new-1");
    // Ensure old IDs are not retained
    expect(sync.getLastPrivateBySeat(0)?.ownRack.some((t) => t.ID === "old-1")).toBe(false);
    unsub();
  });
});
