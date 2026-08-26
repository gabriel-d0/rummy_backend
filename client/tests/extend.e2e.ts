import { test, expect } from '@playwright/test';

test('Day 37 — Extend meld HasOpened selected>=1 meldId via TableBoard OpClientExtendMeld 8', async ({
	page
}) => {
	await page.goto('/demo/extend');
	await expect(page.getByText('Extend Meld Demo — Day 37')).toBeVisible();

	const extendBtn = page.getByTestId('extend-btn');
	await expect(extendBtn).toBeDisabled();

	// Set HasOpened with run 5-6-7
	await page.getByRole('button', { name: 'Set HasOpened + Run 5-6-7' }).click();
	await expect(page.getByTestId('extend-info')).toContainText('hasOpened:true');
	await expect(extendBtn).toBeDisabled(); // no meld selected, no tile selected

	// Click meld on TableBoard
	const meldBtn = page.getByTestId('meld-m1');
	await expect(meldBtn).toBeVisible();
	await meldBtn.click();
	await expect(page.getByTestId('extend-info')).toContainText('selectedMeld:m1');
	await expect(extendBtn).toBeDisabled(); // no tile selected yet

	// Select 1 tile in Rack
	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	// The rack has tile-8 (rack-8) which can extend run 5-6-7 to 5-6-7-8
	// Find any tile and click; easiest click first tile (rack-8)
	await tileButtons.first().click();
	await expect(extendBtn).toBeEnabled();

	// Click extend → sends op 8
	await extendBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":8');
	await expect(page.getByTestId('last-sent')).toContainText('m1');
	await expect(page.getByTestId('last-sent')).toContainText('tileIds');

	// Set NotOpened → should disable even with meld selected and tile selected
	await page.getByRole('button', { name: 'Set HasOpened + Run 5-6-7' }).click();
	await page.getByTestId('meld-m1').click();
	await tileButtons.first().click();
	await expect(extendBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set NotOpened' }).click();
	await expect(page.getByTestId('extend-info')).toContainText('hasOpened:false');
	await expect(extendBtn).toBeDisabled();
});
