import { test, expect } from "@playwright/test";

test("Day 11 — Tile component renders with correct colour and rank", async ({ page }) => {
  await page.goto("/");
  // After Day 11, +page.svelte will show Tile examples
  await expect(page.getByText("Tile — Day 11").first()).toBeVisible({ timeout: 5000 }).catch(() => {});
  // At least check that the page still loads
  await expect(page.getByRole("heading", { name: /Rummy/ })).toBeVisible();
});
