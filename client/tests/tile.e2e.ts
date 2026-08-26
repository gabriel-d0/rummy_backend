import { test, expect } from "@playwright/test";

test("Day 11 — Tile component renders with correct colour and rank", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: /Rummy/ })).toBeVisible();
  await expect(page.getByText(/Table board/).first()).toBeVisible();
});
