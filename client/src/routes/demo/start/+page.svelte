<script lang="ts">
	import TopBar from '../../../components/TopBar.svelte';
	import { onPrivateSnapshot, _resetForTest } from '$lib/game/store';
	import { lastSent } from '$lib/game/actions';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	function makeWaiting(seat: number, playerCount: number): PrivateSnapshot {
		const players = Array.from({ length: playerCount }, (_, i) => ({
			id: i === 0 ? 'alice' : i === 1 ? 'bob' : `p${i}`,
			seat: i,
			hasOpened: false,
			rackCount: 14
		}));
		return {
			v: 1,
			gamePhase: 'Waiting',
			turnPhase: '',
			currentSeat: 0,
			players,
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: [],
			ownSeat: seat
		};
	}

	function setWaitingHost2p() {
		onPrivateSnapshot(makeWaiting(0, 2));
	}

	function setWaitingGuest2p() {
		onPrivateSnapshot(makeWaiting(1, 2));
	}

	function setWaiting1p() {
		onPrivateSnapshot(makeWaiting(0, 1));
	}

	function setPlaying() {
		const snap: PrivateSnapshot = {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MustDraw',
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened: false, rackCount: 14 },
				{ id: 'bob', seat: 1, hasOpened: false, rackCount: 14 }
			],
			stockCount: 70,
			discardRow: [],
			tableMelds: [],
			winner: -1,
			ownRack: [],
			ownSeat: 0
		};
		onPrivateSnapshot(snap);
	}

	function clear() {
		_resetForTest();
		lastSent.set(null);
		localStorage.removeItem('rummy_matchId');
	}

	let sentJson = $derived.by(() => {
		const v = $lastSent;
		return v ? JSON.stringify(v) : '';
	});
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<TopBar />
	<h1 class="mt-4 font-bold text-white">Start Demo — Day 29</h1>
	<p class="text-sm text-white/60">
		Start visible only if Waiting ownSeat==0 players&gt;=2 → OpClientStart 1 sendMatchState
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={setWaitingHost2p} class="rounded bg-emerald-600 px-3 py-1 text-sm text-white"
			>Set Waiting Host 2p</button
		>
		<button onclick={setWaitingGuest2p} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Set Waiting Guest 2p</button
		>
		<button onclick={setWaiting1p} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Set Waiting 1p</button
		>
		<button onclick={setPlaying} class="rounded bg-zinc-600 px-3 py-1 text-sm text-white"
			>Set Playing</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="start-info">
		TopBar Start host 2p Waiting
	</div>
	<div class="mt-2 font-mono text-xs break-all text-white/60" data-testid="last-sent">
		{sentJson}
	</div>
	<div class="mt-2 text-xs text-white/40">
		Click START in TopBar when visible to send OpClientStart 1
	</div>
</div>
