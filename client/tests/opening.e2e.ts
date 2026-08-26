import { test, expect } from '@playwright/test';

test('Day 30 — Opening discard 15→14 via OpClientDiscard 2', async ({ page }) => {
	await page.goto('/demo/opening');
	await expect(page.getByText('Opening Discard Demo — Day 30')).toBeVisible();

	// Set Opening Host 15 — should enable discard
	await page.getByRole('button', { name: 'Set Opening Host 15' }).click();
	await expect(page.getByTestId('opening-info')).toContainText('isOpening:true');
	await expect(page.getByTestId('opening-info')).toContainText('rack:15');
	await expect(page.getByText('Mâna ta • 15 cărți')).toBeVisible();
	// Discard button should be enabled (amber) when selected 1
	const discardBtn = page.getByTestId('discard-btn');
	await expect(discardBtn).toBeVisible();
	await expect(discardBtn).toBeDisabled(); // no selection yet

	// Select one tile — find Tile buttons in rack
	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	await tileButtons.first().click();
	await expect(discardBtn).toBeEnabled();

	// Click discard
	await discardBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":2');
	await expect(page.getByTestId('last-sent')).toContainText('tile-1');
	// Simulate server after discard 15→14
	await page.getByRole('button', { name: 'Simulate 15→14' }).click();
	await expect(page.getByTestId('opening-info')).toContainText('rack:14');
	await expect(page.getByText('Mâna ta • 14 cărți')).toBeVisible();

	// Set Opening Guest 15 — but ownSeat 1, currentSeat 1 → should still enable for that seat
	// For this page, we are seat 1 when setOpeningGuest15
	await page.getByRole('button', { name: 'Set Opening Guest 15' }).click();
	await expect(page.getByTestId('opening-info')).toContainText('rack:15');
	await expect(page.getByTestId('opening-info')).toContainText('ownSeat:1');
	await expect(page.getByTestId('opening-info')).toContainText('currentSeat:1');
	await expect(discardBtn).toBeDisabled(); // no selection
	await tileButtons.first().click();
	await expect(discardBtn).toBeEnabled();

	// Not your turn — should disable
	await page.getByRole('button', { name: 'Set Not Your Turn' }).click();
	await expect(page.getByTestId('opening-info')).toContainText('ownSeat:0');
	await expect(page.getByTestId('opening-info')).toContainText('currentSeat:1');
	// Need to clear selection first (discard still disabled because not my turn)
	// Click a tile to select, but canDiscard false due to not my turn
	await tileButtons.first().click();
	await expect(discardBtn).toBeDisabled();
});
