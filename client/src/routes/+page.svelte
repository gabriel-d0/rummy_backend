<script lang="ts">
  import Tile from "../components/Tile.svelte";
  import TableBoard from "../components/TableBoard.svelte";
  import { colors } from "$lib/ui/tokens";
</script>

<div class="flex min-h-screen flex-col items-center bg-[#0a2e1a] p-4 sm:p-6">
  <div class="flex w-full max-w-5xl flex-col gap-4">
    <header class="flex items-center justify-between">
      <h1 class="text-xl font-bold tracking-tight text-white" style="font-family: Inter, system-ui, sans-serif">
        Rummy — SvelteKit Day 13
      </h1>
      <span class="rounded-full border border-white/10 bg-white/10 px-2.5 py-1 font-mono text-xs text-white/70">Table board</span>
    </header>

    <div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
      <h2 class="text-lg font-semibold text-white">Day 13 — Table board</h2>
      <p class="mt-2 text-sm leading-relaxed text-white/70">
        TableBoard <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">bg-[#f5f1e8] rounded-[18px]</code> <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">ETALĂRI PE MASĂ • n SETURI</code>
        <span style="font-family: 'JetBrains Mono', monospace"> flex-wrap gap-2.5 no scroll auto-scale</span>
      </p>
    </div>

    <TableBoard />

    <div class="rounded-2xl border border-white/10 bg-[#0f2a1a]/50 p-6">
      <p class="text-sm text-white">Table with 3 melds — 66 pct run, 53 pct set, 55 pct set — no scroll.</p>
      <p class="mt-2 inline-block rounded bg-black/40 px-2 py-1 font-mono text-xs text-white/70">
        Day 13 • TableBoard.svelte • Tile 52×72
      </p>
    </div>
  </div>
</div>
