import { writable, derived, get } from 'svelte/store';
import {
	isValidPrivateSnapshot,
	isValidPublicSnapshot,
	type PrivateSnapshot,
	type PublicSnapshot
} from './snapshot';

// Day 22 — Game store — private — Svelte writable PrivateSnapshot|null, onPrivateSnapshot, lastPrivate, derived isMyTurn
// Day 23 — public — Svelte writable PublicSnapshot|null, onPublicSnapshot, TableBoard subscribes

export const privateStore = writable<PrivateSnapshot | null>(null);
export const publicStore = writable<PublicSnapshot | null>(null);

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
	return true;
}

export function getPrivateSnapshot(): PrivateSnapshot | null {
	return get(privateStore);
}

export function onPublicSnapshot(snap: unknown): boolean {
	if (!isValidPublicSnapshot(snap)) return false;
	const pub = snap as PublicSnapshot;
	publicStore.set(pub);
	return true;
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

// For tests: reset stores
export function _resetForTest(): void {
	lastPrivate = null;
	privateStore.set(null);
	publicStore.set(null);
}
