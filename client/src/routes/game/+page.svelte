<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import TopBar from '../../components/TopBar.svelte';
	import TableBoard from '../../components/TableBoard.svelte';
	import Rack from '../../components/Rack.svelte';
	import WinnerOverlay from '../../components/WinnerOverlay.svelte';
	import { authStore, authenticate } from '$lib/nakama/auth';
	import { privateStore, publicStore, _resetForTest as resetGame } from '$lib/game/store';
	import { matchStore, getMatchId, getStoredMatchId } from '$lib/nakama/match';
	import { reconnect } from '$lib/nakama/reconnect';
	import { colors } from '$lib/ui/tokens';

	const priv = $derived($privateStore);
	const pub = $derived($publicStore);
	const gamePhase = $derived(priv?.gamePhase ?? pub?.gamePhase ?? '');
	const isWaiting = $derived(gamePhase === 'Waiting');
	const isPlaying = $derived(priv?.gamePhase === 'Playing' || pub?.gamePhase === 'Playing');
	const isOpening = $derived(gamePhase === 'OpeningDiscard');
	const players = $derived(priv?.players ?? pub?.players ?? []);
	const currentSeat = $derived(priv?.currentSeat ?? pub?.currentSeat ?? -1);
	const mySeat = $derived(priv?.ownSeat ?? -1);
	const stockCount = $derived(priv?.stockCount ?? pub?.stockCount ?? 0);
	const discardLen = $derived(priv?.discardRow.length ?? pub?.discardRow.length ?? 0);
	const matchId = $derived($matchStore ?? getStoredMatchId() ?? null);
	const isAuthed = $derived($authStore !== null);

	onMount(() => {
		(async () => {
			if (!$authStore) {
				try {
					await authenticate();
				} catch (_err) {
					void _err;
				}
			}
			const last = getStoredMatchId() ?? getMatchId();
			if (last && $authStore && !priv && !pub) {
				try {
					await reconnect();
				} catch (_err) {
					void _err;
				}
			}
		})();
	});

	// If no game after a short delay, redirect to lobby
	$effect(() => {
		if (!isAuthed) return;
		if (!priv && !pub && !matchId) {
			// no match, stay on game page will show message, but we can also redirect after 2s
		}
	});

	function handleLeave() {
		try {
			localStorage.removeItem('rummy_matchId');
			resetGame();
			import('$lib/nakama/match').then(({ clearMatchId }) => clearMatchId());
			goto(resolve('/'));
		} catch (_err) {
			void _err;
		}
	}

	function handleBackToLobby() {
		goto(resolve('/'));
	}
</script>

