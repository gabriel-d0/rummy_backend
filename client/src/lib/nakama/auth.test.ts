import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { authStore, getStoredDeviceId } from './auth';

describe('Auth store — Day 19', () => {
	beforeEach(() => {
		const store = new Map<string, string>();
		(globalThis as unknown as Record<string, unknown>).localStorage = {
			getItem: (k: string) => (store.has(k) ? store.get(k)! : null),
			setItem: (k: string, v: string) => {
				store.set(k, v);
			},
			removeItem: (k: string) => {
				store.delete(k);
			},
			key: (i: number) => [...store.keys()][i] ?? null,
			get length() {
				return store.size;
			},
			clear: () => store.clear()
		} as unknown as Storage;
		authStore.set(null);
		try {
			localStorage.removeItem('rummy_device_id');
			localStorage.removeItem('rummy_token');
			localStorage.removeItem('rummy_userId');
		} catch (_err) {
			void _err;
		}
	});

	it('authStore is writable null initially', () => {
		expect(get(authStore)).toBeNull();
	});

	it('getStoredDeviceId reads from localStorage', () => {
		localStorage.setItem('rummy_device_id', 'test-id-123');
		expect(getStoredDeviceId()).toBe('test-id-123');
		expect(localStorage.getItem('rummy_device_id')).toBe('test-id-123');
	});

	it('authStore can be set', () => {
		const fakeSession = {
			token: 'tok',
			user_id: 'uid'
		} as unknown as import('@heroiclabs/nakama-js').Session;
		authStore.set(fakeSession);
		expect(get(authStore)).toBe(fakeSession);
	});
});
