import { writable, derived, get } from 'svelte/store';
import { isValidPrivateSnapshot, type PrivateSnapshot } from './snapshot';

// Day 22 — Game store — private — Svelte writable PrivateSnapshot|null, onPrivateSnapshot, lastPrivate, derived isMyTurn

export const privateStore = writable<PrivateSnapshot | null>(null);

export let lastPrivate: PrivateSnapshot | null = null;

export function onPrivateSnapshot(snap: unknown): boolean {
	if (!isValidPrivateSnapshot(snap)) return false;
	const ps = snap as PrivateSnapshot;
	lastPrivate = ps;
	privateStore.set(ps);
	return true;
}

export function getPrivateSnapshot(): PrivateSnapshot | null {
	return get(privateStore);
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
}
