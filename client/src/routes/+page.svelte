<script lang="ts">
  import Tile from "../components/Tile.svelte";
  import { colors } from "$lib/ui/tokens";
</script>

<div class="flex min-h-screen flex-col items-center bg-[#0a2e1a] p-4 sm:p-6">
  <div class="flex w-full max-w-5xl flex-col gap-4">
    <header class="flex items-center justify-between">
      <h1 class="text-xl font-bold tracking-tight text-white" style="font-family: Inter, system-ui, sans-serif">
        Rummy — SvelteKit Day 12
      </h1>
      <span class="rounded-full border border-white/10 bg-white/10 px-2.5 py-1 font-mono text-xs text-white/70">Tile modern</span>
    </header>

    <div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
      <h2 class="text-lg font-semibold text-white">Day 12 — Tile modern</h2>
      <p class="mt-2 text-sm leading-relaxed text-white/70">
        Tile <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">48×64 white rounded-lg border + shadow</code>
        <span style="font-family: 'JetBrains Mono', monospace"> — Joly amber-50, selected ring-2 sky-500 scale 1.06 + ✓, draggable</span>
      </p>
      <div class="mt-4 flex flex-wrap gap-3">
        <Tile colour={1} rank={5} />
        <Tile colour={2} rank={7} />
        <Tile colour={3} rank={12} />
        <Tile colour={4} rank={1} />
        <Tile isJoker={true} colour={1} rank={1} />
        <Tile colour={1} rank={5} selected={true} />
        <Tile colour={2} rank={7} draggable={true} />
      </div>
      <p class="mt-3 text-xs font-mono text-white/50">Tile — Day 12 • 48×64 • selected ring-2 • draggable</p>
    </div>

    <div class="rounded-2xl border border-white/10 bg-[#0f2a1a]/50 p-6">
      <p class="text-sm text-white">Welcome to SvelteKit — Romanian Tile Rummy will appear here.</p>
      <p class="mt-2 inline-block rounded bg-black/40 px-2 py-1 font-mono text-xs text-white/70">
        Day 12 • Tile.svelte 48×64 • draggable
      </p>
    </div>
  </div>
</div>
