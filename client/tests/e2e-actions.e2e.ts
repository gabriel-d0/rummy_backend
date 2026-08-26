import { test, expect } from '@playwright/test';

test('Day 41 — E2E actions 2 browsers alice/bob draw/discard/meld/win via requestId', async ({
	browser
}) => {
	const aliceContext = await browser.newContext();
	const bobContext = await browser.newContext();
	const alicePage = await aliceContext.newPage();
	const bobPage = await bobContext.newPage();

	// Alice: draw stock MustDraw → op 3 with requestId
	await alicePage.goto('/demo/draw');
	await alicePage.getByRole('button', { name: 'Set MustDraw My Turn' }).click();
	await alicePage.getByTestId('draw-btn').click();
	await expect(alicePage.getByTestId('last-sent')).toContainText('"op":3');
	const aliceDrawSent = await alicePage.getByTestId('last-sent').textContent();
	expect(aliceDrawSent).toContain('requestId');
	expect(aliceDrawSent).toContain('mock-match');

	// Bob: normal discard MeldOrDiscard → op 2
	await bobPage.goto('/demo/discard');
	await bobPage.getByRole('button', { name: 'Set MeldOrDiscard My Turn 0' }).click();
	const bobTiles = bobPage.locator('div.flex.min-h-\\[110px\\] button');
	await bobTiles.first().click();
	await bobPage.getByTestId('discard-btn').click();
	await expect(bobPage.getByTestId('last-sent')).toContainText('"op":2');
	const bobDiscardSent = await bobPage.getByTestId('last-sent').textContent();
	expect(bobDiscardSent).toContain('requestId');
	// requestIds should be distinct
	expect(aliceDrawSent).not.toEqual(bobDiscardSent);

	// Alice: meld initial !HasOpened → op 6
	await alicePage.goto('/demo/meld');
	await alicePage.getByRole('button', { name: 'Set NotOpened MeldOrDiscard' }).click();
	const aliceMeldTiles = alicePage.locator('div.flex.min-h-\\[110px\\] button');
	await aliceMeldTiles.nth(0).click();
	await aliceMeldTiles.nth(1).click();
	await aliceMeldTiles.nth(2).click();
	await aliceMeldTiles.nth(3).click();
	await aliceMeldTiles.nth(4).click();
	await aliceMeldTiles.nth(5).click();
	await alicePage.getByTestId('meld-initial-btn').click();
	await expect(alicePage.getByTestId('last-sent')).toContainText('"op":6');
	await expect(alicePage.getByTestId('last-sent')).toContainText('melds');

	// Bob: meld new HasOpened → op 7
	await bobPage.goto('/demo/meld-new');
	await bobPage.getByRole('button', { name: 'Set HasOpened MeldOrDiscard' }).click();
	const bobMeldTiles = bobPage.locator('div.flex.min-h-\\[110px\\] button');
	await bobMeldTiles.nth(0).click();
	await bobMeldTiles.nth(1).click();
	await bobMeldTiles.nth(2).click();
	await bobPage.getByTestId('meld-initial-btn').click();
	await expect(bobPage.getByTestId('last-sent')).toContainText('"op":7');

	// Alice: win → RoundComplete Winner 0
	await alicePage.goto('/demo/winner');
	await alicePage.getByRole('button', { name: 'Set Winner 0' }).click();
	await expect(alicePage.getByTestId('winner-overlay')).toBeVisible();
	await expect(alicePage.getByTestId('winner-text')).toContainText('Winner 0');

	// Bob: win → Winner 0 as well (public winner)
	await bobPage.goto('/demo/winner');
	await bobPage.getByRole('button', { name: 'Set Winner 0' }).click();
	await expect(bobPage.getByTestId('winner-overlay')).toBeVisible();
	await expect(bobPage.getByTestId('winner-text')).toContainText('Winner 0');

	await aliceContext.close();
	await bobContext.close();
});
