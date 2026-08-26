<script lang="ts">
	let {
		players = 4,
		masa = 1,
		seconds = 0
	} = $props<{ players?: number; masa?: number; seconds?: number }>();
</script>

<header
	class="flex h-12 items-center justify-between border-b border-white/10 bg-black/90 px-3 backdrop-blur sm:px-4"
>
	<div class="flex items-center gap-3">
		<div
			class="grid h-8 w-8 place-items-center rounded-lg bg-amber-400 text-sm font-black text-black"
		>
			R
		</div>
		<div>
			<div
				class="flex items-center gap-2 text-sm leading-none font-black tracking-tight text-white"
			>
				REMI ETALAT <span
					class="rounded bg-amber-400 px-1.5 py-0.5 text-[10px] font-bold text-black">ETALAT</span
				>
			</div>
			<div class="text-[10px] leading-none tracking-widest text-white/60">PREMIUM • ONLINE</div>
		</div>
	</div>
	<div class="flex items-center gap-2">
		<div
			class="hidden items-center gap-1.5 rounded-full border border-white/10 bg-white/10 px-3 py-1.5 text-xs sm:flex"
		>
			<span class="h-2 w-2 animate-pulse rounded-full bg-emerald-500"></span> MASA {masa} • {players}
			JUCĂTORI
		</div>
		<div class="rounded-full bg-white/10 px-2.5 py-1.5 text-xs">🕒 {seconds}s</div>
		<button class="hidden rounded-full bg-white px-3 py-1.5 text-xs font-bold text-black sm:block"
			>REGULI</button
		>
		<button class="rounded-full bg-amber-400 px-3 py-1.5 text-xs font-bold text-black"
			>JOC NOU</button
		>
	</div>
</header>
