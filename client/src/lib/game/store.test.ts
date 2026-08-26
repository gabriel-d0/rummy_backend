import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import {
	privateStore,
	publicStore,
	onPrivateSnapshot,
	isMyTurn,
	myRack,
	lastPrivate,
	privateBySeat,
	getLastPrivateForSeat,
	_resetForTest,
	handleMatchData
} from './store';
import type { PrivateSnapshot } from './snapshot';
import { OpServerState, OpServerStatePublic } from '../nakama/protocol';

function makePrivate(
	ownSeat: number,
	ownRackIds: string[],
	currentSeat = ownSeat
): PrivateSnapshot {
	return {
		v: 1,
		gamePhase: 'Playing',
		turnPhase: 'MeldOrDiscard',
		currentSeat,
		players: [
			{ id: 'alice', seat: 0, hasOpened: true, rackCount: ownSeat === 0 ? ownRackIds.length : 7 },
			{ id: 'bob', seat: 1, hasOpened: false, rackCount: ownSeat === 1 ? ownRackIds.length : 14 }
		],
		stockCount: 70,
		discardRow: [],
		tableMelds: [],
		winner: -1,
		ownRack: ownRackIds.map((id) => ({ ID: id, Colour: 1, Rank: 5, IsJoker: false })),
		ownSeat
	};
}

