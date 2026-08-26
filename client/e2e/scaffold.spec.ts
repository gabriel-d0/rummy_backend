import { test, expect } from "@playwright/test";

test("Day 2 — scaffold shows Phaser canvas and no console errors", async ({ page }) => {
  const errors: string[] = [];
  page.on("pageerror", (e) => errors.push(e.message));
  await page.goto("/");
  await expect(page.locator("#game")).toBeVisible();
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  await expect(page.locator("#info")).toContainText("Day");
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
