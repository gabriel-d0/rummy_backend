import { test, expect } from '@playwright/test';

test('Day 43 — Turn indicator TopBar Current seat-0 Playing/MustDraw ← current and Draw disabled', async ({
	page
}) => {
	await page.goto('/demo/turn');
	await expect(page.getByText('Turn Indicator Demo — Day 43')).toBeVisible();

	// Initially no turn indicator (no snapshot)
	await expect(page.getByTestId('turn-indicator')).toHaveCount(0);

	// Set MustDraw My Turn S0 — should show Current: seat-0 Playing/MustDraw ← rândul tău
	await page.getByRole('button', { name: 'MustDraw MyTurn S0' }).click();
	await expect(page.getByTestId('turn-indicator')).toBeVisible();
	await expect(page.getByTestId('turn-indicator')).toContainText('Current: seat-0');
	await expect(page.getByTestId('turn-indicator')).toContainText('Playing/MustDraw');
	await expect(page.getByTestId('turn-indicator')).toContainText('← rândul tău');
	await expect(page.getByTestId('turn-info')).toContainText('isMyTurn:true');
	await expect(page.getByTestId('turn-info')).toContainText('canDraw:true');
	await expect(page.getByTestId('draw-btn')).toBeEnabled();

	// Set MustDraw Not My Turn S1 — should show Current: seat-1 ← seat-1 and Draw disabled
	await page.getByRole('button', { name: 'MustDraw NotMyTurn S1' }).click();
	await expect(page.getByTestId('turn-indicator')).toContainText('Current: seat-1');
	await expect(page.getByTestId('turn-indicator')).toContainText('Playing/MustDraw');
	await expect(page.getByTestId('turn-indicator')).toContainText('← seat-1');
	await expect(page.getByTestId('turn-info')).toContainText('isMyTurn:false');
	await expect(page.getByTestId('turn-info')).toContainText('canDraw:false');
	await expect(page.getByTestId('draw-btn')).toBeDisabled();

	// Set MeldOrDiscard MyTurn — Draw disabled even my turn
	await page.getByRole('button', { name: 'MeldOrDiscard MyTurn' }).click();
	await expect(page.getByTestId('turn-indicator')).toContainText('Current: seat-0');
	await expect(page.getByTestId('turn-indicator')).toContainText('Playing/MeldOrDiscard');
	await expect(page.getByTestId('turn-indicator')).toContainText('← rândul tău');
	await expect(page.getByTestId('turn-info')).toContainText('canDraw:false');
	await expect(page.getByTestId('draw-btn')).toBeDisabled();

	// Set OpeningDiscard — Turn indicator still shows Current
	await page.getByRole('button', { name: 'OpeningDiscard MyTurn' }).click();
	await expect(page.getByTestId('turn-indicator')).toContainText('Current: seat-0');
	await expect(page.getByTestId('turn-indicator')).toContainText('OpeningDiscard');

	// Mobile indicator exists (hidden on desktop via sm:hidden, visible on mobile)
	await expect(page.getByTestId('turn-indicator-mobile')).toHaveCount(1);
});
