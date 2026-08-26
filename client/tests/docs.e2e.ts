import { test, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';

test('Day 5 — README contains VITE_NAKAMA_HOST', async () => {
	const readmePath = path.resolve('README.md');
	const content = fs.readFileSync(readmePath, 'utf-8');
	expect(content).toContain('VITE_NAKAMA_HOST');
	expect(content).toContain('VITE_NAKAMA_PORT');
	expect(content).toContain('VITE_NAKAMA_KEY');
	expect(content).toContain('defaultkey');
});

test('Day 5 — .env.example has VITE_NAKAMA_*', async () => {
	const envPath = path.resolve('.env.example');
	const content = fs.readFileSync(envPath, 'utf-8');
	expect(content).toContain('VITE_NAKAMA_HOST=127.0.0.1');
	expect(content).toContain('VITE_NAKAMA_PORT=7350');
	expect(content).toContain('VITE_NAKAMA_KEY=defaultkey');
});
