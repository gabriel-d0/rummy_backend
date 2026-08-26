import { test, expect } from '@playwright/test';

test('Day 7 — smoke dev server serves 200 at / and SvelteKit loads', async ({ page }) => {
	const resp = await page.goto('/');
	expect(resp?.status()).toBe(200);
	await expect(page.getByText('Welcome to SvelteKit')).toBeVisible();
});

test('Day 7 — no console errors on load', async ({ page }) => {
	const errors: string[] = [];
	page.on('pageerror', (e) => errors.push(e.message));
	await page.goto('/');
	await page.waitForTimeout(800);
	expect(errors).toHaveLength(0);
});
