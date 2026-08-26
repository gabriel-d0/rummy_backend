import { test, expect } from '@playwright/test';

test('Day 15 — TopBar shows REMI ETALAT and MASA 1', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByText('REMI ETALAT').first()).toBeVisible();
	await expect(page.getByText('MASA 1').first()).toBeVisible();
	await expect(page.getByText('JOC NOU').first()).toBeVisible();
});
