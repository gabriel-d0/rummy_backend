import { test, expect } from '@playwright/test';

test('Day 45 — Error toast OpServerError code/message/details/requestId/op 3s bg #dc2626', async ({
	page
}) => {
	await page.goto('/demo/error');
	await expect(page.getByText('Error Toast Demo — Day 45')).toBeVisible();

	// Initially no toast
	await expect(page.getByTestId('toast')).toHaveCount(0);
	await expect(page.getByTestId('error-info')).toContainText('code:none');

	// Trigger bad_payload via handleMatchData OpServerError 102
	await page.getByRole('button', { name: 'Trigger bad_payload' }).click();
	const toast = page.getByTestId('toast');
	await expect(toast).toBeVisible();
	await expect(toast).toHaveAttribute('data-error-code', 'bad_payload');
	await expect(toast).toContainText('invalid tileId format');
	await expect(toast).toContainText('bad_payload');
	// bg #dc2626 check via computed style
	const bg = await toast.evaluate((el) => getComputedStyle(el).backgroundColor);
	// #dc2626 => rgb(220, 38, 38)
	expect(bg).toBe('rgb(220, 38, 38)');
	await expect(page.getByTestId('error-info')).toContainText('code:bad_payload');
	await expect(page.getByTestId('error-info')).toContainText('requestId:req-123');
	await expect(page.getByTestId('error-info')).toContainText('op:2');
	await expect(page.getByTestId('error-details')).toContainText('tileId');

	// Trigger LEAKED should replace toast and show LEAKED code
	await page.getByRole('button', { name: 'Trigger LEAKED' }).click();
	await expect(toast).toBeVisible();
	await expect(toast).toHaveAttribute('data-error-code', 'LEAKED');
	await expect(toast).toContainText('LEAKED');
	await expect(toast).toContainText('private data leak');

	// Trigger bad_request and check details
	await page.getByRole('button', { name: 'Trigger bad_request' }).click();
	await expect(toast).toHaveAttribute('data-error-code', 'bad_request');
	await expect(toast).toContainText('not in rack');

	// Clear button hides toast
	await page.getByRole('button', { name: 'Clear' }).click();
	await expect(page.getByTestId('toast')).toHaveCount(0);

	// Trigger with details shows details line
	await page.getByRole('button', { name: 'Trigger with details' }).click();
	await expect(toast).toBeVisible();
	await expect(toast).toHaveAttribute('data-error-code', 'bad_payload');
	await expect(toast).toContainText('meld invalid');
	await expect(toast).toContainText('duplicate tile');
	await expect(page.getByTestId('error-info')).toContainText('req-999');
});

test('Day 45 — toast auto-dismiss after 3s', async ({ page }) => {
	await page.goto('/demo/error');
	await page.getByRole('button', { name: 'Trigger bad_payload' }).click();
	await expect(page.getByTestId('toast')).toBeVisible();
	// wait 3.5s
	await page.waitForTimeout(3500);
	await expect(page.getByTestId('toast')).toHaveCount(0);
});
