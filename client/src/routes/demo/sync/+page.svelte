<script lang="ts">
	import Rack from '../../../components/Rack.svelte';
	import TableBoard from '../../../components/TableBoard.svelte';
	import { privateStore, publicStore, onPrivateSnapshot } from '$lib/game/store';
	import type { PrivateSnapshot } from '$lib/game/snapshot';

	// Day 28 — Visual sync demo — PrivateSnapshot Rack only local, PublicSnapshot Table for all

	function makePrivate(seat: number, rackIds: string[], tableIds: string[]): PrivateSnapshot {
		return {
			v: 1,
			gamePhase: 'Playing',
			turnPhase: 'MeldOrDiscard',
			currentSeat: 0,
			players: [
				{ id: 'alice', seat: 0, hasOpened: true, rackCount: seat === 0 ? rackIds.length : 7 },
				{ id: 'bob', seat: 1, hasOpened: true, rackCount: seat === 1 ? rackIds.length : 7 }
			],
			stockCount: 70,
			discardRow: [],
			tableMelds: tableIds.map((id, idx) => ({
				ID: `m${idx + 1}-${id}`,
				Kind: idx === 0 ? 'run' : 'set',
				Tiles: [
					{ ID: `${id}-t1`, Colour: 1, Rank: 5, IsJoker: false },
					{ ID: `${id}-t2`, Colour: 1, Rank: 6, IsJoker: false },
					{ ID: `${id}-t3`, Colour: 1, Rank: 7, IsJoker: false }
				],
				JokerReps: {},
				OwnerSeat: 0
			})),
			winner: -1,
			ownRack: rackIds.map((id) => ({ ID: id, Colour: 1, Rank: 5, IsJoker: false })),
			ownSeat: seat
		};
	}

	function setAlice() {
		const snap = makePrivate(
			0,
			['alice-secret-1', 'alice-secret-2', 'alice-secret-3'],
			['shared-meld-1', 'shared-meld-2']
		);
		onPrivateSnapshot(snap);
	}

	function setBob() {
		const snap = makePrivate(
			1,
			['bob-secret-1', 'bob-secret-2'],
			['shared-meld-1', 'shared-meld-2']
		);
		onPrivateSnapshot(snap);
	}

	function clear() {
		// reset via direct store clear for demo
		import('$lib/game/store').then(({ _resetForTest }) => _resetForTest());
	}

	const rackTiles = $derived.by(() => {
		const priv = $privateStore;
		if (!priv) return [];
		return priv.ownRack.map((t) => ({
			id: t.ID,
			colour: t.Colour,
			rank: t.Rank,
			isJoker: t.IsJoker
		}));
	});

	const tableCount = $derived($publicStore?.tableMelds.length ?? 0);
	const currentSeat = $derived($privateStore?.ownSeat ?? -1);
	const currentRackIds = $derived($privateStore?.ownRack.map((t) => t.ID).join(',') ?? '');
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Sync Demo — Day 28</h1>
	<p class="text-sm text-white/60">PrivateSnapshot Rack only local, PublicSnapshot Table for all</p>
	<div class="mt-3 flex gap-2">
		<button onclick={setAlice} class="rounded bg-emerald-500 px-3 py-1 text-sm font-bold text-black"
			>Set Alice (3 tiles)</button
		>
		<button onclick={setBob} class="rounded bg-sky-500 px-3 py-1 text-sm font-bold text-white"
			>Set Bob (2 tiles)</button
		>
		<button onclick={clear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="current-rack">
		currentSeat:{currentSeat} rack:{currentRackIds} table:{tableCount}
	</div>
	<div class="mt-4 space-y-3">
		<div data-testid="rack-section">
			<Rack tiles={rackTiles} />
		</div>
		<div data-testid="table-section">
			<TableBoard />
		</div>
	</div>
	<div class="mt-3 text-xs text-white/60" data-testid="sync-info">
		Rack only local (alice vs bob different), Table shared (same melds for both)
	</div>
</div>
