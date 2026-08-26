<script lang="ts">
	import TopBar from '../../../components/TopBar.svelte';
	import Rack from '../../../components/Rack.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeSnap(
		gamePhase: string,
		turnPhase: string,
		currentSeat: number,
		ownSeat: number
	): PrivateSnapshot {
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: gamePhase as PrivateSnapshot['gamePhase'],
			turnPhase: turnPhase as PrivateSnapshot['turnPhase'],
			currentSeat,
			players: [
				{ id: 'alice', seat: 0, hasOpened: false, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat
		};
	}

	function setMustDrawMyTurn() {
		onPrivateSnapshot(makeSnap('Playing', 'MustDraw', 0, 0));
	}
	function setMustDrawNotMyTurn() {
		onPrivateSnapshot(makeSnap('Playing', 'MustDraw', 1, 0));
	}
	function setMeldOrDiscardMyTurn() {
		onPrivateSnapshot(makeSnap('Playing', 'MeldOrDiscard', 0, 0));
	}
	function setMeldOrDiscardNotMyTurn() {
		onPrivateSnapshot(makeSnap('Playing', 'MeldOrDiscard', 1, 0));
	}
	function setOpeningMyTurn() {
		onPrivateSnapshot(makeSnap('OpeningDiscard', '', 0, 0));
	}
	function setPlayingSeat1() {
		onPrivateSnapshot(makeSnap('Playing', 'MustDraw', 1, 1));
	}
	function clear() {
		_resetForTest();
	}

	let gamePhase = $derived($privateStore?.gamePhase ?? '');
	let turnPhase = $derived($privateStore?.turnPhase ?? '');
	let currentSeat = $derived($privateStore?.currentSeat ?? -1);
	let ownSeat = $derived($privateStore?.ownSeat ?? -1);
	let isMyTurn = $derived(
		$privateStore ? $privateStore.currentSeat === $privateStore.ownSeat : false
	);
	let canDraw = $derived(
		$privateStore?.gamePhase === 'Playing' && $privateStore?.turnPhase === 'MustDraw' && isMyTurn
	);
</script>

<div class="min-h-screen bg-[#0a2e1a]">
	<TopBar />
	<div class="p-4">
		<h1 class="font-bold text-white">Turn Indicator Demo — Day 43</h1>
		<p class="text-sm text-white/60">
			TopBar Current: seat-0 Playing/MustDraw ← current • Draw disabled if not MustDraw/CurrentSeat
		</p>
		<div class="mt-3 flex flex-wrap gap-2">
			<button
				onclick={setMustDrawMyTurn}
				class="rounded bg-emerald-600 px-3 py-1 text-sm text-white">MustDraw MyTurn S0</button
			>
			<button
				onclick={setMustDrawNotMyTurn}
				class="rounded bg-amber-600 px-3 py-1 text-sm text-white">MustDraw NotMyTurn S1</button
			>
			<button
				onclick={setMeldOrDiscardMyTurn}
				class="rounded bg-sky-600 px-3 py-1 text-sm text-white">MeldOrDiscard MyTurn</button
			>
			<button
				onclick={setMeldOrDiscardNotMyTurn}
				class="rounded bg-zinc-600 px-3 py-1 text-sm text-white">MeldOrDiscard NotMyTurn</button
			>
			<button onclick={setOpeningMyTurn} class="rounded bg-orange-600 px-3 py-1 text-sm text-white"
				>OpeningDiscard MyTurn</button
			>
			<button onclick={setPlayingSeat1} class="rounded bg-violet-600 px-3 py-1 text-sm text-white"
				>Playing Seat1 MyTurn</button
			>
			<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button
			>
		</div>
		<div class="mt-3 text-xs text-white/80" data-testid="turn-info">
			gamePhase:{gamePhase} turnPhase:{turnPhase} currentSeat:{currentSeat} ownSeat:{ownSeat}
			isMyTurn:{String(isMyTurn)} canDraw:{String(canDraw)}
		</div>
		<div class="mt-4">
			<Rack />
		</div>
		<div class="mt-4 text-xs text-white/40">
			Expected TopBar: "Current: seat-0" "Playing/MustDraw" "← rândul tău" when my turn, else "←
			seat-1". Draw button in Rack disabled when not MustDraw or not your turn.
		</div>
	</div>
</div>
