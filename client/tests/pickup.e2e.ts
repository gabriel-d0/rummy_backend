import { test, expect } from '@playwright/test';

test('Day 33 — Pickup for meld selected 2 + discardIndex via TableBoard OpClientPickupDiscardForMeld 5', async ({
	page
}) => {
	await page.goto('/demo/pickup');
	await expect(page.getByText('Pickup Demo — Day 33')).toBeVisible();

	const pickupBtn = page.getByTestId('pickup-btn');
	// Initially no store → disabled (needs 2 selected + discard)
	await expect(pickupBtn).toBeDisabled();

	// Set valid pickup (HasOpened + 2 discards, last not opening)
	await page.getByRole('button', { name: 'Set Valid Pickup' }).click();
	await expect(page.getByTestId('pickup-info')).toContainText('hasOpened:true');
	await expect(pickupBtn).toBeDisabled(); // still needs selected 2 + discardIndex

	// Click discard tile via TableBoard (index 0)
	const discard0 = page.getByTestId('discard-tile-0');
	await expect(discard0).toBeVisible();
	await discard0.click();
	await expect(page.getByTestId('pickup-info')).toContainText('selectedDiscard:0');
	await expect(pickupBtn).toBeDisabled(); // still needs 2 tiles selected

	// Select 2 tiles in Rack
	const tileButtons = page.locator('div.flex.min-h-\\[110px\\] button');
	await expect(tileButtons.first()).toBeVisible();
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	// Now should be enabled (HasOpened, MustDraw, myTurn, 2 selected, discardIndex 0 not opening)
	await expect(pickupBtn).toBeEnabled();

	// Click pickup → sends op 5 with discardIndex and tileIds
	await pickupBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":5');
	await expect(page.getByTestId('last-sent')).toContainText('discardIndex');
	await expect(page.getByTestId('last-sent')).toContainText('tile-1');

	// Set Not Opened → should disable even with 2 selected + discard
	await page.getByRole('button', { name: 'Set Valid Pickup' }).click();
	await page.getByTestId('discard-tile-0').click();
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await expect(pickupBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set Not Opened' }).click();
	await expect(pickupBtn).toBeDisabled();

	// Set Opening Discard (last discard is opening) → pick that opening index should disable
	await page.getByRole('button', { name: 'Set Opening Discard' }).click();
	// Click second discard which is opening (index 1)
	const discard1 = page.getByTestId('discard-tile-1');
	await expect(discard1).toBeVisible();
	await discard1.click();
	await expect(page.getByTestId('pickup-info')).toContainText('selectedDiscard:1');
	// Select 2 tiles again (need to re-select after snapshot reset)
	await tileButtons.nth(0).click();
	await tileButtons.nth(1).click();
	await expect(pickupBtn).toBeDisabled(); // opening discard selected
});
