<script lang="ts">
  import TopBar from "../components/TopBar.svelte";
  import TableBoard from "../components/TableBoard.svelte";
  import Rack from "../components/Rack.svelte";
  import { colors } from "$lib/ui/tokens";
</script>

<div class="min-h-screen bg-[#0a2e1a] flex flex-col">
  <TopBar />
  <div class="flex-1 flex flex-col gap-4 p-3 max-w-[1600px] w-full mx-auto">
    <div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
      <h1 class="text-xl font-bold tracking-tight text-white" style="font-family: Inter, system-ui, sans-serif">
        Rummy — SvelteKit Day 15
      </h1>
      <p class="mt-2 text-sm leading-relaxed text-white/70">
        TopBar <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">REMI ETALAT MASA 1 • 4 JUCĂTORI</code>
        <span style="font-family: 'JetBrains Mono', monospace"> — h-12 bg-black/90</span>
      </p>
    </div>

    <TableBoard />
    <Rack />

    <div class="rounded-2xl border border-white/10 bg-[#0f2a1a]/50 p-6">
      <p class="text-sm text-white">TopBar + TableBoard + Rack — Day 15</p>
      <p class="mt-2 inline-block rounded bg-black/40 px-2 py-1 font-mono text-xs text-white/70">
        Day 15 • TopBar.svelte • TableBoard • Rack
      </p>
    </div>
  </div>
</div>
