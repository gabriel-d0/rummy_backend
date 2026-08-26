<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeSnapshot(opts: {
		seat: number;
		hasOpened: boolean;
		discardIsOpening: boolean;
		discardEmpty: boolean;
		currentSeat: number;
	}): PrivateSnapshot {
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		const discardRow = opts.discardEmpty
			? []
			: [
					{
						Tile: { ID: 'disc-1', Colour: 1, Rank: 7, IsJoker: false },
						IsOpeningDiscard: opts.discardIsOpening,
						Index: 0
					}
				];
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: opts.currentSeat,
			players: [
				{
					id: 'alice',
					seat: 0,
					hasOpened: opts.seat === 0 ? opts.hasOpened : false,
					rackCount: 14
				},
				{ id: 'bob', seat: 1, hasOpened: opts.seat === 1 ? opts.hasOpened : false, rackCount: 14 }
			],
			stockCount: 70,
			discardRow,
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: opts.seat
		};
	}

	function setHasOpenedWithDiscard() {
		onPrivateSnapshot(
			makeSnapshot({
				seat: 0,
				hasOpened: true,
				discardIsOpening: false,
				discardEmpty: false,
				currentSeat: 0
			})
		);
	}

	function setNotOpened() {
		onPrivateSnapshot(
			makeSnapshot({
				seat: 0,
				hasOpened: false,
				discardIsOpening: false,
				discardEmpty: false,
				currentSeat: 0
			})
		);
	}

	function setOpeningDiscard() {
		onPrivateSnapshot(
			makeSnapshot({
				seat: 0,
				hasOpened: true,
				discardIsOpening: true,
				discardEmpty: false,
				currentSeat: 0
			})
		);
	}

	function setEmptyDiscard() {
		onPrivateSnapshot(
			makeSnapshot({
				seat: 0,
				hasOpened: true,
				discardIsOpening: false,
				discardEmpty: true,
				currentSeat: 0
			})
		);
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
	}

	let isMyTurn = $derived(
		$privateStore ? $privateStore.currentSeat === $privateStore.ownSeat : false
	);
	let hasOpened = $derived.by(() => {
		const p = $privateStore;
		if (!p) return false;
		return p.players.find((pl) => pl.seat === p.ownSeat)?.hasOpened ?? false;
	});
	let canDrawPrev = $derived.by(() => {
		const p = $privateStore;
		if (!p) return false;
		if (p.gamePhase !== 'Playing' || p.turnPhase !== 'MustDraw' || !isMyTurn) return false;
		if (!hasOpened) return false;
		if (!p.discardRow.length) return false;
		return !p.discardRow[p.discardRow.length - 1].IsOpeningDiscard;
	});
	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Draw Previous Demo — Day 32</h1>
	<p class="text-sm text-white/60">
		HasOpened + discardRow not empty + !IsOpeningDiscard → IA ULTIMA OpClientDrawPreviousDiscard 4
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button
			onclick={setHasOpenedWithDiscard}
			class="rounded bg-emerald-600 px-3 py-1 text-sm text-white">Set HasOpened + Discard</button
		>
		<button onclick={setNotOpened} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set Not Opened</button
		>
		<button onclick={setOpeningDiscard} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Set Opening Discard</button
		>
		<button onclick={setEmptyDiscard} class="rounded bg-zinc-600 px-3 py-1 text-sm text-white"
			>Set Empty Discard</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="prev-info">
		hasOpened:{String(hasOpened)} discardLen:{$privateStore?.discardRow.length ?? 0} isOpeningDiscard:{String(
			$privateStore?.discardRow[0]?.IsOpeningDiscard ?? false
		)} canDrawPrev:{String(canDrawPrev)}
	</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-4">
		<Rack />
	</div>
</div>