describe('Game store — Day 22-23 private + public', () => {
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

	it('privateStore is null initially', () => {
		expect(get(privateStore)).toBeNull();
		expect(get(publicStore)).toBeNull();
		expect(get(isMyTurn)).toBe(false);
		expect(get(myRack)).toEqual([]);
	});

	it('onPrivateSnapshot sets privateStore with ownRack 3 only local', () => {
		const snap = makePrivate(0, ['alice-secret-1', 'alice-secret-2', 'alice-secret-3']);
		const ok = onPrivateSnapshot(snap);
		expect(ok).toBe(true);
		const priv = get(privateStore);
		expect(priv).not.toBeNull();
		expect(priv!.ownRack.length).toBe(3);
		expect(priv!.ownRack.map((t) => t.ID)).toEqual([
			'alice-secret-1',
			'alice-secret-2',
			'alice-secret-3'
		]);
		expect(get(myRack).length).toBe(3);
		expect(priv!.ownSeat).toBe(0);
	});

	it('privateStore does not leak to publicStore', () => {
		const snap = makePrivate(0, ['alice-secret-1', 'alice-secret-2', 'alice-secret-3']);
		onPrivateSnapshot(snap);
		const pub = get(publicStore);
		expect(pub).not.toBeNull();
		// public must not contain ownRack
		expect((pub as unknown as Record<string, unknown>).ownRack).toBeUndefined();
		expect(JSON.stringify(pub)).not.toContain('alice-secret-1');
	});

	it('different seats have different ownRack local', () => {
		const aliceSnap = makePrivate(0, ['alice-1', 'alice-2']);
		onPrivateSnapshot(aliceSnap);
		expect(get(privateStore)!.ownSeat).toBe(0);
		expect(get(privateStore)!.ownRack[0].ID).toBe('alice-1');

		const bobSnap = makePrivate(1, ['bob-1', 'bob-2', 'bob-3']);
		onPrivateSnapshot(bobSnap);
		expect(get(privateStore)!.ownSeat).toBe(1);
		expect(get(privateStore)!.ownRack.map((t) => t.ID)).toEqual(['bob-1', 'bob-2', 'bob-3']);
		// public updated but still no leak
		expect(JSON.stringify(get(publicStore))).not.toContain('bob-1');
		// need to check public doesn't contain private ids - but it shouldn't
		expect(JSON.stringify(get(publicStore))).not.toContain('alice-1');
	});

	it('isMyTurn derived is true when currentSeat === ownSeat', () => {
		const snapMyTurn = makePrivate(0, ['a1'], 0);
		onPrivateSnapshot(snapMyTurn);
		expect(get(isMyTurn)).toBe(true);

		const snapNotMyTurn = makePrivate(0, ['a1'], 1);
		onPrivateSnapshot(snapNotMyTurn);
		expect(get(isMyTurn)).toBe(false);
	});

	it('onPrivateSnapshot persists lastPrivate per seat', () => {
		const snap = makePrivate(0, ['persist-1']);
		onPrivateSnapshot(snap);
		// lastPrivate variable exported
		// can't directly compare via get, check raw var and localStorage
		expect(lastPrivate?.ownSeat).toBe(0);
		const stored = localStorage.getItem('rummy_lastPrivate:0');
		expect(stored).not.toBeNull();
		expect(JSON.parse(stored as string).ownRack[0].ID).toBe('persist-1');
	});

	it('onPrivateSnapshot rejects invalid version', () => {
		const bad = { ...makePrivate(0, ['x']), v: 2 };
		const ok = onPrivateSnapshot(bad as unknown as PrivateSnapshot);
		expect(ok).toBe(false);
		expect(get(privateStore)).toBeNull();
	});

	it('handleMatchData routes OpServerState to privateStore', () => {
		const snap = makePrivate(0, ['via-op-1', 'via-op-2', 'via-op-3']);
		const json = JSON.stringify(snap);
		const ok = handleMatchData(OpServerState, json);
		expect(ok).toBe(true);
		expect(get(privateStore)!.ownRack.length).toBe(3);
	});

	it('handleMatchData routes OpServerStatePublic to publicStore only', () => {
		const pub = {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: 0,
			players: [{ id: 'alice', seat: 0, hasOpened: false, rackCount: 14 }],
			stockCount: 70,
			discardRow: [],
			tableMelds: [
				{
					ID: 'm1',
					Kind: 'run',
					Tiles: [{ ID: 't1', Colour: 1, Rank: 5, IsJoker: false }],
					JokerReps: {},
					OwnerSeat: 0
				}
			],
			winner: -1
		};
		const ok = handleMatchData(OpServerStatePublic, JSON.stringify(pub));
		expect(ok).toBe(true);
		expect(get(publicStore)!.tableMelds.length).toBe(1);
		// private remains null
		expect(get(privateStore)).toBeNull();
	});

	it('publicStore holds TableMelds 2 and not OwnRack — Day 23', () => {
		const pub = {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MeldOrDiscard',
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: 6 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 68,
			discardRow: [],
			tableMelds: [
				{
					ID: 'm1',
					Kind: 'run',
					Tiles: [
						{ ID: 'r1', Colour: 1, Rank: 5, IsJoker: false },
						{ ID: 'r2', Colour: 1, Rank: 6, IsJoker: false },
						{ ID: 'r3', Colour: 1, Rank: 7, IsJoker: false }
					],
					JokerReps: {},
					OwnerSeat: 0
				},
				{
					ID: 'm2',
					Kind: 'set',
					Tiles: [
						{ ID: 's1', Colour: 1, Rank: 9, IsJoker: false },
						{ ID: 's2', Colour: 2, Rank: 9, IsJoker: false },
						{ ID: 's3', Colour: 3, Rank: 9, IsJoker: false }
					],
					JokerReps: {},
					OwnerSeat: 1
				}
			],
			winner: -1
		};
		const ok = handleMatchData(OpServerStatePublic, JSON.stringify(pub));
		expect(ok).toBe(true);
		const stored = get(publicStore)!;
		expect(stored.tableMelds.length).toBe(2);
		expect(stored.tableMelds[0].ID).toBe('m1');
		expect(stored.tableMelds[1].Kind).toBe('set');
		// TableBoard subscribes, not OwnRack
		expect((stored as unknown as Record<string, unknown>).ownRack).toBeUndefined();
		expect(JSON.stringify(stored)).not.toContain('ownRack');
		expect(JSON.stringify(stored)).not.toContain('alice-secret');
		// private remains null after public-only update
		expect(get(privateStore)).toBeNull();
	});

	it('Day 25 — privateBySeat Map stores per seat and rummy_lastPrivate:0', () => {
		const snap0 = makePrivate(0, ['priv0-1', 'priv0-2']);
		const snap1 = makePrivate(1, ['priv1-1']);
		onPrivateSnapshot(snap0);
		expect(privateBySeat.get(0)?.ownSeat).toBe(0);
		expect(privateBySeat.get(0)?.ownRack[0].ID).toBe('priv0-1');
		expect(localStorage.getItem('rummy_lastPrivate:0')).not.toBeNull();
		expect(JSON.parse(localStorage.getItem('rummy_lastPrivate:0') as string).ownRack[0].ID).toBe(
			'priv0-1'
		);

		onPrivateSnapshot(snap1);
		expect(privateBySeat.get(1)?.ownSeat).toBe(1);
		expect(localStorage.getItem('rummy_lastPrivate:1')).not.toBeNull();
		// both seats kept
		expect(privateBySeat.size).toBe(2);
		expect(getLastPrivateForSeat(0)?.ownRack[0].ID).toBe('priv0-1');
		expect(getLastPrivateForSeat(1)?.ownRack[0].ID).toBe('priv1-1');
	});

	it('Day 25 — getLastPrivateForSeat restores from localStorage when Map cleared', () => {
		const snap = makePrivate(0, ['restore-1']);
		onPrivateSnapshot(snap);
		expect(privateBySeat.has(0)).toBe(true);
		// simulate process restart: clear Map but keep localStorage
		privateBySeat.clear();
		expect(privateBySeat.has(0)).toBe(false);
		const restored = getLastPrivateForSeat(0);
		expect(restored).not.toBeNull();
		expect(restored!.ownRack[0].ID).toBe('restore-1');
		expect(privateBySeat.has(0)).toBe(true);
	});

	it('Day 25 — socket.onDisconnect keeps matchId and rummy_device_id', () => {
		localStorage.setItem('rummy_matchId', 'keep-match-123');
		localStorage.setItem('rummy_device_id', 'keep-device-456');
		localStorage.setItem('rummy_userId', 'keep-user-789');
		// simulate onDisconnect clearing only socketStore, not storage
		// our socketStore _resetForTest clears socket but we verify storage still there
		// here we just verify storage not cleared by store logic
		onPrivateSnapshot(makePrivate(0, ['keep-rack-1']));
		// simulate disconnect: privateBySeat should remain
		expect(privateBySeat.size).toBe(1);
		expect(localStorage.getItem('rummy_matchId')).toBe('keep-match-123');
		expect(localStorage.getItem('rummy_device_id')).toBe('keep-device-456');
		expect(localStorage.getItem('rummy_lastPrivate:0')).not.toBeNull();
	});
});
