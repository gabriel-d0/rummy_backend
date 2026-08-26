<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeOpening(seat: number, tileCount: number, currentSeat = seat): PrivateSnapshot {
		const rack = Array.from({ length: tileCount }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: 'OpeningDiscard',
			turnPhase: '',
			currentSeat,
			players: [
				{ id: 'alice', seat: 0, hasOpened: false, rackCount: seat === 0 ? tileCount : 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: seat === 1 ? tileCount : 14 }
			],
			stockCount: 77,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: seat
		};
	}

	function makeAfterDiscard(seat: number): PrivateSnapshot {
		// simulate server after discarding one tile: rack 14, turn advances to next seat, phase Playing MustDraw
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 2}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: 1,
			players: [
				{ id: 'alice', seat: 0, hasOpened: false, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 77,
			discardRow: [
				{
					Tile: { ID: 'tile-1', Colour: 1, Rank: 1, IsJoker: false },
					IsOpeningDiscard: true,
					Index: 0
				}
			],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: seat
		};
	}

	function setOpeningHost15() {
		onPrivateSnapshot(makeOpening(0, 15, 0));
	}

	function setOpeningGuest15() {
		onPrivateSnapshot(makeOpening(1, 15, 1));
	}

	function setOpeningNotYourTurn() {
		onPrivateSnapshot(makeOpening(0, 15, 1));
	}

	function simulateAfterDiscard() {
		// for demo, manually set to after discard to show 15→14
		const priv = $privateStore;
		if (priv) onPrivateSnapshot(makeAfterDiscard(priv.ownSeat));
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
	}

	let isOpening = $derived($privateStore?.gamePhase === 'OpeningDiscard');
	let rackCount = $derived($privateStore?.ownRack.length ?? 0);
	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Opening Discard Demo — Day 30</h1>
	<p class="text-sm text-white/60">
		OpeningDiscard ownSeat==currentSeat ownRack 15 → OpClientDiscard 2 &#123;tileId&#125; → 15→14
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setOpeningHost15} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set Opening Host 15</button
		>
		<button onclick={setOpeningGuest15} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Set Opening Guest 15</button
		>
		<button
			onclick={setOpeningNotYourTurn}
			class="rounded bg-amber-600 px-3 py-1 text-sm text-white">Set Not Your Turn</button
		>
		<button onclick={simulateAfterDiscard} class="rounded bg-zinc-600 px-3 py-1 text-sm text-white"
			>Simulate 15→14</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="opening-info">
		isOpening:{String(isOpening)} rack:{rackCount} currentSeat:{$privateStore?.currentSeat ?? -1} ownSeat:{$privateStore?.ownSeat ??
			-1}
	</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-4">
		<Rack />
	</div>
</div>
