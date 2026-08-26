import { test, expect } from "@playwright/test";

test("Day 2 — scaffold shows Phaser canvas and no console errors", async ({ page }) => {
  const errors: string[] = [];
  page.on("pageerror", (e) => errors.push(e.message));
  await page.goto("/");
  await expect(page.locator("#game")).toBeVisible();
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  const hasPhaser = await page.evaluate(() => typeof (window as unknown as Record<string, unknown>).Phaser !== "undefined");
  expect(hasPhaser).toBeTruthy();
  await page.waitForTimeout(800);
  expect(errors).toHaveLength(0);
});

test("Day 2 — Preload complete log", async ({ page }) => {
  const logs: string[] = [];
  page.on("console", (m) => logs.push(m.text()));
  await page.goto("/");
  await page.waitForTimeout(1000);
  const hasPreloadLog = logs.some((l) => l.includes("Preload complete"));
  expect(hasPreloadLog).toBeTruthy();
});

test("Day 4 — uses entire width on desktop and no top info bar", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game")).toBeVisible();
  const hasNoInfo = await page.locator("#info").count();
  expect(hasNoInfo).toBe(0);
  const gameBox = await page.locator("#game").boundingBox();
  expect(gameBox).toBeTruthy();
  // On desktop, game should use entire viewport width (flex:1, no max-width)
  const viewport = page.viewportSize();
  if (viewport && viewport.width >= 1024) {
    expect(gameBox!.width).toBeGreaterThan(1000);
  }
});
