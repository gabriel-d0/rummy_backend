<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { flip } from 'svelte/animate';
	import { fade } from 'svelte/transition';
	import TopBar from '../components/TopBar.svelte';
	import { authStore, authenticate } from '$lib/nakama/auth';
	import { privateStore, publicStore, _resetForTest as resetGame } from '$lib/game/store';
	import {
		matchStore,
		getMatchId,
		getStoredMatchId,
		createMatch,
		joinMatch,
		listAvailableMatches,
		type AvailableMatch
	} from '$lib/nakama/match';
	import { goto } from '$app/navigation';
	import { reconnect } from '$lib/nakama/reconnect';
	import { colors } from '$lib/ui/tokens';

	let username = $state('');
	let joinIdInput = $state('');
	let authError = $state('');
	let matchError = $state('');
	let isAuthenticating = $state(true);
	let isCreating = $state(false);
	let isJoining = $state(false);
	let availableMatches = $state<AvailableMatch[]>([]);
	let isListing = $state(false);

	const isAuthed = $derived($authStore !== null);
	const matchId = $derived($matchStore ?? getStoredMatchId() ?? null);
	const priv = $derived($privateStore);
	const pub = $derived($publicStore);
	const gamePhase = $derived(priv?.gamePhase ?? pub?.gamePhase ?? '');
	const currentSeat = $derived(priv?.currentSeat ?? pub?.currentSeat ?? -1);
	const mySeat = $derived(priv?.ownSeat ?? -1);

	async function refreshMatches() {
		if (!$authStore) return;
		if (priv || pub) return; // don't poll when in game
		isListing = true;
		try {
			availableMatches = await listAvailableMatches();
		} catch (_err) {
			void _err;
		} finally {
			isListing = false;
		}
	}

	onMount(() => {
		(async () => {
			// try to authenticate with stored deviceId, fallback to device auth
			try {
				if (!$authStore) {
					await authenticate(username || undefined);
				}
			} catch (err) {
				authError = (err as Error)?.message ?? 'Auth failed';
			} finally {
				isAuthenticating = false;
			}
			// try to reconnect to last match
			const last = getStoredMatchId() ?? getMatchId();
			if (last && $authStore) {
				try {
					await reconnect();
				} catch (_err) {
					void _err;
				}
			}
			// initial poll — no auto interval to avoid hard reload, user clicks reîmprospătează
			// Svelte flip/fade will animate just that list when refreshed
			refreshMatches();
		})();
	});

	async function handleCreate() {
		isCreating = true;
		matchError = '';
		try {
			const id = await createMatch();
			joinIdInput = id;
			await goto(resolve('/game'));
		} catch (err) {
			matchError = (err as Error)?.message ?? 'Create failed';
		} finally {
			isCreating = false;
		}
	}

	async function handleJoin() {
		if (!joinIdInput) {
			matchError = 'Introdu ID masă';
			return;
		}
		isJoining = true;
		matchError = '';
		try {
			await joinMatch(joinIdInput.trim());
			await goto(resolve('/game'));
		} catch (err) {
			matchError = (err as Error)?.message ?? 'Join failed';
		} finally {
			isJoining = false;
		}
	}

	async function handleJoinListed(id: string) {
		joinIdInput = id;
		await handleJoin();
	}

	// Auto-redirect to game when already in a match (priv/pub set)
	$effect(() => {
		if ((priv || pub) && isAuthed) {
			goto(resolve('/game'));
		}
	});

	function handleCopyMatchId() {
		if (matchId) {
			try {
				navigator.clipboard.writeText(matchId);
			} catch (_err) {
				void _err;
			}
		}
	}

	function handleLeave() {
		try {
			localStorage.removeItem('rummy_matchId');
			resetGame();
			// clear matchStore without full reload
			import('$lib/nakama/match').then(({ clearMatchId }) => clearMatchId());
		} catch (_err) {
			void _err;
		}
	}
</script>

