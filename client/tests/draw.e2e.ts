import { test, expect } from '@playwright/test';

test('Day 31 — Draw stock visible only if Playing MustDraw ownSeat==currentSeat and sends OpClientDrawStock 3', async ({
	page
}) => {
	await page.goto('/demo/draw');
	await expect(page.getByText('Draw Stock Demo — Day 31')).toBeVisible();

	const drawBtn = page.getByTestId('draw-btn');
	// Initially no store → disabled
	await expect(drawBtn).toBeDisabled();
	await expect(page.getByTestId('draw-info')).toContainText('canDraw:false');

	// Set MustDraw My Turn — should enable
	await page.getByRole('button', { name: 'Set MustDraw My Turn' }).click();
	await expect(page.getByTestId('draw-info')).toContainText('isMustDraw:true');
	await expect(page.getByTestId('draw-info')).toContainText('isMyTurn:true');
	await expect(page.getByTestId('draw-info')).toContainText('canDraw:true');
	await expect(drawBtn).toBeEnabled();

	// Click draw — should send OpClientDrawStock 3
	await drawBtn.click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":3');
	await expect(page.getByTestId('last-sent')).toContainText('mock-match');

	// Set MustDraw Not My Turn — should disable
	await page.getByRole('button', { name: 'Set MustDraw Not My Turn' }).click();
	await expect(page.getByTestId('draw-info')).toContainText('isMyTurn:false');
	await expect(page.getByTestId('draw-info')).toContainText('canDraw:false');
	await expect(drawBtn).toBeDisabled();

	// Set MeldOrDiscard — should disable until MustDraw (even my turn)
	await page.getByRole('button', { name: 'Set MustDraw My Turn' }).click();
	await expect(drawBtn).toBeEnabled();
	await page.getByRole('button', { name: 'Set MeldOrDiscard' }).click();
	await expect(page.getByTestId('draw-info')).toContainText('isMustDraw:false');
	await expect(drawBtn).toBeDisabled();
});
