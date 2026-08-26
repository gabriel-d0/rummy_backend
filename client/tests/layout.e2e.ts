import { test, expect } from "@playwright/test";

test("Day 16 — Visual layout TopBar + TableBoard + Rack + Jurnal at 1280x800 + 375x667", async ({ page }) => {
  // Desktop 1280x800
  await page.setViewportSize({ width: 1280, height: 800 });
  await page.goto("/");
  await expect(page.getByText("REMI ETALAT").first()).toBeVisible();
  await expect(page.getByText("ETALĂRI PE MASĂ").first()).toBeVisible();
  await expect(page.getByText("Mâna ta").first()).toBeVisible();
  await expect(page.getByText("JURNAL DE JOC").first()).toBeVisible();
  await page.screenshot({ path: "test-results/layout-1280.png", fullPage: true });

  // Mobile 375x667
  await page.setViewportSize({ width: 375, height: 667 });
  await page.goto("/");
  await expect(page.getByText("REMI ETALAT").first()).toBeVisible();
  await expect(page.getByText("Mâna ta").first()).toBeVisible();
  await page.screenshot({ path: "test-results/layout-375.png", fullPage: true });
});