<div class="flex min-h-screen flex-col bg-[#0a2e1a]">
	<TopBar />
	<div class="mx-auto flex w-full max-w-[1600px] flex-1 flex-col gap-3 p-3 lg:flex-row">
		<div class="flex min-w-0 flex-1 flex-col gap-4">
			<!-- Lobby only — gameplay moved to /game -->
			<div class="rounded-2xl border border-white/10 p-4" style="background: {colors.felt}">
				{#if isAuthenticating}
					<h1
						class="text-xl font-bold tracking-tight text-white"
						style="font-family: Inter, system-ui, sans-serif"
					>
						Rummy — Remi Etalat
					</h1>
					<div class="mt-2 text-sm text-white">Se conectează la server…</div>
					<div class="mt-1 text-xs text-white/60">Nakama 127.0.0.1:7350 defaultkey</div>
				{:else if !isAuthed}
					<h1
						class="text-xl font-bold tracking-tight text-white"
						style="font-family: Inter, system-ui, sans-serif"
					>
						Rummy — Remi Etalat
					</h1>
					<div class="mt-2 text-sm font-bold text-white">Neautentificat</div>
					<div class="mt-2 flex gap-2">
						<input
							bind:value={username}
							placeholder="Nume (opțional)"
							class="rounded bg-white/10 px-3 py-1.5 text-sm text-white placeholder:text-white/40"
						/>
						<button
							onclick={async () => {
								isAuthenticating = true;
								try {
									await authenticate(username || undefined);
									authError = '';
								} catch (err) {
									authError = (err as Error)?.message ?? 'Auth failed';
								} finally {
									isAuthenticating = false;
								}
							}}
							class="rounded bg-emerald-500 px-4 py-1.5 text-sm font-bold text-white">Intră</button
						>
					</div>
					{#if authError}
						<div class="mt-2 text-xs text-red-300">{authError}</div>
					{/if}
				{:else if !priv && !pub}
					<!-- No game yet, show lobby create/join -->
					<h1
						class="text-xl font-bold tracking-tight text-white"
						style="font-family: Inter, system-ui, sans-serif"
					>
						Rummy — Remi Etalat — Lobby
					</h1>
					<p class="mt-2 text-sm leading-relaxed text-white/70">
						Autentificat ca <span class="font-mono text-white"
							>{$authStore?.user_id?.slice(0, 8) ?? 'anon'}</span
						>
						{#if isAuthed}
							<span class="ml-2 text-emerald-300">● conectat</span>
						{:else}
							<span class="ml-2 text-white/40">○ neconectat</span>
						{/if}
					</p>
					<!-- Primary: one click create -->
					<div class="mt-4">
						<button
							onclick={handleCreate}
							disabled={isCreating}
							class="w-full rounded-2xl bg-emerald-500 px-6 py-4 text-lg font-black text-white shadow-lg hover:bg-emerald-400 disabled:opacity-50"
							>{isCreating ? 'Se creează…' : '＋ Creează cameră nouă'}</button
						>
						<p class="mt-2 text-center text-xs text-white/60">
							2-4 jucători • primul creează, ceilalți văd camera instant la “Camere disponibile”
						</p>
					</div>
					<!-- Available rooms list — the simple path: no typing -->
					<div class="mt-4 rounded-xl border border-white/10 bg-black/20 p-3">
						<div class="flex items-center justify-between">
							<div class="text-xs font-bold tracking-widest text-white/80">CAMERE DISPONIBILE</div>
							<button
								onclick={refreshMatches}
								class="text-[10px] text-white/60 underline hover:text-white">reîmprospătează</button
							>
						</div>
						{#if isListing}
							<div class="mt-2 text-xs text-white/60">Se caută camere…</div>
						{:else if availableMatches.length === 0}
							<div class="mt-2 text-xs text-white/60">
								Nicio cameră deschisă — creează una nouă mai sus.
							</div>
						{:else}
							<div class="mt-2 space-y-2">
								{#each availableMatches as m (m.matchId)}
									<div
										animate:flip={{ duration: 300 }}
										in:fade={{ duration: 200 }}
										class="flex items-center justify-between rounded-lg bg-white/10 px-3 py-2"
									>
										<div>
											<div class="font-mono text-xs font-bold text-white">
												{m.matchId.slice(0, 8)}…
											</div>
											<div class="text-[10px] text-white/60">
												{m.label || 'rummy'} • {m.size}/4 jucători
											</div>
										</div>
										<button
											onclick={() => handleJoinListed(m.matchId)}
											disabled={isJoining}
											class="rounded-full bg-white px-3 py-1.5 text-xs font-bold text-black hover:bg-white/90 disabled:opacity-50"
											>Intră</button
										>
									</div>
								{/each}
							</div>
						{/if}
					</div>
					{#if matchId}
						<div class="mt-3 flex items-center gap-2 text-xs text-white/80">
							<span>Masa curentă:</span>
							<code class="rounded bg-black/30 px-2 py-1 font-mono text-xs text-white"
								>{matchId}</code
							>
							<button
								onclick={handleCopyMatchId}
								class="rounded bg-white/10 px-2 py-1 text-xs text-white">Copiază</button
							>
							<button
								onclick={handleLeave}
								class="rounded bg-red-500/20 px-2 py-1 text-xs text-red-200">Ieși</button
							>
						</div>
					{/if}
					{#if matchError}
						<div class="mt-2 text-xs text-red-300">{matchError}</div>
					{/if}
					<details class="mt-3 text-xs text-white/40">
						<summary class="cursor-pointer underline">Avansat — alătură-te manual cu ID</summary>
						<div class="mt-2 flex gap-1">
							<input
								bind:value={joinIdInput}
								placeholder="ID masă (ex: abc123)"
								class="w-48 rounded-xl bg-white/10 px-3 py-2 text-sm text-white placeholder:text-white/40"
							/>
							<button
								onclick={handleJoin}
								disabled={isJoining}
								class="rounded-xl bg-white px-4 py-2 text-sm font-bold text-black hover:bg-white/90 disabled:opacity-50"
								>Alătură-te</button
							>
						</div>
					</details>
					<div class="mt-3 text-xs text-white/50">
						Demo-uri izolate (fără server): <a href={resolve('/demo/sync')} class="underline"
							>sync</a
						>
						•
						<a href={resolve('/demo/start')} class="underline">start</a>
						• <a href={resolve('/demo/opening')} class="underline">opening</a> •
						<a href={resolve('/demo/draw')} class="underline">draw</a>
						• <a href={resolve('/demo/winner')} class="underline">winner</a>
					</div>
				{:else}
					<!-- In a match (Waiting/Playing) — go to game -->
					<div class="text-center">
						<div class="text-sm font-bold text-white">Ești în masa {matchId?.slice(0, 8)}…</div>
						<p class="mt-2 text-xs text-white/60">
							Jocul e pe pagina de joc — vei fi redirecționat automat.
						</p>
						<a
							href={resolve('/game')}
							class="mt-3 inline-block rounded-xl bg-emerald-500 px-6 py-2 text-sm font-bold text-white"
							>Mergi la joc →</a
						>
						<button
							onclick={handleLeave}
							class="ml-2 rounded bg-white/10 px-3 py-1 text-xs text-white">Ieși</button
						>
					</div>
				{/if}
			</div>
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
						Journal placeholder — conectează-te și creează/alătură-te unei mese. Demo-urile izolate
						rămân la /demo/*.
					</p>
				{/if}
				{#if matchId}
					<div class="mt-3 rounded bg-white/5 p-2 font-mono text-[10px] break-all text-white/50">
						matchId: {matchId}
					</div>
				{/if}
			</div>
			<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
				<h3 class="text-xs font-bold tracking-widest text-white/80">Acțiuni</h3>
				<p class="mt-2 text-[11px] leading-relaxed text-white/60">
					TRAGE DIN TALON → IA ULTIMA → RIDICĂ PENTRU ETALARE (2 + aruncare) → ETALEAZĂ SELECTATE
					(run, 50+ pct la prima) → EXTINDE → ÎNLOCUIEȘTE JOLY (3 + joker) → ARUNCĂ → WIN.
				</p>
			</div>
		</aside>
	</div>
	<div class="px-3 pb-3 lg:hidden">
		<div class="rounded-2xl border border-white/10 bg-black/20 p-3">
			<h3 class="text-xs font-bold tracking-widest text-white/80">JURNAL DE JOC</h3>
			{#if priv || pub}
				<div class="mt-2 text-xs text-white/60">Faza {gamePhase} • Tur Seat {currentSeat}</div>
			{/if}
		</div>
	</div>
</div>
