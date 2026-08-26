<script lang="ts">
  import Tile from "./Tile.svelte";

  type RackTile = { id: string; colour: number; rank: number; isJoker?: boolean };

  let { tiles = [
    { id: "r1", colour: 2, rank: 3 },
    { id: "r2", colour: 2, rank: 3 },
    { id: "r3", colour: 3, rank: 4 },
    { id: "r4", colour: 3, rank: 9 },
    { id: "r5", colour: 3, rank: 10 },
    { id: "r6", colour: 3, rank: 12 },
    { id: "r7", colour: 3, rank: 13 },
    { id: "r8", colour: 4, rank: 4 },
    { id: "r9", colour: 4, rank: 5 },
    { id: "r10", colour: 4, rank: 7 },
    { id: "r11", colour: 4, rank: 9 }
  ] } = $props<{ tiles?: RackTile[] }>();

  let selected = $state(new Set<string>());

  function toggle(id: string) {
    if (selected.has(id)) selected.delete(id);
    else selected.add(id);
    selected = new Set(selected);
  }
</script>

<div class="w-full rounded-2xl bg-[#1a1a1a] border border-white/10 p-3 sm:p-4 shadow-xl">
  <div class="flex items-center justify-between mb-3">
    <div class="flex items-center gap-2">
      <div class="w-7 h-7 rounded-full bg-white grid place-items-center text-xs font-bold text-black">TU</div>
      <div>
        <div class="text-xs font-bold text-white">Mâna ta • {tiles.length} cărți <span class="ml-1 text-[10px] bg-amber-400 text-black px-1.5 py-0.5 rounded font-bold">TREBUIE SĂ TRAGI</span></div>
        <div class="text-[11px] text-white/60">Click pe carte pentru selectare • Drag pe set pentru lipire</div>
      </div>
    </div>
    <div class="hidden sm:flex gap-1.5">
      <button class="text-xs bg-white/10 text-white/80 px-3 py-1 rounded-full font-medium">SORTEAZĂ CULOARE</button>
      <button class="text-xs bg-white/10 text-white/80 px-3 py-1 rounded-full font-medium">SORTEAZĂ NUMĂR</button>
    </div>
  </div>

  <div class="flex flex-wrap gap-1.5 sm:gap-2 justify-center sm:justify-start content-start min-h-[110px]">
    {#each tiles as t (t.id)}
      <button onclick={() => toggle(t.id)} onkeydown={(e) => e.key === "Enter" && toggle(t.id)} class="p-0 bg-transparent border-0">
        <Tile colour={t.colour} rank={t.rank} isJoker={t.isJoker} selected={selected.has(t.id)} />
      </button>
    {/each}
  </div>

  <div class="mt-3 flex flex-wrap gap-2">
    <button class="flex-1 sm:flex-none bg-emerald-500 hover:bg-emerald-400 text-black text-xs font-bold px-4 py-2.5 rounded-xl">▶ TRAGE DIN TALON</button>
    <button class="flex-1 sm:flex-none bg-white/10 text-white/40 text-xs font-bold px-4 py-2.5 rounded-xl cursor-not-allowed">ETALEAZĂ SELECTATE</button>
    <button class="flex-1 sm:flex-none bg-white/10 text-white/40 text-xs font-bold px-4 py-2.5 rounded-xl cursor-not-allowed">ARUNCĂ CARTEA</button>
    <button class="hidden sm:block ml-auto text-xs text-white/50" onclick={() => (selected = new Set())}>ANULEAZĂ SELECȚIA</button>
  </div>
</div>
