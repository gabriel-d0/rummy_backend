import { test, expect } from '@playwright/test';

test('Day 32 — DrawPrevious HasOpened discardRow not empty !IsOpeningDiscard sends OpClientDrawPreviousDiscard 4', async ({
	page
}) => {
	await page.goto('/demo/draw-previous');
	await expect(page.getByText('Draw Previous Demo — Day 32')).toBeVisible();

	const drawPrevBtn = page.getByTestId('draw-prev-btn');
	// Initially no store → disabled
	await expect(drawPrevBtn).toBeDisabled();

	// Set HasOpened + Discard (valid) → enabled
	await page.getByRole('button', { name: 'Set HasOpened + Discard' }).click();
	await expect(page.getByTestId('prev-info')).toContainText('hasOpened:true');
	await expect(page.getByTestId('prev-info')).toContainText('canDrawPrev:true');
	await expect(drawPrevBtn).toBeEnabled();

	// Click → sends op 4
	await drawPrevBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":4');
	await expect(page.getByTestId('last-sent')).toContainText('mock-match');

	// Set Not Opened → should disable (HasOpened false)
	await page.getByRole('button', { name: 'Set Not Opened' }).click();
	await expect(page.getByTestId('prev-info')).toContainText('hasOpened:false');
	await expect(drawPrevBtn).toBeDisabled();

	// Set HasOpened but Opening Discard → disable (IsOpeningDiscard true)
	await page.getByRole('button', { name: 'Set HasOpened + Discard' }).click();
	await expect(drawPrevBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set Opening Discard' }).click();
	await expect(page.getByTestId('prev-info')).toContainText('isOpeningDiscard:true');
	await expect(drawPrevBtn).toBeDisabled();

	// Set Empty Discard → disable (discardLen 0)
	await page.getByRole('button', { name: 'Set HasOpened + Discard' }).click();
	await expect(drawPrevBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set Empty Discard' }).click();
	await expect(page.getByTestId('prev-info')).toContainText('discardLen:0');
	await expect(drawPrevBtn).toBeDisabled();
});
