<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeMustDraw(seat: number, currentSeat = seat): PrivateSnapshot {
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
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
			ownSeat: seat
		};
	}

	function makeMeldOrDiscard(seat: number): PrivateSnapshot {
		const rack = Array.from({ length: 15 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MeldOrDiscard',
			currentSeat: seat,
			players: [
				{ id: 'alice', seat: 0, hasOpened: false, rackCount: 15 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 69,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: seat
		};
	}

	function setMustDrawMyTurn() {
		onPrivateSnapshot(makeMustDraw(0, 0));
	}

	function setMustDrawNotMyTurn() {
		onPrivateSnapshot(makeMustDraw(0, 1));
	}

	function setMeldOrDiscard() {
		onPrivateSnapshot(makeMeldOrDiscard(0));
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
	}

	let isMustDraw = $derived($privateStore?.turnPhase === 'MustDraw');
	let isMyTurn = $derived(
		$privateStore ? $privateStore.currentSeat === $privateStore.ownSeat : false
	);
	let canDraw = $derived($privateStore?.gamePhase === 'Playing' && isMustDraw && isMyTurn);
	let rackCount = $derived($privateStore?.ownRack.length ?? 0);
	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Draw Stock Demo — Day 31</h1>
	<p class="text-sm text-white/60">
		Playing MustDraw ownSeat==currentSeat → Draw visible OpClientDrawStock 3, disable until
		MeldOrDiscard
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setMustDrawMyTurn} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set MustDraw My Turn</button
		>
		<button onclick={setMustDrawNotMyTurn} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set MustDraw Not My Turn</button
		>
		<button onclick={setMeldOrDiscard} class="rounded bg-zinc-600 px-3 py-1 text-sm text-white"
			>Set MeldOrDiscard</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="draw-info">
		isMustDraw:{String(isMustDraw)} isMyTurn:{String(isMyTurn)} canDraw:{String(canDraw)} rack:{rackCount}
	</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-4">
		<Rack />
	</div>
</div>
