<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import TableBoard from '../../../components/TableBoard.svelte';
	import {
		privateStore,
		onPrivateSnapshot,
		_resetForTest,
		replaceTargetMeldId
	} from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeJokerMeld(hasOpened: boolean): PrivateSnapshot {
		const rack = [
			{ ID: 'rack-7red', Colour: 1, Rank: 7, IsJoker: false },
			{ ID: 'rack-8red', Colour: 1, Rank: 8, IsJoker: false },
			{ ID: 'rack-9red', Colour: 1, Rank: 9, IsJoker: false },
			{ ID: 'rack-5blue', Colour: 3, Rank: 5, IsJoker: false },
			{ ID: 'rack-6blue', Colour: 3, Rank: 6, IsJoker: false },
			{ ID: 'rack-other', Colour: 4, Rank: 4, IsJoker: false },
			{ ID: 'rack-2', Colour: 2, Rank: 2, IsJoker: false },
			{ ID: 'rack-3', Colour: 2, Rank: 3, IsJoker: false },
			{ ID: 'rack-4', Colour: 2, Rank: 4, IsJoker: false },
			{ ID: 'rack-10', Colour: 1, Rank: 10, IsJoker: false },
			{ ID: 'rack-11', Colour: 1, Rank: 11, IsJoker: false },
			{ ID: 'rack-12', Colour: 1, Rank: 12, IsJoker: false },
			{ ID: 'rack-13', Colour: 1, Rank: 13, IsJoker: false },
			{ ID: 'rack-14', Colour: 4, Rank: 7, IsJoker: false }
		];
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
			tableMelds: hasOpened
				? [
						{
							ID: 'm1',
							Kind: 'run',
							Tiles: [
								{ ID: 'mt-5', Colour: 1, Rank: 5, IsJoker: false },
								{ ID: 'mt-6', Colour: 1, Rank: 6, IsJoker: false },
								{ ID: 'joker-1', Colour: 1, Rank: 7, IsJoker: true }
							],
							JokerReps: { 'joker-1': { ID: 'rep-7red', Colour: 1, Rank: 7, IsJoker: false } },
							OwnerSeat: 0
						}
					]
				: [],
			winner: -1,
			ownRack: rack,
			ownSeat: 0
		};
	}

	function setHasOpenedJoker() {
		onPrivateSnapshot(makeJokerMeld(true));
		replaceTargetMeldId.set(null);
	}

	function setNotOpened() {
		onPrivateSnapshot(makeJokerMeld(false));
		replaceTargetMeldId.set(null);
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
	let target = $derived($replaceTargetMeldId);
	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Joker Replace Demo — Day 38</h1>
	<p class="text-sm text-white/60">
		TableBoard onJokerClicked + Rack replaceSelected(targetMeldId, tileId, new1, new2) →
		OpClientReplaceJoker 9 &#123;targetMeldId, tileId, newMeldTiles[2]&#125;
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setHasOpenedJoker} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set HasOpened + Joker Run 5-6-J</button
		>
		<button onclick={setNotOpened} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set NotOpened</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="joker-info">
		hasOpened:{String(hasOpened)} target:{String(target)}
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
