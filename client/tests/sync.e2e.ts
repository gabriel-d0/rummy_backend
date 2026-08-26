import { test, expect } from '@playwright/test';

test('Day 28 — Visual sync PrivateSnapshot Rack only local, PublicSnapshot Table for all', async ({
	page
}) => {
	await page.goto('/demo/sync');
	await expect(page.getByText('Sync Demo — Day 28')).toBeVisible();
	await expect(page.getByTestId('current-rack')).toBeVisible();

	// Initially no rack (clear)
	await expect(page.getByTestId('current-rack')).toContainText('rack:');
	// Set Alice (3 tiles)
	await page.getByRole('button', { name: 'Set Alice' }).click();
	await expect(page.getByTestId('current-rack')).toContainText('alice-secret-1');
	await expect(page.getByTestId('current-rack')).toContainText('alice-secret-2');
	await expect(page.getByTestId('current-rack')).toContainText('table:2');
	// Rack should show 3 tiles (alice's) and not bob's
	await expect(page.getByTestId('rack-section')).toContainText('3 cărți');
	await expect(page.getByTestId('rack-section')).not.toContainText('bob-secret-1');
	// Table should show 2 melds (shared)
	await expect(page.getByTestId('table-section')).toContainText('2 SETURI');

	// Set Bob (2 tiles) — Rack should now show bob's tiles, not alice's, Table still 2
	await page.getByRole('button', { name: 'Set Bob' }).click();
	await expect(page.getByTestId('current-rack')).toContainText('bob-secret-1');
	await expect(page.getByTestId('current-rack')).toContainText('bob-secret-2');
	await expect(page.getByTestId('current-rack')).not.toContainText('alice-secret-1');
	await expect(page.getByTestId('rack-section')).toContainText('2 cărți');
	await expect(page.getByTestId('table-section')).toContainText('2 SETURI');
	// Ensure public table melds are shared (same IDs for both)
	await expect(page.getByTestId('table-section')).toContainText('ETALĂRI PE MASĂ');

	// Clear and verify no leak
	await page.getByRole('button', { name: 'Clear' }).click();
	await expect(page.getByTestId('current-rack')).toContainText('rack:');
	await expect(page.getByTestId('current-rack')).not.toContainText('alice-secret-1');
	await expect(page.getByTestId('current-rack')).not.toContainText('bob-secret-1');
});
