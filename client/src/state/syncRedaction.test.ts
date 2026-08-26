import { describe, it, expect, beforeEach, vi } from "vitest";
import {
  checkNoLeak,
  onPrivateSnapshot,
  onPublicSnapshot,
  clearPrivateListeners,
  clearPublicListeners,
} from "./sync";
import type { PrivateSnapshot, PublicSnapshot } from "./snapshot";

// Day 32: Redaction check (client) — json.stringify(PublicSnapshot) string search for OwnRack IDs

describe("Sync redaction — Day 32", () => {
  beforeEach(() => {
    clearPrivateListeners();
    clearPublicListeners();
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: () => null,
        setItem: () => {},
        removeItem: () => {},
      } as unknown as Storage;
    }
  });

  it("checkNoLeak returns true when no private IDs in publicJson", () => {
    const publicJson = JSON.stringify({ v: 1, players: [{ rackCount: 14 }], stockCount: 77 });
    expect(checkNoLeak(publicJson, ["alice-secret", "bob-secret"])).toBe(true);
  });

  it("checkNoLeak returns false when private ID leaks", () => {
    const publicJson = JSON.stringify({
      v: 1,
      players: [{ rackCount: 14 }],
      ownRack: [{ ID: "alice-secret" }],
    });
    expect(checkNoLeak(publicJson, ["alice-secret"])).toBe(false);
  });

  it("checkNoLeak handles empty privateIds", () => {
    const publicJson = JSON.stringify({ v: 1, stockCount: 77 });
    expect(checkNoLeak(publicJson, [])).toBe(true);
  });

  it("onPublicSnapshot logs no leak when privateIds not in public", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const priv: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 3 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "alice-secret-1", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    onPrivateSnapshot(priv);
    const pub: PublicSnapshot = {
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
    onPublicSnapshot(pub);
    const hasNoLeakLog = consoleSpy.mock.calls.some((args) => String(args[0]).includes("no leak"));
    expect(hasNoLeakLog).toBe(true);
    const hasLeakedLog = consoleSpy.mock.calls.some((args) => String(args[0]).includes("LEAKED"));
    expect(hasLeakedLog).toBe(false);
    consoleSpy.mockRestore();
  });

  it("onPublicSnapshot logs LEAKED when public contains private ID (simulated)", () => {
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const priv: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 1 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: [{ ID: "leaked-id-xyz", Colour: 1, Rank: 5, IsJoker: false }],
      ownSeat: 0,
    };
    onPrivateSnapshot(priv);
    // Create a public snapshot that incorrectly contains the private ID (simulate leak)
    const pubWithLeak = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 1 }],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      // leak field that shouldn't be there
      leaked: "leaked-id-xyz",
    } as unknown as PublicSnapshot;
    onPublicSnapshot(pubWithLeak);
    const hasLeakedLog = consoleSpy.mock.calls.some((args) => String(args[0]).includes("LEAKED"));
    expect(hasLeakedLog).toBe(true);
    consoleSpy.mockRestore();
  });

  it("checkNoLeak mirrors snapshot.ts helper", () => {
    const publicJson = JSON.stringify({ v: 1, stockCount: 77, discardRow: [] });
    expect(checkNoLeak(publicJson, ["private-1"])).toBe(true);
    const publicWithPrivate = JSON.stringify({
      v: 1,
      stockCount: 77,
      discardRow: [{ Tile: { ID: "private-1" } }],
    });
    // Note: discardRow is public, so if private ID appears there it would be a leak only if it's ownRack ID
    // For this test, we just verify the string search works
    expect(checkNoLeak(publicWithPrivate, ["private-1"])).toBe(false);
  });
});
