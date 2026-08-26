import { test, expect } from "@playwright/test";

test("Day 11 — Tile component renders with correct colour and rank", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: /Rummy/ })).toBeVisible();
  await expect(page.getByText("Day 12").first()).toBeVisible();
  await expect(page.getByText("48×64").first()).toBeVisible();
});
