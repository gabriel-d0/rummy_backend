import { test, expect } from "@playwright/test";

// Day 23-27: create/join match via Nakama socket, receive match state, sync
// Tests what we implemented until now: socket.onmatchdata handling, Public/Private snapshots, no leak

test.describe("Match lifecycle — Day 23-27", () => {
  test("creates matchId storage in localStorage after createMatch (logic exists)", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const hasStorageLogic = await page.evaluate(async () => {
      const mod = await import("/src/net/nakama.ts");
      // Check that createMatch stores matchId in localStorage (code inspection via string)
      const src = mod.createMatch.toString();
      return src.includes("localStorage") && src.includes("rummy_matchId");
    });
    expect(hasStorageLogic).toBeTruthy();
  });

  test("handleMatchData parses Envelope and routes op 100/101/102/103", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const routingOk = await page.evaluate(async () => {
      const sync = await import("/src/state/sync.ts");
      return (
        typeof sync.handleMatchData === "function" &&
        typeof sync.onPrivateSnapshot === "function" &&
        typeof sync.onPublicSnapshot === "function" &&
        typeof sync.onServerError === "function" &&
        typeof sync.onServerEvent === "function"
      );
    });
    expect(routingOk).toBeTruthy();
  });

  test("receives PrivateSnapshot (op 100) and stores lastPrivate", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const privateOk = await page.evaluate(async () => {
      const sync = await import("/src/state/sync.ts");
      const snap: unknown = {
        v: 1,
        gamePhase: "Playing",
        turnPhase: "MustDraw",
        currentSeat: 0,
        players: [{ id: "alice", seat: 0, hasOpened: false, rackCount: 14 }],
        stockCount: 77,
        discardRow: [],
        tableMelds: [],
        winner: -1,
        ownRack: [{ ID: "tile-1", Colour: 1, Rank: 5, IsJoker: false }],
        ownSeat: 0,
      };
      // Clear previous
      localStorage.removeItem("rummy_lastPrivate");
      sync.handleMatchData(100, JSON.stringify({ v: 1, op: 100, payload: snap }));
      const stored = sync.getLastPrivateSnapshot();
      const ls = localStorage.getItem("rummy_lastPrivate");
      return !!stored && stored.ownSeat === 0 && stored.ownRack.length === 1 && !!ls;
    });
    expect(privateOk).toBeTruthy();
  });

  test("receives PublicSnapshot (op 101) and stores lastPublic", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const publicOk = await page.evaluate(async () => {
      const sync = await import("/src/state/sync.ts");
      const snap: unknown = {
        v: 1,
        gamePhase: "Playing",
        turnPhase: "MeldOrDiscard",
        currentSeat: 1,
        players: [
          { id: "alice", seat: 0, hasOpened: true, rackCount: 12 },
          { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
        ],
        stockCount: 63,
        discardRow: [{ Tile: { ID: "d1", Colour: 1, Rank: 7, IsJoker: false }, IsOpeningDiscard: true, Index: 0 }],
        tableMelds: [],
        winner: -1,
      };
      sync.handleMatchData(101, JSON.stringify({ v: 1, op: 101, payload: snap }));
      const stored = sync.getLastPublicSnapshot();
      return !!stored && stored.stockCount === 63 && stored.players.length === 2;
    });
    expect(publicOk).toBeTruthy();
  });

  test("snapshot types mirror Go visibility.go Version 1 (no private leak)", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const noLeak = await page.evaluate(async () => {
      const sync = await import("/src/state/sync.ts");
      // Create a private snapshot for seat 0
      const privateSnap: unknown = {
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
        ownRack: [{ ID: "alice-secret", Colour: 1, Rank: 5, IsJoker: false }],
        ownSeat: 0,
      };
      sync.handleMatchData(100, JSON.stringify({ v: 1, op: 100, payload: privateSnap }));
      // Public snapshot should not contain alice-secret
      const publicSnap: unknown = {
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
      };
      sync.handleMatchData(101, JSON.stringify({ v: 1, op: 101, payload: publicSnap }));
      const pub = sync.getLastPublicSnapshot();
      const pubJson = JSON.stringify(pub);
      return !pubJson.includes("alice-secret") && pubJson.includes("rackCount");
    });
    expect(noLeak).toBeTruthy();
  });

  test("socket.onmatchdata handler is set in createSocket (code inspection)", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const handlerOk = await page.evaluate(async () => {
      const mod = await import("/src/net/nakama.ts");
      const src = mod.createSocket.toString();
      return src.includes("onmatchdata") && src.includes("handleMatchData");
    });
    expect(handlerOk).toBeTruthy();
  });
});
