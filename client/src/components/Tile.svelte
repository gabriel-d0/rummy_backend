<script lang="ts">
  let { colour = 1, rank = 1, isJoker = false, size = "rack", selected = false } = $props<{
    colour?: number;
    rank?: number;
    isJoker?: boolean;
    size?: "rack" | "table";
    selected?: boolean;
  }>();

  const colourMap: Record<number, string> = {
    1: "text-red-600 border-red-500",
    2: "text-amber-500 border-amber-400",
    3: "text-sky-600 border-sky-500",
    4: "text-slate-800 border-slate-700"
  };

  function rankLabel(r: number): string {
    if (r === 1) return "A";
    if (r === 11) return "J";
    if (r === 12) return "Q";
    if (r === 13) return "K";
    return String(r);
  }
</script>

{#snippet tileContent()}
  {@const label = isJoker ? "JOLY" : rankLabel(rank)}
  {@const base = colourMap[colour] ?? colourMap[4]}
  <div
    class="relative flex flex-col items-center justify-between rounded-lg border bg-white select-none transition-all
      {size === 'table' ? 'w-[52px] h-[72px] sm:w-[56px] sm:h-[78px]' : 'w-[64px] h-[90px] sm:w-[72px] sm:h-[100px]'}
      {isJoker ? 'border-amber-500 bg-amber-50' : base}
      {selected ? 'ring-2 ring-sky-500 ring-offset-1 scale-[1.03] shadow-lg' : 'shadow-sm hover:shadow-md hover:-translate-y-0.5'}
      cursor-pointer"
  >
    <div class="flex w-full justify-between px-1.5 pt-1 text-[10px] font-bold leading-none {isJoker ? 'text-amber-600' : base.split(' ')[0]}">
      <span>{isJoker ? "J" : label}</span>
      <span class="text-[7px] opacity-60">{isJoker ? "J" : colour === 1 ? "♥" : colour === 2 ? "♦" : colour === 3 ? "♣" : "♠"}</span>
    </div>
    <div class="flex flex-1 flex-col items-center justify-center py-1">
      {#if isJoker}
        <div class="text-center">
          <div class="text-[18px] font-black leading-none text-amber-600">J</div>
          <div class="text-[7px] font-bold tracking-widest text-amber-700">JOKER</div>
        </div>
      {:else}
        <div class="text-center">
          <div class="text-[22px] font-black leading-none {base.split(' ')[0]}">{label}</div>
          <div class="mt-0.5 text-xs {base.split(' ')[0]}">●</div>
        </div>
      {/if}
    </div>
    <div class="flex w-full justify-between px-1.5 pb-1 text-[10px] font-bold leading-none rotate-180 {isJoker ? 'text-amber-600' : base.split(' ')[0]}">
      <span>{isJoker ? "J" : label}</span>
      <span class="text-[7px] opacity-60">{isJoker ? "J" : colour === 1 ? "♥" : colour === 2 ? "♦" : colour === 3 ? "♣" : "♠"}</span>
    </div>
    {#if selected}
      <div class="absolute -right-1.5 -top-1.5 grid h-5 w-5 place-items-center rounded-full bg-sky-500 text-[10px] font-bold text-white shadow">✓</div>
    {/if}
  </div>
{/snippet}

{@render tileContent()}
