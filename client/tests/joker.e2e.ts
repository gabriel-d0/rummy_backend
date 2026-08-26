import { test, expect } from '@playwright/test';

test('Day 38 — Replace joker HasOpened via TableBoard joker click and Rack 3 tiles OpClientReplaceJoker 9', async ({
	page
}) => {
	await page.goto('/demo/joker');
	await expect(page.getByText('Joker Replace Demo — Day 38')).toBeVisible();

	const replaceBtn = page.getByTestId('replace-btn');
	await expect(replaceBtn).toBeDisabled();

	// Set HasOpened with Joker run 5-6-J
	await page.getByRole('button', { name: 'Set HasOpened + Joker Run 5-6-J' }).click();
	await expect(page.getByTestId('joker-info')).toContainText('hasOpened:true');
	await expect(replaceBtn).toBeDisabled(); // no joker selected, no tiles

	// Click joker on TableBoard
	const jokerBtn = page.getByTestId('joker-m1-joker-1');
	await expect(jokerBtn).toBeVisible();
	await jokerBtn.click();
	await expect(page.getByTestId('joker-info')).toContainText('target:m1');
	await expect(replaceBtn).toBeDisabled(); // still need 3 tiles

	// Select 3 tiles in Rack: first is replacement tile (7red), next two for new meld (8red, 9red)
	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	// rack has rack-7red at index 0, rack-8red at 1, rack-9red at 2
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await tileButtons.nth(2).click();
	await expect(replaceBtn).toBeEnabled();

	// Click replace → sends op 9 with targetMeldId, tileId, newMeldTiles
	await replaceBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":9');
	await expect(page.getByTestId('last-sent')).toContainText('m1');
	await expect(page.getByTestId('last-sent')).toContainText('tileId');
	await expect(page.getByTestId('last-sent')).toContainText('newMeldTiles');

	// Set NotOpened → should disable even with target and 3 tiles
	await page.getByRole('button', { name: 'Set HasOpened + Joker Run 5-6-J' }).click();
	await page.getByTestId('joker-m1-joker-1').click();
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await tileButtons.nth(2).click();
	await expect(replaceBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set NotOpened' }).click();
	await expect(page.getByTestId('joker-info')).toContainText('hasOpened:false');
	await expect(replaceBtn).toBeDisabled();
});
