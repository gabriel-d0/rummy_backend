<script lang="ts">
	import { SvelteSet } from 'svelte/reactivity';
	import Tile from './Tile.svelte';
	import { privateStore } from '$lib/game/store';
	import { sendDiscard, sendDrawStock } from '$lib/game/actions';

	type RackTile = { id: string; colour: number; rank: number; isJoker?: boolean };

	let {
		tiles = [
			{ id: 'r1', colour: 2, rank: 3 },
			{ id: 'r2', colour: 2, rank: 3 },
			{ id: 'r3', colour: 3, rank: 4 },
			{ id: 'r4', colour: 3, rank: 9 },
			{ id: 'r5', colour: 3, rank: 10 },
			{ id: 'r6', colour: 3, rank: 12 },
			{ id: 'r7', colour: 3, rank: 13 },
			{ id: 'r8', colour: 4, rank: 4 },
			{ id: 'r9', colour: 4, rank: 5 },
			{ id: 'r10', colour: 4, rank: 7 },
			{ id: 'r11', colour: 4, rank: 9 }
		]
	} = $props<{ tiles?: RackTile[] }>();

	let selected = new SvelteSet<string>();

	function toggle(id: string) {
		if (selected.has(id)) selected.delete(id);
		else selected.add(id);
	}

	// Day 30-31 — Opening discard + Draw stock: derived from privateStore when available, else fallback to prop tiles
	const isOpeningDiscard = $derived($privateStore?.gamePhase === 'OpeningDiscard');
	const isPlaying = $derived($privateStore?.gamePhase === 'Playing');
	const isMustDraw = $derived($privateStore?.turnPhase === 'MustDraw');
	const isMyTurn = $derived(
		$privateStore ? $privateStore.currentSeat === $privateStore.ownSeat : false
	);
	const rackCount = $derived($privateStore?.ownRack.length ?? tiles.length);
	const displayTiles = $derived.by(() => {
		const priv = $privateStore;
		if (priv && priv.ownRack.length > 0) {
			return priv.ownRack.map((t) => ({
				id: t.ID,
				colour: t.Colour,
				rank: t.Rank,
				isJoker: t.IsJoker
			}));
		}
		return tiles;
	});

	const canDiscard = $derived.by(() => {
		if (!$privateStore) return selected.size === 1;
		if (isOpeningDiscard && isMyTurn) return selected.size === 1 && rackCount === 15;
		// later normal discard will also allow, but for Day 30 only opening
		return false;
	});

	const canDraw = $derived.by(() => {
		if (!$privateStore) return false;
		return isPlaying && isMustDraw && isMyTurn;
	});

	let discarding = $state(false);
	let drawing = $state(false);

	async function discardSelected() {
		if (!canDiscard || selected.size !== 1) return;
		const tileId = [...selected][0];
		if (!tileId) return;
		discarding = true;
		try {
			await sendDiscard(tileId);
			selected.clear();
		} catch (_err) {
			void _err;
		} finally {
			discarding = false;
		}
	}

	async function drawStock() {
		if (!canDraw || drawing) return;
		drawing = true;
		try {
			await sendDrawStock();
		} catch (_err) {
			void _err;
		} finally {
			drawing = false;
		}
	}
</script>

<div class="w-full rounded-2xl border border-white/10 bg-[#1a1a1a] p-3 shadow-xl sm:p-4">
	<div class="mb-3 flex items-center justify-between">
		<div class="flex items-center gap-2">
			<div
				class="grid h-7 w-7 place-items-center rounded-full bg-white text-xs font-bold text-black"
			>
				TU
			</div>
			<div>
				<div class="text-xs font-bold text-white">
					Mâna ta • {rackCount} cărți
					{#if isOpeningDiscard && isMyTurn}
						<span class="ml-1 rounded bg-amber-400 px-1.5 py-0.5 text-[10px] font-bold text-black"
							>ARUNCĂ CARTEA</span
						>
					{:else}
						<span class="ml-1 rounded bg-amber-400 px-1.5 py-0.5 text-[10px] font-bold text-black"
							>TREBUIE SĂ TRAGI</span
						>
					{/if}
				</div>
				<div class="text-[11px] text-white/60">
					Click pe carte pentru selectare • Drag pe set pentru lipire
				</div>
			</div>
		</div>
		<div class="hidden gap-1.5 sm:flex">
			<button class="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-white/80"
				>SORTEAZĂ CULOARE</button
			>
			<button class="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-white/80"
				>SORTEAZĂ NUMĂR</button
			>
		</div>
	</div>

	<div
		class="flex min-h-[110px] flex-wrap content-start justify-center gap-1.5 sm:justify-start sm:gap-2"
	>
		{#each displayTiles as t (t.id)}
			<button
				onclick={() => toggle(t.id)}
				onkeydown={(e) => e.key === 'Enter' && toggle(t.id)}
				class="border-0 bg-transparent p-0"
			>
				<Tile colour={t.colour} rank={t.rank} isJoker={t.isJoker} selected={selected.has(t.id)} />
			</button>
		{/each}
	</div>

	<div class="mt-3 flex flex-wrap gap-2">
		<button
			onclick={drawStock}
			disabled={!canDraw || drawing}
			data-testid="draw-btn"
			class="flex-1 rounded-xl px-4 py-2.5 text-xs font-bold sm:flex-none
				{canDraw
				? 'bg-emerald-500 text-black hover:bg-emerald-400'
				: 'cursor-not-allowed bg-white/10 text-white/40'}">▶ TRAGE DIN TALON</button
		>
		<button
			class="flex-1 cursor-not-allowed rounded-xl bg-white/10 px-4 py-2.5 text-xs font-bold text-white/40 sm:flex-none"
			>ETALEAZĂ SELECTATE</button
		>
		<button
			onclick={discardSelected}
			disabled={!canDiscard || discarding}
			data-testid="discard-btn"
			class="flex-1 rounded-xl px-4 py-2.5 text-xs font-bold sm:flex-none
				{canDiscard
				? 'bg-amber-400 text-black hover:bg-amber-300'
				: 'cursor-not-allowed bg-white/10 text-white/40'}">ARUNCĂ CARTEA</button
		>
		<button class="ml-auto hidden text-xs text-white/50 sm:block" onclick={() => selected.clear()}
			>ANULEAZĂ SELECȚIA</button
		>
	</div>
</div>
