import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { reconnect, simulatePrivateAfterReconnect, hasStoredMatch } from './reconnect';
import { createMatch } from './match';
import { privateStore, _resetForTest as resetGame } from '../game/store';
import { _resetForTest as resetSocket } from './socket';
import { _resetForTest as resetMatch } from './match';

describe('Reconnection — rejoin — Day 26', () => {
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
		resetGame();
		resetSocket();
		resetMatch();
	});

	it('hasStoredMatch false initially', () => {
		expect(hasStoredMatch()).toBe(false);
	});

	it('reconnect returns null when no stored match', async () => {
		const id = await reconnect();
		expect(id).toBeNull();
	});

	it('reconnect rejoins stored matchId', async () => {
		const created = await createMatch();
		expect(created).toBe('mock-match');
		expect(hasStoredMatch()).toBe(true);
		expect(localStorage.getItem('rummy_matchId')).toBe('mock-match');

		// simulate disconnect keeping matchId (socket cleared but storage kept)
		// socket _reset clears socketStore but not localStorage rummy_matchId (tested Day 25)
		// Here we keep it, then reconnect
		const rejoined = await reconnect();
		expect(rejoined).toBe('mock-match');
		expect(localStorage.getItem('rummy_matchId')).toBe('mock-match');
	});

	it('reconnect ownRack new-1 not old-1 — Rack rehydrates from OwnRack', async () => {
		await createMatch();
		// initial old snapshot
		const oldSnap = simulatePrivateAfterReconnect(['old-1', 'old-2'], 0, 0);
		expect(get(privateStore)!.ownRack.map((t) => t.ID)).toEqual(['old-1', 'old-2']);
		expect(oldSnap.ownRack[0].ID).toBe('old-1');

		// reconnect
		const rejoined = await reconnect();
		expect(rejoined).toBe('mock-match');

		// server sends new PrivateSnapshot after rejoin (new rack)
		const newSnap = simulatePrivateAfterReconnect(['new-1', 'new-2', 'new-3'], 0, 0);
		const cur = get(privateStore)!;
		expect(cur.ownRack.map((t) => t.ID)).toEqual(['new-1', 'new-2', 'new-3']);
		expect(cur.ownRack.map((t) => t.ID)).not.toContain('old-1');
		expect(newSnap.ownRack[0].ID).toBe('new-1');
		expect(localStorage.getItem('rummy_lastPrivate:0')).toContain('new-1');
		expect(localStorage.getItem('rummy_lastPrivate:0')).not.toContain('old-1');
	});

	it('reconnect keeps rummy_lastPrivate per seat, new overwrites old', async () => {
		await createMatch();
		simulatePrivateAfterReconnect(['old-1'], 0);
		expect(localStorage.getItem('rummy_lastPrivate:0')).toContain('old-1');

		await reconnect();
		simulatePrivateAfterReconnect(['new-1'], 0);
		const stored = JSON.parse(localStorage.getItem('rummy_lastPrivate:0') as string);
		expect(stored.ownRack[0].ID).toBe('new-1');
		expect(get(privateStore)!.ownRack[0].ID).toBe('new-1');
	});

	it('reconnect after disconnect keeps matchId and can rehydrate new rack', async () => {
		const id = await createMatch();
		simulatePrivateAfterReconnect(['old-1'], 0);
		// simulate socket disconnect: clear socket but keep storage
		resetSocket(); // clears socketStore but not rummy_matchId
		expect(localStorage.getItem('rummy_matchId')).toBe(id);
		// reconnect should succeed
		const rejoined = await reconnect();
		expect(rejoined).toBe(id);
		simulatePrivateAfterReconnect(['new-1'], 0);
		expect(get(privateStore)!.ownRack[0].ID).toBe('new-1');
		expect(get(privateStore)!.ownRack[0].ID).not.toBe('old-1');
	});
});
