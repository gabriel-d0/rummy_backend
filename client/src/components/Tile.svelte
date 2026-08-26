<script lang="ts">
  let {
    colour = 1,
    rank = 1,
    isJoker = false,
    size = "rack",
    selected = false,
    draggable = false
  } = $props<{
    colour?: number;
    rank?: number;
    isJoker?: boolean;
    size?: "rack" | "table";
    selected?: boolean;
    draggable?: boolean;
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
      w-[48px] h-[64px]
      {isJoker ? 'border-amber-500 bg-amber-50' : base}
      {selected ? 'ring-2 ring-sky-500 ring-offset-1 scale-[1.06] shadow-lg' : 'shadow-sm hover:shadow-md hover:-translate-y-0.5'}
      {draggable ? 'cursor-grab active:cursor-grabbing' : 'cursor-pointer'}"
    draggable={draggable}
  >
    <div class="flex w-full justify-between px-1 pt-1 text-[9px] font-bold leading-none {isJoker ? 'text-amber-600' : base.split(' ')[0]}">
      <span>{isJoker ? "J" : label}</span>
      <span class="text-[6px] opacity-60">{isJoker ? "J" : colour === 1 ? "♥" : colour === 2 ? "♦" : colour === 3 ? "♣" : "♠"}</span>
    </div>
    <div class="flex flex-1 flex-col items-center justify-center py-0.5">
      {#if isJoker}
        <div class="text-center">
          <div class="text-[16px] font-black leading-none text-amber-600">J</div>
          <div class="text-[6px] font-bold tracking-widest text-amber-700">JOKER</div>
        </div>
      {:else}
        <div class="text-center">
          <div class="text-[18px] font-black leading-none {base.split(' ')[0]}">{label}</div>
          <div class="text-[8px] {base.split(' ')[0]}">●</div>
        </div>
      {/if}
    </div>
    <div class="flex w-full justify-between px-1 pb-1 text-[9px] font-bold leading-none rotate-180 {isJoker ? 'text-amber-600' : base.split(' ')[0]}">
      <span>{isJoker ? "J" : label}</span>
      <span class="text-[6px] opacity-60">{isJoker ? "J" : colour === 1 ? "♥" : colour === 2 ? "♦" : colour === 3 ? "♣" : "♠"}</span>
    </div>
    {#if selected}
      <div class="absolute -right-1 -top-1 grid h-4 w-4 place-items-center rounded-full bg-sky-500 text-[8px] font-bold text-white shadow">✓</div>
    {/if}
  </div>
{/snippet}

{@render tileContent()}
