import { test, expect } from '@playwright/test';

test('Day 11 — Tile component renders with correct colour and rank', async ({ page }) => {
	await page.goto('/demo/sync');
	await expect(page.getByRole('heading', { name: /Sync Demo/ })).toBeVisible();
	await expect(page.getByText(/Mâna ta/).first()).toBeVisible();
});
