<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeMeldOrDiscard(seat: number, currentSeat = seat): PrivateSnapshot {
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MeldOrDiscard',
			currentSeat,
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: true, rackCount: 14 }
			],
			stockCount: 70,
			discardRow: [
				{
					Tile: { ID: 'disc-0', Colour: 1, Rank: 5, IsJoker: false },
					IsOpeningDiscard: false,
					Index: 0
				}
			],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: seat
		};
	}

	function makeMustDraw(seat: number): PrivateSnapshot {
		return {
			...makeMeldOrDiscard(seat, seat),
			turnPhase: 'MustDraw'
		};
	}

	function makeAfterDiscard(seat: number): PrivateSnapshot {
		// after discard, currentSeat advances (0→1) %2, turn MustDraw for next
		const rack = Array.from({ length: 13 }, (_, i) => ({
			ID: `tile-${i + 2}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: (seat + 1) % 2,
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: 13 },
				{ id: 'bob', seat: 1, hasOpened: true, rackCount: 14 }
			],
			stockCount: 70,
			discardRow: [
				{
					Tile: { ID: 'disc-0', Colour: 1, Rank: 5, IsJoker: false },
					IsOpeningDiscard: false,
					Index: 0
				},
				{
					Tile: { ID: 'tile-1', Colour: 1, Rank: 5, IsJoker: false },
					IsOpeningDiscard: false,
					Index: 1
				}
			],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: seat
		};
	}

	function setMeldOrDiscardMyTurn() {
		onPrivateSnapshot(makeMeldOrDiscard(0, 0));
	}

	function setMeldOrDiscardNotMyTurn() {
		onPrivateSnapshot(makeMeldOrDiscard(0, 1));
	}

	function setMustDraw() {
		onPrivateSnapshot(makeMustDraw(0));
	}

	function simulateDiscard() {
		const priv = $privateStore;
		if (priv) onPrivateSnapshot(makeAfterDiscard(priv.ownSeat));
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
	}

	let isMeldOrDiscard = $derived($privateStore?.turnPhase === 'MeldOrDiscard');
	let currentSeat = $derived($privateStore?.currentSeat ?? -1);
	let ownSeat = $derived($privateStore?.ownSeat ?? -1);
	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Normal Discard Demo — Day 34</h1>
	<p class="text-sm text-white/60">
		MeldOrDiscard ownSeat==currentSeat → OpClientDiscard 2 &#123;tileId&#125; →
		CurrentSeat→(current+1)%n
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button
			onclick={setMeldOrDiscardMyTurn}
			class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set MeldOrDiscard My Turn 0</button
		>
		<button
			onclick={setMeldOrDiscardNotMyTurn}
			class="rounded bg-amber-600 px-3 py-1 text-sm text-white">Set Not My Turn 0≠1</button
		>
		<button onclick={setMustDraw} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Set MustDraw</button
		>
		<button onclick={simulateDiscard} class="rounded bg-zinc-600 px-3 py-1 text-sm text-white"
			>Simulate Discard 0→1</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="discard-info">
		isMeldOrDiscard:{String(isMeldOrDiscard)} currentSeat:{currentSeat} ownSeat:{ownSeat}
	</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-4">
		<Rack />
	</div>
</div>
