<script lang="ts">
	import { privateStore, publicStore } from '$lib/game/store';
	import { sendStart } from '$lib/game/actions';

	let {
		players = 4,
		masa = 1,
		seconds = 0
	} = $props<{ players?: number; masa?: number; seconds?: number }>();

	// Day 29 — Start visible only if Waiting ownSeat==0 players>=2
	const canStart = $derived.by(() => {
		const priv = $privateStore;
		if (priv) return priv.gamePhase === 'Waiting' && priv.ownSeat === 0 && priv.players.length >= 2;
		// fallback when no store (static demo): never show Start unless explicitly in Waiting demo via store
		return false;
	});

	const displayPlayers = $derived($privateStore?.players.length ?? players);
	const displayMasa = $derived(masa);

	// Day 43 — Turn indicator: Current seat + gamePhase/turnPhase + isMyTurn arrow
	const gamePhase = $derived($privateStore?.gamePhase ?? $publicStore?.gamePhase ?? '');
	const turnPhase = $derived($privateStore?.turnPhase ?? $publicStore?.turnPhase ?? '');
	const currentSeat = $derived($privateStore?.currentSeat ?? $publicStore?.currentSeat ?? -1);
	const isMyTurn = $derived(
		$privateStore ? $privateStore.currentSeat === $privateStore.ownSeat : false
	);
	const showTurn = $derived(gamePhase !== '' && currentSeat >= 0);

	let sending = $state(false);

	async function handleStart() {
		if (!canStart || sending) return;
		sending = true;
		try {
			await sendStart();
		} catch (_err) {
			void _err;
		} finally {
			sending = false;
		}
	}
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
			<span class="h-2 w-2 animate-pulse rounded-full bg-emerald-500"></span> MASA {displayMasa} • {displayPlayers}
			JUCĂTORI
		</div>
		{#if showTurn}
			<div
				data-testid="turn-indicator"
				aria-label="Turn indicator: current seat {currentSeat}, phase {gamePhase}/{turnPhase}, {isMyTurn
					? 'your turn'
					: 'others turn'}"
				class="hidden items-center gap-1.5 rounded-full border px-3 py-1.5 text-[11px] font-bold sm:flex {isMyTurn
					? 'border-emerald-400 bg-emerald-500 text-white'
					: 'border-white/10 bg-white/10 text-white/80'}"
			>
				<span data-testid="turn-current">Current: seat-{currentSeat}</span>
				<span class="opacity-60">•</span>
				<span>{gamePhase}/{turnPhase || '—'}</span>
				<span>{isMyTurn ? '← rândul tău' : `← seat-${currentSeat}`}</span>
			</div>
			<div
				data-testid="turn-indicator-mobile"
				class="flex items-center gap-1 rounded-full border px-2 py-1 text-[10px] font-bold sm:hidden {isMyTurn
					? 'border-emerald-400 bg-emerald-500 text-white'
					: 'border-white/10 bg-white/10 text-white/70'}"
			>
				<span>S{currentSeat}</span>
				<span>{isMyTurn ? '← tu' : '←'}</span>
			</div>
		{/if}
		<div class="rounded-full bg-white/10 px-2.5 py-1.5 text-xs">🕒 {seconds}s</div>
		<button class="hidden rounded-full bg-white px-3 py-1.5 text-xs font-bold text-black sm:block"
			>REGULI</button
		>
		{#if canStart}
			<button
				data-testid="start-btn"
				onclick={handleStart}
				disabled={sending}
				class="rounded-full bg-emerald-500 px-4 py-1.5 text-xs font-bold text-white hover:bg-emerald-400 disabled:opacity-50"
				>START</button
			>
		{/if}
		<button class="rounded-full bg-amber-400 px-3 py-1.5 text-xs font-bold text-black"
			>JOC NOU</button
		>
	</div>
</header>
