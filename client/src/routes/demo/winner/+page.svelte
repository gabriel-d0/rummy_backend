<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import TableBoard from '../../../components/TableBoard.svelte';
	import WinnerOverlay from '../../../components/WinnerOverlay.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeWinner(winner: number): PrivateSnapshot {
		return {
			v: 1,
			gamePhase: 'RoundComplete',
			turnPhase: '',
			currentSeat: winner,
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: 0 },
				{ id: 'bob', seat: 1, hasOpened: true, rackCount: 14 }
			],
			stockCount: 50,
			discardRow: [],
			tableMelds: [
				{
					ID: 'm1',
					Kind: 'run',
					Tiles: [
						{ ID: 'mt-5', Colour: 1, Rank: 5, IsJoker: false },
						{ ID: 'mt-6', Colour: 1, Rank: 6, IsJoker: false },
						{ ID: 'mt-7', Colour: 1, Rank: 7, IsJoker: false }
					],
					JokerReps: {},
					OwnerSeat: 0
				}
			],
			winner,
			ownRack: [],
			ownSeat: 0
		};
	}

	function setWinner0() {
		onPrivateSnapshot(makeWinner(0));
	}

	function setWinner1() {
		onPrivateSnapshot(makeWinner(1));
	}

	function setPlaying() {
		onPrivateSnapshot({
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MeldOrDiscard',
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: 6 },
				{ id: 'bob', seat: 1, hasOpened: true, rackCount: 14 }
			],
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: [],
			ownSeat: 0
		});
	}

	function clear() {
		_resetForTest();
	}

	let isWinner = $derived($privateStore?.gamePhase === 'RoundComplete');
	let winner = $derived($privateStore?.winner ?? -1);
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<WinnerOverlay />
	<h1 class="font-bold text-white">Winner Demo — Day 39</h1>
	<p class="text-sm text-white/60">
		PrivateSnapshot GamePhase==RoundComplete Winner overlay RESTART MASA
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setWinner0} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set Winner 0</button
		>
		<button onclick={setWinner1} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Set Winner 1</button
		>
		<button onclick={setPlaying} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set Playing</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="winner-info">
		isWinner:{String(isWinner)} winner:{winner}
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
