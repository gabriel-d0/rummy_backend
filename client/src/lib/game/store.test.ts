import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { privateStore, onPrivateSnapshot, getPrivateSnapshot, lastPrivate, _resetForTest } from './store';
import type { PrivateSnapshot } from './snapshot';

describe('Game store — private — Day 22', () => {
	beforeEach(() => {
		_resetForTest();
	});

	it('privateStore is null initially', () => {
		expect(get(privateStore)).toBeNull();
		expect(getPrivateSnapshot()).toBeNull();
		expect(lastPrivate).toBeNull();
	});

	it('onPrivateSnapshot sets privateStore and lastPrivate', () => {
		const snap: PrivateSnapshot = {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: 0,
			players: [{ id: 'alice', seat: 0, hasOpened: false, rackCount: 3 }],
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: [
				{ ID: 't1', Colour: 1, Rank: 5, IsJoker: false },
				{ ID: 't2', Colour: 2, Rank: 5, IsJoker: false },
				{ ID: 't3', Colour: 3, Rank: 5, IsJoker: false }
			],
			ownSeat: 0
		};
		expect(onPrivateSnapshot(snap)).toBe(true);
		expect(get(privateStore)).toEqual(snap);
		expect(getPrivateSnapshot()).toEqual(snap);
		expect(lastPrivate).toEqual(snap);
	});

	it('onPrivateSnapshot rejects invalid', () => {
		expect(onPrivateSnapshot({ v: 2, gamePhase: 'Playing' } as unknown as PrivateSnapshot)).toBe(false);
		expect(get(privateStore)).toBeNull();
	});

	it('privateStore ownRack only local', () => {
		const snap: PrivateSnapshot = {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: 0,
			players: [{ id: 'alice', seat: 0, hasOpened: false, rackCount: 3 }],
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: [
				{ ID: 'a1', Colour: 1, Rank: 1, IsJoker: false },
				{ ID: 'a2', Colour: 1, Rank: 2, IsJoker: false },
				{ ID: 'a3', Colour: 1, Rank: 3, IsJoker: false }
			],
			ownSeat: 0
		};
		onPrivateSnapshot(snap);
		const priv = get(privateStore)!;
		expect(priv.ownRack.map((t) => t.ID)).toEqual(['a1', 'a2', 'a3']);
		expect(priv.ownRack.length).toBe(3);
	});
});
