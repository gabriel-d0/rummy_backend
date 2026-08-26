import { writable, derived, get } from 'svelte/store';
import {
	isValidPrivateSnapshot,
	isValidPublicSnapshot,
	type PrivateSnapshot,
	type PublicSnapshot
} from './snapshot';
import { OpServerState, OpServerStatePublic } from '../nakama/protocol';
import { setMatchDataHandler } from '../nakama/socket';

// Day 22 — Game store — private — Svelte writable PrivateSnapshot|null, onPrivateSnapshot, lastPrivate, derived isMyTurn

export const privateStore = writable<PrivateSnapshot | null>(null);
export const publicStore = writable<PublicSnapshot | null>(null);

// lastPrivate mirrors privateStore but also kept as plain variable for quick sync access and reconnection persistence
export let lastPrivate: PrivateSnapshot | null = null;

export function onPrivateSnapshot(snap: unknown): boolean {
	if (!isValidPrivateSnapshot(snap)) return false;
	const ps = snap as PrivateSnapshot;
	lastPrivate = ps;
	privateStore.set(ps);
	// also update public part from private (private contains public fields)
	publicStore.set({
		v: ps.v,
		gamePhase: ps.gamePhase,
		turnPhase: ps.turnPhase,
		currentSeat: ps.currentSeat,
		players: ps.players,
		stockCount: ps.stockCount,
		discardRow: ps.discardRow,
		tableMelds: ps.tableMelds,
		winner: ps.winner
	});
	try {
		localStorage.setItem(`rummy_lastPrivate:${ps.ownSeat}`, JSON.stringify(ps));
	} catch (_err) {
		void _err;
	}
	return true;
}

export function onPublicSnapshot(snap: unknown): boolean {
	if (!isValidPublicSnapshot(snap)) return false;
	const pub = snap as PublicSnapshot;
	publicStore.set(pub);
	return true;
}

export function getPrivateSnapshot(): PrivateSnapshot | null {
	return get(privateStore);
}

export function getPublicSnapshot(): PublicSnapshot | null {
	return get(publicStore);
}

// derived: is it this player's turn?
export const isMyTurn = derived(privateStore, ($priv) => {
	if (!$priv) return false;
	return $priv.currentSeat === $priv.ownSeat;
});

export const myRack = derived(privateStore, ($priv) => $priv?.ownRack ?? []);

export const mySeat = derived(privateStore, ($priv) => $priv?.ownSeat ?? -1);

// Handle raw Nakama MatchData (op_code + data Uint8Array|string)
export function handleMatchData(opCode: number, rawData: unknown): boolean {
	try {
		let jsonStr: string | null = null;
		if (typeof rawData === 'string') jsonStr = rawData;
		else if (rawData instanceof Uint8Array) jsonStr = new TextDecoder().decode(rawData);
		else if (Array.isArray(rawData) && rawData.length > 0 && rawData[0] instanceof Uint8Array) {
			jsonStr = new TextDecoder().decode(rawData[0] as Uint8Array);
		} else if (rawData && typeof rawData === 'object') {
			// already parsed object (tests)
			jsonStr = JSON.stringify(rawData);
		}
		if (!jsonStr) return false;
		const parsed = JSON.parse(jsonStr) as unknown;
		if (opCode === OpServerState) return onPrivateSnapshot(parsed);
		if (opCode === OpServerStatePublic) return onPublicSnapshot(parsed);
		return false;
	} catch (_err) {
		void _err;
		return false;
	}
}

// Wire socket onmatchdata → store (call once at app start in browser)
export function initGameStore(): void {
	try {
		setMatchDataHandler((result: unknown) => {
			const r = result as Record<string, unknown>;
			// nakama-js MatchData shape: { match_id, op_code, data, presence }
			const op = (r.op_code ?? r.opCode ?? r.op) as number | undefined;
			const data = (r.data ?? r.payload) as unknown;
			if (typeof op === 'number') handleMatchData(op, data);
		});
	} catch (_err) {
		void _err;
	}
}

// Auto-wire in browser (not in vitest node)
// eslint-disable-next-line @typescript-eslint/no-explicit-any
if (typeof window !== 'undefined' && typeof (globalThis as any).process === 'undefined') {
	// not test env
	try {
		initGameStore();
	} catch (_err) {
		void _err;
	}
}

// For tests: reset stores
export function _resetForTest(): void {
	lastPrivate = null;
	privateStore.set(null);
	publicStore.set(null);
	try {
		// clear persisted lastPrivate keys
		for (let i = 0; i < 4; i++) {
			localStorage.removeItem(`rummy_lastPrivate:${i}`);
		}
	} catch (_err) {
		void _err;
	}
}
