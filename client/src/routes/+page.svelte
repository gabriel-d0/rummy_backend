<script lang="ts">
	import { colors } from '$lib/ui/tokens';
</script>

<div class="flex min-h-screen flex-col items-center bg-[#0a2e1a] p-4 sm:p-6">
	<div class="flex w-full max-w-5xl flex-col gap-4">
		<header class="flex items-center justify-between">
			<h1
				class="text-xl font-bold tracking-tight text-white"
				style="font-family: Inter, system-ui, sans-serif"
			>
				Rummy — SvelteKit Day 3
			</h1>
			<span
				class="rounded-full border border-white/10 bg-white/10 px-2.5 py-1 font-mono text-xs text-white/70"
				>Tailwind + Tokens</span
			>
		</header>

		<div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
			<h2 class="text-lg font-semibold text-white">Day 3 — Tailwind & design tokens</h2>
			<p class="mt-2 text-sm leading-relaxed text-white/70">
				Tailwind via <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs"
					>@tailwindcss/vite</code
				>
				+
				<code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs"
					>app.css @import "tailwindcss"</code
				>, tokens
				<code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">felt {colors.felt}</code>
				<code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs">wood {colors.wood}</code>,
				<span style="font-family: Inter, sans-serif">Inter</span> +
				<span style="font-family: 'JetBrains Mono', monospace">JetBrains Mono</span>
			</p>
			<div class="mt-4 flex flex-wrap gap-2">
				<span
					class="h-6 w-16 rounded-full border"
					style="background: {colors.felt}; border-color: {colors.wood}"
				></span>
				<span class="h-6 w-16 rounded-full border bg-white" style="border-color: {colors.tile.red}"
				></span>
				<span
					class="h-6 w-16 rounded-full border bg-white"
					style="border-color: {colors.tile.yellow}"
				></span>
				<span class="h-6 w-16 rounded-full border bg-white" style="border-color: {colors.tile.blue}"
				></span>
				<span
					class="h-6 w-16 rounded-full border bg-white"
					style="border-color: {colors.tile.black}"
				></span>
			</div>
		</div>

		<div class="rounded-2xl border border-white/10 bg-[#0f2a1a]/50 p-6">
			<p class="text-sm text-white">Welcome to SvelteKit — Romanian Tile Rummy will appear here.</p>
			<p class="mt-2 inline-block rounded bg-black/40 px-2 py-1 font-mono text-xs text-white/70">
				Day 3 • Tailwind via @tailwindcss/vite + app.css
			</p>
		</div>
	</div>
</div>
