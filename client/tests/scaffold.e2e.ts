import { test, expect } from '@playwright/test';

test('Day 2 — SvelteKit skeleton shows Rummy and no console errors', async ({ page }) => {
	const errors: string[] = [];
	page.on('pageerror', (e) => errors.push(e.message));
	await page.goto('/');
	await expect(page.getByRole('heading', { name: /Rummy.*Day [23]/ })).toBeVisible();
	await expect(page.getByText('Tailwind & design tokens')).toBeVisible();
	await page.waitForTimeout(500);
	expect(errors).toHaveLength(0);
});
