<script lang="ts">
	import Tile from './Tile.svelte';
	import { publicStore } from '$lib/game/store';

	type Meld = {
		id: string;
		kind: 'run' | 'set';
		tiles: { colour: number; rank: number; isJoker?: boolean; id: string }[];
		points?: number;
	};

	let {
		melds = [
			{
				id: 'm1',
				kind: 'run',
				tiles: [
					{ id: 'a1', colour: 1, rank: 5, isJoker: false },
					{ id: 'a2', colour: 1, rank: 6, isJoker: false },
					{ id: 'a3', colour: 1, rank: 7, isJoker: false },
					{ id: 'a4', colour: 1, rank: 8, isJoker: false },
					{ id: 'a5', colour: 1, rank: 9, isJoker: false },
					{ id: 'a6', colour: 1, rank: 10, isJoker: false }
				],
				points: 66
			},
			{
				id: 'm2',
				kind: 'set',
				tiles: [
					{ id: 'b1', colour: 2, rank: 7, isJoker: false },
					{ id: 'b2', colour: 1, rank: 7, isJoker: false },
					{ id: 'b3', colour: 3, rank: 7, isJoker: false }
				],
				points: 53
			},
			{
				id: 'm3',
				kind: 'set',
				tiles: [
					{ id: 'c1', colour: 1, rank: 7, isJoker: false },
					{ id: 'c2', colour: 1, rank: 7, isJoker: false },
					{ id: 'c3', colour: 2, rank: 7, isJoker: true },
					{ id: 'c4', colour: 3, rank: 7, isJoker: false }
				],
				points: 55
			}
		]
	} = $props<{ melds?: Meld[] }>();

	// Day 23 — TableBoard subscribes to publicStore, maps TableMeld -> Meld, no OwnRack leak
	const displayMelds = $derived.by(() => {
		const pub = $publicStore;
		if (pub && Array.isArray(pub.tableMelds)) {
			if (pub.tableMelds.length === 0 && pub.players.length === 0) return melds;
			return pub.tableMelds.map((tm) => ({
				id: tm.ID,
				kind: (tm.Kind === 'set' ? 'set' : 'run') as 'run' | 'set',
				tiles: tm.Tiles.map((t) => ({
					id: t.ID,
					colour: t.Colour,
					rank: t.Rank,
					isJoker: t.IsJoker
				})),
				points: undefined
			}));
		}
		return melds;
	});
</script>

<div
	class="flex min-h-[280px] w-full flex-col gap-2.5 rounded-[18px] border border-[#e8e0c8] bg-[#f5f1e8] p-3 shadow-inner sm:p-4"
>
	<div
		class="flex items-center justify-between text-[11px] font-semibold tracking-widest text-[#8a7a5a]"
	>
			<span>ETALĂRI PE MASĂ • {displayMelds.length} SETURI</span>
		<span class="hidden text-[10px] font-normal sm:inline">Prima etalare min 45 pct</span>
	</div>
	<div class="flex flex-1 flex-col content-start gap-2.5">
		<div class="flex flex-wrap gap-2.5">
			{#each displayMelds.slice(0, 2) as meld (meld.id)}
				<div
					class="flex flex-wrap items-center gap-1 rounded-xl border border-black/5 bg-white/90 px-2 py-1.5 shadow-sm"
				>
					<div class="flex flex-wrap gap-1">
						{#each meld.tiles as t (t.id)}
							<Tile colour={t.colour} rank={t.rank} isJoker={t.isJoker} size="table" />
						{/each}
					</div>
					{#if meld.points}
						<div
							class="ml-1 rounded-full bg-slate-900 px-1.5 py-0.5 text-[10px] font-bold whitespace-nowrap text-white sm:ml-2"
						>
							{meld.points} pct
						</div>
					{/if}
				</div>
			{/each}
		</div>
		<div class="flex flex-wrap gap-2.5">
			{#each displayMelds.slice(2) as meld (meld.id)}
				<div
					class="flex flex-wrap items-center gap-1 rounded-xl border border-black/5 bg-white/90 px-2 py-1.5 shadow-sm"
				>
					<div class="flex flex-wrap gap-1">
						{#each meld.tiles as t (t.id)}
							<Tile colour={t.colour} rank={t.rank} isJoker={t.isJoker} size="table" />
						{/each}
					</div>
					{#if meld.points}
						<div
							class="ml-1 rounded-full bg-slate-900 px-1.5 py-0.5 text-[10px] font-bold whitespace-nowrap text-white sm:ml-2"
						>
							{meld.points} pct
						</div>
					{/if}
				</div>
			{/each}
		</div>
		<div class="min-h-[40px] flex-1"></div>
	</div>
</div>
