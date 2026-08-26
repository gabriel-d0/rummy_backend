<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import TableBoard from '../../../components/TableBoard.svelte';
	import {
		privateStore,
		onPrivateSnapshot,
		_resetForTest,
		pickupDiscardIndex
	} from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeMustDrawWithDiscard(
		hasOpened: boolean,
		discardEmpty: boolean,
		opening: boolean
	): PrivateSnapshot {
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		const discardRow = discardEmpty
			? []
			: [
					{
						Tile: { ID: 'disc-0', Colour: 1, Rank: 5, IsJoker: false },
						IsOpeningDiscard: false,
						Index: 0
					},
					{
						Tile: { ID: 'disc-1', Colour: 2, Rank: 7, IsJoker: false },
						IsOpeningDiscard: opening,
						Index: 1
					}
				];
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 70,
			discardRow,
			tableMelds: [],
			winner: -1,
			ownRack: rack,
			ownSeat: 0
		};
	}

	function setValidPickup() {
		onPrivateSnapshot(makeMustDrawWithDiscard(true, false, false));
		pickupDiscardIndex.set(null);
	}

	function setNotOpened() {
		onPrivateSnapshot(makeMustDrawWithDiscard(false, false, false));
	}

	function setOpeningDiscard() {
		onPrivateSnapshot(makeMustDrawWithDiscard(true, false, true));
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
	}

	let canPickupInfo = $derived.by(() => {
		const p = $privateStore;
		if (!p) return 'no priv';
		const hasOpened = p.players.find((pl) => pl.seat === p.ownSeat)?.hasOpened ?? false;
		return `hasOpened:${hasOpened} discardLen:${p.discardRow.length} isOpening:${p.discardRow[1]?.IsOpeningDiscard ?? false}`;
	});

	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
	let selectedIdx = $derived($pickupDiscardIndex);
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Pickup Demo — Day 33</h1>
	<p class="text-sm text-white/60">
		Rack selected 2 + discardIndex via TableBoard → OpClientPickupDiscardForMeld 5
		&#123;discardIndex, tileIds:[2]&#125;
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setValidPickup} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set Valid Pickup (HasOpened + 2 discards)</button
		>
		<button onclick={setNotOpened} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set Not Opened</button
		>
		<button onclick={setOpeningDiscard} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Set Opening Discard</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="pickup-info">
		{canPickupInfo} selectedDiscard:{String(selectedIdx)}
	</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-4 space-y-3">
		<div data-testid="table-section">
			<TableBoard />
		</div>
		<div data-testid="rack-section">
			<Rack />
		</div>
	</div>
</div>
