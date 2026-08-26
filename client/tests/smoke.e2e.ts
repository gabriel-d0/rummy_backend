import { test, expect } from '@playwright/test';

test('Day 16 — visual layout TopBar + TableBoard + Rack at 1280x800', async ({ page }) => {
	const resp = await page.goto('/');
	expect(resp?.status()).toBe(200);
	await expect(page.getByText('REMI ETALAT').first()).toBeVisible();
	await expect(page.getByText('ETALĂRI PE MASĂ').first()).toBeVisible();
	await expect(page.getByText('Mâna ta').first()).toBeVisible();
});

test('Day 7 — no console errors on load', async ({ page }) => {
	const errors: string[] = [];
	page.on('pageerror', (e) => errors.push(e.message));
	await page.goto('/');
	await page.waitForTimeout(800);
	expect(errors).toHaveLength(0);
});
