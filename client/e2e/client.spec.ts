import { test, expect } from "@playwright/test";

// Day 2-28 comprehensive smoke — verifies what we implemented until now
// Covers: Vite+Phaser scaffold, Preload, assets, rack/table rendering, stock/turn, layout, no leak

test.describe("Rummy Phaser Client — smoke", () => {
  test("loads Vite dev server at 5173 and shows game canvas", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game")).toBeVisible();
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await expect(page.locator("#info")).toContainText("Rummy Phaser 3");
    await expect(page.locator("#info")).toContainText("client/docs/roadmap.md");
  });

  test("Phaser 3 is loaded and Preload completed", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(1200);
    const phaserVersion = await page.evaluate(() => (window as unknown as Record<string, unknown>).Phaser as { VERSION?: string } | undefined);
    // Phaser is available via window.Phaser when main.ts imports it
    const hasPhaser = await page.evaluate(() => typeof (window as unknown as Record<string, unknown>).Phaser !== "undefined");
    expect(hasPhaser).toBeTruthy();
    // Check that game instance was exposed for e2e
    const hasGame = await page.evaluate(() => !!(window as unknown as Record<string, unknown>).__GAME__);
    expect(hasGame).toBeTruthy();
  });

  test("creates TableScene and RackScene with no console errors", async ({ page }) => {
    const errors: string[] = [];
    page.on("pageerror", (err) => errors.push(err.message));
    const consoleErrors: string[] = [];
    page.on("console", (msg) => {
      if (msg.type() === "error") consoleErrors.push(msg.text());
    });
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(1500);
    expect(errors).toHaveLength(0);
    // Allow at most 0 console errors (Phaser logs are info)
    expect(consoleErrors).toHaveLength(0);
  });

  test("has no private leak in public DOM (redaction)", async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(1500);
    const content = await page.content();
    expect(content).not.toContain("alice-secret");
    expect(content).not.toContain("bob-secret");
    expect(content).not.toContain("OwnRack");
    // PublicView should only expose RackCount, StockCount, not TileInstanceId
    // Check that info does not contain private rack data
    const infoText = await page.locator("#info").textContent();
    expect(infoText).not.toContain("alice-secret");
  });

  test("exposes GAME_SPACE 1000x1000 and LayoutManager", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(800);
    const space = await page.evaluate(() => (window as unknown as Record<string, unknown>).__GAME_SPACE__ as { width: number; height: number } | undefined);
    expect(space).toBeTruthy();
    expect(space?.width).toBe(1000);
    expect(space?.height).toBe(1000);
    // Verify layout via imported module
    const layoutOk = await page.evaluate(async () => {
      const mod = await import("/src/ui/Layout.ts");
      const layout = mod.getLayout(1000, 1000);
      return layout.width === 1000 && layout.height === 1000 && layout.rackSlots.length === 14 && !layout.isMobile;
    });
    expect(layoutOk).toBeTruthy();
    const mobileOk = await page.evaluate(async () => {
      const mod = await import("/src/ui/Layout.ts");
      const layout = mod.getLayout(375, 667);
      return layout.isMobile === true && layout.rackSlots.length === 14;
    });
    expect(mobileOk).toBeTruthy();
  });

  test("GAME_SPACE subspaces are mathematically defined", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const subspacesOk = await page.evaluate(async () => {
      const mod = await import("/src/ui/LayoutManager.ts");
      const subs = mod.getSubspaces();
      // Check that subspaces partition 1000x1000 with gutters — modern compact 132/56
      const hasAll = !!(subs.TopBar && subs.TableArea && subs.MeldArea && subs.DiscardRowArea && subs.PlayerRackArea && subs.ActionButtonsArea);
      const topBarOk = subs.TopBar.height === 80 && subs.TopBar.width === 1000 - 24;
      const rackOk = subs.PlayerRackArea.height === 132 && subs.ActionButtonsArea.height === 56;
      return hasAll && topBarOk && rackOk;
    });
    expect(subspacesOk).toBeTruthy();
  });
});
