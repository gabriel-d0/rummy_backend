<script lang="ts">
  import TopBar from "../components/TopBar.svelte";
  import TableBoard from "../components/TableBoard.svelte";
  import Rack from "../components/Rack.svelte";
  import { colors } from "$lib/ui/tokens";
</script>

<div class="min-h-screen bg-[#0a2e1a] flex flex-col">
  <TopBar />
  <div class="flex-1 flex flex-col lg:flex-row gap-3 p-3 max-w-[1600px] w-full mx-auto">
    <div class="flex-1 flex flex-col gap-4 min-w-0">
      <div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
        <h1 class="text-xl font-bold tracking-tight text-white" style="font-family: Inter, system-ui, sans-serif">
          Rummy — SvelteKit Day 16
        </h1>
        <p class="mt-2 text-sm leading-relaxed text-white/70">
          Visual layout <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">TopBar + TableBoard + Rack + Jurnal</code>
          <span style="font-family: 'JetBrains Mono', monospace"> — 1280×800 & 375×667</span>
        </p>
      </div>

      <TableBoard />
      <Rack />
    </div>
    <aside class="hidden lg:flex w-[300px] shrink-0 flex-col gap-3">
      <div class="rounded-2xl bg-black/20 border border-white/10 p-3">
        <h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
        <p class="mt-2 text-xs text-white/60">Journal placeholder — Day 16</p>
      </div>
    </aside>
  </div>
  <div class="lg:hidden px-3 pb-3">
    <div class="rounded-2xl bg-black/20 border border-white/10 p-3">
      <h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
    </div>
  </div>
</div>