<div class="flex min-h-screen flex-col bg-[#0a2e1a]">
	<TopBar />
	<WinnerOverlay />
	<div class="mx-auto flex w-full max-w-[1600px] flex-1 flex-col gap-3 p-3 lg:flex-row">
		<div class="flex min-w-0 flex-1 flex-col gap-4">
			<div class="rounded-2xl border border-white/10 p-4" style="background: {colors.felt}">
				{#if !isAuthed}
					<div class="text-sm text-white">Se conectează…</div>
				{:else if !priv && !pub}
					{#if matchId}
						<div class="text-center">
							<div class="text-sm font-bold text-white">
								Se conectează la masa {matchId.slice(0, 8)}…
							</div>
							<p class="mt-2 text-xs text-white/60">Aștepți snapshot de la server (1-2s)…</p>
							<div class="mx-auto mt-3 h-1 w-32 overflow-hidden rounded bg-white/10">
								<div class="h-full w-1/3 animate-pulse bg-emerald-500"></div>
							</div>
							<button onclick={handleBackToLobby} class="mt-3 text-xs text-white/60 underline"
								>Înapoi la Lobby</button
							>
						</div>
					{:else}
						<div class="text-center">
							<div class="text-sm font-bold text-white">Nu ești într-o masă</div>
							<p class="mt-2 text-xs text-white/60">
								Întoarce-te în lobby pentru a crea sau intra într-o cameră.
							</p>
							<button
								onclick={handleBackToLobby}
								class="mt-3 rounded-xl bg-white px-4 py-2 text-sm font-bold text-black"
								>Înapoi la Lobby</button
							>
						</div>
					{/if}
				{:else if isWaiting}
					<div class="flex items-center justify-between">
						<div>
							<div class="text-sm font-bold text-white">
								Lobby — Așteptare jucători ({players.length}/4)
							</div>
							<div class="mt-1 text-xs text-white/60">
								Masa: <code class="font-mono text-white">{matchId}</code>
							</div>
							<div class="mt-2 flex flex-wrap gap-1.5">
								{#each players as p (p.seat)}
									<div class="rounded-full bg-white/10 px-2.5 py-1 text-xs text-white">
										{p.id.slice(0, 6)} • Seat {p.seat}
										{p.seat === mySeat ? '(tu)' : ''}
										{p.hasOpened ? '✓ deschis' : ''}
									</div>
								{/each}
							</div>
						</div>
						<button onclick={handleLeave} class="rounded bg-white/10 px-3 py-1 text-xs text-white"
							>Ieși</button
						>
					</div>
					<div class="mt-3 text-xs text-white/60">
						Host (Seat 0) apasă <span class="font-bold text-emerald-300">START</span> în bara de sus când
						sunt 2-4 jucători.
					</div>
				{:else if isPlaying || isOpening}
					<div class="flex flex-wrap items-center gap-3 text-xs text-white/80">
						<span class="rounded bg-white/10 px-2 py-1"
							>Faza: {gamePhase}
							{isOpening ? '• Aruncă 1 din 15' : ''} / {priv?.turnPhase ??
								pub?.turnPhase ??
								''}</span
						>
						<span class="rounded bg-white/10 px-2 py-1"
							>Tu: Seat {mySeat} {mySeat === currentSeat ? '← rândul tău' : ''}</span
						>
						<span class="rounded bg-white/10 px-2 py-1">Rând: Seat {currentSeat}</span>
						<span class="rounded bg-white/10 px-2 py-1">Stoc: {stockCount}</span>
						<span class="rounded bg-white/10 px-2 py-1">Aruncări: {discardLen}</span>
						<span class="rounded bg-white/10 px-2 py-1"
							>Masă: <code class="font-mono">{matchId?.slice(0, 8) ?? ''}</code></span
						>
						<button onclick={handleLeave} class="ml-auto rounded bg-white/10 px-2 py-1 text-xs"
							>Ieși</button
						>
					</div>
				{:else}
					<div class="text-sm text-white">Joc: {gamePhase} {priv?.turnPhase ?? ''}</div>
				{/if}
			</div>

			{#if isPlaying || isOpening || gamePhase === 'RoundComplete'}
				<div class="flex items-center gap-3 rounded-xl border border-white/10 bg-black/20 p-3">
					<div
						class="grid h-16 w-12 place-items-center rounded-lg bg-[#c9a86a] text-center text-[10px] leading-none font-black text-[#3e2a10] shadow"
					>
						<span>TALON</span><span class="text-[14px] leading-none">{stockCount}</span>
					</div>
					<div class="text-xs text-white/80">
						<div class="font-bold text-white">Talon • {stockCount} piese</div>
						<div class="text-[11px] text-white/60">
							Aruncări: {discardLen} • Tu Seat {mySeat} • Rând Seat {currentSeat}
						</div>
					</div>
				</div>
			{/if}
			<TableBoard />
			<Rack />
			{#if isWaiting}
				<div class="rounded-2xl border border-dashed border-white/20 bg-black/10 p-6 text-center">
					<div class="text-sm font-bold text-white">
						Masa e gata — aștepți START de la host (Seat 0)
					</div>
					<p class="mx-auto mt-2 max-w-md text-xs leading-relaxed text-white/60">
						Când hostul apasă <span class="font-bold text-emerald-300">START</span> vei primi 14 cărți
						(host 15) și talonul va apărea.
					</p>
				</div>
			{/if}
		</div>
		<aside class="hidden w-[300px] shrink-0 flex-col gap-3 lg:flex">
			<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
				<h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
				{#if priv || pub}
					<div class="mt-2 space-y-1 text-xs text-white/60">
						<div>Faza: {gamePhase}</div>
						<div>Tur: {priv?.turnPhase ?? pub?.turnPhase ?? ''}</div>
						<div>Current: Seat {currentSeat} ← {currentSeat === mySeat ? 'tu' : 'altul'}</div>
						<div>
							HasOpened: {priv?.players.find((p) => p.seat === mySeat)?.hasOpened ? 'da' : 'nu'}
						</div>
						<div>Winner: {priv?.winner ?? pub?.winner ?? -1}</div>
					</div>
				{:else}
					<p class="mt-2 text-xs text-white/60">
						Conectat la masă {matchId?.slice(0, 8) ?? ''}. Aștepți snapshot de la server.
					</p>
				{/if}
				{#if matchId}
					<div class="mt-3 rounded bg-white/5 p-2 font-mono text-[10px] break-all text-white/50">
						matchId: {matchId}
					</div>
				{/if}
			</div>
		</aside>
	</div>
</div>
