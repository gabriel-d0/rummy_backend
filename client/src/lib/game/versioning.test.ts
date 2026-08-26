import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import {
	SnapshotVersion,
	isValidPublicSnapshot,
	isValidPrivateSnapshot,
	type PrivateSnapshot,
	type PublicSnapshot
} from './snapshot';
import {
	privateStore,
	publicStore,
	privateBySeat,
	handleMatchData,
	_resetForTest,
	lastPrivate
} from './store';
import { OpServerState, OpServerStatePublic } from '../nakama/protocol';

function makePublic(v: number): PublicSnapshot {
	return {
		v,
		gamePhase: 'Playing',
		turnPhase: 'MustDraw',
		currentSeat: 0,
		players: [{ id: 'alice', seat: 0, hasOpened: false, rackCount: 14 }],
		stockCount: 70,
		discardRow: [],
		tableMelds: [],
		winner: -1
	};
}

function makePrivate(v: number): PrivateSnapshot {
	return {
		...makePublic(v),
		ownRack: [{ ID: 't1', Colour: 1, Rank: 5, IsJoker: false }],
		ownSeat: 0
	};
}

describe('Versioning — Day 27 SnapshotVersion=1 bad_version ignore', () => {
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

	it('SnapshotVersion is 1', () => {
		expect(SnapshotVersion).toBe(1);
	});

	it('isValidPublicSnapshot accepts v:1', () => {
		expect(isValidPublicSnapshot(makePublic(1))).toBe(true);
	});

	it('isValidPublicSnapshot rejects v:2 bad_version', () => {
		expect(isValidPublicSnapshot(makePublic(2))).toBe(false);
		expect(isValidPublicSnapshot({ ...makePublic(1), v: 2 } as unknown as PublicSnapshot)).toBe(
			false
		);
	});

	it('isValidPublicSnapshot rejects v:0 and missing v', () => {
		expect(isValidPublicSnapshot(makePublic(0))).toBe(false);
		expect(
			isValidPublicSnapshot({ ...makePublic(1), v: undefined } as unknown as PublicSnapshot)
		).toBe(false);
		expect(isValidPublicSnapshot({ ...makePublic(1), v: null } as unknown as PublicSnapshot)).toBe(
			false
		);
	});

	it('isValidPrivateSnapshot rejects v:2', () => {
		expect(isValidPrivateSnapshot(makePrivate(2))).toBe(false);
		expect(isValidPrivateSnapshot(makePrivate(1))).toBe(true);
	});

	it('handleMatchData ignores v:2 private OpServerState 100', () => {
		const bad = makePrivate(2);
		const ok = handleMatchData(OpServerState, JSON.stringify(bad));
		expect(ok).toBe(false);
		expect(get(privateStore)).toBeNull();
		expect(get(publicStore)).toBeNull();
		expect(lastPrivate).toBeNull();
		expect(privateBySeat.has(0)).toBe(false);
		expect(localStorage.getItem('rummy_lastPrivate:0')).toBeNull();
	});

	it('handleMatchData ignores v:2 public OpServerStatePublic 101', () => {
		const badPub = makePublic(2);
		const ok = handleMatchData(OpServerStatePublic, JSON.stringify(badPub));
		expect(ok).toBe(false);
		expect(get(publicStore)).toBeNull();
		expect(get(privateStore)).toBeNull();
	});

	it('handleMatchData accepts v:1 after v:2 ignored', () => {
		const bad = makePrivate(2);
		expect(handleMatchData(OpServerState, JSON.stringify(bad))).toBe(false);
		const good = makePrivate(1);
		expect(handleMatchData(OpServerState, JSON.stringify(good))).toBe(true);
		expect(get(privateStore)).not.toBeNull();
		expect(get(privateStore)!.v).toBe(1);
		expect(get(publicStore)!.v).toBe(1);
	});

	it('v:2 ignored does not overwrite existing v:1 snapshot', () => {
		const good = makePrivate(1);
		handleMatchData(OpServerState, JSON.stringify(good));
		expect(get(privateStore)!.ownRack[0].ID).toBe('t1');
		const badOverwrite = {
			...makePrivate(2),
			ownRack: [{ ID: 'bad-2', Colour: 1, Rank: 9, IsJoker: false }]
		};
		expect(handleMatchData(OpServerState, JSON.stringify(badOverwrite))).toBe(false);
		// still good
		expect(get(privateStore)!.ownRack[0].ID).toBe('t1');
		expect(get(privateStore)!.v).toBe(1);
	});

	it('v:2 via Uint8Array also ignored', () => {
		const bad = makePrivate(2);
		const bytes = new TextEncoder().encode(JSON.stringify(bad));
		expect(handleMatchData(OpServerState, bytes)).toBe(false);
		expect(get(privateStore)).toBeNull();
	});
});
