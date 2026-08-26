import { test, expect } from '@playwright/test';

test('Day 34 — Normal discard MeldOrDiscard Discard 0→1 via OpClientDiscard 2', async ({
	page
}) => {
	await page.goto('/demo/discard');
	await expect(page.getByText('Normal Discard Demo — Day 34')).toBeVisible();

	const discardBtn = page.getByTestId('discard-btn');
	// Initially no store → fallback selected 1 would enable but we have no selection, so disabled
	await expect(discardBtn).toBeDisabled();

	// Set MeldOrDiscard My Turn 0 — should enable when selected 1
	await page.getByRole('button', { name: 'Set MeldOrDiscard My Turn 0' }).click();
	await expect(page.getByTestId('discard-info')).toContainText('isMeldOrDiscard:true');
	await expect(page.getByTestId('discard-info')).toContainText('currentSeat:0');
	await expect(discardBtn).toBeDisabled(); // no selection

	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	await tileButtons.first().click();
	await expect(discardBtn).toBeEnabled();

	// Click discard → sends op 2
	await discardBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":2');
	await expect(page.getByTestId('last-sent')).toContainText('tile-1');

	// Simulate server after discard: CurrentSeat 0→1
	await page.getByRole('button', { name: 'Simulate Discard 0→1' }).click();
	await expect(page.getByTestId('discard-info')).toContainText('currentSeat:1');
	await expect(page.getByTestId('discard-info')).toContainText('isMeldOrDiscard:false');

	// Not My Turn — should disable even with selection
	await page.getByRole('button', { name: 'Set Not My Turn 0≠1' }).click();
	await expect(page.getByTestId('discard-info')).toContainText('currentSeat:1');
	await expect(page.getByTestId('discard-info')).toContainText('ownSeat:0');
	await expect(discardBtn).toBeDisabled();
	await tileButtons.first().click();
	await expect(discardBtn).toBeDisabled();

	// MustDraw — should disable discard (needs MeldOrDiscard)
	await page.getByRole('button', { name: 'Set MeldOrDiscard My Turn 0' }).click();
	// selected still has tile-1 from previous step, now my turn MeldOrDiscard → enabled directly
	await expect(discardBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set MustDraw' }).click();
	await expect(discardBtn).toBeDisabled();
});
