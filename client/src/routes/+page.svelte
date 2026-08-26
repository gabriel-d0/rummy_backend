<script lang="ts">
	import TopBar from '../components/TopBar.svelte';
	import TableBoard from '../components/TableBoard.svelte';
	import Rack from '../components/Rack.svelte';
</script>

<div class="flex min-h-screen flex-col bg-[#0a2e1a]">
	<TopBar />
	<div class="mx-auto flex w-full max-w-[1600px] flex-1 flex-col gap-3 p-3 lg:flex-row">
		<div class="flex min-w-0 flex-1 flex-col gap-4">
			<TableBoard />
			<Rack />
		</div>
		<aside class="hidden w-[300px] shrink-0 flex-col gap-3 lg:flex">
			<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
				<h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
				<p class="mt-2 text-xs text-white/60">Jurnal placeholder — Day 16 visual layout.</p>
			</div>
			<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
				<h3 class="text-xs font-bold tracking-widest text-white/80">Acțiuni</h3>
				<p class="mt-2 text-[11px] leading-relaxed text-white/60">
					TRAGE DIN TALON → IA ULTIMA → RIDICĂ PENTRU ETALARE → ETALEAZĂ SELECTATE → EXTINDE →
					ÎNLOCUIEȘTE JOLY → ARUNCĂ → WIN.
				</p>
			</div>
		</aside>
	</div>
	<div class="px-3 pb-3 lg:hidden">
		<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
			<h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
			<p class="mt-2 text-xs text-white/60">Jurnal placeholder — Day 16 visual layout mobile.</p>
		</div>
	</div>
</div>
