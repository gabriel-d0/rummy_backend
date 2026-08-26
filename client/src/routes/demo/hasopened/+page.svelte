<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import TableBoard from '../../../components/TableBoard.svelte';
	import { privateStore, onPrivateSnapshot, _resetForTest, selectedMeldId } from '$lib/game/store';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeSnap(
		hasOpened: boolean,
		turnPhase: 'MustDraw' | 'MeldOrDiscard',
		withDiscard = true,
		withMeld = false
	): PrivateSnapshot {
		const rack = Array.from({ length: 14 }, (_, i) => ({
			ID: `tile-${i + 1}`,
			Colour: (i % 4) + 1,
			Rank: (i % 13) + 1,
			IsJoker: false
		}));
		const discardRow = withDiscard
			? [
					{
						Tile: { ID: 'disc-1', Colour: 1, Rank: 7, IsJoker: false },
						IsOpeningDiscard: false,
						Index: 0
					}
				]
			: [];
		const tableMelds = withMeld
			? [
					{
						ID: 'meld-1',
						Kind: 'run' as const,
						Tiles: [
							{ ID: 'a1', Colour: 1, Rank: 5, IsJoker: false },
							{ ID: 'a2', Colour: 1, Rank: 6, IsJoker: false },
							{ ID: 'a3', Colour: 1, Rank: 7, IsJoker: false }
						],
						JokerReps: {},
						OwnerSeat: 0
					}
				]
			: [];
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase,
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 70,
			discardRow,
			tableMelds,
			winner: -1,
			ownRack: rack,
			ownSeat: 0
		};
	}

	function setNotOpenedMustDraw() {
		onPrivateSnapshot(makeSnap(false, 'MustDraw', true, false));
	}
	function setOpenedMustDraw() {
		onPrivateSnapshot(makeSnap(true, 'MustDraw', true, false));
	}
	function setNotOpenedMeldOrDiscard() {
		onPrivateSnapshot(makeSnap(false, 'MeldOrDiscard', true, false));
	}
	function setOpenedMeldOrDiscard() {
		onPrivateSnapshot(makeSnap(true, 'MeldOrDiscard', true, true));
	}
	function setOpenedMeldOrDiscardNoMeld() {
		onPrivateSnapshot(makeSnap(true, 'MeldOrDiscard', true, false));
	}
	function selectMeld() {
		selectedMeldId.set('meld-1');
	}
	function clearMeld() {
		selectedMeldId.set(null);
	}
	function clear() {
		_resetForTest();
		selectedMeldId.set(null);
	}

	let hasOpened = $derived.by(() => {
		const p = $privateStore;
		if (!p) return false;
		return p.players.find((pl) => pl.seat === p.ownSeat)?.hasOpened ?? false;
	});
	let turnPhase = $derived($privateStore?.turnPhase ?? '');
	let canDrawPrevious = $derived.by(() => {
		const p = $privateStore;
		if (!p) return false;
		if (p.gamePhase !== 'Playing' || p.turnPhase !== 'MustDraw') return false;
		if (p.currentSeat !== p.ownSeat) return false;
		if (!hasOpened) return false;
		if (!p.discardRow.length) return false;
		return !p.discardRow[p.discardRow.length - 1].IsOpeningDiscard;
	});
	let canPickup = $derived(hasOpened && turnPhase === 'MustDraw');
	let canExtend = $derived(hasOpened && turnPhase === 'MeldOrDiscard');
	let canReplace = $derived(hasOpened && turnPhase === 'MeldOrDiscard');
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">HasOpened Demo — Day 44</h1>
	<p class="text-sm text-white/60">
		Rack HasOpened per PublicPlayer disables Prev/Pickup/Extend/Replace if !HasOpened • 50+ RUN
		needed
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setNotOpenedMustDraw} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Not Opened MustDraw</button
		>
		<button onclick={setOpenedMustDraw} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Opened MustDraw</button
		>
		<button
			onclick={setNotOpenedMeldOrDiscard}
			class="rounded bg-orange-600 px-3 py-1 text-sm text-white">Not Opened MeldOrDiscard</button
		>
		<button
			onclick={setOpenedMeldOrDiscard}
			class="rounded bg-teal-600 px-3 py-1 text-sm text-white">Opened MeldOrDiscard + Meld</button
		>
		<button
			onclick={setOpenedMeldOrDiscardNoMeld}
			class="rounded bg-sky-600 px-3 py-1 text-sm text-white">Opened MeldOrDiscard no Meld</button
		>
		<button onclick={selectMeld} class="rounded bg-violet-600 px-3 py-1 text-sm text-white"
			>Select Meld</button
		>
		<button onclick={clearMeld} class="rounded bg-zinc-600 px-3 py-1 text-sm text-white"
			>Clear Meld</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="hasopened-info">
		hasOpened:{String(hasOpened)} turnPhase:{turnPhase} canDrawPrev:{String(canDrawPrevious)}
		canPickup:{String(canPickup)} canExtend:{String(canExtend)} canReplace:{String(canReplace)}
	</div>
	<div class="mt-4 space-y-4">
		<Rack />
		<TableBoard />
	</div>
	<div class="mt-4 text-xs text-white/40">
		When !HasOpened, IA ULTIMA / RIDICĂ / EXTINDE / ÎNLOCUIEȘTE disabled. When HasOpened &&
		MustDraw, Prev enabled. When HasOpened && MeldOrDiscard + selected tiles + meld selected, Extend
		enabled.
	</div>
</div>
