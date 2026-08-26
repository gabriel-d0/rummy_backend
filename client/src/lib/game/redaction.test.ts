import { describe, it, expect } from 'vitest';
import { checkNoLeak, type PrivateSnapshot, type PublicSnapshot } from './snapshot';

function makePublic(overrides: Partial<PublicSnapshot> = {}): PublicSnapshot {
	return {
		v: 1,
		gamePhase: 'Playing',
		turnPhase: 'MeldOrDiscard',
		currentSeat: 0,
		players: [
			{ id: 'alice', seat: 0, hasOpened: true, rackCount: 14 },
			{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
		],
		stockCount: 70,
		discardRow: [],
		tableMelds: [],
		winner: -1,
		...overrides
	};
}

function makePrivateWithSecrets(seat: number, secretIds: string[]): PrivateSnapshot {
	const pub = makePublic();
	return {
		...pub,
		ownRack: secretIds.map((id) => ({ ID: id, Colour: 1, Rank: 5, IsJoker: false })),
		ownSeat: seat
	};
}

describe('Redaction — Day 24 checkNoLeak', () => {
	it('publicJson does not contain alice-secret when only private has it', () => {
		const priv = makePrivateWithSecrets(0, ['alice-secret-123', 'alice-secret-456']);
		const pubJson = JSON.stringify({
			v: priv.v,
			gamePhase: priv.gamePhase,
			turnPhase: priv.turnPhase,
			currentSeat: priv.currentSeat,
			players: priv.players,
			stockCount: priv.stockCount,
			discardRow: priv.discardRow,
			tableMelds: priv.tableMelds,
			winner: priv.winner
		} as PublicSnapshot);
		expect(
			checkNoLeak(
				pubJson,
				priv.ownRack.map((t) => t.ID)
			)
		).toBe(true);
		expect(pubJson).not.toContain('alice-secret-123');
		expect(pubJson).not.toContain('alice-secret-456');
	});

	it('checkNoLeak detects leak when publicJson contains privateId', () => {
		const leakJson = JSON.stringify({ v: 1, ownRack: [{ ID: 'alice-secret' }] });
		expect(checkNoLeak(leakJson, ['alice-secret'])).toBe(false);
	});

	it('public snapshot with TableMelds does not leak ownRack IDs', () => {
		const priv = makePrivateWithSecrets(0, ['secret-1', 'secret-2']);
		// public snapshot derived from private (mimics onPrivateSnapshot)
		const pub: PublicSnapshot = {
			v: priv.v,
			gamePhase: priv.gamePhase,
			turnPhase: priv.turnPhase,
			currentSeat: priv.currentSeat,
			players: priv.players,
			stockCount: priv.stockCount,
			discardRow: priv.discardRow,
			tableMelds: [
				{
					ID: 'm1',
					Kind: 'run',
					Tiles: [{ ID: 't-public-1', Colour: 1, Rank: 5, IsJoker: false }],
					JokerReps: {},
					OwnerSeat: 0
				}
			],
			winner: priv.winner
		};
		const pubJson = JSON.stringify(pub);
		expect(checkNoLeak(pubJson, ['secret-1', 'secret-2'])).toBe(true);
		expect(pubJson).toContain('t-public-1');
		expect(pubJson).not.toContain('secret-1');
	});

	it('exhaustive: secrets across 3 seats never leak to public JSON', () => {
		const seats = [0, 1, 2];
		for (const seat of seats) {
			const secrets = [`seat-${seat}-secret-a`, `seat-${seat}-secret-b`];
			const priv = makePrivateWithSecrets(seat, secrets);
			const pubJson = JSON.stringify({
				v: priv.v,
				gamePhase: priv.gamePhase,
				turnPhase: priv.turnPhase,
				currentSeat: priv.currentSeat,
				players: priv.players.map((p) => ({ ...p, rackCount: 14 })),
				stockCount: priv.stockCount,
				discardRow: priv.discardRow,
				tableMelds: priv.tableMelds,
				winner: priv.winner
			});
			expect(checkNoLeak(pubJson, secrets)).toBe(true);
			for (const id of secrets) expect(pubJson).not.toContain(id);
		}
	});

	it('public snapshot only exposes rackCount, not Tile IDs', () => {
		const pub = makePublic({
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: 6 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 68,
			tableMelds: [
				{
					ID: 'm1',
					Kind: 'set',
					Tiles: [
						{ ID: 'pub-tile-1', Colour: 2, Rank: 7, IsJoker: false },
						{ ID: 'pub-tile-2', Colour: 1, Rank: 7, IsJoker: false }
					],
					JokerReps: {},
					OwnerSeat: 0
				}
			]
		});
		const pubJson = JSON.stringify(pub);
		expect(pubJson).toContain('rackCount');
		expect(pubJson).toContain('stockCount');
		expect(pubJson).toContain('pub-tile-1');
		expect(pubJson).not.toContain('ownRack');
		expect(pubJson).not.toContain('alice-secret');
	});

	it('checkNoLeak handles empty privateIds', () => {
		const pubJson = JSON.stringify(makePublic());
		expect(checkNoLeak(pubJson, [])).toBe(true);
	});
});
