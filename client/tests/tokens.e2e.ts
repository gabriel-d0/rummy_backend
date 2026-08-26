import { test, expect } from '@playwright/test';

test('Day 3 — Tailwind and design tokens visible', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByRole('heading', { name: /Rummy.*Day 3/ })).toBeVisible();
	await expect(page.getByText('Tailwind & design tokens')).toBeVisible();
	await expect(page.getByText('Inter').first()).toBeVisible();
});
