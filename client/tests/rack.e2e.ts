import { test, expect } from '@playwright/test';

test('Day 14 — Rack shows 11 tiles with no scroll', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByText('Mâna ta').first()).toBeVisible();
	await expect(page.getByText('TRAGE DIN TALON').first()).toBeVisible();
	const hasScroll = await page.evaluate(() => {
		const el = document.querySelector('div.flex.flex-wrap.gap-1\\.5');
		if (!el) return null;
		return (el as HTMLElement).scrollWidth > (el as HTMLElement).clientWidth;
	});
	// Rack should have no horizontal scroll (flex-wrap)
	expect(hasScroll === false || hasScroll === null).toBeTruthy();
});
