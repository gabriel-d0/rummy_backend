import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { matchStore, createMatch, joinMatch, getMatchId, _resetForTest } from './match';
import { _resetForTest as resetSocket } from './socket';

describe('Match store — Day 21', () => {
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
		resetSocket();
	});

	it('matchStore is null initially', () => {
		expect(get(matchStore)).toBeNull();
		expect(getMatchId()).toBeNull();
		expect(localStorage.getItem('rummy_matchId')).toBeNull();
	});

	it('createMatch creates match and persists rummy_matchId', async () => {
		const id = await createMatch();
		expect(typeof id).toBe('string');
		expect(id.length).toBeGreaterThan(0);
		expect(get(matchStore)).toBe(id);
		expect(getMatchId()).toBe(id);
		expect(localStorage.getItem('rummy_matchId')).toBe(id);
	});

	it('createMatch returns mock-match in test env', async () => {
		const id = await createMatch();
		expect(id).toBe('mock-match');
	});

	it('joinMatch stores given matchId and persists', async () => {
		const joinId = 'test-join-123';
		const id = await joinMatch(joinId);
		expect(id).toBe(joinId);
		expect(get(matchStore)).toBe(joinId);
		expect(localStorage.getItem('rummy_matchId')).toBe(joinId);
	});

	it('alice createMatch and bob joinMatch share same matchId', async () => {
		const aliceId = await createMatch();
		expect(aliceId).toBe('mock-match');
		// simulate bob joining same matchId
		const bobId = await joinMatch(aliceId);
		expect(bobId).toBe(aliceId);
		expect(get(matchStore)).toBe(aliceId);
	});

	it('joinMatch throws if empty', async () => {
		await expect(joinMatch('')).rejects.toThrow('matchId required');
	});

	it('second createMatch returns same singleton mock', async () => {
		const a = await createMatch();
		const b = await createMatch();
		expect(b).toBe(a);
	});
});
