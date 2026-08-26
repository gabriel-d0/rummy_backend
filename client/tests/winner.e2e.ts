import { test, expect } from '@playwright/test';

test('Day 39 — Winner RoundComplete Winner 0 overlay RESTART MASA', async ({ page }) => {
	await page.goto('/demo/winner');
	await expect(page.getByText('Winner Demo — Day 39')).toBeVisible();

	// Initially no winner
	await expect(page.getByTestId('winner-info')).toContainText('isWinner:false');
	await expect(page.getByTestId('winner-overlay')).toHaveCount(0);

	// Set Winner 0 — should show overlay with Winner 0 and RESTART MASA
	await page.getByRole('button', { name: 'Set Winner 0' }).click();
	await expect(page.getByTestId('winner-info')).toContainText('isWinner:true');
	await expect(page.getByTestId('winner-info')).toContainText('winner:0');
	await expect(page.getByTestId('winner-overlay')).toBeVisible();
	await expect(page.getByTestId('winner-text')).toContainText('Winner 0');
	await expect(page.getByTestId('restart-btn')).toContainText('RESTART MASA');

	// Set Winner 1 — overlay should update to Winner 1
	await page.getByRole('button', { name: 'Set Winner 1' }).click();
	await expect(page.getByTestId('winner-text')).toContainText('Winner 1');

	// Set Playing — overlay should hide
	await page.getByRole('button', { name: 'Set Playing' }).click();
	await expect(page.getByTestId('winner-info')).toContainText('isWinner:false');
	await expect(page.getByTestId('winner-overlay')).toHaveCount(0);

	// Back to Winner 0
	await page.getByRole('button', { name: 'Set Winner 0' }).click();
	await expect(page.getByTestId('winner-overlay')).toBeVisible();
});
