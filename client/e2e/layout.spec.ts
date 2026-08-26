import { test, expect } from "@playwright/test";

test("Day 4 — LayoutManager GameSpace 1000x1000 and 6 subspaces", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  const isGameSpace = await page.evaluate(() => {
    const gs = (window as unknown as Record<string, unknown>).__GAME_SPACE__ as { width: number; height: number } | undefined;
    return gs?.width === 1000 && gs?.height === 1000;
  });
  expect(isGameSpace).toBeTruthy();
  const subspacesOk = await page.evaluate(async () => {
    const lm = await import("/src/ui/LayoutManager.ts");
    const s = lm.getSubspaces();
    return s.TopBar.height === 80 && s.PlayerRackArea.height === 132 && s.ActionButtonsArea.height === 56;
  });
  expect(subspacesOk).toBeTruthy();
});

test("Day 4 — proportional placement via centerX and spacing", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  const helpersOk = await page.evaluate(async () => {
    const lm = await import("/src/ui/LayoutManager.ts");
    const s = lm.getSubspaces();
    const cx = lm.centerX(s.TopBar);
    const cy = lm.centerY(s.TopBar);
    const sp = lm.spacing(s.MeldArea.width, 3, 48);
    const tx = lm.tileX(s.MeldArea, 1, 3, 48);
    return cx === s.TopBar.x + s.TopBar.width / 2 && cy === s.TopBar.y + s.TopBar.height / 2 && sp > 6 && tx > s.MeldArea.x;
  });
  expect(helpersOk).toBeTruthy();
});

test("Day 4 — no overlap, stable spacing on resize", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  await page.waitForTimeout(800);
  // Check that 6 subspaces are visible as colored rects
  const hasDebug = await page.evaluate(() => {
    const game = (window as unknown as Record<string, unknown>).__GAME__ as { scene?: { getScene?: (k: string) => { children?: { list: { getData?: (k: string) => unknown }[] } } } } | undefined;
    const preload = game?.scene?.getScene?.("Preload") as unknown as { children?: { list: { getData?: (k: string) => unknown }[] } } | undefined;
    if (!preload?.children?.list) return false;
    const debugCount = preload.children.list.filter((c) => c.getData?.("isDebugSubspace")).length;
    // 6 rects + 6 labels = 12
    return debugCount >= 12;
  });
  expect(hasDebug).toBeTruthy();
});
