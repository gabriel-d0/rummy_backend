<script lang="ts">
  import { colors } from "$lib/ui/tokens";
  import { getClient } from "$lib/nakama/client";
  // Day 4: verify getClient can be called
  const hasClient = typeof getClient === "function";
</script>

<div class="flex min-h-screen flex-col items-center bg-[#0a2e1a] p-4 sm:p-6">
  <div class="flex w-full max-w-5xl flex-col gap-4">
    <header class="flex items-center justify-between">
      <h1 class="text-xl font-bold tracking-tight text-white" style="font-family: Inter, system-ui, sans-serif">
        Rummy — SvelteKit Day 4
      </h1>
      <span class="rounded-full border border-white/10 bg-white/10 px-2.5 py-1 font-mono text-xs text-white/70">Nakama JS</span>
    </header>

    <div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
      <h2 class="text-lg font-semibold text-white">Day 4 — Nakama JS client</h2>
      <p class="mt-2 text-sm leading-relaxed text-white/70">
        <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">Client("127.0.0.1","7350","defaultkey")</code>
        <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">authenticateDevice</code> <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">uuid</code>
        <span style="font-family: 'JetBrains Mono', monospace"> token→localStorage</span>
      </p>
      <p class="mt-3 text-xs font-mono text-white/50">getClient() exists: {String(hasClient)} • defaultkey • 127.0.0.1:7350</p>
    </div>

    <div class="rounded-2xl border border-white/10 bg-[#0f2a1a]/50 p-6">
      <p class="text-sm text-white">SvelteKit + Nakama — Day 4 ready for auth.</p>
      <p class="mt-2 inline-block rounded bg-black/40 px-2 py-1 font-mono text-xs text-white/70">
        Day 4 • Nakama JS 2.8 • $env/static/public
      </p>
    </div>
  </div>
</div>
