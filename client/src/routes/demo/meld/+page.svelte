<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeMeldState(hasOpened: boolean): PrivateSnapshot {
		// Rack must contain high-value tiles to reach 50: 8-9-10-11-12-13 run = 50 with 6 tiles
		const baseTiles = [
			{ ID: 't8', Colour: 1, Rank: 8, IsJoker: false },
			{ ID: 't9', Colour: 1, Rank: 9, IsJoker: false },
			{ ID: 't10', Colour: 1, Rank: 10, IsJoker: false },
			{ ID: 't11', Colour: 1, Rank: 11, IsJoker: false },
			{ ID: 't12', Colour: 1, Rank: 12, IsJoker: false },
			{ ID: 't13', Colour: 1, Rank: 13, IsJoker: false },
			{ ID: 't1a', Colour: 2, Rank: 5, IsJoker: false },
			{ ID: 't2a', Colour: 3, Rank: 5, IsJoker: false }
		];
		const rack = Array.from({ length: 14 }, (_, i) => {
			const base = baseTiles[i % baseTiles.length];
			return {
				ID: `rack-${i}-${base.ID}`,
				Colour: base.Colour,
				Rank: base.Rank,
				IsJoker: base.IsJoker
			};
		});
		// Ensure at least selectedCount tiles are the high run tiles
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MeldOrDiscard',
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: 0
		};
	}

	function setNotOpened() {
		onPrivateSnapshot(makeMeldState(false));
	}

	function setOpened() {
		onPrivateSnapshot(makeMeldState(true));
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
	}

	let hasOpened = $derived.by(() => {
		const p = $privateStore;
		if (!p) return false;
		return p.players.find((pl) => pl.seat === p.ownSeat)?.hasOpened ?? false;
	});
	let canMeldInfo = $derived.by(() => {
		const p = $privateStore;
		if (!p) return 'no priv';
		return `hasOpened:${hasOpened} turn:${p.turnPhase} game:${p.gamePhase}`;
	});
	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Meld Initial Demo — Day 35</h1>
	<p class="text-sm text-white/60">
		NotOpened selected≥3 → OpClientMeldInitial 6 &#123;melds:[&#123;kind:run, tileIds&#125;]&#125;
		total≥50 ≥1 run
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setNotOpened} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set NotOpened MeldOrDiscard</button
		>
		<button onclick={setOpened} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set HasOpened</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="meld-info">{canMeldInfo}</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-4">
		<Rack />
	</div>
</div>
