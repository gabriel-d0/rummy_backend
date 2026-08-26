<script lang="ts">
	import TopBar from '../components/TopBar.svelte';
	import TableBoard from '../components/TableBoard.svelte';
	import Rack from '../components/Rack.svelte';
	import { colors } from '$lib/ui/tokens';
</script>

<div class="flex min-h-screen flex-col bg-[#0a2e1a]">
	<TopBar />
	<div class="mx-auto flex w-full max-w-[1600px] flex-1 flex-col gap-3 p-3 lg:flex-row">
		<div class="flex min-w-0 flex-1 flex-col gap-4">
			<div class="rounded-2xl border border-white/10 p-6" style="background: {colors.felt}">
				<h1
					class="text-xl font-bold tracking-tight text-white"
					style="font-family: Inter, system-ui, sans-serif"
				>
					Rummy — SvelteKit Day 16
				</h1>
				<p class="mt-2 text-sm leading-relaxed text-white/70">
					Visual layout <code class="rounded bg-black/20 px-1 py-0.5 font-mono text-xs"
						>TopBar + TableBoard + Rack + Jurnal</code
					>
					<span style="font-family: 'JetBrains Mono', monospace"> — 1280×800 & 375×667</span>
				</p>
			</div>

			<TableBoard />
			<Rack />
		</div>
		<aside class="hidden w-[300px] shrink-0 flex-col gap-3 lg:flex">
			<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
				<h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
				<p class="mt-2 text-xs text-white/60">Journal placeholder — Day 16</p>
			</div>
		</aside>
	</div>
	<div class="px-3 pb-3 lg:hidden">
		<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
			<h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
		</div>
	</div>
</div>
