import { test, expect } from '@playwright/test';

test('Day 35 — Meld initial !HasOpened selected>=3 sends OpClientMeldInitial 6 total>=50 run', async ({
	page
}) => {
	await page.goto('/demo/meld');
	await expect(page.getByText('Meld Initial Demo — Day 35')).toBeVisible();

	const meldBtn = page.getByTestId('meld-initial-btn');
	// Initially no store → fallback selected>=3? but no selection → disabled
	await expect(meldBtn).toBeDisabled();

	// Set NotOpened MeldOrDiscard — should allow meld when selected 3+
	await page.getByRole('button', { name: 'Set NotOpened MeldOrDiscard' }).click();
	await expect(page.getByTestId('meld-info')).toContainText('hasOpened:false');
	await expect(meldBtn).toBeDisabled(); // no selection yet
	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	// Select 3 tiles
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await tileButtons.nth(2).click();
	await expect(meldBtn).toBeEnabled();

	// Select more to reach 50 — choose 6 to ensure 50 via run 8-13
	await tileButtons.nth(3).click();
	await tileButtons.nth(4).click();
	await tileButtons.nth(5).click();
	await expect(meldBtn).toBeEnabled();

	// Click meld initial → sends op 6
	await meldBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":6');
	await expect(page.getByTestId('last-sent')).toContainText('melds');
	await expect(page.getByTestId('last-sent')).toContainText('tileIds');

	// Set HasOpened → should disable meld initial (needs !HasOpened)
	await page.getByRole('button', { name: 'Set HasOpened' }).click();
	await expect(page.getByTestId('meld-info')).toContainText('hasOpened:true');
	// Need to select again after snapshot reset (clears? but Rack selected persists, now 6 selected but hasOpened true → disabled)
	await expect(meldBtn).toBeDisabled();
});
