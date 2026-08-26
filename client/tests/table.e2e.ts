import { test, expect } from '@playwright/test';

test('Day 13 — Table board shows 3 melds with no scroll', async ({ page }) => {
	await page.goto('/demo/sync');
	await expect(page.getByText('ETALĂRI PE MASĂ').first()).toBeVisible();
	await expect(page.getByText('66 pct').first()).toBeVisible();
	await expect(page.getByText('53 pct').first()).toBeVisible();
	await expect(page.getByText('55 pct').first()).toBeVisible();
	// Check that table has no horizontal scroll
	const hasScroll = await page.evaluate(() => {
		const el = document.querySelector('div.rounded-\\[18px\\]');
		if (!el) return null;
		return (el as HTMLElement).scrollWidth > (el as HTMLElement).clientWidth;
	});
	expect(hasScroll).toBe(false);
});
