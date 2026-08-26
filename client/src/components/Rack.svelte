<script lang="ts">
	import Tile from './Tile.svelte';

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
					Mâna ta • {tiles.length} cărți
					<span class="ml-1 rounded bg-amber-400 px-1.5 py-0.5 text-[10px] font-bold text-black"
						>TREBUIE SĂ TRAGI</span
					>
				</div>
				<div class="text-[11px] text-white/60">Click pe carte pentru selectare</div>
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
		{#each tiles as t (t.id)}
			<Tile colour={t.colour} rank={t.rank} isJoker={t.isJoker} />
		{/each}
	</div>

	<div class="mt-3 flex flex-wrap gap-2">
		<button
			class="flex-1 rounded-xl bg-emerald-500 px-4 py-2.5 text-xs font-bold text-black hover:bg-emerald-400 sm:flex-none"
			>▶ TRAGE DIN TALON</button
		>
		<button
			class="flex-1 rounded-xl bg-white/10 px-4 py-2.5 text-xs font-bold text-white/40 sm:flex-none"
			>↩ IA ULTIMA</button
		>
		<button
			class="flex-1 rounded-xl bg-white/10 px-4 py-2.5 text-xs font-bold text-white/40 sm:flex-none"
			>ETALEAZĂ SELECTATE</button
		>
		<button
			class="flex-1 rounded-xl bg-amber-400 px-4 py-2.5 text-xs font-bold text-black hover:bg-amber-300 sm:flex-none"
			>ARUNCĂ CARTEA</button
		>
	</div>
</div>
