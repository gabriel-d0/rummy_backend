import { test, expect } from '@playwright/test';

test('Day 36 — Meld new HasOpened selected>=3 sends OpClientMeldNew 7', async ({ page }) => {
	await page.goto('/demo/meld-new');
	await expect(page.getByText('Meld New Demo — Day 36')).toBeVisible();

	const meldBtn = page.getByTestId('meld-initial-btn');
	await expect(meldBtn).toBeDisabled();

	// Set HasOpened MeldOrDiscard — should allow meld new when selected 3+
	await page.getByRole('button', { name: 'Set HasOpened MeldOrDiscard' }).click();
	await expect(page.getByTestId('meld-new-info')).toContainText('hasOpened:true');
	await expect(meldBtn).toBeDisabled(); // no selection yet
	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await tileButtons.nth(2).click();
	await expect(meldBtn).toBeEnabled();

	// Click meld new → sends op 7
	await meldBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":7');
	await expect(page.getByTestId('last-sent')).toContainText('melds');
	await expect(page.getByTestId('last-sent')).toContainText('tileIds');

	// Set NotOpened — now same button should send op 6 for initial, but still enabled with 3 selected? Need to re-select after snapshot
	await page.getByRole('button', { name: 'Set NotOpened' }).click();
	await expect(page.getByTestId('meld-new-info')).toContainText('hasOpened:false');
	// After snapshot change, selected cleared? Actually Rack selected persists, but we cleared after meld? After meld click we cleared selected, so now 0 selected
	await expect(meldBtn).toBeDisabled();
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await tileButtons.nth(2).click();
	await expect(meldBtn).toBeEnabled();
	// This should send op 6 (initial) not 7
	await meldBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":6');
});
