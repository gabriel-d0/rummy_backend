import { writable } from 'svelte/store';
import { getClient } from './client';
import type { Session } from '@heroiclabs/nakama-js';

// Day 19 — Auth store — Svelte writable Session|null, deviceId→localStorage rummy_device_id, token→rummy_token

function getOrCreateDeviceId(): string {
	const key = 'rummy_device_id';
	try {
		let id = localStorage.getItem(key);
		if (id) return id;
		if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
			id = (crypto as unknown as { randomUUID: () => string }).randomUUID();
		} else {
			id = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
				const r = (Math.random() * 16) | 0;
				const v = c === 'x' ? r : (r & 0x3) | 0x8;
				return v.toString(16);
			});
		}
		localStorage.setItem(key, id!);
		return id!;
	} catch (_err) {
		void _err;
		return 'test-device-id';
	}
}

export const authStore = writable<Session | null>(null);

export async function authenticate(username?: string): Promise<Session> {
	const client = getClient();
	const deviceId = getOrCreateDeviceId();
	const uname = username ?? `rummy-${deviceId.slice(0, 8)}`;
	const session = await client.authenticateDevice(deviceId, true, uname);
	try {
		localStorage.setItem('rummy_token', session.token as string);
		localStorage.setItem('rummy_userId', session.user_id as string);
		localStorage.setItem('rummy_device_id', deviceId);
	} catch (_err) {
		void _err;
	}
	authStore.set(session);
	return session;
}

export function getStoredDeviceId(): string | null {
	try {
		return localStorage.getItem('rummy_device_id');
	} catch (_err) {
		void _err;
		return null;
	}
}
