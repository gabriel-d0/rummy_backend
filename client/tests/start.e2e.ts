import { test, expect } from '@playwright/test';

test('Day 29 — Start visible only if Waiting ownSeat==0 players>=2 and sends OpClientStart 1', async ({
	page
}) => {
	await page.goto('/demo/start');
	await expect(page.getByText('Start Demo — Day 29')).toBeVisible();
	// Initially no Start (no snapshot)
	await expect(page.getByTestId('start-btn')).toHaveCount(0);

	// Set Waiting Host 2p — should show START
	await page.getByRole('button', { name: 'Set Waiting Host 2p' }).click();
	await expect(page.getByTestId('start-btn')).toBeVisible();
	await expect(page.getByTestId('start-btn')).toContainText('START');
	await expect(page.getByTestId('start-info')).toBeVisible();

	// Click START — should send OpClientStart 1 via lastSent
	await page.getByTestId('start-btn').click();
	await expect(page.getByTestId('last-sent')).toContainText('"op":1');
	await expect(page.getByTestId('last-sent')).toContainText('mock-match');

	// Set Waiting Guest 2p — host is seat 0, guest seat 1 should NOT see START
	await page.getByRole('button', { name: 'Set Waiting Guest 2p' }).click();
	await expect(page.getByTestId('start-btn')).toHaveCount(0);

	// Set Waiting 1p host — only 1 player, no START
	await page.getByRole('button', { name: 'Set Waiting Host 2p' }).click();
	await expect(page.getByTestId('start-btn')).toBeVisible();
	await page.getByRole('button', { name: 'Set Waiting 1p' }).click();
	await expect(page.getByTestId('start-btn')).toHaveCount(0);

	// Set Playing — no START even host 2p
	await page.getByRole('button', { name: 'Set Waiting Host 2p' }).click();
	await expect(page.getByTestId('start-btn')).toBeVisible();
	await page.getByRole('button', { name: 'Set Playing' }).click();
	await expect(page.getByTestId('start-btn')).toHaveCount(0);
});
