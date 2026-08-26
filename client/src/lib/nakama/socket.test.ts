import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { socketStore, createSocket, getSocket, _resetForTest, setMatchDataHandler } from './socket';

describe('Socket store — Day 20', () => {
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
		_resetForTest();
	});

	it('socketStore is writable null initially', () => {
		expect(get(socketStore)).toBeNull();
		expect(getSocket()).toBeNull();
	});

	it('createSocket resolves to a socket', async () => {
		const sock = await createSocket();
		expect(sock).toBeTruthy();
		expect(typeof (sock as unknown as { connect: unknown }).connect).toBe('function');
	});

	it('createSocket sets socketStore and is singleton', async () => {
		const s1 = await createSocket();
		expect(get(socketStore)).toBe(s1);
		expect(getSocket()).toBe(s1);
		const s2 = await createSocket();
		expect(s2).toBe(s1);
	});

	it('setMatchDataHandler wires onmatchdata forwarding', async () => {
		let received: unknown = null;
		setMatchDataHandler((data) => {
			received = data;
		});
		const sock = await createSocket();
		// simulate onmatchdata
		const handler = (sock as unknown as { onmatchdata: (d: unknown) => void }).onmatchdata;
		expect(typeof handler).toBe('function');
		handler({ op_code: 101, data: '{}' });
		expect(received).toEqual({ op_code: 101, data: '{}' });
	});

	it('getSocket returns same as store', async () => {
		expect(getSocket()).toBeNull();
		await createSocket();
		expect(getSocket()).not.toBeNull();
	});
});
