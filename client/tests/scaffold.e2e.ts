import { test, expect } from '@playwright/test';

test('Day 2 — SvelteKit skeleton shows Rummy and no console errors', async ({ page }) => {
	const errors: string[] = [];
	page.on('pageerror', (e) => errors.push(e.message));
	await page.goto('/');
	await expect(page.getByRole('heading', { name: /Rummy.*SvelteKit Day 2/ })).toBeVisible();
	await expect(page.getByText('Vite + SvelteKit + TypeScript')).toBeVisible();
	await expect(page.getByText('Welcome to SvelteKit')).toBeVisible();
	await page.waitForTimeout(500);
	expect(errors).toHaveLength(0);
});
