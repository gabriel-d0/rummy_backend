import { test, expect } from '@playwright/test';

test('Day 40 — Visual actions start→opening→draw→meld→extend→prev/pickup→replace→discard→win', async ({
	page
}) => {
	// Start → Op 1
	await page.goto('/demo/start');
	await page.getByRole('button', { name: 'Set Waiting Host 2p' }).click();
	await page.getByTestId('start-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":1');

	// Opening discard 15→14 → Op 2
	await page.goto('/demo/opening');
	await page.getByRole('button', { name: 'Set Opening Host 15' }).click();
	const tileBtnsOpening = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsOpening.first().click();
	await page.getByTestId('discard-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":2');
	await page.getByRole('button', { name: 'Simulate 15→14' }).click();
	await expect(page.getByTestId('opening-info')).toContainText('rack:14');

	// Draw stock MustDraw → Op 3
	await page.goto('/demo/draw');
	await page.getByRole('button', { name: 'Set MustDraw My Turn' }).click();
	await page.getByTestId('draw-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":3');

	// Meld initial !HasOpened → Op 6
	await page.goto('/demo/meld');
	await page.getByRole('button', { name: 'Set NotOpened MeldOrDiscard' }).click();
	const tileBtnsMeld = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsMeld.nth(0).click();
	await tileBtnsMeld.nth(1).click();
	await tileBtnsMeld.nth(2).click();
	await tileBtnsMeld.nth(3).click();
	await tileBtnsMeld.nth(4).click();
	await tileBtnsMeld.nth(5).click();
	await page.getByTestId('meld-initial-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":6');

	// Meld new HasOpened → Op 7
	await page.goto('/demo/meld-new');
	await page.getByRole('button', { name: 'Set HasOpened MeldOrDiscard' }).click();
	const tileBtnsMeldNew = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsMeldNew.nth(0).click();
	await tileBtnsMeldNew.nth(1).click();
	await tileBtnsMeldNew.nth(2).click();
	await page.getByTestId('meld-initial-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":7');

	// Extend HasOpened → Op 8
	await page.goto('/demo/extend');
	await page.getByRole('button', { name: 'Set HasOpened + Run 5-6-7' }).click();
	await page.getByTestId('meld-m1').click();
	const tileBtnsExtend = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsExtend.first().click();
	await page.getByTestId('extend-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":8');

	// Draw previous HasOpened → Op 4
	await page.goto('/demo/draw-previous');
	await page.getByRole('button', { name: 'Set HasOpened + Discard' }).click();
	await page.getByTestId('draw-prev-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":4');

	// Pickup for meld → Op 5
	await page.goto('/demo/pickup');
	await page.getByRole('button', { name: 'Set Valid Pickup (HasOpened + 2 discards)' }).click();
	await page.getByTestId('discard-tile-0').click();
	const tileBtnsPickup = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsPickup.nth(0).click();
	await tileBtnsPickup.nth(1).click();
	await page.getByTestId('pickup-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":5');

	// Replace joker → Op 9
	await page.goto('/demo/joker');
	await page.getByRole('button', { name: 'Set HasOpened + Joker Run 5-6-J' }).click();
	await page.getByTestId('joker-m1-joker-1').click();
	const tileBtnsJoker = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsJoker.nth(0).click();
	await tileBtnsJoker.nth(1).click();
	await tileBtnsJoker.nth(2).click();
	await page.getByTestId('replace-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":9');

	// Normal discard MeldOrDiscard → Op 2, 0→1
	await page.goto('/demo/discard');
	await page.getByRole('button', { name: 'Set MeldOrDiscard My Turn 0' }).click();
	const tileBtnsDiscard = page.locator('div.flex.min-h-\\[110px\\] button');
	await tileBtnsDiscard.first().click();
	await page.getByTestId('discard-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":2');
	await page.getByRole('button', { name: 'Simulate Discard 0→1' }).click();
	await expect(page.getByTestId('discard-info')).toContainText('currentSeat:1');

	// Win → RoundComplete Winner 0
	await page.goto('/demo/winner');
	await page.getByRole('button', { name: 'Set Winner 0' }).click();
	await expect(page.getByTestId('winner-overlay')).toBeVisible();
	await expect(page.getByTestId('winner-text')).toContainText('Winner 0');
});
