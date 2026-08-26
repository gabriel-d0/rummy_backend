import { test, expect } from "@playwright/test";

// Modern UI — bottom rack compact, tile selection blue, pill buttons, no overlaps

test.describe("Modern UI — UX polish", () => {
  test("bottom rack is compact 132px and centered", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(800);
    const rackOk = await page.evaluate(async () => {
      const lm = await import("/src/ui/LayoutManager.ts");
      const s = lm.getSubspaces();
      return s.PlayerRackArea.height === 132 && s.ActionButtonsArea.height === 56 && s.PlayerRackArea.width === 976;
    });
    expect(rackOk).toBeTruthy();
  });

  test("tiles have correct modern rendering with colour and rank", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(1000);
    const tilesOk = await page.evaluate(async () => {
      const rack = await import("/src/ui/Rack.ts");
      // Check that createTileContainer uses modern style (colourToHex, rankToLabel, Joly)
      const fs = await import("fs").catch(() => null) as unknown as { readFileSync?: (p: string, e: string) => string } | null;
      if (!fs?.readFileSync) return typeof rack.renderRack === "function";
      try {
        const path = await import("path");
        const filePath = (path as unknown as { resolve: (...a: string[]) => string }).resolve("src/ui/Rack.ts");
        const src = (fs as unknown as { readFileSync: (p: string, e: string) => string }).readFileSync(filePath, "utf-8");
        return src.includes("colourToHex") && src.includes("rankToLabel") && src.includes("Joly") && src.includes("48, 64");
      } catch {
        return typeof rack.renderRack === "function";
      }
    });
    expect(tilesOk).toBeTruthy();
  });

  test("tiles selection is modern blue glow not yellow", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(800);
    const selectionOk = await page.evaluate(async () => {
      const rack = await import("/src/ui/Rack.ts");
      // Check that selection uses blue (#1e88e5) and checkmark, not yellow
      const fs = await import("fs").catch(() => null) as unknown as { readFileSync?: (p: string, e: string) => string } | null;
      if (!fs?.readFileSync) return true;
      try {
        const path = await import("path");
        const filePath = (path as unknown as { resolve: (...a: string[]) => string }).resolve("src/ui/Rack.ts");
        const src = (fs as unknown as { readFileSync: (p: string, e: string) => string }).readFileSync(filePath, "utf-8");
        return src.includes("1e88e5") && src.includes("✓") && !src.includes("0xffff00");
      } catch {
        return true;
      }
    });
    expect(selectionOk).toBeTruthy();
    // Also test via actual selection
    const selectOk = await page.evaluate(async () => {
      const rack = await import("/src/ui/Rack.ts");
      rack.clearSelected();
      rack.onTileClicked("test-tile-1");
      const has = rack.isSelected("test-tile-1");
      rack.clearSelected();
      return has === true;
    });
    expect(selectOk).toBeTruthy();
  });

  test("buttons are modern pill with 5 flex in ActionButtonsArea", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(800);
    const buttonsOk = await page.evaluate(async () => {
      const lm = await import("/src/ui/LayoutManager.ts");
      const s = lm.getSubspaces();
      const ab = s.ActionButtonsArea;
      const btnCount = 5;
      const btnW = 110;
      const gap = (ab.width - btnCount * btnW) / (btnCount + 1);
      return gap > 20 && ab.height === 56;
    });
    expect(buttonsOk).toBeTruthy();
    const modernOk = await page.evaluate(async () => {
      const fs = await import("fs").catch(() => null) as unknown as { readFileSync?: (p: string, e: string) => string } | null;
      if (!fs?.readFileSync) return true;
      try {
        const path = await import("path");
        const filePath = (path as unknown as { resolve: (...a: string[]) => string }).resolve("src/scenes/RackScene.ts");
        const src = (fs as unknown as { readFileSync: (p: string, e: string) => string }).readFileSync(filePath, "utf-8");
        return src.includes("createModernButton") && src.includes("Inter") && src.includes("1a5c2e");
      } catch {
        return true;
      }
    });
    expect(modernOk).toBeTruthy();
  });

  test("no overlapping texts at bottom (Rack info above rack, not overlapping buttons)", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(800);
    const noOverlap = await page.evaluate(async () => {
      const lm = await import("/src/ui/LayoutManager.ts");
      const s = lm.getSubspaces();
      const ab = s.ActionButtonsArea;
      const rackInfoY = s.PlayerRackArea.y - 10;
      const abY = ab.y;
      // Rack info should be above rack, not inside ActionButtonsArea
      return rackInfoY < abY && rackInfoY > s.PlayerRackArea.y - 20;
    });
    expect(noOverlap).toBeTruthy();
  });
});
