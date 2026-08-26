import { test, expect } from '@playwright/test';

test('Day 3 — Tailwind and design tokens visible', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByRole('heading', { name: /Rummy/ })).toBeVisible();
	await expect(page.getByText(/Day 12/).first()).toBeVisible();
});
